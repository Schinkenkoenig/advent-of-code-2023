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

func readSpaceSeperatedNumbers(numberStr string) []int {

	numbersAsStrList := strings.Split(numberStr, " ")

	numbers := make([]int, 0)

	for _, n := range numbersAsStrList {
		num, err := strconv.Atoi(n)
		panicIfErr(err)
		numbers = append(numbers, num)
	}

	return numbers
}

func readSeeds(s string) []int {
	numberStr := strings.ReplaceAll(s, "seeds: ", "")
	return readSpaceSeperatedNumbers(numberStr)
}

type Mappings map[int][]int

type SeedMap struct {
	mapName           string
	mappings          Mappings
	sourceDestMapping map[int]int
}

func GetFinalValueFromMappings(number int, maps []SeedMap) int {
	currNumber := number
	for _, smap := range maps {
		for _, m := range smap.mappings {
			destStart := m[0]
			sourceStart := m[1]
			mapRange := m[2]

			if sourceStart <= currNumber && currNumber < sourceStart+mapRange {
				currNumber = destStart + (currNumber - sourceStart)

				break
			}
		}

	}

	return currNumber
}

func GetSeedNumberFromLocation(location int, maps []SeedMap) int {

	currNumber := location
	for i := 0; i < len(maps); i++ {
		smap := maps[len(maps)-1-i]
		for _, m := range smap.mappings {
			destStart := m[1]
			sourceStart := m[0]
			mapRange := m[2]

			if sourceStart <= currNumber && currNumber < sourceStart+mapRange {
				currNumber = destStart + (currNumber - sourceStart)

				break
			}
		}

	}

	return currNumber
}

func main() {
	var filename Filepath = "day5/input"

	lines, err := filename.lines()

	panicIfErr(err)

	seeds := readSeeds(lines[0])
	maps := make([]SeedMap, 0)

	for i := 2; i < len(lines); {

		mapName := strings.ReplaceAll(lines[i], " map:", "")
		mappings := make(Mappings, 0)
		i++
		for lines[i] != "" {
			numbers := readSpaceSeperatedNumbers(lines[i])
			mappings[len(mappings)] = numbers
			i++
		}
		sm := SeedMap{mapName: mapName, mappings: mappings, sourceDestMapping: make(map[int]int, 0)}

		maps = append(maps, sm)
		i++
	}

	fmt.Println("Populated seed maps.")

	lowestLocation := math.MaxInt32
	for _, s := range seeds {
		location := GetFinalValueFromMappings(s, maps)

		if lowestLocation > location {
			lowestLocation = location
		}
		fmt.Printf("Lowest loc: %d\n", lowestLocation)
	}

	fmt.Println("Computed location")
	fmt.Printf("Lowest location is %d\n", lowestLocation)

	lowestLocation = math.MaxInt32
	location := 0
while:
	for {
		current := location
		location += 1000
		for i := 0; i < len(seeds); i += 2 {
			seed := GetSeedNumberFromLocation(current, maps)
			if seeds[i] <= seed && seed < seeds[i]+seeds[i+1] && GetFinalValueFromMappings(seed, maps) == current {
				lowestLocation = current
				break while
			}
		}
	}
	fmt.Printf("Lowest location estimate: %d\n", lowestLocation)

	for current := lowestLocation - 1000; current < lowestLocation; current++ {
		for i := 0; i < len(seeds); i += 2 {
			seed := GetSeedNumberFromLocation(current, maps)
			if seeds[i] <= seed && seed < seeds[i]+seeds[i+1] && GetFinalValueFromMappings(seed, maps) == current {
				if current < lowestLocation {
					lowestLocation = current
				}
			}
		}
	}

	fmt.Printf("Lowest location estimate: %d\n", lowestLocation)
}
