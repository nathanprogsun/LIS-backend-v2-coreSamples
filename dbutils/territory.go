package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/countrylist"
	"coresamples/ent/salesterritory"
	"coresamples/ent/zipcode"
)

type LocationType int

const (
	LocationCountry LocationType = iota
	LocationZipcode
	LocationState
)

func GetCountryByName(countryName string, client *ent.Client, ctx context.Context) (*ent.CountryList, error) {
	return client.CountryList.Query().Where(
		countrylist.Or(
			countrylist.CountryName(countryName),
			countrylist.Alpha2Code(countryName),
			countrylist.Alpha3Code(countryName)),
	).First(ctx)
}

func GetAreaByZipcode(code int, client *ent.Client, ctx context.Context) (*ent.Zipcode, error) {
	return client.Zipcode.Query().Where(zipcode.ID(code)).First(ctx)
}

func GetSalesByLocation(location any, locationType LocationType, client *ent.Client, ctx context.Context) (*ent.SalesTerritory, error) {
	query := client.SalesTerritory.Query()
	switch locationType {
	case LocationCountry:
		return query.Where(
			salesterritory.CountryEQ(location.(string)),
		).First(ctx)
	case LocationZipcode:
		return query.Where(
			salesterritory.Zipcode(location.(int)),
		).First(ctx)
	case LocationState:
		return query.Where(
			salesterritory.State(location.(string)),
		).First(ctx)
	default:
		return nil, nil
	}
}
