// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error)
	UpdateUser(ctx context.Context, name, phoneNumber *string, id int64) (err error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (output User, err error)
	FindUserById(ctx context.Context, id int64) (output User, err error)
}
