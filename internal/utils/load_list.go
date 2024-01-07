package utils

import "context"

type ListLoadFunc[T any] func(ctx context.Context, cursor string) ([]*T, string, error)

// LoadList iterates over results of a given load function and appends them to a given list.
func LoadList[T any](ctx context.Context, loadFunc ListLoadFunc[T]) ([]*T, error) {
	var result []*T
	var (
		items  interface{}
		cursor string
		err    error
	)
	for {
		if items, cursor, err = loadFunc(ctx, cursor); err != nil {
			return nil, err
		}
		for _, item := range items.([]*T) {
			result = append(result, item)
		}
		// valuePtr.Elem().Set(reflect.AppendSlice(valuePtr.Elem(), reflect.ValueOf(items)))
		if cursor == "" {
			break
		}
	}
	return result, nil
}
