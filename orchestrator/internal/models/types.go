package models

type Query struct {
	Query  string  `json:"query"`
	UserID *string `json:"userid"`
}

type Response struct {
	Response string
	Status   bool
}

type Decision struct {
	QueryHash string
	QueryText string

	Verdict    string
	Confidence float64
	Decider    string
}

type Observation struct {
	Source           string
	ProcessingTimeMs int64
}

type VerifyResponse struct {
	QueryHash string `json:"query_hash"` // Unique Identifier
	QueryText string `json:"query_text"` // The original question

	DecisionVerdict    string  `json:"decision_verdict"`    // Likely_True | Likely_False
	DecisionConfidence float64 `json:"decision_confidence"` // Between 0 to 1
	DecisionDecider    string  `json:"decision_decider"`    // The verdict giver BERT | LLM | human

	ObsSource           string `json:"obs_source"`             // Where we got the data | CACHE | DB | BERT
	ObsProcessingTimeMs int64  `json:"obs_processing_time_ms"` // Time it took to process the request
}

type CacheDecision struct {
	Verdict    string
	Confidence float64
	Decider    string
}

func (d Decision) ToCacheModel() CacheDecision {
	return CacheDecision{
		Verdict:    d.Verdict,
		Confidence: d.Confidence,
		Decider:    d.Decider,
	}
}
