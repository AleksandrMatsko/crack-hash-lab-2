package balance

import (
	tasks2 "distributed.systems.labs/manager/internal/tasks"
)

type BalancedTasks map[string][]tasks2.Task

func Balance(workers []string, tasks []tasks2.Task) BalancedTasks {
	balanced := make(BalancedTasks)
	for i := range tasks {
		worker := workers[i%len(workers)]
		val, ok := balanced[worker]
		if !ok {
			tasksForWorker := make([]tasks2.Task, 1)
			tasksForWorker[0] = tasks[i]
			balanced[worker] = tasksForWorker
		} else {
			balanced[worker] = append(val, tasks[i])
		}
	}
	return balanced
}
