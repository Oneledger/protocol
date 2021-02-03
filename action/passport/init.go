package passport

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
)

func init() {
	serialize.RegisterConcrete(new(CreateHospitalAdmin), "action_createHospAdmin")
	serialize.RegisterConcrete(new(UploadTestInfo), "action_uploadTestInfo")
	serialize.RegisterConcrete(new(ReadTestInfo), "action_readTestInfo")
}

func EnablePassport(r action.Router) error {
	err := r.AddHandler(action.PASSPORT_HOSP_ADMIN, createHospitalAdminTx{})
	if err != nil {
		return errors.Wrap(err, "createHospitalAdminTx")
	}
	err = r.AddHandler(action.PASSPORT_SCR_ADMIN, createScreenerAdminTx{})
	if err != nil {
		return errors.Wrap(err, "createScreenerAdminTx")
	}
	err = r.AddHandler(action.PASSPORT_UPLOAD_TEST, uploadTestInfoTx{})
	if err != nil {
		return errors.Wrap(err, "uploadTestInfoTx")
	}
	err = r.AddHandler(action.PASSPORT_READ_TEST, readTestInfoTx{})
	if err != nil {
		return errors.Wrap(err, "readTestInfoTx")
	}
	err = r.AddHandler(action.PASSPORT_UPDATE_TEST, updateTestInfoTx{})
	if err != nil {
		return errors.Wrap(err, "updateTestInfoTx")
	}
	return nil
}
