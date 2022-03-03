package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rpisharody/amfi"
)

func main() {
	var filename = flag.String("j", "~/finances/prices.journal", "The Journal file to read")
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "amfi-fetch: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	commodities := amfi.ReadJournal(file)
	commodity_codes := make([]int, len(commodities))
	var ii int
	for code, _ := range commodities {
		commodity_codes[ii] = int(code)
		ii++
	}

	resp, err := http.Get(amfi.AMFI_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "amfi-fetch: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	amfi_data := amfi.GetAMFIData(resp.Body, commodity_codes)

	for code, name := range commodities {
		fmt.Printf("P %s \"%s\" ₹%.5f ;%d\n", amfi_data[code].Date, name, amfi_data[code].Value, code)
	}
}
