package main

func clamp(stat float64) float64 {
	if stat < minStat {
		return minStat
	}
	if stat > maxStat {
		return maxStat
	}
	return stat
}

func daysBetween(from, to int64) float64 {
	return float64(to-from) / secondsPerDay
}

func decayed(stat float64, anchor, at int64, ratePerDay float64) float64 {
	if anchor == 0 || at <= anchor {
		return clamp(stat)
	}
	return clamp(stat - ratePerDay*daysBetween(anchor, at))
}
