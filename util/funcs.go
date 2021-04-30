package util

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	ic "github.com/branogarbo/imgcli/util"
	gt "github.com/buger/goterm"
	yt "github.com/kkdai/youtube/v2"
)

// GetVidFromYT downloads YouTube video by video ID.
// func DLVidFromYT(videoID string, dst string) error {
// 	var (
// 		client  yt.Client
// 		video   *yt.Video
// 		resp    *http.Response
// 		vidFile *os.File
// 		err     error
// 	)

// 	video, err = client.GetVideo(videoID)
// 	if err != nil {
// 		return err
// 	}

// 	resp, err = client.GetStream(video, &video.Formats[0])
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	//////////////////////////

// 	vidFile, err = os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer vidFile.Close()

// 	_, err = io.Copy(vidFile, resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

//
func GetVidBytesYT(videoID string) ([]byte, error) {
	var (
		client   yt.Client
		video    *yt.Video
		resp     *http.Response
		vidBytes []byte
		err      error
	)

	video, err = client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	resp, err = client.GetStream(video, &video.Formats[0])
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	vidBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return vidBytes, nil
}

//
func GetVidBytesFile(filePath string) ([]byte, error) {
	var (
		vidBytes []byte
		err      error
	)

	vidBytes, err = os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return vidBytes, nil
}

//
func GenFrameImages(pc PlayConfig) error {
	var (
		vidBytes []byte
		err      error
	)

	if pc.IsYouTube {
		vidBytes, err = GetVidBytesYT(pc.Src)
	} else {
		vidBytes, err = GetVidBytesFile(pc.Src)
	}
	if err != nil {
		return err
	}

	err = os.Mkdir("./tmp-frames", 0777)
	if err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg", "-i", "-", "-vf", fmt.Sprintf("fps=%v", pc.Fps), "./tmp-frames/out%d.png")
	cmd.Stdin = bytes.NewBuffer(vidBytes)

	err = cmd.Run()
	if err != nil {
		return errors.New("ffmpeg returned error")
	}

	return nil
}

//
func ConvertFrames(pc PlayConfig) (FrameMap, error) {
	var (
		wg         sync.WaitGroup
		frameFiles []os.DirEntry
		frameChars string
		err        error
		frames     = make(FrameMap)
	)

	frameFiles, err = os.ReadDir("./tmp-frames")
	if err != nil {
		return nil, err
	}

	wg.Add(len(frameFiles))

	errChan := make(chan error, len(frameFiles))
	frameChan := make(chan Frame, len(frameFiles))

	for i, frameFile := range frameFiles {
		go func(i int, ffName string) {
			frameChars, err = ic.OutputImage(ic.OutputConfig{
				Src:          "./tmp-frames/" + ffName,
				OutputMode:   pc.OutputMode,
				AsciiPattern: pc.AsciiPattern,
				OutputWidth:  pc.OutputWidth,
				IsInverted:   pc.IsInverted,
			})
			if err != nil {
				errChan <- err
			}

			frame := Frame{i, frameChars}

			frameChan <- frame
			wg.Done()
		}(i, frameFile.Name())
	}

	wg.Wait()
	close(errChan)

	if err = <-errChan; err != nil {
		return nil, err
	}

	for i := 0; i < len(frameFiles); i++ {
		frame := <-frameChan

		frames[frame.Num] = frame.Chars
	}

	err = CleanUpTmps("./tmp-frames")
	if err != nil {
		return nil, err
	}

	return frames, nil
}

func PlayFrames(pc PlayConfig) (FrameMap, error) {
	var (
		frames FrameMap
		err    error
	)

	err = GenFrameImages(pc)
	if err != nil {
		return nil, err
	}

	frames, err = ConvertFrames(pc)
	if err != nil {
		return nil, err
	}

	gt.Clear()

	for _, fChars := range frames {
		gt.MoveCursor(1, 1)
		gt.Print(fChars)
		gt.Flush()

		time.Sleep(time.Duration(1/pc.Fps) * time.Second)
	}

	return frames, err
}

func CleanUpTmps(files ...string) error {
	var (
		wg      sync.WaitGroup
		err     error
		errChan = make(chan error, len(files))
	)

	wg.Add(len(files))

	for _, file := range files {
		go func(file string) {
			err = os.RemoveAll(file)
			if err != nil {
				errChan <- err
			}

			wg.Done()
		}(file)
	}

	wg.Wait()
	close(errChan)

	return <-errChan
}
