package workflow

import (
	"context"
	"frr-playground/internal/handlers"
	"frr-playground/internal/models"

	"google.golang.org/grpc"
)

type Handler interface {
	Execute(
		ctx context.Context,
		conn *grpc.ClientConn,
		flow map[string]any,
		step models.Step,
	) error
}

func CreateHandlersRegistry() map[string]Handler {

	return map[string]Handler{
		"get":                     handlers.NewGetHandler(),
		"getCapabilities":         handlers.NewGetCapabilitiesHandler(),
		"createCandidate":         handlers.NewCreateCandidateHandler(),
		"editCandidate":           handlers.NewEditCandidateHandler(),
		"commitCandidate":         handlers.NewCommitCandidateHandler(),
		"getTransactions":         handlers.NewGetTransactionsHandler(),
		"getTransaction":          handlers.NewGetTransactionHandler(),
		"sleep":                   handlers.NewSleepHandler(),
		"assertJsonExists":        handlers.NewAssertJSONExistsHandler(),
		"assertJsonMissing":       handlers.NewAssertJSONMissingHandler(),
		"assertCapabilityExists":  handlers.NewAssertCapabilityExistsHandler(),
		"assertCapabilityMissing": handlers.NewAssertCapabilityMissingHandler(),
	}

}
