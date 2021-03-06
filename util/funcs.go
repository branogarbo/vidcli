package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	ic "github.com/branogarbo/imgcli/util"
	gt "github.com/buger/goterm"
	yt "github.com/kkdai/youtube/v2"
	gb "github.com/thecodeteam/goodbye"
)

func PlayFrames(pc PlayConfig) (FrameMap, error) {
	var (
		frames FrameMap
		err    error
	)

	gb.Register(func(ctx context.Context, s os.Signal) {
		err = cleanUpTmps(pc.TmpDirName)
		if err != nil {
			fmt.Println(err)
		}
	})

	_, err = exec.LookPath("ffmpeg")
	if err != nil {
		return nil, errors.New("ffmpeg is not installed")
	}

	err = genFrameImages(&pc)
	if err != nil {
		return nil, err
	}

	frames, err = convertFrames(pc)
	if err != nil {
		return nil, err
	}

	gt.Clear()

	for i := 1; i < len(frames)+1; i++ {
		gt.MoveCursor(1, 1)
		gt.Print(string(frames[i]))
		gt.Flush()

		time.Sleep(time.Second / time.Duration(pc.Fps))
	}

	return frames, err
}

func getVidBytesYT(videoID string) ([]byte, error) {
	var (
		client   yt.Client
		video    *yt.Video
		resp     io.ReadCloser
		vidBytes []byte
		err      error
	)

	video, err = client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	resp, _, err = client.GetStream(video, &video.Formats[0])
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	vidBytes, err = ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	return vidBytes, nil
}

func getVidBytesFile(filePath string) ([]byte, error) {
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

func genFrameImages(pc *PlayConfig) error {
	var (
		vidBytes   []byte
		tmpDirName string
		argList    []string
		cmd        *exec.Cmd
		err        error
	)

	fmt.Println("Loading video...")

	if pc.IsYouTube {
		vidBytes, err = getVidBytesYT(pc.Src)
	} else {
		vidBytes, err = getVidBytesFile(pc.Src)
	}
	if err != nil {
		return err
	}

	tmpDirName, err = ioutil.TempDir(".", "frames")
	if err != nil {
		return err
	}

	fmt.Println("Extracting frames...")

	if pc.Duration == -1 {
		argList = strings.Split(fmt.Sprintf("-i - -vf fps=%v ./%v/%%d.png", pc.Fps, tmpDirName), " ")
	} else {
		argList = strings.Split(fmt.Sprintf("-i - -vf fps=%v -t %v ./%v/%%d.png", pc.Fps, pc.Duration, tmpDirName), " ")
	}

	cmd = exec.Command("ffmpeg", argList...)
	cmd.Stdin = bytes.NewBuffer(vidBytes)

	err = cmd.Run()
	if err != nil {
		return errors.New("ffmpeg returned error")
	}

	pc.TmpDirName = tmpDirName

	return nil
}

func convertFrames(pc PlayConfig) (FrameMap, error) {
	var (
		wg         sync.WaitGroup
		frameFiles []os.DirEntry
		frameChars string
		err        error
		frames     = make(FrameMap)
	)

	if pc.OutputMode == "color" {
		return nil, errors.New("vidcli does not support color")
	}

	frameFiles, err = os.ReadDir(pc.TmpDirName)
	if err != nil {
		return nil, err
	}

	wg.Add(len(frameFiles))

	errChan := make(chan error, len(frameFiles))
	frameChan := make(chan Frame, len(frameFiles))

	fmt.Println("Converting frames...")

	for _, frameFile := range frameFiles {
		go func(ffName string) {
			frameNum, err := strconv.Atoi(ffName[:len(ffName)-4])
			if err != nil {
				errChan <- err
				frameChan <- Frame{}
				return
			}

			frameChars, err = ic.OutputImage(ic.OutputConfig{
				Src:          fmt.Sprintf("./%v/%v", pc.TmpDirName, ffName),
				OutputMode:   pc.OutputMode,
				AsciiPattern: pc.AsciiPattern,
				OutputWidth:  pc.OutputWidth,
				IsInverted:   pc.IsInverted,
			})
			errChan <- err

			frame := Frame{frameNum, []byte(frameChars)}
			frameChan <- frame

			wg.Done()
		}(frameFile.Name())
	}

	wg.Wait()
	close(frameChan)
	close(errChan)

	if err = <-errChan; err != nil {
		return nil, err
	}

	for i := 0; i < len(frameFiles); i++ {
		frame := <-frameChan

		frames[frame.Num] = frame.Chars
	}

	return frames, nil
}

func cleanUpTmps(files ...string) error {
	var (
		wg      sync.WaitGroup
		err     error
		errChan = make(chan error, len(files))
	)

	wg.Add(len(files))

	for _, file := range files {
		go func(file string) {
			err = os.RemoveAll(file)
			errChan <- err

			wg.Done()
		}(file)
	}

	wg.Wait()
	close(errChan)

	return <-errChan
}
