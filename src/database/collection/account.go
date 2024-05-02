package collection

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.ctisoftware.vn/back-end/base/config"
	"git.ctisoftware.vn/back-end/base/src/utilities"
	"git.ctisoftware.vn/back-end/utilities/data/mongo"
)

const AccountCollection = "accounts"

const (
	AccountIndexUsername = "username"
	AccountIndexEmail    = "email"
	AccountIndexUserID   = "user_id"
)

var (
	_accountCollection        *AccountMongoCollection
	loadAccountRepositoryOnce sync.Once
)

type AccountMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadAccountCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadAccountRepositoryOnce.Do(func() {
		_accountCollection, err = NewAccountMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Account() *AccountMongoCollection {
	if _accountCollection == nil {
		panic("database: like account collection is not initiated")
	}
	return _accountCollection
}

func NewAccountMongoCollection(client *mongo.MongoDB, databaseName string) (*AccountMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewAccountMongoCollection] client nil pointer")
	}
	repo := &AccountMongoCollection{
		client:         client,
		collectionName: AccountCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	repo.SetIndex()
	return repo, nil
}

func (repo *AccountMongoCollection) SetIndex() {
	col := repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)

	indexes := []mongoDriver.IndexModel{
		{
			Keys: bson.M{
				AccountIndexUsername: 1,
			},
			Options: &options.IndexOptions{
				Name:   utilities.SetString(AccountIndexUsername),
				Unique: utilities.SetBool(true),
			},
		},
		{
			Keys: bson.M{
				AccountIndexEmail: 1,
			},
			Options: &options.IndexOptions{
				Name: utilities.SetString(AccountIndexEmail),
			},
		},
		{
			Keys: bson.M{
				AccountIndexUserID: 1,
			},
			Options: &options.IndexOptions{
				Name: utilities.SetString(AccountIndexUserID),
			},
		},
	}

	if !repo.needIndex(col) {
		return
	}

	col.Indexes().CreateMany(context.Background(), indexes)
}

func (repo *AccountMongoCollection) needIndex(col *mongoDriver.Collection) bool {
	keyIndexes := []string{
		AccountIndexUsername,
		AccountIndexEmail,
		AccountIndexUserID,
	}

	listIndexes, err := col.Indexes().ListSpecifications(context.Background())
	if err != nil {
		return true
	}
	indexed := make([]string, 0)
	for i := 0; i < len(listIndexes); i++ {
		indexed = append(indexed, listIndexes[i].Name)
	}

	for i := 0; i < len(keyIndexes); i++ {
		if !utilities.StringIntArray(keyIndexes[i], indexed) {
			return true
		}
	}

	return false
}

func (repo *AccountMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
