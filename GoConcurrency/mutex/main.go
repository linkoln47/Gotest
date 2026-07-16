package main

import (
	"fmt"
	"sync"
)

type ScoreboardManager struct {
	l  sync.RWMutex
	sb map[string]int
}

func NewScoreboardManager() *ScoreboardManager {
	return &ScoreboardManager{
		sb: map[string]int{},
	}
}

func (msm *ScoreboardManager) Update(name string, val int) {
	msm.l.Lock()
	defer msm.l.Unlock()
	msm.sb[name] = val
}

func (msm *ScoreboardManager) Read(name string) (int, bool) {
	msm.l.RLock()
	defer msm.l.RUnlock()
	val, ok := msm.sb[name]
	return val, ok
}

func main() {
	msm := NewScoreboardManager()
	teams := []string{"Lions", "Tigers", "Bears"}
	var wg sync.WaitGroup
	wg.Add(len(teams))
	for _, v := range teams {
		go func(team string) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				curScore, found := msm.Read(team)
				if !found {
					curScore = 10
				} else {
					curScore += len(team)
				}
				msm.Update(team, curScore)
			}
		}(v)
	}

	wg.Wait()
	for _, v := range teams {
		score, found := msm.Read(v)
		fmt.Println(v, score, found)
	}
}
