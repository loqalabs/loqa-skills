# Home Assistant Voice Integration Guide

This guide explains how to set up Loqa as a Home Assistant Voice Preview Edition integration, allowing you to use Loqa's local voice processing with your existing Home Assistant smart home setup.

## Overview

Loqa's Home Assistant integration provides:

- **Local Speech-to-Text**: Voice processing happens on your local Loqa hub
- **Multi-device Support**: Connect multiple voice relay devices throughout your home
- **Fallback Integration**: Forward unhandled commands to Home Assistant
- **Full Observability**: Track all voice interactions in the Loqa timeline
- **Simple Setup**: Works with existing HA Voice Preview Edition setup

## Architecture

```
Voice Relay → Loqa Hub → [Local STT] → [Skills System] → Home Assistant
     ↓              ↓                       ↓                ↓
  Audio Stream   Process Audio        Route Commands    Execute Actions
                                         ↓
                                   [Loqa Commander UI]
```

## Prerequisites

### Home Assistant Requirements

1. **Home Assistant Core** version 2023.9 or later
2. **Voice Preview Edition** enabled
3. **Long-lived Access Token** created
4. **HTTP API** accessible from Loqa hub

### Loqa Requirements

1. **Loqa Hub** running with skills system enabled
2. **Voice Relays** configured and connected
3. Network connectivity between Loqa hub and Home Assistant

## Setup Instructions

### Step 1: Create Home Assistant Access Token

