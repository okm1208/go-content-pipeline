package module

import (
	"content-pipline-example/common"
	"content-pipline-example/worker"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)
func getSampleData()[]int{
	csvFile, err := os.Open("/Users/okm12/go/src/content-pipline-example/sample/track_list")
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	result := make([]int,0)
	for _, record := range records {
		trackId, _ :=strconv.Atoi(record[0])
		result = append(result , trackId)
	}

	return result
}
/**

    30개의 고루틴을 기준
	download : 10,  convert : 10 , extract : 10 -> 73 초
	download : 5 , convert : 5 , extract : 20 -> 71 초
	download : 10, convert : 5 , extract : 15 -> 78초

	40개의 고루틴 기준

	download : 10, convert : 10 , extract : 20 -> 69초
	download : 20, convert : 20 , extract : 20 -> 69초
	dpwnload : 30 , convert : 30, extract : 30 -> 82초
	// download 서비스 : 60초
	// convert wav 서비스 : 6초
	// extract 서비스 : 92초
 */

func BenchmarkConcurrencyWorker(b *testing.B){

	var downloadWorkerCount , convertWorkerCount, extractWorkerCount int
	downloadWorkerCount = 10
	convertWorkerCount = 10
	extractWorkerCount = 10
	b.StopTimer()

	sampleData := getSampleData()
	testJobSize := len(sampleData)

	var wg sync.WaitGroup
	wg.Add(1)

	jobInitStream := make(chan common.JobContext)
	jobErrorStream := make(chan common.JobContext)

	downloadFunc := func(jobContext common.JobContext)(common.JobContext,error){

		downloadResultPath, err := DownloadTrackAudioFile(jobContext.TrackId)
		if err != nil {
			return jobContext, err
		}
		jobContext.DownloadSuccess = true
		jobContext.DownloadResultPath = downloadResultPath

		return jobContext, nil
	}

	downloadWorker, err := worker.NewCWorker(jobInitStream, jobErrorStream, downloadFunc, "DownloadWorker" , downloadWorkerCount)
	if err != nil {
		b.Errorf("create download worker")
	}
	downloadWorker.Run()



	convertFunc := func(jobContext common.JobContext) (common.JobContext ,error ){
		convertWavResultPath, err := ConvertWavFile(jobContext.TrackId,jobContext.DownloadResultPath,"44100")
		if err != nil {
			return jobContext, err
		}
		jobContext.ConvertWavResultPath = convertWavResultPath
		jobContext.ConvertWavSuccess = true

		return jobContext , nil
	}

	convertWorker , err := worker.NewCWorker(downloadWorker.GetOutputStream(), jobErrorStream, convertFunc, "ConvertWorker",convertWorkerCount)
	if err != nil {
		b.Errorf("create convert worker")
	}
	convertWorker.Run()

	extractFunc := func(jobContext common.JobContext) (common.JobContext, error){
		extractHashResultPath,err := ExtractorHashFileAndSave(jobContext.TrackId, jobContext.ConvertWavResultPath)
		if err != nil {
			return jobContext, err
		}
		jobContext.ExtractHashResultPath = extractHashResultPath
		jobContext.ExtractHashSuccess = true

		return jobContext, nil
	}
	extractWorker , err := worker.NewCWorker(convertWorker.GetOutputStream(), jobErrorStream, extractFunc, "ExtractWorker",extractWorkerCount)
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
					if !job.DownloadSuccess || !job.ConvertWavSuccess || !job.ExtractHashSuccess {
						b.Errorf("download : %t, convert : %t, extract : %t\n",
							job.DownloadSuccess,job.ConvertWavSuccess,job.ExtractHashSuccess)
					}
					atomic.AddInt32(&doneJobSize,1)
					if doneJobSize == int32(testJobSize) {
						//return
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
		jobInitStream <- common.JobContext{TrackId: sampleData[i]}
	}
	wg.Wait()
}

/**
	100 개 음원 파일을 처리 할때의 속도 측정 ( single worker ) : 936 초
 */
func BenchmarkSingleWorker(b *testing.B){

	b.StopTimer()
	sampleData := getSampleData()
	testJobSize := len(sampleData)

	var wg sync.WaitGroup
	wg.Add(1)

	jobInitStream := make(chan common.JobContext)
	jobErrorStream := make(chan common.JobContext)

	downloadFunc := func(jobContext common.JobContext)(common.JobContext,error){

		downloadResultPath, err := DownloadTrackAudioFile(jobContext.TrackId)
		if err != nil {
			return jobContext, err
		}
		jobContext.DownloadSuccess = true
		jobContext.DownloadResultPath = downloadResultPath

		return jobContext, nil
	}
	downloadWorker , err := worker.NewWorker(jobInitStream, jobErrorStream, downloadFunc, "DownloadWorker")
	if err != nil {
		b.Errorf("create download worker")
	}
	downloadWorker.Run()

	convertFunc := func(jobContext common.JobContext) (common.JobContext ,error ){
		convertWavResultPath, err := ConvertWavFile(jobContext.TrackId,jobContext.DownloadResultPath,"44100")
		if err != nil {
			return jobContext, err
		}
		jobContext.ConvertWavResultPath = convertWavResultPath
		jobContext.ConvertWavSuccess = true

		return jobContext , nil
	}

	convertWorker , err := worker.NewWorker(downloadWorker.GetOutputStream(), jobErrorStream, convertFunc, "ConvertWorker")
	if err != nil {
		b.Errorf("create convert worker")
	}
	convertWorker.Run()

	extractFunc := func(jobContext common.JobContext) (common.JobContext, error){
		extractHashResultPath,err := ExtractorHashFileAndSave(jobContext.TrackId, jobContext.ConvertWavResultPath)
		if err != nil {
			return jobContext, err
		}
		jobContext.ExtractHashResultPath = extractHashResultPath
		jobContext.ExtractHashSuccess = true

		return jobContext, nil
	}
	extractWorker , err := worker.NewWorker(convertWorker.GetOutputStream(), jobErrorStream, extractFunc, "ExtractWorker")
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
					if !job.DownloadSuccess || !job.ConvertWavSuccess || !job.ExtractHashSuccess {
						b.Errorf("download : %t, convert : %t, extract : %t\n",
							job.DownloadSuccess,job.ConvertWavSuccess,job.ExtractHashSuccess)
					}
					atomic.AddInt32(&doneJobSize,1)
					if doneJobSize == int32(testJobSize) {
						return
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
		jobInitStream <- common.JobContext{TrackId: sampleData[i]}
	}
	wg.Wait()

}
