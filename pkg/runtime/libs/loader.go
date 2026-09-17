package libs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/neura-flow/go-magika/pkg/util"
)

const (
	version = "onnxruntime-v1.27.0"
)

type Loader struct {
	tempDir string
	libFile string
	loaded  bool
	lock    sync.Mutex
}

func NewLoader() *Loader {
	return &Loader{
		tempDir: util.ExpandPath("~/.go-magika"),
	}
}

func (l *Loader) Load(libName string) (string, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.loaded {
		return l.libFile, nil
	}

	libFile := filepath.Join(l.getTargetFolder(), libName)
	if exists, _ := util.Exists(libFile); exists {
		l.loaded = true
		l.libFile = libFile
		return libFile, nil
	}

	data, err := libFS.ReadFile(filepath.Join(l.getSourceFolder(), libName))
	if err != nil {
		return "", err
	}

	if err = os.MkdirAll(filepath.Dir(libFile), 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	if err = os.WriteFile(libFile, data, 0755); err != nil {
		return "", err
	}

	l.loaded = true
	l.libFile = libFile
	return libFile, nil
}

func (l *Loader) getSourceFolder() string {
	return fmt.Sprintf("%s-%s/%s", runtime.GOOS, runtime.GOARCH, version)
}

func (l *Loader) getTargetFolder() string {
	return filepath.Join(l.tempDir, fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH), version)
}
