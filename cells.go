package main

import "strings"

type cell struct {
	T        terrain
	Explored bool
}

type terrain int

const (
	WallCell terrain = iota
	GroundCell
	DoorCell
	FoliageCell
	BarrelCell
	StairCell
	StoneCell
	MagaraCell
	BananaCell
	LightCell
	ExtinguishedLightCell
	TableCell
	TreeCell
	HoledWallCell
	ScrollCell
	StoryCell
	ItemCell
	BarrierCell
	WindowCell
	ChasmCell
	WaterCell
	RubbleCell
	CavernCell
	FakeStairCell
	PotionCell
	QueenRockCell
)

func (c cell) IsPassable() bool {
	switch c.T {
	case WallCell, DoorCell, BarrelCell, TableCell, TreeCell, HoledWallCell, BarrierCell, WindowCell, StoryCell, ChasmCell, WaterCell:
		return false
	default:
		return true
	}
}

func (c cell) IsNormalPatrolWay() bool {
	switch c.T {
	case GroundCell, ScrollCell, DoorCell, StairCell, LightCell, ItemCell, ExtinguishedLightCell, StoneCell, MagaraCell, FakeStairCell:
		return true
	default:
		return false
	}
}

func (c cell) IsLevitatePassable() bool {
	switch c.T {
	case ChasmCell:
		return true
	default:
		return c.IsPassable()
	}
}

func (c cell) IsDoorPassable() bool {
	switch c.T {
	case DoorCell:
		return true
	default:
		return c.IsPassable()
	}
}

func (c cell) IsSwimPassable() bool {
	switch c.T {
	case WaterCell:
		return true
	default:
		return c.IsPassable()
	}
}

func (c cell) AllowsFog() bool {
	switch c.T {
	case WallCell, HoledWallCell, WindowCell, StoryCell:
		return false
	default:
		return true
	}
}

func (c cell) CoversPlayer() bool {
	switch c.T {
	case WallCell, BarrelCell, TableCell, TreeCell, HoledWallCell, BarrierCell, WindowCell:
		return true
	default:
		return false
	}
}

func (t terrain) IsPlayerPassable() bool {
	switch t {
	case WallCell, BarrierCell, WindowCell, ChasmCell:
		return false
	default:
		return true
	}
}

func (t terrain) IsDiggable() bool {
	switch t {
	case WallCell, WindowCell, HoledWallCell:
		return true
	default:
		return false
	}
}

func (c cell) BlocksRange() bool {
	switch c.T {
	case WallCell, TreeCell, BarrierCell, WindowCell, StoryCell:
		return true
	default:
		return false
	}
}

func (c cell) Hides() bool {
	switch c.T {
	case WallCell, BarrelCell, TableCell, TreeCell, WindowCell, StoryCell:
		return true
	default:
		return false
	}
}

func (c cell) IsIlluminable() bool {
	switch c.T {
	case WallCell, BarrelCell, TableCell, TreeCell, HoledWallCell, BarrierCell, WindowCell, ChasmCell, RubbleCell:
		return false
	}
	return true
}

func (c cell) IsDestructible() bool {
	switch c.T {
	case WallCell, BarrelCell, DoorCell, TableCell, TreeCell, HoledWallCell, WindowCell:
		return true
	default:
		return false
	}
}

func (c cell) IsWall() bool {
	switch c.T {
	case WallCell:
		return true
	default:
		return false
	}
}

func (c cell) Flammable() bool {
	switch c.T {
	case FoliageCell, DoorCell, BarrelCell, TableCell, TreeCell, WindowCell:
		return true
	default:
		return false
	}
}

func (c cell) IsGround() bool {
	switch c.T {
	case GroundCell, CavernCell, BananaCell, PotionCell, QueenRockCell:
		return true
	default:
		return false
	}
}

