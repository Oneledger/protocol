package tx

import (
	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	pspt "github.com/Oneledger/protocol/action/passport"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (svc *Service) CreateHospitalAdminToken(req client.CreateTokenRequest, reply *client.CreateTokenReply) error {

	message := pspt.CreateHospitalAdmin{
		User:             req.User,
		TokenTypeID:      req.TokenTypeID,
		TokenType:        req.TokenType,
		TokenSubType:     req.TokenSubType,
		OwnerAddress:     req.OwnerAddress,
		SuperUserAddress: req.SuperUserAddress,
		SuperUser:        req.SuperUser,
		CreationTime:     req.CreationTime,
	}

	data, err := message.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{}
	tx := &action.RawTx{
		Type: action.PASSPORT_HOSP_ADMIN,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTokenReply{
		RawTx: packet,
	}

	return nil
}

func (s *Service) CreateScreenerAdminToken(req client.CreateTokenRequest, reply *client.CreateTokenReply) error {

	message := pspt.CreateScreenerAdmin{
		User:             req.User,
		TokenTypeID:      req.TokenTypeID,
		TokenType:        req.TokenType,
		TokenSubType:     req.TokenSubType,
		OwnerAddress:     req.OwnerAddress,
		SuperUserAddress: req.SuperUserAddress,
		SuperUser:        req.SuperUser,
		CreationTime:     req.CreationTime,
	}

	data, err := message.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{}
	tx := &action.RawTx{
		Type: action.PASSPORT_SCR_ADMIN,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTokenReply{
		RawTx: packet,
	}

	return nil
}

func (svc *Service) AddTestInfo(req client.AddTestInfoRequest, reply *client.CreateTxReply) error {
	upload := pspt.UploadTestInfo{
		TestID:       req.TestID,
		Person:       req.Person,
		Test:         req.Test,
		SubTest:      req.SubTest,
		Manufacturer: req.Manufacturer,
		Result:       req.Result,

		TestOrg:  req.TestOrg,
		TestedAt: req.TestedAt,
		TestedBy: req.TestedBy,

		AnalysisOrg: req.AnalysisOrg,
		AnalyzedAt:  req.AnalyzedAt,
		AnalyzedBy:  req.AnalyzedBy,

		Admin:        req.Admin,
		AdminAddress: req.AdminAddress,
		UploadedAt:   req.UploadedAt,
		Notes:        req.Notes,
	}

	data, err := upload.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{}
	tx := &action.RawTx{
		Type: action.PASSPORT_UPLOAD_TEST,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{RawTx: packet}

	return nil
}

func (svc *Service) UpdateTestInfo(req client.UpdateTestInfoRequest, reply *client.CreateTxReply) error {
	upload := pspt.UpdateTestInfo{
		TestID:       req.TestID,
		Person:       req.Person,
		Test:         req.Test,
		Result:       req.Result,

		AnalysisOrg: req.AnalysisOrg,
		AnalyzedAt:  req.AnalyzedAt,
		AnalyzedBy:  req.AnalyzedBy,

		Admin:        req.Admin,
		AdminAddress: req.AdminAddress,
		Notes:        req.Notes,
	}

	data, err := upload.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{}
	tx := &action.RawTx{
		Type: action.PASSPORT_UPDATE_TEST,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{RawTx: packet}

	return nil
}

func (svc *Service) ReadTestInfo(req client.ReadTestInfoRequest, reply *client.CreateTxReply) error {
	read := pspt.ReadTestInfo{
		Org:          req.Org,
		Admin:        req.Admin,
		AdminAddress: req.AdminAddress,
		Person:       req.Person,
		Address:      req.Address,
		Test:         req.Test,
		ReadAt:       req.ReadAt,
	}

	data, err := read.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{}
	tx := &action.RawTx{
		Type: action.PASSPORT_READ_TEST,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.CreateTxReply{RawTx: packet}

	return nil
}
