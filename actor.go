package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Actor struct {
	name         string
	position     rl.Vector3
	bounding_box rl.BoundingBox
	texture      rl.Texture2D
}

func New_actor(name string, pos rl.Vector3) Actor {

	actor := Actor{}

	actor.name = name
	actor.position = pos
	actor.bounding_box = Utils_MakeBoundingBox(rl.NewVector3(actor.position.X, 0, actor.position.Z), rl.NewVector3(1, 1, 1))
	actor.texture = Load_texture("./assets/models/billboard.png")
	//actor.alpha_shader = rl.LoadShader("", "./assets/shaders/discard_alpha.fs")

	return actor

}

func (a *Actor) Draw(camera rl.Camera) {

	rl.DrawBillboard(camera, a.texture, a.position, 1, rl.White)

}
