package game

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// makes sure input is int by converting. afterwards check its a 4 digit number
// return the int if valid else returns an error msg
func ValidateGuess(input string) (int, error) {
	input = strings.TrimSpace(input) //  trim leading/trailing spaces

	guess, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("input is not a valid number")
	}
	if guess < 1000 || guess > 9999 {
		return 0, errors.New("number must be a 4-digit integer (1000–9999)")
	}
	return guess, nil
}

// the true function - calls a determinstic version inside to be able to test it
func GenerateSecretCode(difficulty string) int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return GenerateSecretCodeWithRand(difficulty, rng)
}

// SecretCode is generated according to a certain logic
// Supports difficulty choice: easy, medium, hard
func GenerateSecretCodeWithRand(difficulty string, rng *rand.Rand) int {
	randomNum := rand.Intn(9000) + 1000 // generate a 4-digit number

	switch strings.ToLower(difficulty) {

	case "easy":
		// If easy, just make sure no repeated digits
		for hasRepeatedDigits(randomNum) {
			randomNum = rng.Intn(9000) + 1000
		}
		return randomNum

	case "medium":
		// If medium, use a certain logic
		sum := sumDigits(randomNum)
		if sum%2 == 0 {
			randomNum = reverseInt(randomNum)
		} else {
			randomNum = increaseEachDigitBy1(randomNum)
		}
		if isPalindrome(randomNum) {
			return 7777
		}
		return randomNum

	case "hard":
		// If hard, ensure  repeated digits && prime digit sum
		for {
			if hasRepeatedDigits(randomNum) && isPrime(sumDigits(randomNum)) {
				break
			}
			randomNum = rng.Intn(9000) + 1000
		}
		return randomNum

	default:
		// Fallback logic: return random number
		return randomNum
	}
}

// checks if num has repeated digit - for easy level
func hasRepeatedDigits(n int) bool {
	digitsSeen := make(map[int]bool)

	for n > 0 {
		digit := n % 10
		if digitsSeen[digit] {
			return true // duplicate found
		}
		digitsSeen[digit] = true
		n /= 10
	}

	return false // all digits are unique
}

// int reverse helper func - for medium lvl
func reverseInt(n int) int {
	reversed := 0
	for n > 0 {
		digit := n % 10
		reversed = reversed*10 + digit
		n /= 10
	}
	return reversed
}

// increase each digit by one helper func- for medium lvl
func increaseEachDigitBy1(n int) int {
	result := 0
	multiplier := 1

	for n > 0 {
		digit := n % 10
		newDigit := (digit + 1) % 10 // wraps 9 → 0
		result += newDigit * multiplier
		multiplier *= 10
		n /= 10
	}

	return result
}

// is palindrom helper func - for medium lvl
func isPalindrome(n int) bool {
	original := n
	reversed := 0

	for n > 0 {
		digit := n % 10
		reversed = reversed*10 + digit
		n /= 10
	}

	return original == reversed
}

// sum the digits of the int - for medium and hard lvl
func sumDigits(n int) int {
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// checks if prime - for hard lvl
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// GenerateTimestampPrefix generates a textual prefix containing the current time
func GenerateTimestampPrefix() string {
	currentTime := time.Now()
	timestamp := currentTime.Unix()
	prefix := "TIME: " + fmt.Sprintf("%-v", timestamp)
	return prefix
}
