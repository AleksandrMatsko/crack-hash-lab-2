package processing

import "time"

func CalcTimeout(partCount uint64) time.Duration {
	return time.Duration(partCount) * time.Microsecond * 20
}

func CalcTimeoutWithTaskCount(partCount uint64, taskCount uint64) time.Duration {
	return CalcTimeout(partCount) * time.Duration(taskCount)
}

func CalcTimeoutsWithNumWorkers(partCount uint64, numWorkers uint64) time.Duration {
	return CalcTimeout(partCount) * time.Duration(numWorkers)
}
