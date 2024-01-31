// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type (
	CreateUserInput struct {
		FullName       string
		PhoneNumber    string
		Password, Salt string
	}

	CreateUserOutput struct {
		Id int64
	}
)

type User struct {
	Id             int64
	FullName       string
	PhoneNumber    string
	Password, Salt string
}
