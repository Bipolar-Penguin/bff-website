package trading_bid

import "github.com/Bipolar-Penguin/bff-website/pkg/domain"

type Repository interface {
	FindMany(tradingSessionID string) ([]domain.TradingBid, error)
	Save(domain.TradingBid) error
}
