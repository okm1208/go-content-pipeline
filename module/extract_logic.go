package module

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var extractPath = "/Users/okm12/Downloads/npy"
var hashExtractorPath = "/Users/okm12/go/src/content-pipline-example/python_lib/hash_extractor/hash_extractor.py"
func ExtractorHashFileAndSave (trackId int,orgFilePath string)(string,error){

	extractorLocalPath := extractPath + string(filepath.Separator) + strconv.Itoa(trackId) +".npy"
	trackIdParam := strconv.Itoa(trackId)

	cmd := exec.Command("python" ,
		hashExtractorPath,
		"--filename="+orgFilePath ,
		"--songid="+trackIdParam,
		"--output="+extractorLocalPath)

	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	err := cmd.Run()
	errScan := bufio.NewScanner(bytes.NewReader(stdErr.Bytes()))
	outScan := bufio.NewScanner(bytes.NewReader(stdOut.Bytes()))

	if err != nil {
		errorStr := make([]string,5)

		for errScan.Scan() {
			errorStr = append(errorStr, errScan.Text())
		}
	}else{
		for outScan.Scan() {
			outText := outScan.Text()

			if strings.HasPrefix(outText,"success") {
				npyResultPath := strings.Replace(outText , "success:","",1)
				if len(npyResultPath) > 0 {
					return npyResultPath, nil
				}
				err = errors.New("ExtractorHashFile npyResultPath Empty")
			}else{
				err = errors.New("ExtractorHashFile Return Failed")
			}
		}
	}
	return "",err
}

