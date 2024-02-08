package postgres

import (
	"CurrencyClient/iternal/httpserver/handlers/allcurrency"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func New(path string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("%s, failed to open db %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(request []allcurrency.Request) error {
	const op = "storage.postgres.Save"

	// dell all dates in tables
	sqlStatementDelete := ` TRUNCATE roi,currency;
	`
	fmt.Println("Start delete cash")

	err := s.db.QueryRow(sqlStatementDelete)
	if err != nil {

		return fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println("Cash successful deleted ")

	//Insert dates to db

	// query currency
	sqlStatementCurrency := `
		INSERT INTO currency (IDName,Symbol,Name,
		                      Image,CurrentPrice,MarketCap,MarketCapRank,
		                      FullyDilutedValuation,TotalVolume,High24H,Low24H,
		                      PriceChange24H,PriceChangePercentage24H,MarketCapChange24H,
		                      MarketCapChangePercentage24H,CirculatingSupply,TotalSupply,
		                      MaxSupply,Ath,AthChangePercentage,AthDate,Atl,
		                      AtlChangePercentage,AtlDate,LastUpdated)
		VALUES ($1, $2, $3,$4,$5, $6, $7, $8, $9, $10, $11, $12, $13,
		        $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25) 
	`

	//query roi
	sqlStatementRoi := `
		INSERT INTO roi (times,currency,percentage)
		VALUES($1, $2, $3)
`

	for idx := range request {

		_, err := s.db.Exec(sqlStatementCurrency, request[idx].IDName,
			request[idx].IDName, request[idx].Symbol, request[idx].Image,
			request[idx].CurrentPrice, request[idx].MarketCap, request[idx].MarketCapRank,
			request[idx].FullyDilutedValuation, request[idx].TotalVolume, request[idx].High24H,
			request[idx].Low24H, request[idx].PriceChange24H, request[idx].PriceChangePercentage24H,
			request[idx].MarketCapChange24H, request[idx].MarketCapChangePercentage24H, request[idx].CirculatingSupply,
			request[idx].TotalSupply, request[idx].MaxSupply, request[idx].Ath,
			request[idx].AthChangePercentage, request[idx].AthDate, request[idx].Atl,
			request[idx].AtlChangePercentage, request[idx].AtlDate, request[idx].LastUpdated)

		if err != nil {
			if psError, ok := err.(pq.PGError); ok {
				return fmt.Errorf("%s: %w", op, psError.Error())
			}

			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = s.db.Exec(sqlStatementRoi, request[idx].Roi.Times,
			request[idx].Roi.Currency, request[idx].Roi.Percentage)

		if err != nil {
			if psError, ok := err.(pq.PGError); ok {
				return fmt.Errorf("%s: %w", op, psError.Error())
			}

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) GetByID(id int) (*allcurrency.Request, error) {
	const op = "storage.postgres.GetByID"

	// Get date from table full name table
	sqlStatement := `
	SELECT * FROM currency WHERE id = ($1)
`

	var request allcurrency.Request
	var tmpId int

	err := s.db.QueryRow(sqlStatement, id).Scan(&tmpId, &request.IDName,
		request.Symbol, request.Image,
		request.CurrentPrice, request.MarketCap, request.MarketCapRank,
		request.FullyDilutedValuation, request.TotalVolume, request.High24H,
		request.Low24H, request.PriceChange24H, request.PriceChangePercentage24H,
		request.MarketCapChange24H, request.MarketCapChangePercentage24H, request.CirculatingSupply,
		request.TotalSupply, request.MaxSupply, request.Ath,
		request.AthChangePercentage, request.AthDate, request.Atl,
		request.AtlChangePercentage, request.AtlDate, request.LastUpdated)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Get all country with this id

	roiSqlStatement := `
		SELECT * FROM roi WHERE IDCurrency = ($1) 
`

	err = s.db.QueryRow(roiSqlStatement, id).Scan(&tmpId, &request.Roi.Times,
		&request.Roi.Currency, &request.Roi.Percentage)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &request, nil

}
