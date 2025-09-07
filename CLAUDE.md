# CLAUDE.md - Loqa Skills System

This file provides Claude Code with specific guidance for working with the Loqa Skills System - the extensible plugin framework for voice capabilities and AI features.

## Service Overview

Loqa Skills System provides:
- **Plugin Architecture**: Go plugin system for extensible voice capabilities
- **Skill Management**: Loading, configuration, and lifecycle management of skills
- **Intent Routing**: Routes parsed voice commands to appropriate skills
- **Skill Framework**: Base interfaces and utilities for skill development
- **Manifest System**: Declarative skill configuration and metadata
- **Sandboxing**: Safe execution environment for third-party skills

## Architecture Role

- **Service Type**: Plugin system and skills framework (Go)
- **Dependencies**: loqa-proto (skill communication protocols)
- **Used By**: loqa-hub (skill execution and management)
- **Plugin Type**: Go plugins (.so files) with standardized interfaces
- **Configuration**: Manifest-driven skill registration and setup

## Development Commands

### Skill Development
```bash
# Create new skill from template
make new-skill NAME=my-new-skill

# Build a skill plugin
cd my-skill/
make build

# Install skill to hub
make install

# Load skill via CLI (from hub directory)
go run ./cmd/skills-cli --action load --path ../loqa-skills/my-skill

# Test skill locally
make test
```

### Skills Management
```bash
# List available skills in repository
ls -la */

# Build all skills
make build-all

# Install all skills
make install-all

# Clean build artifacts
make clean
```

### Quality & Testing
```bash
# Test specific skill
cd homeassistant-skill/
go test ./...

# Lint skill code
make lint

# Validate skill manifest
make validate-manifest

# Integration test with hub
make test-integration
```

## Skill Development

### Skill Structure
```bash
my-skill/
├── main.go              # Plugin entry point
├── skill.go             # Skill implementation
├── manifest.json        # Skill metadata and configuration
├── config.json          # Default configuration
├── go.mod               # Go module definition
├── Makefile            # Build and install commands
├── README.md           # Skill documentation
└── tests/              # Skill-specific tests
    ├── unit/
    └── integration/
```

### Skill Interface
```go
// All skills must implement the SkillPlugin interface
type SkillPlugin interface {
    // Initialize the skill with configuration
    Initialize(config map[string]interface{}) error
    
    // Execute a command/intent
    Execute(ctx context.Context, intent *Intent) (*Response, error)
    
    // Get skill metadata
    GetManifest() *Manifest
    
    // Health check
    Health() error
    
    // Cleanup resources
    Cleanup() error
}
```

### Skill Manifest
```json
{
  "id": "my-skill",
  "name": "My Custom Skill",
  "version": "1.0.0",
  "description": "Description of what this skill does",
  "author": "Your Name",
  "license": "AGPL-3.0",
  "intents": [
    {
      "name": "my_intent",
      "patterns": [
        "do something",
        "perform action",
        "execute command"
      ],
      "parameters": ["target", "value"]
    }
  ],
  "configuration": {
    "api_url": {
      "type": "string",
      "required": true,
      "description": "API endpoint URL"
    },
    "timeout": {
      "type": "integer",
      "default": 30,
      "description": "Request timeout in seconds"
    }
  },
  "permissions": [
    "network.http",
    "device.control"
  ]
}
```

## Example Skills

### Home Assistant Skill
```bash
cd homeassistant-skill/

# Configure Home Assistant connection
cp config.example.json config.json
# Edit config.json with your Home Assistant URL and token

# Build and install
make build
make install

# Load into hub
go run ../loqa-hub/cmd/skills-cli --action load --path $(pwd)

# Test commands
echo "turn on the living room lights" | go run ../loqa-hub/cmd --test-mode
```

### Example Skill Template
```bash
cd example-skill/

# This is a template skill showing:
- Basic skill structure
- Intent handling
- Configuration management
- Error handling
- Testing patterns

# Use as starting point for new skills
cp -r example-skill/ ../my-new-skill/
cd ../my-new-skill/
# Modify manifest.json, skill.go, etc.
```

## Skill Configuration

### Runtime Configuration
```json
{
  "skill_id": "homeassistant",
  "name": "Home Assistant Integration",
  "version": "1.2.0",
  "config": {
    "base_url": "http://homeassistant:8123",
    "access_token": "your-long-lived-access-token",
    "timeout": 30,
    "entities": {
      "lights": ["light.living_room", "light.bedroom"],
      "switches": ["switch.coffee_maker"]
    }
  },
  "enabled": true,
  "timeout": "30s",
  "max_retries": 3
}
```

### Environment-Based Configuration
```bash
# Home Assistant skill configuration
export HASS_URL=http://homeassistant:8123
export HASS_TOKEN=your-token

# Weather skill configuration  
export WEATHER_API_KEY=your-api-key
export WEATHER_DEFAULT_LOCATION="San Francisco"

# Media skill configuration
export SPOTIFY_CLIENT_ID=your-client-id
export SPOTIFY_CLIENT_SECRET=your-client-secret
```

## Intent System

### Intent Definition
```go
type Intent struct {
    Name       string                 `json:"name"`
    Text       string                 `json:"text"`
    Confidence float64               `json:"confidence"`
    Entities   map[string]interface{} `json:"entities"`
    Context    map[string]interface{} `json:"context"`
}
```

