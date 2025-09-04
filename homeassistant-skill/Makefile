# Home Assistant Skill Makefile

SKILL_NAME = homeassistant-skill
PLUGIN_FILE = $(SKILL_NAME).so

.PHONY: build clean test install

# Build the skill plugin
build: $(PLUGIN_FILE)

$(PLUGIN_FILE): main.go go.mod go.sum
	go build -buildmode=plugin -o $(PLUGIN_FILE) main.go

# Clean build artifacts
clean:
	rm -f $(PLUGIN_FILE)
	go clean

# Test the skill
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build for production
release: clean
	CGO_ENABLED=1 go build -buildmode=plugin -ldflags="-s -w" -o $(PLUGIN_FILE) main.go
	strip $(PLUGIN_FILE)

# Development build with debug symbols
debug:
	go build -buildmode=plugin -gcflags="-N -l" -o $(PLUGIN_FILE) main.go

# Install the skill to Loqa hub
install: build
	@echo "Installing $(SKILL_NAME) plugin..."
	@mkdir -p ../../loqa-hub/skills/$(SKILL_NAME)
	@cp $(PLUGIN_FILE) ../../loqa-hub/skills/$(SKILL_NAME)/
	@cp skill.json ../../loqa-hub/skills/$(SKILL_NAME)/
	@echo "Skill installed successfully"

# Load the skill via API (requires running Loqa hub)
load:
	curl -X POST http://localhost:3000/api/skills \
		-H "Content-Type: application/json" \
		-d '{"skill_path": "./skills/$(SKILL_NAME)"}'

# Enable the skill
enable:
	curl -X POST http://localhost:3000/api/skills/com.loqalabs.homeassistant/enable

# Check skill status  
status:
	curl http://localhost:3000/api/skills/com.loqalabs.homeassistant

# View skill logs
logs:
	curl http://localhost:3000/api/voice-events?skill_id=com.loqalabs.homeassistant&limit=10