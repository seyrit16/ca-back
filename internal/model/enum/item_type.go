package enum

type ItemLocationType int16

const (
	ItemLocationWarehouse   ItemLocationType = iota // склад игрока
	ItemLocationPocket                              // карман юнита
	ItemLocationEquipHead                           // голова
	ItemLocationEquipBody                           // тело
	ItemLocationEquipBoots                          // ботинки
	ItemLocationEquipWeapon                         // оружие
)
