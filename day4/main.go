package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type Filepath string

func (f *Filepath) lines() ([]string, error) {
	bytes, err := os.ReadFile(string(*f))

	if err != nil {
		return nil, err
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")
	return lines, nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Tokenizer struct {
	input        string
	pos          int
	readPosition int
	ch           byte
}

func (t *Tokenizer) skipDot() {
	for t.ch == '.' {
		t.readChar()
	}
}

func (t *Tokenizer) skipWhitespace() {
	for t.ch == ' ' {
		t.readChar()
	}
}

func New(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.readChar()
	return t
}

func (t *Tokenizer) readChar() {
	if t.readPosition >= len(t.input) {
		t.ch = 0
	} else {
		t.ch = t.input[t.readPosition]
	}

	t.pos = t.readPosition
	t.readPosition += 1
}

type Token struct {
	Type    TokenType
	Literal string
	lne     int
	start   int
}

func (t *Token) End() int {
	return t.start + len(t.Literal) - 1
}

type TokenType string

func (t *Tokenizer) readNumber() string {
	pos := t.pos
	for unicode.IsDigit(rune(t.ch)) {
		t.readChar()
	}

	return t.input[pos:t.pos]
}

type Card struct {
	number int
}

func (t *Tokenizer) readCard() Card {
	t.skipWhitespace()
	for t.ch != ' ' {
		t.readChar()
	}

	t.skipWhitespace()

	num := t.readNumber()

	for t.ch != ':' {
		t.readChar()
	}
	t.readChar()
	t.skipWhitespace()

	converted, err := strconv.Atoi(num)
	panicIfErr(err)
	return Card{number: converted}

}

func (t *Tokenizer) readWinningNumbers() []int {
	t.skipWhitespace()

	pos := t.pos

	for t.ch != '|' {
		t.readChar()
	}

	numberString := strings.TrimSpace(t.input[pos:t.pos])

	numberStrings := strings.Split(numberString, " ")

	numbers := make([]int, 0)
	for _, n := range numberStrings {
		if strings.TrimSpace(n) == "" {
			continue
		}
		num, err := strconv.Atoi(strings.TrimSpace(n))
		panicIfErr(err)

		numbers = append(numbers, num)
	}

	t.readChar()
	return numbers

}

func (t *Tokenizer) readScratchNumbers() []int {

	t.skipWhitespace()

	pos := t.pos

	for t.ch != 0 {
		t.readChar()
	}

	numberString := strings.TrimSpace(t.input[pos:t.pos])

	numberStrings := strings.Split(numberString, " ")

	numbers := make([]int, 0)
	for _, n := range numberStrings {
		if strings.TrimSpace(n) == "" {
			continue
		}
		num, err := strconv.Atoi(strings.TrimSpace(n))
		panicIfErr(err)

		numbers = append(numbers, num)
	}

	return numbers

}

func star1(lines []string) {

	total_points := 0
	for _, l := range lines {
		t := New(l)
		if t.ch == 0 {
			continue
		}

		t.readCard()

		winningNumbers := t.readWinningNumbers()

		for _, w := range winningNumbers {
			fmt.Printf("%d ", w)
		}
		fmt.Println()

		scratchNumbers := t.readScratchNumbers()

		winners := make([]int, 0)

		for _, sn := range scratchNumbers {
			if slices.Contains(winningNumbers, sn) {
				winners = append(winners, sn)
			}
		}

		points := 0

		if len(winners) > 0 {

			points = 1

			i := 0
			for i < len(winners)-1 {
				points *= 2
				i += 1
			}

		}
		total_points += points
	}

	fmt.Printf("Total points for the elves is %d\n", total_points)

}

func star2(lines []string) {

	copies := make(map[int]int, 0)

	noWhitespaceLines := make([]string, 0)

	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			noWhitespaceLines = append(noWhitespaceLines, l)
		}
	}

	for lineNumber, l := range noWhitespaceLines {
		if l == "" {
			continue
		}
		t := New(l)
		if t.ch == 0 {
			continue
		}

		t.readCard()

		winningNumbers := t.readWinningNumbers()

		scratchNumbers := t.readScratchNumbers()

		winners := make([]int, 0)

		for _, sn := range scratchNumbers {
			if slices.Contains(winningNumbers, sn) {
				winners = append(winners, sn)
			}
		}

		// Create copies

		if len(winners) > 0 {

			if _, ok := copies[lineNumber]; !ok {
				copies[lineNumber] = 1
			} else {

				copies[lineNumber] += 1
			}
		}

		i := 1
		for i <= len(winners) && i+lineNumber < len(lines) {
			_, ok := copies[lineNumber+i]
			if ok {
				copies[lineNumber+i] += copies[lineNumber]
			} else {

				copies[lineNumber+i] = 1
			}
			i++

		}

	}

	sum := 0
	for _, v := range copies {
		sum += v
	}
	fmt.Printf("Total points for the elves is %d\n", sum)
}

func main() {
	var filename Filepath = "day4/input_test"

	lines, err := filename.lines()

	panicIfErr(err)

	star1(lines)
	star2(lines)
}
