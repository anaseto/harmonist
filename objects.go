package main

import (
	"errors"
	"sort"

	"github.com/anaseto/gruid"
)

type objects struct {
	Stairs     map[gruid.Point]stair
	Stones     map[gruid.Point]stone
	Magaras    map[gruid.Point]magara // TODO simplify? (there's never more than one)
	Barrels    map[gruid.Point]bool
	Bananas    map[gruid.Point]bool
	Lights     map[gruid.Point]bool // true: on, false: off
	Scrolls    map[gruid.Point]scroll
	Story      map[gruid.Point]story
	Lore       map[gruid.Point]int // TODO simplify? (there's never more than one)
	Items      map[gruid.Point]item
	FakeStairs map[gruid.Point]bool
	Potions    map[gruid.Point]potion
}

type stair int

const (
	NormalStair stair = iota
	WinStair
	BlockedStair
)

func (st stair) String() (desc string) {
	switch st {
	case NormalStair:
		desc = "stairs"
	case WinStair:
		desc = "monolith portal"
	case BlockedStair:
		desc = "sealed stairs"
	}
	return desc
}

const normalStairShortDesc = "stairs downwards"
const deepStairShortDesc = "deep stairs downwards"

func (st stair) ShortString(g *game) (desc string) {
	switch st {
	case NormalStair:
		if g.Depth == WinDepth {
			desc = deepStairShortDesc
		} else {
			desc = normalStairShortDesc
		}
	case WinStair:
		desc = "monolith portal"
	case BlockedStair:
		desc = "blocked " + NormalStair.ShortDesc(g)
	}
	return desc
}

func (st stair) ShortDesc(g *game) (desc string) {
	switch st {
	case NormalStair:
		if g.Depth == WinDepth {
			desc = deepStairShortDesc
		} else {
			desc = normalStairShortDesc
		}
	case WinStair:
		desc = "a monolith portal"
	case BlockedStair:
		desc = "blocked " + NormalStair.ShortDesc(g)
	}
	return desc
}

const normalStairDesc = "Stairs lead to the next level of Dayoriah Clan's domain in Hareka's Underground. You will not be able to come back, because an oric barrier seals the stairs when they are traversed by intruders. The upside of this is that ennemies cannot follow you either."
const deepStairDesc = "Those very deep stairs lead to the next level of Dayoriah Clan's domain in Hareka's Underground. You will not be able to come back, because an oric barrier seals the stairs when they are traversed by intruders. The upside of this is that ennemies cannot follow you either."

func (st stair) Desc(g *game) (desc string) {
	switch st {
	case WinStair:
		desc = "Going through this portal will make you escape from this place, going back to the Surface."
		if g.Depth < MaxDepth {
			desc += " If you're courageous enough, you may skip this portal and continue going deeper in the dungeon, to find Marevor's magara, finishing Shaedra's failed mission."
		}
	case NormalStair:
		desc = normalStairDesc
		if g.Depth == WinDepth {
			desc = deepStairDesc
			desc += " You may want to take those after freeing Shaedra from her cell."
		}
	case BlockedStair:
		desc = "Stairs lead to the next level of the Dayoriah Clan's domain in Hareka's Underground. These are sealed by an oric magical barrier that you have to disable by activating a corresponding seal stone. You will not be able to come back, because an oric barrier seals the stairs again when they are traversed by intruders. The upside of this is that ennemies cannot follow you either."
	}
	return desc
}

func (st stair) Style(g *game) (r rune, fg gruid.Color) {
	r = '>'
	switch st {
	case WinStair:
		fg = ColorFgMagicPlace
		r = 'Δ'
	case NormalStair:
		fg = ColorFgPlace
		if g.Depth == WinDepth {
			fg = ColorViolet
		}
	case BlockedStair:
		fg = ColorFgMagicPlace
	}
	return r, fg
}

type stone int

const (
	InertStone stone = iota
	BarrelStone
	FogStone
	QueenStone
	NightStone
	TreeStone
	TeleportStone
	MappingStone
	SensingStone
	// special
	SealStone
)

