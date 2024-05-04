// package main

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// )

// // processCSV function reads the CSV file and processes each row
// func processCSV(filePath string) (float64, float64, error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(file)
// 	var totalReceived, totalDonated float64
// 	var faultyDonation float64

// 	records, err := reader.ReadAll()
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	for _, record := range records {
// 		amount, err := strconv.ParseFloat(record[0], 64)
// 		if err != nil {
// 			faultyDonation += 0
// 			continue
// 		}
// 		if amount > 0 {
// 			totalReceived += amount
// 			totalDonated += amount
// 			fmt.Printf("Donation of THB %.2f processed successfully.\n", amount)
// 		} else {
// 			faultyDonation += amount
// 		}
// 	}

// 	return totalReceived, totalDonated, nil
// }

// func main() {
// 	if len(os.Args) != 2 {
// 		fmt.Println("Usage: go-tamboon <csv-file>")
// 		os.Exit(1)
// 	}

// 	// Get the current directory
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println("Error getting current directory:", err)
// 		os.Exit(1)
// 	}

// 	// Construct the file path relative to the current directory
// 	filePath := filepath.Join(dir, "data", "fng.1000.csv.rot128")

// 	// Process the CSV file
// 	totalReceived, totalDonated, err := processCSV(filePath)
// 	if err != nil {
// 		fmt.Println("Error processing CSV file:", err)
// 		os.Exit(1)
// 	}

// 	// Summary
// 	fmt.Println("performing donations...")
// 	fmt.Println("done.")
// 	fmt.Printf("\ttotal received: THB  %.2f\n", totalReceived)
// 	fmt.Printf("\tsuccessfully donated: THB  %.2f\n", totalDonated)
// }

/******** working ******************************************************************/
// package main

// import (
// 	"encoding/csv"
// 	"log"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/omise/omise-go"
// 	"github.com/omise/omise-go/operations"
// )

// // Rot128Reader implements io.Reader that transforms
// type Rot128Reader struct{ reader *os.File }

// func NewRot128Reader(r *os.File) (*Rot128Reader, error) {
// 	return &Rot128Reader{reader: r}, nil
// }

// func (r *Rot128Reader) Read(p []byte) (int, error) {
// 	n, err := r.reader.Read(p)
// 	if err != nil {
// 		return n, err
// 	}
// 	rot128(p[:n])
// 	return n, nil
// }

// const (
// 	OmisePublicKey = "pkey_test_5zmos7kagpjnx5wrlmj"
// 	OmiseSecretKey = "skey_test_5zmos7legpeb61bd5sg"
// )

// func main() {
// 	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Open the encrypted CSV file
// 	file, err := os.Open("data/fng.1000.csv.rot128") // Update with your file path
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	// Create a ROT128 reader
// 	rot128Reader, err := NewRot128Reader(file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create a CSV reader using ROT128 reader
// 	reader := csv.NewReader(rot128Reader)

// 	// Read the header row
// 	_, err = reader.Read()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Read remaining rows from the CSV
// 	rows, err := reader.ReadAll()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Iterate over each row
// 	for _, row := range rows {
// 		// Decrypt the card details from the row
// 		cardName := row[0] // Assuming the card name is in the first column
// 		cardNum := row[2]  // Assuming the card number is in the second column
// 		expMonth := row[4] // Assuming the expiration month is in the third column
// 		// expYear := row[5]  // Assuming the expiration year is in the fourth column
// 		// You may need to decrypt other columns as well, depending on your CSV structure

// 		// Convert string to int for expYear
// 		// expYearInt, err := strconv.Atoi(expYear)
// 		// if err != nil {
// 		// 	log.Fatal(err)
// 		// }

// 		// Convert string to time.Month for expMonth
// 		expMonthInt, err := strconv.Atoi(expMonth)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		expMonthTime := time.Month(expMonthInt)
//         // log.Println("ABC", cardName)
// 		// Create a token for the decrypted card details
// 		token, createToken := &omise.Token{}, &operations.CreateToken{
// 			Name:            cardName,
// 			Number:          cardNum,
// 			ExpirationMonth: expMonthTime,
// 			ExpirationYear:  time.Now().AddDate(1, 0, 0).Year(),
// 			// Add other card details as needed
// 		}
// 		if err := client.Do(token, createToken); err != nil {
// 			log.Fatal(err,token)
// 		}

// 		log.Println("Token created:", token, token.ID)

// 		// Retrieve the created token
// 		retrievedToken, retrieveTokenOp := &omise.Token{}, &operations.RetrieveToken{
// 			ID: token.ID,
// 		}
// 		if err := client.Do(retrievedToken, retrieveTokenOp); err != nil {
// 			log.Fatal(err)
// 		}

// 		// log.Println("Retrieved token:", retrievedToken)

// 		// Proceed with creating a charge using the retrieved token
// 		charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
// 			Amount:   100000, // ฿ 1,000.00
// 			Currency: "thb",
// 			Card:     retrievedToken.ID,
// 			// Add other charge details as needed
// 		}
// 		if err := client.Do(charge, createCharge); err != nil {
// 			log.Fatal(err)
// 		}

// 		log.Printf("Charge created: ID %s, Amount %d %s\n", charge.ID, charge.Amount, charge.Currency)
// 	}
// }

// // rot128 implements the ROT128 encryption algorithm
// func rot128(buf []byte) {
// 	for idx := range buf {
// 		buf[idx] += 128
// 	}
// }
/*************************************************************************************/

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
const (
	DefaultCurrency = "thb"
	OmisePublicKey  = "pkey_test_5zmos7kagpjnx5wrlmj"
	OmiseSecretKey  = "skey_test_5zmos7legpeb61bd5sg"
	MinChargeAmount = 2000 // Minimum charge amount in satangs (฿20)
)

// Rot128Reader implements io.Reader that transforms
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

// Function to process a single CSV row
func processCSVRow(row []string, client *omise.Client, wg *sync.WaitGroup, mu *sync.Mutex, totalReceived, successfullyDonated, faultyDonation *int64, topDonors *map[string]int64) {
	defer wg.Done()

	// Decrypt the card details from the row
	cardName := row[0] // Assuming the card name is in the first column
	cardNum := row[2]  // Assuming the card number is in the second column
	expMonth := row[4] // Assuming the expiration month is in the third column

	// Convert string to time.Month for expMonth
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
    log.Println("Token created:", token, token.ID)
    // Retrieve the created token
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
	file, err := os.Open("data/fng.1000.csv.rot128") // Update with your file path
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a ROT128 reader
	rot128Reader, err := NewRot128Reader(file)
	if err != nil {
		log.Fatal(err)
	}

	// Create a CSV reader using ROT128 reader
	reader := csv.NewReader(rot128Reader)

	// Read the header row
	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}

	// Variables to store summary metrics
	var (
		totalReceived     int64
		successfullyDonated int64
		faultyDonation     int64
		topDonors       = make(map[string]int64)
		wg              sync.WaitGroup
		mu              sync.Mutex
	)

	// Process CSV rows concurrently
	for {
		// Read a single row from the CSV
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // Reached end of file
			}
			log.Fatal(err)
		}

		// Increment WaitGroup counter
		wg.Add(1)

		// Process CSV row concurrently
		go processCSVRow(row, client, &wg, &mu, &totalReceived, &successfullyDonated, &faultyDonation, &topDonors)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Print summary
	fmt.Println("performing donations...")
	fmt.Println("done.")
	fmt.Printf("\n\ttotal received: THB %10.2f\n", float64(totalReceived)/100)
	fmt.Printf("\tsuccessfully donated: THB %10.2f\n", float64(successfullyDonated)/100)
	fmt.Printf("\tfaulty donation: THB %10.2f\n", float64(faultyDonation)/100)

	// Compute average per person based on successfully donated amounts
	var averagePerPerson float64
	if successfullyDonated > 0 {
		averagePerPerson = float64(successfullyDonated) / float64(totalReceived)
	}
	fmt.Printf("\n\taverage per person: THB %10.2f\n", averagePerPerson)

	fmt.Println("\ttop donors:")
	for name := range topDonors {
		fmt.Printf(name)
	}
}

