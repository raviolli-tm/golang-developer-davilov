package main

import (
	"errors"
	"flag"
	"fmt"
)

var (
	from, to                    string
	limit, offset               int64
	ErrFromVarIsNotDefinedError = errors.New("from param is not defined")
	ErrToVarIsNotDefinedError   = errors.New("to param is not defined")
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	if from == "" {
		fmt.Println("Fatal error: ", ErrFromVarIsNotDefinedError)
		return
	}

	if to == "" {
		fmt.Println("Fatal error: ", ErrToVarIsNotDefinedError)
		return
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println("Fatal error: ", err)
		return
	}
	fmt.Println("Success!")
}
