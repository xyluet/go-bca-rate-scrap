package scrap

import "context"

type Exchange struct {
	Buy  float64
	Sell float64
}

type ExchangeRate struct {
	Currency  string
	ERate     Exchange
	TTCounter Exchange
	BankNotes Exchange
}

type Service interface {
	ListExchangeRates(context.Context) ([]*ExchangeRate, error)
}

type Scrapper interface {
	Scrap(context.Context) ([]*ExchangeRate, error)
}

type service struct {
	scrapper Scrapper
}

func NewService(scrapper Scrapper) Service {
	s := service{
		scrapper: scrapper,
	}
	return &s
}

func (s *service) ListExchangeRates(ctx context.Context) ([]*ExchangeRate, error) {
	rates, err := s.scrapper.Scrap(ctx)
	if err != nil {
		return nil, err
	}
	return rates, nil
}
