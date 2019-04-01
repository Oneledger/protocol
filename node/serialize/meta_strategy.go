package serialize

// metastrategy has some common methods which is used to compose
// encoding strategies
type metaStrategy struct {

}

// a function type which resembles json.Unmarshal
type deserializeHandler func([]byte, interface{}) error

//meta function to manage data adapter stuff during deserialization
func (m *metaStrategy) wrapDataAdapter(src []byte, dest interface{},
	des deserializeHandler,
	apr DataAdapter) error {

	// create new data instance of obj
	ad := apr.NewDataInstance()

	// deserialize
	err := des(src, ad)
	if err != nil {
		return err
	}

	// set data to primitive
	err = apr.SetData(ad)
	if err != nil {
		return err
	}

	return nil
}

func (m *metaStrategy) serializeString(obj interface{}) {
	return
}
