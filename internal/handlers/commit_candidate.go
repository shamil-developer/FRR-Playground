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
	phase := frrpb.CommitRequest_ALL

	if rawComment, ok := step.Params["comment"]; ok && rawComment != nil {
		comment = fmt.Sprint(rawComment)
	}
	if rawPhase, ok := step.Params["phase"]; ok && rawPhase != nil {
		switch fmt.Sprint(rawPhase) {
		case "VALIDATE":
			phase = frrpb.CommitRequest_VALIDATE
		case "PREPARE":
			phase = frrpb.CommitRequest_PREPARE
		case "ABORT":
			phase = frrpb.CommitRequest_ABORT
		case "APPLY":
			phase = frrpb.CommitRequest_APPLY
		case "ALL":
			phase = frrpb.CommitRequest_ALL
		default:
			return fmt.Errorf("unknown commit phase %q", rawPhase)
		}
	}

	response, err := client.Commit(
		ctx,
		&frrpb.CommitRequest{
			CandidateId: candidateId,

			Phase: phase,

			Comment: comment,
		},
	)

	if err != nil {
		return err
	}

	return utils.PrintAny(response)
}
