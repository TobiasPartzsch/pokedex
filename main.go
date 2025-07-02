package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scan := scanner.Scan()
		if !scan {
			log.Fatal("scanner finished")
		}
		text := scanner.Text()
		cleanInput := cleanInput(text)
		fmt.Printf("Your command was: %s\n", cleanInput[0])
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

// func cleanInputManualSinglePass(text string) []string {
// 	var words []string
// 	var currentWord strings.Builder // Efficiently builds the current word
// 	inWord := false                 // State: Are we currently building a word?

// 	for _, r := range text {
// 		// Always lowercase the character
// 		lowerRune := unicode.ToLower(r)

// 		if unicode.IsSpace(lowerRune) {
// 			// If we were in a word, and now hit space, the word is complete
// 			if inWord {
// 				words = append(words, currentWord.String())
// 				currentWord.Reset() // Clear the builder for the next word
// 				inWord = false
// 			}
// 			// If we were already in space, do nothing (skip consecutive spaces)
// 		} else {
// 			// If we hit a non-space character
// 			currentWord.WriteRune(lowerRune)
// 			inWord = true // We are now in a word (or continuing one)
// 		}
// 	}

// 	// After the loop, if we were still in a word (no trailing space), add it
// 	if inWord {
// 		words = append(words, currentWord.String())
// 	}

// 	return words
// }
