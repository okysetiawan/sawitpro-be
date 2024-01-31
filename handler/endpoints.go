package handler

import (
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/shared/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) UsersLogin(ctx echo.Context) error {
	var (
		req  = generated.UserLoginRequest{}
		res  = generated.UserLoginResponse{}
		rctx = ctx.Request().Context()
	)

	err := ctx.Bind(&req)
	if err != nil {
		return err
	}

	if err := util.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return err
	}

	output, err := s.Repository.FindUserByPhoneNumber(rctx, req.PhoneNumber)
	if err != nil {
		return err
	}
	res.Id = output.Id

	res.AccessToken, err = s.jwt.CreateAccessToken(output.Id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (s *Server) GetUsersProfile(ctx echo.Context, _ generated.GetUsersProfileParams) error {
	var (
		rctx        = ctx.Request().Context()
		userId, err = util.GetUserIDFromContext(rctx)
	)

	if err != nil {
		return err
	}

	user, err := s.Repository.FindUserById(rctx, userId)
	if err != nil {
		return err
	}

	res := &generated.GetProfileResponse{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (s *Server) UpdateUsersProfile(ctx echo.Context, _ generated.UpdateUsersProfileParams) error {
	var (
		req         = generated.UpdateUsersProfileJSONRequestBody{}
		res         = generated.RegisterProfileResponse{}
		rctx        = ctx.Request().Context()
		userId, err = util.GetUserIDFromContext(rctx)
	)

	if err != nil {
		return err
	}

	err = ctx.Bind(&req)
	if err != nil {
		return err
	}

	if req.PhoneNumber == nil && req.FullName == nil {
		return ctx.JSON(http.StatusNoContent, nil)
	}
	if req.PhoneNumber != nil {
		if err := util.ValidatePhoneNumber(*req.PhoneNumber); err != nil {
			return err
		}

		output, err := s.Repository.FindUserByPhoneNumber(rctx, *req.PhoneNumber)
		if err != nil {
			return err
		}
		if output.Id != 0 && output.Id != userId {
			return echo.NewHTTPError(http.StatusConflict, "phone number already existed")
		}
	}

	_, err = s.Repository.FindUserById(rctx, userId)
	if err != nil {
		return err
	}

	err = s.Repository.UpdateUser(rctx, req.FullName, req.PhoneNumber, userId)
	if err != nil {
		return err
	}

	res.Id = userId

	return ctx.JSON(http.StatusOK, res)
}

func (s *Server) CreateUsersProfile(ctx echo.Context) error {
	var (
		req  = generated.RegisterProfileRequest{}
		res  = generated.RegisterProfileResponse{}
		rctx = ctx.Request().Context()
	)

	err := ctx.Bind(&req)
	if err != nil {
		return err
	}

	if err := util.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return err
	}

	createReq := repository.CreateUserInput{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		Salt:        util.RandomSalt(),
	}
	createReq.Password = util.HashPassword(req.Password + ":" + createReq.Salt)

	output, err := s.Repository.CreateUser(rctx, createReq)
	if err != nil {
		return err
	}

	res.Id = output.Id
	return ctx.JSON(http.StatusOK, res)
}
