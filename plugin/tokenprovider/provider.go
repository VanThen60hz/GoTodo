package tokenprovider

import (
	"GoTodo/common"
	"errors"
)

type Provider interface {
	Generate(data TokenPayload, expiry int) (Token, error)
	Validate(token string) (TokenPayload, error)
	SecretKey() string // ko nen de vao day
}

type TokenPayload interface {
	UserId() int
	Role() string
}

type Token interface {
	GetToken() string
}

var (
	ErrNotFound = common.NewCustomError(
		errors.New("token not found"),
		"token not found",
		"ErrNotFound",
	)
	ErrEncodingToken = common.NewCustomError(
		errors.New("error encoding the token"),
		"error encoding the token",
		"ErrEncodingToken",
	)
	ErrInvalidToken = common.NewCustomError(
		errors.New("invalid token provided"),
		"invalid token provided",
		"ErrInvalidToken",
	)
)
