# Go library

This directory contains the Go library for Magika.

The inference relies on the [ONNX Runtime](https://onnxruntime.ai/), and interfacing with the ONNX Runtime
[onnxruntime-purego](github.com/shota3506/onnxruntime-purego).

## Usage
As illustrated in [`example/main.go`](./example/main.go), calling magika from go boils down
to creating a scanner associated with a given model, and scanning the content.

```golang

// This package illustrates the usage of the Magika go binding.
//
// It requires the onnxruntime and the Magika assets to be accessible.
// onnxruntime is available on https://github.com/microsoft/onnxruntime/releases
// Magika asserts are available on https://github.com/google/magika/tree/main/assets
//

package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/neura-flow/go-magika/pkg/magika"
)

func main() {
	// Create a scanner.
	s, err := magika.NewScanner()
	if err != nil {
		log.Fatalf("NewScanner failed: %v", err)
	}
	defer s.Close()

	// Scan
	ct, err := s.Scan(context.Background(), strings.NewReader("Hello world"), 11)
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}
	fmt.Printf("%+v\n", ct)
}


```