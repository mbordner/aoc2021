package main

import "fmt"

const (
	winningScore = 1000
)

type Player struct {
	pos   int
	score int
	turns int
}

func newPlayer(start int) *Player {
	p := new(Player)
	p.pos = start
	return p
}

func (p *Player) roll(val int) {
	p.pos = (val+p.pos-1)%10 + 1
	p.score += p.pos
	p.turns++
}

type Players []*Player

func (ps *Players) maxScore() int {
	val := 0
	for i := 0; i < len(*ps); i++ {
		if (*ps)[i].score > val {
			val = (*ps)[i].score
		}
	}
	return val
}

func (ps *Players) nextPlayer(turns int) *Player {
	return (*ps)[turns%len(*ps)]
}

func sum(v int) int {
	val := v
	for i := 1; i <= 2; i++ {
		val += v + i
	}
	return val
}

func main() {

	players := Players{newPlayer(4), newPlayer(6)}

	turns := 0

	for players.maxScore() < winningScore {
		p := players.nextPlayer(turns)

		p.roll(sum(turns*3 + 1))

		fmt.Println("Player ", turns%len(players)+1, " pos ", p.pos, " with score ", p.score)

		turns++
	}

	loser := players.nextPlayer(turns)
	fmt.Println(loser.score * turns * 3)

}
