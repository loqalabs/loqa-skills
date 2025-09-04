module github.com/loqalabs/loqa-skills/homeassistant-skill

go 1.23.0

require (
	github.com/loqalabs/loqa-hub v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
)

replace github.com/loqalabs/loqa-hub => ../../loqa-hub

require go.uber.org/multierr v1.10.0 // indirect
