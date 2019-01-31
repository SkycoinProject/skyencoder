# skyencoder
Code-generation based encoder for Skycoin

[![Build Status](https://travis-ci.com/skycoin/skyencoder.svg?branch=master)](https://travis-ci.com/skycoin/skyencoder)
[![GoDoc](https://godoc.org/github.com/skycoin/skyencoder?status.svg)](https://godoc.org/github.com/skycoin/skyencoder)

## Introduction

`skyencoder` generates a file with encode and decode methods for a struct, using the [Skycoin encoding format](github.com/skycoin/skycoin/wiki/encoder).

Skycoin's [`package encoder`](https://godoc.org/github.com/skycoin/skycoin/src/cipher/encoder) has a reflect-based encoder that can be used at runtime.

## Installation

```sh
go install github.com/skycoin/skyencoder
```

This installs `skyencoder` to `$GOPATH/bin`.  Make sure `$GOPATH/bin` is in
your shell environment's `$PATH` variable in order to invoke it in the shell.

## go:generate

To use `go generate` to generate the code, add a directive like this in the file where the struct is defined:

```go
// go:generate skyencoder -struct Foo
```

Then, use `go:generate` to generate it:

```sh
go generate github.com/foo/foo
```

## CLI Usage

```
» go run cmd/skyencoder/skyencoder.go --help
Usage of skyencoder:
	skyencoder [flags] -struct T [go import path e.g. github.com/skycoin/skycoin/src/coin]
	skyencoder [flags] -struct T files... # Must be a single package
Flags:
  -output-file string
    	output file name; default <struct_name>_skyencoder.go
  -output-path string
    	output path; defaults to the package's path, or the file's containing folder
  -package string
    	package name for the output; if not provided, defaults to the struct's package
  -struct string
    	struct name, must be set
  -tags string
    	comma-separated list of build tags to apply
exit status 2
```

`skyencoder` generates a file with encode and decode methods for a struct, using the [Skycoin encoding format](github.com/skycoin/skycoin/wiki/encoder).

By default, the generated file is written to the same package as the source struct.

If you wish to have the file written to a different location, use `-package` to control the name of the destination package,
`-output-path` to control the destination path, and `-output-file` to control the destination filename.

Build tags can be applied to the loaded package with `-tags`.

## CLI Examples

Generate code for struct `coin.SignedBlock` in `github.com/skycoin/skycoin/src/coin`:

```sh
go run cmd/skyencoder/skyencoder.go -struct SignedBlock github.com/skycoin/skycoin/src/coin
```

Generate code for struct `Foo` in `/tmp/foo/foo.go`:

```sh
go run cmd/skyencoder/skyencoder.go -struct Foo /tmp/foo/foo.go
```

*Note: absolute paths can only point to a Go file. If there are multiple Go files in that same path, all of them must be included.*

Generate code for struct `coin.SignedBlock` in `github.com/skycoin/skycoin/src/coin`, but sent to an external package:

```sh
go run cmd/skyencoder/skyencoder.go -struct SignedBlock -package foo -output-path /tmp/foo github.com/skycoin/skycoin/src/coin
```

*Note: do not use `-package` if the generated file is going to be in the same package as the struct*
