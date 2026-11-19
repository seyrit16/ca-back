package mapgen

type Range struct {
	Min int
	Max int
}

type GenerationSettings struct {
	WaterCount     Range
	WaterSize      Range
	MountainsCount Range
	MountainsSize  Range

	MonstersCount int
	TreesCount    int
	ChestsCount   int
}

func DefaultGenerationSettings() GenerationSettings {
	return GenerationSettings{
		WaterCount:     Range{Min: 3, Max: 5},
		WaterSize:      Range{Min: 4, Max: 6},
		MountainsCount: Range{Min: 5, Max: 7},
		MountainsSize:  Range{Min: 1, Max: 3},

		MonstersCount: 10,
		TreesCount:    30,
		ChestsCount:   20,
	}
}
