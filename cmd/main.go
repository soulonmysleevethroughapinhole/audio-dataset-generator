package main

import (
	"github.com/soulonmysleevethroughapinhole/audio-dataset-generator/pkg/dataset_generator"

	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/options"
	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/synthesizer"
	// "github.com\soulonmysleevethroughapinhole\audio-generator"
	// "github.com/soulonmysleevethroughapinhole/audio-generator/pkg/synthesizer"
)

var ws_options options.Option // maybe it'd be better as global variable, passing this everywhere feels dirty now

func main() {
	// TODO: load options from here (check up)
	ws_options.Load() // metaconfig path static & hardcoded (good), maybe replace with env?

	// CHANGE TO THIS !!!!!!!!!!!
	// synth includews audio processor, keyboard, and synth capabilities
	synth := synthesizer.New(ws_options)
	dataset_generator.App(synth, ws_options)
}
