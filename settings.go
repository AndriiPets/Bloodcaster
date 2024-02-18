package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	//game
	WALLS_HEIGHT  float32 = 0.5
	SCREEN_WIDTH  int32   = 1280
	SCREEN_HEIGHT int32   = 720
	HALF_WIDTH    int32   = SCREEN_WIDTH / 2
	HALF_HEIGHT   int32   = SCREEN_HEIGHT / 2
	FPS           int32   = 60
	//map
	MAP_CELL_SIZE = 3.0 //shoud be dividable by map length and width witout reminder
	//camera
	FOV                           = 60.0
	PLAYER_CAMERA_MIN_CLAMP       = 89.0
	PLAYER_CAMERA_MAX_CLAMP       = -89.0
	PLAYER_CAMERA_PANNING_DIVIDER = 5.1
	PLAYER_CAMERA_OFFSET_Y        = 1.85
	//controls
	MOUSE_SENSITIVITY       = 0.35
	KEY_MOVE_FORWARD  int32 = rl.KeyW
	KEY_MOVE_BACKWARD int32 = rl.KeyS
	KEY_MOVE_LEFT     int32 = rl.KeyA
	KEY_MOVE_RIGHT    int32 = rl.KeyD
	KEY_MOVE_UP       int32 = rl.KeySpace
	KEY_MOVE_DOWN     int32 = rl.KeyZ
	KEY_SHIFT         int32 = rl.KeyLeftShift
	KEY_FIRE          int32 = rl.MouseButtonLeft
	KEY_WEAPON_ONE    int32 = rl.KeyOne
	KEY_WEAPON_TWO    int32 = rl.KeyTwo
	KEY_WEAPON_THREE  int32 = rl.KeyThree
	//ID's
	WEAPON_MELEE_ID  int32 = 0
	WEAPON_PISTOL_ID int32 = 1
	WEAPON_RIFLE_ID  int32 = 2
)

// item and weapon settings
var WEAPON_PISTOL Weapon = Weapon{
	name:                   "pistol",
	input_key:              KEY_WEAPON_TWO,
	weapon_id:              WEAPON_PISTOL_ID,
	damage:                 3,
	ammo:                   30,
	fire_rate:              0.8,
	fire_range:             12.0,
	picked_up:              true,
	sprite_speed:           10,
	sprite_fire_frame:      2,
	sprites_total:          5,
	sprite_position_offset: rl.NewVector2(-0.10, 0.75),
}

var WEAPON_MELEE Weapon = Weapon{
	name:                   "melee",
	input_key:              KEY_WEAPON_ONE,
	weapon_id:              WEAPON_MELEE_ID,
	damage:                 5,
	ammo:                   99,
	fire_rate:              1.25,
	fire_range:             2.0,
	picked_up:              true,
	sprite_speed:           10,
	sprite_fire_frame:      1,
	sprites_total:          7,
	sprite_position_offset: rl.NewVector2(0.5, 1.0),
}

var WEAPON_RIFLE Weapon = Weapon{
	name:                   "rifle",
	input_key:              KEY_WEAPON_THREE,
	weapon_id:              WEAPON_RIFLE_ID,
	damage:                 3,
	ammo:                   99,
	fire_rate:              0.5,
	fire_range:             20.0,
	picked_up:              true,
	sprite_speed:           10,
	sprite_fire_frame:      2,
	sprites_total:          5,
	sprite_position_offset: rl.NewVector2(-0.10, 1.0),
}

func Create_weapons_map() map[string]*Weapon {
	w_map := map[string]*Weapon{"pistol": &WEAPON_PISTOL, "melee": &WEAPON_MELEE, "rifle": &WEAPON_RIFLE}

	return w_map
}
