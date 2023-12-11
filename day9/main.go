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
	for unicode.IsDigit(rune(t.ch)) || t.ch == '-' {
		t.readChar()
	}

	return t.input[pos:t.pos]
}

func readSpaceSepNumbers(line string) []int {
	t := New(line)
	numbers := make([]int, 0)

	for t.ch != 0 {
		for t.ch == ' ' {
			t.readChar()
		}

		numStr := t.readNumber()
		num, err := strconv.Atoi(numStr)
		panicIfErr(err)
		numbers = append(numbers, num)

	}
	return numbers
}

func computeDiffArray(numbers []int) []int {
	diffAry := make([]int, len(numbers)-1)

	for i := 0; i < len(numbers)-1; i++ {
		diff := numbers[i+1] - numbers[i]
		diffAry[i] = diff
	}

	return diffAry
}

func onlyZeros(numbers []int) bool {
	for _, n := range numbers {
		if n != 0 {
			return false
		}
	}

	return true
}

func extrapolateNextValue(lastValue int, diffArrays [][]int) int {
	nextValue := lastValue
	for i := len(diffArrays) - 1; i >= 0; i-- {
		ary := diffArrays[i]
		nextValue += ary[len(ary)-1]
	}
	return nextValue
}

func extrapolatePreviousValue(firstValue int, diffArrays [][]int) int {
	lastDiff := 0
	for i := len(diffArrays) - 1; i >= 0; i-- {
		ary := diffArrays[i]
		diff := ary[0] - lastDiff
		fmt.Printf("%d-%d-%d\n", ary[0], lastDiff, diff)
		lastDiff = diff
	}
	return firstValue - lastDiff
}

func main() {
	var filename Filepath = "day9/input"

	lines, err := filename.lines()
	panicIfErr(err)
	total := 0
	total_prev := 0

	for _, l := range lines {
		if l == "" {
			continue
		}

		numbers := readSpaceSepNumbers(l)
		diffs := make([][]int, 0)
		var currDiff *[]int = &numbers

		for !onlyZeros(*currDiff) {
			diff := computeDiffArray(*currDiff)
			currDiff = &diff
			diffs = append(diffs, diff)
		}
		next := extrapolateNextValue(numbers[len(numbers)-1], diffs)
		previous := extrapolatePreviousValue(numbers[0], diffs)
		total += next
		total_prev += previous

		fmt.Printf("%v->%v->%d,%d\n", numbers, diffs, next, previous)

	}

	fmt.Printf("Total is %d\n", total)
	fmt.Printf("Total previous is %d\n", total_prev)

}
