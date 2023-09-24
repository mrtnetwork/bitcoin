package formating

import "reflect"

// ToInterfaceSlice converts a slice of any type into a slice of interfaces.
// It takes a slice 'slice' of any type as input and utilizes reflection to create
// a new slice of interfaces where each element in the original slice is converted
// to an interface{} type. The resulting slice of interfaces is returned.
// If the input is not a slice, it raises a panic.
func ToInterfaceSlice(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("Input is not a slice")
	}

	interfaceSlice := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		interfaceSlice[i] = v.Index(i).Interface()
	}

	return interfaceSlice
}

// FlattenList flattens a nested list structure represented by an interface{}.
// It takes an input of type interface{} which can contain nested slices, and it recursively
// flattens the nested structure into a single []interface{}.
// The function uses reflection to inspect the input's type and structure.
// If the input is a slice, it recursively flattens its elements.
// If the input is not a slice, it appends it as an element to the result slice.
// The flattened elements are returned as a []interface{} slice.
func FlattenList(input interface{}) []interface{} {
	var result []interface{}

	// Use reflection to check the type of input
	switch reflect.TypeOf(input).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(input)
		for i := 0; i < slice.Len(); i++ {
			item := slice.Index(i).Interface()
			result = append(result, FlattenList(item)...)
		}
	default:
		result = append(result, input)
	}

	return result
}
