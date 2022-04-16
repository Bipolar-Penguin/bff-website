package domain

import "time"

type TradingSession struct {
	ID          string `json:"id" bson:"_id"`
	Status      string `json:"status" bson:"status"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	MaxPrice    int    `json:"max_price" bson:"max_price"`
	Date        struct {
		Start time.Time `json:"start" bson:"start"`
		End   time.Time `json:"end" bson:"end"`
	} `json:"date" bson:"date"`
	ImageURLs []string `json:"image_urls" bson:"image_urls"`
	UserID    string   `json:"user_id" bson:"user_id"`
}
