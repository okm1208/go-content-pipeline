package module

import (
	"content-pipline-example/common"
	"log"
	"sync"
	"time"
)

type ExtractHash struct {
	inputStream <-chan common.JobContext
	outputStream chan common.JobContext
	errorStream chan common.JobContext
}

func NewExtractHash(inputStream <-chan common.JobContext, errorStream chan common.JobContext)(*ExtractHash , error ){
	return &ExtractHash{
		inputStream: inputStream,
		outputStream: make(chan common.JobContext),
		errorStream: errorStream,
	}, nil
}

func (extractHash *ExtractHash)Run() {
	log.Println("ExtractHash Run ===")

	go func() {
		defer func() {
			log.Println("ExtractHash End ===")
			close(extractHash.outputStream)
		}()
		var wg sync.WaitGroup
		wg.Add(common.WorkerSize)
		for i := 0; i < common.WorkerSize; i++ {
			go func(workerId int) {
				defer func() {
					wg.Done()
					log.Printf("Extract Worker [%d] End\n",workerId)
				}()
				log.Printf("Extract Worker [%d] Start\n", workerId)
				for {
					select {
					case job, more := <- extractHash.inputStream:
						if more {
							log.Printf("Worker [%d] Extracting..........\n",workerId)
							time.Sleep(time.Duration(common.RandSec())*time.Second)
							job.ExtractHashSuccess = true
							log.Printf("Worker [%d] Extract end.\n",workerId)
							extractHash.outputStream <- job
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


func (extractHash *ExtractHash)GetOutputStream()chan common.JobContext{
	return extractHash.outputStream
}