/*
Copyright © 2021 Brian Longmore branodev@gmail.com

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
	ic "github.com/branogarbo/imgcli/util"
	u "github.com/branogarbo/vidcli/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "vidcli",
	Version: "v1.2.0",
	Short:   "Plays videos in the command line as ascii lol",
	Example: "vidcli do later",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vidSrc := args[0]

		vidFPS, _ := cmd.Flags().GetInt("fps")
		isVidYT, _ := cmd.Flags().GetBool("isYT")
		outputMode, _ := cmd.Flags().GetString("mode")
		asciiPattern, _ := cmd.Flags().GetString("ascii")
		outputWidth, _ := cmd.Flags().GetInt("width")
		duration, _ := cmd.Flags().GetInt("duration")
		isInverted, _ := cmd.Flags().GetBool("invert")

		_, err := u.PlayFrames(u.PlayConfig{
			Src:          vidSrc,
			Fps:          vidFPS,
			IsYouTube:    isVidYT,
			OutputMode:   outputMode,
			AsciiPattern: asciiPattern,
			OutputWidth:  outputWidth,
			Duration:     duration,
			IsInverted:   isInverted,
		})

		return err
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntP("fps", "r", 10, "do later")
	rootCmd.Flags().BoolP("isYT", "y", false, "do later")
	rootCmd.Flags().StringP("mode", "m", ic.DefaultMode, "do later")
	rootCmd.Flags().StringP("ascii", "p", ic.DefaultPattern, "do later")
	rootCmd.Flags().IntP("width", "w", 75, "do later")
	rootCmd.Flags().IntP("duration", "d", -1, "do later")
	rootCmd.Flags().BoolP("invert", "i", false, "do later")
}
