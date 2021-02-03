package passport

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type AuthTokenStore struct {
	State  *storage.State
	szlr   serialize.Serializer
	prefix []byte
}

func NewAuthTokenStore(prefix string, state *storage.State) *AuthTokenStore {
	return &AuthTokenStore{
		State:  state,
		prefix: storage.Prefix(prefix),
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (ats *AuthTokenStore) WithState(state *storage.State) *AuthTokenStore {
	ats.State = state
	return ats
}

func (ats *AuthTokenStore) Exists(org TokenTypeID, id UserID) bool {
	key := ats.getAuthTokenKey(org, id)
	return ats.State.Exists(key)
}

func (ats *AuthTokenStore) CreateAuthToken(token *AuthToken) (err error) {
	key := ats.getAuthTokenKey(token.TokTypeId, token.OwnerId)
	if ats.State.Exists(key) {
		err = ErrDuplicateIdentifier
		return
	}
	err = ats.set(key, token)
	if err != nil {
		return
	}
	return
}

func (ats *AuthTokenStore) RemoveAuthToken(org TokenTypeID, id UserID) error {
	// get token
	key := ats.getAuthTokenKey(org, id)
	if !ats.State.Exists(key) {
		return ErrAuthTokenNotFound
	}

	// remove token
	ok, err := ats.State.Delete(key)
	if err != nil || !ok {
		return ErrAuthTokenRemove
	}
	return err
}

func (ats *AuthTokenStore) GetAuthToken(org TokenTypeID, id UserID) (token *AuthToken, err error) {
	key := ats.getAuthTokenKey(org, id)
	token, err = ats.get(key)
	return
}

func (ats *AuthTokenStore) HasPermission(org TokenTypeID, id UserID, require TokenPermission) (permited bool, err error) {
	token, err := ats.GetAuthToken(org, id)
	if err != nil {
		return
	}
	permited = token.TokPermit.HasPermission(require)
	return
}

func (ats *AuthTokenStore) IterateOrgTokens(org TokenTypeID, fn func(token *AuthToken) bool) (stopped bool) {
	prefix := append(ats.prefix, (org.String() + storage.DB_PREFIX)...)
	return ats.State.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			token := &AuthToken{}
			err := ats.szlr.Deserialize(value, token)
			if err != nil {
				logger.Error("failed to deserialize auth token")
				return false
			}
			return fn(token)
		},
	)
}

func (ats *AuthTokenStore) IterateAll(fn func(token *AuthToken) bool) (stopped bool) {
	return ats.State.IterateRange(
		ats.prefix,
		storage.Rangefix(string(ats.prefix)),
		true,
		func(key, value []byte) bool {
			token := &AuthToken{}
			err := ats.szlr.Deserialize(value, token)
			if err != nil {
				logger.Error("failed to deserialize auth token")
				return false
			}
			return fn(token)
		},
	)
}

func (ats *AuthTokenStore) GetWithAddress(address keys.Address) (authToken *AuthToken) {
	authToken = &AuthToken{}
	ats.IterateAll(func(token *AuthToken) bool {
		if token.OwnerAddress.Equal(address) {
			authToken = token
		}
		return false
	})
	return authToken
}

//-----------------------------helper functions
// Get auth token key
func (ats *AuthTokenStore) getAuthTokenKey(org TokenTypeID, id UserID) []byte {
	key := string(ats.prefix) + org.String() + storage.DB_PREFIX + id.String()
	return storage.StoreKey(key)
}

// Set object
func (ats *AuthTokenStore) set(key storage.StoreKey, token *AuthToken) error {
	dat, err := ats.szlr.Serialize(token)
	if err != nil {
		return err
	}
	err = ats.State.Set(key, dat)
	return err
}

// Get auth token
func (ats *AuthTokenStore) get(key storage.StoreKey) (token *AuthToken, err error) {
	dat, err := ats.State.Get(key)
	if err != nil {
		return
	}
	if len(dat) == 0 {
		err = ErrAuthTokenNotFound
		return
	}
	token = &AuthToken{}
	err = ats.szlr.Deserialize(dat, token)
	return
}
