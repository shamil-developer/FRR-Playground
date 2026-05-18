package handlers

import (
	"context"
	"fmt"

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

	comment := "workflow commit"

	if rawComment, ok := step.Params["comment"]; ok && rawComment != nil {
		comment = fmt.Sprint(rawComment)
	}

	response, err := client.Commit(
		ctx,
		&frrpb.CommitRequest{
			CandidateId: candidateId,

			Phase: frrpb.CommitRequest_ALL,

			Comment: comment,
		},
	)

	if err != nil {
		return err
	}

	return utils.PrintAny(response)
}
