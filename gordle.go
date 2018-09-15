package gordle

import (
	"errors"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font/gofont/goregular"
)

type word struct {
	font    *truetype.Font
	content string
	size    float64
	count   int
	color   colorful.Color
}

var (
	letter     = "[@+\\p{L}\\p{N}]"
	joiner     = "[-.:/''\\p{M}\\x{2032}\\x{00A0}\\x{200C}\\x{200D}~]"
	wordexp, _ = regexp.Compile(letter + "+(" + joiner + "+" + letter + "+)*")

	uselessWords = map[string]bool{
		"and":    true,
		"the":    true,
		"then":   true,
		"in":     true,
		"on":     true,
		"around": true,
		"at":     true,
		"beside": true}
)

// GenerateCloud is the main function for creating a word cloud. Simply hand it a string of words and it will give
// back a PNG image as a byte array
func GenerateCloud(text string) ([]byte, error) {
	dc := gg.NewContext(1000, 1000)
	words, err := extractSortedWords(text)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	generateFonts(words)
	xpos := float64(512)
	ypos := float64(10)
	for i := 0; i < len(words); i++ {
		face := truetype.NewFace(words[i].font, &truetype.Options{Size: words[i].size})
		dc.SetFontFace(face)
		dc.SetRGB(words[i].color.R, words[i].color.G, words[i].color.B)
		dc.DrawStringAnchored(words[i].content, xpos, ypos, 0.5, 0.5)
		ypos += 36
	}
	dc.SavePNG("out.png")
	return nil, nil
}

// extractSortedWords parses the provided string looking for words, counting them as it goes and finally returning
// a slice of pointers to word structs
func extractSortedWords(text string) ([]word, error) {
	// Use a mpa of pointers so that we can count the words as we go. This keeps us from having to iterate over the
	// list of words more than we need to
	counts := make(map[string]*word)
	// Find and count the words in the text
	if plainTextWords := wordexp.FindAllString(text, -1); plainTextWords != nil {
		for _, sWord := range plainTextWords {
			lw := strings.ToLower(sWord)
			if !uselessWords[lw] {
				w, ok := counts[lw]
				if !ok {
					counts[lw] = &word{content: lw, count: 1}
				} else {
					w.count++
				}
			}
		}
		words := make([]word, len(counts))
		i := 0
		for _, v := range counts {
			words[i] = *v
			i++
		}
		sort.Slice(words, func(i, j int) bool {
			return words[i].count > words[j].count
		})
		return words, nil
	}
	// This should only happen if the regex fails
	return nil, errors.New("Failed to find words in provided text")
}

func generateFonts(words []word) {
	// Generate the colors for the fonts from a random pallette
	numColors := 10
	palette, err := colorful.SoftPalette(numColors)
	if err != nil {
		log.Fatal("Failed to generate color palette. All text will be black: " + err.Error())
		numColors = 1
		palette = []colorful.Color{{R: 0, G: 0, B: 0}}
	}
	// Pick a font face
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	// Set the font for each, including size and color
	for i := 0; i < len(words); i++ {
		words[i].font = font
		words[i].color = palette[rand.Intn(numColors)]
		words[i].size = 24
	}
}