func (c cell) IsNotable() bool {
	switch c.T {
	case StairCell, StoneCell, BarrelCell, MagaraCell, BananaCell,
		ScrollCell, ItemCell, FakeStairCell, PotionCell:
		return true
	default:
		return false
	}
}

func (c cell) ShortDesc(g *game, pos position) (desc string) {
	switch c.T {
	case WallCell:
		desc = "a wall"
	case GroundCell:
		desc = "paved ground"
	case DoorCell:
		desc = "a door"
	case FoliageCell:
		desc = "foliage"
	case BarrelCell:
		desc = "a barrel"
	case StoneCell:
		desc = g.Objects.Stones[pos].ShortDesc(g)
	case StairCell:
		desc = g.Objects.Stairs[pos].ShortDesc(g)
	case MagaraCell:
		desc = g.Objects.Magaras[pos].ShortDesc()
	case BananaCell:
		desc = "a banana"
	case LightCell:
		desc = "a campfire"
	case ExtinguishedLightCell:
		desc = "an extinguished campfire"
	case TableCell:
		desc = "a table"
	case TreeCell:
		desc = "a banana tree"
	case HoledWallCell:
		desc = "a holed wall"
	case ScrollCell:
		desc = g.Objects.Scrolls[pos].ShortDesc(g)
	case StoryCell:
		desc = g.Objects.Story[pos].ShortDesc(g, pos)
	case ItemCell:
		desc = g.Objects.Items[pos].ShortDesc(g)
	case BarrierCell:
		desc = "a temporal magical barrier"
	case WindowCell:
		desc = "a window"
	case ChasmCell:
		desc = "a chasm"
	case WaterCell:
		desc = "shallow water"
	case RubbleCell:
		desc = "rubblestone"
	case CavernCell:
		desc = "cave ground"
	case FakeStairCell:
		desc = NormalStairShortDesc
	case PotionCell:
		desc = g.Objects.Potions[pos].ShortDesc(g)
	case QueenRockCell:
		desc = "queen rock"
	}
	return desc
}

func (c cell) Desc(g *game, pos position) (desc string) {
	switch c.T {
	case WallCell:
		desc = "A wall is a pile of rocks."
	case GroundCell:
		desc = "This is paved ground."
	case DoorCell:
		desc = "A closed door blocks your line of sight. Doors open automatically when you or a creature stand on them."
	case FoliageCell:
		desc = "Blue dense foliage grows in Hareka's Underground. It is difficult to see through."
	case BarrelCell:
		desc = "A barrel. You can hide yourself inside it when no creatures see you. It is a safe place for resting and recovering."
	case StoneCell:
		desc = g.Objects.Stones[pos].Desc(g)
	case StairCell:
		desc = g.Objects.Stairs[pos].Desc(g)
	case MagaraCell:
		desc = g.Objects.Magaras[pos].Desc(g)
	case BananaCell:
		desc = "A gawalt monkey cannot enter a healthy sleep without eating one of those bananas before."
	case LightCell:
		desc = "A campfire illuminates surrounding cells. Creatures can spot you in illuminated cells from a greater range."
	case ExtinguishedLightCell:
		desc = "An extinguished campfire can be lighted again by some creatures."
	case TableCell:
		desc = "You can hide under the table so that only adjacent creatures can see you. Most creatures cannot walk accross the table."
	case TreeCell:
		desc = "Underground banana trees grow with nearly no light sources. Their rare bananas are very appreciated by many creatures, specially some harpy species. You may find some bananas dropped by them while exploring. You can climb trees to see farther. Moreover, only big, flying or jumping creatures will be able to attack you while you stand on a tree. The top is never illuminated."
	case HoledWallCell:
		desc = "Only very small creatures can pass there. It is difficult to see through."
	case ScrollCell:
		desc = g.Objects.Scrolls[pos].Desc(g)
	case StoryCell:
		desc = g.Objects.Story[pos].Desc(g, pos)
	case ItemCell:
		desc = g.Objects.Items[pos].Desc(g)
	case BarrierCell:
		desc = "A temporal magical barrier created by oric energies. It may have been created by an oric magara or an oric celmist. Sometimes, natural oric energies may produce such barriers too in energetically unstable Underground areas."
	case WindowCell:
		desc = "A transparent window in the wall."
	case ChasmCell:
		desc = "A chasm. If you jump into it, you'll be seriously injured."
	case WaterCell:
		desc = "Shallow water."
	case RubbleCell:
		desc = "Rubblestone is a collection of rocks broken into smaller stones."
	case CavernCell:
		desc = "This is natural cave ground."
	case FakeStairCell:
		desc = NormalStairDesc
	case PotionCell:
		desc = g.Objects.Potions[pos].Desc(g)
	case QueenRockCell:
		desc = "Queen rock amplifies sounds. Even though you are usually very silent, monsters may hear your footsteps when walking on those rocks."
	}
	var autodesc string
	if !c.T.IsPlayerPassable() {
		autodesc += " It is impassable."
	}
	if c.Flammable() {
		autodesc += " It is flammable."
	}
	if c.IsLevitatePassable() && !c.IsPassable() {
		autodesc += " It can be traversed with levitation."
	}
	if c.T.IsDiggable() && !c.IsPassable() {
		autodesc += " It is diggable by oric destructive magic."
	}
	if c.IsSwimPassable() && !c.IsPassable() {
		autodesc += " It is possible to traverse by swimming."
	}
	if c.BlocksRange() {
		autodesc += " It blocks ranged attacks from foes."
	}
	if c.Hides() {
		autodesc += " You can hide just behind it."
	}
	if autodesc != "" {
		desc += "\n\n" + strings.TrimSpace(autodesc)
	}
	return desc
}

