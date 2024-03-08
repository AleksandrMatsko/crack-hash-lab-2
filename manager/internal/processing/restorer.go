package processing

import (
	"context"
	"distributed.systems.labs/manager/internal/config"
	"distributed.systems.labs/manager/internal/storage"
	"fmt"
	"log"
)

func Restorer(ctx context.Context, S storage.Storage) {
	defaultLogger := log.Default()
	logger := log.New(
		defaultLogger.Writer(),
		fmt.Sprintf("restorer: "),
		defaultLogger.Flags()|log.Lmsgprefix,
	)

	workers, err := config.GetWorkers()
	if err != nil {
		logger.Printf("failed to get workers: %s", err)
		return
	}
	numWorkers := len(workers)
	if numWorkers == 0 {
		logger.Printf("no workers")
		return
	}

	var reqs []storage.RequestMetadata
	switch s := S.(type) {
	case *storage.MongoStorage:
		reqs, err = s.GetInProgressRequests()
		if err != nil {
			logger.Printf("error occured while getting in progress requests: %s", err)
			return
		}
	default:
		return
	}
	logger.Printf("found %v requests with status = IN_PROGRESS", len(reqs))
	for i := range reqs {
		timeout := CalcTimeoutsWithNumWorkers(reqs[i].Tasks[0].PartCount, uint64(numWorkers))
		go RequestChecker(ctx, logger, reqs[i].ID, timeout, S)
	}
}
