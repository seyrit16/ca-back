package mapgen

import (
	"fmt"
	"math"
	"sort"
	"time"

	. "ca/internal/model/enum"
)

type GenerationItem struct {
	Type        BiomeType
	Coordinates []Coordinate
	Content     ContentType

	Radius int16
	Center Coordinate
}

type SeededRandom struct {
	seed uint32
}

func NewSeededRandom(seed uint32) *SeededRandom {
	if seed == 0 {
		seed = 1
	}
	return &SeededRandom{seed: seed}
}

func (r *SeededRandom) nextUint32() uint32 {
	r.seed = uint32(1664525)*r.seed + 1013904223
	return r.seed
}

func (r *SeededRandom) Next(minInclusive, maxExclusive int) int {
	if maxExclusive <= minInclusive {
		return minInclusive
	}
	rang := maxExclusive - minInclusive
	return minInclusive + int(r.nextUint32())%rang
}

type MapGenerator struct {
	// кол-во клеток = 3 * r² + 3 * r + 1
	Radius      int
	SpawnRadius int
	Settings    GenerationSettings

	Tiles  map[[2]int16]*TileData
	Coords map[[2]int16]Coordinate
	Items  []GenerationItem
	rng    *SeededRandom
}

func NewMapGenerator(radius, spawnRadius int) *MapGenerator {
	return &MapGenerator{
		Radius:      radius,
		SpawnRadius: spawnRadius,
		Settings:    DefaultGenerationSettings(),
	}
}

func (g *MapGenerator) Generate() {
	g.Tiles = make(map[[2]int16]*TileData)
	g.Coords = make(map[[2]int16]Coordinate)
	g.Items = []GenerationItem{}
	g.rng = NewSeededRandom(uint32(time.Now().UnixNano() % 1000))

	for q := -g.Radius; q <= g.Radius; q++ {
		rMin := max(-g.Radius, -q-g.Radius)
		rMax := min(g.Radius, -q+g.Radius)
		for r := rMin; r <= rMax; r++ {
			c := Coordinate{Q: int16(q), R: int16(r)}
			g.Tiles[c.Key()] = &TileData{}
			g.Coords[c.Key()] = c
		}
	}

	countWater := g.rng.Next(g.Settings.WaterCount.Min, g.Settings.WaterCount.Max+1)
	countMountains := g.rng.Next(g.Settings.MountainsCount.Min, g.Settings.MountainsCount.Max+1)

	g.generateLandscape(BiomeTypeWater, countWater, g.Settings.WaterSize.Min, g.Settings.WaterSize.Max)
	g.generateLandscape(BiomeTypeMountain, countMountains, g.Settings.MountainsSize.Min, g.Settings.MountainsSize.Max)

	g.generateContent(ContentTypeTree, g.Settings.TreesCount)
	g.generateContent(ContentTypeSecret, g.Settings.ChestsCount)
	g.generateMonsters(g.Settings.MonstersCount, [2]int16{3, 5})
}

func (g *MapGenerator) GetTile(c Coordinate) *TileData {
	return g.Tiles[c.Key()]
}

func (g *MapGenerator) CheckInRedZone(c Coordinate) bool {
	corner1 := Coordinate{Q: int16(g.Radius), R: int16(-g.Radius)}
	corner2 := Coordinate{Q: int16(-g.Radius), R: int16(g.Radius)}
	return g.hexDistance(c, corner1) <= float64(g.SpawnRadius) ||
		g.hexDistance(c, corner2) <= float64(g.SpawnRadius)
}

func (g *MapGenerator) hexDistance(a, b Coordinate) float64 {
	dq := int(a.Q) - int(b.Q)
	dr := int(a.R) - int(b.R)
	return float64(abs(dq)+abs(dr)+abs(dq+dr)) / 2
}

func (g *MapGenerator) generateLandscape(biome BiomeType, count, sizeMin, sizeMax int) {
	for i := 0; i < count; i++ {
		gi := GenerationItem{Type: biome}
		size := g.rng.Next(sizeMin, sizeMax+1)

		isFirstRun := true
		foundSpot := false
		c := InvalidCoord

		for j := 0; j < size; j++ {
			if isFirstRun || sizeMax < 2 {
				foundSpot = false
				for k := 0; k < 20; k++ {
					candidate := g.randomCoordInMap()
					if g.CheckInRedZone(candidate) {
						continue
					}
					tile := g.GetTile(candidate)
					if tile != nil && tile.Biome == BiomeTypeNone {
						c = candidate
						foundSpot = true
						break
					}
				}

				if !foundSpot {
					break
				}

				gi.Coordinates = append(gi.Coordinates, c)
				g.GetTile(c).Biome = biome
				isFirstRun = false
			}

			if sizeMax < 2 {
				break
			}

			c = g.newDirection(gi.Coordinates)
			if c.IsInvalid() {
				continue
			}

			gi.Coordinates = append(gi.Coordinates, c)
			g.GetTile(c).Biome = biome
		}

		g.Items = append(g.Items, gi)
	}
}

