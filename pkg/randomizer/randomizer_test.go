package randomizer

import "testing"

func TestRandomReturnsCharsWithTheLengthOfSeven(t *testing.T) {
	random := New()

	expectedLength := 7
	actual := len(random.Random(7, Chars))

	if expectedLength != actual {
		t.Errorf("Expected %d, got %d", expectedLength, actual)
	}
}
