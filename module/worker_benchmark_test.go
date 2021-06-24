package module

import (
	"content-pipline-example/common"
	"content-pipline-example/worker"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)


func BenchmarkDownloadConcurrencyWorker(b *testing.B){
	b.StopTimer()

	//downloadWorkerCount := 5
	//convertWorkerCount := 5
	extractWorkerCount := 5
	sampleData := getSampleData()
	testJobSize := len(sampleData)

	var wg sync.WaitGroup
	wg.Add(1)


	jobInitStream := make(chan common.JobContext)
	jobErrorStream := make(chan common.JobContext)

	//downloadFunc := func(jobContext common.JobContext)(common.JobContext,error){
	//
	//	downloadResultPath, err := DownloadTrackAudioFile(jobContext.TrackId)
	//	if err != nil {
	//		return jobContext, err
	//	}
	//	jobContext.DownloadSuccess = true
	//	jobContext.DownloadResultPath = downloadResultPath
	//
	//	return jobContext, nil
	//}
	//
	//downloadWorker, err := worker.NewCWorker(jobInitStream, jobErrorStream, downloadFunc, "DownloadWorker" , downloadWorkerCount)
	//if err != nil {
	//	b.Errorf("create download worker")
	//}
	//downloadWorker.Run()
	//convertFunc := func(jobContext common.JobContext) (common.JobContext ,error ){
	//	convertWavResultPath, err := ConvertWavFile(jobContext.TrackId,jobContext.DownloadResultPath,"44100")
	//	if err != nil {
	//		return jobContext, err
	//	}
	//	jobContext.ConvertWavResultPath = convertWavResultPath
	//	jobContext.ConvertWavSuccess = true
	//
	//	return jobContext , nil
	//}
	//
	//convertWorker , err := worker.NewCWorker(jobInitStream, jobErrorStream, convertFunc, "ConvertWorker",convertWorkerCount)
	//if err != nil {
	//	b.Errorf("create convert worker")
	//}
	//convertWorker.Run()
	extractFunc := func(jobContext common.JobContext) (common.JobContext, error){
		extractHashResultPath,err := ExtractorHashFileAndSave(jobContext.TrackId, jobContext.ConvertWavResultPath)
		if err != nil {
			return jobContext, err
		}
		jobContext.ExtractHashResultPath = extractHashResultPath
		jobContext.ExtractHashSuccess = true

		return jobContext, nil
	}
	extractWorker , err := worker.NewCWorker(jobInitStream, jobErrorStream, extractFunc, "ExtractWorker",extractWorkerCount)
	if err != nil {
		b.Errorf("create extract worker")
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
					atomic.AddInt32(&doneJobSize,1)
					if doneJobSize == int32(testJobSize) {
						close(jobInitStream)
					}
				}else{
					return
				}
			case job, more := <- jobErrorStream:
				if more {
					log.Printf("Worker error : [%d] , %v\n",job.TrackId, job.Error)
				}
			}
		}
	}()

	b.StartTimer()
	for i := 0 ; i < testJobSize ; i++ {
		trackIdStr := strconv.Itoa(sampleData[i])
		jobInitStream <- common.JobContext{
			TrackId: sampleData[i],
			//DownloadResultPath : "/Users/okm12/Downloads/audio/"+trackIdStr+".m4a",
			ConvertWavResultPath: "/Users/okm12/Downloads/wav/"+trackIdStr+".wav",

		}
	}

	wg.Wait()
}