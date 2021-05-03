/*
Copyright © 2021 Brian Longmore brianl.ext@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	u "github.com/branogarbo/vidcli/util"
	"github.com/spf13/cobra"
)

var (
	vidSrc       string
	vidFPS       int
	isVidYT      bool
	outputMode   string
	asciiPattern string
	outputWidth  int
	duration     int
	isInverted   bool
	err          error
)

var rootCmd = &cobra.Command{
	Use:     "golcli",
	Short:   "Plays videos in the command line as ascii lol",
	Example: "vidcli do later",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vidSrc = args[0]

		vidFPS, err = cmd.Flags().GetInt("fps")
		isVidYT, err = cmd.Flags().GetBool("isYT")
		outputMode, err = cmd.Flags().GetString("mode")
		asciiPattern, err = cmd.Flags().GetString("ascii")
		outputWidth, err = cmd.Flags().GetInt("width")
		duration, err = cmd.Flags().GetInt("duration")
		isInverted, err = cmd.Flags().GetBool("invert")

		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = u.PlayFrames(u.PlayConfig{
			Src:          vidSrc,
			Fps:          vidFPS,
			IsYouTube:    isVidYT,
			OutputMode:   outputMode,
			AsciiPattern: asciiPattern,
			OutputWidth:  outputWidth,
			Duration:     duration,
			IsInverted:   isInverted,
		})

		if err != nil {
			fmt.Println(err)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVarP(&vidFPS, "fps", "r", 10, "do later")
	rootCmd.Flags().BoolVarP(&isVidYT, "isYT", "y", false, "do later")
	rootCmd.Flags().StringVarP(&outputMode, "mode", "m", "ascii", "do later")
	rootCmd.Flags().StringVarP(&asciiPattern, "ascii", "p", " .:-=+*#%@", "do later")
	rootCmd.Flags().IntVarP(&outputWidth, "width", "w", 75, "do later")
	rootCmd.Flags().IntVarP(&duration, "duration", "d", -1, "do later")
	rootCmd.Flags().BoolVarP(&isInverted, "invert", "i", false, "do later")
}
