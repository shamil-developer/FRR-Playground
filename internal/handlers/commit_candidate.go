package handlers

import (
	"context"

	"google.golang.org/grpc"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"
)

type CommitCandidateHandler struct{}

func NewCommitCandidateHandler() *CommitCandidateHandler {
	return &CommitCandidateHandler{}
}

func (h *CommitCandidateHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	client := frrpb.NewNorthboundClient(conn)

	candidateId := flow["candidateId"].(uint32)

	response, err := client.Commit(
		ctx,
		&frrpb.CommitRequest{
			CandidateId: candidateId,

			Phase: frrpb.CommitRequest_ALL,

			Comment: "static route commit",
		},
	)

	if err != nil {
		return err
	}

	return utils.PrintAny(response)
}
