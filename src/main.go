package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
)

var (
	width  uint = 1200
	height uint = 1200

	ex = m.Vector{1, 0, 0}
	ey = m.Vector{0, 1, 0}
	ez = m.Vector{0, 0, 1}

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Println("Creating scene...")
	scene := CornellBox()
	fmt.Println("Rendering...")

	params := render.Params{
		Scene:        scene,
		NumWorkers:   10,
		NumSamples:   200,
		AntiAliasing: true,
		TracerType:   m.WhittedStyle,
	}

	// aw := render.NewAVI("out.avi", width, height)
	film := render.Render(params)
	film.SaveAsPNG("out.png")

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		defer f.Close()
		runtime.GC()
		pprof.Lookup("allocs").WriteTo(f, 0)
	}
}
