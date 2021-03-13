package main

import (
	"content-pipline-example/module"
	"content-pipline-example/common"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)


var mqUrl = "amqp://admin:admin@localhost//mcp"

func main(){
	log.Println("Main Process Start ===")

	var wg sync.WaitGroup
	wg.Add(1)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT)

	// MQ Connection
	var mqConn common.MqConn
	err := mqConn.InitConn(mqUrl)
	if err != nil {
		panic(err)
	}

	jobInitStream := make(chan common.JobContext)
	jobErrorStream := make(chan common.JobContext)

	jobOutputStream := CreateDefaultPipLine(jobInitStream,jobErrorStream)

	go func(){
		defer func() {
			log.Println("All Job InputStream Close.")
			wg.Done()
		}()
		for{
			select {
			case job, more := <-jobOutputStream:
				if more {
					log.Printf("job success : download -> %t, convert -> %t , hash extract -> %t\n",
						job.DownloadSuccess,job.ConvertWavSuccess, job.ExtractHashSuccess)
				}else{
					return
				}
			case job ,more := <-jobErrorStream:
				if more {
					log.Printf("job error : download -> %t, convert -> %t, hash extract -> %t\n",
						job.DownloadSuccess, job.ConvertWavSuccess, job.ExtractHashSuccess)
				}else{
					return
				}
			}

		}
	}()

	mqDeliveries , err := mqConn.Channel.Consume("WINTER.Q.DI_TO_TC",
		"content-pipeline",
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		panic(err)
	}
	module.MsgDelivery(mqDeliveries, jobInitStream)

	go func(){
		signal := <-osSignal
			log.Printf("Receive Os signal : %s\n",signal.String())
			mqConn.CloseConn()
	}()

	wg.Wait()
	log.Println("Main Process End ===")
}



func CreateDownloadPipLine(jobInitStream, jobErrorStream chan common.JobContext) chan common.JobContext{
	downloadModule,err := module.NewDownloadModule(jobInitStream,jobErrorStream)
	if err != nil {
		panic(err)
	}
	downloadModule.Run()
	return downloadModule.GetOutputStream()
}
func CreateDefaultPipLine(jobInitStream, jobErrorStream chan common.JobContext)chan common.JobContext{
	downloadModule,err := module.NewDownloadModule(jobInitStream,jobErrorStream)
	if err != nil {
		panic(err)
	}
	downloadModule.Run()

	convertWavModule, err := module.NewConvertWavModule(downloadModule.GetOutputStream(),jobErrorStream)
	if err != nil {
		panic(err)
	}
	convertWavModule.Run()

	extractHashModule, err := module.NewExtractHash(convertWavModule.GetOutputStream(),jobErrorStream)
	if err != nil {
		panic(err)
	}
	extractHashModule.Run()

	return extractHashModule.GetOutputStream()
}