func (stn stone) String() (text string) {
	switch stn {
	case InertStone:
		text = "inert stone"
	case BarrelStone:
		text = "barrel stone"
	case FogStone:
		text = "fog stone"
	case QueenStone:
		text = "queenstone"
	case NightStone:
		text = "night stone"
	case TreeStone:
		text = "tree stone"
	case TeleportStone:
		text = "teleport stone"
	case MappingStone:
		text = "mapping stone"
	case SensingStone:
		text = "sensing stone"
	case SealStone:
		text = "seal stone"
	}
	return text
}

func (stn stone) Desc(g *game) (text string) {
	switch stn {
	case InertStone:
		text = "This magical stone has been depleted of magical energies."
	case BarrelStone:
		text = "Activating this magical stone will teleport you away to a barrel in the same level."
	case FogStone:
		text = "Activating this magical stone will produce fog in a 4-radius area using harmonic energies."
	case QueenStone:
		text = "This magical stone is made from queen rock. Activating it will produce an amplified harmonic sound confusing enemies in a quite large area. This can also attract monsters."
	case NightStone:
		text = "Activating this magical stone will produce hypnotic harmonic sounds and illusions inducing sleep in all the monsters in sight."
	case TreeStone:
		text = "Activating this magical stone will lignify monsters in sight."
	case TeleportStone:
		text = "Activating this magical stone will teleport away monsters in sight."
	case MappingStone:
		text = "Activating this magical stone shows you the map layout and item locations in a wide area."
	case SensingStone:
		text = "Activating this magical stone shows you the current position of monsters in a wide area."
	case SealStone:
		text = "Activating this magical stone will disable a magical barrier somewhere in the same level, usually one blocking stairs."
	}
	return text
}

func (stn stone) ShortDesc(g *game) string {
	return Indefinite(stn.String(), false)
}

func (stn stone) Style(g *game) (r rune, fg gruid.Color) {
	r = '∩'
	switch stn {
	case InertStone:
		fg = ColorFgPlace
	case SealStone:
		fg = ColorFgPlayer
	case MappingStone, SensingStone:
		fg = ColorViolet
	case BarrelStone:
		fg = ColorFgObject
	default:
		fg = ColorFgMagicPlace
	}
	return r, fg
}

func (g *game) UseStone(p gruid.Point) {
	g.StoryPrintf("Activated %s", g.Objects.Stones[p])
	g.Objects.Stones[p] = InertStone
	g.Stats.UsedStones++
	g.Print("The stone becomes inert.")
}

func (g *game) ActivateStone() (err error) {
	stn, ok := g.Objects.Stones[g.Player.P]
	if !ok {
		return errors.New("No stone to activate here.")
	}
	op := g.Player.P
	switch stn {
	case InertStone:
		err = errors.New("The stone is inert.")
	case BarrelStone:
		g.Print("The stone teleports you away.")
		g.TeleportToBarrel()
	case FogStone:
		g.Print("The stone releases fog.")
		g.Fog(g.Player.P, FogStoneDistance)
	case QueenStone:
		g.ActivateQueenStone()
	case NightStone:
		err = g.ActivateNightStone()
	case TreeStone:
		err = g.ActivateTreeStone()
	case TeleportStone:
		err = g.ActivateTeleportStone()
	case MappingStone:
		err = g.MagicMapping(MappingDistance)
	case SensingStone:
		err = g.Sensing()
	case SealStone:
		err = g.BarrierStone()
	}
	if err != nil {
		return err
	}
	g.UseStone(op)
	return nil
}

const (
	FogStoneDistance   = 4
	QueenStoneDistance = 15
	MappingDistance    = 32
)

func (g *game) ActivateQueenStone() {
	g.MakeNoise(QueenStoneNoise, g.Player.P)
	dij := &noisePath{g: g}
	g.PR.BreadthFirstMap(dij, []gruid.Point{g.Player.P}, QueenStoneDistance)
	targets := []*monster{}
	for _, m := range g.Monsters {
		if !m.Exists() {
			continue
		}
		if m.State == Resting {
			continue
		}
		c := g.PR.BreadthFirstMapAt(m.P)
		if c > QueenStoneDistance {
			continue
		}
		targets = append(targets, m)
	}
	g.Print("The stone releases confusing sounds.")
	g.md.LOSWavesAnimation(DefaultLOSRange, WaveMagicNoise, g.Player.P)
	for _, m := range targets {
		m.EnterConfusion(g)
		if m.Search == invalidPos {
			m.Search = m.P
		}
	}
}

