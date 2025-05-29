package game

import (
	"math/rand"
	"strings"
	"testing"

	"strconv"

	"github.com/stretchr/testify/assert"
)

/* Example inputs:
valid: "007123" - not sure if valid, needs 4 digits input?, "1181", " 1022  "
invalid: "$", "-15", " "
*/

// guess tests
func TestValidateGuess_Valid(t *testing.T) {
	guess, err := ValidateGuess("2000")
	assert.NoError(t, err)
	assert.Equal(t, 2000, guess)

	guess, err = ValidateGuess(" 1022  ")
	assert.NoError(t, err)
	assert.Equal(t, 1022, guess)
}

func TestValidateGuess_Invalid(t *testing.T) {
	invalidInputs := []string{"abc", "-15", " ", "999", "10000"}
	for _, input := range invalidInputs {
		_, err := ValidateGuess(input)
		assert.Error(t, err, "expected error for input: %q", input)
	}
}

func TestGenerateSecretCode_Easy(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	code := GenerateSecretCodeWithRand("easy", rng)
	assert.False(t, hasRepeatedDigits(code))
	assert.GreaterOrEqual(t, code, 1000)
	assert.LessOrEqual(t, code, 9999)
}

func TestGenerateSecretCode_Hard(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	code := GenerateSecretCodeWithRand("hard", rng)
	assert.True(t, hasRepeatedDigits(code))
	assert.True(t, isPrime(sumDigits(code)))
}

func TestGenerateSecretCode_Medium(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	code := GenerateSecretCodeWithRand("medium", rng)
	assert.GreaterOrEqual(t, code, 1000)
	assert.LessOrEqual(t, code, 9999)
}

// prefix tests
func TestGenerateTimestampPrefix_Format(t *testing.T) {
	prefix := GenerateTimestampPrefix()

	assert.True(t, strings.HasPrefix(prefix, "TIME: "), "Prefix should start with 'TIME: '")
	assert.Greater(t, len(prefix), len("TIME: "), "Prefix should include a timestamp")
}

func TestGenerateTimestampPrefix_TimestampValue(t *testing.T) {
	prefix := GenerateTimestampPrefix()

	// Extract the number after "TIME: "
	parts := strings.Split(prefix, "TIME: ")
	if len(parts) != 2 {
		t.Fatalf("Expected prefix to contain 'TIME: ' followed by timestamp, got: %s", prefix)
	}

	timestampStr := strings.TrimSpace(parts[1])

	// Parse the timestamp as an int64
	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	assert.NoError(t, err, "Timestamp part should be a valid integer")
	assert.Greater(t, timestampInt, int64(0))
}
