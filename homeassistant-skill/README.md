# Home Assistant Voice Skill

A Loqa skill that integrates with Home Assistant Voice Preview Edition, allowing Loqa to act as a local voice processing frontend for your Home Assistant smart home.

## Features

- **Local Voice Processing**: All STT happens on your Loqa hub for privacy
- **Multi-device Support**: Handle voice commands from multiple Loqa pucks
- **Fallback Integration**: Forwards unhandled commands to Home Assistant
- **Full Observability**: Track all voice interactions in the Loqa timeline
- **Easy Configuration**: Simple setup with HA access token

## Installation

### Prerequisites

- Running Loqa Hub with skills system enabled
- Home Assistant instance with Voice Preview Edition enabled
- Home Assistant long-lived access token

### Build and Install

```bash
# Build the plugin
make build

# Install to Loqa hub
make install

# Load the skill via API
make load

# Configure the skill (see Configuration section)
# ...

# Enable the skill
make enable
```

### Manual Installation

```bash
# Build the plugin
go build -buildmode=plugin -o homeassistant-skill.so main.go

# Copy to Loqa hub skills directory
mkdir -p /path/to/loqa-hub/skills/homeassistant-skill
cp homeassistant-skill.so /path/to/loqa-hub/skills/homeassistant-skill/
cp skill.json /path/to/loqa-hub/skills/homeassistant-skill/

# Load via Loqa CLI
cd /path/to/loqa-hub
go run ./cmd/skills-cli --action load --path ./skills/homeassistant-skill
```

## Configuration

### 1. Create Home Assistant Access Token

1. In Home Assistant, go to Profile → Security
2. Create a "Long-Lived Access Token"
3. Name it "Loqa Voice Integration"
4. Copy the generated token

### 2. Configure the Skill

```bash
# Configure via API
curl -X PUT http://localhost:3000/api/skills/com.loqalabs.homeassistant \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "base_url": "http://homeassistant.local:8123",
      "access_token": "your-long-lived-access-token",
      "device_id": "loqa-voice-assistant",
      "device_name": "Loqa Voice Assistant",
      "timeout_seconds": 30
    }
  }'

# Or create a config file
mkdir -p /path/to/loqa-hub/config/skills
cat > /path/to/loqa-hub/config/skills/com.loqalabs.homeassistant.json << EOF
{
  "skill_id": "com.loqalabs.homeassistant",
  "name": "Home Assistant Voice Integration",
  "version": "1.0.0",
  "config": {
    "base_url": "http://homeassistant.local:8123",
    "access_token": "your-long-lived-access-token",
    "device_id": "loqa-voice-assistant",
    "device_name": "Loqa Voice Assistant",
    "timeout_seconds": 30
  },
  "enabled": true
}
EOF
```

### Configuration Options

| Parameter | Required | Description | Default |
|-----------|----------|-------------|---------|
| `base_url` | Yes | Home Assistant base URL | `http://homeassistant.local:8123` |
| `access_token` | Yes | Long-lived access token | - |
| `device_id` | No | Device ID for HA API | `loqa-voice-assistant` |
| `device_name` | No | Device name in HA | `Loqa Voice Assistant` |
| `timeout_seconds` | No | Request timeout | `30` |
| `mqtt_enabled` | No | Enable MQTT (future) | `false` |
| `mqtt_topic` | No | MQTT topic | `homeassistant/voice` |

## Usage

Once configured and enabled, the skill will:

1. **Act as Fallback**: Handle voice commands not processed by other skills
2. **Forward to HA**: Send commands to Home Assistant Voice API
3. **Return Responses**: Relay HA's text responses back to the user
4. **Track Events**: Log all interactions in the Loqa timeline

### Example Voice Commands

- "Turn on the living room lights"
- "Set the thermostat to 72 degrees"
- "Lock the front door"
- "What's the temperature in the bedroom?"
- "Play music in the kitchen"

## Development

### Building

```bash
# Development build with debug symbols
make debug

# Production build with optimizations
make release

# Test the code
make test

# Update dependencies
make deps
```

### Plugin Architecture

The skill implements the Loqa `SkillPlugin` interface:

- `Initialize()` - Sets up HA connection and validates config
- `CanHandle()` - Always returns true (acts as fallback)
- `HandleIntent()` - Forwards commands to HA and processes responses
- `GetManifest()` - Returns skill metadata and configuration schema
- `UpdateConfig()` - Handles configuration updates
- `HealthCheck()` - Tests HA connectivity

### Testing

```bash
# Check skill status
make status

# View recent skill usage
make logs

# Test a command (requires skill to be loaded)
curl -X POST http://localhost:3000/api/voice-events \
  -H "Content-Type: application/json" \
  -d '{
    "puck_id": "test-puck",
    "transcription": "turn on the lights",
    "intent": "lighting_control"
  }'
```

## Troubleshooting

### Skill Not Loading

```bash
# Check skill status
curl http://localhost:3000/api/skills/com.loqalabs.homeassistant

# View Loqa hub logs
journalctl -u loqa-hub -f

# Check plugin build
ldd homeassistant-skill.so
```

### Connection Issues

```bash
# Test HA connectivity
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://homeassistant.local:8123/api/

# Check configuration
curl http://localhost:3000/api/skills/com.loqalabs.homeassistant/config
```

### Voice Commands Not Working

1. Verify skill is enabled and healthy
2. Check that other skills aren't intercepting commands
3. Ensure HA Voice Preview Edition is working
4. Review voice event timeline for errors

## Security

- **Access Token**: Store securely and rotate regularly
- **Network**: Use HTTPS for HA if internet-accessible
- **Firewall**: Restrict access to HA API port (8123)
- **Monitoring**: Watch for unusual API usage patterns

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This skill is licensed under the GNU Affero General Public License v3.0. See [LICENSE](../../LICENSE) for details.