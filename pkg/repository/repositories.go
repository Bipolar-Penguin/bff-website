package repository

import (
	"context"
	"time"

	"github.com/Bipolar-Penguin/bff-website/pkg/repository/trading_bid"
	"github.com/Bipolar-Penguin/bff-website/pkg/repository/trading_session"
	"github.com/Bipolar-Penguin/bff-website/pkg/repository/user"
	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	tradingDatabase string = "trading"

	userCollection           string = "users"
	tradingSessionCollection string = "trading_sessions"
	tradingBids              string = "trading_bids"
)

type Repositories struct {
	User           user.Repository
	TradingSession trading_session.Repository
	TradingBid     trading_bid.Repository
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
	r.TradingSession = trading_session.NewMongoRepository(client.Database(tradingDatabase).Collection(tradingSessionCollection))
	r.TradingBid = trading_bid.NewMongoRepository(client.Database(tradingDatabase).Collection(tradingBids))

	return r, nil
}
