package types

import "time"

// Token represents a valid token inside the chain
type Token struct {
	Name  string      `yaml:"name"`
	Units []TokenUnit `yaml:"units"`
}

// TokenUnit represents a unit of a token
type TokenUnit struct {
	Denom    string   `yaml:"denom"`
	Exponent int      `yaml:"exponent"`
	Aliases  []string `yaml:"aliases,omitempty"`
	PriceID  string   `yaml:"price_id,omitempty"`
}

// TokenPrice represents the price at a given moment in time of a token unit
type TokenPrice struct {
	UnitName  string
	Price     float64
	MarketCap int64
	Timestamp time.Time
}

// NewTokenPrice returns a new TokenPrice instance containing the given data
func NewTokenPrice(unitName string, price float64, marketCap int64, timestamp time.Time) TokenPrice {
	return TokenPrice{
		UnitName:  unitName,
		Price:     price,
		MarketCap: marketCap,
		Timestamp: timestamp,
	}
}
