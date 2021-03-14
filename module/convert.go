package module

import (
	"content-pipline-example/common"
	"log"
	"sync"
	"time"
)

type ConvertWavModule struct {
	inputStream <-chan common.JobContext
	outputStream chan common.JobContext
	errorStream chan common.JobContext
}

func NewConvertWavModule(inputStream <-chan common.JobContext, errorStream chan common.JobContext)(*ConvertWavModule,error){
	return &ConvertWavModule{
		inputStream: inputStream,
		outputStream: make(chan common.JobContext),
		errorStream: errorStream,
	}, nil
}
func (convertWavModule *ConvertWavModule)Run(){

	log.Println("ConvertWavModule Run ===")

	go func() {
		defer func() {
			log.Println("ConvertWavModule End ===")
			close(convertWavModule.outputStream)
		}()

		var wg sync.WaitGroup
		wg.Add(common.WorkerSize)

		for i := 0 ; i< common.WorkerSize; i++ {
			go func(workerId int){
				defer func() {
					wg.Done()
					log.Printf("Download Worker [%d] End\n", workerId)
				}()
				for {
					select {
						case job, more := <- convertWavModule.inputStream:
							if more {
								log.Printf("Worker [%d] Converting..........\n",workerId)
								time.Sleep(time.Duration(common.RandSec())*time.Second)
								job.ConvertWavSuccess = true
								log.Printf("Worker [%d] Convert end.\n",workerId)
								convertWavModule.outputStream <- job
							}else{
								return
							}
					}
				}
			}(i)
		}
		log.Println("Worker Create end")
		wg.Wait()
	}()
}

func (convertWavModule *ConvertWavModule)GetOutputStream()chan common.JobContext{
	return convertWavModule.outputStream
}



