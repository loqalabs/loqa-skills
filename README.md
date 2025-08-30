[![Sponsor](https://img.shields.io/badge/Sponsor-Loqa-ff69b4?logo=githubsponsors&style=for-the-badge)](https://github.com/sponsors/annabarnes1138)
[![Ko-fi](https://img.shields.io/badge/Buy%20me%20a%20coffee-Ko--fi-FF5E5B?logo=ko-fi&logoColor=white&style=for-the-badge)](https://ko-fi.com/annabarnes)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL--3.0-blue?style=for-the-badge)](LICENSE)
[![Made with ❤️ by LoqaLabs](https://img.shields.io/badge/Made%20with%20%E2%9D%A4%EF%B8%8F-by%20LoqaLabs-ffb6c1?style=for-the-badge)](https://loqalabs.com)

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