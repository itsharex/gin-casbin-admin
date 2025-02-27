//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/wxlbd/gin-casbin-admin/internal/handler"
	"github.com/wxlbd/gin-casbin-admin/internal/repository"
	"github.com/wxlbd/gin-casbin-admin/internal/server"
	"github.com/wxlbd/gin-casbin-admin/internal/service"
	"github.com/wxlbd/gin-casbin-admin/pkg/casbinx"
	"github.com/wxlbd/gin-casbin-admin/pkg/config"
	"github.com/wxlbd/gin-casbin-admin/pkg/gormx"
	"github.com/wxlbd/gin-casbin-admin/pkg/jwtx"
	"github.com/wxlbd/gin-casbin-admin/pkg/log"
	"github.com/wxlbd/gin-casbin-admin/pkg/redisx"
)

var ServerSet = wire.NewSet(server.NewServerHTTP)

var RepositorySet = wire.NewSet(
	repository.NewRepository,
)

var ServiceSet = wire.NewSet(
	service.NewService,
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
)

func NewWire(cfg *config.Config, logger *log.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
		casbinx.New,
		gormx.NewDB,
		redisx.New,
		jwtx.New,
		ServerSet,
		RepositorySet,
		ServiceSet,
		HandlerSet,
	))

	return nil, nil, nil
}
