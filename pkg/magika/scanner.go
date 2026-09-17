package magika

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/neura-flow/go-magika/pkg/runtime"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

import (
	"embed"
)

const (
	modelName string = "standard_v3_3"
)

//go:embed all:models
var modelFS embed.FS

type Scanner struct {
	modelConfig *ModelConfig
	ckb         map[string]ContentType

	runtime *runtime.Runtime
	session *ort.Session
}

func NewScanner() (*Scanner, error) {
	s := &Scanner{}

	var err error
	if s.runtime, err = runtime.New(); err != nil {
		return nil, err
	}

	if err = s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Scanner) init() error {
	// load model config
	configBytes, err := modelFS.ReadFile(fmt.Sprintf("models/%s/config.min.json", modelName))
	if err != nil {
		return err
	}
	var modelConfig *ModelConfig
	if err = json.Unmarshal(configBytes, &modelConfig); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	s.modelConfig = modelConfig

	// load content type kb
	ckbBytes, err := modelFS.ReadFile("models/content_types_kb.min.json")
	if err != nil {
		return err
	}
	var ckb map[string]ContentType
	if err = json.Unmarshal(ckbBytes, &ckb); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	for label, ct := range ckb {
		ct.Label = label
		ckb[label] = ct
	}
	s.ckb = ckb

	// load model file
	file, err := modelFS.Open(fmt.Sprintf("models/%s/model.onnx", modelName))
	if err != nil {
		return fmt.Errorf("failed to open model file: %w", err)
	}
	if s.session, err = s.runtime.NewSession(file, &ort.SessionOptions{
		IntraOpNumThreads: 0,
	}); err != nil {
		return fmt.Errorf("failed to create detect session: %w", err)
	}
	return nil
}

// Scan scans the given reader containing the given size of bytes, and
// returns the inferred content type.
// It is safe for concurrent use.
func (s *Scanner) Scan(ctx context.Context, r io.ReaderAt, size int) (ContentType, error) {
	ct, _, err := s.scanScore(ctx, r, size)
	return ct, err
}

// scanScore scans the given reader containing the given size of bytes, and
// returns the inferred content type and its score.
// It is safe for concurrent use.
func (s *Scanner) scanScore(ctx context.Context, r io.ReaderAt, size int) (ContentType, float32, error) {
	if size == 0 {
		return s.ckb[contentTypeLabelEmpty], 1, nil
	}
	ft, err := ExtractFeatures(s.modelConfig, r, size)
	if err != nil {
		return ContentType{}, 0, fmt.Errorf("extract features: %w", err)
	}
	// Do not use the model for small files.
	if ft.Beg[s.modelConfig.MinFileSizeForDl-1] == int32(s.modelConfig.PaddingToken) {
		if utf8.Valid(ft.firstBlock) {
			return s.ckb[contentTypeLabelTxt], 1, nil
		} else {
			return s.ckb[contentTypeLabelUnknown], 1, nil
		}
	}
	scores, err := s.run(ctx, ft.Flatten())
	if err != nil {
		return ContentType{}, 0, fmt.Errorf("run onnx: %w", err)
	}
	if len(scores) == 0 {
		return ContentType{}, 0, errors.New("run onnx: empty result")
	}
	best := 0
	for i, v := range scores {
		if v > scores[best] {
			best = i
		}
	}
	ct, err := s.contentType(best, scores[best])
	if err != nil {
		return ContentType{}, 0, fmt.Errorf("get content type: %w", err)
	}
	return ct, scores[best], nil
}

func (s *Scanner) run(ctx context.Context, features []int32) ([]float32, error) {
	inputValue, err := ort.NewTensorValue(s.runtime.RT(), features, []int64{1, 2048})
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess input data, err: %v", err)
	}

	inputNames := s.session.InputNames()
	inputValues := make(map[string]*ort.Value)
	inputValues[inputNames[0]] = inputValue
	outputs, err := s.runtime.Run(ctx, s.session, inputValues)
	if err != nil {
		return nil, fmt.Errorf("failed to run infer, err: %v", err)
	}

	outputNames := s.session.OutputNames()
	data0, _, err0 := ort.GetTensorData[float32](outputs[outputNames[0]])
	if err0 != nil {
		return nil, err0
	}

	return data0, nil
}

func (s *Scanner) contentType(best int, score float32) (ContentType, error) {
	l := s.modelConfig.TargetLabelsSpace[best]
	ct, ok := s.ckb[l]
	if !ok {
		return ContentType{}, fmt.Errorf("no content type found for %q", l)
	}
	th := s.modelConfig.MediumConfidenceThreshold
	if t, ok := s.modelConfig.Thresholds[l]; ok {
		th = t
	}
	// Return the inferred content type if the threshold is met, otherwise
	// falls back to a relevant default.
	switch {
	case score >= th:
	case ct.IsText:
		l = contentTypeLabelTxt
	default:
		l = contentTypeLabelUnknown
	}
	ct, ok = s.ckb[l]
	if !ok {
		return ContentType{}, fmt.Errorf("no content type found for %q", l)
	}
	if l, ok = s.modelConfig.Overwrite[l]; ok {
		if ct, ok = s.ckb[l]; !ok {
			return ContentType{}, fmt.Errorf("no content type found for %q", l)
		}
	}
	return ct, nil
}

func (s *Scanner) Close() error {
	if s.session != nil {
		s.session.Close()
	}
	return nil
}
