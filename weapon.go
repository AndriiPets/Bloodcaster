package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"maze.io/x/math32"
)

type Weapon struct {
	name                   string
	input_key              int32
	weapon_id              int32
	damage                 int32
	ammo                   int32
	fire_rate              float32
	fire_range             float32
	picked_up              bool
	sprite_speed           float32
	sprite_fire_frame      int32
	sprites_total          int32
	sprite_texture         rl.Texture2D
	sprite_position_offset rl.Vector2
}

type WeaponHolder struct {
	current_weapon string
	weapons        map[string]*Weapon
	curr_frame     int32
	frame_counter  int32
	is_active      bool
	weapon_pos     rl.Vector2
	//bob
	current_bob float32
	angle       float32
	angle_a     float32
	angle_v     float32
	grav        float32
}

func Weapons_init() WeaponHolder {

	wh := WeaponHolder{}

	wh.weapons = Create_weapons_map()
	wh.load_weapon_textures()
	wh.current_weapon = "pistol"
	wh.weapon_pos = rl.NewVector2(0, 0)
	wh.angle = math32.Pi / 2
	wh.angle_a = 0
	wh.angle_v = 0
	wh.grav = 0.6

	return wh
}

func (w *WeaponHolder) load_weapon_textures() {

	path := "./assets/weapons/"

	for _, weapon := range w.weapons {
		weapon.sprite_texture = Load_texture(path + weapon.name + ".png")
		//fmt.Println(weapon.sprite_texture)
	}
}

func (w *WeaponHolder) Get_switch_input() {

	pressed_key := rl.GetKeyPressed()

	for _, weapon := range w.weapons {
		if pressed_key == weapon.input_key {
			w.weapon_change(weapon.name)
		}
	}
}

func (w *WeaponHolder) weapon_change(name string) {

	if !w.is_active {
		w.current_weapon = name
	}
}

func (w *WeaponHolder) Weapon_fire(camera *rl.Camera, game_map *Map, nextFire float32) float32 {

	if nextFire > 0 {
		nextFire -= rl.GetFrameTime()
	} else {

		weapon := w.weapons[w.current_weapon]

		if !w.is_active {

			w.is_active = true
			w.curr_frame = weapon.sprite_fire_frame

			ray_cast := rl.GetMouseRay(rl.NewVector2(float32(HALF_WIDTH), float32(HALF_HEIGHT)), *camera)
			game_map.Test_wall_hit(ray_cast)
			//rl.DrawRay(ray_cast, rl.Red)
		}
		nextFire = weapon.fire_rate
	}

	return nextFire
}

func (w *WeaponHolder) Draw(swing_delta float32) {

	weapon := w.weapons[w.current_weapon]

	frame_width := float32(weapon.sprite_texture.Width) / float32(weapon.sprites_total)
	frame_height := float32(weapon.sprite_texture.Height)
	origin := rl.NewVector2(frame_width/2, frame_height)

	scale := math32.Min(frame_width*2.0/frame_width, frame_height*2.0/frame_height)
	pos_x := float32(HALF_WIDTH) - (frame_width * weapon.sprite_position_offset.X)
	pos_y := float32(SCREEN_HEIGHT) - (frame_height * weapon.sprite_position_offset.Y)

	//weapon bob
	if swing_delta > 0 {
		var len float32 = 30
		force := w.grav * math32.Sin(w.angle)
		w.angle_a = (-1 * force) / len
		w.angle_v += w.angle_a
		w.angle += w.angle_v

		w.weapon_pos.X = len*math32.Sin(w.angle) + pos_x
		w.weapon_pos.Y = len*math32.Cos(w.angle) + pos_y

	} else {
		w.current_bob = 0
		w.weapon_pos.X = pos_x
		w.weapon_pos.Y = pos_y
	}

	//rl.SetWindowTitle(fmt.Sprintln(rl.GetFPS()))

	source_rect := rl.NewRectangle(0, 0, frame_width, frame_height)
	dest_rect := rl.NewRectangle(w.weapon_pos.X, w.weapon_pos.Y, frame_width*scale, frame_height*scale) //make bobble like camera

	if w.is_active {
		w.frame_counter++

		if w.frame_counter >= rl.GetFPS()/int32(weapon.sprite_speed/weapon.fire_rate) {
			w.curr_frame++

			if w.curr_frame >= weapon.sprites_total {

				w.curr_frame = 0
				w.is_active = false

			}
			w.frame_counter = 0
		}
	}

	source_rect.X = frame_width * float32(w.curr_frame)
	rl.DrawTexturePro(weapon.sprite_texture, source_rect, dest_rect, origin, 0, rl.White)
}
