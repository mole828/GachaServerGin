package src

type RarityCounter = map[int]int

type Analysis struct {
	Summary RarityCounter            `json:"summary"`
	Pools   map[string]RarityCounter `json:"pools"`
}

type Analyst interface {
	Analyze(uid string)
	Analysis(uid string) Analysis
}
