package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"

	ic "github.com/branogarbo/imgcli/util"
	yt "github.com/kkdai/youtube/v2"
)

// GetVidFromFile gets video by file path.
func GetVidFromFile(path string) (io.ReadCloser, error) {
	return ic.GetFileByPath(path)
}

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
func GenFramesFromVid(bc BuildConfig) error {
	var (
		cmd = exec.Command("ffmpeg", "-i", bc.Dst, "-vf", fmt.Sprintf("fps=%v", bc.Fps), "vidcli-tmp/frame%%06d.png")
		err = cmd.Run()
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

//
// func BuildFrames(bc BuildConfig) error {

// }

func CleanUpFiles(files ...string) error {
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