func (g *game) ActivateNightStone() error {
	targets := []*monster{}
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.P) {
			continue
		}
		if mons.State != Resting {
			targets = append(targets, mons)
		}
	}
	if len(targets) == 0 {
		return errors.New("There are no suitable monsters in sight.")
	}
	g.Print("The stone releases hypnotic harmonies.")
	g.md.LOSWavesAnimation(DefaultLOSRange, WaveSleeping, g.Player.P)
	for _, mons := range targets {
		g.Printf("%s falls asleep.", mons.Kind.Definite(true))
		mons.State = Resting
		mons.Dir = ZP
		mons.ExhaustTime(g, 4+RandInt(2))
	}
	return nil
}

func (g *game) ActivateTreeStone() error {
	targets := []*monster{}
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.P) || mons.Kind.ResistsLignification() {
			continue
		}
		targets = append(targets, mons)
	}
	if len(targets) == 0 {
		return errors.New("There are no suitable monsters in sight.")
	}
	g.Print("The stone releases magical spores.")
	g.md.LOSWavesAnimation(DefaultLOSRange, WaveTree, g.Player.P)
	for _, mons := range targets {
		mons.EnterLignification(g)
		if mons.Search == invalidPos {
			mons.Search = mons.P
		}
	}
	return nil
}

func (g *game) ActivateTeleportStone() error {
	targets := g.MonstersInLOS()
	if len(targets) == 0 {
		return errors.New("There are no suitable monsters in sight.")
	}
	g.Print("The stone releases oric teleport energies.")
	for _, mons := range targets {
		if mons.Search == invalidPos && mons.Kind.CanOpenDoors() {
			mons.Search = mons.P
		}
		mons.TeleportAway(g)
		if mons.Kind.ReflectsTeleport() {
			g.Printf("The %s reflected back some energies.", mons.Kind)
			g.Teleportation()
			break
		}
	}
	return nil
}

func (g *game) TeleportToBarrel() {
	barrels := []gruid.Point{}
	for p := range g.Objects.Barrels {
		barrels = append(barrels, p)
	}
	p := barrels[RandInt(len(barrels))]
	op := g.Player.P
	g.Print("You teleport away.")
	g.md.TeleportAnimation(op, p, true)
	g.PlacePlayerAt(p)
}

func (g *game) MagicMapping(maxdist int) error {
	dp := &mappingPath{g: g}
	nodes := g.PR.BreadthFirstMap(dp, []gruid.Point{g.Player.P}, maxdist)
	cdists := make(map[int][]int)
	for _, n := range nodes {
		cdists[n.Cost] = append(cdists[n.Cost], idx(n.P))
	}
	var dists []int
	for dist := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	for _, d := range dists {
		if maxdist > 0 && d > maxdist {
			continue
		}
		draw := false
		for _, i := range cdists[d] {
			p := idxtopos(i)
			c := g.Dungeon.Cell(p)
			if !explored(c) {
				g.Dungeon.SetExplored(p)
				g.SeeNotable(c, p)
				draw = true
			}
		}
		if draw {
			g.md.MagicMappingAnimation()
		}
	}
	g.Printf("You feel aware of your surroundings..")
	return nil
}

func (g *game) Sensing() error {
	for _, mons := range g.Monsters {
		if mons.Exists() && !g.Player.Sees(mons.P) && distance(mons.P, g.Player.P) <= MappingDistance {
			mons.UpdateKnowledge(g, mons.P)
		}
	}
	g.Printf("You briefly sense monsters around.")
	return nil
}

func (g *game) BarrierStone() error {
	if g.Depth == MaxDepth {
		g.Objects.Story[g.Places.Artifact] = StoryArtifact
		g.Print("You feel oric energies dissipating.")
		return nil
	}
	for p, st := range g.Objects.Stairs {
		// actually there is at most only such stair
		if st == BlockedStair {
			g.Objects.Stairs[p] = NormalStair
		}
	}
	g.Print("You feel oric energies dissipating.")
	return nil
}

