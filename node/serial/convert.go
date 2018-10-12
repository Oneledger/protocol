package serial

// ConvertMap takes a structure and return a map of its elements
func ConvertMap(structure interface{}) map[string]interface{} {
	var result map[string]interface{}

	children := GetChildren(structure)
	result = make(map[string]interface{}, len(children))

	for _, child := range children {
		result[child.Name] = child.Value
	}
	return result
}
