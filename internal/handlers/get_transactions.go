// internal/handlers/get_transactions.go

package handlers

import (
	"context"
	"io"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"

	"google.golang.org/grpc"
)

type GetTransactionsHandler struct{}

func NewGetTransactionsHandler() *GetTransactionsHandler {
	return &GetTransactionsHandler{}
}

func (h *GetTransactionsHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	client := frrpb.NewNorthboundClient(conn)

	stream, err := client.ListTransactions(
		ctx,
		&frrpb.ListTransactionsRequest{},
	)
	if err != nil {
		return err
	}

	for {

		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		err = utils.PrintAny(resp)

		if err != nil {
			return err
		}
	}

	return nil
}
