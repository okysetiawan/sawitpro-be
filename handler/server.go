package handler

import (
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/shared/jwt"
)

type Server struct {
	Repository repository.RepositoryInterface
	jwt        jwt.Signer
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	JWT        jwt.Signer
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		jwt:        opts.JWT,
	}
}
