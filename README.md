[![Sponsor](https://img.shields.io/badge/Sponsor-Loqa-ff69b4?logo=githubsponsors&style=for-the-badge)](https://github.com/sponsors/annabarnes1138)
[![Ko-fi](https://img.shields.io/badge/Buy%20me%20a%20coffee-Ko--fi-FF5E5B?logo=ko-fi&logoColor=white&style=for-the-badge)](https://ko-fi.com/annabarnes)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL--3.0-blue?style=for-the-badge)](LICENSE)
[![Made with ❤️ by Loqa Labs](https://img.shields.io/badge/Made%20with%20%E2%9D%A4%EF%B8%8F-by%20Loqa Labs-ffb6c1?style=for-the-badge)](https://loqalabs.com)

# 🧩 Loqa Skills

[![CI/CD Pipeline](https://github.com/loqalabs/loqa-skills/actions/workflows/ci.yml/badge.svg)](https://github.com/loqalabs/loqa-skills/actions/workflows/ci.yml)

Modular skill plugin system for the Loqa voice assistant platform.

## Overview

Loqa Skills implements a comprehensive plugin architecture introduced in Milestone 4a, providing:
- **Plugin Architecture**: Manifest-driven skill system with lifecycle hooks
- **Security Model**: Trust levels and sandboxing for safe skill execution  
- **Development Tools**: CLI tools and web UI for skill management
- **Official Skills**: Curated skills for common use cases
- **Sample Skills**: Examples and templates for developers

## Available Skills

### Official Skills
- **Home Assistant Bridge**: Control Home Assistant devices
- **Media Control**: Music playback and audio management
- **Timer & Reminders**: Voice-activated timers and reminders
- **Weather**: Local weather information and forecasts

### Sample Skills
- **Example Skill**: Comprehensive skill template demonstrating all features

## Features

### 🆕 Milestone 4a: Modular Plugin Architecture

- 🧩 **SkillPlugin Interface**: Comprehensive lifecycle hooks (Initialize, Teardown, CanHandle, HandleIntent)
- 📋 **Skill Manifests**: JSON-based configuration with permissions, intents, and metadata
- 🔄 **Dynamic Loading**: Runtime skill loading, unloading, and reloading
- 🛡️ **Security Model**: Trust levels (system, verified, community, unknown) and sandbox modes
- 🎛️ **Management Tools**: CLI tool and web UI for skill administration
- 🌐 **REST API**: Complete skill management via `/api/skills` endpoints
- 🔧 **Multi-Format Support**: Go plugins, process-based skills, future WASM support

### Core Capabilities

- 🏠 **Smart Home Integration**: Ready-to-use Home Assistant connectivity
- 🎵 **Media Control**: Music and audio playbook management
- ⏰ **Productivity**: Timers, reminders, and task management
- 🛠️ **Development Framework**: Tools and templates for custom skills
- 📦 **Easy Deployment**: Containerized skills with Docker support

## Skill Development

### Skill Manifest (`skill.json`)

Each skill requires a manifest file defining:
- **Basic Info**: ID, name, version, author, license
- **Intent Patterns**: Voice commands the skill handles with examples
- **Permissions**: Required system access (microphone, network, etc.)
- **Configuration**: Schema for user-configurable settings
- **Security**: Trust level and sandbox mode requirements

### SkillPlugin Interface

Skills implement the following lifecycle:
- `Initialize(ctx, config)` - Setup with configuration
- `CanHandle(intent)` - Determine if skill handles the intent
- `HandleIntent(ctx, request)` - Process the voice command
- `GetManifest()` - Return skill metadata
- `UpdateConfig(ctx, config)` - Handle configuration changes
- `Teardown(ctx)` - Clean shutdown

### Skill Types

- **Go Plugins**: Compiled `.so` files loaded dynamically
- **Process Skills**: Executable binaries run as separate processes
- **WASM Skills**: WebAssembly modules (future)

## Architecture

The plugin system provides:
- **Dynamic Loading**: Skills loaded/unloaded at runtime via CLI or web UI
- **Sandboxing**: Process isolation and permission enforcement
- **Intent Routing**: Automatic routing based on confidence and priority
- **Error Handling**: Graceful fallback and recovery

## Getting Started

See the main [Loqa documentation](https://github.com/loqalabs/loqa) for setup and usage instructions.

## License

Licensed under the GNU Affero General Public License v3.0. See [LICENSE](LICENSE) for details.