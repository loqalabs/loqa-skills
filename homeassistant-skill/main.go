/*
 * This file is part of Loqa (https://github.com/loqalabs/loqa).
 * Copyright (C) 2025 Loqa Labs
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"github.com/loqalabs/loqa-hub/internal/skills"
)

// HomeAssistantSkill implements the Home Assistant Voice Preview Edition integration
type HomeAssistantSkill struct {
	logger     *zap.Logger
	config     *skills.SkillConfig
	status     skills.SkillStatus
	httpClient *http.Client
	haConfig   *HAConfig
}

// HAConfig holds Home Assistant specific configuration
type HAConfig struct {
	BaseURL     string `json:"base_url"`
	AccessToken string `json:"access_token"`
	MQTTEnabled bool   `json:"mqtt_enabled"`
	MQTTTopic   string `json:"mqtt_topic"`
	DeviceID    string `json:"device_id"`
	DeviceName  string `json:"device_name"`
	Timeout     int    `json:"timeout_seconds"`
}

// HAVoiceRequest represents a Home Assistant voice request
type HAVoiceRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
	DeviceID string `json:"device_id"`
}

// HAVoiceResponse represents a Home Assistant voice response
type HAVoiceResponse struct {
	Success      bool   `json:"success"`
	Text         string `json:"text"`
	SpeechText   string `json:"speech_text"`
	AudioURL     string `json:"audio_url,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// NewSkillPlugin is the plugin entry point
func NewSkillPlugin(logger *zap.Logger) skills.SkillPlugin {
	return &HomeAssistantSkill{
		logger: logger,
		status: skills.SkillStatus{
			State:   skills.SkillStateLoading,
			Healthy: false,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Initialize initializes the Home Assistant skill
func (s *HomeAssistantSkill) Initialize(ctx context.Context, config *skills.SkillConfig) error {
	s.config = config
	
	// Parse HA-specific configuration
	if err := s.parseHAConfig(); err != nil {
		s.status.State = skills.SkillStateError
		s.status.LastError = fmt.Sprintf("Failed to parse HA config: %v", err)
		return err
	}
	
	// Set HTTP client timeout from config
	if s.haConfig.Timeout > 0 {
		s.httpClient.Timeout = time.Duration(s.haConfig.Timeout) * time.Second
	}
	
	// Test connection to Home Assistant
	if err := s.testConnection(ctx); err != nil {
		s.logger.Warn("Failed to connect to Home Assistant", zap.Error(err))
		s.status.LastError = fmt.Sprintf("HA connection test failed: %v", err)
		// Don't fail initialization - allow skill to load and retry later
	}
	
	s.status.State = skills.SkillStateReady
	s.status.Healthy = true
	
	s.logger.Info("Initialized Home Assistant skill", 
		zap.String("skill", config.SkillID),
		zap.String("version", config.Version),
		zap.String("ha_url", s.haConfig.BaseURL),
		zap.String("device_id", s.haConfig.DeviceID))
	
	return nil
}

// parseHAConfig extracts Home Assistant configuration from skill config
func (s *HomeAssistantSkill) parseHAConfig() error {
	configMap := s.config.Config
	
	s.haConfig = &HAConfig{
		BaseURL:     "http://homeassistant.local:8123",
		DeviceID:    "loqa-voice-assistant",
		DeviceName:  "Loqa Voice Assistant",
		MQTTTopic:   "homeassistant/voice",
		Timeout:     30,
	}
	
	if baseURL, ok := configMap["base_url"].(string); ok {
		s.haConfig.BaseURL = baseURL
	}
	
	if token, ok := configMap["access_token"].(string); ok {
		s.haConfig.AccessToken = token
	}
	
	if deviceID, ok := configMap["device_id"].(string); ok {
		s.haConfig.DeviceID = deviceID
	}
	
	if deviceName, ok := configMap["device_name"].(string); ok {
		s.haConfig.DeviceName = deviceName
	}
	
	if mqttEnabled, ok := configMap["mqtt_enabled"].(bool); ok {
		s.haConfig.MQTTEnabled = mqttEnabled
	}
	
	if mqttTopic, ok := configMap["mqtt_topic"].(string); ok {
		s.haConfig.MQTTTopic = mqttTopic
	}
	
	if timeout, ok := configMap["timeout_seconds"].(float64); ok {
		s.haConfig.Timeout = int(timeout)
	}
	
	// Validate required fields
	if s.haConfig.AccessToken == "" {
		return fmt.Errorf("access_token is required for Home Assistant integration")
	}
	
	return nil
}

// testConnection tests the connection to Home Assistant
func (s *HomeAssistantSkill) testConnection(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/", strings.TrimSuffix(s.haConfig.BaseURL, "/"))
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.haConfig.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to HA: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HA API returned status %d", resp.StatusCode)
	}
	
	return nil
}

// Teardown shuts down the Home Assistant skill
func (s *HomeAssistantSkill) Teardown(ctx context.Context) error {
	s.status.State = skills.SkillStateShutdown
	s.status.Healthy = false
	
	s.logger.Info("Shutdown Home Assistant skill")
	return nil
}

// CanHandle determines if this skill can handle the given intent
func (s *HomeAssistantSkill) CanHandle(intent skills.VoiceIntent) bool {
	// This skill acts as a fallback for all intents not handled by other skills
	// It forwards everything to Home Assistant for processing
	
	// Skip if the skill is not healthy
	if !s.status.Healthy {
		return false
	}
	
	// Always return true for enabled HA skill to act as fallback
	// The priority system will ensure other skills get first chance
	return true
}

// HandleIntent processes a voice intent by forwarding it to Home Assistant
func (s *HomeAssistantSkill) HandleIntent(ctx context.Context, intent *skills.VoiceIntent) (*skills.SkillResponse, error) {
	s.logger.Info("Handling HA intent", 
		zap.String("intent", intent.Intent),
		zap.String("transcript", intent.Transcript),
		zap.String("device_id", intent.DeviceID))
	
	// Update usage stats
	s.status.LastUsed = time.Now()
	s.status.UsageCount++
	
	startTime := time.Now()
	
	// Prepare request for Home Assistant
	haRequest := HAVoiceRequest{
		Text:     intent.Transcript,
		Language: "en", // TODO: Extract from intent or config
		DeviceID: s.haConfig.DeviceID,
	}
	
	// Send to Home Assistant
	haResponse, err := s.sendToHomeAssistant(ctx, haRequest)
	if err != nil {
		s.logger.Error("Failed to send to Home Assistant", zap.Error(err))
		return &skills.SkillResponse{
			Success:      false,
			Message:      "Failed to process with Home Assistant",
			SpeechText:   "Sorry, I couldn't connect to Home Assistant to process that request.",
			Error:        err.Error(),
			ErrorCode:    "ha_connection_error",
			ResponseTime: time.Since(startTime),
		}, nil
	}
	
	// Convert HA response to skill response
	response := &skills.SkillResponse{
		Success:      haResponse.Success,
		Message:      haResponse.Text,
		SpeechText:   haResponse.SpeechText,
		AudioURL:     haResponse.AudioURL,
		ResponseTime: time.Since(startTime),
		Actions: []skills.SkillAction{
			{
				Type:       "home_assistant_command",
				Target:     "homeassistant",
				Parameters: map[string]interface{}{
					"text":      intent.Transcript,
					"device_id": intent.DeviceID,
				},
				Success: haResponse.Success,
			},
		},
		Metadata: map[string]interface{}{
			"ha_device_id": s.haConfig.DeviceID,
			"ha_base_url":  s.haConfig.BaseURL,
		},
	}
	
	if !haResponse.Success {
		response.Error = haResponse.ErrorMessage
		response.ErrorCode = "ha_processing_error"
	}
	
	return response, nil
}

// sendToHomeAssistant sends a voice request to Home Assistant
func (s *HomeAssistantSkill) sendToHomeAssistant(ctx context.Context, request HAVoiceRequest) (*HAVoiceResponse, error) {
	// Use the voice assistant API endpoint
	url := fmt.Sprintf("%s/api/voice/process", strings.TrimSuffix(s.haConfig.BaseURL, "/"))
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.haConfig.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return &HAVoiceResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("HA API error: %d - %s", resp.StatusCode, string(body)),
		}, nil
	}
	
	var haResponse HAVoiceResponse
	if err := json.Unmarshal(body, &haResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &haResponse, nil
}

// GetManifest returns the skill manifest
func (s *HomeAssistantSkill) GetManifest() (*skills.SkillManifest, error) {
	return &skills.SkillManifest{
		ID:          "com.loqalabs.homeassistant",
		Name:        "Home Assistant Voice Integration",
		Version:     "1.0.0",
		Description: "Integrates with Home Assistant Voice Preview Edition as a fallback skill",
		Author:      "Loqa Labs",
		License:     "AGPL-3.0",
		
		IntentPatterns: []skills.IntentPattern{
			{
				Name:       "ha_fallback",
				Examples:   []string{"*"}, // Catch-all pattern
				Confidence: 0.1, // Very low confidence to act as fallback
				Priority:   100, // Low priority - other skills should match first
				Enabled:    true,
				Categories: []string{"fallback"},
			},
		},
		
		Languages:   []string{"en"},
		Categories:  []string{"smart_home", "fallback", "integration"},
		
		Permissions: []skills.Permission{
			{
				Type:        skills.PermissionNetwork,
				Resource:    "*",
				Actions:     []string{"http_request"},
				Description: "Connect to Home Assistant API",
			},
			{
				Type:        skills.PermissionDeviceControl,
				Resource:    "*",
				Actions:     []string{"*"},
				Description: "Control devices via Home Assistant",
			},
			{
				Type:        skills.PermissionSpeaker,
				Resource:    "*",
				Actions:     []string{"play"},
				Description: "Play TTS responses from Home Assistant",
			},
		},
		
		ConfigSchema: &skills.ConfigSchema{
			Properties: map[string]skills.ConfigProperty{
				"base_url": {
					Type:        "string",
					Description: "Home Assistant base URL (e.g., http://homeassistant.local:8123)",
					Default:     "http://homeassistant.local:8123",
					Format:      "url",
				},
				"access_token": {
					Type:        "string",
					Description: "Home Assistant long-lived access token",
					Sensitive:   true,
				},
				"device_id": {
					Type:        "string",
					Description: "Device ID to use when communicating with HA",
					Default:     "loqa-voice-assistant",
				},
				"device_name": {
					Type:        "string",
					Description: "Device name displayed in Home Assistant",
					Default:     "Loqa Voice Assistant",
				},
				"mqtt_enabled": {
					Type:        "boolean",
					Description: "Enable MQTT integration (future feature)",
					Default:     false,
				},
				"mqtt_topic": {
					Type:        "string",
					Description: "MQTT topic for voice commands",
					Default:     "homeassistant/voice",
				},
				"timeout_seconds": {
					Type:        "integer",
					Description: "Request timeout in seconds",
					Default:     30,
				},
			},
			Required: []string{"base_url", "access_token"},
		},
		
		LoadOnStartup: true,
		Singleton:     true,
		Timeout:       "30s",
		SandboxMode:   skills.SandboxNone,
		TrustLevel:    skills.TrustSystem,
		
		Keywords: []string{"homeassistant", "ha", "smart home", "fallback", "integration"},
		Tags:     []string{"integration", "smart-home", "fallback"},
	}, nil
}

// GetStatus returns the current skill status
func (s *HomeAssistantSkill) GetStatus() skills.SkillStatus {
	return s.status
}

// GetConfig returns the skill configuration
func (s *HomeAssistantSkill) GetConfig() (*skills.SkillConfig, error) {
	return s.config, nil
}

// UpdateConfig updates the skill configuration
func (s *HomeAssistantSkill) UpdateConfig(ctx context.Context, config *skills.SkillConfig) error {
	s.config = config
	
	// Re-parse HA configuration
	if err := s.parseHAConfig(); err != nil {
		return fmt.Errorf("failed to parse updated HA config: %w", err)
	}
	
	// Update HTTP client timeout
	if s.haConfig.Timeout > 0 {
		s.httpClient.Timeout = time.Duration(s.haConfig.Timeout) * time.Second
	}
	
	// Test new connection
	if err := s.testConnection(ctx); err != nil {
		s.logger.Warn("Failed to test updated HA connection", zap.Error(err))
		s.status.LastError = fmt.Sprintf("HA connection test failed: %v", err)
	} else {
		s.status.LastError = ""
		s.status.Healthy = true
	}
	
	s.logger.Info("Updated Home Assistant skill configuration", 
		zap.String("skill_id", config.SkillID),
		zap.String("version", config.Version),
		zap.String("ha_url", s.haConfig.BaseURL))
	
	return nil
}

// HealthCheck performs a health check
func (s *HomeAssistantSkill) HealthCheck(ctx context.Context) error {
	if s.status.State != skills.SkillStateReady {
		return fmt.Errorf("skill not ready, current state: %s", s.status.State)
	}
	
	// Test connection to Home Assistant
	if err := s.testConnection(ctx); err != nil {
		s.status.Healthy = false
		s.status.LastError = fmt.Sprintf("HA health check failed: %v", err)
		return err
	}
	
	s.status.Healthy = true
	s.status.LastError = ""
	return nil
}

// Plugin entry point - required for Go plugins
var SkillPlugin HomeAssistantSkill