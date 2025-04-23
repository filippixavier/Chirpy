package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuth(t *testing.T) {
	// cases := []struct {
	// 	input    string
	// 	expected []string
	// }{
	// 	{
	// 		input:    "  hello  world  ",
	// 		expected: []string{"hello", "world"},
	// 	},
	// 	{
	// 		input:    " this is     sparta !",
	// 		expected: []string{"this", "is", "sparta", "!"},
	// 	},
	// 	{
	// 		input:    "  hello  ",
	// 		expected: []string{"hello"},
	// 	},
	// 	{
	// 		input:    "  hello  world  ",
	// 		expected: []string{"hello", "world"},
	// 	},
	// 	{
	// 		input:    "  HellO  World  ",
	// 		expected: []string{"hello", "world"},
	// 	},
	// }

	// for _, c := range cases {
	// 	actual := cleanInput(c.input)
	// 	// Check the length of the actual slice against the expected slice
	// 	// if they don't match, use t.Errorf to print an error message
	// 	// and fail the test
	// 	if len(actual) != len(c.expected) {
	// 		t.Errorf("Expected %v, got %v", c.expected, actual)
	// 	}
	// 	for i := range actual {
	// 		word := actual[i]
	// 		expectedWord := c.expected[i]
	// 		// Check each word in the slice
	// 		// if they don't match, use t.Errorf to print an error message
	// 		// and fail the test
	// 		if word != expectedWord {
	// 			t.Errorf("%v is not the same as expected %v", word, expectedWord)
	// 		}
	// 	}
	// }

	t.Run("Test case", func(t *testing.T) {
		id := uuid.New()
		token, err := MakeJWT(id, "secretoken", 30*time.Minute)

		if err != nil {
			t.Error(err)
			return
		}

		dcoded, err := ValidateJWT(token, "secretoken")

		if err != nil {
			t.Error(err)
			return
		}

		if dcoded != id {
			t.Errorf("original id and decoded one does not match")
		}

	})

}
