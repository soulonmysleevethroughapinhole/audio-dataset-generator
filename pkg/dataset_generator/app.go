package dataset_generator

import (
	"log"

	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/audio_handler"
	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/emulator"
	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/options"
	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/patch_generator"
	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/patch_handler"
	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/synthesizer/audiomodule"
)

/*
called by main, passed options
instantiates a synthesizer -- ALTERATION, done in main.go
gets synthesizer capabilities
passes opt. & synth cap. to patch_generator
loops through patches returned, passes them to synthesizer, writes to dataset
dataset name / folder generated from options & synth capabilities

one function: generated datasets
*/

var audiohandler audio_handler.AudioHandler

func App(audioProcessor audiomodule.AudioProcessor, synth_capabilities map[string]map[int]map[string]emulator.AccessParam, ws_options options.Option) {
	// get synth. capabilities means accessparameters, basically, i think
	audiohandler.Load(ws_options)

	patches, patch_access := patch_generator.App(synth_capabilities, ws_options)

	var patchhandler patch_handler.PatchHandler
	patchhandler.Load(ws_options)

	// maybe can remove soon?
	// if len(patches) != 6561 {
	// if len(patches) != 59049 {
	if len(patches) != 2268 { // bad
		log.Println(len(patches))
		log.Println(synth_capabilities)
		log.Fatal("maybe the order of ranges is different? next time this comes up save to debug")
	}

	path_list, filename_list := patchhandler.GetFilesInPersistence()

	if len(filename_list) == 0 { // TEMP FIX
		switch ws_options.Preference.PatchPersistenceFolderOptions {
		case "subfolders1000":
			// nuke folder
			patchhandler.PersistPatches(patches) // this is done in persistpatch, so fix it up
			// patchhandler.PersistPatchesSubfolders(patches, 1000)
		case "single_folder":
			log.Fatal("not implemented, just use the fn from old commits")
		default:
			log.Fatal("not implemented")
		}
	} else {
		log.Println("patches already saved... moving on") // add rewrite option
	}

	// reload patches instead,
	reloaded_patches := make(map[string]patch_generator.Patch, len(patches))
	path_list, filename_list = patchhandler.GetFilesInPersistence()

	if len(path_list) != len(filename_list) {
		log.Fatal("len path list != len filename list")
	}
	for i := 0; i < len(path_list); i++ {
		// refactor so that new audio name is gotten here, and passed into doesfile exist & like

		// maybe binary search -- cached now testing
		if audiohandler.DoesFileExistByPatchName(filename_list[i]) != false {
			log.Println(filename_list[i], "already saved as audio")
			//TODO: compare patch contents
			// weird behavior here?

			continue
		}

		// load patch
		reloaded_patches[filename_list[i]] = patchhandler.LoadPatchByPath(path_list[i])

		// load into audioproc
		patchhandler.LoadPatchIntoEmulator(reloaded_patches[filename_list[i]], &synth_capabilities, patch_access)

		// TODO
		// PLAY  HERE
		// mods.keyboard.NoteEventFromVirtual(60, true)

		// (&audioProcessor).keyboard.NoteEventFromVirtual(60, true) // need to change audiomodule interface
		signal := audiohandler.GetSignal(audioProcessor, 5)

		audiohandler.PersistAudioWithName(signal, filename_list[i])

	}

	// log.Fatal("end")

	// for i := 0; i < len(patches); i++ {
	// 	log.Println("loading patch ", i)

	// 	log.Fatal("FAULTY BELOW THIS CHECK FUNCTION")
	// 	patch_name, err := patchhandler.PersistPatch(patches[i], 0)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	log.Fatal("dataset_generation finished")

}
