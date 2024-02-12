package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameHud struct {
	player    string
	crosshair rl.Texture2D
}

func (h *GameHud) Draw() {
	rl.DrawTexture(h.crosshair, HALF_WIDTH, HALF_HEIGHT, rl.White)
}

func Hud_init() GameHud {
	hud := GameHud{}
	hud.player = "player"
	hud.crosshair = Load_texture("./assets/textures/crosshair.png")

	return hud
}
