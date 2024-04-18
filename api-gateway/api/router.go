package api

import (
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"api-gateway/api/handlers"
	v1 "api-gateway/api/handlers/v1"

	// "api-gateway/api/middleware"

	// "api-gateway/api/middleware"
	grpcClients "api-gateway/internal/infrastructure/grpc_service_client"
	redisrepo "api-gateway/internal/infrastructure/repository/redis"
	"api-gateway/internal/pkg/config"
	"api-gateway/internal/usecase/app_version"
	"api-gateway/internal/usecase/event"
	"api-gateway/internal/usecase/refresh_token"
	// "api-gateway/internal/usecase/refresh_token"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Cache          redisrepo.Cache
	Enforcer       *casbin.CachedEnforcer
	Service        grpcClients.ServiceClient
	RefreshToken   refresh_token.RefreshToken
	BrokerProducer event.BrokerProducer
	AppVersion     app_version.AppVersion
}

// NewRoute
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) http.Handler {
	handleOption := &handlers.HandlerOption{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Cache:          option.Cache,
		Enforcer:       option.Enforcer,
		Service:        option.Service,
		RefreshToken:   option.RefreshToken,
		AppVersion:     option.AppVersion,
	}

	router := chi.NewRouter()
	router.Use(chimiddleware.RealIP, chimiddleware.Logger, chimiddleware.Recoverer)
	router.Use(chimiddleware.Timeout(option.ContextTimeout))
	// router.Use(middleware.Tracing)
	router.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Route("/v1", func(r chi.Router) {
		// r.Use(middleware.AuthContext(option.Config.Token.Secret))
		r.Mount("/users", v1.NewUserHandler(handleOption))
	})

	// declare swagger api route
	router.Get("/swagger/*", httpSwagger.Handler())
	return router
}
