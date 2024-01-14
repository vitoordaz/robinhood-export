package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
)

type arguments struct {
	username string // -u, robinhood account username or email
	output   string // -o, path to output file
	verbose  bool   // -v, enable verbose messages
	all      bool   // -a, return everything
}

const (
	exitCodeError = 2
	exitCodeOk    = 0
)

var (
	logVerbose = log.New(os.Stdout, "D: ", 0)
	logError   = log.New(os.Stderr, "ERROR: ", 0)

	helpCmd = flag.NewFlagSet("help", flag.ExitOnError)

	dividendsCmd         = flag.NewFlagSet("dividends", flag.ExitOnError)
	dividendsCmdUsername = dividendsCmd.String("u", "", "Robinhood account username or email.") // optional
	dividendsCmdOutput   = dividendsCmd.String("o", "", "path to output file.")                 // optional
	dividendsCmdVerbose  = dividendsCmd.Bool("v", false, "enable verbose messages.")            // optional

	ordersCmd         = flag.NewFlagSet("orders", flag.ExitOnError)
	ordersCmdUsername = ordersCmd.String("u", "", "Robinhood account username or email.") // optional
	ordersCmdOutput   = ordersCmd.String("o", "", "path to output file.")                 // optional
	ordersCmdVerbose  = ordersCmd.Bool("v", false, "enable verbose messages.")            // optional

	positionsCmd         = flag.NewFlagSet("positions", flag.ExitOnError)
	positionsCmdUsername = positionsCmd.String("u", "", "Robinhood account username or email.") // optional
	positionsCmdOutput   = positionsCmd.String("o", "", "path to output file.")                 // optional
	positionsCmdAll      = positionsCmd.Bool("a", false, "return all positions (even closed).") // optional
	positionsCmdVerbose  = positionsCmd.Bool("v", false, "enable verbose messages.")            // optional
)

func main() {
	flag.Usage = printUsage
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(exitCodeError)
	}

	dividendsCmd.Usage = printDividendsUsage
	helpCmd.Usage = printHelpUsage
	ordersCmd.Usage = printOrdersUsage
	positionsCmd.Usage = printPositionsUsage

	switch flag.Arg(0) {
	case "dividends":
		if err := dividendsCmd.Parse(os.Args[2:]); err != nil {
			logError.Println(err)
			dividendsCmd.Usage()
			os.Exit(exitCodeError)
		}
		doDividends(arguments{
			username: *dividendsCmdUsername,
			verbose:  *dividendsCmdVerbose,
			output:   *dividendsCmdOutput,
		})
		os.Exit(exitCodeError)
	case "help":
		if err := helpCmd.Parse(os.Args[2:]); err != nil || len(helpCmd.Args()) < 1 {
			if err != nil {
				logError.Println(err)
			}
			helpCmd.Usage()
			os.Exit(exitCodeError)
		}
		doHelp(helpCmd.Arg(0))
		os.Exit(exitCodeOk)
	case "orders":
		if err := ordersCmd.Parse(os.Args[2:]); err != nil {
			logError.Println(err)
			ordersCmd.Usage()
			os.Exit(exitCodeError)
		}
		doOrders(arguments{username: *ordersCmdUsername, verbose: *ordersCmdVerbose, output: *ordersCmdOutput})
		os.Exit(exitCodeOk)
	case "positions":
		if err := positionsCmd.Parse(os.Args[2:]); err != nil {
			if err != nil {
				logError.Println(err)
			}
			positionsCmd.Usage()
			os.Exit(exitCodeError)
		}
		doPositions(arguments{
			username: *positionsCmdUsername,
			verbose:  *positionsCmdVerbose,
			output:   *positionsCmdOutput,
			all:      *positionsCmdAll,
		})
		os.Exit(exitCodeOk)
	default:
		flag.Usage()
		os.Exit(exitCodeError)
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", fmt.Errorf("line is too long")
	}
	return string(line), nil
}

func getAuthToken(
	ctx context.Context,
	client robinhood.Client,
	username string,
) (*robinhood.ResponseToken, error) {
	password := ""
	reader := bufio.NewReader(os.Stdin)
	for username == "" || password == "" {
		if username == "" {
			fmt.Print("Enter username (email): ")
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			username = strings.TrimSpace(line)
			if username == "" {
				logError.Println("username (email) is required")
				continue
			}
		}
		if password == "" {
			fmt.Print("Enter password: ")
			line, err := term.ReadPassword(0)
			if err != nil {
				return nil, fmt.Errorf("ERROR: %w", err)
			}
			password = strings.TrimSpace(string(line))
			fmt.Println() // NOTE: term.ReadPassword doesn't add new line after enter
			if password == "" {
				logError.Println("password is required")
				continue
			}
		}
	}
	var (
		err  error
		resp *robinhood.ResponseToken
		mfa  = ""
	)
	for resp == nil || resp.AccessToken == "" {
		msg := "Trying to log in using username, password"
		if mfa != "" {
			msg += " and OTP code"
		}
		logVerbose.Println(msg)
		resp, err = client.GetToken(ctx, username, password, mfa)
		if err != nil {
			return nil, err
		}
		if resp.MFARequired {
			fmt.Print("Enter OTP code: ")
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			mfa = strings.TrimSpace(line)
			if mfa == "" {
				logError.Println("OTP code is required")
				continue
			}
		}
	}
	logVerbose.Println("Successfully logged in")
	return resp, nil
}
