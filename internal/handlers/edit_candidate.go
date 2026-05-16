package handlers

import (
	"context"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type EditCandidateHandler struct{}

func NewEditCandidateHandler() *EditCandidateHandler {
	return &EditCandidateHandler{}
}

func (h *EditCandidateHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	client := frrpb.NewNorthboundClient(conn)

	candidateId := flow["candidateId"].(uint32)

	var updates []*frrpb.PathValue

	rawUpdates := step.Params["update"].([]any)

	for _, item := range rawUpdates {

		update := item.(map[string]any)

		updates = append(
			updates,
			&frrpb.PathValue{
				Path: update["path"].(string),

				Value: update["value"].(string),
			},
		)
	}

	response, err := client.EditCandidate(
		ctx,
		&frrpb.EditCandidateRequest{
			CandidateId: candidateId,

			Update: updates,
		},
	)

	if err != nil {
		return err
	}

	return utils.PrintAny(response)
}
