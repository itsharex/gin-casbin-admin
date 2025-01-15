package server

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/nunu-casbin-admin/internal/handler"
	"github.com/wxlbd/nunu-casbin-admin/internal/middleware"
	"github.com/wxlbd/nunu-casbin-admin/internal/service"
	"github.com/wxlbd/nunu-casbin-admin/pkg/config"
	"github.com/wxlbd/nunu-casbin-admin/pkg/jwtx"
	"github.com/wxlbd/nunu-casbin-admin/pkg/log"
)

func NewServerHTTP(
	cfg *config.Config,
	logger *log.Logger,
	jwt *jwtx.JWT,
	handler *handler.Handler,
	enforcer *casbin.Enforcer,
	svc service.Service,
) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(
		middleware.CORSMiddleware(),
		middleware.RequestLogger(logger),
	)
	api := r.Group("api")
	{
		auth := api.Group("auth")
		// 完全公开的接口
		auth.POST("/login", handler.User().Login)
		auth.POST("/refresh-token", handler.User().RefreshToken)

		// 需要JWT认证的接口
		profile := api.Group("profile").Use(middleware.JWTAuth(jwt))
		{
			profile.GET("", handler.User().Current)
			profile.GET("menus", handler.Menu().GetUserMenus)
			// profile.GET("/menu/tree", handler.Menu().GetMenuTree)
			profile.GET("roles", handler.User().GetCurrentUserRoles)
		}

		// 需要完整权限控制的接口
		authorized := api.Group("")
		authorized.Use(
			middleware.JWTAuth(jwt),
			middleware.CasbinMiddleware(enforcer, logger, svc),
		)
		// 权限控制
		{
			permission := authorized.Group("permission")

			// 用户管理 system:user:xxx
			userGroup := permission.Group("user")
			{
				userGroup.GET("", handler.User().List)                         // permission:user:list
				userGroup.POST("", handler.User().Create)                      // permission:user:create
				userGroup.PUT("/:id", handler.User().Update)                   // permission:user:update
				userGroup.DELETE("/:ids", handler.User().Delete)               // permission:user:delete
				userGroup.GET("/:id", handler.User().Detail)                   // permission:user:detail
				userGroup.GET("/:id/roles", handler.User().GerUserRoles)       // permission:user:get:roles
				userGroup.PATCH(":id/password", handler.User().UpdatePassword) // permission:user:set:password
				userGroup.PUT(":id/roles", handler.User().AssignRoles)         // permission:user:set:roles
			}

			// 角色管理 permission:role:xxx
			roleGroup := permission.Group("role")
			{
				roleGroup.GET("", handler.Role().List)                        // permission:role:list
				roleGroup.POST("", handler.Role().Create)                     // permission:role:create
				roleGroup.PUT("/:id", handler.Role().Update)                  // permission:role:update
				roleGroup.DELETE("/:ids", handler.Role().Delete)              // permission:role:delete
				roleGroup.GET("/:id", handler.Role().Detail)                  // permission:role:detail
				roleGroup.GET("/:id/menus", handler.Role().GetPermittedMenus) // permission:role:get:menus
				roleGroup.PUT("/:id/menus", handler.Role().AssignMenus)       // permission:role:set:menus
			}

			// 菜单管理 permission:menu:xxx
			menuGroup := permission.Group("menu")
			{
				// menuGroup.GET("", handler.Menu().List)             // permission:menu:list
				menuGroup.POST("", handler.Menu().Create)          // permission:menu:create
				menuGroup.PUT("/:id", handler.Menu().Update)       // permission:menu:update
				menuGroup.DELETE("/:ids", handler.Menu().Delete)   // permission:menu:delete
				menuGroup.GET("/tree", handler.Menu().GetMenuTree) // permission:menu:tree
			}
		}
	}

	return r
}
