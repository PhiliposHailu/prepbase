package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/philipos/prepbase/domain"
	"google.golang.org/api/option"
)

type geminiService struct {
	apiKey string
}

func NewAIService() domain.AIService {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		fmt.Println("⚠️ WARNING: GEMINI_API_KEY is not set in .env")
	}
	return &geminiService{apiKey: key}
}

func (s *geminiService) GenerateHint(content string) (string, error) {
	if s.apiKey == "" {
		return "", errors.New("AI service is disabled (missing API key)")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.apiKey))

	// models := client.ListModels(ctx)

	// for {
	// 	m, err := models.Next()
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Println(m.Name)
	// }

	if err != nil {
		return "", err
	}
	defer client.Close()

	// Use the lightweight text model
	model := client.GenerativeModel("gemini-2.5-flash")

	// The prompt we send to the AI
	prompt := fmt.Sprintf("Act as an expert technical interviewer. Give me a short, 2-sentence hint to solve this problem, but DO NOT give me the exact code solution. Here is the problem: %s", content)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	// Extract the text from the complex Google response object
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
	}

	return "", errors.New("AI returned an empty response")
}
