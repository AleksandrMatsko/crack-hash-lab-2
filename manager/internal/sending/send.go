package sending

import (
	"distributed.systems.labs/manager/internal/config"
	"distributed.systems.labs/manager/internal/storage"
	"distributed.systems.labs/manager/internal/tasks"
	"distributed.systems.labs/shared/pkg/contracts"
	"encoding/json"
	"log"
)

func sendToWorker(req contracts.TaskRequest) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	config.GetToSendChan() <- reqBytes
	return nil
}

func SendLoop(logger *log.Logger, toSendTasks []tasks.Task, m storage.RequestMetadata) {
	for i, t := range toSendTasks {
		req := contracts.TaskRequest{
			StartIndex: t.StartIndex,
			PartCount:  t.PartCount,
			Alphabet:   m.Alphabet.ToOneLine(),
			MaxLength:  m.MaxLength,
			ToCrack:    m.Hash,
			RequestID:  m.ID,
		}
		err := sendToWorker(req)
		if err != nil {
			logger.Printf("failed to send task %v to workers: %s", i, err)
		}
	}
}
