package models

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/ajeetraina/aiwatch/pkg/logger"
)

// DebugInfo contains information about the Docker environment
type DebugInfo struct {
	DockerVersion     string `json:"dockerVersion"`
	DockerPath        string `json:"dockerPath"`
	UserID            string `json:"userID"`
	GroupID           string `json:"groupID"`
	DockerSocketPath  string `json:"dockerSocketPath"`
	DockerSocketPerms string `json:"dockerSocketPerms"`
	Error             string `json:"error,omitempty"`
}

// HandleDebugDocker provides debugging information about Docker
func HandleDebugDocker(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	debugInfo := DebugInfo{}
	
	// Get Docker version
	versionCmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	versionOutput, err := versionCmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get Docker version")
		debugInfo.Error = err.Error()
	} else {
		debugInfo.DockerVersion = string(versionOutput)
	}
	
	// Get Docker path
	whichCmd := exec.Command("which", "docker")
	whichOutput, err := whichCmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Failed to find Docker path")
	} else {
		debugInfo.DockerPath = string(whichOutput)
	}
	
	// Get current user and group ID
	idCmd := exec.Command("id")
	idOutput, err := idCmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user ID info")
	} else {
		debugInfo.UserID = string(idOutput)
	}
	
	// Check Docker socket
	lsCmd := exec.Command("ls", "-la", "/var/run/docker.sock")
	lsOutput, err := lsCmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Failed to check Docker socket")
	} else {
		debugInfo.DockerSocketPath = "/var/run/docker.sock"
		debugInfo.DockerSocketPerms = string(lsOutput)
	}
	
	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(debugInfo)
}