func (g *MapGenerator) newDirection(coords []Coordinate) Coordinate {
	for i := 0; i < len(coords); i++ {
		indexFrom := g.rng.Next(0, len(coords))
		cFrom := coords[indexFrom]
		blocked := [6]bool{}

		for k := 0; k < 12; k++ {
			dir := g.rng.Next(0, 6)
			if blocked[dir] {
				continue
			}

			dq := Directions[dir][0]
			dr := Directions[dir][1]
			newC := Coordinate{Q: cFrom.Q + dq, R: cFrom.R + dr}

			if g.isValidStep(newC) {
				return newC
			}
			blocked[dir] = true
		}
	}

	return InvalidCoord
}

func (g *MapGenerator) isValidStep(c Coordinate) bool {
	tile := g.GetTile(c)
	if tile == nil {
		return false
	}
	if g.CheckInRedZone(c) {
		return false
	}
	return tile.Biome == BiomeTypeNone
}

func (g *MapGenerator) randomCoordInMap() Coordinate {
	q := g.rng.Next(-g.Radius, g.Radius+1)
	rMin := max(-g.Radius, -q-g.Radius)
	rMax := min(g.Radius, -q+g.Radius)
	r := g.rng.Next(rMin, rMax+1)
	return Coordinate{Q: int16(q), R: int16(r)}
}

func (g *MapGenerator) generateContent(contentType ContentType, count int) {
	var candidates []Coordinate
	for _, c := range g.Coords {
		tile := g.GetTile(c)
		if tile.Biome == BiomeTypeNone && tile.Content == ContentTypeNone && !g.CheckInRedZone(c) {
			candidates = append(candidates, c)
		}
	}

	if len(candidates) == 0 {
		return
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Q != candidates[j].Q {
			return candidates[i].Q < candidates[j].Q
		}
		return candidates[i].R < candidates[j].R
	})

	actualCount := min(count, len(candidates))
	step := len(candidates) / actualCount

	for i := 0; i < actualCount; i++ {
		bucketStart := i * step
		bucketEnd := bucketStart + step
		if i == actualCount-1 {
			bucketEnd = len(candidates)
		}
		pick := g.rng.Next(bucketStart, bucketEnd)
		g.GetTile(candidates[pick]).Content = contentType
	}
}

