package main

import (
	"maze.io/x/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerCustomCamera struct {
	target_distance float32
	forward         rl.Vector3
	right           rl.Vector3

	//mouse controls
	mouse_sensitivity float32
	use_mouse         bool
	invert_y          bool

	//special options
	focused      bool
	allow_flight bool

	//player stuff
	player_eyes_pos float32
	angle           rl.Vector2
	move_speed      rl.Vector3
	turn_speed      rl.Vector2

	//view bobble stuff
	view_bobble_freq      float32
	view_bobble_mag       float32
	curr_bobble           float32
	view_bobble_waver_mag float32
	swing_delta           float32
}

func (c *PlayerCustomCamera) custom_camera_init() {
	c.target_distance = 0
	c.player_eyes_pos = 0.4
	c.mouse_sensitivity = 0.003
	c.move_speed = rl.NewVector3(1.5, 1, 2.2)
	c.turn_speed = rl.NewVector2(90, 90)
	c.allow_flight = false
	c.view_bobble_freq = 40
	c.view_bobble_mag = 0.03
	c.curr_bobble = 0
	c.view_bobble_waver_mag = 0.005
	c.use_mouse = true
	c.focused = rl.IsWindowFocused()
	c.invert_y = false
}

func (c *PlayerCustomCamera) Player_camera_init(pos_x, pos_z float32) rl.Camera {

	camera := rl.Camera{}
	c.custom_camera_init()

	camera.Position = rl.NewVector3(pos_x, 0.4, pos_z)
	camera.Target = rl.NewVector3(0.0, 0.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = FOV
	camera.Projection = rl.CameraPerspective

	/* //Distance
	v1 := camera.Position
	v2 := camera.Target
	dx := v2.X - v1.X
	dy := v2.Y - v1.Y
	dz := v2.Z - v1.Z

	c.target_distance = math32.Sqrt(dx*dx + dy*dy + dz*dz)

	// Camera angle calculation
	// Camera angle in plane XZ (0 aligned with Z, move positive CCW)
	c.angle.X = math32.Atan2(dx, dz)
	// Camera angle in plane XY (0 aligned with X, move positive CW)
	c.angle.Y = math32.Atan2(dy, math32.Sqrt(dx*dx+dz*dz)) */

	c.player_eyes_pos = camera.Position.Y
	c.mouse_sensitivity = MOUSE_SENSITIVITY

	return camera
}

func If_bool(b bool) float32 {
	if b {
		return 1.0
	}
	return 0.0
}

func (c *PlayerCustomCamera) Get_speed_for_axis(axis int32, speed float32) float32 {
	factor := 1.0
	if rl.IsKeyDown(KEY_SHIFT) {
		factor = 2.0
	}
	if rl.IsKeyDown(axis) {
		return speed * rl.GetFrameTime() * float32(factor)
	}
	return 0
}

func (c *PlayerCustomCamera) Player_update_camera(camera *rl.Camera) float32 {

	mouse_pos_delta := rl.GetMouseDelta()
	direction := map[string]float32{
		"forward":  c.Get_speed_for_axis(KEY_MOVE_FORWARD, c.move_speed.Z),
		"backward": c.Get_speed_for_axis(KEY_MOVE_BACKWARD, c.move_speed.Z),
		"left":     c.Get_speed_for_axis(KEY_MOVE_LEFT, c.move_speed.X),
		"right":    c.Get_speed_for_axis(KEY_MOVE_RIGHT, c.move_speed.X),
		"up":       c.Get_speed_for_axis(KEY_MOVE_UP, c.move_speed.Y),
		"down":     c.Get_speed_for_axis(KEY_MOVE_DOWN, c.move_speed.Y),
	}

	// Camera orientation calculation
	c.angle.X -= mouse_pos_delta.X * c.mouse_sensitivity * rl.GetFrameTime()
	c.angle.Y -= mouse_pos_delta.Y * c.mouse_sensitivity * rl.GetFrameTime()

	//Angle clamp
	if c.angle.Y > PLAYER_CAMERA_MIN_CLAMP*rl.Deg2rad {
		c.angle.Y = PLAYER_CAMERA_MIN_CLAMP * rl.Deg2rad
	} else if c.angle.Y < PLAYER_CAMERA_MAX_CLAMP*rl.Deg2rad {
		c.angle.Y = PLAYER_CAMERA_MAX_CLAMP * rl.Deg2rad
	}

	//Recalculate camera target based on translation and rotation
	target := rl.Vector3Transform(rl.NewVector3(0, 0, 1), rl.MatrixRotateXYZ(rl.NewVector3(c.angle.Y, -c.angle.X, 0)))

	if c.allow_flight {
		c.forward = target
	} else {
		c.forward = rl.Vector3Transform(rl.NewVector3(0, 0, 1), rl.MatrixRotateXYZ(rl.NewVector3(0, -c.angle.X, 0)))
	}

	c.right = rl.NewVector3(c.forward.Z*-1.0, 0, c.forward.X)

	camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(c.forward, direction["forward"]-direction["backward"]))
	camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(c.right, direction["right"]-direction["left"]))
	camera.Position.Y = direction["up"] - direction["down"]

	var eye_offset float32 = c.player_eyes_pos

	swing_delta := math32.Max(
		math32.Abs(direction["forward"]-direction["backward"]),
		math32.Abs(direction["right"]-direction["left"]))

	//Head bobble calculation
	if swing_delta > 0 {

		c.curr_bobble += swing_delta * c.view_bobble_freq

		var view_bobble_dampen float32 = 8.0

		eye_offset -= math32.Sin(c.curr_bobble/view_bobble_dampen) * c.view_bobble_mag

		camera.Up.X = math32.Sin(c.curr_bobble/(view_bobble_dampen*2)) * c.view_bobble_waver_mag
		camera.Up.Z = -math32.Sin(c.curr_bobble/(view_bobble_dampen*2)) * c.view_bobble_waver_mag

	} else {
		c.curr_bobble = 0
		camera.Up.X = 0
		camera.Up.Z = 0
	}

	camera.Position.Y += eye_offset

	// Move camera according to matrix position
	camera.Target.X = camera.Position.X + target.X
	camera.Target.Y = camera.Position.Y + target.Y
	camera.Target.Z = camera.Position.Z + target.Z

	//camera.Position.Y = c.player_eyes_pos
	return swing_delta

}
