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

type BertClient struct {
	BaseURL	string
	HttpClient	*http.Client
}

func NewBertClient(url string) *BertClient {
	return &BertClient{
		BaseURL: url,
		HttpClient: &http.Client{Timeout: 5*time.Second},
	}
}

type bertRequest struct {
	Query string `json:"query_text"`
}

type bertResponse struct {
	Verdict					string	`json:"verdict"`
	Confidence			float64	`json:"confidence"`
	Decider					string	`json:"decider"`
}

func (b *BertClient) GetVerdict(ctx context.Context, text string) (models.Decision, error) {
	reqBody, err := json.Marshal(bertRequest{Query: text})
	if err != nil {
		return models.Decision{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		b.BaseURL+"/predict",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return models.Decision{}, err
	}

	req.Header.Set("Content-Type","application/json")
	
	resp, err := b.HttpClient.Do(req)
	if err != nil {
		return models.Decision{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Decision{}, fmt.Errorf("bert returned %d", resp.StatusCode)
	}

	var bertResp bertResponse

	if err := json.NewDecoder(resp.Body).Decode(&bertResp); err != nil {
		return models.Decision{}, err
	}

	return models.Decision{
		Verdict: bertResp.Verdict,
		Confidence: bertResp.Confidence,
		Decider:	bertResp.Decider,
	}, nil
}