func (g *MapGenerator) generateMonsters(count int, radiusRange [2]int16) {
	for i := 0; i < count; i++ {
		var candidates []Coordinate
		for _, c := range g.Coords {
			tile := g.GetTile(c)
			if tile.Biome != BiomeTypeNone || tile.Content != ContentTypeNone {
				continue
			}
			if g.CheckInRedZone(c) {
				continue
			}
			candidates = append(candidates, c)
		}

		if len(candidates) == 0 {
			break
		}

		sort.Slice(candidates, func(a, b int) bool {
			if candidates[a].Q != candidates[b].Q {
				return candidates[a].Q < candidates[b].Q
			}
			return candidates[a].R < candidates[b].R
		})

		step := len(candidates) / count
		if step == 0 {
			step = 1
		}
		bucketStart := i * step
		bucketEnd := bucketStart + step
		if bucketEnd > len(candidates) {
			bucketEnd = len(candidates)
		}
		if bucketStart >= len(candidates) {
			bucketStart = 0
			bucketEnd = len(candidates)
		}

		radius := radiusRange[0] + int16(g.rng.Next(0, int(radiusRange[1]-radiusRange[0]+1)))

		placed := false
		for attempt := bucketStart; attempt < bucketEnd; attempt++ {
			pick := g.rng.Next(bucketStart, bucketEnd)
			c := candidates[pick]

			validCenter := true
			for _, gi2 := range g.Items {
				if gi2.Radius == 0 {
					continue
				}
				dist := g.hexDistance(c, gi2.Center)
				minDist := float64(gi2.Radius + radius)
				if dist < minDist {
					validCenter = false
					break
				}
			}

			if !validCenter {
				continue
			}

			groupIndex := int16(len(g.Items) + 1)
			tile := g.GetTile(c)
			tile.Content = ContentTypeMonster
			tile.Group = groupIndex
			tile.RadiusContent = radius

			gi := GenerationItem{Type: BiomeTypeNone, Radius: radius, Center: c, Content: ContentTypeMonster}

			for _, coord := range g.Coords {
				if g.hexDistance(c, coord) <= float64(radius) {
					gi.Coordinates = append(gi.Coordinates, coord)
					zoneTile := g.GetTile(coord)
					if zoneTile != nil {
						zoneTile.Group = groupIndex
					}
				}
			}

			g.Items = append(g.Items, gi)
			placed = true
			break
		}

		if !placed {
			// не нашли место в бакете, пробуем весь список
			for _, c := range candidates {
				validCenter := true
				for _, gi2 := range g.Items {
					if gi2.Radius == 0 {
						continue
					}
					dist := g.hexDistance(c, gi2.Center)
					minDist := float64(gi2.Radius + radius)
					if dist < minDist {
						validCenter = false
						break
					}
				}

				if !validCenter {
					continue
				}

				groupIndex := int16(len(g.Items) + 1)
				tile := g.GetTile(c)
				tile.Content = ContentTypeMonster
				tile.Group = groupIndex
				tile.RadiusContent = radius

				gi := GenerationItem{Type: BiomeTypeNone, Radius: radius, Center: c, Content: ContentTypeMonster}

				for _, coord := range g.Coords {
					if g.hexDistance(c, coord) <= float64(radius) {
						gi.Coordinates = append(gi.Coordinates, coord)
						zoneTile := g.GetTile(coord)
						if zoneTile != nil {
							zoneTile.Group = groupIndex
						}
					}
				}

				g.Items = append(g.Items, gi)
				placed = true

				for _, gi2 := range g.Items[:len(g.Items)-1] {
					if gi2.Radius == 0 {
						continue
					}
					dist := g.hexDistance(c, gi2.Center)
					minDist := float64(radius + gi2.Radius)
					if dist < minDist {
						fmt.Printf("OVERLAP: new monster at %v r=%d, existing at %v r=%d, dist=%.1f minDist=%.1f\n",
							c, radius, gi2.Center, gi2.Radius, dist, minDist)
					}
				}

				break
			}
		}
	}
}

func (g *MapGenerator) ToHTML() string {
	if g.Tiles == nil {
		return "Map is empty"
	}

	lairZone := make(map[[2]int16]bool)
	for _, gi := range g.Items {
		if gi.Radius == 0 {
			continue
		}
		for _, coord := range gi.Coordinates {
			dist := g.hexDistance(gi.Center, coord)
			if dist >= float64(gi.Radius)-0.5 && dist <= float64(gi.Radius)+0.5 {
				lairZone[coord.Key()] = true
			}
		}
	}

	result := `<html><body style="background:#111;font-family:monospace;white-space:pre;line-height:1.2">`

	for r := -g.Radius; r <= g.Radius; r++ {
		offset := int(math.Abs(float64(r))) / 2
		line := fmt.Sprintf("%*s", offset*2, "")

		qMin := max(-g.Radius, -r-g.Radius)
		qMax := min(g.Radius, -r+g.Radius)

		for q := qMin; q <= qMax; q++ {
			c := Coordinate{Q: int16(q), R: int16(r)}
			tile := g.GetTile(c)
			if tile != nil {
				line += g.htmlSymbol(tile, lairZone[c.Key()])
			}
		}
		result += line + "\n"
	}

	result += `</body></html>`
	return result
}

func (g *MapGenerator) htmlSymbol(tile *TileData, inLairZone bool) string {
	sym, color := g.colorFor(tile, inLairZone)
	return fmt.Sprintf(`<span style="color:%s">%s </span>`, color, sym)
}

func (g *MapGenerator) colorFor(tile *TileData, inLairZone bool) (string, string) {
	if tile.Content == ContentTypeTree {
		return "T", "#2d8a2d"
	}
	if tile.Content == ContentTypeSecret {
		return "C", "#d4af37"
	}
	if tile.Content == ContentTypeMonster {
		return "M", "#cc2222"
	}
	switch tile.Biome {
	case BiomeTypeWater:
		return "~", "#2255cc"
	case BiomeTypeMountain:
		return "^", "#888888"
	default:
		if inLairZone {
			return "○", "#cc4444"
		}
		return "○", "#ffffff"
	}
}

func (g *MapGenerator) symbolFor(tile *TileData, inLairZone bool) string {
	if tile.Content == ContentTypeTree {
		return "T"
	}
	if tile.Content == ContentTypeSecret {
		return "C"
	}
	if tile.Content == ContentTypeMonster {
		return "M"
	}
	switch tile.Biome {
	case BiomeTypeWater:
		return "~"
	case BiomeTypeMountain:
		return "^"
	default:
		if inLairZone {
			return "."
		}
		return " "
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
