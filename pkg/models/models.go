package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/ajeetraina/aiwatch/pkg/logger"
)

// Model represents a Docker model
type Model struct {
	Name         string `json:"name"`
	Parameters   string `json:"parameters"`
	Quantization string `json:"quantization"`
	Architecture string `json:"architecture"`
	ModelID      string `json:"modelId"`
	Created      string `json:"created"`
	Size         string `json:"size"`
}

// GetAvailableModels retrieves the list of available models from Docker Model Runner
func GetAvailableModels() ([]Model, error) {
	log := logger.GetLogger()
	
	// Execute docker model ls command
	cmd := exec.Command("docker", "model", "ls", "--format", "{{.NAME}}\t{{.PARAMETERS}}\t{{.QUANTIZATION}}\t{{.ARCHITECTURE}}\t{{.MODEL_ID}}\t{{.CREATED}}\t{{.SIZE}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute docker model ls command")
		return nil, fmt.Errorf("failed to list docker models: %v", err)
	}

	// Parse the output
	models := []Model{}
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		fields := strings.Split(line, "\t")
		if len(fields) < 7 {
			log.Warn().Str("line", line).Msg("Invalid model line format")
			continue
		}
		
		model := Model{
			Name:         fields[0],
			Parameters:   fields[1],
			Quantization: fields[2],
			Architecture: fields[3],
			ModelID:      fields[4],
			Created:      fields[5],
			Size:         fields[6],
		}
		
		models = append(models, model)
	}

	log.Info().Int("count", len(models)).Msg("Retrieved available Docker models")
	return models, nil
}

// HandleListModels returns the list of available models as JSON
func HandleListModels(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	models, err := GetAvailableModels()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get available models")
		http.Error(w, "Failed to retrieve models", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}