func (c cell) Style(g *game, pos position) (r rune, fg uicolor) {
	switch c.T {
	case WallCell:
		r, fg = '#', ColorFgLOS
	case GroundCell:
		r, fg = '.', ColorFgLOS
	case DoorCell:
		r, fg = '+', ColorFgPlace
	case FoliageCell:
		r, fg = '"', ColorFgLOS
	case BarrelCell:
		r, fg = '&', ColorFgObject
	case StoneCell:
		r, fg = g.Objects.Stones[pos].Style(g)
	case StairCell:
		st := g.Objects.Stairs[pos]
		r, fg = st.Style(g)
	case MagaraCell:
		r, fg = '/', ColorFgObject
	case BananaCell:
		r, fg = ')', ColorFgObject
	case LightCell:
		r, fg = '☼', ColorFgObject
	case ExtinguishedLightCell:
		r, fg = '○', ColorFgLOS
	case TableCell:
		r, fg = 'π', ColorFgObject
	case TreeCell:
		r, fg = '♣', ColorFgConfusedMonster
	case HoledWallCell:
		r, fg = 'Π', ColorViolet
	case ScrollCell:
		r, fg = g.Objects.Scrolls[pos].Style(g)
	case StoryCell:
		r, fg = g.Objects.Story[pos].Style(g)
	case ItemCell:
		r, fg = g.Objects.Items[pos].Style(g)
	case BarrierCell:
		r, fg = 'Ξ', ColorFgMagicPlace
	case WindowCell:
		r, fg = 'Θ', ColorViolet
	case ChasmCell:
		r, fg = '◊', ColorFgLOS
	case WaterCell:
		r, fg = '≈', ColorFgLOS
	case RubbleCell:
		r, fg = '^', ColorFgLOS
	case CavernCell:
		r, fg = ',', ColorFgLOS
	case FakeStairCell:
		r, fg = '>', ColorFgPlace
	case PotionCell:
		r, fg = g.Objects.Potions[pos].Style(g)
	case QueenRockCell:
		r, fg = '‗', ColorFgLOS
	}
	return r, fg
}
