package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	name        string
	game_map    Map
	camera_type PlayerCustomCamera
	camera      rl.Camera3D
	yaw         float32
	pitch       float32
}

func (g *Game) draw() {

	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode3D(g.camera)

	rl.DrawGrid(100, 1)

	g.game_map.Update()

	rl.EndMode3D()

	rl.EndDrawing()

}

func (g *Game) game_init() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "BloodCaster")
	rl.SetTargetFPS(FPS)
	rl.DisableCursor()
	g.camera_type = PlayerCustomCamera{}

	g.game_map = Map{}
	g.game_map.Map_init()
	g.camera = g.camera_type.Player_camera_init(g.game_map.player_pos.X, g.game_map.player_pos.Y)
}

func (g *Game) init_camera() {
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(-4.0, 0.4, -4.0)
	camera.Target = rl.NewVector3(0.0, 0.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = FOV
	camera.Projection = rl.CameraPerspective

	g.pitch = -0.6 // mouseDelta y
	g.yaw = -2.45  // mouseDelta x
	camera_front := Utils_calculate_camera_vector(g.yaw, g.pitch)
	camera.Target = rl.Vector3Add(camera.Position, camera_front)

	g.camera = camera
}

func (g *Game) run() {
	for !rl.WindowShouldClose() {
		oldCamPos := g.camera.Position
		//caculate new camera target based on cursor position each frame
		/* delta_time := rl.GetFrameTime()
		mouse_delta := rl.GetMouseDelta()

		g.yaw += mouse_delta.X * delta_time
		g.pitch += -mouse_delta.Y * delta_time

		if g.pitch > 1.5 {

			g.pitch = 1.5

		} else if g.pitch < -1.5 {

			g.pitch = -1.5

		}

		camera_front := Utils_calculate_camera_vector(g.yaw, g.pitch)
		g.camera.Target = rl.Vector3Add(g.camera.Position, camera_front) */

		//rl.UpdateCamera(&g.camera, rl.CameraFirstPerson)
		g.camera_type.Player_update_camera(&g.camera)

		playerPos := rl.NewVector2(g.camera.Position.X, g.camera.Position.Z)
		if g.game_map.Check_wall_collision(playerPos, 0.2) {
			g.camera.Position = oldCamPos
		}

		//player check collision
		rl.SetWindowTitle(fmt.Sprintln("player pos: ", g.camera.Position.X, g.camera.Position.Z))

		g.draw()
		defer g.game_map.Unload()
	}
	rl.CloseWindow()
}

func main() {

	game := Game{}
	game.game_init()
	game.run()

}
