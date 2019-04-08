package serialize


type DataAdapter interface{
	NewDataInstance() Data
	Data() Data
	SetData(interface{}) error
}

type Data interface {
	Primitive() DataAdapter
	SerialTag() string
}