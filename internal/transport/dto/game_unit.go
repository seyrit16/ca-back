package dto

import "ca/internal/model"

type GameUnitDTO struct {
	PlayerId string `json:"player_id"`
	UnitId   string `json:"unit_id"`

	Seat      int16 `json:"seat"`
	PositionQ int16 `json:"position_q"`
	PositionR int16 `json:"position_r"`

	CharacterTemplateId int16  `json:"character_template_id"`
	Name                string `json:"name"`

	Hp      int16 `json:"hp"`
	MaxHp   int16 `json:"max_hp"`
	Attack  int16 `json:"attack"`
	Defense int16 `json:"defense"`
	Luck    int16 `json:"luck"`
	Speed   int16 `json:"speed"`

	CriticalChance int16 `json:"critical_chance"`
	CriticalDamage int16 `json:"critical_damage"`

	Accuracy int16 `json:"accuracy"`
	Dodge    int16 `json:"dodge"`

	EquipHead   *int16 `json:"equip_head"`
	EquipBody   *int16 `json:"equip_body"`
	EquipBoots  *int16 `json:"equip_boots"`
	EquipWeapon *int16 `json:"equip_weapon"`

	IsAlive bool `json:"is_alive"`
}

func ToUnitDTO(unit *model.GameUnit) *GameUnitDTO {
	return &GameUnitDTO{
		PlayerId:            unit.PlayerId,
		UnitId:              unit.UnitId,
		Seat:                unit.Seat,
		PositionQ:           unit.PositionQ,
		PositionR:           unit.PositionR,
		CharacterTemplateId: unit.CharacterTemplateId,
		Name:                unit.Name,
		Hp:                  unit.Hp,
		MaxHp:               unit.MaxHp,
		Attack:              unit.Attack,
		Defense:             unit.Defense,
		Luck:                unit.Luck,
		Speed:               unit.Speed,
		CriticalChance:      unit.CriticalChance,
		CriticalDamage:      unit.CriticalDamage,
		Accuracy:            unit.Accuracy,
		Dodge:               unit.Dodge,
		EquipHead:           unit.EquipHead,
		EquipBody:           unit.EquipBody,
		EquipBoots:          unit.EquipBoots,
		EquipWeapon:         unit.EquipWeapon,
		IsAlive:             unit.IsAlive,
	}
}

func UnitsToDTO(units map[string]*model.GameUnit) map[string]*GameUnitDTO {
	result := make(map[string]*GameUnitDTO)

	for key, unit := range units {
		result[key] = ToUnitDTO(unit)
	}
	return result
}
