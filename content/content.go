package content

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

const linkRoot string = "https://raw.githubusercontent.com/first20hours/google-10000-english/master/"
const mediumWords string = "google-10000-english-usa-no-swears-medium.txt"

// Contains all word input elements and the current position
type Content struct {
	WordInputs []*WordInput
	CurrentIndex int
}

// Content element containing word and input the user typed so far
type WordInput struct {
	Word  string
	Typed string
}

func (content *Content) HasNext() bool {
	return content.CurrentIndex + 1 < len(content.WordInputs)
}

func (content *Content) Next() bool{
	content.CurrentIndex++
	return false
}

func (wordInput *WordInput) IsNextChar(char byte) bool {
	nextIndex := len(wordInput.Typed)
	if nextIndex < len(wordInput.Word) {
		return char == wordInput.Word[nextIndex]
	}

	return false
}

func (wordInput *WordInput) IsCorrect() bool {
	return strings.EqualFold(wordInput.Word, wordInput.Typed)
}

func FetchAndCreateContent(count int) *Content {
	// Get 10k most used medium length english words from public github repository
	res, err := http.Get(linkRoot + mediumWords)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	byteArr, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	words := strings.Split(string(byteArr), "\n")
	words = words[:500] // Limit to top 500 words only
	
	// Take random words from the response for the desired length
	selected := make([]string, count)
	for i := 0; i < count; i++ {
		index := rand.Intn(len(words))
		selected[i] = words[index]
	}

	return CreateContent(selected)
}

func CreateContent(words []string) *Content {
	wordInputs := make([]*WordInput, len(words))
	for i, word := range words {
		wordInputs[i] = &WordInput{word, ""}
	}

	return &Content{
		WordInputs: wordInputs,
	}
}
