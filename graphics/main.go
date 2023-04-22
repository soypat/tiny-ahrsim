package main

import (
	"encoding/json"
	"log"
	"net/http"
	"syscall/js"
	"time"

	"github.com/soypat/ahrs"
	"github.com/soypat/gwasm"
	"github.com/soypat/three"
)

func main() {
	gwasm.AddScript("https://threejs.org/build/three.js", "THREE", time.Second)
	err := three.Init()
	if err != nil {
		panic(err)
	}
	scale := 200. // Scale will define overall size. All objects will be scaled accordingly.
	log.Println("begin program")
	// Get size of window.
	width := js.Global().Get("innerWidth").Float()
	height := js.Global().Get("innerHeight").Float()

	// Bring in the heavyweight renderer.
	// This will render scenes passed into it.
	// You may add it to DOM using appendChild (see end of this program).
	log.Println("create renderer 1")
	renderer := three.NewWebGLRenderer(three.WebGLRendererParam{})
	log.Println("create renderer")
	renderer.SetSize(width, height, true)
	log.Println("finish creating renderer")
	// setup camera and scene
	camera := three.NewPerspectiveCamera(70, width/height, 1, scale*5)
	camera.SetUp(three.NewVector3(0, 0, 1))
	// put camera along x=y=z line to get nice ISO-view
	camera.SetPosition(three.NewVector3(scale*3, scale*3, scale*3))

	camera.LookAt(three.NewVector3(0, 0, 0)) // Look at origin so cube is inside view.

	// Scene is passed as an argument to renderer. Scene contains 3D objects.
	scene := three.NewScene()
	// lights, without lights everything will be dark! second and last argument to renderer.
	light := three.NewDirectionalLight(three.NewColor("white"), 1)
	light.SetPosition(three.NewVector3(0, scale*1.3, scale*1.5))

	scene.Add(light) // This is the idiom to add objects to scene.

	// Create Axis lines

	// xline := three.NewLine(lineGeom(zero, xnorm.scale(scale*1.5)), RedLine(1))
	// yline := three.NewLine(lineGeom(zero, ynorm.scale(scale*1.5)), GreenLine(1))
	// zline := three.NewLine(lineGeom(zero, znorm.scale(scale*1.5)), BlueLine(1))
	rigidBody := three.NewGroup()
	rigidBody.Add(three.NewAxesHelper(scale * 3))
	// rigidBody.Add(yline)
	// rigidBody.Add(zline)
	// cube object
	geom := three.NewBoxGeometry(three.BoxGeometryParameters{
		Width:  scale,
		Height: scale,
		Depth:  scale,
	})
	boxmat := ColoredLambertSurface("skyblue")
	mesh := three.NewMesh(geom, boxmat)
	rigidBody.Add(mesh)

	scene.Add(rigidBody)

	// We create a recursive callback to continuously animate our project.
	var animate func()
	ang := ahrs.EulerAngles{}
	animate = func() {
		rigidBody.SetRotationFromEuler(three.NewEuler(ang.Q, ang.R, ang.S, ""))
		renderer.Render(scene, camera)
		js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			animate()
			return nil
		}))
	}
	// Add renderer to DOM (HTML page).
	js.Global().Get("document").Get("body").Call("appendChild", renderer.Get("domElement"))
	// start animation using recursive callback method.
	// Each time a frame is finished rendering a new request to animate is called.
	animate()

	wait := make(chan bool)
	for {
		go func() (err error) {
			defer func() {
				if err != nil {
					time.Sleep(time.Second)
				}
				wait <- true
			}()
			resp, err := http.Get("http://localhost:8081/attitude")
			if err != nil {
				return err
			}
			d := json.NewDecoder(resp.Body)
			err = d.Decode(&ang)
			if err != nil {
				return err
			}
			return nil
		}()
		<-wait
		time.Sleep(100 * time.Millisecond)
	}
}

// Line geometry creation.

type vec struct {
	x, y, z float64
}

func (v *vec) scale(a float64) *vec {
	v.x *= a
	v.y *= a
	v.z *= a
	return v
}

var (
	zero  = &vec{}
	xnorm = &vec{1, 0, 0}
	ynorm = &vec{0, 1, 0}
	znorm = &vec{0, 0, 1}
)

func lineGeom(to, from *vec) three.Geometry {
	bufGeom := three.NewBufferGeometry()
	bufGeom.SetAttribute("positions", three.NewBufferAttribute([]float32{
		float32(to.x), float32(to.y), float32(to.z),
		float32(from.x), float32(from.y), float32(from.z),
	}, 3))
	return bufGeom
}

type rigidBody struct {
	subBodies []*three.Euler
}

func NewRigidBody(rotations ...*three.Euler) rigidBody {
	return rigidBody{
		subBodies: rotations,
	}
}
