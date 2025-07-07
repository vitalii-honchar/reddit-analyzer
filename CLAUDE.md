# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

**Target**: Intelligent Reddit Research Agent for indie hackers to discover hidden market opportunities. The system analyzes Reddit to identify underdeveloped markets with low competition but high demand based on user's project direction.

**Key Innovation**: The AI Agent intelligently selects niche subreddits based on user input, avoiding obvious entrepreneurship communities (r/entrepreneur, r/startups) to find genuine, unserved problems in specific domains.

**User Experience**: 
- User provides project direction (e.g., "cybersecurity project", "company public speaking logging")
- AI Agent autonomously selects relevant niche subreddits
- System analyzes fresh posts (last 1 week) for hidden opportunities
- Outputs ranked opportunities with analysis, avoiding business idea generation

**Important**: The system should NOT generate business ideas. Instead, it should analyze trends and clearly explain user needs with providing links to relevant subreddits.

**Scope**: This is a PoC (Proof of Concept) - no hosting or CLI delivery planned initially. Run directly from VS Code using main.go. If successful, delivery method will be considered later.

## Tech Stack & Principles

**Language**: Go - chosen for the entire project
**AI Agent Library**: Uses external go-agent library (github.com/vitalii-honchar/go-agent) for agent orchestration and tool management
**LLM Integration**: OpenAI Go SDK for LLM tool implementation and API communication
**Authentication**: OpenAI API key stored in .env file and accessed via OPENAI_API_KEY environment variable
**Code Principles**: Clean code practices for Go development

## Planned Architecture

**Intelligent Single-Agent Workflow**: 
- **Master Agent**: "Reddit Indie Hacker Research Agent" orchestrates the entire analysis workflow
- **Tools**: Specialized tools for each analysis phase
- **Inspiration**: Architecture based on Anthropic's multi-agent research system (https://www.anthropic.com/engineering/built-multi-agent-research-system)

**Key Tools to Build**:
1. **SubredditSelectorTool**: AI intelligently selects niche subreddits based on user's project direction
2. **TopSubredditPostsTool**: Fetches recent posts from selected subreddits using Reddit API
3. **FilterPosts**: Filters posts by engagement metrics and 1-week timeframe
4. **EvaluatePost**: Uses LLM to rate posts 1-5 for hidden indie hacker opportunities

**Intelligence Focus**:
- Avoid obvious entrepreneurship subreddits (r/entrepreneur, r/startups)
- Target niche communities where problems exist but solutions are underdeveloped
- Focus on genuine pain points mentioned by domain professionals

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
- **go-reddit**: Reddit API client for fetching posts without authentication
- **go-agent**: External AI agent library for orchestration and tool management
- **openai-go**: OpenAI Go SDK for LLM API communication and tool implementations

### Implementation Plan
Refined development roadmap with detailed workflow analysis:

**Completed Tasks** ✅
- Analyze go-agent library structure and capabilities
- Create comprehensive implementation plan for Reddit analysis system
- Update CLAUDE.md with refined project description and implementation plan

**Phase 0: Foundation Setup** (Critical)
1. ⏳ Add go-agent dependency to project (go mod tidy)
2. ⏳ Set up environment configuration (.env file with OPENAI_API_KEY)
3. ⏳ Define core data structures (RedditPost, OpportunityAnalysis, AgentResult)

**Phase 1: Core Tools Implementation** (High Priority)
4. ⏳ Implement SubredditSelectorTool - AI intelligently selects 3-5 niche subreddits
5. ⏳ Implement TopSubredditPostsTool - fetches posts from multiple subreddits (last 7 days)
6. ⏳ Implement FilterPostsTool - filters by engagement metrics (min upvotes/comments)
7. ⏳ Implement EvaluatePostTool - LLM rates posts 1-5 for hidden indie opportunities

**Phase 2: Agent Orchestration** (High Priority)  
8. ⏳ Create master agent with all tools and workflow coordination
9. ⏳ Update main application with user input and console output formatting

**Phase 3: Testing & Refinement** (Medium Priority)
10. ⏳ Test with real Reddit data and refine subreddit selection logic
11. ⏳ Add comprehensive error handling and rate limiting
12. ⏳ Optimize console output with structured opportunity ranking

**Expected Application Behavior:**
```
$ go run ./cmd/redditanalyzer
Enter project direction: cybersecurity project

🔍 Analyzing: "cybersecurity project"
🎯 Selected subreddits: r/sysadmin, r/cybersecurity, r/msp, r/ITCareerQuestions
📥 Fetching posts (last 7 days): 156 posts found
🔍 Filtering by engagement: 23 posts selected
🤖 Evaluating opportunities...

HIDDEN OPPORTUNITIES FOUND:
🥇 SCORE: 5/5 - Password Management for Small MSPs
   Problem: MSPs struggling with client password management
   Subreddit: r/msp | Upvotes: 47 | Comments: 23
   Analysis: Market gap for affordable password solutions <50 employees
   Link: reddit.com/r/msp/comments/xyz123

🥈 SCORE: 4/5 - Automated Security Compliance Reporting
   Problem: Manual SOC 2 compliance report generation
   Subreddit: r/sysadmin | Upvotes: 34 | Comments: 18
   Analysis: Mid-size companies need automated reporting tools
   Link: reddit.com/r/sysadmin/comments/abc456

Analysis complete: 3 opportunities found from 23 posts analyzed
```

## Development Notes
- Project follows standard Go module structure
- No tests, documentation, or build automation currently exist
- Git repository on main branch with pending changes to go.mod and new cmd/ directory