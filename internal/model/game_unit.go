package model

import "github.com/google/uuid"

type GameUnit struct {
	// PRIMARY KEY (game_id, player_id, unit_id)
	GameId uuid.UUID `db:"game_id"`
	// PRIMARY KEY (game_id, player_id, unit_id)
	PlayerId string `db:"player_id"`
	// PRIMARY KEY (game_id, player_id, unit_id)
	UnitId string `db:"unit_id"` // unitId внутри игры "unit_0", "unit_1"

	Seat      int16 `db:"seat"`
	PositionQ int16 `db:"position_q"`
	PositionR int16 `db:"position_r"`

	CharacterTemplateId int16 `db:"character_template_id"`
	// ==================================
	// template stats
	Name string `db:"name"`

	// main stats
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

	// equip
	EquipHead   *int16 `db:"equip_head"`   //item_id или nil
	EquipBody   *int16 `db:"equip_body"`   //item_id или nil
	EquipBoots  *int16 `db:"equip_boots"`  //item_id или nil
	EquipWeapon *int16 `db:"equip_weapon"` //item_id или nil

	IsAlive bool `db:"is_alive"`

	Pocket []*GameItem `db:"-"`
}
