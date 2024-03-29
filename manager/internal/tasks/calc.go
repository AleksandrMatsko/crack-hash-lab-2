package tasks

import (
	"time"
)

type Task struct {
	StartIndex uint64
	PartCount  uint64
	Done       bool
	StartedAt  time.Time
}

func pow[T uint | uint8 | uint16 | uint32 | uint64](base, exp T) T {
	var result T = 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}

func calcTotalWordsCount(lenAlphabet uint64, maxLength uint64) uint64 {
	return lenAlphabet * (pow[uint64](lenAlphabet, maxLength) - 1) / (lenAlphabet - 1)
}

func CalcTasksWithFixedNumParts(
	alphabetLength int,
	maxLength int,
	numParts uint64,
) []Task {
	totalWordsCount := calcTotalWordsCount(uint64(alphabetLength), uint64(maxLength))

	basePartCount := totalWordsCount / numParts
	rest := totalWordsCount % numParts

	var startIndex uint64
	tasks := make([]Task, numParts)
	for i := range tasks {
		tasks[i] = Task{
			StartIndex: startIndex,
			PartCount:  basePartCount,
			StartedAt:  time.Now(),
		}
		if uint64(i) < rest {
			tasks[i].PartCount += 1
		}
		startIndex += tasks[i].PartCount
	}
	return tasks
}
