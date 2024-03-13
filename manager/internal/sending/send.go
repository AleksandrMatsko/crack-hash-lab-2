package sending

import (
	"distributed.systems.labs/manager/internal/balance"
	"distributed.systems.labs/manager/internal/config"
	"distributed.systems.labs/manager/internal/storage"
	"distributed.systems.labs/manager/internal/tasks"
	"distributed.systems.labs/shared/pkg/contracts"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
)

func sendToWorker(logger *log.Logger, hostPort string, req contracts.TaskRequest) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	/*r, err := http.Post(
		fmt.Sprintf("http://%s/internal/api/worker/hash/crack/task", hostPort),
		"application/json",
		bytes.NewReader(reqBytes))
	if err != nil {
		// TODO better handling
		return err
	}
	logger.Printf("worker %s response: %s", hostPort, r.Status)
	if r.StatusCode != http.StatusOK {
		// TODO better handling
	}*/
	config.GetToSendChan() <- reqBytes
	return nil
}

// SendToWorkers consumes BalancedTasks and sends them to related workers.
// Returns the workers to which send succeed and task that has NOT been sent/
func SendToWorkers(
	logger *log.Logger,
	balanced balance.BalancedTasks,
	m storage.RequestMetadata,
	S storage.Storage,
) ([]string, []tasks.Task) {
	wg := sync.WaitGroup{}
	wg.Add(len(balanced))

	workersOK := make([]string, 0)
	workersOKMtx := sync.Mutex{}

	notSentTasks := make([]tasks.Task, 0)
	notSentTasksMtx := sync.Mutex{}

	for hostPort, tasksQuota := range balanced {
		go func(hostPort string, tsks []tasks.Task) {
			defer wg.Done()
			for _, t := range tsks {
				req := contracts.TaskRequest{
					StartIndex: t.StartIndex,
					PartCount:  t.PartCount,
					Alphabet:   m.Alphabet.ToOneLine(),
					MaxLength:  m.MaxLength,
					ToCrack:    m.Hash,
					RequestID:  m.ID,
				}
				err := sendToWorker(logger, hostPort, req)
				if err != nil {
					logger.Printf("send to worker %s failed with err: %s", hostPort, err)
					notSentTasksMtx.Lock()
					for _, tsk := range tsks {
						_tsk := tsk
						notSentTasks = append(notSentTasks, _tsk)
					}
					notSentTasksMtx.Unlock()
					return
				}
				/*_, err = S.Atomically(m.ID, func(r *storage.RequestMetadata) error {
					r.Tasks[t.TaskIdx].StartedAt = time.Now()
					return nil
				})
				logger.Printf("Error after updating start time %s", err)*/
			}

			workersOKMtx.Lock()
			defer workersOKMtx.Unlock()
			workersOK = append(workersOK, hostPort)

		}(hostPort, tasksQuota)
	}
	wg.Wait()
	return workersOK, notSentTasks
}

func BalanceAndSendLoop(
	logger *log.Logger,
	workers []string,
	preparedTasks []tasks.Task,
	S storage.Storage,
	m storage.RequestMetadata,
) {
	builder := strings.Builder{}
	for {
		balanced := balance.Balance(workers, preparedTasks)
		builder.Reset()
		for worker, tids := range balanced {
			builder.WriteString(fmt.Sprintf("\n\tworker %s: %v", worker, tids))
		}
		logger.Println(builder.String())

		workers, preparedTasks = SendToWorkers(logger, balanced, m, S)
		if len(preparedTasks) == 0 {
			break
		} else if len(workers) == 0 {
			logger.Printf("all workers are down")
			storage.SetStatusErrAndSave(logger, S, m.ID)
			return
		}

		builder.Reset()
		builder.WriteString(fmt.Sprintf("alive workers: %v\n", workers))
		for i, t := range preparedTasks {
			_, _ = builder.WriteString(fmt.Sprintf("task %v: StartIndex = %v PartCount %v\n", i, t.StartIndex, t.PartCount))
		}
		logger.Println(builder.String())
	}
}
