package model

type CharacterTemplate struct {
	Id int16 `db:"id"`

	// ==================================
	// template
	Name string `db:"name"`

	// stats
	Hp      int16 `db:"hp"`
	MaxHp   int16 `db:"max_hp"`
	Attack  int16 `db:"attack"`
	Defense int16 `db:"defense"`
	Luck    int16 `db:"luck"`
	Speed   int16 `db:"speed"`

	// alternative stats
	CriticalChance int16 `db:"critical_chance"`
	CriticalDamage int16 `db:"critical_damage"`

	Accuracy int16 `db:"accuracy"`
	Dodge    int16 `db:"dodge"`
	// ==================================
}
