package extractor

import (
	"fmt"
	"path/filepath"

	"github.com/lijo-jose/gffmpeg/pkg/gffmpeg"
)

const (
	ImageFileExtention string = "png"
)

func GetFFMPEG(binPath string) gffmpeg.GFFmpeg {
	GFFMPEG, err := gffmpeg.NewGFFmpeg(binPath)
	if err != nil {
		panic(err)
	}
	return GFFMPEG
}

func ExtractFrames(gff gffmpeg.GFFmpeg, videoPath, outputPath string, fps int) {
	builder := gffmpeg.NewBuilder()
	builder = builder.SrcPath(videoPath).VideoFilters(fmt.Sprintf("fps=%v", fps)).DestPath(outputPath + "%10d." + ImageFileExtention)
	gff.Set(builder)
	output := gff.Start(nil)
	if output.Err != nil {
		panic(output.Err)
	}

}

func ExtractAudio(gff gffmpeg.GFFmpeg, videoPath, outputPath string) {
	builder := gffmpeg.NewBuilder()
	builder = builder.SrcPath(videoPath).VideoFilters("-q:a -map a").DestPath(filepath.Join(outputPath, "extractedAudio.mp3"))
	gff.Set(builder)
	output := gff.Start(nil)
	if output.Err != nil {
		panic(output.Err)
	}
}
