package src

type RarityCounter = map[int]int

type HasDraw struct {
	Limited int `json:"limited"`
	Normal  int `json:"normal"`
}

type Analysis struct {
	Summary RarityCounter            `json:"summary"`
	Pools   map[string]RarityCounter `json:"pools"`
	HasDraw HasDraw                  `json:"hasDraw"`
}

type Analyst interface {
	Analyze(uid string)
	Analysis(uid string) Analysis
}
