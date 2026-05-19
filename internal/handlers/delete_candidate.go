package handlers

import (
	"context"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type DeleteCandidateHandler struct{}

func NewDeleteCandidateHandler() *DeleteCandidateHandler {
	return &DeleteCandidateHandler{}
}

func (h *DeleteCandidateHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {
	client := frrpb.NewNorthboundClient(conn)
	candidateId := flow["candidateId"].(uint32)

	response, err := client.DeleteCandidate(
		ctx,
		&frrpb.DeleteCandidateRequest{CandidateId: candidateId},
	)
	if err != nil {
		return err
	}

	delete(flow, "candidateId")
	return utils.PrintAny(response)
}
