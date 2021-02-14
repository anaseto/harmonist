package main

import (
	"strings"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

type cell rl.Cell

const (
	WallCell cell = iota
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
	Explored = 0b10000000
)

func terrain(c cell) cell {
	return c &^ Explored
}

func explored(c cell) bool {
	return c&Explored != 0
}

func (c cell) IsPassable() bool {
	switch terrain(c) {
	case WallCell, DoorCell, BarrelCell, TableCell, TreeCell, HoledWallCell, BarrierCell, WindowCell, StoryCell, ChasmCell, WaterCell:
		return false
	default:
		return true
	}
}

func (c cell) IsJumpPassable() bool {
	switch terrain(c) {
	case TableCell, ChasmCell, WaterCell, StoryCell:
		return true
	default:
		return c.IsPassable()
	}
}

func (c cell) IsNormalPatrolWay() bool {
	switch terrain(c) {
	case GroundCell, ScrollCell, DoorCell, StairCell, LightCell, ItemCell, ExtinguishedLightCell, StoneCell, MagaraCell, FakeStairCell:
		return true
	default:
		return false
	}
}

func (c cell) IsLevitatePassable() bool {
	switch terrain(c) {
	case ChasmCell:
		return true
	default:
		return c.IsPassable()
	}
}

func (c cell) IsDoorPassable() bool {
	switch terrain(c) {
	case DoorCell:
		return true
	default:
		return c.IsPassable()
	}
}

func (c cell) IsSwimPassable() bool {
	switch terrain(c) {
	case WaterCell:
		return true
	default:
		return c.IsPassable()
	}
}

func (c cell) IsJumpPropulsion() bool {
	switch terrain(c) {
	case WallCell, WindowCell:
		return true
	default:
		return false
	}
}

func (c cell) IsEnclosing() bool {
	switch terrain(c) {
	case BarrelCell, TableCell, HoledWallCell:
		return true
	default:
		return false
	}
}

func (c cell) AllowsFog() bool {
	switch terrain(c) {
	case WallCell, HoledWallCell, WindowCell, StoryCell:
		return false
	default:
		return true
	}
}

func (c cell) CoversPlayer() bool {
	switch terrain(c) {
	case WallCell, BarrelCell, TableCell, TreeCell, HoledWallCell, BarrierCell, WindowCell:
		return true
	default:
		return false
	}
}

func (c cell) IsPlayerPassable() bool {
	switch terrain(c) {
	case WallCell, BarrierCell, WindowCell, ChasmCell:
		return false
	default:
		return true
	}
}

func (c cell) IsDiggable() bool {
	switch terrain(c) {
	case WallCell, WindowCell, HoledWallCell:
		return true
	default:
		return false
	}
}

func (c cell) BlocksRange() bool {
	switch terrain(c) {
	case WallCell, TreeCell, BarrierCell, WindowCell, StoryCell:
		return true
	default:
		return false
	}
}

func (c cell) Hides() bool {
	switch terrain(c) {
	case WallCell, BarrelCell, TableCell, TreeCell, WindowCell, StoryCell:
		return true
	default:
		return false
	}
}

func (c cell) IsIlluminable() bool {
	switch terrain(c) {
	case WallCell, BarrelCell, TableCell, TreeCell, HoledWallCell, BarrierCell, WindowCell, ChasmCell, RubbleCell:
		return false
	}
	return true
}

func (c cell) IsDestructible() bool {
	switch terrain(c) {
	case WallCell, BarrelCell, DoorCell, TableCell, TreeCell, HoledWallCell, WindowCell:
		return true
	default:
		return false
	}
}

func (c cell) IsWall() bool {
	switch terrain(c) {
	case WallCell:
		return true
	default:
		return false
	}
}

func (c cell) Flammable() bool {
	switch terrain(c) {
	case FoliageCell, DoorCell, BarrelCell, TableCell, TreeCell, WindowCell:
		return true
	default:
		return false
	}
}

func (c cell) IsGround() bool {
	switch terrain(c) {
	case GroundCell, CavernCell, BananaCell, PotionCell, QueenRockCell:
		return true
	default:
		return false
	}
}

func (c cell) IsNotable() bool {
	switch terrain(c) {
	case StairCell, StoneCell, BarrelCell, MagaraCell, BananaCell,
		ScrollCell, ItemCell, FakeStairCell, PotionCell:
		return true
	default:
		return false
	}
}

func (c cell) ShortString(g *game, p gruid.Point) (desc string) {
	switch terrain(c) {
	case WallCell:
		desc = "wall"
	case GroundCell:
		desc = "paved ground"
	case DoorCell:
		desc = "door"
	case FoliageCell:
		desc = "foliage"
	case BarrelCell:
		desc = "barrel"
	case StoneCell:
		desc = g.Objects.Stones[p].String()
	case StairCell:
		desc = g.Objects.Stairs[p].ShortString(g)
	case MagaraCell:
		desc = g.Objects.Magaras[p].String()
	case BananaCell:
		desc = "banana"
	case LightCell:
		desc = "campfire"
	case ExtinguishedLightCell:
		desc = "extinguished campfire"
	case TableCell:
		desc = "table"
	case TreeCell:
		desc = "banana tree"
	case HoledWallCell:
		desc = "holed wall"
	case ScrollCell:
		desc = g.Objects.Scrolls[p].String()
	case StoryCell:
		desc = g.Objects.Story[p].String()
	case ItemCell:
		desc = g.Objects.Items[p].String()
	case BarrierCell:
		desc = "magical barrier"
	case WindowCell:
		desc = "closed window"
	case ChasmCell:
		if g.Depth == MaxDepth || g.Depth == WinDepth {
			desc = "deep chasm"
		} else {
			desc = "chasm"
		}
	case WaterCell:
		desc = "shallow water"
	case RubbleCell:
		desc = "rubblestone"
	case CavernCell:
		desc = "cave ground"
	case FakeStairCell:
		if g.Depth == WinDepth {
			desc = DeepStairShortDesc
		} else {
			desc = NormalStairShortDesc
		}
	case PotionCell:
		desc = g.Objects.Potions[p].String()
	case QueenRockCell:
		desc = "queen rock"
	}
	return desc
}

func (c cell) ShortDesc(g *game, p gruid.Point) (desc string) {
	switch terrain(c) {
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
		desc = g.Objects.Stones[p].ShortDesc(g)
	case StairCell:
		desc = g.Objects.Stairs[p].ShortDesc(g)
	case MagaraCell:
		desc = g.Objects.Magaras[p].ShortDesc()
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
		desc = g.Objects.Scrolls[p].ShortDesc()
	case StoryCell:
		desc = g.Objects.Story[p].String()
	case ItemCell:
		desc = g.Objects.Items[p].ShortDesc(g)
	case BarrierCell:
		desc = "a magical barrier"
	case WindowCell:
		desc = "a closed window"
	case ChasmCell:
		if g.Depth == MaxDepth || g.Depth == WinDepth {
			desc = "a deep chasm"
		} else {
			desc = "a chasm"
		}
	case WaterCell:
		desc = "shallow water"
	case RubbleCell:
		desc = "rubblestone"
	case CavernCell:
		desc = "cave ground"
	case FakeStairCell:
		if g.Depth == WinDepth {
			desc = DeepStairShortDesc
		} else {
			desc = NormalStairShortDesc
		}
	case PotionCell:
		desc = g.Objects.Potions[p].ShortDesc(g)
	case QueenRockCell:
		desc = "queen rock"
	}
	return desc
}

func (c cell) Desc(g *game, p gruid.Point) (desc string) {
	switch terrain(c) {
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
		desc = g.Objects.Stones[p].Desc(g)
	case StairCell:
		desc = g.Objects.Stairs[p].Desc(g)
	case MagaraCell:
		desc = g.Objects.Magaras[p].Desc(g)
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
		desc = g.Objects.Scrolls[p].Desc(g)
	case StoryCell:
		desc = g.Objects.Story[p].Desc(g, p)
	case ItemCell:
		desc = g.Objects.Items[p].Desc(g)
	case BarrierCell:
		desc = "A temporal magical barrier created by oric energies. It may have been created by an oric magara or an oric celmist. Sometimes, natural oric energies may produce such barriers too in energetically unstable Underground areas."
	case WindowCell:
		desc = "A transparent window in the wall."
	case ChasmCell:
		if g.Depth == MaxDepth || g.Depth == WinDepth {
			desc = "A deep chasm. If you jump into it, you'll be dead."
		} else {
			desc = "A chasm. If you jump into it, you'll be seriously injured."
		}
	case WaterCell:
		desc = "Shallow water."
	case RubbleCell:
		desc = "Rubblestone is a collection of rocks broken into smaller stones."
	case CavernCell:
		desc = "This is natural cave ground."
	case FakeStairCell:
		if g.Depth == WinDepth {
			desc = DeepStairDesc
		} else {
			desc = NormalStairDesc
		}
	case PotionCell:
		desc = g.Objects.Potions[p].Desc(g)
	case QueenRockCell:
		desc = "Queen rock amplifies sounds. Even though you are usually very silent, monsters may hear your footsteps when walking on those rocks."
	}
	var autodesc string
	if !c.IsPlayerPassable() {
		autodesc += " It is impassable."
	}
	if c.Flammable() {
		autodesc += " It is flammable."
	}
	if c.IsLevitatePassable() && !c.IsPassable() {
		autodesc += " It can be traversed with levitation."
	}
	if c.IsDiggable() && !c.IsPassable() {
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
	if c.IsJumpPropulsion() {
		autodesc += " You can use it to propulse yourself with a jump."
	}
	if autodesc != "" {
		desc += "\n\n" + strings.TrimSpace(autodesc)
	}
	return desc
}

func (c cell) Style(g *game, p gruid.Point) (r rune, fg gruid.Color) {
	switch terrain(c) {
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
		r, fg = g.Objects.Stones[p].Style(g)
	case StairCell:
		st := g.Objects.Stairs[p]
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
		r, fg = g.Objects.Scrolls[p].Style(g)
	case StoryCell:
		r, fg = g.Objects.Story[p].Style(g)
	case ItemCell:
		r, fg = g.Objects.Items[p].Style(g)
	case BarrierCell:
		r, fg = 'Ξ', ColorFgMagicPlace
	case WindowCell:
		r, fg = 'Θ', ColorViolet
	case ChasmCell:
		r, fg = '◊', ColorFgLOS
		if g.Depth == MaxDepth || g.Depth == WinDepth {
			fg = ColorViolet
		}
	case WaterCell:
		r, fg = '≈', ColorFgLOS
	case RubbleCell:
		r, fg = '^', ColorFgLOS
	case CavernCell:
		r, fg = ',', ColorFgLOS
	case FakeStairCell:
		r, fg = '>', ColorFgPlace
		if g.Depth == WinDepth {
			fg = ColorViolet
		}
	case PotionCell:
		r, fg = g.Objects.Potions[p].Style(g)
	case QueenRockCell:
		r, fg = '‗', ColorFgLOS
	}
	return r, fg
}
