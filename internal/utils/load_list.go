package utils

import (
	"context"
	"reflect"
)

type ListLoadFunc func(ctx context.Context, cursor string) (interface{}, string, error)

// LoadList iterates over results of a given load function and appends them to a given list.
func LoadList(ctx context.Context, list interface{}, loadFunc ListLoadFunc) error {
	valuePtr := reflect.ValueOf(list)
	var (
		items  interface{}
		cursor string
		err    error
	)
	for {
		if items, cursor, err = loadFunc(ctx, cursor); err != nil {
			return err
		}
		valuePtr.Elem().Set(reflect.AppendSlice(valuePtr.Elem(), reflect.ValueOf(items)))
		if cursor == "" {
			break
		}
	}
	return nil
}
