package handler

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

func NewRegHandler(registrator Registrator) *regHandler {
	return &regHandler{
		registrator: registrator,
	}
}

type regHandler struct {
	proto.UnimplementedRegistrationServiceServer
	registrator Registrator
}

func (rh *regHandler) CheckEMail(ctx context.Context, req *proto.CheckEMailRequest) (*proto.CheckEMailResponse, error) {

	status, err := rh.registrator.GetEMailStatus(ctx, req.Email)
	if err != nil {
		return nil, WrapErr(err)
	}

	resp := &proto.CheckEMailResponse{}
	switch status {
	case domain.EMailAvailable:
		resp.Status = proto.CheckEMailResponse_AVAILABLE
	case domain.EMailBusy:
		resp.Status = proto.CheckEMailResponse_BUSY
	default:
		return nil, fmt.Errorf("%w - unknown email status", domain.ErrServerInternal)
	}
	return resp, nil
}

func (rh *regHandler) Registrate(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	data := &domain.EMailData{
		EMail:    req.Email,
		Password: req.Password,
	}
	sID, err := rh.registrator.Registrate(ctx, data)

	if err != nil {
		return nil, WrapErr(err)
	}

	resp := &proto.RegistrationResponse{
		SessionId: string(sID),
	}

	return resp, nil
}
func (rh *regHandler) PassOTP(ctx context.Context, req *proto.PassOTPRequest) (*proto.PassOTPResponse, error) {

	sID, err := rh.registrator.PassOTP(ctx, domain.SessionID(req.SessionId), req.Password)
	if err != nil {
		return nil, WrapErr(err)
	}

	resp := &proto.PassOTPResponse{
		SessionId: string(sID),
	}

	return resp, nil
}
func (rh *regHandler) SetMasterKey(ctx context.Context, req *proto.MasterKeyRequest) (*proto.MasterKeyResponse, error) {

	mKeyData := &domain.MasterKeyData{
		EncryptedMasterKey: req.EncryptedMasterKey,
		MasterKeyHint:      req.MasterKeyPassHint,
		HelloEncrypted:     req.HelloEncrypted,
	}

	err := rh.registrator.InitMasterKey(ctx, domain.SessionID(req.SessionId), mKeyData)
	if err != nil {
		return nil, WrapErr(err)
	}

	return nil, nil
}
