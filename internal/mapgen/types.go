package mapgen

import . "ca/internal/model/enum"

type TileData struct {
	Biome     BiomeType
	Content   ContentType
	OpenFaces uint8

	Group         int16
	RadiusContent int16
}

func (t *TileData) IsFaceOpen(dir int) bool {
	return (t.OpenFaces & (1 << dir)) != 0
}

func (t *TileData) SetFaceOpen(dir int, open bool) {
	if open {
		t.OpenFaces |= 1 << dir
	} else {
		t.OpenFaces &^= 1 << dir
	}
}
