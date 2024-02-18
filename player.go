package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	next_fire     float32
	camera_type   PlayerCustomCamera
	camera        rl.Camera3D
	pos_x         float32
	pos_y         float32
	game_map      *Map
	weapons       *WeaponHolder
	swing_delta   float32
	player_heigth float32
}

func Player_init(game_map *Map, weapons *WeaponHolder) Player {

	p := Player{}

	p.camera_type = PlayerCustomCamera{}
	p.pos_x = game_map.player_pos.X
	p.pos_y = game_map.player_pos.Y
	p.camera = p.camera_type.Player_camera_init(p.pos_x, p.pos_y)
	p.game_map = game_map
	p.weapons = weapons
	p.player_heigth = WALLS_HEIGHT
	p.next_fire = 0.0

	return p
}

func (p *Player) Update() {
	oldCamPos := p.camera.Position

	p.swing_delta = p.camera_type.Player_update_camera(&p.camera)

	playerPos := rl.NewVector3(p.camera.Position.X, p.player_heigth, p.camera.Position.Z)
	if p.game_map.Check_wall_collision(playerPos, 0.3) {
		p.camera.Position = oldCamPos
	}

	p.weapons.Get_switch_input()
	p.fire_weapon(&p.camera)

}

func (p *Player) fire_weapon(camera *rl.Camera3D) {
	p.next_fire -= rl.GetFrameTime()
	//fmt.Println(p.next_fire)

	if rl.IsMouseButtonDown(KEY_FIRE) {
		p.next_fire = p.weapons.Weapon_fire(camera, p.game_map, p.next_fire)
	}
}
