package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type Ranker interface {
	Ranking() []string
}

type Team struct {
	Name    string
	Players []string
}

type League struct {
	Teams []Team
	Wins  map[string]int
}

func (l *League) MatchResult(t1 string, t1S int, t2 string, t2S int) {
	if t1S > t2S {
		l.Wins[t1]++
	}
	if t1S < t2S {
		l.Wins[t2]++
	}
}

func (l League) Ranking() []string {
	names := make([]string, 0, len(l.Teams))
	for _, t := range l.Teams {
		names = append(names, t.Name)
	}

	sort.Slice(names, func(i, j int) bool {
		wi := l.Wins[names[i]]
		wj := l.Wins[names[j]]

		if wi != wj {
			return wi > wj
		}

		return names[i] < names[j]
	})

	return names
}

func RankPrinter(p Ranker, w io.Writer) error {
	for _, name := range p.Ranking() {
		if _, err := io.WriteString(w, name+"\n"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	l := League{
		Teams: []Team{
			{Name: "Bulls", Players: []string{"A", "B"}},
			{Name: "Lakers", Players: []string{"C", "D"}},
			{Name: "Celtics", Players: []string{"E", "F"}},
		},
		Wins: make(map[string]int),
	}

	l.MatchResult("Bulls", 1, "Lakers", 2)
	l.MatchResult("Celtics", 1, "Lakers", 0)
	l.MatchResult("Bulls", 1, "Celtics", 1) // ничья

	fmt.Println("Wins map:", l.Wins)
	fmt.Println("Ranking:")
	if err := RankPrinter(l, os.Stdout); err != nil {
		fmt.Println("write error:", err)
	}
}