type scroll int

const (
	ScrollStory scroll = iota
	ScrollExtended
	ScrollDayoriahMessage
	ScrollLore
)

func (sc scroll) String() (desc string) {
	switch sc {
	case ScrollLore:
		desc = "message"
	default:
		desc = "story message"
	}
	return desc
}

func (sc scroll) ShortDesc() (desc string) {
	switch sc {
	case ScrollLore:
		desc = "a message"
	default:
		desc = "a story message"
	}
	return desc
}

func (sc scroll) Text(g *game) (desc string) {
	switch sc {
	case ScrollStory:
		desc = "Your friend Shaedra got captured by nasty people from the Dayoriah Clan while she was trying to retrieve a powerful magara artifact that was stolen from the great magara-specialist Marevor Helith.\n\nAs a gawalt monkey, you don't understand much why people complicate so much their lives caring about artifacts and the like, but one thing is clear: you have to rescue your friend, somewhere to be found in this Underground area controlled by the Dayoriah Clan. If what you heard the guards say is true, Shaedra's imprisoned on the eighth floor.\n\nYou are small and have good night vision, so you hope the infiltration will go smoothly..."
	case ScrollExtended:
		desc = "Now that Shaedra's back to safety, you can either follow her advice, and get away from here too using the monolith portal, or you can finish the original mission: going deeper to find Marevor's powerful magara, before the Dayoriah Clan does bad experiments with it. You honestly didn't understand why it was dangerous, but Shaedra and Marevor had seemed truly concerned.\n\nMarevor said that he'll be able to create a new portal for you when you activate the artifact upon finding it."
	case ScrollDayoriahMessage:
		desc = `“The thief that infiltrated our turf and tried to retrieve our new acquisition has been captured. However, it is possible that she has an accomplice. Please be careful and stop every suspect.”

[Order of a Dayoriah Clan's foreman, the paper is rather recent and you can't help but think the thief is actually Shaedra. So they really did capture her! You have to hurry and save her.]`
	case ScrollLore:
		i, ok := g.Objects.Lore[g.Player.P]
		if !ok {
			// should not happen
			desc = "Some unintelligible notes."
			break
		}
		if i < len(LoreMessages) {
			desc = LoreMessages[i]
		}
	default:
		desc = "a message"
	}
	return desc
}

func (sc scroll) Desc(g *game) (desc string) {
	desc = "A message. It can be read."
	return desc
}

func (sc scroll) Style(g *game) (r rune, fg gruid.Color) {
	r = '?'
	fg = ColorFgMagicPlace
	if sc == ScrollLore {
		fg = ColorViolet
	}
	return r, fg
}

type story int

const (
	NoStory story = iota // just a normal ground cell but free
	StoryShaedra
	StoryMarevor
	StoryArtifact
	StoryArtifactSealed
)

func (st story) Desc(g *game, p gruid.Point) (desc string) {
	switch st {
	case NoStory:
		desc = GroundCell.Desc(g, p)
	case StoryShaedra:
		desc = "Shaedra is the friend you came here to rescue, a human-like creature with claws, a ternian. Many other human-like creatures consider them as savages."
	case StoryMarevor:
		desc = "Marevor Helith is an ancient undead nakrus very fond of teleporting people away. He is a well-known expert in the field of magaras - items that many people simply call magical objects. His current research focus is monolith creation. Marevor, a repentant necromancer, is now searching for his old disciple Jaixel in the Underground to help him overcome the past."
	case StoryArtifact:
		desc = "This is the magara that you have to retrieve: the Gem Portal Artifact that was stolen to Marevor Helith."
	case StoryArtifactSealed:
		desc = "This is the magara that you have to retrieve: the Gem Portal Artifact that was stolen to Marevor Helith. Before taking it, you have to release the magical barrier that protects it activating the corresponding protective barrier magical stone."
	}
	return desc
}

