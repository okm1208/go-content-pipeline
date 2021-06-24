package module

import (
	"os/exec"
	"path/filepath"
	"strconv"
)

var convertPath = "/Users/okm12/Downloads/wav"

func ConvertWavFile(trackId int, orgFilePath string, sampleRate string)(string,error){
	wavFileName := strconv.Itoa(trackId) + ".wav"
	convertLocalPath := convertPath + string(filepath.Separator) + wavFileName

	cmd := exec.Command("ffmpeg",
		"-i",
		orgFilePath,
		"-sample_fmt",
		"s16",
		"-ar" ,
		sampleRate,
		"-y",
		convertLocalPath)
	err := cmd.Run()
	if err != nil {
		return "", err
	} else {
		return convertLocalPath, nil
	}
}


