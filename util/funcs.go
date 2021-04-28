package util

import (
	"fmt"
	"os"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/cgo/ffmpeg"
)

func GetFramesFromVid(filePath string) error {
	var (
		file     []byte
		err      error
		decoder  *ffmpeg.VideoDecoder
		stream   av.CodecData
		vidFrame *ffmpeg.VideoFrame
	)

	file, err = os.ReadFile(filePath)
	if err != nil {
		return err
	}

	decoder, err = ffmpeg.NewVideoDecoder(stream)
	if err != nil {
		return err
	}

	vidFrame, err = decoder.Decode(file)
	if err != nil {
		return err
	}

	fmt.Println(vidFrame.Image)

	return nil
}
