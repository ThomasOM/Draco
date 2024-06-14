package content

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

const linkRoot string = "https://raw.githubusercontent.com/first20hours/google-10000-english/master/"

const (
	ShortWords WordOptions = iota
	MediumWords
	LongWords
)

// Word options type
type WordOptions int

// Contains all word input elements and the current position
type Content struct {
	Description string
	WordInputs []*WordInput
	CurrentIndex int
}

// Content element containing word and input the user typed so far
type WordInput struct {
	Word  string
	Typed string
}

func ParseWordOption(text string) (WordOptions, error) {
	switch strings.ToUpper(text) {
	case "SHORT":
		return ShortWords, nil
	case "MEDIUM":
		return MediumWords, nil
	case "LONG":
		return LongWords, nil
	}

	return 0, errors.New("Invalid word option")
}

func (wordOptions WordOptions) String() string {
	return []string{"SHORT", "MEDIUM", "LONG"}[wordOptions]
}

func (wordOptions WordOptions) source() string {
	sources := []string{
		linkRoot + "google-10000-english-usa-no-swears-short.txt",
		linkRoot + "google-10000-english-usa-no-swears-medium.txt",
		linkRoot + "google-10000-english-usa-no-swears-long.txt",
	}
	return sources[wordOptions]
}

func (content *Content) HasNext() bool {
	return content.CurrentIndex + 1 < len(content.WordInputs)
}

func (content *Content) Next() bool{
	content.CurrentIndex++
	return false
}

func (content *Content) Reset() {
	for _, w := range content.WordInputs {
		w.Typed = ""
	}

	content.CurrentIndex = 0
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

func FetchAndCreateContent(options WordOptions, count int) *Content {
	// Get 10k most used medium length english words from public github repository
	res, err := http.Get(options.source())
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

	// Create new content with generated description
	description := strings.Title(strings.ToLower(options.String())) + " " + strconv.Itoa(count) + " words"
	return CreateContent(description, selected)
}

func CreateContent(description string, words []string) *Content {
	wordInputs := make([]*WordInput, len(words))
	for i, word := range words {
		wordInputs[i] = &WordInput{word, ""}
	}

	return &Content{
		Description: description,
		WordInputs: wordInputs,
	}
}
