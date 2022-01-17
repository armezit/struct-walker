package walker

import (
	"reflect"
	"strconv"
	"unicode"
)

// Visitor is a function that will be called on each visited node.
// value is a non-map value, corresponding to the above path.
// branch is a slice containing consecutive interfaces used to arrive at the given value.
// path is a slice containing consecutive keys used to arrive at the given value.
type Visitor func(value reflect.Value, branch []interface{}, path []string, field *reflect.StructField)

// Walk walks the given struct interface recursively and calls the visitor at every node
func Walk(s interface{}, visitor Visitor) {
	walk(s, []interface{}{}, []string{}, visitor)
}

func getFieldKey(f reflect.StructField) string {
	jsonTag := f.Tag.Get("json")
	if jsonTag == "-" || jsonTag == "" {
		// lowercase first char
		for i, v := range f.Name {
			return string(unicode.ToLower(v)) + f.Name[i+1:]
		}
	}
	return jsonTag
}

func walk(node interface{}, branch []interface{}, path []string, visitor Visitor) {
	nodeType := reflect.TypeOf(node)
	nodeValue := reflect.ValueOf(node)
	branch = append(branch, node)

	switch nodeValue.Kind() {

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the node we have to call Elem()
		// At the same time this unwraps the pointer, so we don't end up in
		// an infinite recursion
		parentValue := nodeValue.Elem()
		// Check if the pointer is nil
		if !parentValue.IsValid() {
			return
		}
		// Unwrap the newly created pointer
		walk(parentValue, branch, path, visitor)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		parentValue := nodeValue.Elem()
		walk(parentValue, branch, path, visitor)

	case reflect.Struct:
		for i := 0; i < nodeValue.NumField(); i += 1 {
			child := nodeValue.Field(i)
			field := nodeType.Field(i)

			childPath := append(path, getFieldKey(field))

			visitor(child, branch, childPath, &field)
			walk(child.Interface(), branch, childPath, visitor)
		}

	case reflect.Slice:
		for i := 0; i < nodeValue.Len(); i += 1 {
			child := nodeValue.Index(i)
			childPath := append(path, strconv.Itoa(i))

			visitor(child, branch, childPath, nil)
			walk(child.Interface(), branch, childPath, visitor)
		}

	case reflect.Map:
		for _, key := range nodeValue.MapKeys() {
			child := nodeValue.MapIndex(key)
			childPath := append(path, key.String())

			visitor(child, branch, childPath, nil)
			walk(child.Interface(), branch, childPath, visitor)
		}
	}
}
