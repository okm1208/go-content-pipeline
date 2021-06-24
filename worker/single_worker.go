package worker

import (
	"content-pipline-example/common"
	"log"
	"time"
)

type action func(jobContext common.JobContext)(common.JobContext,error)

type Worker struct {
	inputStream <-chan common.JobContext
	outputStream chan common.JobContext
	errorStream chan common.JobContext
	fn action
	name string
}
func NewWorker(inputStream <-chan common.JobContext,
	errorStream chan common.JobContext,
	fn action,
	name string,
	)(*Worker , error){
	return &Worker{
		inputStream: inputStream,
		outputStream: make(chan common.JobContext),
		errorStream: errorStream,
		fn: fn,
		name : name,
	}, nil
}
func (worker *Worker) Run() {
	go func() {
		defer func() {
			log.Println("Worker End ===")
			close(worker.outputStream)
		}()

		for {
			select{
				case job, more := <- worker.inputStream:
					if more {
						log.Printf("Worker[%s] running.\n",worker.name)
						time.Sleep(time.Duration(common.RandSec())*time.Second)
						job,err := worker.fn(job)
						if err != nil {
							job.Error = err
							worker.errorStream <- job
						}else{
							log.Printf("Worker[%s] end.\n",worker.name)
							worker.outputStream <- job
						}
					} else {
						return
					}
			}
		}
	}()
}

func (worker *Worker)GetOutputStream()chan common.JobContext{
	return worker.outputStream
}
