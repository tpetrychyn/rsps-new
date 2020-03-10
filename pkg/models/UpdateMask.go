package models

type UpdateMask struct {
	Hitmark        bool
	Graphic        bool
	NeedsPlacement bool
	ForcedMovement bool
	ForcedChat     bool
	FaceTile       bool
	Appearance     bool
	FaceActor      bool
	PublicChat     bool
	Animation      bool
}

func (u *UpdateMask) UpdateRequired() uint {
	if u.Hitmark || u.Graphic || u.NeedsPlacement || u.ForcedMovement || u.ForcedChat || u.FaceTile || u.Appearance	|| u.FaceActor || u.PublicChat || u.Animation {
		return 1
	}
	return 0
}

func (u *UpdateMask) Clear() {
	u.Hitmark = false
	u.Graphic = false
	u.NeedsPlacement = false

	u.Appearance = false
	u.Animation = false
}