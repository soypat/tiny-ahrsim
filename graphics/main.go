package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/soypat/ahrs"
	three "github.com/soypat/gthree"
)

func main() {
	scale := 200. // Scale will define overall size. All objects will be scaled accordingly.

	// Get size of window.
	width := js.Global.Get("innerWidth").Float()
	height := js.Global.Get("innerHeight").Float()
	// Bring in the heavyweight renderer.
	// This will render scenes passed into it.
	// You may add it to DOM using appendChild (see end of this program).
	renderer := three.NewWebGLRenderer()
	renderer.SetSize(width, height, true)

	// setup camera and scene
	camera := three.NewPerspectiveCamera(70, width/height, 1, scale*5)
	camera.Up.Set(0, 0, 1)
	// put camera along x=y=z line to get nice ISO-view
	camera.Position.Set(scale*2, scale*2, scale*2)
	camera.LookAt(0, 0, 0) // Look at origin so cube is inside view.

	// Scene is passed as an argument to renderer. Scene contains 3D objects.
	scene := three.NewScene()
	// lights, without lights everything will be dark! second and last argument to renderer.
	light := three.NewDirectionalLight(three.NewColor("white"), 1)
	light.Position.Set(0, scale*1.3, scale*1.5)
	scene.Add(light) // This is the idiom to add objects to scene.

	// Create Axis lines
	xline := three.NewLine(lineGeom(zero, xnorm.scale(scale*1.5)), RedLine(1))
	yline := three.NewLine(lineGeom(zero, ynorm.scale(scale*1.5)), GreenLine(1))
	zline := three.NewLine(lineGeom(zero, znorm.scale(scale*1.5)), BlueLine(1))
	scene.Add(xline)
	scene.Add(yline)
	scene.Add(zline)

	// cube object
	geom := three.NewBoxGeometry(&three.BoxGeometryParameters{
		Width:  scale,
		Height: scale,
		Depth:  scale,
	})
	boxmat := ColoredLambertSurface("skyblue")
	mesh := three.NewMesh(geom, boxmat)
	scene.Add(mesh)

	// Generate a rigid body to rotate
	body := NewRigidBody(xline.Rotation, yline.Rotation, zline.Rotation, mesh.Rotation)

	// We create a recursive callback to continuously animate our project.
	var animate func()
	ang := ahrs.EulerAngles{}
	animate = func() {
		body.set(ang.Q, ang.R, ang.S)
		renderer.Render(scene, camera)
		js.Global.Call("requestAnimationFrame", animate)
	}
	// Add renderer to DOM (HTML page).
	js.Global.Get("document").Get("body").Call("appendChild", renderer.Get("domElement"))
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
			resp, err := http.Get("http://localhost:8080/attitude")
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

func lineGeom(to, from *vec) *three.BasicGeometry {
	geom := three.NewBasicGeometry(three.BasicGeometryParams{})
	geom.AddVertex(to.x, to.y, to.z)
	geom.AddVertex(from.x, from.y, from.z)
	return &geom
}

type rigidBody struct {
	subBodies []*three.Euler
}

func NewRigidBody(rotations ...*three.Euler) rigidBody {
	return rigidBody{
		subBodies: rotations,
	}
}

func (r *rigidBody) rotate(x, y, z float64) {
	for i := range r.subBodies {
		r.subBodies[i].Set(r.subBodies[i].X+x, r.subBodies[i].Y+y, r.subBodies[i].Z+z, "XYZ")
	}
}
func (r *rigidBody) set(x, y, z float64) {
	for i := range r.subBodies {
		r.subBodies[i].Set(x, y, z, "XYZ")
	}
}
