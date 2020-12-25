package gopool

type Worker struct {
	Name  string
	State string
}

func (wr *Worker) Run(jobCh TaskQueue, resultCh ResultQueue, firstTask Task) {
	if firstTask != nil {
		firstTask.Do()
		resultCh <- firstTask
	}

	for job := range jobCh {
		job.Do()
		resultCh <- job
	}
}
