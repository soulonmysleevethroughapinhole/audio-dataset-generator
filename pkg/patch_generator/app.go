package patch_generator

import (
	"log"

	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/emulator"
	"github.com/soulonmysleevethroughapinhole/audio-generator/pkg/options"

	// "github.com/soulonmysleevethroughapinhole/audio-generator/pkg/patch"
	patchpkg "github.com/soulonmysleevethroughapinhole/audio-generator/pkg/patch"
)

/* called by dataset generator, which passes the synthesizer's capabilities / options
according to which the patches will be generated

synth capabilies could just be the accessparams, since it includes
all that can be modified about the synth

and perhaps at time goes on dynamically adding processors will change this behavior

options should include step size for params, etc
*/

type Patch_Generator struct {
	patch_accesses []patchpkg.Patch_access
	step_size      float32

	Patches []patchpkg.Patch
	params  map[string]map[int]map[string]emulator.AccessParam
}

func (p Patch_Generator) recTraverseParamMap(param_i int, patches *[]patchpkg.Patch) {
	layer, n_layer, param := p.patch_accesses[param_i].Values()
	accessparam := p.params[layer][n_layer][param]
	// ranges := accessparam.GetRanges()
	ranges := accessparam.GetRangesPtr() // return to pointers
	p.recursiveStepAccessParam(ranges, 0, param_i, patches)
}

func (p Patch_Generator) recursiveStepAccessParam(ranges []*emulator.Range, i int, param_i int, patches *[]patchpkg.Patch) {

	ranges[i].Value = ranges[i].Min

	for (ranges)[i].Value <= (ranges)[i].Max {
		if ranges[i].Value > ranges[i].Max {
			log.Fatal("this should not happen anymore, i think this is a pointer problem again")
			break
		}

		// if this isnt the last param, go to the next param so that its ranges can be recursed
		if i != len(ranges)-1 {
			p.recursiveStepAccessParam(ranges, i+1, param_i, patches)
		} else if param_i != len(p.patch_accesses)-1 {
			p.recTraverseParamMap(param_i+1, patches)
		}

		// if this is the last range and last param --> save
		if param_i == len(p.patch_accesses)-1 && i == len(ranges)-1 {
			// log.Println(layer_debug, n_layer_debug, param_debug, accessparam_name_debug)
			// log.Println("saved here")

			copy_params := make(map[string]map[int]map[string]emulator.AccessParam)
			copy_patch_param_values := make(map[string]map[int]map[string]map[string]float32)

			// Copy from the original map to the target map
			for layer := range p.params {
				ln_map := make(map[int]map[string]emulator.AccessParam) // old

				ln_val_map := make(map[int]map[string]map[string]float32) // new

				for layer_n := range p.params[layer] {
					ap_map := make(map[string]emulator.AccessParam) // old

					ap_val_map := make(map[string]map[string]float32) // new

					for ap_name := range p.params[layer][layer_n] {
						val_map := make(map[string]float32)

						ap_map[ap_name] = p.params[layer][layer_n][ap_name].CopyParam() // old

						for _, paramvalue := range ap_map[ap_name].GetParamValues() {
							val_map[paramvalue.Name] = paramvalue.Value
						}
						ap_val_map[ap_name] = val_map // new
					}
					ln_map[layer_n] = ap_map         // old
					ln_val_map[layer_n] = ap_val_map // new
				}
				copy_params[layer] = ln_map                 // old
				copy_patch_param_values[layer] = ln_val_map // new
			}
			foo_patch := patchpkg.Patch{Accessparams: copy_params, Accessvalues: copy_patch_param_values}
			*patches = append(*patches, foo_patch)
		}
		if ranges[i].Value > ranges[i].Max {
			log.Fatal("what")
		}

		if float32(int(ranges[i].Value)+int((ranges[i].Max-ranges[i].Min)/5)) <= ranges[i].Max {
			// ranges[i].Value = float32(int(ranges[i].Value) + int((ranges[i].Max-ranges[i].Min)/2))
			ranges[i].Value = float32(int(ranges[i].Value) + int((ranges[i].Max-ranges[i].Min)/5))
			// ranges[i].Value = float32(int(ranges[i].Value) + int((ranges[i].Max-ranges[i].Min)/8))
			// log.Println(layer_debug, n_layer_debug, param_debug, accessparam_name_debug, " value increased by ", (ranges[i].Max-ranges[i].Min)/2, "to ", ranges[i].Value)
		} else {
			break
		}
	}
}

func (p Patch_Generator) demoTraverseParamMap(foo []int, i int, bar *[][]int) {
	foo[i] = 0
	for foo[i] < 10 {
		if i == len(foo)-1 {
			cpy := make([]int, len(foo))
			copy(cpy, foo)

			*bar = append(*bar, cpy)
		} else {
			p.demoTraverseParamMap(foo, i+1, bar)
		}
		foo[i] += 1
	}
}

func Load(synth_capabilities map[string]map[int]map[string]emulator.AccessParam) Patch_Generator {
	gen := Patch_Generator{step_size: 0.10}

	for layer := range synth_capabilities {
		for layer_n := range synth_capabilities[layer] {
			for ap_name := range synth_capabilities[layer][layer_n] {
				gen.patch_accesses = append(gen.patch_accesses, patchpkg.Patch_access{Layer: layer, N_layer: layer_n, Param: ap_name})
			}
		}
	}

	return gen
}

func App(synth_capabilities map[string]map[int]map[string]emulator.AccessParam, ws_options options.Option) ([]patchpkg.Patch, []patchpkg.Patch_access) {
	gen := Load(synth_capabilities)
	gen.params = synth_capabilities
	var patches []patchpkg.Patch
	// log.Fatal("problem is here")
	gen.recTraverseParamMap(0, &patches)
	return patches, gen.patch_accesses
}
