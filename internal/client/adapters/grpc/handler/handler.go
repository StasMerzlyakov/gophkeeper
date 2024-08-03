package handler

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func NewHandler(conf *config.ClientConf) (*handler, error) {
	cred, err := loadTLSCredentials(conf.CACert)
	if err != nil {
		return nil, fmt.Errorf("%w can't read CACert - %v", domain.ErrClientInternal, err.Error())
	}

	h := &handler{}

	client, err := grpc.NewClient(conf.ServerAddress,
		grpc.WithTransportCredentials(cred),
		grpc.WithUnaryInterceptor(h.JWTInterceptor),
		grpc.WithStreamInterceptor(h.JWTStreamInterceptor))

	if err != nil {
		return nil, fmt.Errorf("%w can't create grpc client %v", domain.ErrClientInternal, err.Error())
	}

	h.conn = client
	h.pinger = proto.NewPingerClient(client)
	h.loginer = proto.NewAuthServiceClient(client)
	h.dataAccessor = proto.NewDataAccessorClient(client)
	h.registrator = proto.NewRegistrationServiceClient(client)
	h.fileAccessor = proto.NewFileAccessorClient(client)

	return h, nil
}

var _ app.AppServer = (*handler)(nil)

type handler struct {
	conn         *grpc.ClientConn
	pinger       proto.PingerClient
	dataAccessor proto.DataAccessorClient
	loginer      proto.AuthServiceClient
	registrator  proto.RegistrationServiceClient
	fileAccessor proto.FileAccessorClient
	sessionID    string
	jwtToken     string
}

func (h *handler) SetJWTToken(jwtToken string) {
	h.jwtToken = jwtToken
}

func (h *handler) SetSessionID(sessionID string) {
	h.sessionID = sessionID
}

func (h *handler) SessionID() string {
	return h.sessionID
}

func (h *handler) JWTToken() string {
	return h.jwtToken
}

func (h *handler) Start() {
}

func (h *handler) Stop() {
	if err := h.conn.Close(); err != nil {
		panic(err)
	}
}

func (h *handler) Ping(ctx context.Context) error {
	_, err := h.pinger.Ping(ctx, nil)
	return err
}

func (h *handler) CheckEMail(ctx context.Context, email string) (domain.EMailStatus, error) {
	req := &proto.CheckEMailRequest{
		Email: email,
	}

	resp, err := h.registrator.CheckEMail(ctx, req)
	if err != nil {
		return domain.EMailBusy, fmt.Errorf("%w: check email err ", err)
	}

	status := domain.EMailStatus(resp.Status.String())
	switch status {
	case domain.EMailAvailable, domain.EMailBusy:
		return status, nil
	default:
		return domain.EMailStatus(""), fmt.Errorf("%w unknown email status", domain.ErrServerInternal)
	}
}

func (h *handler) Registrate(ctx context.Context, data *domain.EMailData) error {
	req := &proto.RegistrationRequest{
		Email:    data.EMail,
		Password: data.Password,
	}
	resp, err := h.registrator.Registrate(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: registration err ", err)
	}
	h.sessionID = resp.SessionId
	return nil
}

func (h *handler) PassRegOTP(ctx context.Context, otpPass string) error {
	req := &proto.PassOTPRequest{
		SessionId: h.sessionID,
		Password:  otpPass,
	}

	resp, err := h.registrator.PassOTP(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: passOTP err ", err)
	}
	h.sessionID = resp.SessionId
	return nil
}

func (h *handler) InitMasterKey(ctx context.Context, mKey *domain.MasterKeyData) error {
	req := &proto.MasterKeyRequest{
		SessionId:          h.sessionID,
		MasterPasswordHint: mKey.MasterPasswordHint,
		HelloEncrypted:     mKey.HelloEncrypted,
	}

	_, err := h.registrator.SetMasterKey(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: init master key err ", err)
	}
	h.sessionID = ""
	return nil
}

func (h *handler) Login(ctx context.Context, data *domain.EMailData) error {
	req := &proto.LoginRequest{
		Email:    data.EMail,
		Password: data.Password,
	}

	resp, err := h.loginer.Login(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: login err ", err)
	}

	h.sessionID = resp.SessionId
	return nil
}

func (h *handler) PassLoginOTP(ctx context.Context, otpPass string) error {
	req := &proto.PassOTPRequest{
		SessionId: h.sessionID,
		Password:  otpPass,
	}

	resp, err := h.loginer.PassOTP(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: passOTP err ", err)
	}
	h.sessionID = ""
	h.jwtToken = resp.Token
	log := app.GetMainLogger()
	log.Errorf(h.jwtToken)
	return nil
}

func (h *handler) JWTInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	if h.jwtToken != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, domain.AuthorizationMetadataTokenName, h.jwtToken)
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (h *handler) JWTStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if h.jwtToken != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, domain.AuthorizationMetadataTokenName, h.jwtToken)
	}

	return streamer(ctx, desc, cc, method, opts...)
}

func loadTLSCredentials(caFile string) (credentials.TransportCredentials, error) {
	file, err := os.Open(caFile)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	pemServerCA, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}
