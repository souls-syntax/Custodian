package service

import (
	"github.com/souls-syntax/Templates/internal/models"
)

func BuildResponse(dec models.Decision, obs models.Observation) models.VerifyResponse {
	return models.VerifyResponse {
		QueryHash:						dec.QueryHash,
		QueryText:						dec.QueryText,

		DecisionVerdict:			dec.Verdict,
		DecisionConfidence:		dec.Confidence,
		DecisionDecider:			dec.Decider,

		ObsSource:						obs.Source,
		ObsProcessingTimeMs:	obs.ProcessingTimeMs,
	}
}
