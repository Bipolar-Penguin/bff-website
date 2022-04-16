package repository

import (
	"context"
	"time"

	"github.com/Bipolar-Penguin/bff-website/pkg/repository/user"
	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	tradingDatabase string = "trading"

	userCollection string = "users"
)

type Repositories struct {
	User user.Repository
}

func MakeRepositories(mongoURL string, logger log.Logger) (*Repositories, error) {
	var r = new(Repositories)

	var err error

	clientOpts := options.Client().ApplyURI(mongoURL)

	clientOpts.SetServerSelectionTimeout(30 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		logger.Log("error", err)
		return nil, err
	}

	r.User = user.NewMongoRepository(client.Database(tradingDatabase).Collection(userCollection))

	return r, nil
}
