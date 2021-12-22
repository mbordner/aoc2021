package main

import (
	"fmt"
	"time"
)

const (
	winningScore = 21
)

var (
	rollSums = make(map[int]int)
	start    = time.Now()
)

func init() {
	// calculate the sums from all permutations of the possible 1,2,3 die values (with repeating allowed)
	vals := []int{1, 2, 3}
	for i := 0; i < len(vals); i++ {
		for j := 0; j < len(vals); j++ {
			for k := 0; k < len(vals); k++ {
				sum := vals[i] + vals[j] + vals[k]
				if c, e := rollSums[sum]; e {
					rollSums[sum] = c + 1
				} else {
					rollSums[sum] = 1
				}
			}
		}
	}
}

type DiracDiceStateCounts map[DiracDiceState]uint64

func (s *DiracDiceStateCounts) add(newStates DiracDiceStateCounts) {
	for ns, counts := range newStates {
		if c, exists := (*s)[ns]; exists {
			(*s)[ns] = c + counts
		} else {
			(*s)[ns] = counts
		}
	}
}

func (s DiracDiceStateCounts) counts() uint64 {
	val := uint64(0)
	for _, c := range s {
		val += c
	}
	return val
}

type DiracDiceState struct {
	positions [2]byte
	scores    [2]byte
}

func newDirState(p1pos, p1score, p2pos, p2score byte) DiracDiceState {
	return DiracDiceState{[2]byte{p1pos, p2pos}, [2]byte{p1score, p2score}}
}

func (s DiracDiceState) clone() DiracDiceState {
	return newDirState(s.positions[0], s.scores[0], s.positions[1], s.scores[1])
}

func (s DiracDiceState) getMaxScore() int {
	if s.scores[0] > s.scores[1] {
		return int(s.scores[0])
	}
	return int(s.scores[1])
}

func (s DiracDiceState) roll(turn int) DiracDiceStateCounts {
	// each roll will generate 3^3 new states with permutations of the dice values 1,2,3...
	// all the sums are already calculated for these 27 new states spawning
	// based on the turn, if it's even, the current player is the 1st player
	// otherwise if it is odd, it is the 2nd player

	i := turn % 2 //current player index

	newStates := make(DiracDiceStateCounts)
	for rollSum, count := range rollSums {
		newState := s.clone()

		newState.positions[i] = (byte(rollSum)+s.positions[i]-1)%10 + 1
		newState.scores[i] += newState.positions[i]

		newStates[newState] = uint64(count)
	}

	return newStates
}

func main() {

	initState := newDirState(4, 0, 6, 0)

	states := make(DiracDiceStateCounts)
	states[initState] = 1

	// after each turn, the sum of # of wins will be appended, so turn number == len(turns)
	turns := make([]uint64, 0, 25)

	for len(states) > 0 {
		wins := uint64(0)

		nextStateCounts := make(DiracDiceStateCounts)

		for s, c := range states {

			newStates := s.roll(len(turns))

			losingStates := make(DiracDiceStateCounts)

			twins := uint64(0)
			for ts, tc := range newStates {
				scaledCount := tc * c

				if ts.getMaxScore() >= winningScore {
					twins += scaledCount
				} else {
					losingStates[ts] = scaledCount
				}
			}

			if len(losingStates) > 0 {
				nextStateCounts.add(losingStates)
			}

			wins += twins
		}

		turns = append(turns, wins)

		fmt.Println("universes after turn ", len(turns), " are ", nextStateCounts.counts(), " with ", len(nextStateCounts), " distinct states")

		states = nextStateCounts

	}

	playerWins := []uint64{0, 0}
	for i, v := range turns {
		playerWins[i%2] += v
	}

	fmt.Println(playerWins)

	duration := time.Since(start)
	fmt.Println(duration)

}