func (st story) String() (desc string) {
	switch st {
	case NoStory:
		desc = "paved ground"
	case StoryShaedra:
		desc = "Shaedra"
	case StoryMarevor:
		desc = "Marevor"
	case StoryArtifact:
		desc = "Gem Portal Artifact"
	case StoryArtifactSealed:
		desc = "Gem Portal Artifact (sealed)"
	}
	return desc
}

func (st story) Style(g *game) (r rune, fg gruid.Color) {
	fg = ColorFgPlayer
	switch st {
	case NoStory:
		fg = ColorFgLOS
		r = '.'
	case StoryShaedra:
		r = 'S'
	case StoryMarevor:
		r = 'M'
	case StoryArtifact:
		r = '='
	case StoryArtifactSealed:
		r = '='
		fg = ColorFgMagicPlace
	}
	return r, fg
}

type item int

const (
	NoItem item = iota
	CloakMagic
	CloakHear
	CloakVitality
	CloakAcrobat // no exhaustion between jumps?
	CloakShadows // reduce monster los?
	CloakSmoke
	CloakConversion
	AmuletTeleport
	AmuletConfusion
	AmuletFog
	AmuletLignification
	AmuletObstruction
	MarevorMagara
)

func (it item) IsCloak() bool {
	switch it {
	case CloakMagic,
		CloakHear,
		CloakVitality,
		CloakAcrobat,
		CloakSmoke,
		CloakShadows,
		CloakConversion:
		return true
	}
	return false
}

func (it item) IsAmulet() bool {
	switch it {
	case AmuletTeleport,
		AmuletConfusion,
		AmuletFog,
		AmuletLignification,
		AmuletObstruction:
		return true
	}
	return false
}

func (it item) String() (desc string) {
	switch it {
	case NoItem:
		desc = "empty slot"
	case CloakMagic:
		desc = "cloak of magic"
	case CloakHear:
		desc = "cloak of hearing"
	case CloakVitality:
		desc = "cloak of vitality"
	case CloakAcrobat:
		desc = "cloak of acrobatics"
	case CloakShadows:
		desc = "cloak of shadows"
	case CloakSmoke:
		desc = "cloak of smoking"
	case CloakConversion:
		desc = "cloak of conversion"
	case AmuletTeleport:
		desc = "amulet of teleport"
	case AmuletConfusion:
		desc = "amulet of confusion"
	case AmuletFog:
		desc = "amulet of fog"
	case AmuletLignification:
		desc = "amulet of lignification"
	case AmuletObstruction:
		desc = "amulet of obstruction"
	case MarevorMagara:
		desc = "Moon Portal Artifact"
	}
	return desc
}

func (it item) ShortDesc(g *game) string {
	return it.String()
}

func (it item) Desc(g *game) (desc string) {
	switch it {
	case NoItem:
		return "You do not have an item equipped on this slot."
	case CloakMagic:
		desc = "increases your magical reserves."
	case CloakHear:
		desc = "improves your hearing skills."
	case CloakVitality:
		desc = "improves your health."
	case CloakAcrobat:
		desc = "removes exhaustion from jumps."
	case CloakShadows:
		desc = "reduces the range at which foes see you in the dark."
	case CloakSmoke:
		desc = "leaves smoke behind as you move, making you difficult to spot."
	case CloakConversion:
		desc = "converts lost health from wounds into magical energy."
	case AmuletTeleport:
		desc = "teleports away foes that critically hit you."
	case AmuletConfusion:
		desc = "confuses foes that critically hit you."
	case AmuletFog:
		desc = "releases fog and makes you swift when critically hurt."
	case AmuletLignification:
		desc = "lignifies foes that critically hit you."
	case AmuletObstruction:
		desc = "uses a magical barrier to blow away monsters that critically hit you."
	case MarevorMagara:
		desc = "magara was given to you by Marevor Helith so that he can create an escape portal when you reach Shaedra. Its sister magara, the Gem Portal Artifact, also crafted by Marevor, is the artifact that was stolen and that Shaedra was trying to retrieve before being captured.\n\nThis magara needs a lot of time to recharge, so you'll only be able to use it once."
	}
	return "The " + it.ShortDesc(g) + " " + desc
}

