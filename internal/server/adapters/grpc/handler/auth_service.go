package handler

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

func NewAuthService(aService AuthService) *authService {
	return &authService{
		authService: aService,
	}
}

type authService struct {
	proto.UnimplementedAuthServiceServer
	authService AuthService
}

func (aS *authService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	action := domain.GetAction(0)

	data := &domain.EMailData{
		EMail:    req.Email,
		Password: req.Password,
	}

	sID, err := aS.authService.Login(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}
	resp := &proto.LoginResponse{
		SessionId: string(sID),
	}

	return resp, nil
}

func (aS *authService) PassOTP(ctx context.Context, req *proto.PassOTPRequest) (*proto.AuthResponse, error) {
	action := domain.GetAction(0)
	jwtToken, err := aS.authService.CheckOTP(ctx, domain.SessionID(req.SessionId), req.Password)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	resp := &proto.AuthResponse{
		Token: string(jwtToken),
	}

	return resp, nil
}
