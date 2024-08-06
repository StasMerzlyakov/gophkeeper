package handler

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

func (h *handler) GetHelloData(ctx context.Context) (*domain.HelloData, error) {

	resp, err := h.dataAccessor.Hello(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: get hello data err ", err)
	}
	data := &domain.HelloData{
		HelloEncrypted:     resp.HelloEncrypted,
		MasterPasswordHint: resp.MasterPasswordHint,
	}

	return data, nil
}

func (h *handler) GetBankCardList(ctx context.Context) ([]domain.EncryptedBankCard, error) {

	list, err := h.dataAccessor.GetBankCardList(ctx, nil)
	if err != nil {
		action := domain.GetAction(1)
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	resp := []domain.EncryptedBankCard{}
	for _, card := range list.Cards {
		resp = append(resp, domain.EncryptedBankCard{
			Number:  card.Number,
			Content: card.Content,
		})
	}
	return resp, nil
}

func (h *handler) CreateBankCard(ctx context.Context, req *domain.EncryptedBankCard) error {
	encrBanckCard := &proto.CreateBankCardRequest{
		Number:  req.Number,
		Content: req.Content,
	}
	_, err := h.dataAccessor.CreateBankCard(ctx, encrBanckCard)
	if err != nil {
		action := domain.GetAction(1)
		return fmt.Errorf("%v err - %w", action, err)
	}
	return nil
}

func (h *handler) DeleteBankCard(ctx context.Context, number string) error {
	_, err := h.dataAccessor.DeleteBankCard(ctx, &proto.DeleteBankCardRequest{
		Number: number,
	})

	if err != nil {
		action := domain.GetAction(1)
		return fmt.Errorf("%v err - %w", action, err)
	}
	return nil
}

func (h *handler) UpdateBankCard(ctx context.Context, card *domain.EncryptedBankCard) error {

	_, err := h.dataAccessor.UpdateBankCard(ctx, &proto.UpdateBankCardRequest{
		Number:  card.Number,
		Content: card.Content,
	})

	if err != nil {
		action := domain.GetAction(1)
		return fmt.Errorf("%v err - %w", action, err)
	}

	return nil
}

func (h *handler) GetUserPasswordDataList(ctx context.Context) ([]domain.EncryptedUserPasswordData, error) {

	list, err := h.dataAccessor.GetUserPasswordDataList(ctx, nil)
	if err != nil {
		action := domain.GetAction(1)
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	resp := []domain.EncryptedUserPasswordData{}
	for _, card := range list.Datas {
		resp = append(resp, domain.EncryptedUserPasswordData{
			Hint:    card.Hint,
			Content: card.Content,
		})
	}
	return resp, nil
}

func (h *handler) CreateUserPasswordData(ctx context.Context, data *domain.EncryptedUserPasswordData) error {

	req := &proto.CreateUserPasswordDataRequest{
		Hint:    data.Hint,
		Content: data.Content,
	}

	_, err := h.dataAccessor.CreateUserPasswordData(ctx, req)
	if err != nil {
		action := domain.GetAction(1)
		return fmt.Errorf("%v err - %w", action, err)
	}
	return nil
}

func (h *handler) DeleteUserPasswordData(ctx context.Context, hint string) error {

	_, err := h.dataAccessor.DeleteUserPasswordData(ctx, &proto.DeleteUserPasswordDataRequest{
		Hint: hint,
	})

	if err != nil {
		action := domain.GetAction(1)
		return fmt.Errorf("%v err - %w", action, err)
	}
	return nil
}

func (h *handler) UpdateUserPasswordData(ctx context.Context, data *domain.EncryptedUserPasswordData) error {

	_, err := h.dataAccessor.UpdateUserPasswordData(ctx, &proto.UpdateUserPasswordDataRequest{
		Hint:    data.Hint,
		Content: data.Content,
	})

	if err != nil {
		action := domain.GetAction(1)
		return fmt.Errorf("%v err - %w", action, err)
	}

	return nil
}
