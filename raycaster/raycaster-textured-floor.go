package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	texWidth  = 128
	texHeight = 128
)

var (
	fullscreen = false
	showMap    = false
	width      = 320
	height     = 200
	scale      = 3.0

	pos   = pixel.V(18.0, 9.5)
	dir   = pixel.V(-1.0, 0.0)
	plane = pixel.V(0.0, 0.66)

	floorTex = floorTexture()
)

var world = [24][24]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 9, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 2, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 5, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 5, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 5, 5, 0, 5, 5, 5, 0, 5, 5, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 9, 0, 4, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 5, 0, 5, 5, 5, 5, 5, 5, 5, 0, 5, 0, 1},
	{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 5, 0, 5, 0, 0, 0, 0, 0, 5, 0, 5, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 5, 5, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

func getColor(x, y int) color.RGBA {
	switch world[x][y] {
	case 0:
		return color.RGBA{64, 64, 64, 255}
	case 1:
		return color.RGBA{244, 115, 33, 255}
	case 2:
		return color.RGBA{54, 124, 43, 255}
	case 3:
		return color.RGBA{0, 125, 198, 255}
	case 4:
		return color.RGBA{160, 32, 240, 255}
	case 5:
		return color.RGBA{235, 235, 235, 255}
	default:
		return color.RGBA{255, 194, 32, 255}
	}
}

func frame() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, image.Rect(0, 0, width, height/2), &image.Uniform{color.RGBA{192, 192, 192, 255}}, image.ZP, draw.Src)

	for x := 0; x < width; x++ {
		var (
			step         image.Point
			sideDist     pixel.Vec
			perpWallDist float64
			hit, side    bool

			rayPos, worldX, worldY = pos, int(pos.X), int(pos.Y)

			cameraX = 2*float64(x)/float64(width) - 1

			rayDir = pixel.V(
				dir.X+plane.X*cameraX,
				dir.Y+plane.Y*cameraX,
			)

			deltaDist = pixel.V(
				math.Sqrt(1.0+(rayDir.Y*rayDir.Y)/(rayDir.X*rayDir.X)),
				math.Sqrt(1.0+(rayDir.X*rayDir.X)/(rayDir.Y*rayDir.Y)),
			)
		)

		if rayDir.X < 0 {
			step.X = -1
			sideDist.X = (rayPos.X - float64(worldX)) * deltaDist.X
		} else {
			step.X = 1
			sideDist.X = (float64(worldX) + 1.0 - rayPos.X) * deltaDist.X
		}

		if rayDir.Y < 0 {
			step.Y = -1
			sideDist.Y = (rayPos.Y - float64(worldY)) * deltaDist.Y
		} else {
			step.Y = 1
			sideDist.Y = (float64(worldY) + 1.0 - rayPos.Y) * deltaDist.Y
		}

		for !hit {
			if sideDist.X < sideDist.Y {
				sideDist.X += deltaDist.X
				worldX += step.X
				side = false
			} else {
				sideDist.Y += deltaDist.Y
				worldY += step.Y
				side = true
			}

			if world[worldX][worldY] > 0 {
				hit = true
			}
		}

		var wallX float64

		if side {
			perpWallDist = (float64(worldY) - rayPos.Y + (1-float64(step.Y))/2) / rayDir.Y
			wallX = rayPos.X + perpWallDist*rayDir.X
		} else {
			perpWallDist = (float64(worldX) - rayPos.X + (1-float64(step.X))/2) / rayDir.X
			wallX = rayPos.Y + perpWallDist*rayDir.Y
		}

		wallX -= math.Floor(wallX)

		lineHeight := int(float64(height) / perpWallDist)

		drawStart := -lineHeight/2 + height/2
		if drawStart < 0 {
			drawStart = 0
		}

		drawEnd := lineHeight/2 + height/2
		if drawEnd >= height {
			drawEnd = height - 1
		}

		c := getColor(worldX, worldY)

		if side {
			c.R = c.R / 2
			c.G = c.G / 2
			c.B = c.B / 2
		}

		for y := drawStart; y < drawEnd+1; y++ {
			if y == drawStart {
				m.Set(x, y, color.RGBA{c.R / 2, c.G / 2, c.B / 2, 255})
			} else if y == drawEnd {
				m.Set(x, y, color.RGBA{c.R / 2, c.G / 2, c.B / 2, 255})
			} else {
				m.Set(x, y, c)
			}
		}

		var floorWall pixel.Vec

		if !side && rayDir.X > 0 {
			floorWall.X = float64(worldX)
			floorWall.Y = float64(worldY) + wallX
		} else if !side && rayDir.X < 0 {
			floorWall.X = float64(worldX) + 1.0
			floorWall.Y = float64(worldY) + wallX
		} else if side && rayDir.Y > 0 {
			floorWall.X = float64(worldX) + wallX
			floorWall.Y = float64(worldY)
		} else {
			floorWall.X = float64(worldX) + wallX
			floorWall.Y = float64(worldY) + 1.0
		}

		distWall, distPlayer := perpWallDist, 0.0

		for y := drawEnd + 1; y < height; y++ {
			currentDist := float64(height) / (2.0*float64(y) - float64(height))

			weight := (currentDist - distPlayer) / (distWall - distPlayer)

			currentFloor := pixel.V(
				weight*floorWall.X+(1.0-weight)*pos.X,
				weight*floorWall.Y+(1.0-weight)*pos.Y,
			)

			m.Set(x, y, floorTex.At(
				int(currentFloor.X*float64(texWidth))%texWidth,
				int(currentFloor.Y*float64(texHeight))%texHeight,
			))
		}
	}

	return m
}

