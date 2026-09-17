package runtime

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/neura-flow/go-magika/pkg/runtime/libs"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

const (
	ApiVersion = 23
)

type Runtime struct {
	rt  *ort.Runtime
	env *ort.Env
}

func New() (*Runtime, error) {
	libFile, err := libs.NewLoader().Load(getLibraryName())
	if err != nil {
		return nil, err
	}
	rt, err := ort.NewRuntime(libFile, ApiVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime: %w", err)
	}
	env, err := rt.NewEnv("go-magika", ort.LoggingLevelError)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}
	r := &Runtime{
		rt:  rt,
		env: env,
	}
	return r, nil
}

func (r *Runtime) RT() *ort.Runtime {
	return r.rt
}

func (r *Runtime) Env() *ort.Env {
	return r.env
}

func (r *Runtime) NewSession(rdr io.Reader, options *ort.SessionOptions) (*ort.Session, error) {
	return r.rt.NewSessionFromReader(r.env, rdr, options)
}

func (r *Runtime) Run(ctx context.Context, session *ort.Session, values map[string]*ort.Value) (map[string]*ort.Value, error) {
	return session.Run(ctx, values)
}

type OutputTensor struct {
	Name   string
	Shape  []int64
	Logits []float32
}

func (r *Runtime) Close() error {
	if r.env != nil {
		r.env.Close()
	}
	return r.rt.Close()
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
