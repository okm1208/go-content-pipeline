package module

import (
	"content-pipline-example/common"
	"log"
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
		for  {
			select {
				case job, more := <- convertWavModule.inputStream :
					if more {
						log.Println("Converting.......")
						time.Sleep(time.Duration(common.RandSec())*time.Second)
						job.ConvertWavSuccess = true
						log.Println("Convert end.")
						convertWavModule.outputStream <- job
					}else{
						return
					}
			}
		}

	}()

}

func (convertWavModule *ConvertWavModule)GetOutputStream()chan common.JobContext{
	return convertWavModule.outputStream
}



