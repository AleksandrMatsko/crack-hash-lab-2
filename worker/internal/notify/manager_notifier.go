package notify

import (
	"bytes"
	"context"
	"distributed.systems.labs/shared/pkg/contracts"
	"distributed.systems.labs/worker/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ManagerNotifier struct {
	ctx     context.Context
	resChan chan contracts.TaskResultRequest
}

func InitManagerNotifier(ctx context.Context) *ManagerNotifier {
	return &ManagerNotifier{
		ctx:     ctx,
		resChan: make(chan contracts.TaskResultRequest),
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
				hostAndPort, err := config.GetManagerHostAndPort()
				if err != nil {
					logger.Printf("failed to get manager host and port: %s", err)
					return
				}

				reqBytes, err := json.Marshal(request)
				if err != nil {
					logger.Printf("failed to marshal request to json: %s", err)
					return
				}

				req, err := http.NewRequest(
					http.MethodPatch,
					fmt.Sprintf("http://%s/internal/api/manager/hash/crack/request", hostAndPort),
					bytes.NewReader(reqBytes))
				if err != nil {
					logger.Printf("failed to create request: %s", err)
					return
				}

				r, err := http.DefaultClient.Do(req)
				if err != nil {
					logger.Printf("failed to PATCH request result: %s", err)
					return
				}
				logger.Printf("response status: %s", r.Status)
				if r.StatusCode != http.StatusOK {
					// TODO
				}
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
