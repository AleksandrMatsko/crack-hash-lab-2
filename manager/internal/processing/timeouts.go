package processing

import "time"

func CalcTimeout(partCount uint64) time.Duration {
	return time.Duration(partCount) * time.Microsecond * 10
}
