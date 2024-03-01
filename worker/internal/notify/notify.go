package notify

import (
	"context"
	"distributed.systems.labs/shared/pkg/contracts"
)

type Notifier interface {
	Close()
	ListenAndNotify()
	Context() context.Context
	GetResChan() chan<- contracts.TaskResultRequest
}
