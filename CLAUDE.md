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

## Enhanced Multi-Agent Architecture

**Hierarchical Multi-Agent System with Validation Loops**:
- **Master Orchestrator Agent**: Coordinates entire workflow and compiles final results
- **Reddit Discovery Agent**: Intelligent subreddit discovery and post collection
- **Problem Validation Agent**: Validates customer problems and measures market demand  
- **Opportunity Assessment Agent**: Evaluates opportunities against indie hacker constraints
- **Inspiration**: Enhanced architecture based on Anthropic's multi-agent research system

**Multi-Agent Workflow**:
1. **Reddit Discovery Agent** - Explores obvious + non-obvious subreddits iteratively
2. **Problem Validation Agent** - Validates problems with 5+ mentions, competitor analysis
3. **Opportunity Assessment Agent** - Scores opportunities with TAM-SAM-SOM analysis
4. **Iterative Refinement Loop** - 2-3 rounds of discovery and validation
5. **Master Orchestrator** - Compiles ranked opportunities with evidence

**Advanced Tools to Build**:
1. **SubredditSelectorTool**: AI selects 3 obvious + 7 non-obvious subreddits
2. **TopSubredditPostsTool**: Fetches posts from multiple subreddits (last 7 days)
3. **FilterPostsTool**: Filters by engagement metrics and relevance scoring
4. **ProblemExtractorTool**: Extracts explicit and implicit problems from posts
5. **ValidationSearchTool**: Searches Reddit for similar problems (5+ mentions)
6. **CompetitorAnalysisTool**: Web search for existing solutions and gaps
7. **OpportunityEvaluatorTool**: Scores opportunities against user constraints
8. **TAM-SAM-SOM Tool**: Market size analysis for promising opportunities
9. **NicheAnalysisTool**: Identifies underdeveloped niches with low competition

**Intelligence Focus**:
- Multi-round discovery: obvious → non-obvious → pattern-based subreddits
- Quantified validation: 5+ problem mentions across different communities
- Cross-domain insights: cybersecurity + MSP tools, infrastructure + remote work
- Evidence-based ranking: every opportunity backed by real Reddit discussions
- Competitive intelligence: identifies gaps in existing solutions

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
1. ✅ Add go-agent dependency to project (go mod tidy)
2. ✅ Set up environment configuration (.env file with OPENAI_API_KEY)
3. ✅ Define core data structures (RedditPost, OpportunityAnalysis, AgentResult)

**Phase 1: Core Tools Implementation** (High Priority)
4. ✅ Implement SubredditSelectorTool - AI intelligently selects 3-5 niche subreddits
5. ✅ Implement TopSubredditPostsTool - fetches posts from multiple subreddits (last 7 days)
6. ✅ Implement FilterPostsTool - filters by engagement metrics (min upvotes/comments)
7. ✅ Implement EvaluatePostTool - LLM rates posts 1-5 for hidden indie opportunities

**Phase 2: Agent Orchestration** (High Priority)  
8. ✅ Create master agent with all tools and workflow coordination **(FIXED - NOW TRUE REACT AGENT)**
9. ✅ Update main application with user input and console output formatting

**Phase 3: ReAct Agent Architecture Fix** (CRITICAL)
10. ✅ Fix agent implementation to use true ReAct agent from go-agent library
11. ✅ Update agent to use agent.NewAgent() with proper configuration
12. ✅ Implement proper tool registration and agent execution workflow
13. ✅ Add typed result structures for agent output
14. ❌ Test ReAct agent with iterative reasoning and tool usage **(JSON PARSING ERROR - NEEDS TOOL FIX)**

**Phase 4: Testing & Refinement** (Medium Priority)
15. ✅ Test with real Reddit data and refine subreddit selection logic
16. ⏳ Add comprehensive error handling and rate limiting
17. ⏳ Optimize console output with structured opportunity ranking

**Phase 5: CLI UI Enhancement** (High Priority)
18. ⏳ Research and integrate BubbleTea TUI framework
19. ⏳ Design interactive UI wireframes for analysis workflow
20. ⏳ Implement progress tracking and real-time status updates
21. ⏳ Create interactive opportunity browser with filtering
22. ⏳ Add keyboard shortcuts and navigation controls
23. ⏳ Enhance visual styling with Lip Gloss components

### ReAct Agent Architecture Fix Details

**Current Problem:**
The existing implementation in `reddit_research_agent.go` is NOT using a true ReAct agent. It's just making direct LLM calls with tools attached, which means:
- No iterative reasoning and acting cycle
- No ability for the agent to reason about tool results and decide next actions
- Limited to single-shot LLM responses instead of multi-step problem solving

**Required Changes:**
1. **Import Agent Package**: Add `github.com/vitalii-honchar/go-agent/pkg/goagent/agent`
2. **Create ReAct Agent**: Use `agent.NewAgent()` with proper configuration
3. **Tool Registration**: Register each tool individually with `WithTool()`
4. **Agent Execution**: Use `agent.Run()` instead of `llm.Call()`
5. **Typed Results**: Define proper result structures for type safety