// rot128 implements the ROT128 encryption algorithm
func rot128(buf []byte) {
	for idx := range buf {
		buf[idx] += 128
	}
}









// package main

// import (
// 	"log"
// 	"time"

// 	"github.com/omise/omise-go"
// 	"github.com/omise/omise-go/operations"
// )

// const (
// 	OmisePublicKey  = "pkey_test_5zmos7kagpjnx5wrlmj"
// 	OmiseSecretKey  = "skey_test_5zmos7legpeb61bd5sg"
// 	DefaultCardName = "Ms. Primrose F Smallburrow"
// 	DefaultCardNum  = "5426958001804693"
// 	DefaultExpMonth = 6
// 	DefaultCity     = "Bangkok"
// 	DefaultPostal   = "10320"
// 	DefaultSecCode  = "140"
// )

// func main() {
// 	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Calculate expiration year
// 	expYear := time.Now().AddDate(1, 0, 0).Year()

// 	// Create a token
// 	cardToken, createToken := &omise.Token{}, &operations.CreateToken{
// 		Name:            DefaultCardName,
// 		Number:          DefaultCardNum,
// 		ExpirationMonth: DefaultExpMonth,
// 		ExpirationYear:  expYear, // Use the calculated expiration year
// 		City:            DefaultCity,
// 		PostalCode:      DefaultPostal,
// 		SecurityCode:    DefaultSecCode,
// 	}
// 	if err := client.Do(cardToken, createToken); err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Token created:", cardToken.ID)

