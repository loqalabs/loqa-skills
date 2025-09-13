# Development Guidance

See [/loqalabs/CLAUDE.md](../CLAUDE.md) for complete development workflow guidance.

## Service Context

**loqa-skills** - Modular skill plugin system with manifest-driven architecture (Go)

- **Role**: Plugin system for voice command handling (Home Assistant, custom skills)
- **Quality Gates**: `make quality-check` (includes go fmt, go vet, plugin validation, hub integration)
- **Development**: Individual skills in subdirectories, `make build && make install` per skill
- **Management**: Skills CLI in loqa-hub (`go run ./cmd/skills-cli --help`)

All workflow rules and development guidance are provided automatically by the MCP server based on repository detection.