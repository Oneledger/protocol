package vpart_rpc

import "github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_data"

type GetVehiclePartRequest struct {
	VIN      vpart_data.Vin `json:"vin"`
	PartType string         `json:"partType"`
}

type GetVehiclePartReply struct {
	VehiclePart vpart_data.VPart `json:"vehiclePart"`
	Height      int64            `json:"height"`
}
