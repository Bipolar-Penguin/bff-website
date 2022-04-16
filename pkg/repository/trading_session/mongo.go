package trading_session

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Bipolar-Penguin/bff-website/pkg/domain"
	"github.com/google/uuid"
)

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(coll *mongo.Collection) *mongoRepository {
	r := &mongoRepository{coll}
	r.dummy()
	return r
}

func (m *mongoRepository) Find(sessionID string) (domain.TradingSession, error) {
	var session domain.TradingSession

	var err error

	err = m.coll.FindOne(context.Background(), bson.M{"_id": sessionID}).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return session, nil
	}

	if err != nil {
		return session, err
	}

	return session, nil
}

func (m *mongoRepository) Save(session domain.TradingSession) error {
	var operations []mongo.WriteModel

	if session.ID != "" {
		session.ID = uuid.NewString()
	}
	mdl := mongo.NewReplaceOneModel().
		SetFilter(bson.M{"_id": session.ID}).
		SetUpsert(true).
		SetReplacement(session)

	operations = append(operations, mdl)

	_, err := m.coll.BulkWrite(context.Background(), operations)

	if err != nil {
		return err
	}

	return nil
}

func (m *mongoRepository) FindAll() ([]domain.TradingSession, error) {
	var sessions []domain.TradingSession

	cursor, err := m.coll.Find(
		context.Background(),
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (m *mongoRepository) dummy() error {
	session := domain.TradingSession{
		ID:          "1345497a-1f76-46a2-9561-b5fdf77b722e",
		Status:      "active",
		Title:       `Флорариум с суккулентами "Нежность S"`,
		Description: "ГОСУДАРСТВЕННОЕ БЮДЖЕТНОЕ ПРОФЕССИОНАЛЬНОЕ ОБРАЗОВАТЕЛЬНОЕ УЧРЕЖДЕНИЕ ДЕПАРТАМЕНТА ЗДРАВООХРАНЕНИЯ ГОРОДА МОСКВЫ «МЕДИЦИНСКИЙ КОЛЛЕДЖ № 7»",
		MaxPrice:    15000000,
		Date: struct {
			Start time.Time "json:\"start\" bson:\"start\""
			End   time.Time "json:\"end\" bson:\"end\""
		}{
			Start: time.Now(),
			End:   time.Now().Add(time.Duration(6 * time.Hour)),
		},
		ImageURLs: []string{"https://zakupki.mos.ru/newapi/api/Core/Thumbnail/2119285468/140/140"},
		UserID:    "",
	}

	_, err := m.coll.InsertOne(context.Background(), session)

	return err
}
