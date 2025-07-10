package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"reddit-analyzer/internal/redditanalyzer/domain"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/vitalii-honchar/go-agent/pkg/goagent/llm"
)

// SubredditSelectorTool uses AI to intelligently select niche subreddits
// based on the user's project direction, avoiding obvious entrepreneurship communities
type SubredditSelectorTool struct {
	client *openai.Client
}

// SubredditSelectorParams represents the parameters for subreddit selection
type SubredditSelectorParams struct {
	ProjectDirection string `json:"project_direction"`
}

// SubredditSelectorResult represents the result of subreddit selection
type SubredditSelectorResult struct {
	llm.BaseLLMToolResult
	Selection domain.SubredditSelection `json:"selection"`
}

// NewSubredditSelectorTool creates a new SubredditSelectorTool
func NewSubredditSelectorTool() (*SubredditSelectorTool, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &SubredditSelectorTool{
		client: &client,
	}, nil
}

// CreateLLMTool creates the LLM tool for subreddit selection
func (s *SubredditSelectorTool) CreateLLMTool() llm.LLMTool {
	return llm.NewLLMTool(
		llm.WithLLMToolName("select_subreddits"),
		llm.WithLLMToolDescription("AI intelligently selects 3-5 niche subreddits based on project direction, avoiding obvious entrepreneurship communities"),
		llm.WithLLMToolParametersSchema[SubredditSelectorParams](),
		llm.WithLLMToolCall(s.selectSubreddits),
	)
}

// selectSubreddits performs the actual subreddit selection using OpenAI
func (s *SubredditSelectorTool) selectSubreddits(id string, params SubredditSelectorParams) (SubredditSelectorResult, error) {
	projectDirection := params.ProjectDirection

	prompt := s.buildSelectionPrompt(projectDirection)

	resp, err := s.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model:       openai.ChatModelGPT4oMini,
		Temperature: openai.Float(0.3),
		MaxTokens:   openai.Int(int64(1000)),
	})
	if err != nil {
		return SubredditSelectorResult{}, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return SubredditSelectorResult{}, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	
	// Clean the response - remove markdown code blocks if present
	originalContent := content
	content = cleanJSONResponse(content)
	
	var selection domain.SubredditSelection
	if err := json.Unmarshal([]byte(content), &selection); err != nil {
		return SubredditSelectorResult{}, fmt.Errorf("failed to parse OpenAI response: %w\nOriginal: %q\nCleaned: %q", err, originalContent, content)
	}

	return SubredditSelectorResult{
		BaseLLMToolResult: llm.BaseLLMToolResult{ID: id},
		Selection:         selection,
	}, nil
}

// buildSelectionPrompt creates the prompt for subreddit selection
func (s *SubredditSelectorTool) buildSelectionPrompt(projectDirection string) string {
	return fmt.Sprintf(`You are an expert Reddit researcher helping indie hackers find hidden market opportunities.

Your task: Select 3-5 niche subreddits for analyzing "%s" opportunities.

CRITICAL RULES:
1. AVOID obvious entrepreneurship subreddits: r/entrepreneur, r/startups, r/SideProject, r/business
2. Target NICHE communities where domain professionals discuss real problems
3. Focus on communities where solutions are underdeveloped
4. Look for subreddits where people complain about tools/processes but solutions aren't widely discussed

For "%s", consider:
- Professional communities in that domain
- Technical support/help subreddits  
- Industry-specific career subreddits
- Tool/software specific communities
- Problem-focused communities

Return ONLY a JSON object with this structure:
{
  "subreddits": ["subreddit1", "subreddit2", "subreddit3", "subreddit4", "subreddit5"],
  "reasoning": "Brief explanation of why these subreddits were chosen for finding hidden opportunities"
}

Examples of GOOD selections:
- For "cybersecurity": ["sysadmin", "cybersecurity", "msp", "ITCareerQuestions", "networking"]
- For "productivity": ["productivity", "ADHD", "getmotivated", "studytips", "TimeManagement"]

Select subreddits that will reveal genuine problems without obvious solutions.`, projectDirection, projectDirection)
}

// cleanJSONResponse removes markdown code blocks and extra formatting from OpenAI responses
func cleanJSONResponse(content string) string {
	original := content
	content = strings.TrimSpace(content)
	
	// Remove all backticks and common markdown patterns
	content = strings.ReplaceAll(content, "```json", "")
	content = strings.ReplaceAll(content, "```", "")
	content = strings.TrimSpace(content)
	
	// Find JSON content between first { and last }
	start := strings.Index(content, "{")
	if start == -1 {
		// If no JSON found, try to extract from original
		start = strings.Index(original, "{")
		if start == -1 {
			return content
		}
		content = original
	}
	
	end := strings.LastIndex(content, "}")
	if end == -1 || end <= start {
		return content
	}
	
	jsonContent := content[start : end+1]
	
	// Additional cleanup - remove any remaining backticks
	jsonContent = strings.ReplaceAll(jsonContent, "`", "")
	
	return strings.TrimSpace(jsonContent)
}