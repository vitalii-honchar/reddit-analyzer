package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"reddit-analyzer/internal/searchagent"
	"reddit-analyzer/internal/searchagent/domain"
)

func main() {
	// Set up logger with human-readable format
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	// Create the Reddit Research Agent
	researchAgent, err := agent.NewRedditResearchAgent(log)
	if err != nil {
		log.WithError(err).Fatal("Failed to create Reddit Research Agent")
	}

	// Hardcode cybersecurity for testing
	projectDirection := "cybersecurity"
	fmt.Print("Enter project direction: ")
	fmt.Println(projectDirection)

	fmt.Printf("\n🔍 Analyzing: \"%s\"\n", projectDirection)

	// Run the analysis
	ctx := context.Background()
	result, err := researchAgent.AnalyzeProject(ctx, projectDirection)
	if err != nil {
		log.WithError(err).Fatal("Analysis failed")
	}

	// Display results
	displayResults(result)
}

func displayResults(result *domain.AnalysisResult) {
	fmt.Printf("\n🎯 Selected subreddits: %s\n", strings.Join(result.SelectedSubreddits, ", "))
	fmt.Printf("📥 Fetching posts (last 7 days): %d posts found\n", result.PostsAnalyzed)
	fmt.Printf("🔍 Filtering by engagement: %d posts selected\n", result.PostsFiltered)
	fmt.Printf("🤖 Evaluating opportunities...\n\n")

	if len(result.Opportunities) == 0 {
		fmt.Println("No opportunities found matching the criteria.")
		return
	}

	// Sort opportunities by score (descending)
	sort.Slice(result.Opportunities, func(i, j int) bool {
		return result.Opportunities[i].Analysis.Score > result.Opportunities[j].Analysis.Score
	})

	fmt.Println("HIDDEN OPPORTUNITIES FOUND:")
	
	medals := []string{"🥇", "🥈", "🥉"}
	for i, opportunity := range result.Opportunities {
		if i >= 3 {
			break // Show only top 3
		}
		
		medal := medals[i]
		if i >= len(medals) {
			medal = "🏅"
		}

		fmt.Printf("%s SCORE: %d/5 - %s\n", medal, opportunity.Analysis.Score, opportunity.Analysis.ProblemSummary)
		fmt.Printf("   Problem: %s\n", opportunity.Post.Title)
		fmt.Printf("   Subreddit: r/%s | Upvotes: %d | Comments: %d\n", opportunity.Post.Subreddit, opportunity.Post.Score, opportunity.Post.NumComments)
		fmt.Printf("   Analysis: %s\n", opportunity.Analysis.MarketAnalysis)
		fmt.Printf("   Link: reddit.com/r/%s/comments/%s\n\n", opportunity.Post.Subreddit, opportunity.Post.ID)
	}

	fmt.Printf("Analysis complete: %d opportunities found from %d posts analyzed\n", len(result.Opportunities), result.PostsFiltered)
}
