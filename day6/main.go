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

func (t *Tokenizer) readNumber() string {
	pos := t.pos
	for unicode.IsDigit(rune(t.ch)) {
		t.readChar()
	}

	return t.input[pos:t.pos]
}

func readValues(valuesLine string) []int {

	values := make([]int, 0)
	t := New(valuesLine)
	for t.ch != ':' {
		t.readChar()
	}
	t.readChar()
	for t.ch != 0 {
		for t.ch == ' ' {
			t.readChar()
		}
		pos := t.pos
		t.readNumber()

		valuesStr := valuesLine[pos:t.pos]
		value, err := strconv.Atoi(valuesStr)
		panicIfErr(err)
		values = append(values, value)
	}

	return values
}

func readValuesBadKerning(valuesLine string) int {

	t := New(valuesLine)
	for t.ch != ':' {
		t.readChar()
	}
	t.readChar()
	pos := t.pos
	for t.ch != 0 {
		t.readChar()
	}
	removedWhitespace := strings.ReplaceAll(t.input[pos:t.pos], " ", "")
	num, err := strconv.Atoi(removedWhitespace)
	panicIfErr(err)

	return num
}

type Race struct {
	time     int
	distance int
}

func (r *Race) computeWaysToWin() int {
	waysToWin := 0
	for t := 1; t < r.time; t++ {
		speed := t
		distance := speed * (r.time - t)

		if distance > r.distance {
			waysToWin++
		}
	}
	return waysToWin
}

func CreateRaces(times []int, distances []int) []Race {
	races := make([]Race, 0)

	for i := 0; i < len(times); i++ {
		race := Race{time: times[i], distance: distances[i]}

		races = append(races, race)
	}

	return races
}

func main() {
	var filename Filepath = "day6/input"

	lines, err := filename.lines()

	panicIfErr(err)

	times := readValues(lines[0])
	distances := readValues(lines[1])

	races := CreateRaces(times, distances)

	total_win_ratio := 1
	for _, r := range races {
		total_win_ratio *= r.computeWaysToWin()
	}

	fmt.Printf("Ratio: %d\n", total_win_ratio)

	r := Race{time: readValuesBadKerning(lines[0]), distance: readValuesBadKerning(lines[1])}
	fmt.Printf("Actual win ways %d\n", r.computeWaysToWin())

}
