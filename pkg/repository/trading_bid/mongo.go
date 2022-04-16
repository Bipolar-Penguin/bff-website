package trading_bid

import (
	"context"

	"github.com/Bipolar-Penguin/bff-website/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(coll *mongo.Collection) *mongoRepository {
	return &mongoRepository{coll}
}

func (m *mongoRepository) FindMany(tradingSessionID string) ([]domain.TradingBid, error) {
	var bids []domain.TradingBid

	var opts options.FindOptions
	opts.SetSort(bson.M{"bid": 1})

	cursor, err := m.coll.Find(context.Background(), bson.M{"trading_session_id": tradingSessionID}, &opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &bids); err != nil {
		return nil, err
	}

	return bids, nil
}

func (m *mongoRepository) Save(bid domain.TradingBid) error {
	_, err := m.coll.InsertOne(context.Background(), bid)
	return err
}
