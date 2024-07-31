package handler

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

func NewDataAccessor(accessor DataAccessor) *dataAccessor {
	return &dataAccessor{
		accessor: accessor,
	}
}

type dataAccessor struct {
	proto.UnimplementedDataAccessorServer
	accessor DataAccessor
}

func (da *dataAccessor) Hello(ctx context.Context, emtpy *empty.Empty) (*proto.HelloResponse, error) {
	action := domain.GetAction(1)

	dd, err := da.accessor.GetHelloData(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	resp := &proto.HelloResponse{
		HelloEncrypted:     dd.HelloEncrypted,
		MasterPasswordHint: dd.MasterPasswordHint,
	}

	return resp, nil
}

func (da *dataAccessor) GetBankCardList(ctx context.Context, emtpy *empty.Empty) (*proto.BankCardListResponse, error) {
	action := domain.GetAction(1)

	list, err := da.accessor.GetBankCardList(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	resp := &proto.BankCardListResponse{}
	for _, card := range list {
		resp.Cards = append(resp.Cards, &proto.BankCard{
			Number:  card.Number,
			Content: card.Content,
		})
	}
	return resp, nil
}

func (da *dataAccessor) CreateBankCard(ctx context.Context, req *proto.CreateBankCardRequest) (*empty.Empty, error) {
	action := domain.GetAction(1)

	encrBanckCard := &domain.EncryptedBankCard{
		Number:  req.Number,
		Content: req.Content,
	}
	err := da.accessor.CreateBankCard(ctx, encrBanckCard)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}
	return nil, nil
}

func (da *dataAccessor) DeleteBankCard(ctx context.Context, req *proto.DeleteBankCardRequest) (*empty.Empty, error) {
	action := domain.GetAction(1)
	err := da.accessor.DeleteBankCard(ctx, req.Number)

	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}
	return nil, nil
}

func (da *dataAccessor) UpdateBankCard(ctx context.Context, req *proto.UpdateBankCardRequest) (*empty.Empty, error) {
	action := domain.GetAction(1)
	err := da.accessor.UpdateBankCard(ctx, &domain.EncryptedBankCard{
		Number:  req.Number,
		Content: req.Content,
	})

	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	return nil, nil
}
func (da *dataAccessor) GetUserPasswordDataList(ctx context.Context, emtpy *empty.Empty) (*proto.UserPasswordDataResponse, error) {
	action := domain.GetAction(1)

	list, err := da.accessor.GetUserPasswordDataList(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	resp := &proto.UserPasswordDataResponse{}
	for _, dt := range list {
		resp.Datas = append(resp.Datas, &proto.UserPasswordData{
			Hint:    dt.Hint,
			Content: dt.Content,
		})
	}
	return resp, nil
}

func (da *dataAccessor) CreateUserPasswordData(ctx context.Context, req *proto.CreateUserPasswordDataRequest) (*empty.Empty, error) {
	action := domain.GetAction(1)

	usedData := &domain.EncryptedUserPasswordData{
		Hint:    req.Hint,
		Content: req.Content,
	}
	err := da.accessor.CreateUserPasswordData(ctx, usedData)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}
	return nil, nil
}

func (da *dataAccessor) DeleteUserPasswordData(ctx context.Context, req *proto.DeleteUserPasswordDataRequest) (*empty.Empty, error) {
	action := domain.GetAction(1)
	err := da.accessor.DeleteUserPasswordData(ctx, req.Hint)

	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}
	return nil, nil
}

func (da *dataAccessor) UpdateUserPasswordData(ctx context.Context, req *proto.UpdateUserPasswordDataRequest) (*empty.Empty, error) {

	action := domain.GetAction(1)
	err := da.accessor.UpdateUserPasswordData(ctx, &domain.EncryptedUserPasswordData{
		Hint:    req.Hint,
		Content: req.Content,
	})

	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	return nil, nil
}
