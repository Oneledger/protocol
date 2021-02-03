package passport

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type AuthToken struct {
	TokId        TokenID         `json:"tokId"`
	TokType      TokenType       `json:"tokType"`
	TokSubType   TokenSubType    `json:"tokSubType"`
	TokTypeId    TokenTypeID     `json:"tokTypeId"`
	TokRole      TokenRole       `json:"tokRole"`
	TokPermit    TokenPermission `json:"tokPermit"`
	OwnerId      UserID          `json:"ownerId"`
	OwnerAddress keys.Address    `json:"ownerAddress"`
	CreatedAt    string          `json:"createdAt"`
}

func NewAuthToken(id TokenID, tokenType TokenType, subType TokenSubType, typeId TokenTypeID,
	role TokenRole, permit TokenPermission, ownerId UserID, addr keys.Address, createdAt string) *AuthToken {
	return &AuthToken{
		TokId:        id,
		TokType:      tokenType,
		TokSubType:   subType,
		TokTypeId:    typeId,
		TokRole:      role,
		TokPermit:    permit,
		OwnerId:      ownerId,
		OwnerAddress: addr,
		CreatedAt:    createdAt,
	}
}

func (token *AuthToken) Err() (err error) {
	if err = token.TokId.Err(); err != nil {
		return
	}
	if err = token.TokTypeId.Err(); err != nil {
		return
	}
	if err = token.OwnerAddress.Err(); err != nil {
		return
	}
	return
}

func (token *AuthToken) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(token)
	if err != nil {
		logger.Error("auth token not serializable", err)
		return []byte{}
	}
	return value
}

func (token *AuthToken) FromBytes(msg []byte) (*AuthToken, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, token)
	if err != nil {
		logger.Error("failed to deserialize auth token from bytes", err)
		return nil, err
	}
	return token, nil
}
