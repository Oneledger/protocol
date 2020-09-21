package vpart_data

type VPart struct {
	VIN           Vin         `json:"vin"`
	PartType      string      `json:"partType"`
	DealerName    string      `json:"dealerName"`
	StockNum      StockNumber `json:"stockNum"`
	DealerAddress string      `json:"dealerAddress"`
	Year          int         `json:"year"`
}

func NewVPart(VIN Vin, partType string, dealerName string, stockNum StockNumber, dealerAddress string, year int) *VPart {
	return &VPart{VIN: VIN, PartType: partType, DealerName: dealerName, StockNum: stockNum, DealerAddress: dealerAddress, Year: year}
}


