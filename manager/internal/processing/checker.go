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
	logger.Printf("starting checker with timeout %v", timeout)
	defer logger.Println("checker exited")
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			logger.Println("checker stopped from context")

		case <-timer.C:
			logger.Printf("checker woke up after %v", timeout)
			timeoutedTasks := make([]tasks.Task, 0)
			var status config.RequestStatus
			m, err := S.Atomically(requestID, func(metadata *storage.RequestMetadata) error {
				status = metadata.Status
				if metadata.Status == config.Ready || metadata.Status == config.Error {
					return nil
				}

				for _, t := range metadata.Tasks {
					if !t.Done && time.Now().Sub(t.StartedAt) > timeout {
						tsk := t
						timeoutedTasks = append(timeoutedTasks, tsk)
					}
				}
				return nil
			})
			if err != nil {
				logger.Printf("checker has err while checking timeouted tasks: %s", err)
				timer.Reset(timeout)
				continue
			}
			if status != config.InProgress {
				logger.Printf("request status = %s, checker stopping...", status)
				return
			}

			if len(timeoutedTasks) == 0 {
				timer.Reset(timeout)
				continue
			}

			logger.Printf("checker needs to rebalance %v tasks", len(timeoutedTasks))

			sending.SendLoop(logger, timeoutedTasks, m)

			timer.Reset(timeout)
		}
	}

}
