package module

import (
	"content-pipline-example/common"
	"log"
	"sync"
	"time"
)


type DownloadModule struct {
	inputStream <-chan common.JobContext
	outputStream chan common.JobContext
	errorStream chan common.JobContext
}
func NewDownloadModule(inputStream <-chan common.JobContext ,
						errorStream chan common.JobContext)(*DownloadModule , error ){
	return &DownloadModule{
		inputStream: inputStream,
		outputStream: make(chan common.JobContext),
		errorStream: errorStream,
	}, nil
}
func (downloadModule *DownloadModule)Run(){
	log.Println("DownloadModule Run ===")

	go func() {
		defer func() {
			log.Println("DownloadModule End ===")
			close(downloadModule.outputStream)
		}()
		var wg sync.WaitGroup
		wg.Add(common.WorkerSize)
		for i := 0; i < common.WorkerSize; i++ {
			go func(workerId int) {
				defer func() {
					wg.Done()
					log.Printf("Download Worker [%d] End\n",workerId)
				}()
				log.Printf("Download Worker [%d] Start\n", workerId)
				for {
					select {
					case job, more := <- downloadModule.inputStream:
						if more {
							log.Printf("Worker [%d] Downloading..........\n",workerId)
							time.Sleep(time.Duration(common.RandSec())*time.Second)
							job.DownloadSuccess = true
							log.Printf("Worker [%d] Download end.\n",workerId)
							downloadModule.outputStream <- job
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
func (downloadModule *DownloadModule)GetOutputStream()chan common.JobContext{
	return downloadModule.outputStream
}
