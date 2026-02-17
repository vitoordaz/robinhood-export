package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
	"github.com/vitoordaz/robinhood-export/internal/utils"
)

func doOrders(args arguments) {
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

	logVerbose.Println("loading orders")
	switch args.format {
	case "json":
		orders, err := loadItems(
			ctx,
			func(ctx context.Context, cursor string) (*robinhood.ResponseList[map[string]any], error) {
				return client.GetItems(ctx, robinhood.EndpointOrders, token, cursor)
			},
		)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d orders\n", len(orders))

		outputFunc = func(f *os.File) error { return outputItemsJSON(f, orders) }
	case "csv":
		orders, err := loadItems(
			ctx,
			func(ctx context.Context, cursor string) (*robinhood.ResponseList[*robinhood.Order], error) {
				return client.GetOrders(ctx, token, cursor)
			},
		)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d orders\n", len(orders))

		logVerbose.Println("loading instruments")
		instruments, err := loadInstruments(ctx, client, getOrdersInstrumentIds(orders))
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d instruments\n", len(instruments))

		logVerbose.Println("loading markets")
		markets, err := loadMarkets(ctx, client, getInstrumentsMarketIds(instruments))
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d markets\n", len(markets))

		outputFunc = func(f *os.File) error { return outputOrdersCSV(f, orders, instruments, markets) }
	default:
		logError.Fatalln("unsupported output format " + args.format)
	}

	outputData(args.output, outputFunc)
}

func getOrdersInstrumentIds(orders []*robinhood.Order) []string {
	return utils.GetIDs(orders, func(order *robinhood.Order) string {
		return order.Instrument
	})
}

func outputOrdersCSV(
	w io.Writer,
	orders []*robinhood.Order,
	instruments []*robinhood.Instrument,
	markets []*robinhood.Market,
) error {
	instrumentByURL := getInstrumentByURL(instruments)
	marketByURL := getMarketByURL(markets)
	header := []string{
		"side",
		"state",
		"market",
		"symbol",
		"settle_date",
		"quantity",
		"price",
		"fee",
		"principal",
		"currency",
	}
	writer := csv.NewWriter(w)
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, order := range orders {
		instrument, ok := instrumentByURL[order.Instrument]
		if !ok {
			return fmt.Errorf("missing instrument: %s", order.Instrument)
		}
		market, ok := marketByURL[instrument.Market]
		if !ok {
			return fmt.Errorf("instrument %s missing market: %s", instrument.Symbol, instrument.Market)
		}
		record := []string{
			order.Side,
			order.State,
			market.Acronym,
			instrument.Symbol,
			order.LastTransactionAt,
			order.CumulativeQuantity,
			order.AveragePrice,
			order.Fees,
			robinhood.GetNotional(order),
			robinhood.GetCurrencyCode(order),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
