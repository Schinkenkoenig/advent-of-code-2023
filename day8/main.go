package main

import (
	"fmt"
	"os"
	"slices"
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
	for unicode.IsDigit(rune(t.ch)) {
		t.readChar()
	}

	return t.input[pos:t.pos]
}

type Move string

const (
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
)

type MoveOption struct {
	left  string
	right string
}

func readMoves(line string) []Move {

	t := New(line)
	moves := make([]Move, 0)
	for t.ch != 0 {
		if t.ch == 'L' {
			moves = append(moves, LEFT)
		}
		if t.ch == 'R' {
			moves = append(moves, RIGHT)
		}
		t.readChar()
	}
	return moves
}

func readLocationMoveLine(line string) (string, MoveOption) {

	t := New(line)
	// read everything until equals
	pos := t.pos
	for t.ch != '=' {
		t.readChar()
	}
	current := strings.TrimSpace(t.input[pos:t.pos])

	for t.ch != '(' {
		t.readChar()
	}
	t.readChar()

	pos = t.pos
	for t.ch != ',' {
		t.readChar()
	}
	left := strings.TrimSpace(t.input[pos:t.pos])
	t.readChar()

	pos = t.pos
	for t.ch != ')' {
		t.readChar()
	}

	right := strings.TrimSpace(t.input[pos:t.pos])

	moveOption := MoveOption{left: left, right: right}
	return current, moveOption
}

func hasOnlyLastCharZ(nodes []string) bool {
	for _, node := range nodes {
		if node[len(node)-1] != 'Z' {
			return false
		}
	}
	return true
}

func GetMovePerLocationMap(lines []string) map[string]MoveOption {

	movesPerLocation := make(map[string]MoveOption, 0)
	for _, l := range lines {
		if l == "" {
			continue
		}
		current, moveOption := readLocationMoveLine(l)

		movesPerLocation[current] = moveOption
	}
	return movesPerLocation
}

func star1() {

	var filename Filepath = "day8/input"

	lines, err := filename.lines()

	panicIfErr(err)

	movesPerLocation := GetMovePerLocationMap(lines[1:])

	moves := readMoves(lines[0])
	location := "AAA"

	moveCount := 0
	for location != "ZZZ" {
		for _, m := range moves {
			if location == "ZZZ" {
				break
			}
			moveCount++
			move := movesPerLocation[location]

			if m == LEFT {
				fmt.Printf("Move from %s to %s\n", location, move.left)
				location = move.left
				continue
			}
			if m == RIGHT {
				fmt.Printf("Move from %s to %s\n", location, move.right)
				location = move.right
				continue
			}
		}
	}

	fmt.Printf("Move count: %d\n", moveCount)
}

func moveAllNodes(nodes []string, perLocation map[string]MoveOption, m Move) {

	for i, node := range nodes {
		possibleMoves := perLocation[node]
		if m == LEFT {
			nodes[i] = possibleMoves.left
		}
		if m == RIGHT {
			nodes[i] = possibleMoves.right
		}
	}
}
func gcd(a int, b int) int {

	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// Returns LCM of array elements
func findlcm(arr []int, n int) int {
	// Initialize result
	ans := arr[0]

	for i := 1; i < n; i++ {

		ans = ((arr[i] * ans) /
			(gcd(arr[i], ans)))
	}

	return ans
}

func main() {
	star1()
	var filename Filepath = "day8/input"

	lines, err := filename.lines()

	panicIfErr(err)

	movesPerLocation := GetMovePerLocationMap(lines[1:])

	moves := readMoves(lines[0])
	moveCount := 0
	// get all nodes that end with A
	// follow every node left and right
	// only process nodes that do not currently end with Z
	// stop when every nodes ends with Z

	currentNodes := make([]string, 0)

	for node := range movesPerLocation {
		if node[len(node)-1] == 'A' {
			currentNodes = append(currentNodes, node)
			fmt.Printf("%v\n", node)
		}
	}
	slices.Sort(currentNodes)
	movesUntilZ := make([]int, len(currentNodes))

while:
	for {
		for _, m := range moves {
			if hasOnlyLastCharZ(currentNodes) {
				break while
			}
			moveCount++
			for i, node := range currentNodes {
				if node[len(node)-1] == 'Z' {
					continue
				}

				possibleMoves := movesPerLocation[node]
				if m == LEFT {
					currentNodes[i] = possibleMoves.left
				}
				if m == RIGHT {
					currentNodes[i] = possibleMoves.right
				}

				if currentNodes[i][len(currentNodes[i])-1] == 'Z' {
					movesUntilZ[i] = moveCount
				}

			}
			fmt.Printf("Current nodes %v\n", currentNodes)
		}

	}

	fmt.Printf("Move count is %d\n", moveCount)
	fmt.Printf("Move count until Z %v\n", findlcm(movesUntilZ, len(movesUntilZ)))
}