1. Log into your Home Assistant web interface
2. Navigate to **Profile** → **Security** → **Long-Lived Access Tokens**
3. Click **Create Token**
4. Enter a name like "Loqa Voice Integration"
5. Copy the generated token (you'll need this later)

### Step 2: Enable Home Assistant Voice Preview

1. Go to **Settings** → **Voice assistants**
2. Enable **Voice Preview Edition** if not already enabled
3. Configure your preferred STT/TTS providers (optional - Loqa handles STT)
4. Test that the voice system works with Home Assistant

### Step 3: Configure Loqa Home Assistant Skill

#### Option A: Using the Skills API

```bash
# Configure the Home Assistant skill
curl -X PUT http://your-loqa-hub:3000/api/skills/builtin.homeassistant \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "base_url": "http://your-homeassistant:8123",
      "access_token": "your-long-lived-access-token",
      "device_id": "loqa-voice-assistant",
      "device_name": "Loqa Voice Assistant",
      "timeout_seconds": 30
    }
  }'

# Enable the skill
curl -X POST http://your-loqa-hub:3000/api/skills/builtin.homeassistant/enable
```

#### Option B: Using Configuration Files

Create a skill configuration file:

```bash
# Create config directory if it doesn't exist
mkdir -p /path/to/loqa-hub/config/skills

# Create the HA skill configuration
cat > /path/to/loqa-hub/config/skills/builtin.homeassistant.json << EOF
{
  "skill_id": "builtin.homeassistant",
  "name": "Home Assistant Integration",
  "version": "1.0.0",
  "config": {
    "base_url": "http://your-homeassistant:8123",
    "access_token": "your-long-lived-access-token",
    "device_id": "loqa-voice-assistant",
    "device_name": "Loqa Voice Assistant",
    "mqtt_enabled": false,
    "mqtt_topic": "homeassistant/voice",
    "timeout_seconds": 30
  },
  "permissions": [
    {
      "type": "network",
      "resource": "*",
      "actions": ["http_request"],
      "description": "Connect to Home Assistant API"
    }
  ],
  "enabled": true,
  "timeout": "30s",
  "max_retries": 3
}
EOF
```

### Step 4: Configure Network Access

Ensure your Loqa hub can reach Home Assistant:

```bash
# Test connectivity from Loqa hub
curl -H "Authorization: Bearer your-long-lived-access-token" \
     http://your-homeassistant:8123/api/

# Expected response: {"message": "API running."}
```

### Step 5: Test the Integration

1. **Check Skill Status**:
   ```bash
   curl http://your-loqa-hub:3000/api/skills/builtin.homeassistant
   ```

2. **Test Voice Command**:
   - Say "Hey Loqa, turn on the living room lights" to your voice relay
   - Check the Loqa Commander timeline at `http://your-loqa-hub:5173`
   - Verify the command was forwarded to Home Assistant

3. **Verify Home Assistant Logs**:
   - Check HA logs for incoming API requests from Loqa
   - Look for device actions being executed

## Configuration Reference

### Required Configuration

| Parameter | Description | Example |
|-----------|-------------|---------|
| `base_url` | Home Assistant base URL | `http://homeassistant.local:8123` |
| `access_token` | Long-lived access token | `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...` |

### Optional Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `device_id` | Device ID for HA API calls | `loqa-voice-assistant` |
| `device_name` | Display name in HA | `Loqa Voice Assistant` |
| `timeout_seconds` | API request timeout | `30` |
| `mqtt_enabled` | Enable MQTT integration | `false` |
| `mqtt_topic` | MQTT topic for commands | `homeassistant/voice` |

## Troubleshooting

### Common Issues

#### "Failed to connect to Home Assistant"

**Symptoms**: Skill status shows connection errors

**Solutions**:
1. Verify Home Assistant is running and accessible
2. Check network connectivity between Loqa and HA
3. Ensure the base_url is correct (include http/https)
4. Verify firewall allows connections on port 8123

#### "Authentication failed"

**Symptoms**: HTTP 401 errors in Loqa logs

**Solutions**:
1. Verify the access token is correct and active
2. Check token permissions in HA settings
3. Regenerate the access token if needed

#### "Commands not reaching Home Assistant"

**Symptoms**: Loqa processes voice but HA doesn't respond

**Solutions**:
1. Check if the HA skill is enabled and healthy
2. Verify the skill priority (HA skill should be low priority/fallback)
3. Enable debug logging to see API requests
4. Test the HA Voice Preview Edition directly

#### "Voice commands timing out"

**Symptoms**: Slow responses or timeout errors

**Solutions**:
1. Increase `timeout_seconds` in skill config
2. Check network latency between Loqa and HA
3. Verify Home Assistant system performance

### Debug Logging

Enable detailed logging for troubleshooting:

```bash
# Set environment variables for debug logging
export LOG_LEVEL=debug
export HA_SKILL_DEBUG=true

# Restart Loqa hub
systemctl restart loqa-hub
```

Check logs:
```bash
# View Loqa hub logs
journalctl -u loqa-hub -f

# View Home Assistant logs
docker logs -f homeassistant
```

### Health Checks

Monitor integration health:

```bash
# Check overall skill status
curl http://your-loqa-hub:3000/api/skills

# Check specific HA skill health
curl http://your-loqa-hub:3000/api/skills/builtin.homeassistant

# View recent voice events
curl http://your-loqa-hub:3000/api/voice-events?limit=10
```

## Advanced Configuration

### Simple Mode Fallback

For constrained environments, you can use simple mode instead of the full skills system:

```json
{
  "simple_mode": {
    "enabled": true,
    "home_assistant": {
      "base_url": "http://your-homeassistant:8123",
      "access_token": "your-token",
      "device_id": "loqa-simple",
      "timeout": 30
    }
  }
}
```

### MQTT Integration (Future)

The skill supports MQTT for future Home Assistant integrations:

```json
{
  "config": {
    "mqtt_enabled": true,
    "mqtt_topic": "homeassistant/voice/loqa",
    "mqtt_host": "your-mqtt-broker",
    "mqtt_port": 1883
  }
}
```

### Multiple Home Assistant Instances

You can configure multiple HA instances by creating additional skill configurations:

```bash
# Load a second HA skill for a different location
curl -X POST http://your-loqa-hub:3000/api/skills \
  -d '{"skill_path": "/path/to/ha-office-skill"}'
```

## Security Considerations

1. **Access Token Security**: Store tokens securely and rotate regularly
2. **Network Isolation**: Use VLANs or firewalls to isolate smart home traffic
3. **HTTPS**: Use HTTPS for Home Assistant if accessible over the internet
4. **Monitoring**: Monitor API usage and watch for unusual patterns

## Performance Optimization

1. **Local Processing**: Most voice processing happens on Loqa, reducing HA load
2. **Caching**: Consider caching HA state locally for faster responses  
3. **Parallel Processing**: Multiple relay devices can send commands simultaneously
4. **Timeout Tuning**: Adjust timeouts based on your network and HA performance

## Integration with Other Skills

The HA skill acts as a fallback, so other skills take priority:

1. Built-in skills (lights, etc.) handle specific commands first
2. Custom skills can override HA integration for specific intents
3. HA skill handles anything not matched by other skills
4. Skill priority can be configured per skill

## Support and Contributing

- **Issues**: Report issues at [github.com/loqalabs/loqa/issues]
- **Documentation**: Contribute improvements to this guide
- **Feature Requests**: Suggest new HA integration features

## Changelog

- **v1.0.0**: Initial HA Voice Preview Edition integration
- **v1.0.1**: Added simple mode fallback support
- **v1.1.0**: Enhanced observability and debug tooling