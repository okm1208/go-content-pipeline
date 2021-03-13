package module

import (
	"content-pipline-example/common"
	"log"
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
		for {
			select{
			case job, more := <- extractHash.inputStream:
				if more {
					log.Println("Extracting.......")
					time.Sleep(time.Duration(common.RandSec())*time.Second)
					job.ExtractHashSuccess = true
					log.Println("Extract End.")
					extractHash.outputStream <- job
				}else{
					return
				}
			}
		}
	}()

}


func (extractHash *ExtractHash)GetOutputStream()chan common.JobContext{
	return extractHash.outputStream
}