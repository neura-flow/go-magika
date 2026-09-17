package libs

import (
	"fmt"
	"runtime"
	"testing"
)

func TestLoad(t *testing.T) {
	loader := NewLoader()
	libPath, err := loader.Load(getLibraryName())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("libPath: %s\n", libPath)
}

func getLibraryName() string {
	if runtime.GOOS == "darwin" {
		return "libonnxruntime.dylib"
	} else if runtime.GOOS == "windows" {
		return "libonnxruntime.dll"
	} else if runtime.GOOS == "linux" {
		return "libonnxruntime.so"
	}
	return "libonnxruntime.so"
}
