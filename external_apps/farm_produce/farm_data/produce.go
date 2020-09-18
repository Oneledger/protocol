package farm_data

type Produce struct {
	BatchID         BatchID `json:"batchId"`
	ItemType        string  `json:"itemType"`
	FarmID          FarmID  `json:"farmId"`
	FarmName        string  `json:"farmName"`
	HarvestLocation string  `json:"harvestLocation"`
	HarvestDate     int64   `json:"harvestDate"`
	Classification  string  `json:"classification"`
	Quantity        int     `json:"quantity"`
	Description     string  `json:"description"`
}

func NewProduce(batchID BatchID, itemType string, farmID FarmID, farmName string, harvestLocation string, harvestDate int64, classification string, quantity int, description string) *Produce {
	return &Produce{BatchID: batchID, ItemType: itemType, FarmID: farmID, FarmName: farmName, HarvestLocation: harvestLocation, HarvestDate: harvestDate, Classification: classification, Quantity: quantity, Description: description}
}