### Intent Patterns
```json
{
  "intents": [
    {
      "name": "turn_on_light",
      "patterns": [
        "turn on {light}",
        "switch on the {light}",
        "turn {light} on",
        "enable {light}"
      ],
      "parameters": ["light"]
    },
    {
      "name": "set_brightness",
      "patterns": [
        "set {light} to {brightness}%",
        "dim {light} to {brightness}",
        "make {light} {brightness} percent bright"
      ],
      "parameters": ["light", "brightness"]
    }
  ]
}
```

## Skill Testing

### Unit Testing
```go
func TestSkillExecution(t *testing.T) {
    skill := &MySkill{}
    
    // Initialize with test config
    config := map[string]interface{}{
        "api_url": "http://localhost:8080",
        "timeout": 10,
    }
    err := skill.Initialize(config)
    assert.NoError(t, err)
    
    // Test intent execution
    intent := &Intent{
        Name: "test_intent",
        Text: "test command",
        Entities: map[string]interface{}{
            "target": "test_device",
        },
    }
    
    response, err := skill.Execute(context.Background(), intent)
    assert.NoError(t, err)
    assert.Equal(t, "Command executed successfully", response.Text)
}
```

### Integration Testing
```bash
# Test skill with running hub
cd homeassistant-skill/
make test-integration

# Test skill loading
go run ../loqa-hub/cmd/skills-cli --action load --path $(pwd)
go run ../loqa-hub/cmd/skills-cli --action test --skill homeassistant

# Test intent processing
echo "turn on kitchen light" | go run ../loqa-hub/cmd --test-mode
```

### Mock Testing
```bash
# Use mock implementations for external dependencies
# tests/mocks/homeassistant_mock.go
type MockHomeAssistant struct{}
func (m *MockHomeAssistant) CallService(domain, service string, data map[string]interface{}) error {
    return nil // Mock implementation
}
```

## Advanced Features

### Skill Chaining
```go
// Skills can invoke other skills
func (s *MySkill) Execute(ctx context.Context, intent *Intent) (*Response, error) {
    // Execute primary action
    result1, err := s.executePrimaryAction(intent)
    if err != nil {
        return nil, err
    }
    
    // Chain to another skill if needed
    if intent.Context["chain_to"] != nil {
        chainIntent := &Intent{
            Name: intent.Context["chain_to"].(string),
            // ... populate intent
        }
        return s.skillManager.ExecuteIntent(ctx, chainIntent)
    }
    
    return result1, nil
}
```

### Async Skill Execution
```go
func (s *MySkill) Execute(ctx context.Context, intent *Intent) (*Response, error) {
    // For long-running operations, return immediate response
    if s.isLongRunningIntent(intent) {
        go func() {
            // Execute asynchronously
            result := s.executeLongRunningAction(intent)
            // Send result via callback or event
            s.publishResult(intent.ID, result)
        }()
        
        return &Response{
            Text: "Operation started, I'll let you know when it's complete",
            Async: true,
        }, nil
    }
    
    // Synchronous execution for quick operations
    return s.executeSync(intent)
}
```

## Skills CLI Integration

### Hub Skills CLI
```bash
# From loqa-hub directory
go run ./cmd/skills-cli --help

# List loaded skills
go run ./cmd/skills-cli --action list

# Get detailed skill information
go run ./cmd/skills-cli --action info --skill homeassistant

# Enable/disable skills
go run ./cmd/skills-cli --action enable --skill homeassistant
go run ./cmd/skills-cli --action disable --skill homeassistant

# Reload skill (for development)
go run ./cmd/skills-cli --action reload --skill homeassistant
```

### Skill Management API
```bash
# Via HTTP API (when hub is running)
curl http://localhost:3000/api/skills
curl -X POST http://localhost:3000/api/skills -d '{"skill_path": "/path/to/skill"}'
curl -X POST http://localhost:3000/api/skills/homeassistant/enable
curl -X DELETE http://localhost:3000/api/skills/homeassistant
```

## Common Tasks

### Creating a New Skill
```bash
# 1. Copy example skill template
cp -r example-skill/ my-weather-skill/
cd my-weather-skill/

# 2. Update manifest.json
{
  "id": "weather",
  "name": "Weather Information",
  "intents": [
    {
      "name": "get_weather",
      "patterns": ["what's the weather", "weather forecast"]
    }
  ]
}

# 3. Implement skill logic in skill.go
func (s *WeatherSkill) Execute(ctx context.Context, intent *Intent) (*Response, error) {
    weather := s.getWeatherData(intent.Entities["location"])
    return &Response{
        Text: fmt.Sprintf("The weather is %s with a temperature of %d°F", 
                         weather.Condition, weather.Temperature),
    }, nil
}

# 4. Build and test
make build
make test
make install
```

### Debugging Skills
```bash
# Enable debug logging for skills
export LOG_LEVEL=debug

# Test skill loading
go run ../loqa-hub/cmd/skills-cli --action load --path $(pwd) --verbose

# Debug intent matching
go run ../loqa-hub/cmd --test-mode --debug-intents

# Monitor skill execution
tail -f ../loqa-hub/logs/skills.log
```

### Skill Deployment
```bash
# Build for production
make build GOOS=linux GOARCH=amd64

# Package skill
make package

# Deploy to skill repository
make deploy REPO_URL=https://skills.loqalabs.com
```

## Related Documentation

- **Master Documentation**: `../loqa/config/CLAUDE.md` - Full ecosystem overview
- **Hub Integration**: `../loqa-hub/CLAUDE.md` - Skills management and execution
- **Protocol Definitions**: `../loqa-proto/CLAUDE.md` - Skills communication protocols
