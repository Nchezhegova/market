package models

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/Nchezhegova/market/internal/service/hash"
	"github.com/Nchezhegova/market/internal/service/jwt"
)

type UserModel struct {
	ID   int    `json:"ID"`
	Name string `json:"login"`
	//	Email          string `json:"email"`
	LoyaltyBalance int64  `json:"loyalty"`
	Address        string `json:"address"`
	Password       string `json:"password"`
}

type User interface {
	Add(context.Context) error
	CheckUser(context.Context) bool
	Login(context.Context) error
	CheckToken(context.Context, string) error
}

func (u *UserModel) Add(ctx context.Context) error {
	var err error = nil
	if u.CheckUser(ctx) {
		err = fmt.Errorf("user with the same name already exists")
		return err
	}
	if u.Name == "" || u.Password == "" {
		err = fmt.Errorf("no required parameters")
		return err
	}
	hashpass := base64.StdEncoding.EncodeToString(hash.CalculateHash(u.Password))
	db.AddUser(ctx, u.Name, hashpass)
	return err
}

func (u *UserModel) CheckUser(ctx context.Context) bool {
	if exists := db.CheckUser(ctx, u.Name); exists {
		log.Logger.Info("User exists")
		return exists
	}
	return false
}

func (u *UserModel) Login(ctx context.Context) (string, error) {
	if !u.CheckUser(ctx) {
		err := fmt.Errorf("user does not exist")
		return "", err
	}

	var pass string
	pass, u.ID = db.CheckPassword(ctx, u.Name)
	if pass == base64.StdEncoding.EncodeToString(hash.CalculateHash(u.Password)) {
		token, err := jwt.BuildJWTString(u.ID)
		if err != nil {
			err := fmt.Errorf("problem with token")
			return "", err
		}
		return token, nil
	} else {
		err := fmt.Errorf("wrong password")
		return "", err
	}
}

func (u *UserModel) CheckToken(ctx context.Context, token string) (int, error) {
	if u.ID = jwt.GetUserID(token); u.ID < 0 {
		err := fmt.Errorf("not valid token")
		return 0, err
	}
	return u.ID, nil
}
