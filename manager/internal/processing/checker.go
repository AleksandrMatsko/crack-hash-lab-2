package processing

import (
	"context"
	"distributed.systems.labs/manager/internal/config"
	"distributed.systems.labs/manager/internal/sending"
	"distributed.systems.labs/manager/internal/storage"
	"distributed.systems.labs/manager/internal/tasks"
	"github.com/google/uuid"
	"log"
	"time"
)

func RequestChecker(ctx context.Context, logger *log.Logger, requestID uuid.UUID, timeout time.Duration, S storage.Storage) {
	logger.Printf("starting request checker with timeout %v", timeout)
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			logger.Println("request checker stops from context")

		case <-timer.C:
			logger.Printf("request checker woke up after %v", timeout)
			timeoutedTasks := make([]tasks.Task, 0)
			var status config.RequestStatus
			m, err := S.Atomically(requestID, func(metadata *storage.RequestMetadata) {
				status = metadata.Status
				if metadata.Status == config.Ready || metadata.Status == config.Error {
					return
				}

				for _, t := range metadata.Tasks {
					if !t.Done && time.Now().Sub(t.StartedAt) > timeout {
						tsk := t
						timeoutedTasks = append(timeoutedTasks, tsk)
					}
				}
			})
			if err != nil {
				logger.Printf("request checker has err while checking timeouted tasks: %s", err)
				return
			}
			if status != config.InProgress {
				logger.Printf("request status = %s, request checker stopping...", status)
				return
			}

			if len(timeoutedTasks) == 0 {
				timer.Reset(timeout)
				continue
			}

			logger.Printf("request checker needs to rebalance %v tasks", len(timeoutedTasks))

			workers, err := config.GetWorkers()
			if err != nil {
				logger.Printf("request checker error while getting workers: %s", err)
				storage.SetStatusErrAndSave(logger, S, requestID)
				return
			}
			if len(workers) == 0 {
				logger.Printf("request checker error: no workers")
				storage.SetStatusErrAndSave(logger, S, requestID)
				return
			}

			sending.BalanceAndSendLoop(logger, workers, timeoutedTasks, S, m)

			_, _ = S.Atomically(requestID, func(req *storage.RequestMetadata) {
				status = req.Status
			})
			if status != config.Error {
				timer.Reset(timeout)
			} else {
				logger.Println("request checker exiting")
				return
			}
		}
	}

}