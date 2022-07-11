package werf

import (
	"time"

	"github.com/mayukh42/logx/logx"
	"github.com/mayukh42/werf/lib"
)

const (
	HOT_PRICE_THRESHOLD = 999.99
	HOT_DATE_THRESHOLD  = "7d"
	SIZE_CAT_SMALL      = iota
	SIZE_CAT_MEDIUM
	SIZE_CAT_LARGE
)

type (
	Item struct {
		Name         string  `json:"name"`
		Type         string  `json:"type"`
		Description  string  `json:"desc"`
		Origin       string  `json:"origin"`
		SizeCategory string  `json:"size_cat"`
		Weight       float64 `json:"weight"`
		Price        float64 `json:"price"`
		Expiry       string  `json:"expiry"`
	}

	Cargo struct {
		// simplest warehouse entry - dynamodb
		Id     int       `json:"id"`
		Weight int       `json:"weight"`
		Value  int       `json:"value"`
		Date   time.Time `json:"date"`
		Hot    bool      `json:"hot"`
	}

	ShipManifest struct {
		Items []Item `json:"items"`
		Date  string `json:"date"`
	}
)

func (i *Item) IsHot() bool {
	pricey := i.Price > HOT_PRICE_THRESHOLD
	et, er := time.Parse(logx.TIME_FORMAT, i.Expiry)
	if er != nil {
		// assume expiry is > 10d
		return pricey
	}
	tl := lib.DateAfter(HOT_DATE_THRESHOLD)
	urgent := et.Before(tl)
	return pricey || urgent
}
