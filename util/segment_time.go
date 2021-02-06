package util

func SegmentTime(start, end int64) (sum, startTime, endTime int64) {
	startTime = start - start%1800
	if end%1800 == 0 {
		endTime = end
	} else {
		endTime = end - end%1800 + 1800
	}
	sum = (endTime - startTime) / 1800
	return sum, startTime, endTime
}
