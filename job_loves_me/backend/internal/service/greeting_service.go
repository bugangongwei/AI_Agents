package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type GreetingService struct {
	apiKey     string
	baseURL    string
	resumeServ *ResumeService
}

func NewGreetingService(resumeServ *ResumeService) *GreetingService {
	return &GreetingService{
		apiKey:     os.Getenv("DEEPSEEK_API_KEY"),
		baseURL:    os.Getenv("DEEPSEEK_BASE_URL"),
		resumeServ: resumeServ,
	}
}

type DeepSeekRequest struct {
	Model    string            `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (s *GreetingService) GenerateGreeting(ctx context.Context, userID uint, jdText string) (<-chan string, error) {
	// Get latest resume
	resume, err := s.resumeServ.GetLatestResume(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest resume: %w", err)
	}

	prompt := fmt.Sprintf(`你是一位求职专家。请根据以下简历内容和职位描述(JD)，写一段发给招聘人员的打招呼推荐语。

要求：
1. 风格简洁、专业、有礼貌。
2. 目的：吸引注意并获得面试机会。
3. 长度：200字以内。
4. 受众：HR、技术负责人、开发人员或猎头。

简历内容：
%s

职位描述(JD)：
%s`, resume.RawText, jdText)

	reqBody := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []DeepSeekMessage{
			{Role: "user", Content: prompt},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("DeepSeek API error: status %d", resp.StatusCode)
		}

	ch := make(chan string)
	go func() {
		defer resp.Body.Close()
		defer close(ch)
		// Placeholder for SSE parsing logic in next iteration
	}()

	return ch, nil
}

// Non-streaming version for MVP stability
func (s *GreetingService) GenerateGreetingSync(ctx context.Context, userID uint, jdText string) (string, error) {
	resume, err := s.resumeServ.GetLatestResume(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get latest resume: %w", err)
	}

	prompt := fmt.Sprintf(`你是一位求职专家。请根据以下简历内容和职位描述(JD)，写一段发给招聘人员的打招呼推荐语。

要求：
1. 风格简洁、专业、有礼貌。
2. 目的：吸引注意并获得面试机会。
3. 长度：200字以内。
4. 受众：HR、技术负责人、开发人员或猎头。

简历内容：
%s

职位描述(JD)：
%s`, resume.RawText, jdText)

	reqBody := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []DeepSeekMessage{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API error: status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from DeepSeek")
	}

	return result.Choices[0].Message.Content, nil
}
