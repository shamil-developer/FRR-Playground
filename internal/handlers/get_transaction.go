// internal/handlers/get_transaction.go

package handlers

import (
	"context"

	"frr-playground/frrpb"
	"frr-playground/internal/models"
	"frr-playground/internal/utils"

	"google.golang.org/grpc"
)

type GetTransactionHandler struct{}

func NewGetTransactionHandler() *GetTransactionHandler {
	return &GetTransactionHandler{}
}

func (h *GetTransactionHandler) Execute(
	ctx context.Context,
	conn *grpc.ClientConn,
	flow map[string]any,
	step models.Step,
) error {

	client := frrpb.NewNorthboundClient(conn)

	transactionID := uint32(
		step.Params["transactionId"].(int),
	)

	resp, err := client.GetTransaction(
		ctx,
		&frrpb.GetTransactionRequest{
			TransactionId: transactionID,
			Encoding:      frrpb.Encoding_JSON,
			WithDefaults:  true,
		},
	)
	if err != nil {
		return err
	}

	err = utils.PrintAny(resp)

	if err != nil {
		return err
	}

	return nil
}
