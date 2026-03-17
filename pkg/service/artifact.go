package service

import (
	"context"
	"encoding/json"

	"github.com/modelpack/modctl/pkg/backend"
	"github.com/modelpack/model-csi-driver/pkg/config/auth"
	"github.com/pkg/errors"
)

type artifactJSON struct {
	artifact *backend.InspectedModelArtifact
}

type artifactResponse struct {
	ID           string                  `json:"id"`
	Digest       string                  `json:"digest"`
	Architecture string                  `json:"architecture"`
	CreatedAt    string                  `json:"created_at"`
	Family       string                  `json:"family"`
	Format       string                  `json:"format"`
	Name         string                  `json:"name"`
	ParamSize    string                  `json:"param_size"`
	Precision    string                  `json:"precision"`
	Quantization string                  `json:"quantization"`
	Layers       []artifactLayerResponse `json:"layers"`
	Test         string                  `json:"test,omitempty"`
}

type artifactLayerResponse struct {
	MediaType string `json:"media_type"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	Filepath  string `json:"filepath"`
}

func newArtifactJSON(artifact *backend.InspectedModelArtifact) artifactJSON {
	return artifactJSON{artifact: artifact}
}

func (a artifactJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(newArtifactResponse(a.artifact))
}

func newArtifactResponse(artifact *backend.InspectedModelArtifact) *artifactResponse {
	if artifact == nil {
		return nil
	}

	response := &artifactResponse{
		ID:           artifact.ID,
		Digest:       artifact.Digest,
		Architecture: artifact.Architecture,
		CreatedAt:    artifact.CreatedAt,
		Family:       artifact.Family,
		Format:       artifact.Format,
		Name:         artifact.Name,
		ParamSize:    artifact.ParamSize,
		Precision:    artifact.Precision,
		Quantization: artifact.Quantization,
		Layers:       make([]artifactLayerResponse, 0, len(artifact.Layers)),
		Test:         "xxx",
	}

	for _, layer := range artifact.Layers {
		response.Layers = append(response.Layers, artifactLayerResponse{
			MediaType: layer.MediaType,
			Digest:    layer.Digest,
			Size:      layer.Size,
			Filepath:  layer.Filepath,
		})
	}

	return response
}

func (s *Service) GetArtifact(ctx context.Context, reference string) (*backend.InspectedModelArtifact, error) {
	keyChain, err := auth.GetKeyChainByRef(reference)
	if err != nil {
		return nil, errors.Wrapf(err, "get auth for model: %s", reference)
	}
	plainHTTP := keyChain.ServerScheme == "http"

	b, err := backend.New("")
	if err != nil {
		return nil, errors.Wrap(err, "create modctl backend")
	}

	modelArtifact := NewModelArtifact(b, reference, plainHTTP)

	artifact, err := modelArtifact.Inspect(ctx, reference)
	if err != nil {
		return nil, err
	}

	return artifact, nil
}
