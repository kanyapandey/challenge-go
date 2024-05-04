package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

// Constants for CSV processing
// See https://dashboard.omise.co/v2/settings for more information.
const (
	DefaultCurrency = "thb"
	OmisePublicKey  = "pkey_test_5zmos7kagpjnx5wrlmj"
	OmiseSecretKey  = "skey_test_5zmos7legpeb61bd5sg"
	MinChargeAmount = 2000
)

// Rot128Reader implements io.Reader that transforms
// Reference: https://github.com/opn-ooo/challenges/blob/master/challenge-go/cipher/rot128.go
type Rot128Reader struct{ reader *os.File }

func NewRot128Reader(r *os.File) (*Rot128Reader, error) {
	return &Rot128Reader{reader: r}, nil
}

func (r *Rot128Reader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if err != nil {
		return n, err
	}
	rot128(p[:n])
	return n, nil
}

func processCSVRow(row []string, client *omise.Client, wg *sync.WaitGroup, mu *sync.Mutex, totalReceived, successfullyDonated, faultyDonation *int64, topDonors *map[string]int64) {
	defer wg.Done()

	cardName := row[0]
	cardNum := row[2]  
	expMonth := row[4] 

	expMonthInt, err := strconv.Atoi(expMonth)
	if err != nil {
		log.Printf("Error converting expiration month: %v\n", err)
		return
	}
	expMonthTime := time.Month(expMonthInt)

	// Create a token for the decrypted card details
	token, createToken := &omise.Token{}, &operations.CreateToken{
		Name:            cardName,
		Number:          cardNum,
		ExpirationMonth: expMonthTime,
		ExpirationYear:  time.Now().AddDate(1, 0, 0).Year(),
	}
	if err := client.Do(token, createToken); err != nil {
		log.Printf("Error creating token: %v\n", token.ID, err)
		mu.Lock()
		defer mu.Unlock()
		*faultyDonation += MinChargeAmount
		return
	}

    // Retrieve a token for the decrypted card details
	retrievedToken, retrieveTokenOp := &omise.Token{}, &operations.RetrieveToken{
		ID: token.ID,
	}
	if err := client.Do(retrievedToken, retrieveTokenOp); err != nil {
		log.Printf("Error retrieving token: %v\n", err)
		return
	}

	// Proceed with creating a charge using the retrieved token
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   MinChargeAmount,
		Currency: DefaultCurrency,
		Card:     retrievedToken.ID,
	}
	if err := client.Do(charge, createCharge); err != nil {
		log.Printf("Error creating charge: %v\n", err)
		mu.Lock()
		defer mu.Unlock()
		*faultyDonation += MinChargeAmount
		return
	}

	mu.Lock()
	defer mu.Unlock()
	*successfullyDonated += MinChargeAmount
	(*topDonors)[cardName] += MinChargeAmount
	*totalReceived += MinChargeAmount
}

func main() {
	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// Open the encrypted CSV file
	file, err := os.Open("data/fng.1000.csv.rot128")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rot128Reader, err := NewRot128Reader(file)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(rot128Reader)

	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}

	// Variables to store summary metrics 
    // wg and mu for Concurrency control Reference: https://medium.com/@nagarjun_nagesh/concurrency-in-go-race-conditions-deadlocks-and-common-pitfalls-52243faf1a2f
	var (
		totalReceived     int64
		successfullyDonated int64
		faultyDonation     int64
		topDonors       = make(map[string]int64)
		wg              sync.WaitGroup
		mu              sync.Mutex
	)

	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}

		wg.Add(1)

		go processCSVRow(row, client, &wg, &mu, &totalReceived, &successfullyDonated, &faultyDonation, &topDonors)
	}

	wg.Wait()

	// Produce Summary
	fmt.Println("performing donations...")
	fmt.Println("done.")
	fmt.Printf("\n\ttotal received: THB %10.2f\n", float64(totalReceived)/100)
	fmt.Printf("\tsuccessfully donated: THB %10.2f\n", float64(successfullyDonated)/100)
	fmt.Printf("\tfaulty donation: THB %10.2f\n", float64(faultyDonation)/100)

	var averagePerPerson float64
	if successfullyDonated > 0 {
		averagePerPerson = float64(successfullyDonated) / float64(totalReceived)
	}
	fmt.Printf("\n\taverage per person: THB %10.2f\n", averagePerPerson)

	fmt.Println("\ttop donors:")
	for name := range topDonors {
		fmt.Printf("\n\t", name)
	}
}

// ROT128 encryption algorithm
func rot128(buf []byte) {
	for idx := range buf {
		buf[idx] += 128
	}
}





