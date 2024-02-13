package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameHud struct {
	player        string //give hud player and weapons object to draw health, ammo ect.
	weapon_holder *WeaponHolder
	crosshair     rl.Texture2D
}

func (h *GameHud) Draw() {
	rl.DrawTexture(h.crosshair, HALF_WIDTH-(h.crosshair.Width/2), HALF_HEIGHT-(h.crosshair.Width/2), rl.White)
	h.weapon_holder.Draw()
}

func Hud_init(weapons *WeaponHolder) GameHud {
	hud := GameHud{}
	hud.player = "player"
	hud.weapon_holder = weapons
	hud.crosshair = Load_texture("./assets/textures/crosshair.png")

	return hud
}
