package service

import (
	"net/http"
	"bytes"
	"context"
	"encoding/json"
	"time"
	"fmt"
	"github.com/souls-syntax/Templates/internal/models"
)

type LlmClient struct {
	BaseURL	string
	HttpClient	*http.Client
}

func NewLlmClient(url string) *LlmClient {
	return &LlmClient{
		BaseURL: url,
		HttpClient: &http.Client{Timeout: 60*time.Second},
	}
}

type llmRequest struct {
	Query string `json:"query_text"`
}

type llmResponse struct {
	Verdict					string	`json:"verdict"`
	Confidence			float64	`json:"confidence"`
	Decider					string	`json:"decider"`
	// Explanation currently as place holder
}

func (l *LlmClient) GetAnalysis(ctx context.Context, text string) (models.Decision, error) {
	reqBody, err := json.Marshal(llmRequest{Query: text})
	if err != nil {
		return models.Decision{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		l.BaseURL+"/predict_llm",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return models.Decision{}, err
	}

	req.Header.Set("Content-Type","application/json")
	
	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return models.Decision{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Decision{}, fmt.Errorf("Llm general returned %d", resp.StatusCode)
	}

	var llmResp llmResponse

	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return models.Decision{}, err
	}

	return models.Decision{
		Verdict: llmResp.Verdict,
		Confidence: llmResp.Confidence,
		Decider:	llmResp.Decider,
	}, nil
}
