package utils

import (
	"context"
	"sync"
)

type ItemLoadFunc[T any] func(ctx context.Context, id string) (T, error)

const maxConcurrency = 10

type valueOrError[T any] struct {
	Value T
	Err   error
}

// LoadDetails concurrently loads items with a given ids using given load function.
func LoadDetails[T any](ctx context.Context, ids []string, loadFunc ItemLoadFunc[T]) ([]T, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	valuesCh := make(chan *valueOrError[T])
	go func() {
		wg := sync.WaitGroup{}
		defer func() {
			wg.Wait() // wait for all requests to complete before closing resultsCh
			close(valuesCh)
		}()
		sem := make(chan interface{}, maxConcurrency)
		for _, id := range ids {
			wg.Add(1)
			select {
			case sem <- nil:
			case <-ctx.Done():
				return
			}
			go func(id string) {
				value, er := loadFunc(ctx, id)
				valuesCh <- &valueOrError[T]{Value: value, Err: er}
				wg.Done()
				<-sem
			}(id)
		}
	}()

	items := make([]T, 0, len(ids))
	for v := range valuesCh {
		if v.Err != nil {
			// in case of error cancel pending requests and drain resultsCh
			cancel()
			for range valuesCh {
			}
			return nil, v.Err
		}
		items = append(items, v.Value)
	}
	return items, nil
}
