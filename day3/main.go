package main

import (
	"fmt"
	"os"
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
	line    int
	start   int
}

func (t *Token) LiteralAsInt() int {
	num, _ := strconv.Atoi(t.Literal)

	return num
}

func (t *Token) End() int {
	return t.start + len(t.Literal) - 1
}

type TokenType string

const (
	NUMBER = "NUMBER"
	SYMBOL = "SYMBOL"
	EOF    = "EOF"
)

func (t *Tokenizer) readNumber() string {
	pos := t.pos
	for unicode.IsDigit(rune(t.ch)) {
		t.readChar()
	}

	return t.input[pos:t.pos]
}

func (t *Tokenizer) NextToken() Token {
	var tok Token
	t.skipDot()

	if unicode.IsDigit(rune(t.ch)) {
		tok.Type = NUMBER
		tok.start = t.pos
		tok.Literal = t.readNumber()
		return tok
	}

	if t.ch == 0 {
		tok.Literal = ""
		tok.Type = EOF
		tok.start = t.pos
		return tok
	}

	tok.Type = SYMBOL
	tok.Literal = string(t.ch)
	tok.start = t.pos
	t.readChar()
	return tok
}

func isAdjacent(number *Token, symbol *Token) bool {

	lineDiff := number.line - symbol.line
	if lineDiff < 0 {
		lineDiff *= -1
	}

	if lineDiff == 0 {

		return number.start-symbol.start == 1 || symbol.start-number.End() == 1
	}

	if lineDiff == 1 {

		return number.start-1 <= symbol.start && symbol.start <= number.End()+1
	}

	return false
}

func main() {
	var filename Filepath = "day3/input"

	lines, err := filename.lines()

	panicIfErr(err)

	total := 0
	numbers := make([]Token, 0)
	symbols := make([]Token, 0)
	for lineNumber, l := range lines {
		t := New(l)

		for t.ch != 0 {
			tok := t.NextToken()
			tok.line = lineNumber
			if tok.Type == SYMBOL {
				symbols = append(symbols, tok)
			}

			if tok.Type == NUMBER {
				numbers = append(numbers, tok)
			}
		}

	}

	// Star 1
	for _, num := range numbers {
		for _, sym := range symbols {
			if isAdjacent(&num, &sym) {
				total += num.LiteralAsInt()
			}
		}
	}

	// Star 2
	totalRatio := 0
	for _, sym := range symbols {
		adjacentNumbers := make([]Token, 0)
		for _, num := range numbers {
			if isAdjacent(&num, &sym) {
				adjacentNumbers = append(adjacentNumbers, num)
			}
		}

		if sym.Literal == "*" && len(adjacentNumbers) == 2 {
			totalRatio += adjacentNumbers[0].LiteralAsInt() * adjacentNumbers[1].LiteralAsInt()
		}
	}

	fmt.Println("The result for star 1 is: ", total)
	fmt.Println("The result for star 2 is: ", totalRatio)
}
