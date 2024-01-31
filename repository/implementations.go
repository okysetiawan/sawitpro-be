package repository

import (
	"context"
	"fmt"
	"github.com/SawitProRecruitment/UserService/shared/util"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func (r *Repository) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
	var (
		query = "INSERT INTO users (name, phone_number, password, salt) VALUES ($1, $2, $3, $4) RETURNING id"
		args  = []any{input.FullName, input.PhoneNumber, input.Password, input.Salt}
	)

	err = r.Db.QueryRowContext(ctx, query, args...).Scan(&output.Id)
	if err != nil {
		err = util.TransformError(err)
		return
	}

	return
}

func (r *Repository) UpdateUser(ctx context.Context, name, phoneNumber *string, id int64) (err error) {
	var (
		query   = "UPDATE users SET %s WHERE id = ?"
		args    []any
		setList []string
	)

	if name != nil {
		setList = append(setList, "name = ?")
		args = append(args, *name)
	}
	if phoneNumber != nil {
		setList = append(setList, "phone_number = ?")
		args = append(args, *phoneNumber)
	}

	if len(setList) == 0 {
		return nil
	}

	query = fmt.Sprintf(query, strings.Join(setList, ", "))
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	args = append(args, id)

	_, err = r.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return
	}

	return
}

func (r *Repository) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (output User, err error) {
	var (
		query = "SELECT id, name, phone_number, password, salt FROM users WHERE phone_number = $1"
		args  = []any{phoneNumber}
	)

	err = r.Db.QueryRowContext(ctx, query, args...).
		Scan(&output.Id, &output.FullName, &output.PhoneNumber, &output.Password, &output.Salt)
	if err != nil {
		err = util.TransformError(err)
		return
	}

	return
}

func (r *Repository) FindUserById(ctx context.Context, id int64) (output User, err error) {
	var (
		query = "SELECT id, name, phone_number, password, salt FROM users WHERE id = $1"
		args  = []any{id}
	)

	err = r.Db.QueryRowContext(ctx, query, args...).
		Scan(&output.Id, &output.FullName, &output.PhoneNumber, &output.Password, &output.Salt)
	if err != nil {
		return
	}

	if output.Id == 0 {
		return output, echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return
}
