package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Entity struct {
	role        string
	position    rl.Vector2
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
	planes         map[string]rl.Model
	player_pos     rl.Vector2
	hit_mark       []rl.Vector3
}

func (m *Map) get_game_map() [][]rune {
	game_map := [][]rune{
		{'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c', 'c'},
		{'c', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', 'c', '_', '_', 'c'},
		{'c', '_', '_', 'c', 'c', 'c', 'c', '_', '_', '_', 'c', 'c', 'c', '_', 'c', '_', '_', 'c'},
		{'c', '_', 'p', '_', '_', '_', 'c', '_', '_', '_', '_', '_', 'c', '_', '_', '_', '_', 'c'},
		{'c', '_', '_', '_', '_', '_', 'c', '_', '_', '_', 'c', '_', 'c', '_', 'c', 'c', '_', 'c'},
		{'c', '_', '_', 'c', 'c', 'c', 'c', '_', '_', '_', 'c', '_', 'c', '_', 'c', '_', '_', 'c'},
		{'c', '_', '_', '_', '_', '_', '_', '_', '_', 'c', 'c', 'c', 'c', '_', 'c', 'c', '_', 'c'},
		{'c', '_', '_', 'c', '_', '_', '_', 'c', '_', '_', 'c', '_', '_', '_', '_', '_', '_', 'c'},
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

	return game_map
}

func (m *Map) Map_parse_regions() {

	game_map := m.get_game_map()
	c_map := make(map[int][]Entity)

	for y, row := range game_map {
		for x, val := range row {
			if val == 'c' {
				//draw walls and create the wall entity
				//position := rl.NewVector2(float32(m.map_position.X-0.5+float32(x)), float32(m.map_position.Z-0.5+float32(y)))
				position := rl.NewVector2(float32(m.map_position.X+0.5+float32(x)), float32(m.map_position.Z+0.5+float32(y)))
				bb := Utils_MakeBoundingBox(rl.NewVector3(position.X, WALLS_HEIGHT, position.Y), rl.NewVector3(1.0, 1.0, 1.0))
				model := m.Create_wall_model()

				map_entity := Entity{role: "wall", position: position, boundingBox: bb, model: model}
				m.walls = append(m.walls, map_entity)

				//calculate to wich cell object belongs and place it in cell_map
				grid_cell := Utils_find_cell_index(position.X, position.Y, MAP_CELL_SIZE)

				num := m.cell_index_map[grid_cell]

				c_map[num] = append(c_map[num], map_entity)
			} else if val == 'p' {
				m.player_pos = rl.NewVector2(float32(m.map_position.X-0.5+float32(x)), float32(m.map_position.Z-0.5+float32(y)))
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
	p_map := make(map[string]rl.Model)

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

	p_map["floor"] = floor_plane
	p_map["ceiling"] = ceiling_plane
	m.planes = p_map

}

func (m Map) Unload() {
	for _, plane := range m.planes {
		rl.UnloadShader(plane.Materials.Shader)
		rl.UnloadModel(plane)
	}

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

func (m *Map) Update() {

	plane, ok := m.planes["floor"]
	if !ok {
		rl.DrawModel(rl.LoadModelFromMesh(rl.GenMeshPlane(m.width, m.height, 1, 1)), m.map_position, 1, rl.White)
	} else {
		rl.DrawModel(plane, rl.NewVector3(-m.width/2, 0.0, -m.height/2), 1, rl.White)
		rl.DrawModelEx(
			m.planes["ceiling"],
			rl.NewVector3(-m.width/2, 1.0, -m.height/2),
			rl.NewVector3(-1, 0, 0), 180.0,
			rl.NewVector3(1, 1, 1),
			rl.White)
	}

	for _, wall := range m.walls {
		//rl.DrawCubeWires(rl.NewVector3(wall.position.X, WALLS_HEIGHT, wall.position.Y), 1.0, 1.0, 1.0, rl.Blue)
		rl.DrawModel(wall.model, rl.NewVector3(wall.position.X, WALLS_HEIGHT, wall.position.Y), 1.0, rl.White)
	}

	for _, point := range m.hit_mark {
		rl.SetWindowTitle(fmt.Sprintln(point))
		rl.DrawSphere(point, 0.1, rl.Black)
	}
}
func (m *Map) Check_wall_collision(entity_pos rl.Vector2, entity_size float32) bool {

	entity_inx := Utils_find_cell_index(entity_pos.X, entity_pos.Y, 3.0)
	entity_cell := m.cell_index_map[entity_inx]

	entity_bb := Utils_MakeBoundingBox(
		rl.NewVector3(entity_pos.X, 0.0, entity_pos.Y),
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
		hit_level := rl.GetRayCollisionMesh(ray, wall.model.GetMeshes()[0], rl.MatrixTranslate(pos.X, WALLS_HEIGHT, pos.Y))
		if hit_level.Hit {

			if hit_level.Distance < distance {
				fmt.Println("distance", hit_level.Distance)
				distance = hit_level.Distance
				if len(m.hit_mark) < 5 {
					m.hit_mark = append(m.hit_mark, hit_level.Point)
				} else {
					m.hit_mark = m.hit_mark[1:] //pop
					m.hit_mark = append(m.hit_mark, hit_level.Point)
				}
			}
		}
	}
}
