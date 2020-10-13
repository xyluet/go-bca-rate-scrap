package colly

import (
	"context"
	"go-scrap/internal/scrap"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Factory interface {
	NewCollector() *colly.Collector
}

type defaultFactory struct{}

func (f *defaultFactory) NewCollector() *colly.Collector {
	return colly.NewCollector()
}

func DefaultFactory() Factory {
	return &defaultFactory{}
}

type scrapper struct {
	factory Factory
}

func NewScrapper(factory Factory) scrap.Scrapper {
	s := scrapper{
		factory: factory,
	}
	return &s
}

func (s *scrapper) Scrap(context.Context) ([]*scrap.ExchangeRate, error) {
	c := s.factory.NewCollector()
	rates := []*scrap.ExchangeRate{}
	c.OnHTML(".kurs-e-rate tbody tr", func(e *colly.HTMLElement) {
		tds := e.DOM.ChildrenFiltered("td")
		rates = append(rates, &scrap.ExchangeRate{
			Currency: tds.Eq(0).Text(),
			ERate: scrap.Exchange{
				Buy:  s.textToFloat(tds.Eq(1).Text()),
				Sell: s.textToFloat(tds.Eq(2).Text()),
			},
			TTCounter: scrap.Exchange{
				Buy:  s.textToFloat(tds.Eq(3).Text()),
				Sell: s.textToFloat(tds.Eq(4).Text()),
			},
			BankNotes: scrap.Exchange{
				Buy:  s.textToFloat(tds.Eq(5).Text()),
				Sell: s.textToFloat(tds.Eq(6).Text()),
			},
		})
	})
	if err := c.Visit("https://www.bca.co.id/en/e-rate"); err != nil {
		return nil, err
	}
	return rates, nil
}

func (s *scrapper) textToFloat(str string) float64 {
	val, _ := strconv.ParseFloat(strings.ReplaceAll(str, ",", ""), 10)
	return val
}
