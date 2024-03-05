package main

import "fmt"

type RailroadWithChecker interface {
	CheckRailsWidth() int
	// CallForHelp() string // if a new method is added, the Train struct needs to implement it
}

type Railroad struct {
	Width int
}

// the Railroad struct will receive a point of RailroadWithChecker interface (Any struct that implements CheckRailsWidth method will be allowed as arg)
func (rr *Railroad) IsTrainWidthValid(rc RailroadWithChecker) bool {
	train := rc.CheckRailsWidth()

	return train <= rr.Width
}

type Train struct {
	Width int
}

// Train struct implements the interface method
func (t *Train) CheckRailsWidth() int {
	return t.Width
}

// func (t *Train) CallForHelp() string {
// 	return "HEEEELP"
// }

func main() {

	railroad := Railroad{Width: 10}

	highSpeedTrain := &Train{Width: 12}
	maglevTrain := &Train{Width: 5}

	fmt.Printf("Can train High Speed Train pass? - %v\n", railroad.IsTrainWidthValid(highSpeedTrain))
	fmt.Printf("Can train Mag Lev Train pass...? - %v\n", railroad.IsTrainWidthValid(maglevTrain))

}
