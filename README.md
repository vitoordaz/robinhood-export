# Robinhood-Export

[![Test](https://github.com/vitoordaz/robinhood-export/actions/workflows/test.yml/badge.svg?branch=mainline)](https://github.com/vitoordaz/robinhood-export/actions/workflows/test.yml)
[![Linter](https://github.com/vitoordaz/robinhood-export/actions/workflows/linter.yml/badge.svg?branch=mainline)](https://github.com/vitoordaz/robinhood-export/actions/workflows/linter.yml)

`robinhood-export` is a CLI tool for exporting data from [Robinhood](https://robinhood.com/).

## Commands

```
robinhood-export <command> [arguments]
```

| Command | Description |
|---|---|
| `accounts` | Export all accounts |
| `dividends` | Export all dividends |
| `instruments` | Export all instruments |
| `options orders` | Export all options orders |
| `orders` | Export all orders |
| `positions` | Export all positions |

Use `robinhood-export help <command>` for more information about a command.

## Arguments

Most commands accept the following arguments:

| Flag | Description | Default |
|---|---|---|
| `-u` | Robinhood account username or email | _(prompted)_ |
| `-t` | Path to a token file (skips interactive login) | — |
| `-o` | Path to output file (defaults to stdout) | — |
| `-f` | Output format: `json` or `csv` | `json` |
| `-v` | Enable verbose messages | `false` |

The `positions` command also accepts `-a` to include closed positions (default: open only).

The `instruments` command does not require authentication and has no `-u` or `-t` flags.

## Authentication

On first run, the tool prompts for your Robinhood username (email) and password. If MFA is enabled on your account, you will also be prompted for an OTP code.

To avoid repeated logins, save a token to a file and pass it with `-t`:

```sh
# Log in interactively and save the token
robinhood-export orders -u user@example.com -o orders.json
# On subsequent runs, reuse the token
robinhood-export orders -t token.json -o orders.json
```

## Output Formats

Both `json` and `csv` output formats are supported via the `-f` flag.

**JSON** (default): each record is written as a newline-delimited JSON object.

**CSV**: a header row followed by one row per record. Available CSV columns per command:

| Command | CSV columns |
|---|---|
| `accounts` | account number, user, type, brokerage account type, state, cash, locked |
| `instruments` | name, symbol, type, market |
| `orders` | side, state, quantity, average price, price, symbol, name, created at, updated at |
| `positions` | account, quantity, average buy price, symbol, name |
| `dividends` | account, instrument, amount, rate, position, withholding, payable date, paid at, state |
| `options orders` | direction, state, quantity, premium, processed quantity, processed premium, created at, updated at |

## Usage Examples

Export all orders to a CSV file:

```sh
robinhood-export orders -f csv -o orders.csv
```

Export open positions as JSON (stdout):

```sh
robinhood-export positions
```

Export all positions including closed ones:

```sh
robinhood-export positions -a -f csv -o positions.csv
```

Export dividends with verbose logging:

```sh
robinhood-export dividends -v -f csv -o dividends.csv
```

Export instruments (no login required):

```sh
robinhood-export instruments -f csv -o instruments.csv
```

## Installation

### From source

```sh
go install github.com/vitoordaz/robinhood-export/cmd/robinhood-export@latest
```

### Docker

```sh
docker build -t robinhood-export .
docker run --rm -it robinhood-export orders -f csv
```
