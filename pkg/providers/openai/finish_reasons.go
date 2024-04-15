package openai

import (
	"glide/pkg/api/schemas"
	"glide/pkg/telemetry"
	"go.uber.org/zap"
)

var (
	// Reference: https://platform.openai.com/docs/api-reference/chat/object

	CompleteReason  = "stop"
	MaxTokensReason = "length"
	FilteredReason  = "content_filter"
)

func NewFinishReasonMapper(tel *telemetry.Telemetry) *FinishReasonMapper {
	return &FinishReasonMapper{
		tel: tel,
	}
}

type FinishReasonMapper struct {
	tel *telemetry.Telemetry
}

func (m *FinishReasonMapper) Map(finishReason string) *schemas.FinishReason {
	if len(finishReason) == 0 {
		return nil
	}

	var reason *schemas.FinishReason

	switch finishReason {
	case CompleteReason:
		reason = &schemas.Complete
	case MaxTokensReason:
		reason = &schemas.MaxTokens
	case FilteredReason:
		reason = &schemas.ContentFiltered
	default:
		m.tel.Logger.Warn(
			"Unknown finish reason, other is going to used",
			zap.String("unknown_reason", finishReason),
		)

		reason = &schemas.OtherReason
	}

	return reason
}
