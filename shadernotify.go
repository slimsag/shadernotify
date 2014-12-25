package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/lmath.v1"

	"davsk.net/procedural"

	"gopkg.in/fsnotify.v1"
)

var objects []*gfx.Object

func watchShaders() (chan string, chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	event := make(chan string)
	quit := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				// expects filenames such as triangle-vert.glsl or triangle-frag.glsl
				if filepath.Ext(ev.Name) == ".glsl" {
					b := filepath.Base(ev.Name)
					i := strings.LastIndex(b, "-")
					b = b[:i]
					event <- b
				}
			case err := <-watcher.Errors:
				log.Println("watch error:", err)
			case <-quit:
				watcher.Close()
				return
			}
		}
	}()

	err = watcher.Add("./")
	if err != nil {
		log.Fatal("Failed to watch folder.", err)
	}

	return event, quit
}

func loadShader(s *gfx.Shader) {
	vert, err := ioutil.ReadFile(fmt.Sprintf("./%s-vert.glsl", s.Name))
	if err != nil {
		log.Fatal("Failed to locate vertex shader.", err)
	}

	frag, err := ioutil.ReadFile(fmt.Sprintf("./%s-frag.glsl", s.Name))
	if err != nil {
		log.Fatal("Failed to located fragment shader.", err)
	}

	s.GLSLVert = vert
	s.GLSLFrag = frag
}

func gfxLoop(w window.Window, r gfx.Renderer) {
	watchEvent, watchQuit := watchShaders()

	camera := gfx.NewCamera()
	camera.SetPersp(r.Bounds(), 75, 0.0001, 1000.0)
	camera.SetPos(lmath.Vec3{0, -20, 0})

	go func() {
		event := make(chan window.Event, 32)
		w.Notify(event, window.ResizedEvents|window.CloseEvents)

		for e := range event {
			switch e.(type) {
			case window.Resized:
				camera.Lock()
				camera.SetPersp(r.Bounds(), 75, 0.0001, 1000.0)
				camera.Unlock()
			case window.Close:
				watchQuit <- true
			}
		}
	}()

	o := gfx.NewObject()
	m := procedural.Sphere(5)
	o.Meshes = []*gfx.Mesh{m}
	o.Shader = gfx.NewShader("rocky")
	o.SetScale(lmath.Vec3{10, 10, 10})
	o.SetPos(lmath.Vec3{0, 20, 0})
	o.SetRot(lmath.Vec3{180, 0, 0})
	o.FaceCulling = gfx.NoFaceCulling
	f, err := os.Open("./rocky.jpg")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	tex := &gfx.Texture{
		Source: img,
		Bounds: img.Bounds(),
		WrapU:  gfx.Repeat,
		WrapV:  gfx.Repeat,
		Format: gfx.DXT1RGBA,
	}
	o.Textures = []*gfx.Texture{tex}
	objects = append(objects, o)

	offset := lmath.Vec3{0, 0.2, 1}
	for {
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{0, 0, 0, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		var reload string
		select {
		case reload = <-watchEvent:
		default:
		}

		for _, obj := range objects {
			if obj.Shader != nil {
				if obj.Shader.Name == reload {
					obj.Shader.Lock()
					name := obj.Shader.Name
					obj.Shader.Reset()
					obj.Shader.Name = name
					obj.Shader.Unlock()
				}
				if !obj.Shader.Loaded {
					loadShader(obj.Shader)
					done := make(chan *gfx.Shader, 1)
					r.LoadShader(obj.Shader, done)
					<-done
				}
				obj.Shader.Inputs["time"] = float32(r.Clock().Time().Seconds())
			}
			obj.SetRot(obj.Rot().Add(offset))
			r.Draw(r.Bounds(), obj, camera)
		}

		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
