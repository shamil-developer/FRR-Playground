package handlers

import (
	"context"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type CreateCandidateHandler struct{}

func NewCreateCandidateHandler() *CreateCandidateHandler {
	return &CreateCandidateHandler{}
}

func (h *CreateCandidateHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	client := frrpb.NewNorthboundClient(conn)

	response, err := client.CreateCandidate(
		ctx,
		&frrpb.CreateCandidateRequest{},
	)

	if err != nil {
		return err
	}

	flow["candidateId"] = response.CandidateId

	return utils.PrintAny(response)
}
