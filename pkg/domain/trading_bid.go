package domain

import "time"

type TradingBid struct {
	TradingSessionID string    `json:"trading_session_id" bson:"trading_session_id"`
	UserID           string    `json:"user_id" bson:"user_id"`
	Bid              int       `json:"bid" bson:"bid"`
	Date             time.Time `json:"date" bson:"date"`
}
