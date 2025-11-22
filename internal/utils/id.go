package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func GetId() string {
	// Capturing the first part of the HEX id, today's date. to the second
	today := time.Now()

	// Formatting the date to YYYY-MM-DD HH:MM:SS format.
	formattedDate := today.Format("2006-01-02 15:04:05")

	formattedDateBytes := []byte(formattedDate)
	// Generating random bytes. 32 bytes to be exact
	magic := make([]byte, 32)

	// creating a hash maker, it what makes hashes.
	hasher := sha256.New()

	// Giving the hasher data that it will hash for us.
	hasher.Write(append(magic, formattedDateBytes...))

	// Here we've calculated the sum, i.e the hash.
	sum := hasher.Sum(nil)

	// Our hash in string data type.
	hashedData := hex.EncodeToString(sum)

	return hashedData
}