**Agent Configuration Pattern:**
```go
redditAgent, err := agent.NewAgent(
    agent.WithName[*domain.AnalysisResult]("reddit_research_agent"),
    agent.WithLLMConfig[*domain.AnalysisResult](llm.LLMConfig{
        Type:        llm.LLMTypeOpenAI,
        APIKey:      apiKey,
        Model:       "gpt-4o",
        Temperature: 0.7,
    }),
    agent.WithBehavior[*domain.AnalysisResult]("You are a Reddit Research Agent..."),
    agent.WithTool[*domain.AnalysisResult]("select_subreddits", subredditTool),
    agent.WithTool[*domain.AnalysisResult]("fetch_reddit_posts", fetchTool),
    agent.WithTool[*domain.AnalysisResult]("filter_posts", filterTool),
    agent.WithTool[*domain.AnalysisResult]("evaluate_post", evaluateTool),
    agent.WithOutputSchema[*domain.AnalysisResult](&domain.AnalysisResult{}),
)
```

**Expected ReAct Workflow:**
1. **Think**: Agent analyzes the project direction and plans approach
2. **Act**: Agent selects and calls appropriate tools (subreddit selection)
3. **Observe**: Agent reviews tool results and decides next action
4. **Think**: Agent reasons about which posts to fetch based on subreddits
5. **Act**: Agent fetches posts from selected subreddits
6. **Observe**: Agent analyzes fetched posts
7. **Think**: Agent determines filtering criteria
8. **Act**: Agent filters posts by engagement
9. **Observe**: Agent reviews filtered results
10. **Think**: Agent decides how to evaluate each post
11. **Act**: Agent evaluates posts for opportunities
12. **Observe**: Agent compiles final analysis results

### BubbleTea TUI Implementation Plan

**Dependencies to Add:**
- `github.com/charmbracelet/bubbletea` - Core TUI framework
- `github.com/charmbracelet/bubbles` - Pre-built components (spinner, progress, list)
- `github.com/charmbracelet/lipgloss` - Styling and layout

**UI Architecture (Model-Update-View):**

**Model State:**
```go
type AppModel struct {
    state        AppState // current view (input, analyzing, results)
    projectInput string
    progress     float64
    currentStep  string
    opportunities []domain.RankedOpportunity
    selectedIndex int
    viewport     viewport.Model
    spinner      spinner.Model
    list         list.Model
}
```

**App States:**
1. **InputState** - Project direction input with validation
2. **AnalyzingState** - Progress tracking with animated spinner
3. **ResultsState** - Interactive opportunity browser with filtering

**Key UI Components:**
- **Input Screen**: Styled text input with validation and help text
- **Progress Screen**: Multi-step progress bar with real-time updates
- **Results Browser**: Scrollable list with detailed opportunity cards
- **Opportunity Detail**: Expandable view with full analysis and links

**Interactive Features:**
- **Navigation**: Arrow keys, vim-style (j/k), tab navigation
- **Filtering**: Filter opportunities by score, subreddit, keywords
- **Sorting**: Sort by score, engagement, subreddit
- **Actions**: Copy links, open in browser, save results

**Visual Design:**
- **Color Scheme**: Professional dark theme with accent colors
- **Typography**: Headers, body text, code styling with Lip Gloss
- **Layout**: Responsive panels with proper spacing and borders
- **Animations**: Smooth transitions, loading spinners

**Expected Enhanced UI Behavior:**
```
┌─ Reddit Opportunity Analyzer ─────────────────────────────────┐
│                                                               │
│  🔍 Enter your project direction:                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ cybersecurity project                                   │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                               │
│  💡 Examples: "SaaS for small businesses", "mobile app"      │
│                                                               │
│  Press Enter to analyze • Ctrl+C to quit                     │
└───────────────────────────────────────────────────────────────┘

┌─ Analyzing: "cybersecurity project" ──────────────────────────┐
│                                                               │
│  🎯 Selecting subreddits...              ████████████ 100%   │
│  📥 Fetching posts (156 found)...        ████████████ 100%   │
│  🔍 Filtering by engagement (23 posts)   ████████████ 100%   │
│  🤖 Evaluating opportunities...          ██████░░░░░░  60%   │
│                                                               │
│  ⏱️  Estimated time remaining: 45 seconds                     │
└───────────────────────────────────────────────────────────────┘

┌─ Hidden Opportunities Found (3 results) ──────────────────────┐
│                                                               │
│  🥇 SCORE: 5/5 - Password Management for Small MSPs          │
│  └─ r/msp • 47 upvotes • 23 comments                         │
│                                                               │
│  🥈 SCORE: 4/5 - Automated Security Compliance               │
│  └─ r/sysadmin • 34 upvotes • 18 comments                    │
│                                                               │
│  🥉 SCORE: 4/5 - Network Monitoring for Remote Teams         │
│  └─ r/networking • 28 upvotes • 15 comments                  │
│                                                               │
│  ↑↓ Navigate • Enter View Details • F Filter • S Sort • Q Quit│
└───────────────────────────────────────────────────────────────┘
```

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