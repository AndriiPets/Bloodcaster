package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	//game
	WALLS_HEIGHT  float32 = 0.5
	SCREEN_WIDTH  int32   = 1280
	SCREEN_HEIGHT int32   = 720
	HALF_WIDTH    int32   = 640
	HALF_HEIGHT   int32   = 360
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
)
