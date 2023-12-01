package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	bytes, err := os.ReadFile("input")

	if err != nil {
		panic(err)
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")
	total := 0
	for _, line := range lines {
		var first rune
		firstSet := false
		var curr rune
		var last rune

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
		twodigit := string(first) + string(last)
		fmt.Printf("%s -> %s \n", line, twodigit)

		num, err := strconv.Atoi(twodigit)

		if err != nil {
			panic(fmt.Sprintf("%s is not a full number.", twodigit))
		}

		total += num
	}

	fmt.Printf("The total result is: %d\n", total)
}
