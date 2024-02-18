package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Entity struct {
	role        string
	position    rl.Vector3
	boundingBox rl.BoundingBox
	model       rl.Model
}

type Map struct {
	name           string
	width          float32
	height         float32
	walls          []Entity
	cell_map       map[int][]Entity
	cell_index_map map[rl.Vector2]int
	map_position   rl.Vector3
	player_pos     rl.Vector2
	hit_mark       []rl.Vector3
	actors         []Actor
	alpha_shader   rl.Shader
}

func (m *Map) get_game_map() [][]rune {
	game_map := [][]rune{
		{'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c'},
		{'c', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', 'c', '_', '_', 'c'},
		{'c', '_', '_', 'c', 'c', 'c', 'c', '_', '_', 'm', 'c', 'c', 'c', '_', 'c', '_', '_', 'c'},
		{'c', '_', 'p', '_', '_', '_', 'c', '_', '_', '_', '_', '_', 'c', '_', '_', '_', '_', 'c'},
		{'c', '_', '_', '_', '_', '_', 'c', '_', '_', '_', 'c', '_', 'c', '_', 'c', 'c', '_', 'c'},
		{'c', '_', '_', 'c', 'c', 'c', 'c', '_', '_', '_', 'c', '_', 'c', '_', 'c', '_', '_', 'c'},
		{'c', '_', '_', '_', '_', '_', '_', '_', '_', 'c', 'c', 'c', 'c', 'm', 'c', 'c', '_', 'c'},
		{'c', '_', '_', 'c', '_', '_', 'm', 'c', '_', '_', 'c', '_', '_', '_', '_', '_', '_', 'c'},
		{'c', '_', '_', 'c', '_', '_', '_', 'c', '_', '_', 'c', '_', '_', '_', '_', '_', '_', 'c'},
		{'c', '_', '_', 'c', '_', '_', '_', 'c', '_', '_', '_', '_', '_', '_', '_', '_', 'c', 'c'},
		{'c', '_', '_', 'c', '_', '_', '_', 'c', '_', '_', '_', '_', '_', '_', 'c', '_', '_', 'c'},
		{'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c'},
	}
	//18x12

	m.height = float32(len(game_map))
	m.width = float32(len(game_map[0]))

	return game_map
}

func Map_init() Map {

	game_map := Map{}

	game_map.map_position = rl.NewVector3(-18.0, 0.0, -12.0)
	game_map.Map_cell_partition()
	game_map.Map_parse_regions()
	game_map.Create_planes()
	game_map.alpha_shader = rl.LoadShader("", "./assets/shaders/discard_alpha.fs")

	return game_map
}

func (m *Map) Map_parse_regions() {

	game_map := m.get_game_map()
	c_map := make(map[int][]Entity)
	model := m.Create_wall_model()

	for y, row := range game_map {
		for x, val := range row {
			position := rl.NewVector3(float32(m.map_position.X+0.5+float32(x)), WALLS_HEIGHT, float32(m.map_position.Z+0.5+float32(y)))
			if val == 'c' {
				//draw walls and create the wall entity
				//position := rl.NewVector2(float32(m.map_position.X-0.5+float32(x)), float32(m.map_position.Z-0.5+float32(y)))

				bb := Utils_MakeBoundingBox(rl.NewVector3(position.X, 0, position.Z), rl.NewVector3(1.0, 1.0, 1.0))

				map_entity := Entity{role: "wall", position: position, boundingBox: bb, model: model}
				m.walls = append(m.walls, map_entity)

				//calculate to wich cell object belongs and place it in cell_map
				grid_cell := Utils_find_cell_index(position.X, position.Z, MAP_CELL_SIZE)

				num := m.cell_index_map[grid_cell]

				c_map[num] = append(c_map[num], map_entity)
			} else if val == 'p' {
				m.player_pos = rl.NewVector2(position.X, position.Z)
			} else if val == 'm' {
				m.actors = append(m.actors, New_actor("bill", position))
			}
		}
	}
	m.cell_map = c_map
}

func (m *Map) Create_wall_model() rl.Model {

	texture := Load_texture("./assets/textures/wall1.png")

	cube := rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))

	//rl.SetMaterialTexture(cube.Materials, rl.MapDiffuse, texture)
	rl.SetMaterialTexture(cube.Materials, rl.MapDiffuse, texture)

	return cube

}

