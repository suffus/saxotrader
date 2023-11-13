package saxotrader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Currency string
type YesNo string

type Instrument struct {
	Type        string
	Subtype     string
	Symbol      string
	Description string
	AssetType   string
	UIC         string
}

type Derivitive struct {
	Primary    *Instrument
	Underlying *Instrument
}

type Asset interface {
	Primary() *Instrument
	Underlying() *Instrument
}

type Booking struct {
	Date       time.Time
	Asset      Asset
	Amount     float64
	AmountType string
	Currency   Currency
	AccountId  string
}

type Transaction struct {
	Date       time.Time
	AssetTo    Asset
	AssetFrom  Asset
	Amount     float64
	AmountType string
	Currency   Currency
	AccountId  string
}

type Amount struct {
	AmountType string
	Amount     float64
}

type BookingDetail struct {
	Date                            time.Time `json:"date"`
	AccountId                       string    `json:"account_id"`
	AccountCurrency                 Currency  `json:"account_currency"`
	ClientCurrency                  Currency  `json:"client_currency"`
	AmountType                      string    `json:"amount_type"`
	AffectsBalance                  bool      `json:"affects_balance?"`
	AssetType                       string    `json:"asset_type"`
	UIC                             string    `json:"uic"`
	UnderlyingInstrumentSubtype     string    `json:"underlying_instrument_subtype"`
	InstrumentSymbol                string    `json:"instrument_symbol"`
	InstrumentDescription           string    `json:"instrument_description"`
	InstrumentSubtype               string    `json:"instrument_subtype"`
	UnderlyingInstrumentAssetType   string    `json:"underlying_instrument_asset_type"`
	UnderlyingInstrumentDescription string    `json:"underlying_instrument_description"`
	UnderlyingInstrumentSymbol      string    `json:"underlying_instrument_symbol"`
	UnderlyingInstrumentUIC         string    `json:"underlying_instrument_uic"`
	Amount                          float64   `json:"amount"`
	AmountAccountCurrency           float64   `json:"account_currency_amount"`
	AmountClientCurrency            float64   `json:"client_currency_amount"`
	CostType                        string    `json:"cost_type"`
	CostSubtype                     string    `json:"cost_subtype"`
}

type HoldingHistory struct {
	StartDateTime time.Time
	EndDateTime   time.Time
	Asset         Asset
	Amount        float64
	AmountType    string
}

type Portfolio struct {
	Bookings    []*BookingDetail
	Instruments map[string]*Instrument
	Holdings    []*HoldingHistory
}

func (port *Portfolio) GetLatestPositions() {
	for _, instr := range port.Instruments {
		fmt.Println(instr)
	}
}

func NewPortfolio() *Portfolio {
	return &Portfolio{Instruments: make(map[string]*Instrument), Bookings: make([]*BookingDetail, 0)}
}

func (port *Portfolio) LoadBookings(bF *os.File) {
	r := csv.NewReader(bF)
	r.Read()
	for {
		booking := BookingDetail{}
		err := Unmarshal(r, &booking)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Fprintln(os.Stderr, "Error while unmarshalling: ", err.Error())
			}
		} else {
			port.Bookings = append(port.Bookings, &booking)
		}
	}
	sort.Slice(port.Bookings, func(i, j int) bool {
		return port.Bookings[i].Date.Before(port.Bookings[j].Date)
	})
	for _, b := range port.Bookings {
		instr, ok := port.Instruments[b.UIC]
		if !ok {
			newInstr := Instrument{b.AssetType, b.InstrumentSubtype, b.InstrumentSymbol, b.InstrumentDescription, b.AssetType, b.UIC}
			instr = &newInstr
			port.Instruments[b.UIC] = instr
		}

	}

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
		case "saxotrader.Currency":
			f.SetString(record[i])
		case "time.Time":
			tm, err := time.Parse("02-01-2006", record[i])
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
