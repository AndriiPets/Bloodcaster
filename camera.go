package main

import (
	"maze.io/x/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerCustomCamera struct {
	target_distance   float32
	player_eyes_pos   float32
	angle             rl.Vector2
	mouse_sensitivity float32
	player_speed      float32
	forward           rl.Vector3
	right             rl.Vector3
	allow_flight      bool
	//view bobble stuff
	view_bobble_freq      float32
	view_bobble_mag       float32
	curr_bobble           float32
	view_bobble_waver_mag float32
}

func (c *PlayerCustomCamera) custom_camera_init() {
	c.target_distance = 0
	c.player_eyes_pos = 0.4
	c.mouse_sensitivity = 0.003
	c.player_speed = 1.75
	c.allow_flight = false
	c.view_bobble_freq = 0.5
	c.view_bobble_mag = 0.02
	c.curr_bobble = 0
	c.view_bobble_waver_mag = 0.002
}

func (c *PlayerCustomCamera) Player_camera_init(pos_x, pos_z float32) rl.Camera {

	camera := rl.Camera{}
	c.custom_camera_init()

	camera.Position = rl.NewVector3(pos_x, 0.4, pos_z)
	camera.Target = rl.NewVector3(0.0, 0.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = FOV
	camera.Projection = rl.CameraPerspective

	//Distance
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
	c.angle.Y = math32.Atan2(dy, math32.Sqrt(dx*dx+dz*dz))

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

func (c *PlayerCustomCamera) Get_speed_for_axis(axis int32) float32 {
	factor := 1.0
	if rl.IsKeyDown(KEY_SHIFT) {
		factor = 2.0
	}
	if rl.IsKeyDown(axis) {
		return c.player_speed * rl.GetFrameTime() * float32(factor)
	}
	return 0
}

func (c *PlayerCustomCamera) Player_update_camera(camera *rl.Camera) {

	mouse_pos_delta := rl.GetMouseDelta()
	direction := map[string]float32{
		"forward":  c.Get_speed_for_axis(KEY_MOVE_FORWARD),
		"backward": c.Get_speed_for_axis(KEY_MOVE_BACKWARD),
		"left":     c.Get_speed_for_axis(KEY_MOVE_LEFT),
		"right":    c.Get_speed_for_axis(KEY_MOVE_RIGHT),
		"up":       c.Get_speed_for_axis(KEY_MOVE_UP),
		"down":     c.Get_speed_for_axis(KEY_MOVE_DOWN),
	}

	/* //Move camera around x pos
	camera.Position.X += ((math32.Sin(c.angle.X)*direction["backward"] -
		math32.Sin(c.angle.X)*direction["forward"] -
		math32.Cos(c.angle.X)*direction["left"] +
		math32.Cos(c.angle.X)*direction["right"]) * c.player_speed) * rl.GetFrameTime()

	//Move camera around y pos
	camera.Position.Y += (math32.Sin(c.angle.Y)*direction["forward"] -
		math32.Sin(c.angle.Y)*direction["backward"]*c.player_speed) * rl.GetFrameTime()

	//Move camera around z pos
	camera.Position.Z += (math32.Cos(c.angle.X)*direction["backward"] -
		math32.Cos(c.angle.X)*direction["forward"] +
		math32.Sin(c.angle.X)*direction["left"] -
		math32.Sin(c.angle.X)*direction["right"]*c.player_speed) * rl.GetFrameTime() */

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
	//translation := rl.MatrixTranslate(0, 0, (c.target_distance / PLAYER_CAMERA_PANNING_DIVIDER))

	//rotation := rl.MatrixInvert(rl.MatrixRotateXYZ(rl.NewVector3(math32.Pi*2-c.angle.Y, math32.Pi*2-c.angle.X, 0)))

	//transform := rl.MatrixMultiply(translation, rotation)
	target := rl.Vector3Transform(rl.NewVector3(0, 0, 1), rl.MatrixRotateXYZ(rl.NewVector3(c.angle.Y, -c.angle.X, 0)))

	if c.allow_flight {
		c.forward = target
	} else {
		c.forward = rl.Vector3Transform(rl.NewVector3(0, 0, 1), rl.MatrixRotateXYZ(rl.NewVector3(0, -c.angle.X, 0)))
	}

	c.right = rl.NewVector3(c.forward.Z*-1.0, 0, c.forward.X)

	camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(c.forward, direction["forward"]-direction["backward"]))
	camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(c.right, direction["right"]-direction["left"]))
	camera.Position.Y = +direction["up"] - direction["down"]

	var eye_offset float32 = c.player_eyes_pos

	if c.view_bobble_freq > 0 {
		swing_delta := math32.Max(
			math32.Abs(direction["forward"]-direction["backward"]),
			math32.Abs(direction["right"]-direction["left"]))

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

}
