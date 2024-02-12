package main

import (
	math "maze.io/x/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Utils_MakeBoundingBox(position, size rl.Vector3) rl.BoundingBox {

	bb := rl.BoundingBox{
		Min: rl.NewVector3(position.X-size.X/2, position.Y-size.Y/2, position.Z-size.Z/2),
		Max: rl.NewVector3(position.X+size.X/2, position.Y+size.Y/2, position.Z+size.Z/2)}

	return bb
}

func Utils_calculate_camera_vector(yaw float32, pitch float32) rl.Vector3 {
	c_yaw := math.Cos(yaw)
	c_pitch := math.Cos(pitch)
	s_yaw := math.Sin(yaw)
	s_pitch := math.Sin(pitch)

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
	index_x := x - math.Mod(x, len)
	index_y := y - math.Mod(y, len)

	return rl.NewVector2(index_x, index_y)
}

func Load_texture(path string) rl.Texture2D {
	img := rl.LoadImage(path)
	tex := rl.LoadTextureFromImage(img)

	rl.UnloadImage(img)

	return tex
}
