package farm_data

type Product struct {
	BatchID         BatchID  `json:"batchId"`
	ItemType        ItemType `json:"itemType"`
	FarmID          FarmID   `json:"farmId"`
	FarmName        string   `json:"farmName"`
	HarvestLocation string   `json:"harvestLocation"`
	HarvestDate     string   `json:"harvestDate"`
	Classification  string   `json:"classification"`
	Quantity        int      `json:"quantity"`
	Description     string   `json:"description"`
}

func newProduct(batchID BatchID, itemType ItemType, farmID FarmID, farmName string, harvestLocation string, harvestDate string, classification string, quantity int, description string) *productBatch {
	return &Product{BatchID: batchID, ItemType: itemType, FarmID: farmID, FarmName: farmName, HarvestLocation: harvestLocation, HarvestDate: harvestDate, Classification: classification, Quantity: quantity, Description: description}
}



