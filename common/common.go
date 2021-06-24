package common

import (
	"math/rand"
)
import "time"


//func DownloadTrackAudioFile(trackId int)(string, error){
//func ConvertWavFile(trackId int, orgFilePath string, sampleRate string)(string,error){
//func ExtractorHashFileAndSave (trackId int,orgFilePath string)(string,error){

type JobContext struct {
	TrackId int

	DownloadResultPath string
	ConvertWavResultPath string
	ExtractHashResultPath string

	DownloadSuccess bool
	ConvertWavSuccess bool
	ExtractHashSuccess bool

	Error error
}
var WorkerSize = 2
func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

type QueueMessage struct {
	TrackId int `json:"trackId"`
}

func RandSec() int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	return n
}