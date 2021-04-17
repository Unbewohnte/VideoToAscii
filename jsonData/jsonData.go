package jsonData

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Data struct {
	MaxGOROUTINES         uint
	WIDTH                 uint
	HEIGHT                uint
	ExtractionFPS         int
	AudioFilePath         string
	AudioPlayback         bool
	FFMPEGbin             string
	VideoFramesOutputPath string
	InputVideo            string
	AsciiFilesPath        string
	AsciiChars            []string
}

const (
	jsonFilename string = "VideoToAsciiSettings.json"
)

var (
	defaults = Data{

		MaxGOROUTINES:         50,
		WIDTH:                 210,
		HEIGHT:                60,
		ExtractionFPS:         30,
		AudioFilePath:         ``,
		AudioPlayback:         false,
		FFMPEGbin:             `(must have)`,
		VideoFramesOutputPath: ``,
		InputVideo:            `(must have)`,
		AsciiFilesPath:        ``,
		AsciiChars:            []string{" ", "░", "▒", "▓", "█"},
	}
	currentDir, _ = os.Getwd()
)

func checkIfJsonExist() bool {
	files, err := ioutil.ReadDir(currentDir)

	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() == jsonFilename {
			return true
		}
	}
	return false
}

// This Json MUST be contained in the same directory with a binary
func createJson() *os.File {
	json, err := os.Create(filepath.Join(currentDir, jsonFilename))
	if err != nil {
		panic(err)
	}
	return json
}

func writeDefaults(jsonSettings *os.File) {
	marshalledDefaults, err := json.MarshalIndent(defaults, "", " ")
	if err != nil {
		panic(err)
	}
	jsonSettings.Write(marshalledDefaults)
}

func readFromJson() *Data {
	file, err := ioutil.ReadFile(filepath.Join(currentDir, jsonFilename))
	if err != nil {
		panic(err)
	}
	var data Data
	json.Unmarshal(file, &data)

	videoFrames := filepath.Join(data.VideoFramesOutputPath)
	asciiTxtPath := filepath.Join(data.AsciiFilesPath)
	audioPath := filepath.Join(data.AudioFilePath)

	if videoFrames == "" || videoFrames == "." || videoFrames == " " {
		data.VideoFramesOutputPath = currentDir
	}
	if asciiTxtPath == "" || asciiTxtPath == "." || asciiTxtPath == " " {
		data.AsciiFilesPath = currentDir
	}
	if audioPath == "" || audioPath == "." || audioPath == " " {
		data.AudioFilePath = currentDir
	}

	return &data
}

func createDirs(data *Data) {
	videoFrames := filepath.Join(data.VideoFramesOutputPath)
	asciiTxtPath := filepath.Join(data.AsciiFilesPath)
	audioPath := filepath.Join(data.AudioFilePath)
	if videoFrames == "" || videoFrames == "." || videoFrames == " " {
		videoFrames = currentDir
	}
	if asciiTxtPath == "" || asciiTxtPath == "." || asciiTxtPath == " " {
		asciiTxtPath = currentDir
	}
	if audioPath == "" || audioPath == "." || audioPath == " " {
		audioPath = currentDir
	}
	err := os.MkdirAll(videoFrames, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(asciiTxtPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(audioPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func GetSettings() (*Data, bool) {
	if checkIfJsonExist() == false {
		jsonFile := createJson()
		defer jsonFile.Close()

		writeDefaults(jsonFile)
		data := readFromJson()
		createDirs(data)

		return data, true
	}
	data := readFromJson()
	createDirs(data)

	return data, false
}
