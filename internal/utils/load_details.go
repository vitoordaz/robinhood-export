package utils

import (
	"context"
	"reflect"
	"sync"
)

type ItemLoadFunc func(ctx context.Context, id string) (interface{}, error)

const maxConcurrency = 10

// LoadDetails concurrently loads items with a given ids using given load function.
func LoadDetails(ctx context.Context, ids []string, out interface{}, loadFunc ItemLoadFunc) error {
	type result struct {
		item interface{}
		err  error
	}

	ctx, cancel := context.WithCancel(ctx)
	resultsCh := make(chan *result)
	go func() {
		wg := sync.WaitGroup{}
		defer func() {
			wg.Wait() // wait for all requests to complete before closing resultsCh
			close(resultsCh)
		}()
		sem := make(chan interface{}, maxConcurrency)
		for _, i := range ids {
			wg.Add(1)
			select {
			case sem <- nil:
			case <-ctx.Done():
				return
			}
			go func(id string) {
				item, er := loadFunc(ctx, id)
				resultsCh <- &result{item, er}
				wg.Done()
				<-sem
			}(i)
		}
	}()

	outPtr := reflect.ValueOf(out)
	for result := range resultsCh {
		if result.err != nil {
			// in case of error cancel pending requests and drain resultsCh
			cancel()
			for range resultsCh {
			}
			return result.err
		}
		outPtr.Elem().Set(reflect.Append(outPtr.Elem(), reflect.ValueOf(result.item)))
	}
	return nil
}
