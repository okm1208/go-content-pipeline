package common

import (
	"math/rand"
)
import "time"

type JobContext struct {
	TrackId int
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