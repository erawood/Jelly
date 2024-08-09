package jelly

import (
	"bufio"
	"log"
	"os"
	"strings"
)

/*

	Version: 0.0.1
	Date: 09/08/2024 UK

	Jelly is just another way of storing variables in a readable file
	basically another version of a .env file but in a single file and different format

	Jelly files are basic, currently only support variable and addition


	I opted to have any vars or strings on the same line to just append
	No need to have a + operator as its implied

	This is very basic right now and only allows strings or variables
	If a number is needed, it must be done as a string



*/

const (
	EQUALS     = "="
	VARIABLE   = '@'
	SPEECHMARK = '"'
	ESCAPE     = '\\'
)

type Store struct {
	// Auto-increments for each validated expression, no need to manually track or
	// waste any time on determining length
	Count int

	// Each of the key values are stored as Expressions
	// To keep things simple, think of an expression as split by the equals sign
	// On the left of the equals sign is the variable name, on the right is the value.
	// Values can also point to other values. E.g. a = 1; b = @a;
	//
	// Used the @ symbol to act as the pointer character, right now just a basic
	// check to see if the first char (whitespace ignored) is an @
	Items []Expression
}

func NewStore(filename string) Store {
	return fetch(filename)
}

type Expression struct {

	// Basic representation of the name to value relationship
	// Where the Name is the items left of the equals
	// And the Value is the items right of the equals

	// The Rust implementation of this called these left/right but I think a better
	// Naming convention actually benefits here
	Name string

	// As of 06/08/24 (dd/mm/yy) this only works for strings however
	// potentially I will turn this into a Generic
	Value string
}

func fetch(filename string) Store {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var expressions []Expression
	var count int = 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		// Decode Cycle
		expression, ok := splitByChar(text, EQUALS)

		if !ok {
			continue
		}
		expressions = append(expressions, expression)
		count++
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return Store{
		Count: count,
		Items: expressions,
	}
}

func splitByChar(data string, split string) (Expression, bool) {
	var expression Expression

	before, after, found := strings.Cut(data, split)

	if !found {
		return Expression{}, false
	}

	beforeTrimmed := strings.TrimSpace(before)
	afterTrimmed := strings.TrimSpace(after)

	expression = Expression{beforeTrimmed, afterTrimmed}

	return expression, true
}

func (store *Store) CreateItems(value string) []string {
	var items []string

	var current string
	var inString bool = false
	var escape bool = false
	var escapeIndex int
	for i, character := range value {

		// check if this is the start of a new variable
		if character == VARIABLE {
			if inString {
				// then dont treat this as a variable
				current += string(character)
			}

			if current != "" {
				items = append(items, current)
				current = string(character)
			} else {
				current += string(character)

			}

		} else if character == SPEECHMARK {

			if escape {
				current += string(character)
				escape = !escape

			} else {
				if inString {
					inString = false
					current += string(character)
					items = append(items, current)
					current = ""
				} else {
					if current != "" {
						items = append(items, current)
					}

					inString = true
					current = string(character)
				}
			}

		} else if character == ESCAPE {
			escapeIndex = i
			escape = !escape
		} else if character == ' ' {
			if inString {
				current += string(' ')
			} else {
				if current != "" {
					items = append(items, current)
					current = ""
				}

			}
		} else {
			current += string(character)
		}

		if escape {
			if escapeIndex == i-1 {
				escape = false
				escapeIndex = -1
			}
		}
	}

	if current != "" {
		items = append(items, current)
	}

	return items
}

func (store *Store) processItems(items []string) string {
	var processed string

	for _, item := range items {

		if item[0] == VARIABLE {
			processed += store.Get(item[1:])
		} else {
			processed += item[1 : len(item)-1]
		}
	}

	return processed
}

func (store *Store) Get(key string) string {
	keyAsLowercase := strings.ToLower(key)
	var exp Expression
	var found bool = false

	for _, expression := range store.Items {

		if !(strings.ToLower(expression.Name) == keyAsLowercase) {
			continue
		}

		exp = expression
		found = true
		break
	}

	if !found {

		return ""
	}

	items := store.CreateItems(exp.Value)
	processedItems := store.processItems(items)

	return processedItems
}