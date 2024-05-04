package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

func TestProcessCSVRow(t *testing.T) {
	tests := []struct {
		name            string
		row             []string
		expectedTokenID string
		expectedError   error
	}{
		{
			name:            "Valid row",
			row:             []string{"Obi-wan Kenobi", "4242424242424242", "123", "456", "12"},
			expectedTokenID: "tok_123456",
			expectedError:   nil,
		},
		{
			name:            "Invalid expiration month",
			row:             []string{"Luke Skywalker", "4242424242424242", "123", "456", "invalid"},
			expectedTokenID: "",
			expectedError:   errors.New("Error converting expiration month: strconv.Atoi: parsing \"invalid\": invalid syntax"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &testClient{}

			var buf bytes.Buffer
			logger := log.New(&buf, "", 0)

			tokenID, err := processCSVRow(tt.row, client, logger)

			if tokenID != tt.expectedTokenID {
				t.Errorf("processCSVRow() tokenID = %s, want %s", tokenID, tt.expectedTokenID)
			}

			if (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError == nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("processCSVRow() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

type testClient struct{}

func (c *testClient) Do(resource interface{}, operation *operations.Op) error {
	if _, ok := resource.(*omise.Token); ok {
		return nil
	}
	return nil
}
