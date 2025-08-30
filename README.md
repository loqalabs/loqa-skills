# 🧩 Loqa Skills

[![CI/CD Pipeline](https://github.com/loqalabs/loqa-skills/actions/workflows/ci.yml/badge.svg)](https://github.com/loqalabs/loqa-skills/actions/workflows/ci.yml)

Official and sample skills packaged as external services for the Loqa platform.

## Overview

Loqa Skills provides:
- Official skill implementations (Home Assistant, media control, etc.)
- Sample skills for developers
- Skill development framework and templates
- Testing and deployment tools

## Available Skills

### Official Skills
- **Home Assistant Bridge**: Control Home Assistant devices
- **Media Control**: Music playback and audio management
- **Timer & Reminders**: Voice-activated timers and reminders
- **Weather**: Local weather information and forecasts

### Sample Skills
- **Hello World**: Basic skill template
- **Device Control**: Simple device command examples
- **Custom Commands**: Advanced skill development patterns

## Features

- 🏠 **Smart Home Integration**: Ready-to-use Home Assistant connectivity
- 🎵 **Media Control**: Music and audio playback management
- ⏰ **Productivity**: Timers, reminders, and task management
- 🛠️ **Development Framework**: Tools and templates for custom skills
- 📦 **Easy Deployment**: Containerized skills with Docker support

## Skill Development

Each skill is an independent service that:
- Subscribes to relevant NATS subjects
- Processes voice commands and intents
- Executes actions (API calls, device control, etc.)
- Reports status and responses

## Architecture

Skills communicate with the Loqa ecosystem via:
- NATS message bus for commands and events
- gRPC for high-performance data exchange
- REST APIs for external service integration

## Getting Started

See the main [Loqa documentation](https://github.com/loqalabs/loqa) for setup and usage instructions.

## License

Licensed under the GNU Affero General Public License v3.0. See [LICENSE](LICENSE) for details.