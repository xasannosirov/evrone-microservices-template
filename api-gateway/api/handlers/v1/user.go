package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	errorapi "api-gateway/api/errors"
	"api-gateway/api/handlers"
	"api-gateway/api/models"

	pbu "api-gateway/genproto/user_service"
	grpcClient "api-gateway/internal/infrastructure/grpc_service_client"
	"api-gateway/internal/pkg/config"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"go.uber.org/zap"
)

type contentHandler struct {
	handlers.BaseHandler
	logger   *zap.Logger
	config   *config.Config
	service  grpcClient.ServiceClient
	enforcer *casbin.CachedEnforcer
}

func NewUserHandler(option *handlers.HandlerOption) http.Handler {
	handler := contentHandler{
		logger:   option.Logger,
		config:   option.Config,
		service:  option.Service,
		enforcer: option.Enforcer,
	}

	handler.Cache = option.Cache
	handler.Client = option.Service
	handler.Config = option.Config

	policies := [][]string{
		{"unauthorized", "/v1/users/create", "POST"},
		{"unauthorized", "/v1/users/update", "PUT"},
		{"unauthorized", "/v1/users/delete/{id}", "DELETE"},
		{"unauthorized", "/v1/users/get/{id}", "GET"},
		{"unauthorized", "/v1/users/list", "GET"},
	}
	for _, policy := range policies {
		_, err := option.Enforcer.AddPolicy(policy)
		if err != nil {
			option.Logger.Error("error during investor enforcer add policies", zap.Error(err))
		}
	}

	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		// auth
		// r.Use(middleware.Authorizer(option.Enforcer, option.Logger))

		// users
		r.Post("/create", handler.CreateUser())
		r.Put("/update", handler.UpdateUser())
		r.Delete("/delete/{id}", handler.DeleteUser())
		r.Get("/get/{id}", handler.GetUserByID())
		r.Get("/list", handler.ListUser())

	})

	return router
}


func (h *contentHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var serverReq pbu.User

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		if err := json.Unmarshal(bytes, &serverReq); err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		resp, err := h.service.UserService().Create(ctx, &serverReq)
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		render.JSON(w, r, resp.Id)
	}
}

func (h *contentHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var serverReq pbu.User

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		if err := json.Unmarshal(bytes, &serverReq); err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		resp, err := h.service.UserService().Update(ctx, &serverReq)
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		render.JSON(w, r, resp)
	}
}

func (h *contentHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		GUID := chi.URLParam(r, "id")

		_, err := h.service.UserService().Delete(ctx, &pbu.GetUserRequest{
			Id: GUID,
		})

		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		render.JSON(w, r, models.DelResp{Status: true})
	}
}

func (h *contentHandler) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		GUID := chi.URLParam(r, "id")

		res, err := h.service.UserService().Get(ctx, &pbu.GetUserRequest{
			Id: GUID,
		})

		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		render.JSON(w, r, models.User{
			Id:        res.Id,
			FirstName: res.FirstName,
			LastName:  res.LastName,
			Username:  res.Username,
			Email:     res.Email,
			Password:  res.Password,
			Bio:       res.Bio,
			Website:   res.Website,
		})
	}
}

func (h *contentHandler) ListUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		page := r.URL.Query().Get("page")
		pageToInt, err := strconv.Atoi(page)
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}
		limit := r.URL.Query().Get("limit")
		limitToInt, err := strconv.Atoi(limit)
		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		res, err := h.service.UserService().GetAll(ctx, &pbu.GetAllUserRequest{
			Page:  int64(pageToInt),
			Limit: int64(limitToInt),
		})

		if err != nil {
			render.Render(w, r, errorapi.Error(err))
			return
		}

		var users models.UserList
		for _, u := range res.AllUsers {
			temp := &models.User{
				Id:        u.Id,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Username:  u.Username,
				Email:     u.Email,
				Password:  u.Password,
				Bio:       u.Bio,
				Website:   u.Website,
			}
			users.Users = append(users.Users, temp)

		}
		render.JSON(w, r, users)
	}
}
