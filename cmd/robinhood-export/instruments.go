package main

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
	"github.com/vitoordaz/robinhood-export/internal/utils"
)

func doInstruments(args arguments) {
	if !args.verbose {
		logVerbose.SetOutput(io.Discard) // disable verbose logging
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := robinhood.New()

	var outputFunc func(f *os.File) error

	logVerbose.Println("loading instruments")
	switch args.format {
	case "json":
		instruments, err := loadItems(
			ctx,
			func(ctx context.Context, cursor string) (*robinhood.ResponseList[map[string]any], error) {
				return client.GetItems(ctx, robinhood.EndpointInstruments, nil, cursor)
			},
		)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d instruments\n", len(instruments))
		outputFunc = func(f *os.File) error { return outputItemsJSON(f, instruments) }
	case "csv":
		instruments, err := listInstruments(ctx, client)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d instruments\n", len(instruments))
		outputFunc = func(f *os.File) error { return outputInstrumentsCSV(f, instruments) }
	default:
		logError.Fatalln("unsupported output format " + args.format)
	}

	outputData(args.output, outputFunc)
}

func listInstruments(ctx context.Context, client robinhood.Client) ([]*robinhood.Instrument, error) {
	return utils.LoadList(
		ctx,
		func(ctx context.Context, cursor string) ([]*robinhood.Instrument, string, error) {
			result, err := client.ListInstruments(ctx, cursor)
			if err != nil {
				return nil, "", err
			}
			return result.Results, result.Next, nil
		},
	)
}

func outputInstrumentsCSV(w io.Writer, instruments []*robinhood.Instrument) error {
	header := []string{"name", "symbol", "type", "market"}
	writer := csv.NewWriter(w)
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, instrument := range instruments {
		record := []string{
			instrument.Name,
			instrument.Symbol,
			instrument.Type,
			instrument.Market,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
