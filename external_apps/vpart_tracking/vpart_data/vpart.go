package vpart_data

type VPart struct {
	VIN           Vin         `json:"vin"`
	PartType      string      `json:"partType"`
	DealerName    string      `json:"dealerName"`
	DealerAddress string      `json:"dealerAddress"`
	StockNum      StockNumber `json:"stockNum"`
	Year          int         `json:"year"`
}

func NewVPart(VIN Vin, partType string, dealerName string, dealerAddress string, stockNum StockNumber, year int) *VPart {
	return &VPart{VIN: VIN, PartType: partType, DealerName: dealerName, DealerAddress: dealerAddress, StockNum: stockNum, Year: year}
}


