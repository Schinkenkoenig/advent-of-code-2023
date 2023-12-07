package main

import (
	"fmt"
	"os"
	"sort"
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

type CardType int

const (
	TWO   = 2
	THREE = 3
	FOUR  = 4
	FIVE  = 5
	SIX   = 6
	SEVEN = 7
	EIGHT = 8
	NINE  = 9
	TEN   = 10
	JACK  = 11
	JOKER = 1
	QUEEN = 12
	KING  = 13
	ACE   = 14
)

func GetCardType(char byte, jokerAllowed bool) CardType {
	switch {
	case char == '2':
		return TWO
	case char == '3':
		return THREE
	case char == '4':
		return FOUR
	case char == '5':
		return FIVE
	case char == '6':
		return SIX
	case char == '7':
		return SEVEN
	case char == '8':
		return EIGHT
	case char == '9':
		return NINE
	case char == 'T':
		return TEN
	case char == 'J':
		if jokerAllowed {
			return JOKER
		}
		return JACK
	case char == 'Q':
		return QUEEN
	case char == 'K':
		return KING
	case char == 'A':
		return ACE
	default:
		return 0
	}
}

type HandType int

const (
	HIGH_CARD  = 0
	ONE_PAIR   = 1
	TWO_PAIR   = 2
	THREE_KIND = 3
	FULL_HOUSE = 4
	FOUR_KIND  = 5
	FIVE_KIND  = 6
)

func (c HandType) ToString() string {
	switch {
	case c == ONE_PAIR:
		return "one pair"
	case c == TWO_PAIR:
		return "two pair"
	case c == THREE_KIND:
		return "three kind"
	case c == FOUR_KIND:
		return "four kind"
	case c == FIVE_KIND:
		return "five kind"
	case c == FULL_HOUSE:
		return "full house"
	}
	return "high card"
}

type HandBid struct {
	Cards    []CardType
	Bid      int
	handType HandType
}

func (h *HandBid) DetermineHandWithJoker() HandType {
	cardcount := make(map[CardType]int, 0)

	for i := 0; i < len(h.Cards); i++ {
		if _, ok := cardcount[h.Cards[i]]; !ok {
			cardcount[h.Cards[i]] = 1
		} else {
			cardcount[h.Cards[i]]++
		}
	}
	if len(cardcount) == 1 {
		return FIVE_KIND
	}

	if len(cardcount) == 2 {
		_, ok := cardcount[JOKER]

		if ok {
			return FIVE_KIND
		}

		for _, v := range cardcount {
			if v == 1 {
				return FOUR_KIND
			}
		}
		return FULL_HOUSE
	}

	if len(cardcount) == 3 {

		jokercount, ok := cardcount[JOKER]

		if ok {
			if jokercount == 3 {
				return FOUR_KIND
			}

			if jokercount == 2 {
				return FOUR_KIND
			}

			if jokercount == 1 {
				for k, v := range cardcount {
					if k == JOKER {
						continue
					}

					if v == 2 {
						return FULL_HOUSE
					}

					if v == 1 || v == 3 {
						return FOUR_KIND
					}
				}
			}
		}

		for _, v := range cardcount {
			if v == 3 {
				return THREE_KIND
			}
		}
		return TWO_PAIR
	}

	if len(cardcount) == 4 {
		_, ok := cardcount[JOKER]
		if ok {
			return THREE_KIND
		}

		return ONE_PAIR
	}

	if len(cardcount) == 5 {
		_, ok := cardcount[JOKER]
		if ok {
			return ONE_PAIR
		}
	}

	return HIGH_CARD
}

func (h *HandBid) DetermineHand() HandType {
	cardcount := make(map[CardType]int, 0)

	for i := 0; i < len(h.Cards); i++ {
		if _, ok := cardcount[h.Cards[i]]; !ok {
			cardcount[h.Cards[i]] = 1
		} else {
			cardcount[h.Cards[i]]++
		}
	}

	if len(cardcount) == 4 {
		return ONE_PAIR
	}

	if len(cardcount) == 1 {
		return FIVE_KIND
	}

	if len(cardcount) == 2 {
		for _, v := range cardcount {
			if v == 1 {
				return FOUR_KIND
			}
		}
		return FULL_HOUSE
	}

	if len(cardcount) == 3 {

		for _, v := range cardcount {
			if v == 3 {
				return THREE_KIND
			}
		}
		return TWO_PAIR
	}

	return HIGH_CARD
}

func CreateHandBid(line string, jokerAllowed bool) HandBid {
	t := New(line)

	cards := make([]CardType, 0)
	for t.ch != ' ' {
		cards = append(cards, GetCardType(t.ch, jokerAllowed))
		t.readChar()
	}
	t.readChar()

	bid, err := strconv.Atoi(t.readNumber())
	panicIfErr(err)

	return HandBid{Cards: cards, Bid: bid}
}

type Bids []HandBid

func (b Bids) Len() int {
	return len(b)
}
func (b Bids) Less(i, j int) bool {
	if b[i].handType != b[j].handType {
		return b[i].handType < b[j].handType
	}

	for k := 0; k < len(b[i].Cards); k++ {
		if b[i].Cards[k] != b[j].Cards[k] {
			return b[i].Cards[k] < b[j].Cards[k]
		}
	}

	return true
}

func (b Bids) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func main() {
	var filename Filepath = "day7/input"

	lines, err := filename.lines()

	panicIfErr(err)

	// assuming input with one final empty line
	handBids := make(Bids, len(lines)-1)
	for i, l := range lines {
		if l == "" {
			continue
		}
		handBids[i] = CreateHandBid(l, false)
	}

	for i := 0; i < len(handBids); i++ {
		c := handBids[i]
		handBids[i].handType = c.DetermineHand()
	}
	sort.Sort(handBids)

	total_winnings := 0
	for i, c := range handBids {
		winning := c.Bid * (i + 1)
		fmt.Printf("%v name %s, winnings: %d\n", c, c.handType.ToString(), winning)
		total_winnings += winning
	}

	fmt.Printf("Total winnings %d\n", total_winnings)

	handBidsJoker := make(Bids, len(lines)-1)
	for i, l := range lines {
		if l == "" {
			continue
		}
		handBidsJoker[i] = CreateHandBid(l, true)
	}

	for i := 0; i < len(handBidsJoker); i++ {
		c := handBidsJoker[i]
		handBidsJoker[i].handType = c.DetermineHandWithJoker()
	}
	sort.Sort(handBidsJoker)

	total_winnings_joker := 0
	for i, c := range handBidsJoker {
		winning := c.Bid * (i + 1)
		fmt.Printf("%v name %s, winnings: %d\n", c, c.handType.ToString(), winning)
		total_winnings_joker += winning
	}

	fmt.Printf("Total winnings %d\n", total_winnings_joker)

}
