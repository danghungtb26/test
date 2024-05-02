package model_account

import (
	"context"
	"fmt"

	"git.ctisoftware.vn/back-end/base/src/database/collection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenerateAccountID(ctx context.Context) (int, string) {
	option := options.FindOne()
	option.SetSort(bson.M{"count_id": -1})

	result := Account{}

	err := collection.Account().Collection().FindOne(ctx, nil, option).Decode(&result)
	if err != nil {
		return 0, ""
	}

	return result.CountID, fmt.Sprintf("Account_%d", result.CountID)
}
