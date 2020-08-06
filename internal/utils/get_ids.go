package utils

import "reflect"

type GetIDFunc func(interface{}) string

func GetIDs(items interface{}, getIDFunc GetIDFunc) []string {
	set := map[string]bool{}
	value := reflect.ValueOf(items)
	var ids []string
	for i := 0; i < value.Len(); i++ {
		id := getIDFunc(value.Index(i).Interface())
		if !set[id] {
			ids = append(ids, id)
			set[id] = true
		}
	}
	return ids
}
