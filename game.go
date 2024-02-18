package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	name          string
	game_map      Map
	hud           GameHud
	player        Player
	weapon_holder WeaponHolder
}

func (g *Game) draw() {

	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode3D(g.player.camera)

	rl.DrawGrid(100, 1)

	g.game_map.Update(g.player.camera)

	rl.EndMode3D()

	g.hud.Draw()

	rl.EndDrawing()

}

func (g *Game) game_init() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "BloodCaster")
	rl.SetTargetFPS(FPS)
	rl.DisableCursor()

	g.init_map()
	g.init_weapons()
	g.init_player()
	g.init_hud()
}

func (g *Game) init_player() {
	g.player = Player_init(&g.game_map, &g.weapon_holder)
}

func (g *Game) init_hud() {
	g.hud = Hud_init(&g.weapon_holder, &g.player)
}

func (g *Game) init_map() {
	g.game_map = Map_init()
}

func (g *Game) init_weapons() {
	g.weapon_holder = Weapons_init()
}

func (g *Game) run() {
	for !rl.WindowShouldClose() {
		//camera update and collision detection
		g.player.Update()
		g.draw()
		defer g.game_map.Unload()
	}
	rl.CloseWindow()
}
