package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

type MonthTransactions struct {
	total int
	sum   int
}

func main() {

	transactionsByMonth := make(map[string]int)
	totalRevenue := 0.0

	// Open the CSV file
	file, err := os.Open("transacciones.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the records
	for _, record := range records {
		date, err := time.Parse("01/02/2016", record[1])
		if err != nil {
			log.Fatal(err)
		}

		month := date.Month().String()

		if _, exists := transactionsByMonth[month]; !exists {
			transactionsByMonth[month] = new MonthTransactions{0, 0}
		}

		productQuantity, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		productPrice, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Fatal(err)
		}

		totalRevenue += productQuantity * productPrice
		transactionsByMonth[month]++

	}

}
