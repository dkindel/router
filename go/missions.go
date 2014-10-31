
package main

import (
	"fmt"
)
var _ = fmt.Print

// TODO: Refactor code so it isn't SR2-specific. But this will work for now.

// Gives an instance of an SR2 game state tracker
func NewStoryState() *StoryState {
	return &StoryState{
		make(map[*Mission]bool),
		NewStory(),
		0,
	}
}

type StoryState struct {
	// If the mission is not in the map it has not been started. If the mission is in the
	// map and the value is false the mission has been started but not fully completed.
	// Else the mission has been completed.
	// TODO: Should this instead be a number? Ie are individual mission linearly-tracked rather
	// than start/finished?
	Missions map[*Mission]bool
	St Story
	Respect int
}

func (s *StoryState) GetAvailable() []*Mission {
	tmp := s.St.All()
	ret := make([]*Mission, 0, len(tmp))
next:
	for _, v := range tmp {
		// Have we already done it?
		if s.Missions[v] {
			continue next
		}
		// Find out what it means to me
		if s.Respect < v.MinRespect {
			continue next
		}
		// Check prereqs
		for _, pre := range v.Pre {
			if !s.Missions[pre] {
				// Didn't do one of them
				continue next
			}
		}
		// We can do this Mission!
		ret = append(ret, v)
	}
	return ret
}

func (s*StoryState) Complete(m *Mission, repReward int) {
	if finished, started := s.Missions[m]; finished {
		fmt.Errorf("%v has already been added to the story state!\n", m)
		return
	} else if !started {
		//fmt.Println(m, " wasn't even started yet. You're so good at this game!")
	}
	s.Missions[m] = true
	s.Respect += repReward
}


func NewStory() Story {
	return sr2{}
}

type sr2 struct {}
func (sr2) All() []*Mission {
	return all
}

var (
	m0 = NM("Mission 1", 0)
	m1 = NM("Optional 1", 0)
	m2 = NM("Mission 2", 0, m0)
	m3 = NM("Optional 2", 10, m1)
	m4 = NM("Sidequest 1", 100, m2, m3)
	m5 = NM("Mission 3", 10, m2)
	all = []*Mission{m0,m1,m2,m3,m4,m5}
)

func NM(name string, minrespect int, pre...*Mission) *Mission {
	return &Mission { name, minrespect, pre }
}

type Mission struct {
	Name string
	MinRespect int
	Pre []*Mission
	//TODO: MissionData : tracks stuff about doing the mission
	// Might be per-story
}
func (m*Mission) String() string {
	return m.Name
}

type Story interface {
	All() []*Mission
}
