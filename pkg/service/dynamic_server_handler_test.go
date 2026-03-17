package service

import (
	"encoding/json"
	"testing"

	"github.com/modelpack/modctl/pkg/backend"
	"github.com/stretchr/testify/require"
)

func TestArtifactJSON_MarshalJSON(t *testing.T) {
	artifact := &backend.InspectedModelArtifact{
		ID:           "sha256:config",
		Digest:       "sha256:manifest",
		Architecture: "amd64",
		CreatedAt:    "2026-03-16T00:00:00Z",
		Family:       "llama",
		Format:       "safetensors",
		Name:         "demo",
		ParamSize:    "7B",
		Precision:    "fp16",
		Quantization: "none",
		Layers: []backend.InspectedModelArtifactLayer{
			{
				MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
				Digest:    "sha256:layer",
				Size:      123,
				Filepath:  "weights/model.safetensors",
			},
		},
	}

	payload, err := json.Marshal(newArtifactJSON(artifact))
	require.NoError(t, err)

	var actual map[string]any
	require.NoError(t, json.Unmarshal(payload, &actual))

	require.Equal(t, "sha256:config", actual["id"])
	require.Equal(t, "sha256:manifest", actual["digest"])
	require.Equal(t, "amd64", actual["architecture"])
	require.Equal(t, "2026-03-16T00:00:00Z", actual["created_at"])
	require.Equal(t, "7B", actual["param_size"])
	require.Equal(t, "xxx", actual["test"])

	_, hasCamelID := actual["Id"]
	_, hasCamelCreatedAt := actual["CreatedAt"]
	require.False(t, hasCamelID)
	require.False(t, hasCamelCreatedAt)

	layers, ok := actual["layers"].([]any)
	require.True(t, ok)
	require.Len(t, layers, 1)

	firstLayer, ok := layers[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "application/vnd.oci.image.layer.v1.tar+gzip", firstLayer["media_type"])
	require.Equal(t, "sha256:layer", firstLayer["digest"])
	require.Equal(t, "weights/model.safetensors", firstLayer["filepath"])

	_, hasCamelMediaType := firstLayer["MediaType"]
	require.False(t, hasCamelMediaType)
}

func TestArtifactJSON_MarshalJSON_NilArtifact(t *testing.T) {
	payload, err := json.Marshal(newArtifactJSON(nil))
	require.NoError(t, err)
	require.Equal(t, "null", string(payload))
}
