package service_account

import (
	"context"
	"fmt"
	"time"

	src_const "git.ctisoftware.vn/back-end/base/src/const"
	"git.ctisoftware.vn/back-end/base/src/database/collection"
	model_account "git.ctisoftware.vn/back-end/base/src/database/model/account"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AccountAddCommand struct {
	WorkspaceID string
	ActionBy    string
	ActionByID  string

	Name     string
	Gender   int
	Birthday time.Time
	Phone    string
	Avatar   string

	Username string
	Email    string
	Password string

	AppSheetID string

	Type int
}

func (c *AccountAddCommand) Valid(ctx context.Context) (codeErr string) {
	if c.Email == "" && c.Username == "" {
		codeErr = src_const.InvalidErr + src_const.ElementErr_Account + src_const.InternalError + src_const.ServiceErr_Auth
		return codeErr
	}

	return ""
}

func AccountAdd(ctx context.Context, c *AccountAddCommand) (result *model_account.Account, err error) {
	condition := make(map[string]interface{})
	condition["$or"] = []bson.M{
		{"username": c.Username},
		{"email": c.Email},
	}

	cnt, err := collection.Account().Collection().CountDocuments(ctx, condition)
	if err == nil && cnt > 0 {
		err = fmt.Errorf("account is exsited")
		return nil, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	count, id := model_account.GenerateAccountID(ctx)
	result = &model_account.Account{
		ID:        id,
		CountID:   int(count),
		Username:  c.Username,
		Email:     c.Email,
		Password:  string(password),
		Type:      c.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    model_account.StatusActive,
		LogsStatus: []model_account.LogStatus{
			{
				Status:    model_account.StatusActive,
				CreatedAt: time.Now(),
			},
		},
	}

	_, err = collection.Account().Collection().InsertOne(ctx, result)
	if err != nil {
		return nil, err
	}
	return
}
