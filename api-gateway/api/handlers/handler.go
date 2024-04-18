package handlers

import (
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	"api-gateway/api/middleware"
	grpcClients "api-gateway/internal/infrastructure/grpc_service_client"
	"api-gateway/internal/infrastructure/repository/redis"
	"api-gateway/internal/pkg/config"
	// "api-gateway/internal/pkg/otlp"
	appV "api-gateway/internal/usecase/app_version"
	"api-gateway/internal/usecase/event"
	"api-gateway/internal/usecase/refresh_token"
)

const (
	InvestorToken = "investor"
)

type HandlerOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Enforcer       *casbin.CachedEnforcer
	Cache          redis.Cache
	Service        grpcClients.ServiceClient
	RefreshToken   refresh_token.RefreshToken
	AppVersion     appV.AppVersion
	BrokerProducer event.BrokerProducer
}

type BaseHandler struct {
	Cache  redis.Cache
	Config *config.Config
	Client grpcClients.ServiceClient
}

func (h *BaseHandler) GetAuthData(ctx context.Context) (map[string]string, bool) {
	// tracing
	// ctx, span := otlp.Start(ctx, "handler", "GetAuthData")
	// defer span.End()

	data, ok := ctx.Value(middleware.RequestAuthCtx).(map[string]string)
	return data, ok
}
