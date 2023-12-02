package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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

type Game struct {
	index        int
	subsets      []Subset
	minimumGreen int
	minimumRed   int
	minimumBlue  int
}

func (g *Game) IsPossible() bool {
	for _, subset := range g.subsets {
		if !subset.IsPossible() {
			return false
		}
	}
	return true
}

func (g *Game) ComputeMinimumCubes() {
	for _, subset := range g.subsets {
		for _, draw := range subset.draws {
			if draw.color == "blue" && g.minimumBlue < draw.count {
				g.minimumBlue = draw.count
			}
			if draw.color == "red" && g.minimumGreen < draw.count {
				g.minimumGreen = draw.count
			}
			if draw.color == "green" && g.minimumRed < draw.count {
				g.minimumRed = draw.count
			}
		}
	}
}

func (g *Game) ComputePowerCubes() int {
	return g.minimumRed * g.minimumBlue * g.minimumGreen
}

type Subset struct {
	draws []Draw
}

func (ss *Subset) IsPossible() bool {
	for _, draw := range ss.draws {
		if !draw.IsPossible() {
			return false
		}
	}
	return true
}

type Draw struct {
	color string
	count int
}

func (d *Draw) IsPossible() bool {
	switch {
	case d.color == "red":
		return d.count <= 12
	case d.color == "green":
		return d.count <= 13
	case d.color == "blue":
		return d.count <= 14
	}
	return false
}

func star1(lines []string) {
	games := make([]Game, 0)
	for _, line := range lines {
		s := strings.Split(line, ":")

		gamePart := strings.TrimSpace(s[0])
		subsetsString := strings.TrimSpace(s[1])

		gameIdx, err := strconv.Atoi(strings.ReplaceAll(gamePart, "Game ", ""))
		panicIfErr(err)

		game := Game{
			index:        gameIdx,
			minimumBlue:  math.MinInt32,
			minimumGreen: math.MinInt32,
			minimumRed:   math.MinInt32,
			subsets:      make([]Subset, 0)}

		ssStrings := strings.Split(subsetsString, ";")

		for _, drawings := range ssStrings {
			subset := Subset{draws: make([]Draw, 0)}
			s := strings.Split(drawings, ",")

			for _, draw := range s {
				trimmedDraw := strings.TrimSpace(draw)
				s := strings.Split(trimmedDraw, " ")
				count, err := strconv.Atoi(s[0])
				panicIfErr(err)

				d := Draw{color: s[1], count: count}
				subset.draws = append(subset.draws, d)
			}

			game.subsets = append(game.subsets, subset)
		}

		game.ComputeMinimumCubes()
		fmt.Println(line)
		fmt.Printf("red: %d, green %d, blue: %d, power: %d\n", game.minimumRed, game.minimumGreen, game.minimumBlue, game.ComputePowerCubes())
		games = append(games, game)
	}

	totalMaxCubes := 0
	total := 0
	for _, game := range games {
		if game.IsPossible() {
			total += game.index
		}
		totalMaxCubes += game.ComputePowerCubes()
	}

	fmt.Printf("The result for star1 is: %d\n", total)
	fmt.Printf("The result for star2 is: %d\n", totalMaxCubes)
}

func main() {
	var filename Filepath = "day2/input"

	lines, err := filename.lines()

	panicIfErr(err)

	star1(lines)
}
