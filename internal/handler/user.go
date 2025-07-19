package handler

import (
	"codematic/internal/domain/user"
)

type User struct {
	service user.Service
	env     *Environment
}

func (h *User) Init(basePath string, env *Environment) error {
	h.env = env

	h.service = user.NewService(
		env.DB,
		env.JWTManager,
		env.Logger,
	)

	return nil

}
