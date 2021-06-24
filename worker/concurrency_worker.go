package worker

import (
	"content-pipline-example/common"
	"log"
	"sync"
)


type CWorker struct {
	inputStream <-chan common.JobContext
	outputStream chan common.JobContext
	errorStream chan common.JobContext
	fn action
	name string
	workerCount int
}

func NewCWorker(inputStream <-chan common.JobContext,
	errorStream chan common.JobContext,
	fn action,
	name string,
	workerCount int ,
)(*CWorker , error){
	return &CWorker{
		inputStream: inputStream,
		outputStream: make(chan common.JobContext),
		errorStream: errorStream,
		fn: fn,
		name : name,
		workerCount: workerCount,
	}, nil
}

func (cWorker *CWorker) Run() {

	go func() {
		defer func() {
			log.Println("CWorker End ===")
			close(cWorker.outputStream)
		}()
		var wg sync.WaitGroup
		wg.Add(cWorker.workerCount)


		for i:= 0 ; i <cWorker.workerCount ; i++ {
			go func(workerId int) {
				defer func() {
					wg.Done()
					log.Printf("%s CWorker [%d] End\n",cWorker.name, workerId)
				}()
				for{
					select{
					case job, more := <- cWorker.inputStream:
						if more {
							//log.Printf("%s CWorker[%d] running.\n",cWorker.name, workerId)
							//time.Sleep(time.Duration(common.RandSec())*time.Second)
							job,err := cWorker.fn(job)
							if err != nil {
								job.Error = err
								cWorker.errorStream <- job
							}else{
								//log.Printf("%s CWorker[%d] end.\n",cWorker.name, workerId)
								cWorker.outputStream <- job
							}
						}else{
							return
						}
					}
				}
			}(i)
		}
		wg.Wait()

	}()
}

func (cWorker *CWorker)GetOutputStream()chan common.JobContext{
	return cWorker.outputStream
}
