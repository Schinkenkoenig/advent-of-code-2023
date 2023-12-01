package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Filepath string

func (f *Filepath) lines() ([]string, error) {
	bytes, err := os.ReadFile("input")

	if err != nil {
		return nil, err
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")
	return lines, nil
}

func numberFromTwoRunes(digit1 rune, digit2 rune) (int, error) {
	if !unicode.IsDigit(digit1) || !unicode.IsDigit(digit2) {
		return 0, errors.New("one of the arguments is not a digit")
	}

	num, _ := strconv.Atoi(string(digit1) + string(digit2))

	return num, nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetFirstAndLastDigit(line string) (rune, rune) {
	var first rune
	var curr rune
	var last rune
	firstSet := false

	for _, c := range line {
		if unicode.IsDigit(c) {
			if !firstSet {
				first = c
				firstSet = true
			}
			curr = c
		}
	}
	last = curr
	return first, last
}

func star1(lines []string) {
	total := 0
	for _, line := range lines {
		first, last := GetFirstAndLastDigit(line)
		num, err := numberFromTwoRunes(first, last)
		panicIfErr(err)
		total += num
	}

	fmt.Printf("The total result for star1 is: %d\n", total)
}

func isAtCurrentPos(s string, i int, search string) bool {

	length := len(search)

	if i+length > len(s) {
		return false
	}

	if search == s[i:i+length] {
		return true
	}

	return false
}

func GetFindingForStringDigit(line string, i int) *DigitFinding {
	var finding *DigitFinding
	switch {
	case isAtCurrentPos(line, i, "one"):
		finding = &DigitFinding{pos: i, num: 1}

	case isAtCurrentPos(line, i, "two"):
		finding = &DigitFinding{pos: i, num: 2}

	case isAtCurrentPos(line, i, "three"):
		finding = &DigitFinding{pos: i, num: 3}

	case isAtCurrentPos(line, i, "four"):
		finding = &DigitFinding{pos: i, num: 4}

	case isAtCurrentPos(line, i, "five"):
		finding = &DigitFinding{pos: i, num: 5}

	case isAtCurrentPos(line, i, "six"):
		finding = &DigitFinding{pos: i, num: 6}

	case isAtCurrentPos(line, i, "seven"):
		finding = &DigitFinding{pos: i, num: 7}

	case isAtCurrentPos(line, i, "eight"):
		finding = &DigitFinding{pos: i, num: 8}

	case isAtCurrentPos(line, i, "nine"):
		finding = &DigitFinding{pos: i, num: 9}
	}
	return finding
}

type DigitFinding struct {
	pos int
	num int
}

type Findings []DigitFinding

func (f Findings) Less(i, j int) bool { return f[i].pos < f[j].pos }
func (f Findings) Len() int           { return len(f) }
func (f Findings) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func star2(lines []string) {

	total := 0
	for _, line := range lines {
		findings := make(Findings, 0)
		for i, c := range line {
			if unicode.IsDigit(c) {
				n, _ := strconv.Atoi(string(c))
				finding := DigitFinding{pos: i, num: n}
				findings = append(findings, finding)
				continue
			}

			finding := GetFindingForStringDigit(line, i)
			if finding != nil {
				findings = append(findings, *finding)
			}

		}
		sort.Sort(findings)

		var first, last int

		if len(findings) > 1 {
			first = findings[0].num
			last = findings[len(findings)-1].num
		} else {
			first = findings[0].num
			last = first
		}

		twoDigit := fmt.Sprintf("%d%d", first, last)
		num, err := strconv.Atoi(twoDigit)

		panicIfErr(err)

		total += num
	}

	fmt.Printf("The result of day 1 - star 2 is %d\n", total)
}

func main() {
	var filename Filepath = "input"

	lines, err := filename.lines()

	panicIfErr(err)

	star1(lines)
	star2(lines)

}
