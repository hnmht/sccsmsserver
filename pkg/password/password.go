package password

import (
	"crypto/rand"
	"math/big"
)

const (
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	symbols      = "!@#$%^&*()_+-=[]{};:'\"|`~<>?"
	allChars     = lowerLetters + upperLetters + digits + symbols
)

// Generate a random password with a variable length of 8 to 16 characters
func GenerateSecurePassword() (string, error) {
	// Step 1: Lenght of the generated random password
	length, err := rand.Int(rand.Reader, big.NewInt(9))
	if err != nil {
		return "", err
	}
	passwordLength := int(length.Int64()) + 8

	// Step 2: Ensure the password contains at least on uppercase letter,
	// one lowercase letter, one number, and one special character
	password := make([]byte, passwordLength)
	password[0], err = getRandomChar(lowerLetters)
	if err != nil {
		return "", err
	}
	password[1], err = getRandomChar(upperLetters)
	if err != nil {
		return "", err
	}
	password[2], err = getRandomChar(digits)
	if err != nil {
		return "", err
	}
	password[3], err = getRandomChar(symbols)
	if err != nil {
		return "", err
	}

	// Step 3: Fill the password with the ramaining random characters
	for i := 4; i < passwordLength; i++ {
		password[i], err = getRandomChar(allChars)
		if err != nil {
			return "", err
		}
	}

	// Step 4: Shuffle the password using the Fisher-Yates algorithm to randomize the order
	for i := len(password) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}

// Randomly select a character from the given character set
func getRandomChar(charset string) (byte, error) {
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return 0, err
	}
	return charset[randomIndex.Int64()], nil
}
