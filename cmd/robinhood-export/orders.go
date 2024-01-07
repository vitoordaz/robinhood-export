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

	ctx := context.Background()
	client := robinhood.New()

	token, err := getAuthToken(ctx, client, args.username)
	if err != nil {
		logError.Fatalln(err)
	}

	logVerbose.Println("loading orders")
	orders, err := loadOrders(ctx, client, token)
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

	var f *os.File
	if args.output == "" {
		f = os.Stdout
	} else {
		if f, err = os.Create(args.output); err != nil {
			logError.Fatalln(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				logError.Println(err) // still exit with code 0, because this error is not critical
			}
		}()
	}
	if err := outputOrders(f, orders, instruments, markets); err != nil {
		logError.Fatalln(err)
	}
}

func loadOrders(
	ctx context.Context,
	client robinhood.Client,
	token *robinhood.ResponseToken,
) ([]*robinhood.Order, error) {
	orders, err := utils.LoadList(ctx, func(c context.Context, cursor string) ([]*robinhood.Order, string, error) {
		result, err := client.GetOrders(c, token, cursor)
		if err != nil {
			return nil, "", err
		}
		return result.Results, result.Next, nil
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func getOrdersInstrumentIds(orders []*robinhood.Order) []string {
	return utils.GetIDs(orders, func(order interface{}) string {
		return order.(*robinhood.Order).Instrument
	})
}

func outputOrders(
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
