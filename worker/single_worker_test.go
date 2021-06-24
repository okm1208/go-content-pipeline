package worker

import (
	"content-pipline-example/common"
	"log"
	"sync"
	"sync/atomic"
	"testing"
)

func TestWorkerMultiple(t *testing.T){
	testJobSize := 5

	var wg sync.WaitGroup
	wg.Add(1)

	jobInitStream := make(chan common.JobContext)
	jobErrorStream := make(chan common.JobContext)

	downloadFunc := func(jobContext common.JobContext) common.JobContext{
		jobContext.DownloadSuccess = true
		return jobContext
	}
	downloadWorker , err := NewWorker(jobInitStream, jobErrorStream, downloadFunc, "DownloadWorker")
	if err != nil {
		t.Errorf("create download worker")
	}
	downloadWorker.Run()

	convertFunc := func(jobContext common.JobContext) common.JobContext{
		jobContext.ConvertWavSuccess = true
		return jobContext
	}

	convertWorker , err := NewWorker(downloadWorker.GetOutputStream(), jobErrorStream, convertFunc, "ConvertWorker")
	if err != nil {
		t.Errorf("create convert worker")
	}
	convertWorker.Run()

	extractFunc := func(jobContext common.JobContext) common.JobContext{
		jobContext.ExtractHashSuccess = true
		jobContext.WorkerSuccess = true
		return jobContext
	}
	extractWorker , err := NewWorker(convertWorker.GetOutputStream(), jobErrorStream, extractFunc, "ExtractWorker")
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
					log.Printf("Worker done : [%d]\n",job.TrackId)
					if !job.DownloadSuccess || !job.ConvertWavSuccess || !job.ExtractHashSuccess || !job.WorkerSuccess {
						t.Errorf("download : %t, convert : %t, extract : %t, work success : %t\n",
							job.DownloadSuccess,job.ConvertWavSuccess,job.ExtractHashSuccess, job.WorkerSuccess)
					}
					atomic.AddInt32(&doneJobSize,1)
					if doneJobSize == int32(testJobSize) {
						return
					}
				}else{
					return
				}
			}
		}
	}()

	for i := 0 ; i < testJobSize ; i++ {
		jobInitStream <- common.JobContext{TrackId: i}
	}

	wg.Wait()
}
func TestWorkerSimple(t *testing.T){
	testJobSize := 10

	var wg sync.WaitGroup
	wg.Add(1)

	jobInitStream := make(chan common.JobContext)
	jobErrorStream := make(chan common.JobContext)

	downloadFunc := func(jobContext common.JobContext)common.JobContext{
		jobContext.DownloadSuccess = true
		jobContext.WorkerSuccess = true
		return jobContext
	}
	downloadWorker, err := NewWorker(jobInitStream, jobErrorStream ,downloadFunc,"DownloadWorker")
	if err != nil {
		t.Errorf("create download worker")
	}
	jobOutputStream := downloadWorker.GetOutputStream()


	downloadWorker.Run()
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
					log.Printf("Worker done : [%d]\n",job.TrackId)
					if !job.WorkerSuccess {
						t.Errorf("WorkerSuccess %t; want true\n", job.WorkerSuccess)
					}
					atomic.AddInt32(&doneJobSize,1)
					if doneJobSize == int32(testJobSize) {
						return
					}
				}else{
					return
				}
			}
		}
	}()

	for i := 0 ; i < testJobSize ; i++ {
		jobInitStream <- common.JobContext{TrackId: i}
	}


	wg.Wait()

}

