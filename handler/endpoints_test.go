package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/shared/jwt"
	"github.com/SawitProRecruitment/UserService/shared/util"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer_UsersLogin(t *testing.T) {
	var (
		ctrl   = gomock.NewController(t)
		router = echo.New()

		repo      = repository.NewMockRepositoryInterface(ctrl)
		jwtSigner = jwt.NewMockSigner(ctrl)
		s         = &Server{
			Repository: repo,
			jwt:        jwtSigner,
		}
	)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		req := generated.UserLoginRequest{
			Password:    "fdafafds",
			PhoneNumber: "+62123132131",
		}
		expectedRes := generated.UserLoginResponse{
			AccessToken: "fdafafdasfasdfafasdf",
			Id:          1,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: req.PhoneNumber,
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber).
			Return(user, nil)

		jwtSigner.EXPECT().CreateAccessToken(user.Id).Return(expectedRes.AccessToken, nil)

		err := s.UsersLogin(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	})

	t.Run("Failed CreateAccessToken", func(t *testing.T) {
		req := generated.UserLoginRequest{
			Password:    "fdafafds",
			PhoneNumber: "+62123132131",
		}
		expectedRes := generated.UserLoginResponse{
			AccessToken: "fdafafdasfasdfafasdf",
			Id:          1,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: req.PhoneNumber,
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber).
			Return(user, nil)

		jwtSigner.EXPECT().CreateAccessToken(user.Id).Return(expectedRes.AccessToken, context.DeadlineExceeded)

		err := s.UsersLogin(ctx)
		assert.Error(t, err)
	})

	t.Run("Failed FindUserByPhoneNumber", func(t *testing.T) {
		req := generated.UserLoginRequest{
			Password:    "fdafafds",
			PhoneNumber: "+62123132131",
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: req.PhoneNumber,
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber).
			Return(user, context.DeadlineExceeded)

		err := s.UsersLogin(ctx)
		assert.Error(t, err)
	})

	t.Run("Failed Bind", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/login", nil)
		ctx := router.NewContext(r, w)

		err := s.UsersLogin(ctx)
		assert.Error(t, err)
	})
}

func TestServer_GetUsersProfile(t *testing.T) {
	var (
		ctrl   = gomock.NewController(t)
		router = echo.New()

		repo      = repository.NewMockRepositoryInterface(ctrl)
		jwtSigner = jwt.NewMockSigner(ctrl)
		s         = &Server{
			Repository: repo,
			jwt:        jwtSigner,
		}
	)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", nil)
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: "+62123132131",
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserById(ctx.Request().Context(), int64(1)).
			Return(user, nil)

		err := s.GetUsersProfile(ctx, generated.GetUsersProfileParams{})
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	})

	t.Run("Failed FindUserById", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", nil)
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: "+62123132131",
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserById(ctx.Request().Context(), int64(1)).
			Return(user, context.DeadlineExceeded)

		err := s.GetUsersProfile(ctx, generated.GetUsersProfileParams{})
		assert.Error(t, err)
	})

	t.Run("Failed Invalid UserID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", nil)
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", "1"))
		ctx := router.NewContext(r, w)

		err := s.GetUsersProfile(ctx, generated.GetUsersProfileParams{})
		assert.Error(t, err)
	})
}

