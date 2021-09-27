package main

import "github.com/soypat/three"

// color can be specified like so:
// "rgb(255, 0, 0)"
// "rgb(100%, 0%, 0%)"
// "skyblue" // X11 color names (without CamelCase), see three.js source
// "hsl(0, 100%, 50%)"
func ColoredLine(color string, width float64) *three.LineBasicMaterial {
	param := three.NewMaterialParameters()
	param.Color = three.NewColor(color)
	param.LineWidth = width
	return three.NewLineBasicMaterial(param)
}

func RedLine(width float64) *three.LineBasicMaterial {
	param := three.NewMaterialParameters()
	param.Color = three.NewColor("rgb(255,0,0)")
	param.LineWidth = width
	return three.NewLineBasicMaterial(param)
}

func GreenLine(width float64) *three.LineBasicMaterial {
	param := three.NewMaterialParameters()
	param.Color = three.NewColor("rgb(0,255,0)")
	param.LineWidth = width
	return three.NewLineBasicMaterial(param)
}

func BlueLine(width float64) *three.LineBasicMaterial {
	param := three.NewMaterialParameters()
	param.Color = three.NewColor("rgb(0,0,255)")
	param.LineWidth = width
	return three.NewLineBasicMaterial(param)
}

func WhiteLine(width float64) *three.LineBasicMaterial {
	param := three.NewMaterialParameters()
	param.LineWidth = width
	return three.NewLineBasicMaterial(param)
}

// color can be specified like so:
// "rgb(255, 0, 0)"
// "rgb(100%, 0%, 0%)"
// "skyblue" // X11 color names (without CamelCase), see three.js source
// "hsl(0, 100%, 50%)"
func ColoredLambertSurface(color string) *three.MeshLambertMaterial {
	boxparam := three.NewMaterialParameters()
	boxparam.Color = three.NewColor(color)
	return three.NewMeshLambertMaterial(boxparam)
}
