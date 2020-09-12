package merger

import (
	"database/sql"
	"fmt"
	"reflect"
)

// IDChanges represents the changed ids of two slices of a model type
// after a merge has happened, so dependent objects can be updated
// accordingly. So if the ID of an object of the left slice
// changed from id 5 to 20, it will be represented as: {5: 20}.
type IDChanges struct {
	Left  map[int]int
	Right map[int]int
}

// UpdateIDs updates a given ID (named by IDName) on the left and right
// slices of *model.Model according to the given IDChanges.
func UpdateIDs(left interface{}, right interface{}, IDName string, changes IDChanges) {
	for _, mSide := range []mergeSide{leftSide, rightSide} {
		var side interface{}
		var chges map[int]int
		if mSide == leftSide {
			side = left
			chges = changes.Left
		} else {
			side = right
			chges = changes.Right
		}

		switch reflect.TypeOf(side).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(side)
			for i := 0; i < s.Len(); i++ {
				elem := s.Index(i)
				if elem.IsNil() {
					continue
				}

				field := elem.Elem().FieldByName(IDName)
				if !field.IsValid() {
					panic(fmt.Sprintf("Given struct does not contain field %s", IDName))
				}

				switch t := field.Interface().(type) {
				case int:
					if new, ok := chges[int(field.Int())]; ok {
						field.SetInt(int64(new))
					}
				case sql.NullInt32:
					val := field.Field(0)
					if new, ok := chges[int(val.Int())]; ok {
						val.SetInt(int64(new))
					}
				default:
					panic(fmt.Sprintf("Type %T of field %s is not supported!", t, IDName))
				}
			}
		default:
			panic("Only slices are supported!")
		}
	}

}
