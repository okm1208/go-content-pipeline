package worker

import (
	"content-pipline-example/common"
	"log"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrencyWorkerSimple(t *testing.T) {
	testJobSize := 10

	var wg sync.WaitGroup
	wg.Add(1)

	jobInitStream := make(chan common.JobContext)
	jobErrorStream := make(chan common.JobContext)

	// download 비즈니스 로직
	downloadFunc := func(jobContext common.JobContext) (common.JobContext ,error ) {
		jobContext.DownloadSuccess = true
		return jobContext, nil
	}
	downloadWorker, err := NewCWorker(jobInitStream, jobErrorStream, downloadFunc, "DownloadWorker" , 5)
	if err != nil {
		t.Errorf("create download worker")
	}
	downloadWorker.Run()

	// convert 비즈니스 로직
	convertFunc := func(jobContext common.JobContext) (common.JobContext,error ) {
		jobContext.ConvertWavSuccess = true
		return jobContext, nil
	}

	convertWorker, err := NewCWorker(downloadWorker.GetOutputStream(), jobErrorStream, convertFunc, "ConvertWorker",5)
	if err != nil {
		t.Errorf("create convert worker")
	}
	convertWorker.Run()

	// extract 비즈니스 로직
	extractFunc := func(jobContext common.JobContext) (common.JobContext,error) {
		jobContext.ExtractHashSuccess = true
		return jobContext, nil
	}
	extractWorker, err := NewCWorker(convertWorker.GetOutputStream(), jobErrorStream, extractFunc, "ExtractWorker", 5)
	if err != nil {
		t.Errorf("create extract worker")
	}
	extractWorker.Run()
	jobOutputStream := extractWorker.GetOutputStream()

	go func() {
		var doneJobSize int32 = 0
		defer func() {
			log.Println("All Job InputStream Close.")
			wg.Done()
		}()
		for {
			select {
			case job, more := <-jobOutputStream:
				if more {
					log.Printf("Worker done : [%d]\n", job.TrackId)
					if !job.DownloadSuccess || !job.ConvertWavSuccess || !job.ExtractHashSuccess {
						t.Errorf("download : %t, convert : %t, extract : %t\n",
							job.DownloadSuccess, job.ConvertWavSuccess, job.ExtractHashSuccess)
					}
					atomic.AddInt32(&doneJobSize, 1)
					if doneJobSize == int32(testJobSize) {
						return
					}
				} else {
					return
				}
			}
		}
	}()

	for i := 0; i < testJobSize; i++ {
		jobInitStream <- common.JobContext{TrackId: i}
	}

	wg.Wait()
}
