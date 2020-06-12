package main

type inventory struct {
	Weapon      *item
	Helmet      *item
	Amulet      *item
	Armour      *item
	Boots       *item
	Consumables [NConsumables]*item
}

const NConsumables = 6

type item struct {
	Class     itemClass
	Active    itemActive
	Charge    int
	MaxCharge int
	Passives  []itemPassive
}

type itemClass int

const (
	Weapon itemClass = iota
	Helmet
	Amulet // background class item
	Armour
	Boots
	Consumable
)

type itemActive int

const (
	NoActive    itemActive = iota
	ActiveBlink            // random blink
	// TODO: other actives
)

type itemPassive int

const (
	NoPassive     itemPassive = iota
	PassiveArmour             // basic dmg absorption
	PassiveCleave             // axe cleave effect
	// TODO: other passives
)

func (it *item) ShortDesc(g *game) string {
	return "item" // TODO item short desc
}

func (it *item) Desc(g *game) string {
	return "item" // TODO item desc
}

func (it *item) Style(g *game) (rune, uicolor) {
	return '[', ColorYellow // TODO item style
}
