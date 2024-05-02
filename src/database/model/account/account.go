package model_account

import (
	"time"

	graph_model "git.ctisoftware.vn/back-end/base/src/graph/generated/model"
)

const (
	StatusActive   = 1
	StatusInactive = 2
	StatusBlock    = 3
	StatusBan      = 4
)

const (
	GenderMale   = 1
	GenderFemale = 2
	GenderOther  = 3
)

const (
	AccountTypeRoot      = 1
	AccountTypeSuperUser = 2
	AccountTypeUser      = 3
)

type Account struct {
	ID      string `json:"id" bson:"_id"`
	CountID int    `json:"count_id" bson:"count_id"`

	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Status   int    `json:"status" bson:"status"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Deleted   bool      `json:"deleted" bson:"deleted,omitempty"`

	LogsStatus []LogStatus `json:"logs_status" bson:"logs_status,omitempty"`

	Type int `json:"type" bson:"type"`
}

type LogStatus struct {
	Status    int32     `json:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func (a *Account) ConvertToModelGraph() *graph_model.Account {
	data := graph_model.Account{
		ID:       a.ID,
		UserName: a.Username,
		Email:    a.Email,
	}

	return &data
}
