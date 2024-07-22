package handler

import (
	"context"

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

	data := &domain.EMailData{
		EMail:    req.Email,
		Password: req.Password,
	}

	sID, err := aS.authService.Login(ctx, data)
	if err != nil {
		return nil, WrapErr(err)
	}
	resp := &proto.LoginResponse{
		SessionId: string(sID),
	}

	return resp, nil
}

func (aS *authService) PassOTP(ctx context.Context, req *proto.PassOTPRequest) (*proto.AuthResponse, error) {
	jwtToken, err := aS.authService.CheckOTP(ctx, domain.SessionID(req.SessionId), req.Password)
	if err != nil {
		return nil, WrapErr(err)
	}

	resp := &proto.AuthResponse{
		Token: string(jwtToken),
	}

	return resp, nil
}