func TestServer_UpdateUsersProfile(t *testing.T) {
	var (
		ctrl   = gomock.NewController(t)
		router = echo.New()

		repo      = repository.NewMockRepositoryInterface(ctrl)
		jwtSigner = jwt.NewMockSigner(ctrl)
		s         = &Server{
			Repository: repo,
			jwt:        jwtSigner,
		}
	)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		phoneNumber := "+6232131"
		req := generated.UpdateProfileRequest{
			PhoneNumber: &phoneNumber,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: phoneNumber,
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserByPhoneNumber(ctx.Request().Context(), *req.PhoneNumber).Return(user, nil)
		repo.EXPECT().FindUserById(ctx.Request().Context(), int64(1)).Return(user, nil)
		repo.EXPECT().UpdateUser(ctx.Request().Context(), req.FullName, req.PhoneNumber, user.Id).Return(nil)

		err := s.UpdateUsersProfile(ctx, generated.UpdateUsersProfileParams{})
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	})

	t.Run("Failed UpdateUser", func(t *testing.T) {
		phoneNumber := "+6232131"
		req := generated.UpdateProfileRequest{
			PhoneNumber: &phoneNumber,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: phoneNumber,
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserByPhoneNumber(ctx.Request().Context(), *req.PhoneNumber).Return(user, nil)
		repo.EXPECT().FindUserById(ctx.Request().Context(), int64(1)).Return(user, nil)
		repo.EXPECT().UpdateUser(ctx.Request().Context(), req.FullName, req.PhoneNumber, user.Id).Return(context.DeadlineExceeded)

		err := s.UpdateUsersProfile(ctx, generated.UpdateUsersProfileParams{})
		assert.Error(t, err)
	})

	t.Run("Failed FindUserById", func(t *testing.T) {
		phoneNumber := "+6232131"
		req := generated.UpdateProfileRequest{
			PhoneNumber: &phoneNumber,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: phoneNumber,
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserByPhoneNumber(ctx.Request().Context(), *req.PhoneNumber).Return(user, nil)
		repo.EXPECT().FindUserById(ctx.Request().Context(), int64(1)).Return(user, context.DeadlineExceeded)

		err := s.UpdateUsersProfile(ctx, generated.UpdateUsersProfileParams{})
		assert.Error(t, err)
	})

	t.Run("Failed Conflict PhoneNumber", func(t *testing.T) {
		phoneNumber := "+6232131"
		req := generated.UpdateProfileRequest{
			PhoneNumber: &phoneNumber,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          2,
			FullName:    "Sulaiman",
			PhoneNumber: phoneNumber,
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserByPhoneNumber(ctx.Request().Context(), *req.PhoneNumber).Return(user, nil)

		err := s.UpdateUsersProfile(ctx, generated.UpdateUsersProfileParams{})
		assert.Error(t, err)
	})

	t.Run("Failed FindUserByPhoneNumber", func(t *testing.T) {
		phoneNumber := "+6232131"
		req := generated.UpdateProfileRequest{
			PhoneNumber: &phoneNumber,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		user := repository.User{
			Id:          1,
			FullName:    "Sulaiman",
			PhoneNumber: phoneNumber,
			Password:    "fddasfa",
			Salt:        "fdasfsa",
		}
		repo.EXPECT().FindUserByPhoneNumber(ctx.Request().Context(), *req.PhoneNumber).Return(user, context.DeadlineExceeded)

		err := s.UpdateUsersProfile(ctx, generated.UpdateUsersProfileParams{})
		assert.Error(t, err)
	})

	t.Run("Failed Prefix PhoneNumber", func(t *testing.T) {
		phoneNumber := "+12123132131"
		req := generated.UpdateProfileRequest{
			PhoneNumber: &phoneNumber,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		err := s.UpdateUsersProfile(ctx, generated.UpdateUsersProfileParams{})
		assert.Error(t, err)
	})

	t.Run("Failed NoContent", func(t *testing.T) {
		req := generated.UpdateProfileRequest{}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		err := s.UpdateUsersProfile(ctx, generated.UpdateUsersProfileParams{})
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusNoContent)
	})

	t.Run("Failed Bind", func(t *testing.T) {
		req := generated.UpdateProfileRequest{
			PhoneNumber: nil,
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile/1", bytes.NewBuffer(buff))
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		err := s.UpdateUsersProfile(ctx, generated.UpdateUsersProfileParams{})
		assert.Error(t, err)
	})
}

func TestServer_CreateUsersProfile(t *testing.T) {
	var (
		ctrl   = gomock.NewController(t)
		router = echo.New()

		repo      = repository.NewMockRepositoryInterface(ctrl)
		jwtSigner = jwt.NewMockSigner(ctrl)
		s         = &Server{
			Repository: repo,
			jwt:        jwtSigner,
		}
	)
	defer ctrl.Finish()
	os.Setenv("TEST", "TRUE")
	defer os.Setenv("TEST", "FALSE")

	t.Run("Success", func(t *testing.T) {
		req := generated.RegisterProfileRequest{
			FullName:    "Sulaiman",
			Password:    "testpassword123",
			PhoneNumber: "+6232131",
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		createReq := repository.CreateUserInput{
			FullName:    req.FullName,
			PhoneNumber: req.PhoneNumber,
			Password:    "hashedpassworduser123123123",
			Salt:        util.RandomSalt(),
		}
		user := repository.CreateUserOutput{Id: 1}
		repo.EXPECT().CreateUser(ctx.Request().Context(), createReq).Return(user, nil)

		err := s.CreateUsersProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	})

	t.Run("Failed CreateUser", func(t *testing.T) {
		req := generated.RegisterProfileRequest{
			FullName:    "Sulaiman",
			Password:    "testpassword123",
			PhoneNumber: "+6232131",
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		createReq := repository.CreateUserInput{
			FullName:    req.FullName,
			PhoneNumber: req.PhoneNumber,
			Password:    "hashedpassworduser123123123",
			Salt:        util.RandomSalt(),
		}
		user := repository.CreateUserOutput{Id: 1}
		repo.EXPECT().CreateUser(ctx.Request().Context(), createReq).Return(user, context.DeadlineExceeded)

		err := s.CreateUsersProfile(ctx)
		assert.Error(t, err)
	})

	t.Run("Failed Prefix PhoneNumber", func(t *testing.T) {
		req := generated.RegisterProfileRequest{
			FullName:    "Sulaiman",
			Password:    "testpassword123",
			PhoneNumber: "+1232131",
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile", bytes.NewBuffer(buff))
		r.Header.Set("Content-Type", echo.MIMEApplicationJSON)
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		err := s.CreateUsersProfile(ctx)
		assert.Error(t, err)
	})

	t.Run("Failed Bind", func(t *testing.T) {
		req := generated.RegisterProfileRequest{
			FullName:    "Sulaiman",
			Password:    "testpassword123",
			PhoneNumber: "+1232131",
		}

		buff, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/users/profile", bytes.NewBuffer(buff))
		r = r.WithContext(context.WithValue(r.Context(), "UserID", int64(1)))
		ctx := router.NewContext(r, w)

		err := s.CreateUsersProfile(ctx)
		assert.Error(t, err)
	})
}
