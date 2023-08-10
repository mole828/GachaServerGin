package implements

import (
	"GachaServerGin/src"
	"GachaServerGin/tools"
	"math"
)

type MemAnalyze struct {
	data     src.GachaData
	analyses tools.DefaultDict[string, src.Analysis]
}

func makeRarityCounter() tools.Counter[int] {
	return tools.NewCounter[int]()
}

func (an MemAnalyze) Analyze(uid string) {
	analyses := an.analyses.Get(uid)

	data := an.data.GetGachasByPage(uid, 0, math.MaxInt)
	pools := tools.NewDefaultDict[string, tools.Counter[int]](makeRarityCounter)
	for _, gacha := range data {
		pool := pools.Get(gacha.Pool)
		for _, char := range gacha.Chars {
			pool.Inc(char.Rarity, 1)
		}
	}

	summary := makeRarityCounter()
	for pool, counter := range pools.Data() {
		analyses.Pools[pool] = counter.Data()
		for key, value := range counter.Data() {
			summary.Inc(key, value)
		}
	}
	analyses.Summary = summary.Data()

	an.analyses.Set(uid, analyses)
}

func (an MemAnalyze) Analysis(uid string) src.Analysis {
	return an.analyses.Get(uid)
}

func NewMemAnalyst(data src.GachaData) src.Analyst {
	return MemAnalyze{
		data,
		tools.NewDefaultDict[string, src.Analysis](func() src.Analysis {
			return src.Analysis{
				Summary: src.RarityCounter{},
				Pools:   map[string]src.RarityCounter{},
			}
		}),
	}
}
