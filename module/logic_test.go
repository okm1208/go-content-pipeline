package module

import (
	"github.com/mewkiz/pkg/osutil"
	"os"
	"testing"
)

func TestDownloadTrackAudioFile(t *testing.T) {
	downloadPath , err := DownloadTrackAudioFile(400000000)
	defer os.Remove(downloadPath)

	if err != nil {
		t.Errorf("download failed : %v",err)
		return
	}
	if !osutil.Exists(downloadPath) {
		t.Errorf("download file not exists")
		return
	}
}
func TestDownloadAndConvertFile(t *testing.T){
	trackId := 400000000
	downloadPath , err := DownloadTrackAudioFile(trackId)
	defer os.Remove(downloadPath)
	defer os.Remove(convertPath)
	if err != nil {
		t.Errorf("download failed : %v",err)
		return
	}
	if !osutil.Exists(downloadPath) {
		t.Errorf("download file not exists")
		return
	}

	convertPath , err := ConvertWavFile(trackId,downloadPath, "44100")

	if err != nil {
		t.Errorf("convert failed : %v", err)
		return
	}
	if !osutil.Exists(convertPath) {
		t.Errorf("convert file not exists")
		return
	}
}
func TestFullContentPipeline(t *testing.T){
	trackId := 400000000
	var downloadPath ,convertPath, extractPath string
	defer os.Remove(downloadPath)
	defer os.Remove(convertPath)
	defer os.Remove(extractPath)

	downloadPath , err := DownloadTrackAudioFile(trackId)

	if err != nil {
		t.Errorf("download failed : %v",err)
		return
	}
	if !osutil.Exists(downloadPath) {
		t.Errorf("download file not exists")
		return
	}
	convertPath , err = ConvertWavFile(trackId,downloadPath, "44100")

	if err != nil {
		t.Errorf("convert failed : %v", err)
		return
	}
	if !osutil.Exists(convertPath) {
		t.Errorf("convert file not exists")
		return
	}

	extractPath, err = ExtractorHashFileAndSave(trackId, convertPath)
	if err != nil {
		t.Errorf("extract failed : %v", err)
		return
	}
	if !osutil.Exists(extractPath) {
		t.Errorf("extract file not exists")
		return
	}

}