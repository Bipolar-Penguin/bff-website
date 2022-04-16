package trading_session

import "github.com/Bipolar-Penguin/bff-website/pkg/domain"

type Repository interface {
	FindAll() ([]domain.TradingSession, error)
	Find(sessionID string) (domain.TradingSession, error)
	Save(domain.TradingSession) error
}
