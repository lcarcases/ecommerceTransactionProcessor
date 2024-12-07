package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type MonthTransactions struct {
	total int
	sum   float64
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	transactionsByMonth := make(map[string]MonthTransactions)
	totalRevenue := 0.0
	report := ""

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresHost := os.Getenv("POSTGRES_HOST")

	// Set up the PostgreSQL connection URL
	// Username is "myuser" and password is "mypassword" and was set up in "Dockerfile"
	// Database name is mydatabase and was set up in Dockerfile
	//connStr := "postgres://postgresUser:postgresPassword@postgresHost:5433/postgresDB"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", postgresUser, postgresPassword, postgresHost, postgresDB)

	// Establish the connection
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	defer conn.Close(context.Background())

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
	for i, record := range records {
		// Avoid read header of the CSV
		if i == 0 {
			continue
		}

		/*transactionId, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}*/

		date, err := time.Parse("01/02/06", record[0])
		if err != nil {
			log.Fatal(err)
		}

		month := date.Month().String()

		productId, err := strconv.Atoi(record[1])

		if _, exists := transactionsByMonth[month]; !exists {
			transactionsByMonth[month] = MonthTransactions{0, 0}
		}

		productQuantity, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		productPrice, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		totalRevenue += productQuantity * productPrice

		// Update the total transactions and revenue for the month
		currentTransactionsMonth := transactionsByMonth[month]
		currentTransactionsMonth.sum += productQuantity * productPrice
		currentTransactionsMonth.total++
		transactionsByMonth[month] = currentTransactionsMonth

		// Generate the report
		report = fmt.Sprintf("Total Revenue: $%.2f\n", totalRevenue)
		for month, transactions := range transactionsByMonth {
			avgTransactionValue := transactions.sum / float64(transactions.total)
			report += fmt.Sprintf("Number of transactions in %s: %d\n", month, transactions.total)
			report += fmt.Sprintf("Average transaction value in %s: $%.2f\n", month, avgTransactionValue)
		}

		// Insert data into the database
		_, err = conn.Exec(context.Background(),
			"INSERT INTO transactions (date, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
			date, productId, productQuantity, productPrice)
		if err != nil {
			log.Fatal("Unable to insert data:", err)
		}

	}

	//Mail sending

	// Set up authentication iformation
	fmt.Println("Sending email...")
	fmt.Println(os.Getenv("gmailPassword"))
	auth := smtp.PlainAuth("", "lcarcases@gmail.com", os.Getenv("gmailPassword"), "smtp.gmail.com")

	// Define the message to be sent and the recipient
	to := []string{"lcarcases@gmail.com"}
	msg := []byte("To: lcarcases@gmail.com\r\n" +
		"Subject: Monthly Report\r\n" +
		"\r\n" +
		report + "\r\n")

	// Send the email
	err = smtp.SendMail("smtp.gmail.com:587", auth, "lcarcases@gmail.com", to, msg)

	if err != nil {
		log.Fatal(err)
	}

}
