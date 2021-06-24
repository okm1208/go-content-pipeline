package module

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var downloadPath = "/Users/okm12/Downloads/audio"

type ModFile struct {
	MultimediaSvcType string `json:"multimediaSvcType"`
	MultimediaFileType string `json:"multimediaFileType"`
	HighlightYn string `json:"highlightYn"`
	BitRate int `json:"bitRate"`
	ModStreamPath string `json:"modStreamPath"`
}

func DownloadTrackAudioFile(trackId int)(string, error){
	modFiles , err := getModFiles(trackId)
	if err != nil {
		return "", err
	}
	var targetModFile *ModFile
	for _ , m := range modFiles {
		if m.MultimediaSvcType == "D" &&
						m.MultimediaFileType == "AAC" &&
							m.BitRate == 128 &&
								m.HighlightYn == "N" {
			targetModFile = &m
			break
		}
	}
	if targetModFile == nil {
		return "", errors.New("ModFile is empty.")
	}
	downloadUrl := genServiceUrlTimestamp(targetModFile.ModStreamPath,genTimeStamp())
	downloadFileName := strconv.Itoa(trackId)+filepath.Ext(downloadUrl)
	downloadFilePath :=
		downloadPath +string(filepath.Separator)+ downloadFileName

	err = downloadFile(downloadUrl, downloadFilePath)

	if err != nil {
		return "", err
	}

	return downloadFilePath, nil
}

func downloadFile(downloadUrl string, downloadPath string) error {
	// Get the data
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return  err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("download statusCode : %d", resp.StatusCode))
	}
	// Create the file
	out, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return err
}
func genTimeStamp() string {
	t := time.Now().Add(time.Hour * 24 *7)
	return strconv.FormatInt( t.Unix() , 10)
}

func genServiceUrlTimestamp(TargetUrl string , timestamp string) string {

	SecretKey := "fLo_AutH"
	Domain := "https://hpd.music-flo.com"

	hashString := fmt.Sprintf("%s%s%s*" , SecretKey , TargetUrl , timestamp)

	h := sha1.New()
	io.WriteString(h, hashString)

	hashKey := fmt.Sprintf("%x", h.Sum(nil))

	splitPath := strings.Split(TargetUrl , "/")
	l := len(splitPath)
	mediaUrl :=	strings.Join(splitPath[2:l] , "/")

	randString := genRandomString(16)

	return fmt.Sprintf("%s/%s/%s/%s/%s/%s" , Domain , splitPath[1], hashKey , timestamp,randString ,mediaUrl)
}

func genRandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func getModFiles(trackId int) ([]ModFile , error){
	url := fmt.Sprintf("http://pri-heartqueen.music-flo.com/track/modfiles?trackId=%d" , trackId)
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		var body []byte
		var modFile []ModFile
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {

			err = json.Unmarshal(body , &modFile)
			if err == nil {
				return modFile , nil
			}
		}
	}
	return nil , err
}