package implements

import (
	"GachaServerGin/src"
	"GachaServerGin/tools"
	"github.com/samber/lo"
	"math"
)

type MemAnalyze struct {
	data         src.GachaData
	analyses     tools.DefaultDict[string, src.Analysis]
	limitedPools []string
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

	hd := src.HasDraw{}
	var (
		endLimited = false
		endNormal  = false
	)
	for _, gacha := range data {
		for _, char := range gacha.Chars {
			if lo.IndexOf(an.limitedPools, gacha.Pool) != -1 {
				if !endLimited {
					if char.Rarity == 5 {
						endLimited = true
						break
					}
					hd.Limited += 1
				}
			} else {
				if !endNormal {
					if char.Rarity == 5 {
						endNormal = true
						break
					}
					hd.Normal += 1
				}
			}
			src.Logger.Infof("%s %s %+v", gacha.Pool, char.Name, hd)
		}
		if endLimited && endNormal {
			break
		}
	}
	analyses.HasDraw = hd

	an.analyses.Set(uid, analyses)
}

func (an MemAnalyze) Analysis(uid string) src.Analysis {
	return an.analyses.Get(uid)
}

func NewMemAnalyst(data src.GachaData) src.Analyst {
	v := MemAnalyze{
		data: data,
		analyses: tools.NewDefaultDict[string, src.Analysis](func() src.Analysis {
			return src.Analysis{
				Summary: src.RarityCounter{},
				Pools:   map[string]src.RarityCounter{},
				HasDraw: src.HasDraw{},
			}
		}),
		limitedPools: src.GetLimitedPools(),
	}
	src.Logger.Infof("limited pools: %+v", v.limitedPools)
	return v
}

func (an MemAnalyze) UpdateLimitedPools() {
	an.limitedPools = src.GetLimitedPools()
	panic("TODO: how to re analyze ")
}
