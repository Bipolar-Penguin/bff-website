package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Bipolar-Penguin/bff-website/pkg/domain"
	"github.com/Bipolar-Penguin/bff-website/pkg/repository"
	"github.com/Bipolar-Penguin/bff-website/pkg/transport/amqp"
	"github.com/golang-jwt/jwt"
)

const (
	salt       = "foobar"
	tokenTTL   = 999 * time.Hour
	signingKey = "foobar"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

type Service struct {
	rep    *repository.Repositories
	broker *amqp.RabbitBroker
}

func NewService(rep *repository.Repositories, broker *amqp.RabbitBroker) *Service {
	return &Service{rep, broker}
}

// Trading bids features
func (s *Service) GetTradingBids(tradingSessionID string) ([]domain.TradingBid, error) {
	return s.rep.TradingBid.FindMany(tradingSessionID)
}

func (s *Service) MakeTradingBid(tradingSessionID, userID string) error {
	tradingSession, err := s.rep.TradingSession.Find(tradingSessionID)
	if err != nil {
		return err
	}
	bids, err := s.rep.TradingBid.FindMany(tradingSessionID)
	if err != nil {
		return err
	}

	if len(bids) == 0 {
		var newBid domain.TradingBid

		newBid.TradingSessionID = tradingSessionID
		newBid.UserID = userID
		newBid.Date = time.Now()
		newBid.Bid = tradingSession.MaxPrice - int(float64(tradingSession.MaxPrice)*0.01)
		fmt.Println(newBid)

		if newBid.Bid < 1 {
			return errors.New("cannot make bid: min price reached")
		}

		if err := s.rep.TradingBid.Save(newBid); err != nil {
			return err
		}

		s.broker.PublishEvent(domain.Event{
			GUID:    newBid.UserID,
			Action:  "update",
			Amount:  newBid.Bid,
			EventID: newBid.TradingSessionID,
		})
		return nil
	}

	if bids[0].UserID == userID {
		return errors.New("cannot make bid: you already make a bid")
	}

	newBid := bids[0]

	newBid.TradingSessionID = tradingSessionID
	newBid.UserID = userID
	newBid.Date = time.Now()
	newBid.Bid = newBid.Bid - int(float64(tradingSession.MaxPrice)*0.01)
	fmt.Println(newBid)

	if newBid.Bid < 1 {
		return errors.New("cannot make bid: min price reached")
	}

	if err := s.rep.TradingBid.Save(newBid); err != nil {
		return err
	}

	s.broker.PublishEvent(domain.Event{
		GUID:    newBid.UserID,
		Action:  "update",
		Amount:  newBid.Bid,
		EventID: newBid.TradingSessionID,
	})
	return nil
}

// Trading session features
func (s *Service) GetSessions() ([]domain.TradingSession, error) {
	return s.rep.TradingSession.FindAll()
}

func (s *Service) ExtendSession(sessionID string) error {
	session, err := s.rep.TradingSession.Find(sessionID)
	if err != nil {
		return err
	}

	session.Date.End.Add(time.Duration(5 * time.Minute))

	err = s.rep.TradingSession.Save(session)

	return nil
}

// User features
func (s *Service) SaveUser(user domain.User) (domain.User, error) {
	return s.rep.User.Save(user)
}

func (s *Service) Authenticate(authHeader string) (string, error) {
	token, err := jwt.ParseWithClaims(authHeader, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func (s *Service) GenerateToken(userID string) (string, error) {

	user, err := s.rep.User.Find(userID)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}