// 	// Retrieve the created token
// 	retrievedToken, retrieveTokenOp := &omise.Token{}, &operations.RetrieveToken{
// 		ID: cardToken.ID,
// 	}
// 	if err := client.Do(retrievedToken, retrieveTokenOp); err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Retrieved token:", retrievedToken)

// 	// Proceed with creating a charge using the retrieved token (same as your existing code)
// 	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
// 		Amount:   100000, // ฿ 1,000.00
// 		Currency: "thb",
// 		Card:     retrievedToken.ID, // Use the retrieved token ID
// 	}
// 	if err := client.Do(charge, createCharge); err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("Charge created: ID %s, Amount %d %s\n", charge.ID, charge.Amount, charge.Currency)
// }


// package main

// import (
// 	"bufio"
// 	"log"
// 	"os"
// 	"strings"
//     "fmt"
// )

// // Decrypt a string using ROT-128 algorithm
// func rot128Decrypt(input string) string {
// 	var result strings.Builder

// 	for _, char := range input {
// 		// Decrypt each character using ROT-128 algorithm
// 		decryptedChar := char - 128
// 		if decryptedChar < 0 {
// 			decryptedChar += 256
// 		}
// 		result.WriteRune(decryptedChar)
// 	}

// 	return result.String()
// }

// func main() {
// 	// Open the encrypted CSV file
// 	encryptedFile, err := os.Open("data/fng.1000.csv.rot128")
// 	if err != nil {
// 		log.Fatal("Error opening encrypted CSV file:", err)
// 	}
// 	defer encryptedFile.Close()

// 	// Create an output file for the decrypted CSV
// 	decryptedFile, err := os.Create("decrypted.csv")
// 	if err != nil {
// 		log.Fatal("Error creating decrypted CSV file:", err)
// 	}
// 	defer decryptedFile.Close()

// 	// Create a scanner to read the encrypted file line by line
// 	scanner := bufio.NewScanner(encryptedFile)

// 	// Create a writer to write decrypted data to the output file
// 	writer := bufio.NewWriter(decryptedFile)

// 	// Read and decrypt each line of the encrypted file
// 	for scanner.Scan() {
// 		encryptedLine := scanner.Text()
// 		decryptedLine := rot128Decrypt(encryptedLine)

// 		// Write the decrypted line to the output file
// 		_, err := writer.WriteString(decryptedLine + "\n")
// 		if err != nil {
// 			log.Fatal("Error writing decrypted line to output file:", err)
// 		}
// 	}

// 	// Check for scanner errors
// 	if err := scanner.Err(); err != nil {
// 		log.Fatal("Error reading encrypted CSV file:", err)
// 	}

//         // Read and decrypt each line of the encrypted file
//     for scanner.Scan() {
//         encryptedLine := scanner.Text()
//         decryptedLine := rot128Decrypt(encryptedLine)
    
//         fmt.Println("Decrypted line:", decryptedLine) // Debugging output
    
//         // Write the decrypted line to the output file
//         _, err := writer.WriteString(decryptedLine + "\n")
//         if err != nil {
//         log.Fatal("Error writing decrypted line to output file:", err)
//         }
//     }
// 	// Flush the writer to ensure all data is written to the output file
// 	if err := writer.Flush(); err != nil {
// 		log.Fatal("Error flushing writer:", err)
// 	}

// 	log.Println("Decryption completed. Decrypted CSV file saved as decrypted.csv")
// }





