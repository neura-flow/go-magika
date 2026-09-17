package magika

const (
	modelConfigFile    = "config.min.json"
	contentTypesKBFile = "content_types_kb.min.json"
	modelFile          = "model.onnx"
	modelsDir          = "models"
)

// ModelConfig holds the portion of Magika's model configuration that is relevant for inference.
type ModelConfig struct {
	BegSize                   int                `json:"beg_size"`
	MidSize                   int                `json:"mid_size"`
	EndSize                   int                `json:"end_size"`
	UseInputsAtOffsets        bool               `json:"use_inputs_at_offsets"`
	MediumConfidenceThreshold float32            `json:"medium_confidence_threshold"`
	MinFileSizeForDl          int64              `json:"min_file_size_for_dl"`
	PaddingToken              int                `json:"padding_token"`
	BlockSize                 int                `json:"block_size"`
	TargetLabelsSpace         []string           `json:"target_labels_space"`
	Thresholds                map[string]float32 `json:"thresholds"`
	Overwrite                 map[string]string  `json:"overwrite_map"`
}

const (
	contentTypeLabelEmpty   = "empty"
	contentTypeLabelTxt     = "txt"
	contentTypeLabelUnknown = "unknown"
)

// ContentType holds the definition of a content type.
type ContentType struct {
	Label       string   // As keyed in the content types KB.
	MimeType    string   `json:"mime_type"`
	Group       string   `json:"group"`
	Description string   `json:"description"`
	Extensions  []string `json:"extensions"`
	IsText      bool     `json:"is_text"`
}
