package common

import (
	"fmt"
	"github.com/bradhe/stopwatch"
	"runtime"
	"testing"
	"time"
)

func TestStopWatch(t *testing.T){
	watch := stopwatch.Start()
	time.Sleep(1 * time.Second)

	time.Sleep(1 * time.Second)
	watch.Stop()
	fmt.Printf("Milliseconds elapsed: %v\n", watch.Seconds())

}

func TestGoMaxProcess(t *testing.T){
	fmt.Printf("%d\n",runtime.NumCPU())
}
