package module

import (
	"content-pipline-example/common"
	"log"
	"time"
)


type DownloadModule struct {
	inputStream <-chan common.JobContext
	outputStream chan common.JobContext
	errorStream chan common.JobContext
	doneStream chan bool
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
		for {
			select{
				case job, more := <- downloadModule.inputStream:
					if more {
						log.Println("Downloading.......")
						time.Sleep(time.Duration(common.RandSec())*time.Second)
						job.DownloadSuccess = true
						log.Println("Download end.")
						downloadModule.outputStream <- job
					} else {
						return
					}
			}
		}
	}()
}
func (downloadModule *DownloadModule)GetOutputStream()chan common.JobContext{
	return downloadModule.outputStream
}
