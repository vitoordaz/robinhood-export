package main

import (
	"context"
	"io"
	"os"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
)

func doOptionsOrders(args arguments) {
	if !args.verbose {
		logVerbose.SetOutput(io.Discard) // disable verbose logging
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := robinhood.New()

	var (
		token *robinhood.ResponseToken
		err   error
	)
	if args.pathToTokenFile == "" {
		token, err = getAuthToken(ctx, client, args.username)
	} else {
		token, err = getAuthTokenFromFile(args.pathToTokenFile)
	}
	if err != nil {
		logError.Fatalln(err)
	}

	var outputFunc func(f *os.File) error

	logVerbose.Println("loading options orders")
	switch args.format {
	case "json":
		orders, err := loadItems(
			ctx,
			func(ctx context.Context, cursor string) (*robinhood.ResponseList[map[string]any], error) {
				return client.GetItems(ctx, robinhood.EndpointOptionsOrders, token, cursor)
			},
		)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d option orders", len(orders))
		outputFunc = func(f *os.File) error { return outputItemsJSON(f, orders) }
	default:
		logError.Fatalln("unsupported output format " + args.format)
	}

	outputData(args.output, outputFunc)
}
