package gordle

import (
	"fmt"
	"testing"
)

func TestExtractSortedWords(t *testing.T) {
	s := "The quick brown fox jumps over the lazy dog dog dog"
	words, err := extractSortedWords(s)
	if err != nil {
		t.Error(err)
	}
	for _, w := range words {
		fmt.Println(w.content)
		if w.content == "dog" && w.count != 3 {
			t.Fatal("dog does not have 3 counts")
		}
	}
}
