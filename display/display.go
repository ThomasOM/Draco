package display

import (
	"fmt"
	"me/thomazz/draco/content"
	"me/thomazz/draco/stats"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const titleColor tcell.Color = tcell.ColorPurple
const correctColor tcell.Color = tcell.ColorGreen
const wrongColor tcell.Color = tcell.ColorRed
const noColor tcell.Color = tcell.ColorGray
const cursorColor tcell.Color = tcell.ColorWhite

type Display struct {
	content *content.Content
	stats   *stats.Stats

	app      *tview.Application
	root     *tview.Flex
	textView *tview.TextView
	statsView   *tview.TextView
	buttonView   *tview.TextView

	finished bool
	closed   bool
}

func CreateDisplay(content *content.Content, stats *stats.Stats) *Display {
	return &Display{
		content: content,
		stats:   stats,
	}
}

func (display *Display) Start() {
	// Set up main text view
	display.textView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(noColor)

	display.textView.SetBorder(true).
		SetTitle(" Draco ").
		SetTitleColor(titleColor).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			display.processInput(event)
			display.textView.SetText(display.renderText())
			return event
		})

	display.textView.SetText(display.renderText())

	// Stats view
	display.statsView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(noColor)

	// Button view
	display.buttonView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(noColor).
		SetText(
			"ESC - Close app\n" +
			"F1 - Restart",
		)

	display.statsView.SetBorder(true).SetTitle(" Stats ")
	display.buttonView.SetBorder(true).SetTitle(" Actions ")

	// Footer to display stats
	footer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(display.statsView, 0, 1, false).
		AddItem(display.buttonView, 0, 1, false)
		
	// Root flex element
	display.root = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(display.textView, 0, 1, false).
		AddItem(footer, 25, 0, false)

	// Set up app
	display.app = tview.NewApplication()
	display.app.SetRoot(display.root, true).SetFocus(display.textView)

	// Start update loop
	go display.update()

	// Locks thread until app is closed
	if err := display.app.Run(); err != nil {
		panic(err)
	}

	// Shuts down update loop instantly after app closes
	display.closed = true
}

func (display *Display) finish() {
	// Marks app as done
	display.finished = true
	display.stats.TimeFinished = time.Now()
}

func (display *Display) processInput(event *tcell.EventKey) {
	// First process action keys
	if display.processActions(event) {
		return
	}

	// Do not allow extra input after finishing
	if display.finished {
		return
	}

	currentInput := display.content.WordInputs[display.content.CurrentIndex]

	if event.Key() == tcell.KeyBackspace {
		typed := currentInput.Typed

		if len(typed) == 0 {
			previousIndex := max(0, display.content.CurrentIndex-1)
			previousInput := display.content.WordInputs[previousIndex]

			// Only allow user to go back to previous word input if it was incorrect
			if !previousInput.IsCorrect() {
				display.content.CurrentIndex = previousIndex
			}
		} else {
			typed = typed[0 : len(typed)-1]
		}

		currentInput.Typed = typed
	} else if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case ' ':
			display.stats.TotalInputsTyped++
			display.stats.TotalSpacesTyped++
			if currentInput.IsCorrect() {
				display.stats.CorrectSpacesTyped++
			}

			// Move to the next word if possible or finish if on the last word
			if display.content.HasNext() {
				display.content.Next()
			} else {
				display.finish()
			}
		default:
			display.stats.TotalCharactersTyped++
			display.stats.TotalInputsTyped++

			// Track correct chars typed
			if currentInput.IsNextChar(byte(event.Rune())) {
				display.stats.CorrectCharactersTyped++
			}

			typed := currentInput.Typed
			currentInput.Typed = typed + string(event.Rune())

			// Only start timer when typing
			if display.stats.TimeStarted.IsZero() {
				display.stats.TimeStarted = time.Now()
			}

			// If no word is next we can stop the app
			if currentInput.IsCorrect() && !display.content.HasNext() {
				display.finish()
			}
		}
	}
}

func (display *Display) processActions(event *tcell.EventKey) bool {
	switch event.Key() {
	case tcell.KeyEscape: // Close the app when pressing escape
		display.app.Stop()
		return true
	case tcell.KeyF1: // Reset app on F1
		display.finished = false
		display.content.Reset()
		display.stats.Reset()
		return true
	}

	return false
}

func (display *Display) update() {
	for !display.closed {
		var builder strings.Builder

		// Dump all stats into string builder
		builder.WriteString("Mode: ")
		builder.WriteString(display.content.Description)
		builder.WriteString("\n")
		builder.WriteString("Total inputs: ")
		builder.WriteString(strconv.Itoa(display.stats.TotalInputsTyped))
		builder.WriteString("\n")
		builder.WriteString("Time passed: ")
		builder.WriteString(fmt.Sprintf("%.1fs", display.stats.SecondsPassed()))
		builder.WriteString("\n")
		builder.WriteString("Accuracy: ")
		builder.WriteString(fmt.Sprintf("%.2f%%", display.stats.Accuracy() * 100))
		builder.WriteString("\n")
		builder.WriteString("Raw WPM: ")
		builder.WriteString(fmt.Sprintf("%.1f", display.stats.WordsPerMinuteRaw()))
		builder.WriteString("\n")
		builder.WriteString("WPM: ")
		builder.WriteString(fmt.Sprintf("%.1f", display.stats.WordsPerMinute()))

		display.statsView.SetText(builder.String())
		display.app.Draw()

		time.Sleep(100) // 10 updates per second
	}
}

func (display *Display) renderText() string {
	var builder strings.Builder

	// Keep track of the last color so we do not insert unnecessary dynamic colors
	lastColor := noColor
	for index, wordInput := range display.content.WordInputs {
		word := []rune(wordInput.Word)
		typed := []rune(wordInput.Typed)
		maxLength := max(len(word), len(typed))

		for i := 0; i < maxLength; i++ {
			color := noColor

			var wordChar rune
			if i < len(word) {
				wordChar = word[i]
			} else {
				wordChar = -1
			}

			var typedChar rune
			if i < len(typed) {
				typedChar = typed[i]
			} else {
				typedChar = -1
			}

			// Select the correct character and color to display
			var char rune
			if wordChar != -1 && typedChar != -1 {
				char = wordChar
				if wordChar == typedChar {
					color = correctColor
				} else {
					color = wrongColor
				}
			} else if wordChar != -1 {
				char = wordChar
				if index < display.content.CurrentIndex {
					color = wrongColor
				} else {
					color = noColor
				}
			} else if typedChar != -1 {
				char = typedChar
				color = wrongColor
			} else {
				char = '?'
				color = noColor
			}

			// Cursor highlight
			if index == display.content.CurrentIndex && i == len(typed) {
				color = cursorColor
			}

			// Update dynamic coloring if changed
			if lastColor != color {
				builder.WriteRune('[')
				builder.WriteString(color.Name())
				builder.WriteRune(']')
				lastColor = color
			}

			builder.WriteString(string(char))
		}

		// Separator character
		if index == display.content.CurrentIndex && len(word) == len(typed) {
			lastColor = cursorColor
			builder.WriteRune('[')
			builder.WriteString(lastColor.Name())
			builder.WriteRune(']')
			builder.WriteRune('_')
		} else {
			builder.WriteRune(' ')
		}
	}

	return builder.String()
}
