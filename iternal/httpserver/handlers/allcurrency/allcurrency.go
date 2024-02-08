package allcurrency

import (
	"CurrencyClient/iternal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	IDName                       string    `json:"id"`
	Symbol                       string    `json:"symbol"`
	Name                         string    `json:"name"`
	Image                        string    `json:"image"`
	CurrentPrice                 float64   `json:"current_price"`
	MarketCap                    int64     `json:"market_cap"`
	MarketCapRank                int       `json:"market_cap_rank"`
	FullyDilutedValuation        int64     `json:"fully_diluted_valuation"`
	TotalVolume                  float64   `json:"total_volume"`
	High24H                      float64   `json:"high_24h"`
	Low24H                       float64   `json:"low_24h"`
	PriceChange24H               float64   `json:"price_change_24h"`
	PriceChangePercentage24H     float64   `json:"price_change_percentage_24h"`
	MarketCapChange24H           float64   `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H float64   `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64   `json:"circulating_supply"`
	TotalSupply                  float64   `json:"total_supply"`
	MaxSupply                    float64   `json:"max_supply"`
	Ath                          float64   `json:"ath"`
	AthChangePercentage          float64   `json:"ath_change_percentage"`
	AthDate                      time.Time `json:"ath_date"`
	Atl                          float64   `json:"atl"`
	AtlChangePercentage          float64   `json:"atl_change_percentage"`
	AtlDate                      time.Time `json:"atl_date"`
	Roi                          Roi       `json:"roi"`
	LastUpdated                  time.Time `json:"last_updated"`
}

type Roi struct {
	Times      float64
	Currency   string
	Percentage float64
}

type RequestAllCurrency interface {
	GetByID(id int) (*Request, error)
}

func New(log *slog.Logger, currency RequestAllCurrency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.allcurrency,New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Start parse url for id")

		idString := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idString)
		if err != nil {
			log.Error("Uncorrected id in request ")

			render.JSON(w, r, "Uncorrected id in request")

			return
		}

		log.Info("start getting date from db")

		request, err := currency.GetByID(id)

		if err != nil {
			log.Error("can't get date from db with this id:", sl.Err(err))

			render.JSON(w, r, "can't get date from db with this id")

			return
		}

		log.Info("finish getting date from db")

		// TODO: return struct

		render.JSON(w, r, request)

	}

}
