package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"

	yt "github.com/kkdai/youtube/v2"
)

// GetVidFromYT downloads YouTube video by video ID.
func DLVidFromYT(videoID string, dst string) error {
	var (
		client  yt.Client
		video   *yt.Video
		resp    *http.Response
		vidFile *os.File
		err     error
	)

	video, err = client.GetVideo(videoID)
	if err != nil {
		return err
	}

	resp, err = client.GetStream(video, &video.Formats[0])
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//////////////////////////

	vidFile, err = os.Create(dst)
	if err != nil {
		return err
	}
	defer vidFile.Close()

	_, err = io.Copy(vidFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//
func GenFrames(pc PlayConfig) error {
	var (
		vidPath  = pc.Src
		tmpFiles []string
		err      error
	)

	if pc.IsYouTube {
		vidPath = "./tmp-vid.mp4"

		err = DLVidFromYT(pc.Src, vidPath)
		if err != nil {
			return err
		}

		tmpFiles = append(tmpFiles, vidPath)
	}

	cmd := exec.Command("ffmpeg", "-i", vidPath, "-vf", fmt.Sprintf("fps=%v", pc.Fps), "tmp-frames/%06d.png")
	err = cmd.Run()
	if err != nil {
		return errors.New("ffmpeg returned an error")
	}

	err = CleanUpTmps(tmpFiles...)
	if err != nil {
		return err
	}

	return nil
}

//
// func

//
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
