package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"sort"
	"sync"

	"github.com/nsf/termbox-go"

	"github.com/unbewohnte/VideoToAscii/audio"
	"github.com/unbewohnte/VideoToAscii/extractor"
	"github.com/unbewohnte/VideoToAscii/jsonData"
	"github.com/unbewohnte/VideoToAscii/processor"
)

var (
	settings, firstLaunch = jsonData.GetSettings()
	WG                    sync.WaitGroup
	// how many goroutines you want to be working on processing frames
	MaxGOROUTINES uint = settings.MaxGOROUTINES
	// width and height of an ascii text file
	WIDTH  uint = settings.WIDTH
	HEIGHT uint = settings.HEIGHT
	// the amount of frames/second you extract from video with ffmpeg
	ExtractionFPS int = settings.ExtractionFPS
	// where the audio file will go when you extract it from a video
	AudioFilePath string = settings.AudioFilePath
	// shows if you want to play audio or not
	AudioPlayback bool = settings.AudioPlayback
	// path to ffmpeg
	FFMPEGbin = settings.FFMPEGbin
	// this is where the frames from a video will be
	VideoFramesOutputPath = settings.VideoFramesOutputPath
	// path to a video
	InputVideo = settings.InputVideo
	// this is where the processed ascii textfiles will be
	AsciiFilesPath = settings.AsciiFilesPath
	// character set for "asciifying" images
	asciiChars = settings.AsciiChars
)

func main() {
	// 1) extract frames && audio (optional) from a video
	// 2) process videoframes into Ascii
	// 3) it`s ready to play
	if firstLaunch {
		fmt.Println("Created settings file; closing in 10 seconds...")
		time.Sleep(time.Second * 10)
		os.Exit(0)
	}
	fmt.Print("Extraction mode (0), playback mode (1), processing images mode (2), Extract audio (3) ? (0,1,2,3) : ")
	var input string
	fmt.Scanln(&input)
	if input == "0" || input == "0 " {
		t0 := time.Now()

		fmt.Println("Extracting images...")
		gff := extractor.GetFFMPEG(FFMPEGbin)
		extractor.ExtractFrames(gff, InputVideo, VideoFramesOutputPath, ExtractionFPS)

		FinishTime := time.Now().Sub(t0)
		fmt.Printf("Done in %v", FinishTime)
		fmt.Scanln()

	} else if input == "1" || input == "1 " {
		asciiFiles, err := ioutil.ReadDir(AsciiFilesPath)
		if err != nil {
			panic(err)
		}
		var Sequence []string
		for _, file := range asciiFiles {
			if file.Name()[len(file.Name())-3:] == "txt" {
				frame, readErr := ioutil.ReadFile(filepath.Join(AsciiFilesPath, file.Name()))
				if readErr != nil {
					panic(readErr)
				}
				Sequence = append(Sequence, string(frame))
			}
		}
		err = termbox.Init()
		if err != nil {
			panic(err)
		}

		timeForEachFrame := time.Duration(time.Second / time.Duration(ExtractionFPS))

		t0 := time.Now()

		if AudioPlayback == true {
			go audio.PlayAudio(filepath.Join(AudioFilePath, "extractedAudio.mp3"))
		}

		var counter uint64 = 0
		var nextFrameTime time.Time = time.Now()
		for {
			if counter < uint64(len(Sequence)) {
				now := time.Now()
				if now.After(nextFrameTime) {
					nextFrameTime = now.Add(timeForEachFrame)
					showFrame(Sequence[counter])
					counter++
				}
			} else {
				termbox.Close()
				break
			}
		}
		fmt.Printf("Took %v", time.Now().Sub(t0))
		fmt.Scanln()

	} else if input == "2" || input == "2 " {
		t0 := time.Now()

		fmt.Println("Processing images...")

		files, err := ioutil.ReadDir(VideoFramesOutputPath)
		if err != nil {
			panic(err)
		}

		var sortedFilenames []string

		for _, f := range files {
			if f.Name()[len(f.Name())-3:] == extractor.ImageFileExtention {
				sortedFilenames = append(sortedFilenames, f.Name())
			}
		}
		sort.Strings(sortedFilenames)

		jobs := make(chan *processor.DataForAscii, len(sortedFilenames))

		for i := 0; i < int(MaxGOROUTINES); i++ {
			WG.Add(1)
			go Worker(jobs, &WG)
		}

		var counter uint64 = 0
		for {

			if counter == uint64(len(sortedFilenames)) {
				break
			}
			if len(jobs) < int(MaxGOROUTINES) {
				img, err := processor.GetImage(filepath.Join(VideoFramesOutputPath, sortedFilenames[counter]))
				if err != nil {
					panic(err)
				}

				jobs <- &processor.DataForAscii{
					Img:      img,
					Width:    WIDTH,
					Height:   HEIGHT,
					Filename: fmt.Sprintf("%010d_ascii.txt", counter),
				}
				img = nil
				counter++
			}
		}

		close(jobs)
		WG.Wait()

		FinishTime := time.Now().Sub(t0)
		fmt.Printf("Done in %v", FinishTime)
		fmt.Scanln()

	} else if input == "3" || input == "3 " {
		t0 := time.Now()

		fmt.Println("Extracting audio...")
		gff := extractor.GetFFMPEG(FFMPEGbin)
		extractor.ExtractAudio(gff, InputVideo, AudioFilePath)

		FinishTime := time.Now().Sub(t0)
		fmt.Printf("Done in %v", FinishTime)
		fmt.Scanln()

	}
}

func Worker(jobs <-chan *processor.DataForAscii, WG *sync.WaitGroup) {
	defer WG.Done()
	for data := range jobs {
		processor.ASCIIfy(asciiChars, data.Img, data.Width, data.Height, filepath.Join(AsciiFilesPath, data.Filename))
		data = nil
	}
}

func showFrame(frame string) {
	termbox.SetCursor(0, 0)
	var x, y int = 0, 0
	for _, char := range frame {
		termbox.HideCursor()
		if string(char) == "\n" {
			y++
			x = 0
		} else {
			termbox.SetCell(x, y, char, termbox.ColorWhite, termbox.ColorBlack)
			x++
		}

	}
	termbox.Flush()
}
