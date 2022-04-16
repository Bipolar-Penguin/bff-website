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
	r := &mongoRepository{coll}
	r.dummy()
	return r
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

func (m *mongoRepository) dummy() {
	var users = []domain.User{
		{
			ID:               "a",
			Username:         "Поставщик А",
			Region:           "Санкт-Петербург",
			OrganizationType: "individual",
			Role:             "supplier",
			Contacts: struct {
				Email       string "json:\"email\" bson:\"email\""
				PhoneNumber string "json:\"phone_number\" bson:\"phone_number\""
				TelegramID  string "json:\"telegram_id\" bson:\"telegram_id\""
			}{
				Email:       "kuwerin@gmail.com",
				PhoneNumber: "79967726643",
				TelegramID:  "528263453",
			},
			Permissions: struct {
				Email    bool "json:\"email\" bson:\"email\""
				Phone    bool "json:\"phone\" bson:\"phone\""
				Telegram bool "json:\"telegram\" bson:\"telegram\""
				Push     bool "json:\"push\" bson:\"push\""
			}{
				Email:    true,
				Phone:    true,
				Telegram: true,
				Push:     true,
			},
		},
		{
			ID:               "b",
			Username:         "Поставщик C",
			Region:           "Санкт-Петербург",
			OrganizationType: "individual",
			Role:             "supplier",
			Contacts: struct {
				Email       string "json:\"email\" bson:\"email\""
				PhoneNumber string "json:\"phone_number\" bson:\"phone_number\""
				TelegramID  string "json:\"telegram_id\" bson:\"telegram_id\""
			}{
				Email:       "",
				PhoneNumber: "79817914985",
				TelegramID:  "528569218",
			},
			Permissions: struct {
				Email    bool "json:\"email\" bson:\"email\""
				Phone    bool "json:\"phone\" bson:\"phone\""
				Telegram bool "json:\"telegram\" bson:\"telegram\""
				Push     bool "json:\"push\" bson:\"push\""
			}{
				Email:    false,
				Phone:    true,
				Telegram: true,
				Push:     true,
			},
		},
		{
			ID:               "c",
			Username:         "Поставщик C",
			Region:           "Санкт-Петербург",
			OrganizationType: "individual",
			Role:             "supplier",
			Contacts: struct {
				Email       string "json:\"email\" bson:\"email\""
				PhoneNumber string "json:\"phone_number\" bson:\"phone_number\""
				TelegramID  string "json:\"telegram_id\" bson:\"telegram_id\""
			}{
				Email:       "",
				PhoneNumber: "",
				TelegramID:  "",
			},
			Permissions: struct {
				Email    bool "json:\"email\" bson:\"email\""
				Phone    bool "json:\"phone\" bson:\"phone\""
				Telegram bool "json:\"telegram\" bson:\"telegram\""
				Push     bool "json:\"push\" bson:\"push\""
			}{
				Email:    false,
				Phone:    false,
				Telegram: false,
				Push:     true,
			},
		},
	}

	for _, user := range users {
		m.coll.InsertOne(context.Background(), user)
	}
}
