package main

import (
	"encoding/json"
	"os"
	"sort"
)

type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

const maxScores = 10

func LoadScores(path string) []Score {
	var s []Score
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &s)
	}
	return s
}

func SaveScores(path string, s []Score) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// Qualifies reports whether score would enter the table.
func Qualifies(s []Score, score int) bool {
	return score > 0 && (len(s) < maxScores || score > s[len(s)-1].Score)
}

// Insert adds a score, keeps the table sorted and trimmed, and returns its rank (0-based).
func Insert(s []Score, e Score) ([]Score, int) {
	s = append(s, e)
	sort.SliceStable(s, func(i, j int) bool { return s[i].Score > s[j].Score })
	rank := -1
	for i := range s {
		if s[i] == e {
			rank = i
			break
		}
	}
	if len(s) > maxScores {
		s = s[:maxScores]
	}
	return s, rank
}
