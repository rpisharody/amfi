package amfi

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const layout = "02-Jan-2006"
const store_layout = "2006/01/02"
const AMFI_URL = "https://www.amfiindia.com/spages/NAVAll.txt"

type AMFI_Data struct {
	Value float64
	Date  string
}

func GetAMFIData(stream io.Reader, amfi_codes []int) map[int64]AMFI_Data {
	data := make(map[int64]AMFI_Data)

	reader := csv.NewReader(stream)
	reader.Comma = ';'

	// Skip the header.
	_, _ = reader.Read()

	// Convert the codes to be read into a map for easier membership check.
	amfi_map := make(map[int]bool)
	for _, amfi_code := range amfi_codes {
		amfi_map[amfi_code] = true
	}

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if errors.Is(err, csv.ErrFieldCount) {
			errMsg(err, false)
			continue
		} else if err != nil {
			errMsg(err, true)
		}
		mf_code, err := strconv.ParseInt(record[0], 10, 32)
		if !amfi_map[int(mf_code)] {
			continue
		}
		if err != nil {
			errMsg(err, true)
		}
		mf_value, err := strconv.ParseFloat(record[len(record)-2], 64)
		if err != nil {
			errMsg(err, false)
			continue
		}
		mf_date, err := time.Parse(layout, record[len(record)-1])
		if err != nil {
			errMsg(err, true)
		}
		data[mf_code] = AMFI_Data{mf_value, mf_date.Format(store_layout)}
	}
	return data
}

func ReadJournal(stream io.Reader) map[int64]string {
	data := make(map[int64]string)

	reader := csv.NewReader(stream)
	reader.Comma = ' '
	reader.FieldsPerRecord = 5

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			errMsg(err, false)
			break
		} else if errors.Is(err, csv.ErrFieldCount) {
			errMsg(err, false)
			continue
		} else if err != nil {
			errMsg(err, true)
		}
		mf_code := record[len(record)-1]
		mf_code = strings.Replace(mf_code, ";", "", 1)
		code, err := strconv.ParseInt(mf_code, 10, 64)
		if err != nil {
			errMsg(err, false)
			continue
		}
		data[code] = record[2]
	}
	return data
}

func errMsg(err error, exit bool) {
	fmt.Fprintf(os.Stderr, "amfi: %v\n", err)
	if exit {
		os.Exit(1)
	}
}
