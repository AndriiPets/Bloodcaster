package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Utils_MakeBoundingBox(position, size rl.Vector3) rl.BoundingBox {

	bb := rl.BoundingBox{
		Min: rl.NewVector3(position.X-size.X/2, position.Y-size.Y/2, position.Z-size.Z/2),
		Max: rl.NewVector3(position.X+size.X/2, position.Y+size.Y/2, position.Z+size.Z/2)}

	return bb
}

func Utils_calculate_camera_vector(yaw float32, pitch float32) rl.Vector3 {
	c_yaw := float32(math.Cos(float64(yaw)))
	c_pitch := float32(math.Cos(float64(pitch)))
	s_yaw := float32(math.Sin(float64(yaw)))
	s_pitch := float32(math.Sin(float64(pitch)))

	x := c_yaw * c_pitch
	y := s_pitch
	z := s_yaw * c_pitch
	direction := rl.NewVector3(x, y, z)

	camera_front_vec := rl.Vector3Normalize(direction)
	return camera_front_vec
}

func Utils_find_grid_cell_center(index_x, index_y, len float32) rl.Vector2 {
	min_x := index_x * len
	min_y := index_y * len

	center_x := min_x + len/2
	center_y := min_y + len/2
	cell_center := rl.NewVector2(center_x, center_y)
	return cell_center
}

func Utils_find_cell_index(x, y, len float32) rl.Vector2 {
	index_x := x - float32(math.Mod(float64(x), float64(len)))
	index_y := y - float32(math.Mod(float64(y), float64(len)))

	return rl.NewVector2(index_x, index_y)
}
