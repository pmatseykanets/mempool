# mempool

[![build](https://github.com/pmatseykanets/mempool/actions/workflows/.build.yml/badge.svg)](https://github.com/pmatseykanets/mempool/actions/workflows/.build.yml)

Priority based mempool (toy) implementation.

## Installing

```sh
go get github.com/pmatseykanets/mempool
```

## Usage

Sample usage:

```go
// Instantiate a mempool.
const capacity = 5000
pool := mempool.New(capacity)

// Add transactions.
pool.Push(tx)
//...

// Use an iterator to list prioritized transactions.
pool.GetPrioritized(func(tx *types.Transaction) bool {
    fmt.Println(tx.Hash)
}
```

## mempool-cli

A utility that reads transactions from a text file, loads them into a mempool, and produces the file with prioritized list of transactions. The transactions are formatted as follows, with one transaction per line.

```txt
TxHash=40E10C7CF56A738C0B8AD4EE30EA8008C7B2334B3ADA195083F8CB18BD3911A0 Gas=729000 FeePerGas=0.11134106816568039 Signature=6386A3893BEB6A5A64E0677F406634E791DEE78D49CF30581AE5281D4094E495E671647EF5E7FD2D207AB8EBA0EA693703E9C368402731BE99E81BDB748EA662
```

mempool-cli options:

```txt
  -cap int
        Mempool capacity (default 5000)
  -in string
        Input file path (default stdin)
  -out string
        Output file path (default stdout)
  -version
        Print version and exit
```

Example usage:

```sh
mempool-cli -in transactions.txt -out prioritized-transactions.txt -cap 5000
```

or

```sh
mempool-cli < transactions.txt > prioritized-transactions.txt
```

Install `mempool-cli`:

```sh
GO111MODULE=on go install github.com/pmatseykanets/mempool/cmd/mempool-cli@latest
```

Build `mempool-cli`:

```sh
make build
```
