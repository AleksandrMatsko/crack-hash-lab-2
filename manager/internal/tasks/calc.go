package tasks

import (
	"log"
	"time"
)

type Task struct {
	StartIndex uint64
	PartCount  uint64
	Done       bool
	StartedAt  time.Time

	// TaskIdx is the index of task in storage.RequestMetadata.Tasks array
	TaskIdx int
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

func CalcTasks(
	logger *log.Logger,
	alphabetLength int,
	maxLength int,
	numParts int,
	startIndex uint64,
	usePartCount bool,
	partCount uint64,
) []Task {

	var totalWords uint64
	if usePartCount {
		totalWords = partCount
	} else {
		maxTotalWords := calcTotalWordsCount(uint64(alphabetLength), uint64(maxLength))
		totalWords = maxTotalWords - startIndex
	}

	logger.Printf("total words count = %v", totalWords)

	basePartCount := totalWords / uint64(numParts)
	rest := totalWords % uint64(numParts)
	tasks := make([]Task, numParts)
	for i := range tasks {
		tasks[i] = Task{
			StartIndex: startIndex,
			PartCount:  basePartCount,
			StartedAt:  time.Now(),
			TaskIdx:    i,
		}
		if uint64(i) < rest {
			tasks[i].PartCount += 1
		}
		startIndex += tasks[i].PartCount
	}
	return tasks
}

func CalcTasksWithFixedLength(
	alphabetLength int,
	maxLength int,
	maxTaskSize uint64,
) []Task {
	totalWordsCount := calcTotalWordsCount(uint64(alphabetLength), uint64(maxLength))

	numTasks := totalWordsCount / maxTaskSize
	rest := totalWordsCount % maxTaskSize
	if rest != 0 {
		numTasks += 1
	}
	var startIndex uint64
	tasks := make([]Task, numTasks)
	var i uint64
	for i = 0; i < numTasks; i++ {
		if i == numTasks-1 {
			tasks[i] = Task{
				StartIndex: startIndex,
				PartCount:  rest,
				Done:       false,
				TaskIdx:    int(i),
				StartedAt:  time.Now(),
			}
		} else {
			tasks[i] = Task{
				StartIndex: startIndex,
				PartCount:  maxTaskSize,
				Done:       false,
				TaskIdx:    int(i),
				StartedAt:  time.Now(),
			}
		}
		startIndex += tasks[i].PartCount
	}
	return tasks
}

func CalcTasksWithNumWorkers(
	alphabetLength int,
	maxLength int,
	numWorkers uint64,
	k uint64,
) []Task {
	totalWordsCount := calcTotalWordsCount(uint64(alphabetLength), uint64(maxLength))

	numParts := numWorkers * k

	basePartCount := totalWordsCount / numParts
	rest := totalWordsCount % numParts

	var startIndex uint64
	tasks := make([]Task, numParts)
	for i := range tasks {
		tasks[i] = Task{
			StartIndex: startIndex,
			PartCount:  basePartCount,
			StartedAt:  time.Now(),
			TaskIdx:    i,
		}
		if uint64(i) < rest {
			tasks[i].PartCount += 1
		}
		startIndex += tasks[i].PartCount
	}
	return tasks
}
