[![GoDoc](https://godoc.org/github.com/noxer/serialtools?status.svg)](https://godoc.org/github.com/noxer/serialtools)
[![Go Report Card](https://goreportcard.com/badge/github.com/noxer/serialtools)](https://goreportcard.com/report/github.com/noxer/serialtools)
[![Build Status](https://travis-ci.org/noxer/serialtools.svg?branch=master)](https://travis-ci.org/noxer/serialtools)

serialtools
===
This package offers functionality to make working with serial interfaces easier.

LFNormalizer
---
Currently the only available function is the normalization of line endings which can differ for device to device. So far I've encountered `\r` (most common), `\n`, `\r\n` (also quite common), and `\n\r`. The normalizer converts all these line endings into `\n` while reading from the source reader.

Here is an example on how to use the reader:
```go
package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/noxer/serialtools"
)

func main() {
	// we just simulate the data from the serial device
	s := strings.NewReader("this\r\nis\n\rdata\rwith\nmixed\n\r\nnewlines")

	// wrap the reader into a LFNormalizer
	r := serialtools.NewLFNormalizer(s)

    // now you can read from the normalizer and just deal with \n
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if n > 0 {
		fmt.Printf("Read from serial: %s\n", strconv.Quote(string(buf[:n])))
	}

	// don't forget to check for errors!
	if err != nil {
		fmt.Printf("Error while reading: %s\n", err)
	}
}
```
this should print the following text:
```
Read from serial: "this\nis\ndata\nwith\nmixed\n\nnewlines"
```