func (m *Map) Create_planes() {

	//generate plane models
	floor_plane := rl.LoadModelFromMesh(rl.GenMeshPlane(m.width, m.height, 1, 1))
	ceiling_plane := rl.LoadModelFromMesh(rl.GenMeshPlane(m.width, m.height, 1, 1))

	texture := Load_texture("./assets/textures/floor1.png")
	texture_ceiling := Load_texture("./assets/textures/ceiling1.png")
	rl.SetMaterialTexture(floor_plane.Materials, rl.MapDiffuse, texture)
	rl.SetMaterialTexture(ceiling_plane.Materials, rl.MapDiffuse, texture_ceiling)

	//texture tiling
	tiling := []float32{m.width / 2, m.height / 2}
	shader := rl.LoadShader("", "./assets/shaders/tiling.fs")
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "tiling"), tiling, rl.ShaderUniformVec2)

	floor_plane.Materials.Shader = shader
	ceiling_plane.Materials.Shader = shader

	floor_pos := rl.NewVector3(-m.width/2, 0.0, -m.height/2)
	floor_bb := Utils_MakeBoundingBox(rl.NewVector3(floor_pos.X, 0, floor_pos.Z), rl.NewVector3(m.width, 0, m.height))

	ceiling_pos := rl.NewVector3(-m.width/2, 1.0, -m.height/2)
	ceiling_bb := Utils_MakeBoundingBox(rl.NewVector3(ceiling_pos.Z, 0, ceiling_pos.Z), rl.NewVector3(m.width, 0, m.height))

	floor_entity := Entity{role: "floor", position: floor_pos, boundingBox: floor_bb, model: floor_plane}
	ceiling_entity := Entity{role: "ceiling", position: ceiling_pos, boundingBox: ceiling_bb, model: ceiling_plane}

	m.walls = append(m.walls, floor_entity, ceiling_entity)

}

func (m Map) Unload() {

	for _, wall := range m.walls {
		rl.UnloadTexture(wall.model.Materials.Maps.Texture)
		rl.UnloadModel(wall.model)

	}
}

func (m *Map) Map_cell_partition() {
	//goes thru the map and figures out to what cell each point belongs and place it in the cell_index_map
	i_map := make(map[rl.Vector2]int)
	cell_num := 0
	game_map := m.get_game_map()
	mapPosition := m.map_position

	for y, row := range game_map {
		for x := range row {

			position := rl.NewVector2(float32(mapPosition.X-0.5+float32(x)), float32(mapPosition.Z-0.5+float32(y)))

			grid_cell := Utils_find_cell_index(position.X, position.Y, MAP_CELL_SIZE)

			_, ok := i_map[grid_cell]
			if !ok {
				i_map[grid_cell] = cell_num
				cell_num += 1
			}

		}
	}
	m.cell_index_map = i_map

}

func (m *Map) Update(camera rl.Camera) {

	for _, wall := range m.walls {

		switch wall.role {
		case "wall":
			//rl.DrawCubeWires(rl.NewVector3(wall.position.X, WALLS_HEIGHT, wall.position.Y), 1.0, 1.0, 1.0, rl.Blue)
			rl.DrawModel(wall.model, wall.position, 1.0, rl.White)
		case "floor":
			rl.DrawModel(wall.model, wall.position, 1, rl.White)
		case "ceiling":
			rl.DrawModelEx(
				wall.model,
				wall.position,
				rl.NewVector3(-1, 0, 0), 180.0,
				rl.NewVector3(1, 1, 1),
				rl.White)
		}
	}

	rl.BeginShaderMode(m.alpha_shader)
	for _, actor := range m.actors {
		actor.Draw(camera)
	}
	rl.EndShaderMode()

	for _, point := range m.hit_mark {
		rl.SetWindowTitle(fmt.Sprintln(point))
		rl.DrawSphere(point, 0.05, rl.Black)
	}
}
func (m *Map) Check_wall_collision(entity_pos rl.Vector3, entity_size float32) bool {

	entity_inx := Utils_find_cell_index(entity_pos.X, entity_pos.Z, 3.0)
	entity_cell := m.cell_index_map[entity_inx]

	entity_bb := Utils_MakeBoundingBox(
		rl.NewVector3(entity_pos.X, 0, entity_pos.Z),
		rl.NewVector3(entity_size, entity_size, entity_size))

	for _, wall := range m.cell_map[entity_cell] {
		if rl.CheckCollisionBoxes(entity_bb, wall.boundingBox) {

			return true
		}
	}
	return false
}

func (m *Map) Test_wall_hit(ray rl.Ray) {

	var distance float32 = 1000

	for _, wall := range m.walls {
		pos := wall.position
		hit_level := rl.GetRayCollisionMesh(ray, wall.model.GetMeshes()[0], rl.MatrixTranslate(pos.X, pos.Y, pos.Z))
		if hit_level.Hit {

			if hit_level.Distance < distance {
				fmt.Println("distance", hit_level.Distance)
				distance = hit_level.Distance
				if len(m.hit_mark) < 10 {
					m.hit_mark = append(m.hit_mark, hit_level.Point)
				} else {
					m.hit_mark = m.hit_mark[1:] //pop
					m.hit_mark = append(m.hit_mark, hit_level.Point)
				}
			}
		}
	}
}
