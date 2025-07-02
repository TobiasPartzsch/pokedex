package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var commands map[string]cliCommand

func main() {
	commands = map[string]cliCommand{
		"help": {
			description: "Displays a help message",
			callback: func() error {
				return printHelp(commands)
			},
		},
		"exit": {
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scan := scanner.Scan()
		if !scan {
			log.Fatal("scanner finished")
		}
		text := scanner.Text()
		cleanInput := cleanInput(text)
		if len(cleanInput) == 0 {
			continue
		}
		command, exists := commands[cleanInput[0]]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}
		command.callback()
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	msg := "Closing the Pokedex... Goodbye!"
	fmt.Println(msg)
	os.Exit(0)
	return nil
}

func printHelp(commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println("")

	for command, definition := range commands {
		fmt.Println(command + ": " + definition.description)
	}
	return nil
}

type cliCommand struct {
	description string
	callback    func() error
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
