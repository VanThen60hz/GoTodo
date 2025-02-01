package jwtplugin

import (
	"GoTodo/common"
	"GoTodo/plugin/tokenprovider"
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTPlugin struct {
	name   string
	secret string
}

func NewJWTPlugin(name string) *JWTPlugin {
	return &JWTPlugin{name: name}
}

func (j *JWTPlugin) GetPrefix() string {
	return j.name
}

func (j *JWTPlugin) Get() interface{} {
	return j
}

func (j *JWTPlugin) Name() string {
	return j.name
}

func (j *JWTPlugin) InitFlags() {
	flag.StringVar(&j.secret, "jwt-secret", "default_secret", "JWT Secret Key")
}

func (j *JWTPlugin) Configure() error {
	if j.secret == "default_secret" {
		return fmt.Errorf("JWT secret key is not configured")
	}
	return nil
}

func (j *JWTPlugin) Run() error {
	return nil
}

func (j *JWTPlugin) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

type myClaims struct {
	Payload common.TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}

type token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token
}

func (j *JWTPlugin) SecretKey() string {
	return j.secret
}

func (j *JWTPlugin) Generate(data tokenprovider.TokenPayload, expiry int) (tokenprovider.Token, error) {
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		common.TokenPayload{
			UId:   data.UserId(),
			URole: data.Role(),
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(expiry))),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})

	myToken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	return &token{
		Token:   myToken,
		Expiry:  expiry,
		Created: now,
	}, nil
}

func (j *JWTPlugin) Validate(myToken string) (tokenprovider.TokenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil || !res.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}

	claims, ok := res.Claims.(*myClaims)
	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}

	return claims.Payload, nil
}
