package saxotrader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Currency string
type YesNo string

type Instrument struct {
	Subtype     string
	Symbol      string
	Description string
	AssetType   string
	UIC         string
}

type Booking struct {
	Date                            time.Time `json:"date"`
	AccountId                       string    `json:"account_id"`
	AccountCurrency                 Currency  `json:"account_currency"`
	ClientCurrency                  Currency  `json:"client_currency"`
	AmountType                      string    `json:"amount_type"`
	AffectsBalance                  bool      `json:"affects_balance?"`
	AssetType                       string    `json:"asset_type"`
	UIC                             string    `json:"uic"`
	UnderlyingInstrumentSubtype     string    `json:"underlying_instrument_subtype"`
	InstrumentSymbol                string
	InstrumentDescription           string
	InstrumentSubtype               string
	UnderlyingInstrumentAssetType   string
	UnderlyingInstrumentDescription string
	UnderlyingInstrumentSymbol      string
	UnderlyingInstrumentUIC         string
	Amount                          float64
	AmountAccountCurrency           float64
	AmountClientCurrency            float64
	CostType                        string
	CostSubtype                     string
}

type FieldMismatch struct {
	Expected int
	Got      int
}

type UnsupportedType string

func (ut *UnsupportedType) Error() string {
	return fmt.Sprintf("Unsupported type in Unmarshal: %s", *ut)
}

func (f *FieldMismatch) Error() string {
	return fmt.Sprintf("Mismatch in expected field length, expected %d got %d", f.Expected, f.Got)
}

func ParseSaxoBool(val string) (bool, error) {
	val = strings.ToLower(val)
	if val == "yes" || val == "true" {
		return true, nil
	} else if val == "no" || val == "false" {
		return false, nil
	}

	return false, errors.New(fmt.Sprintf("Could not interpret %s as a bool", val))
}

func Unmarshal(reader *csv.Reader, v interface{}) error {
	record, err := reader.Read()
	if err != nil {
		return err
	}
	s := reflect.ValueOf(v).Elem()
	if s.NumField() != len(record) {
		return &FieldMismatch{s.NumField(), len(record)}
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		switch f.Type().String() {
		case "string":
			f.SetString(record[i])
		case "int":
			ival, err := strconv.ParseInt(record[i], 10, 0)
			if err != nil {
				return err
			}
			f.SetInt(ival)
		case "float64":
			fval, err := strconv.ParseFloat(record[i], 64)
			if err != nil {
				return err
			}
			f.SetFloat(fval)
		case "bool":
			bval, err := ParseSaxoBool(record[i])
			if err != nil {
				return err
			}
			f.SetBool(bval)
		case "Currency":
			f.SetString(record[i])
		case "time.Time":
			tm, err := time.Parse("22-10-2018", record[i])
			if err != nil {
				return err
			}
			f.Set(reflect.ValueOf(tm))

		default:
			errx := UnsupportedType(f.Type().String())
			return &errx
		}
	}
	return nil
}
