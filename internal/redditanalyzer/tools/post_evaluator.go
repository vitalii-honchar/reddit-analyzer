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

// EvaluatePostTool uses LLM to rate posts for indie hacker opportunities
type EvaluatePostTool struct {
	client *openai.Client
}

// EvaluatePostResult represents the result of post evaluation
type EvaluatePostResult struct {
	llm.BaseLLMToolResult
	Analysis domain.OpportunityAnalysis `json:"analysis"`
}

// NewEvaluatePostTool creates a new EvaluatePostTool
func NewEvaluatePostTool() (*EvaluatePostTool, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &EvaluatePostTool{
		client: &client,
	}, nil
}

// CreateLLMTool creates the LLM tool for post evaluation
func (e *EvaluatePostTool) CreateLLMTool() llm.LLMTool {
	return llm.NewLLMTool(
		llm.WithLLMToolName("evaluate_post"),
		llm.WithLLMToolDescription("Uses LLM to rate a Reddit post 1-5 for hidden indie hacker opportunities"),
		llm.WithLLMToolParametersSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"post": map[string]any{
					"type":        "object",
					"description": "Reddit post to evaluate for opportunity potential",
				},
				"project_direction": map[string]any{
					"type":        "string",
					"description": "The user's project direction for context",
				},
			},
			"required": []string{"post", "project_direction"},
		}),
		llm.WithLLMToolCall(e.evaluatePost),
	)
}

// evaluatePost performs the actual post evaluation using OpenAI
func (e *EvaluatePostTool) evaluatePost(id string, args map[string]any) (EvaluatePostResult, error) {
	postData, ok := args["post"].(map[string]interface{})
	if !ok {
		return EvaluatePostResult{}, fmt.Errorf("post must be an object")
	}

	projectDirection, ok := args["project_direction"].(string)
	if !ok {
		return EvaluatePostResult{}, fmt.Errorf("project_direction must be a string")
	}

	// Extract post information
	title, _ := postData["title"].(string)
	content, _ := postData["content"].(string)
	subreddit, _ := postData["subreddit"].(string)
	score := 0
	if s, ok := postData["score"].(float64); ok {
		score = int(s)
	}
	comments := 0
	if c, ok := postData["num_comments"].(float64); ok {
		comments = int(c)
	}

	prompt := e.buildEvaluationPrompt(title, content, subreddit, score, comments, projectDirection)

	resp, err := e.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model:       openai.ChatModelGPT4oMini,
		Temperature: openai.Float(0.2),
		MaxTokens:   openai.Int(int64(800)),
	})
	if err != nil {
		return EvaluatePostResult{}, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return EvaluatePostResult{}, fmt.Errorf("no response from OpenAI")
	}

	content_response := resp.Choices[0].Message.Content
	
	// Clean the response to extract JSON using more robust cleaning
	originalResponse := content_response
	cleanedResponse := strings.TrimSpace(content_response)
	
	// Remove all backticks and common markdown patterns
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "```json", "")
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "```", "")
	cleanedResponse = strings.TrimSpace(cleanedResponse)
	
	// Find JSON content between first { and last }
	start := strings.Index(cleanedResponse, "{")
	if start == -1 {
		start = strings.Index(originalResponse, "{")
		if start == -1 {
			return EvaluatePostResult{}, fmt.Errorf("no JSON found in response: %q", originalResponse)
		}
		cleanedResponse = originalResponse
	}
	
	end := strings.LastIndex(cleanedResponse, "}")
	if end == -1 || end <= start {
		return EvaluatePostResult{}, fmt.Errorf("invalid JSON structure in response: %q", cleanedResponse)
	}
	
	jsonContent := cleanedResponse[start : end+1]
	// Additional cleanup - remove any remaining backticks
	jsonContent = strings.ReplaceAll(jsonContent, "`", "")
	jsonContent = strings.TrimSpace(jsonContent)
	
	var analysis domain.OpportunityAnalysis
	if err := json.Unmarshal([]byte(jsonContent), &analysis); err != nil {
		return EvaluatePostResult{}, fmt.Errorf("failed to parse OpenAI response: %w\nOriginal: %q\nCleaned: %q", err, originalResponse, jsonContent)
	}

	// Set the post ID from the original post data
	if postID, ok := postData["id"].(string); ok {
		analysis.PostID = postID
	}

	return EvaluatePostResult{
		BaseLLMToolResult: llm.BaseLLMToolResult{ID: id},
		Analysis:          analysis,
	}, nil
}

// buildEvaluationPrompt creates the prompt for post evaluation
func (e *EvaluatePostTool) buildEvaluationPrompt(title, content, subreddit string, score, comments int, projectDirection string) string {
	return fmt.Sprintf(`You are an expert indie hacker opportunity analyst. Evaluate this Reddit post for hidden business opportunities.

PROJECT DIRECTION: "%s"

REDDIT POST TO ANALYZE:
- Subreddit: r/%s
- Title: %s
- Content: %s
- Score: %d upvotes
- Comments: %d

EVALUATION CRITERIA:
1. SCORE (1-5): Rate the opportunity potential for an indie hacker
   - 5: Exceptional opportunity - clear problem, underserved market, feasible solution
   - 4: Strong opportunity - good problem with some existing solutions but gaps remain
   - 3: Moderate opportunity - valid problem but competitive/complex market
   - 2: Weak opportunity - problem exists but difficult to monetize or serve
   - 1: Poor opportunity - no clear problem or oversaturated market

2. Focus on INDIE HACKER viability:
   - Can be built by 1-2 people
   - Has clear monetization potential
   - Targets underserved niche
   - Problem mentioned by domain professionals
   - Not just complaints but actual workflow pain points

3. Avoid rating high if:
   - Problem has obvious enterprise solutions
   - Requires massive infrastructure
   - Highly regulated industries
   - Purely consumer complaints without B2B potential

Return ONLY a JSON object with this structure:
{
  "score": 3,
  "problem_summary": "Brief description of the core problem mentioned",
  "market_analysis": "Analysis of market gap and opportunity size",
  "target_audience": "Who specifically has this problem",
  "competition_note": "Notes about existing solutions or lack thereof",
  "reasoning": "Why you gave this score - focus on indie hacker feasibility"
}

Focus on finding problems that professionals complain about in their daily work that could be solved with simple tools or services.`, 
		projectDirection, subreddit, title, content, score, comments)
}