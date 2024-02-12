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
	hud         GameHud
}

func (g *Game) draw() {

	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode3D(g.camera)

	rl.DrawGrid(100, 1)

	g.game_map.Update()

	rl.EndMode3D()

	g.hud.Draw()

	rl.EndDrawing()

}

func (g *Game) game_init() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "BloodCaster")
	rl.SetTargetFPS(FPS)
	rl.DisableCursor()

	g.init_map()
	g.init_camera()
	g.init_hud()
}

func (g *Game) init_camera() {
	g.camera_type = PlayerCustomCamera{}
	g.camera = g.camera_type.Player_camera_init(g.game_map.player_pos.X, g.game_map.player_pos.Y)
}

func (g *Game) init_hud() {
	g.hud = Hud_init()
}

func (g *Game) init_map() {
	g.game_map = Map_init()
}

func (g *Game) run() {
	for !rl.WindowShouldClose() {
		//camera update and collision detection
		oldCamPos := g.camera.Position

		g.camera_type.Player_update_camera(&g.camera)

		playerPos := rl.NewVector2(g.camera.Position.X, g.camera.Position.Z)
		if g.game_map.Check_wall_collision(playerPos, 0.3) {
			g.camera.Position = oldCamPos
		}

		rl.SetWindowTitle(fmt.Sprintln("player pos: ", g.camera.Position.X, g.camera.Position.Z))

		g.draw()
		defer g.game_map.Unload()
	}
	rl.CloseWindow()
}
