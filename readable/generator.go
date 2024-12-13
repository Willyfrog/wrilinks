package readable

import (
	"strings"
)

var vowels = map[rune]bool{
	'a': true, 'e': true, 'i': true, 'o': true, 'u': true,
}

// GenerateReadableString creates a readable string from input bytes
// ensuring vowel/consonant patterns and lowercase letters only
func GenerateReadableString(input []byte) string {
	// Convert to base26 (a-z)
	var result strings.Builder
	value := 0
	consonantCount := 0

	// Use input bytes to generate sequence
	for _, b := range input {
		value = (value*256 + int(b)) % 26
		char := rune('a' + value)
		
		// If we have 2 consonants already, force a vowel
		if consonantCount >= 2 {
			char = []rune{'a', 'e', 'i', 'o', 'u'}[value%5]
			consonantCount = 0
		} else if !vowels[char] {
			consonantCount++
		} else {
			consonantCount = 0
		}
		
		result.WriteRune(char)
		
		// Keep length reasonable
		if result.Len() >= 12 {
			break
		}
	}

	return result.String()
}
