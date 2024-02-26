package ginrouter

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"try-go-clickhouse/internal/model"
	"try-go-clickhouse/internal/util/envconf"
)

func New(repo model.Repository, env *envconf.Spec, log *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	NewHandler(repo, r, env, log)

	return r
}

type Handler struct {
	Router *gin.Engine
	Repo   model.Repository
	Env    *envconf.Spec
	Log    *zap.Logger
}

func NewHandler(repo model.Repository, r *gin.Engine, env *envconf.Spec, log *zap.Logger) {
	h := &Handler{
		Router: r,
		Repo:   repo,
		Env:    env,
		Log:    log,
	}
	r.GET("/health", h.HealthCheck())

	r.GET("/users", h.ListUsers())
	r.GET("/users/:id", h.GetUser())
	r.POST("/users", h.AddUser())
}

func (h *Handler) HealthCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	}
}

func (h *Handler) ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		users, err := h.Repo.ListUsers(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func (h *Handler) AddUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user model.User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		user, err = h.Repo.AddUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (h *Handler) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		user, err := h.Repo.GetUser(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, nil)
				return
			}
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
