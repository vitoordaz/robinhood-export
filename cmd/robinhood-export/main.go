package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/vitoordaz/robinhood-export/internal/robinhood"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
)

type arguments struct {
	username string // -u, robinhood account username or email
	output   string // -o, path to output file
	verbose  bool   // -v, enable verbose messages
}

var (
	logVerbose = log.New(os.Stdout, "D: ", 0)
	logError   = log.New(os.Stderr, "ERROR: ", 0)

	helpCmd = flag.NewFlagSet("help", flag.ExitOnError)

	ordersCmd         = flag.NewFlagSet("orders", flag.ExitOnError)
	ordersCmdUsername = ordersCmd.String("u", "", "Robinhood account username or email.") // optional
	ordersCmdOutput   = ordersCmd.String("o", "", "path to output file.")                 // optional
	ordersCmdVerbose  = ordersCmd.Bool("v", false, "enable verbose messages.")            // optional

	positionsCmd         = flag.NewFlagSet("positions", flag.ExitOnError)
	positionsCmdUsername = positionsCmd.String("u", "", "Robinhood account username or email.") // optional
	positionsCmdOutput   = positionsCmd.String("o", "", "path to output file.")                 // optional
	positionsCmdVerbose  = positionsCmd.Bool("v", false, "enable verbose messages.")            // optional
)

func main() {
	flag.Usage = printUsage
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	helpCmd.Usage = printHelpUsage
	ordersCmd.Usage = printOrdersUsage
	positionsCmd.Usage = printPositionsUsage

	switch flag.Arg(0) {
	case "help":
		if err := helpCmd.Parse(os.Args[2:]); err != nil || len(helpCmd.Args()) < 1 {
			helpCmd.Usage()
			os.Exit(2)
		}
		doHelp(helpCmd.Arg(0))
	case "orders":
		if err := ordersCmd.Parse(os.Args[2:]); err != nil {
			ordersCmd.Usage()
			os.Exit(2)
		}
		doOrders(arguments{username: *ordersCmdUsername, verbose: *ordersCmdVerbose, output: *ordersCmdOutput})
	case "positions":
		if err := positionsCmd.Parse(os.Args[2:]); err != nil {
			positionsCmd.Usage()
			os.Exit(2)
		}
		doPositions(arguments{
			username: *positionsCmdUsername,
			verbose:  *positionsCmdVerbose,
			output:   *positionsCmdOutput,
		})
	default:
		flag.Usage()
		os.Exit(2)
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
) (resp *robinhood.ResponseToken, err error) {
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
			line, err := terminal.ReadPassword(0)
			if err != nil {
				return nil, fmt.Errorf("ERROR: %v", err)
			}
			password = strings.TrimSpace(string(line))
			fmt.Println() // NOTE: terminal.ReadPassword doesn't add new line after enter
			if password == "" {
				logError.Println("password is required")
				continue
			}
		}
	}
	mfa := ""
	for resp == nil || resp.AccessToken == "" {
		msg := "Trying to log in using username, password"
		if mfa != "" {
			msg += " and OTP code"
		}
		logVerbose.Println(msg)
		resp, err = client.GetToken(ctx, username, password, mfa)
		if err != nil {
			return
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
	return
}
