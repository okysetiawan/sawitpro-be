package main

import (
	"github.com/SawitProRecruitment/UserService/middleware"
	"github.com/SawitProRecruitment/UserService/shared/jwt"
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Auth("POST:/v1/users/login", "POST:/v1/users/profile"))
	e.HTTPErrorHandler = e.DefaultHTTPErrorHandler

	generated.RegisterHandlers(e, server)

	routes := e.Routes()
	e.Logger.Infof("routes: %v", routes)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	var dbDsn = os.Getenv("DATABASE_URL")

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	jwtSigner := jwt.GetSigner()

	opts := handler.NewServerOptions{
		Repository: repo,
		JWT:        jwtSigner,
	}
	return handler.NewServer(opts)
}
