package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"fmt"
	"strconv"
	"strings"
)

const (
	BadSales = "BAD SALES"
)

type ISalesService interface {
	GetSalesByTerritory(zipcode string, countryName string, state string, ctx context.Context) (string, error)
}

type SalesService struct {
	Service
}

func newSalesService(dbClient *ent.Client, redisClient *common.RedisClient) ISalesService {
	s := &SalesService{
		Service: InitService(dbClient, redisClient),
	}
	return s
}

func (s *SalesService) GetSalesByTerritory(zipcode string, countryName string, state string, ctx context.Context) (string, error) {
	if zipcode == "" || countryName == "" {
		return BadSales, fmt.Errorf("missing zipcode and country name")
	}
	country, err := dbutils.GetCountryByName(countryName, s.dbClient, ctx)
	if err != nil || country == nil {
		return BadSales, err
	}
	// first look up by country name, if it's in the US then we wouldn't find any and so that we can then look up by zipcode
	countryNameStandard := country.CountryName
	sales, err := dbutils.GetSalesByLocation(countryNameStandard, dbutils.LocationCountry, s.dbClient, ctx)
	if sales != nil {
		return sales.Sales, nil
	}
	// Can't find sales by country, try finding them by zipcode instead,
	// we can only look up by zipcode when in US
	if countryNameStandard != "United States" {
		return BadSales, err
	}
	code, err := strconv.Atoi(strings.TrimSpace(zipcode)[:5])
	if err != nil {
		code = 0
	}
	sales, err = dbutils.GetSalesByLocation(code, dbutils.LocationZipcode, s.dbClient, ctx)
	if sales != nil {
		return sales.Sales, nil
	}
	// Can't find sales by zipcode, try finding them by state instead
	zipcodeArea, err := dbutils.GetAreaByZipcode(code, s.dbClient, ctx)
	if err != nil {
		return BadSales, fmt.Errorf("invalid zipcode")
	}
	sales, err = dbutils.GetSalesByLocation(zipcodeArea.State, dbutils.LocationState, s.dbClient, ctx)
	if sales != nil {
		return sales.Sales, nil
	}

	return BadSales, err
}