func minimap() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, 24, 24))

	for x, row := range world {
		for y, _ := range row {
			m.Set(x, y, getColor(x, y))
		}
	}

	m.Set(int(pos.X), int(pos.Y), color.RGBA{255, 0, 0, 255})

	return m
}

func floorTexture() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, texWidth, texHeight))

	for x := 0; x < texWidth; x++ {
		for y := 0; y < texHeight; y++ {
			c := uint8(x ^ y)

			m.Set(x, y, color.RGBA{c, c, c, 255})
		}
	}

	return m
}

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:       true,
		Undecorated: true,
	}

	if fullscreen {
		cfg.Monitor = pixelgl.PrimaryMonitor()
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	c := win.Bounds().Center()

	last := time.Now()

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			moveForward(5 * dt)
		}

		if win.Pressed(pixelgl.KeyA) {
			moveLeft(5 * dt)
		}

		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			moveBackwards(5 * dt)
		}

		if win.Pressed(pixelgl.KeyD) {
			moveRight(5 * dt)
		}

		if win.Pressed(pixelgl.KeyRight) {
			turnRight(2 * dt)
		}

		if win.Pressed(pixelgl.KeyLeft) {
			turnLeft(2 * dt)
		}

		if win.JustPressed(pixelgl.KeyM) {
			showMap = !showMap
		}

		p := pixel.PictureDataFromImage(frame())

		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		if showMap {
			m := pixel.PictureDataFromImage(minimap())

			mc := m.Bounds().Min.Add(pixel.V(-m.Rect.W()/2, m.Rect.H()/2))

			pixel.NewSprite(m, m.Bounds()).
				Draw(win, pixel.IM.
					Moved(mc).
					Rotated(mc, -1.57453626).
					ScaledXY(pixel.ZV, pixel.V(-scale*2, scale*2)))
		}

		win.Update()
	}
}

func moveForward(s float64) {
	if world[int(pos.X+dir.X*s)][int(pos.Y)] == 0 {
		pos.X += dir.X * s
	}

	if world[int(pos.X)][int(pos.Y+dir.Y*s)] == 0 {
		pos.Y += dir.Y * s
	}
}

func moveLeft(s float64) {
	if world[int(pos.X-plane.X*s)][int(pos.Y)] == 0 {
		pos.X -= plane.X * s
	}

	if world[int(pos.X)][int(pos.Y-plane.Y*s)] == 0 {
		pos.Y -= plane.Y * s
	}
}

func moveBackwards(s float64) {
	if world[int(pos.X-dir.X*s)][int(pos.Y)] == 0 {
		pos.X -= dir.X * s
	}

	if world[int(pos.X)][int(pos.Y-dir.Y*s)] == 0 {
		pos.Y -= dir.Y * s
	}
}

func moveRight(s float64) {
	if world[int(pos.X+plane.X*s)][int(pos.Y)] == 0 {
		pos.X += plane.X * s
	}

	if world[int(pos.X)][int(pos.Y+plane.Y*s)] == 0 {
		pos.Y += plane.Y * s
	}
}

func turnRight(s float64) {
	dir.Y = dir.X*math.Sin(-s) + dir.Y*math.Cos(-s)
	dir.X = dir.X*math.Cos(-s) - dir.Y*math.Sin(-s)

	plane.Y = plane.X*math.Sin(-s) + plane.Y*math.Cos(-s)
	plane.X = plane.X*math.Cos(-s) - plane.Y*math.Sin(-s)
}

func turnLeft(s float64) {
	dir.Y = dir.X*math.Sin(s) + dir.Y*math.Cos(s)
	dir.X = dir.X*math.Cos(s) - dir.Y*math.Sin(s)

	plane.Y = plane.X*math.Sin(s) + plane.Y*math.Cos(s)
	plane.X = plane.X*math.Cos(s) - plane.Y*math.Sin(s)
}

func main() {
	flag.BoolVar(&fullscreen, "f", fullscreen, "fullscreen")
	flag.IntVar(&width, "w", width, "width")
	flag.IntVar(&height, "h", height, "height")
	flag.Float64Var(&scale, "s", scale, "scale")
	flag.Parse()

	pixelgl.Run(run)
}
