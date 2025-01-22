package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/wxlbd/gin-casbin-admin/internal/handler"
	"github.com/wxlbd/gin-casbin-admin/pkg/jwtx"
	"github.com/wxlbd/gin-casbin-admin/pkg/log"
)

type service struct {
	user    handler.UserService
	role    handler.RoleService
	menu    handler.MenuService
	dict    handler.DictService
	captcha handler.CaptchaService
}

func NewService(logger *log.Logger, repo Repository, enforcer *casbin.Enforcer, jwt *jwtx.JWT, redisClient *redis.Client) handler.Service {
	return &service{
		user:    NewUserService(logger, repo, jwt),
		role:    NewRoleService(repo, enforcer),
		menu:    NewMenuService(repo, enforcer),
		dict:    NewDictService(logger, repo),
		captcha: NewCaptchaService(redisClient),
	}
}

func (s *service) User() handler.UserService {
	return s.user
}

func (s *service) Role() handler.RoleService {
	return s.role
}

func (s *service) Menu() handler.MenuService {
	return s.menu
}

func (s *service) Dict() handler.DictService {
	return s.dict
}

func (s *service) Captcha() handler.CaptchaService {
	return s.captcha
}