func (it item) Style(g *game) (r rune, fg gruid.Color) {
	fg = ColorFgObject
	if it.IsAmulet() {
		r = '='
	} else if it.IsCloak() {
		r = '['
	}
	return r, fg
}

func (g *game) EquipItem() error {
	it, ok := g.Objects.Items[g.Player.P]
	if !ok {
		return errors.New("Nothing to equip here.")
	}
	var oitem item
	switch {
	case it.IsCloak():
		oitem = g.Player.Inventory.Body
		g.Player.Inventory.Body = it
	case it.IsAmulet():
		oitem = g.Player.Inventory.Neck
		g.Player.Inventory.Neck = it
	}
	if oitem != NoItem {
		g.Objects.Items[g.Player.P] = oitem
		g.Printf("You equip the %s.", it.ShortDesc(g))
		g.Printf("You leave the %s.", oitem.ShortDesc(g))
		g.StoryPrintf("Equipped %s", it.ShortDesc(g))
	} else {
		delete(g.Objects.Items, g.Player.P)
		g.Dungeon.SetCell(g.Player.P, GroundCell)
		g.Printf("You equip the %s.", it.ShortDesc(g))
		g.StoryPrintf("Equipped %s", it.ShortDesc(g))
	}
	switch {
	case it.IsCloak():
		AchCloak.Get(g)
	case it.IsAmulet():
		AchAmulet.Get(g)
	}
	return nil
}

func (g *game) RandomCloak() (it item) {
	cloaks := []item{CloakMagic,
		CloakHear,
		CloakVitality,
		CloakAcrobat,
		CloakSmoke,
		CloakShadows,
		CloakConversion}
loop:
	for {
		it = cloaks[RandInt(len(cloaks))]
		for _, cl := range g.GeneratedCloaks {
			if cl == it {
				continue loop
			}
		}
		break
	}
	return it
}

func (g *game) RandomAmulet() (it item) {
	amulets := []item{AmuletTeleport,
		AmuletConfusion,
		AmuletFog,
		AmuletLignification,
		AmuletObstruction}
loop:
	for {
		it = amulets[RandInt(len(amulets))]
		for _, cl := range g.GeneratedAmulets {
			if cl == it {
				continue loop
			}
		}
		break
	}
	return it
}

type potion int

const (
	HealthPotion potion = iota
	MagicPotion
)

func (ptn potion) String() (desc string) {
	switch ptn {
	case HealthPotion:
		desc = "health potion"
	case MagicPotion:
		desc = "magic potion"
	}
	return desc
}

func (ptn potion) ShortDesc(g *game) (desc string) {
	return Indefinite(ptn.String(), false)
}

func (ptn potion) Desc(g *game) (desc string) {
	switch ptn {
	case HealthPotion:
		desc = "Drinking a health potion will cure 1 HP."
	case MagicPotion:
		desc = "Drinking a magic potion will replenish 1 MP."
	}
	desc += "\n\nYou cannot carry around potions, you drink them when you move onto them."
	return desc
}

func (ptn potion) Style(g *game) (r rune, fg gruid.Color) {
	r = '!'
	switch ptn {
	case HealthPotion:
		fg = ColorFgHPok
	case MagicPotion:
		fg = ColorFgMPok
	}
	return r, fg
}

func (g *game) DrinkPotion(p gruid.Point) {
	ptn, ok := g.Objects.Potions[p]
	if !ok {
		// should not happen
		g.Dungeon.SetCell(p, GroundCell)
		g.PrintStyled("Unexpected potion.", logError)
		return
	}
	switch ptn {
	case HealthPotion:
		if g.Player.HP >= g.Player.HPMax() {
			return
		}
		g.Player.HP++
		g.StoryPrintf("Drank %s (HP: %d).", ptn, g.Player.HP)
	case MagicPotion:
		if g.Player.MP >= g.Player.MPMax() {
			return
		}
		g.Player.MP++
		g.StoryPrintf("Drank %s (MP: %d).", ptn, g.Player.MP)
	}
	g.Printf("You drink %s.", ptn.ShortDesc(g))
	g.Dungeon.SetCell(p, GroundCell)
	delete(g.Objects.Potions, p)
}
