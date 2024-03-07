package notify

import (
	"context"
	"distributed.systems.labs/shared/pkg/contracts"
	"encoding/json"
	"fmt"
	"log"
)

type ManagerNotifier struct {
	ctx      context.Context
	resChan  chan contracts.TaskResultRequest
	sendChan chan<- []byte
}

func InitManagerNotifier(ctx context.Context, sendChan chan<- []byte) *ManagerNotifier {
	return &ManagerNotifier{
		ctx:      ctx,
		resChan:  make(chan contracts.TaskResultRequest),
		sendChan: sendChan,
	}
}

func (mn *ManagerNotifier) Close() {
	close(mn.resChan)
}

func (mn *ManagerNotifier) ListenAndNotify() {
	for {
		select {
		case <-mn.ctx.Done():
			log.Println("shut down notifier")
			return

		case res := <-mn.resChan:
			defaultLogger := log.Default()
			logger := log.New(
				defaultLogger.Writer(),
				fmt.Sprintf("request-id: %s startIndex = %v ", res.RequestID, res.StartIndex),
				defaultLogger.Flags()|log.Lmsgprefix)

			logger.Printf("received res.Cracks: %v", res.Cracks)
			go func(request contracts.TaskResultRequest, logger *log.Logger) {
				reqBytes, err := json.Marshal(request)
				if err != nil {
					logger.Printf("failed to marshal request to json: %s", err)
					return
				}

				mn.sendChan <- reqBytes
			}(res, logger)
		}
	}
}

func (mn *ManagerNotifier) Context() context.Context {
	return mn.ctx
}

func (mn *ManagerNotifier) GetResChan() chan<- contracts.TaskResultRequest {
	return mn.resChan
}
