package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Bipolar-Penguin/bff-website/pkg/domain"
	"github.com/google/uuid"
)

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(coll *mongo.Collection) *mongoRepository {
	return &mongoRepository{coll}
}

func (m *mongoRepository) Find(userID string) (domain.User, error) {
	var user domain.User

	var err error

	err = m.coll.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, nil
	}

	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *mongoRepository) Save(user domain.User) (domain.User, error) {
	var operations []mongo.WriteModel

	if user.ID != "" {
		mdl := mongo.NewReplaceOneModel().
			SetFilter(bson.M{"_id": user.ID}).
			SetUpsert(true).
			SetReplacement(user)

		operations = append(operations, mdl)

		_, err := m.coll.BulkWrite(context.Background(), operations)

		if err != nil {
			return user, err
		}

		return user, nil
	}

	user.ID = uuid.NewString()
	_, err := m.coll.InsertOne(context.Background(), user)
	if err != nil {
		return user, err
	}

	return user, nil
}
