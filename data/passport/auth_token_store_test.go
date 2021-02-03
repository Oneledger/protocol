package passport

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	tokStore *AuthTokenStore
	userIds  []UserID
	tokIds   []TokenID
	addrs    []keys.Address
	tokens   []*AuthToken
)

func init() {
	for i := 1; i <= 5; i++ {
		userIds = append(userIds, UserID(fmt.Sprintf("person%v", i)))
		tokIds = append(tokIds, TokenID(fmt.Sprintf("tokenID%v", i)))
		pub, _, _ := keys.NewKeyPairFromTendermint()
		h, _ := pub.GetHandler()
		addrs = append(addrs, h.Address())
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)
	tokens = append(tokens, NewAuthToken(tokIds[0], TokenSuperAdmin, ScreenerInvalid, TypeIDSuperAdmin, RoleSuperAdmin, PermitSuper, userIds[0], addrs[0], createdAt))
	tokens = append(tokens, NewAuthToken(tokIds[1], TokenHospital, ScreenerGeneral, TokenTypeID("SunnyBrook Hospital"), RoleHospitalAdmin, PermitUpload, userIds[1], addrs[1], createdAt))
	tokens = append(tokens, NewAuthToken(tokIds[2], TokenHospital, ScreenerGeneral, TokenTypeID("LakeRidge Hospital"), RoleHospitalAdmin, PermitUpload, userIds[2], addrs[2], createdAt))
	tokens = append(tokens, NewAuthToken(tokIds[3], TokenScreener, ScreenerBorderService, TokenTypeID("Windsor Port of Entry"), RoleScreenerAdmin, PermitScan, userIds[3], addrs[3], createdAt))
	tokens = append(tokens, NewAuthToken(tokIds[4], TokenScreener, ScreenerBorderService, TokenTypeID("Windsor Port of Entry"), RoleScreenerAdmin, PermitScan, userIds[4], addrs[4], createdAt))
}

func setupTokenStore() {
	fmt.Println("####### Testing auth token store #######")
	memDb = db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("cs", memDb))
	tokStore = NewAuthTokenStore("token", cs)
}

func TestNewAuthTokenStore(t *testing.T) {
	setupTokenStore()
	assert.NotNil(t, tokStore)
}

func TestAuthTokenStore_CreateAuthToken(t *testing.T) {
	setupTokenStore()

	// create tokens
	for _, token := range tokens {
		err := tokStore.CreateAuthToken(token)
		assert.Nil(t, err)
		err = tokStore.CreateAuthToken(token)
		assert.Equal(t, ErrDuplicateIdentifier, err)
	}
	tokStore.State.Commit()

	// verify
	for _, token := range tokens {
		actual, err := tokStore.GetAuthToken(token.TokTypeId, token.OwnerId)
		assert.Nil(t, err)
		assert.EqualValues(t, token, actual)
	}
}

func TestAuthTokenStore_RemoveAuthToken(t *testing.T) {
	setupTokenStore()

	// create tokens
	for _, token := range tokens {
		err := tokStore.CreateAuthToken(token)
		assert.Nil(t, err)
	}
	tokStore.State.Commit()

	// remove tokens
	for _, token := range tokens {
		err := tokStore.RemoveAuthToken(token.TokTypeId, "fakeOwner")
		assert.Equal(t, ErrAuthTokenNotFound, err)
		err = tokStore.RemoveAuthToken(token.TokTypeId, token.OwnerId)
		assert.Nil(t, err)
	}
	tokStore.State.Commit()

	// verify
	for _, token := range tokens {
		tok, err := tokStore.GetAuthToken(token.TokTypeId, token.OwnerId)
		assert.Equal(t, ErrAuthTokenNotFound, err)
		assert.Nil(t, tok)
	}
}

func TestInfoStore_IterateOrgTokens(t *testing.T) {
	setupTokenStore()

	// create tokens
	for _, token := range tokens {
		err := tokStore.CreateAuthToken(token)
		assert.Nil(t, err)
	}
	tokStore.State.Commit()
	index := 0

	// iterates superadmin
	tokStore.IterateOrgTokens(TypeIDSuperAdmin, func(token *AuthToken) bool {
		assert.EqualValues(t, tokens[index], token)
		index++
		return false
	})
	assert.Equal(t, 1, index)

	// iterates SunnyBrook
	tokStore.IterateOrgTokens(TokenTypeID("SunnyBrook Hospital"), func(token *AuthToken) bool {
		assert.EqualValues(t, tokens[index], token)
		index++
		return false
	})
	assert.Equal(t, 2, index)

	// iterates LakeRidge
	tokStore.IterateOrgTokens(TokenTypeID("LakeRidge Hospital"), func(token *AuthToken) bool {
		assert.EqualValues(t, tokens[index], token)
		index++
		return false
	})
	assert.Equal(t, 3, index)

	// iterates border service
	tokStore.IterateOrgTokens(TokenTypeID("Windsor Port of Entry"), func(token *AuthToken) bool {
		assert.EqualValues(t, tokens[index], token)
		index++
		return false
	})
	assert.Equal(t, 5, index)
}

func TestAuthTokenStore_HasPermission(t *testing.T) {
	setupTokenStore()

	// create tokens
	for _, token := range tokens {
		err := tokStore.CreateAuthToken(token)
		assert.Nil(t, err)
	}
	tokStore.State.Commit()

	// verify super admin permission
	permitted, err := tokStore.HasPermission(tokens[0].TokTypeId, tokens[0].OwnerId, PermitSuper)
	assert.Nil(t, err)
	assert.True(t, permitted)
	token, err := tokStore.GetAuthToken(tokens[0].TokTypeId, tokens[0].OwnerId)
	assert.Nil(t, err)
	assert.True(t, token.TokPermit.HasPermission(PermitUpload))
	assert.True(t, token.TokPermit.HasPermission(PermitScan))
	assert.True(t, token.TokPermit.HasPermission(PermitQueryTest))

	// verify hospital admin permission
	permitted, err = tokStore.HasPermission(tokens[1].TokTypeId, tokens[1].OwnerId, PermitUpload)
	assert.Nil(t, err)
	assert.True(t, permitted)
	token, err = tokStore.GetAuthToken(tokens[1].TokTypeId, tokens[1].OwnerId)
	assert.Nil(t, err)
	assert.False(t, token.TokPermit.HasPermission(PermitSuper))
	assert.False(t, token.TokPermit.HasPermission(PermitScan))
	assert.False(t, token.TokPermit.HasPermission(PermitQueryTest))

	// verify screener admin permission
	permitted, err = tokStore.HasPermission(tokens[3].TokTypeId, tokens[3].OwnerId, PermitScan)
	assert.Nil(t, err)
	assert.True(t, permitted)
	token, err = tokStore.GetAuthToken(tokens[3].TokTypeId, tokens[3].OwnerId)
	assert.Nil(t, err)
	assert.False(t, token.TokPermit.HasPermission(PermitSuper))
	assert.False(t, token.TokPermit.HasPermission(PermitUpload))
	assert.False(t, token.TokPermit.HasPermission(PermitQueryTest))
}
