package database

import (
	"context"
	"time"

	"git.ctisoftware.vn/back-end/utilities/data/mongo"

	"git.ctisoftware.vn/back-end/base/config"
	src_const "git.ctisoftware.vn/back-end/base/src/const"
	"git.ctisoftware.vn/back-end/base/src/database/collection"
)

func ConnectDatabse(ctx context.Context) error {
	var mongoClient *mongo.MongoDB
	var err error
	numberRetry := config.Get().NumberRetry
	if numberRetry == 0 {
		numberRetry = src_const.DEFAULTNUMBERRETRY
	}

	for i := 1; i <= config.Get().NumberRetry; i++ {
		mongoClient, err = mongo.NewMongoDBFromUrl(ctx, config.Get().MongoURL, time.Second*10)
		if err != nil {
			if i == config.Get().NumberRetry {
				return err
			}
			time.Sleep(10 * time.Second)
		}

		if mongoClient != nil {
			break
		}
	}

	if err := collection.LoadAccountCollectionMongo(mongoClient); err != nil {
		return err
	}
	return nil
}
