package utils

import "context"

type ListLoadFunc[T any] func(ctx context.Context, cursor string) ([]T, string, error)

// LoadList iterates over results of a given load function and appends them to a given list.
func LoadList[T any](ctx context.Context, loadFunc ListLoadFunc[T]) ([]T, error) {
	var (
		result []T
		items  []T
		cursor string
		err    error
	)
	for {
		items, cursor, err = loadFunc(ctx, cursor)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
		if cursor == "" {
			break
		}
	}
	return result, nil
}
