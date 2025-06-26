# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

**Target**: Analyze Reddit to retrieve latest trends and user needs for indiehacker/solo developer business ideation. The system identifies trends and user needs in underdeveloped markets with low competition but high demand.

**Important**: The system should NOT generate business ideas. Instead, it should analyze trends and clearly explain user needs with providing links to relevant subreddits.

**Scope**: This is a PoC (Proof of Concept) - no hosting or CLI delivery planned initially. Run directly from VS Code using main.go. If successful, delivery method will be considered later.

## Tech Stack & Principles

**Language**: Go - chosen for the entire project
**AI Agent Library**: Custom-built simple AI Agent library (avoiding existing Go AI libs)
**Code Principles**: Clean code practices for Go development

## Planned Architecture

**Multi-Agent Workflow**: 
- **Master Agent**: Coordinates overall workflow and sub-agent execution
- **Sub-Agents**: Specialized agents for specific analysis tasks
- **Inspiration**: Architecture based on Anthropic's multi-agent research system (https://www.anthropic.com/engineering/built-multi-agent-research-system)

**Key Components to Build**:
- Custom AI Agent execution framework
- Reddit data retrieval and analysis system
- Trend identification and user need extraction logic
- Multi-agent coordination system

## Common Commands

### Development
- **Build**: `make build` (outputs to out/redditanalyzer)
- **Lint**: `make lint` (uses go vet)
- **Run**: `go run ./cmd/redditanalyzer`
- **Test**: `go test ./...`
- **Dependencies**: `go mod tidy` to install/clean dependencies, then `go mod vendor` to update vendor directory

### VS Code Debug
- Use F5 or "Launch Reddit Analyzer" configuration to debug
- Program entry point: `cmd/redditanalyzer/main.go`

## Current Architecture

### Project Structure
- **Entry Point**: `cmd/redditanalyzer/main.go` - Main application with logrus logging setup
- **Module**: `reddit-analyzer` using Go 1.24.4
- **Dependencies**: Uses `github.com/sirupsen/logrus` for structured logging

### Current State
The project is in early development phase with basic logging infrastructure. The `cmd/` directory structure follows Go project layout conventions, preparing for the multi-agent architecture implementation.

### Key Dependencies
- **logrus**: Structured logging with colored output and timestamps configured in main.go

### Future Architecture (To Be Implemented)
The system will evolve into a multi-agent workflow where:
- Master agent orchestrates Reddit analysis workflow
- Sub-agents handle specialized tasks (data retrieval, trend analysis, user need extraction)
- Custom AI agent framework manages agent coordination and execution
- Results focus on trend analysis and user needs (not business idea generation)

## Development Notes
- Project follows standard Go module structure
- No tests, documentation, or build automation currently exist
- Git repository on main branch with pending changes to go.mod and new cmd/ directory