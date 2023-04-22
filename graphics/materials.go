package main

import "github.com/soypat/three"

// color can be specified like so:
// "rgb(255, 0, 0)"
// "rgb(100%, 0%, 0%)"
// "skyblue" // X11 color names (without CamelCase), see three.js source
// "hsl(0, 100%, 50%)"
func ColoredLine(color string, width float64) three.LineBasicMaterial {
	return three.NewLineBasicMaterial(&three.MaterialParameters{
		Color:     three.NewColor(color),
		LineWidth: width,
	})
}

func RedLine(width float64) three.LineBasicMaterial {
	return ColoredLine("rgb(255,0,0)", width)
}

func GreenLine(width float64) three.LineBasicMaterial {
	return ColoredLine("rgb(0,255,0)", width)

}

func BlueLine(width float64) three.LineBasicMaterial {
	return ColoredLine("rgb(0,0,255)", width)
}

func WhiteLine(width float64) three.LineBasicMaterial {
	return ColoredLine("rgb(255,255,255)", width)
}

// color can be specified like so:
// "rgb(255, 0, 0)"
// "rgb(100%, 0%, 0%)"
// "skyblue" // X11 color names (without CamelCase), see three.js source
// "hsl(0, 100%, 50%)"
func ColoredLambertSurface(color string) three.MeshLambertMaterial {
	return three.NewMeshLambertMaterial(&three.MaterialParameters{
		Color: three.NewColor(color),
	})
}
