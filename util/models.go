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
package util

type FrameMap map[int]string

type Frame struct {
	Num   int
	Chars string
}

// DA PLAY PLAN

// P. DL YT vid (opt.), have mp4 available

// X 1. split vid into frames as images in ./tmp-frames/
// X 2. generate ascii img from each frame image
// X 3. play converted frames in order at certain fps

type PlayConfig struct {
	Src          string
	Fps          int
	IsYouTube    bool
	OutputMode   string
	AsciiPattern string
	OutputWidth  int
	IsPrinted    bool
	IsQuiet      bool
	IsInverted   bool
}
