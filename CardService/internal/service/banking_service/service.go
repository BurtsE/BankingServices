package banking_service

import (
	"CardService/generated/protobuf"
	"CardService/internal/service"
	"context"
	"time"
)

var _ service.IBankingService = (*BankingService)(nil)

type BankingService struct {
	client protobuf.BankingServiceClient
}

func NewBankingService(client protobuf.BankingServiceClient) *BankingService {
	return &BankingService{client: client}
}

func (b *BankingService) AccountIsActive(accountID string) (ok bool, err error) {

	isActiveRequest := &protobuf.IsActiveRequest{AccountId: accountID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := b.client.AccountIsActive(ctx, isActiveRequest)
	if err != nil {
		return false, err
	}

	return resp.GetIsActive(), nil
}
