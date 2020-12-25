package gopool

import "log"

type Worker struct {
	Name  string
	State string
}

func (wr *Worker) Run(jobCh TaskQueue, resultCh ResultQueue, firstTask Task) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("task run failed! | err | %v", err)
		}
	}()

	if firstTask != nil {
		log.Printf("begin consum | task | %v", firstTask)
		firstTask.Do()
		resultCh <- firstTask
		log.Printf("consum | task | %v", firstTask)
	}

	for job := range jobCh {
		log.Printf("begin consum | task | %v", firstTask)
		job.Do()
		resultCh <- job
		log.Printf("consum | task | %v", job)
	}
}
