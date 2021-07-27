package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pmatseykanets/mempool/mempool"
	"github.com/pmatseykanets/mempool/types"
	"github.com/shopspring/decimal"
)

var version, buildVersion string // Set by go build.

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var (
		inPath, outPath string
		showVersion     bool
		cap             = 5000 // Default capacity.
	)
	flag.IntVar(&cap, "cap", cap, "Mempool capacity")
	flag.StringVar(&inPath, "in", "", "Input file path (default stdin)")
	flag.StringVar(&outPath, "out", "", "Output file path (default stdout)")
	flag.BoolVar(&showVersion, "version", showVersion, "Print version and exit")
	flag.Parse()

	if showVersion {
		fmt.Println("mempool-cli version ", version)
		if buildVersion != "" {
			fmt.Println(buildVersion)
		}

		os.Exit(0)
	}

	var (
		inFile, outFile *os.File
		err             error
	)

	if inPath != "" {
		inFile, err = os.Open(inPath)
		if err != nil {
			return err
		}
		defer inFile.Close()
	} else {
		inFile = os.Stdin
	}

	if outPath != "" {
		outFile, err = os.Create(outPath)
		if err != nil {
			return err
		}
		defer outFile.Close()
	} else {
		outFile = os.Stdout
	}

	var (
		pool    = mempool.New(cap)
		scanner = bufio.NewScanner(inFile)
		txt     string
		line    int
		tx      *types.Transaction
	)

	for scanner.Scan() {
		line++
		txt = scanner.Text()
		if txt == "" {
			continue
		}

		tx, err = parseTx(txt)
		if err != nil {
			return err
		}

		pool.Push(tx)
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	pool.GetPrioritized(func(tx *types.Transaction) bool {
		_, err = writeTx(outFile, tx)
		if err != nil {
			return false
		}

		outFile.WriteString("\n")

		return true
	})

	return err
}

// parseTx attempts to parse a line with transaction data
// in the form of TxHash=... Gas=... FeePerGas=... Signature=...
func parseTx(s string) (*types.Transaction, error) {
	var (
		hash, perGas, sig string
		tx                = &types.Transaction{}
	)
	n, err := fmt.Sscanf(s, "TxHash=%s Gas=%d FeePerGas=%s Signature=%s", &hash, &tx.Gas, &perGas, &sig)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, fmt.Errorf("unexpected number of elements %d", n)
	}

	err = tx.Hash.SetString(hash)
	if err != nil {
		return nil, err
	}

	err = tx.Signature.SetString(sig)
	if err != nil {
		return nil, err
	}

	tx.FeePerGas, err = decimal.NewFromString(perGas)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// writeTx writes transaction data to w.
func writeTx(w io.Writer, tx *types.Transaction) (int, error) {
	return fmt.Fprintf(w,
		"TxHash=%s Gas=%d FeePerGas=%s Signature=%s",
		strings.ToUpper(tx.Hash.String()),
		tx.Gas,
		tx.FeePerGas,
		strings.ToUpper(tx.Signature.String()),
	)
}
