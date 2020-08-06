package utils

import (
	"context"
	"reflect"
)

type ListLoadFunc func(ctx context.Context, cursor string) (interface{}, string, error)

// LoadList iterates over results of a given load function and appends them to a given list.
func LoadList(ctx context.Context, list interface{}, loadFunc ListLoadFunc) (err error) {
	valuePtr := reflect.ValueOf(list)
	var items interface{}
	var cursor string
	for {
		if items, cursor, err = loadFunc(ctx, cursor); err != nil {
			return
		}
		valuePtr.Elem().Set(reflect.AppendSlice(valuePtr.Elem(), reflect.ValueOf(items)))
		if cursor == "" {
			break
		}
	}
	return
}
