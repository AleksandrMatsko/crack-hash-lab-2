package processing

import (
	"context"
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

	var reqs []storage.RequestMetadata
	var err error
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
		timeout := CalcTimeout(reqs[i].Tasks[0].PartCount)
		go RequestChecker(
			ctx,
			log.New(
				logger.Writer(),
				fmt.Sprintf("restorer for %v: ", reqs[i].ID),
				logger.Flags()),
			reqs[i].ID,
			timeout,
			S)
	}
}
