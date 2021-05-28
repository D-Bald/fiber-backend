package utils

import (
	"reflect"
)

func MakeQueryFilterFromStruct(s interface{}) (map[string]interface{}, error) {
	// initialise output map
	queryFields := make(map[string]interface{})
	// iterate over struct fields of s
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		// Leave out zero values (important, if the bson flag does not contain "omitempty" option)
		if !v.Field(i).IsZero() {
			switch v.Field(i).Kind() {
			//If the field is a slice and contains just one value, just add the single value not as slice. ONLY USED TO OPTIMIZE MONGO QUERIES
			case reflect.Slice:
				if v.Field(i).Len() == 1 {
					queryFields[string(v.Type().Field(i).Tag.Get("json"))] = v.Field(i).Index(0).Interface()
				} else {
					queryFields[string(v.Type().Field(i).Tag.Get("json"))] = v.Field(i).Interface()
				}
			// If the field is a map add all of its key:value pairs as inline pairs to the current filter
			case reflect.Map:
				iter := v.Field(i).MapRange()
				for iter.Next() {
					queryFields[iter.Key().String()] = iter.Value().Interface()
				}
			// If the field is a struct make a new QueryFilter from it and add it to the current filter
			case reflect.Struct:
				nestedQuery, err := MakeQueryFilterFromStruct(v.Field(i))
				if err != nil {
					return nil, err
				}
				iter := reflect.ValueOf(nestedQuery).MapRange()
				for iter.Next() {
					queryFields[iter.Key().String()] = iter.Value().Interface()
				}
			default:
				queryFields[string(v.Type().Field(i).Tag.Get("json"))] = v.Field(i).Interface()
			}
		}
		// Checks for boolean types because the zero value of this type `false` can be relevant for queries
		if v.Type().Field(i).Type.Kind() == reflect.Bool {
			queryFields[string(v.Type().Field(i).Tag.Get("json"))] = v.Field(i).Interface()
		}
	}
	return queryFields, nil
}
