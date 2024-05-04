# GO-TAMBOON ไปทำบุญ

This Go application processes donations from a CSV file using the Omise payment gateway. It decrypts the file using a  ROT-128 algorithm, creates a charge via the Omise Charge API for each row in the decrypted CSV, and produces a summary of the donation process at the end.

## Installation

1. Clone the repository:

    git clone https://github.com/kanyapandey/challenge-go.git

2. Install dependency:

    go mod tidy

3. Replace the values of OmisePublicKey and OmiseSecretKey in main.go with your actual Omise API keys.

4. go run main.go   


