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
	
	// Debug logging
	log.Info().Msg("Attempting to execute 'docker model ls' command")
	
	// Execute docker model ls command
	cmd := exec.Command("docker", "model", "ls", "--format", "{{.NAME}}\t{{.PARAMETERS}}\t{{.QUANTIZATION}}\t{{.ARCHITECTURE}}\t{{.MODEL_ID}}\t{{.CREATED}}\t{{.SIZE}}")
	output, err := cmd.CombinedOutput()
	
	// Log the output and any error for debugging
	if err != nil {
		log.Error().Err(err).Str("output", string(output)).Msg("Failed to execute docker model ls command")
		
		// Try basic docker command to check connectivity
		checkCmd := exec.Command("docker", "version")
		checkOutput, checkErr := checkCmd.CombinedOutput()
		if checkErr != nil {
			log.Error().Err(checkErr).Str("output", string(checkOutput)).Msg("Docker connectivity check failed")
		} else {
			log.Info().Msg("Docker connection verified but 'docker model ls' command failed")
		}
		
		// Check for the docker executable
		whichCmd := exec.Command("which", "docker")
		whichOutput, whichErr := whichCmd.CombinedOutput()
		if whichErr != nil {
			log.Error().Err(whichErr).Msg("Docker executable not found in PATH")
		} else {
			log.Info().Str("path", strings.TrimSpace(string(whichOutput))).Msg("Docker executable found")
		}
		
		return nil, fmt.Errorf("failed to list docker models: %v, output: %s", err, string(output))
	}

	// Parse the output
	models := []Model{}
	lines := strings.Split(string(output), "\n")
	
	log.Info().Str("output", string(output)).Msg("Docker model ls output")
	
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

// GetFallbackModels returns hard-coded models for when Docker commands fail
func GetFallbackModels() []Model {
	return []Model{
		{
			Name:         "ai/llama3.2:1B-Q8_0",
			Parameters:   "1.24 B",
			Quantization: "Q8_0",
			Architecture: "llama",
			ModelID:      "a15c3117eeeb",
			Created:      "5 weeks ago",
			Size:         "1.22 GiB",
		},
		{
			Name:         "ai/qwen3",
			Parameters:   "8.19 B",
			Quantization: "IQ2_XXS/Q4_K_M",
			Architecture: "qwen3",
			ModelID:      "79fa56c07429",
			Created:      "3 days ago",
			Size:         "4.68 GiB",
		},
	}
}

// HandleListModels returns the list of available models as JSON
func HandleListModels(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	models, err := GetAvailableModels()
	if err != nil || len(models) == 0 {
		log.Error().Err(err).Msg("Failed to get available models, using fallback")
		
		// Use fallback models
		fallbackModels := GetFallbackModels()
		
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(fallbackModels)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(models)
}