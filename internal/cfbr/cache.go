package cfbr

import "github.com/RecursionExcursion/cfbr-core-go/cfbrcore"

type seasonCache struct {
	year       int
	seasonData ApiSeasonData
}

func (sc *seasonCache) getSeasonInfo() cfbrcore.SeasonInfo {
	return sc.seasonData.SeasonInfo
}

func (sc *seasonCache) getGameData(ids []string) (map[string]cfbrcore.GameData, []string) {
	cache := map[string]cfbrcore.GameData{}
	missing := []string{}

	for _, id := range ids {
		gm, ok := sc.seasonData.GameData[id]
		if ok {
			cache[id] = gm
		} else {
			missing = append(missing, id)
		}
	}

	return cache, missing
}

func (sc *seasonCache) set(sd ApiSeasonData) {
	sc.seasonData = sd
	sc.year = sd.SeasonInfo.Year
}

func (sc *seasonCache) setIfEmpty(sd ApiSeasonData) {
	if sc.year == 0 {
		sc.set(sd)
	}
}

func (sc *seasonCache) get(yr int) (ApiSeasonData, bool) {
	return sc.seasonData, sc.year != 0
}

func (sc *seasonCache) isSet() bool {
	return sc.year != 0
}
