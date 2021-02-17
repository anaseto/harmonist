package main

import (
	"errors"
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

type player struct {
	HP      int
	HPbonus int
	MP      int
	Bananas int
	Magaras []magara
	Dir     gruid.Point
	//Aptitudes map[aptitude]bool
	Statuses  map[status]int
	Expire    map[status]int
	P         gruid.Point
	Target    gruid.Point
	LOS       map[gruid.Point]bool
	FOV       *rl.FOV
	Inventory inventory
}

type inventory struct {
	Body item
	Neck item
	Misc item
}

const DefaultHealth = 5

func (pl *player) HPMax() int {
	hpmax := DefaultHealth
	if pl.Inventory.Body == CloakVitality {
		hpmax += 2
	}
	if hpmax < 2 {
		hpmax = 2
	}
	return hpmax
}

const DefaultMPmax = 6

func (pl *player) MPMax() int {
	mpmax := DefaultMPmax
	if pl.Inventory.Body == CloakMagic {
		mpmax += 2
	}
	return mpmax
}

func (pl *player) HasStatus(st status) bool {
	return pl.Statuses[st] > 0
}

func (g *game) AutoToDir() bool {
	if g.MonsterInLOS() == nil {
		p := g.Player.P.Add(g.AutoDir)
		if g.PlayerCanPass(p) {
			again, err := g.PlayerBump(p)
			if err != nil {
				g.Print(err.Error())
				g.AutoDir = ZP
				return false
			}
			if again {
				return false
			}
		} else {
			g.AutoDir = ZP
			return false
		}
		return true
	}
	g.AutoDir = ZP
	return false
}

func (g *game) GoToDir(dir gruid.Point) (again bool, err error) {
	if g.MonsterInLOS() != nil {
		g.AutoDir = ZP
		return again, errors.New("You cannot travel while there are monsters in view.")
	}
	p := g.Player.P.Add(dir)
	if !g.PlayerCanPass(p) {
		return again, errors.New("You cannot move in that direction.")
	}
	again, err = g.PlayerBump(p)
	if err != nil || again {
		return again, err
	}
	g.AutoDir = dir
	return again, err
}

func (g *game) MoveToTarget() bool {
	if !valid(g.AutoTarget) {
		return false
	}
	path := g.PlayerPath(g.Player.P, g.AutoTarget)
	if g.MonsterInLOS() != nil {
		g.AutoTarget = InvalidPos
	}
	if len(path) < 1 {
		g.AutoTarget = InvalidPos
		return false
	}
	var err error
	var again bool
	if len(path) > 1 {
		again, err = g.PlayerBump(path[1])
		if g.ExclusionsMap[path[1]] {
			g.AutoTarget = InvalidPos
		}
	} else {
		g.WaitTurn()
	}
	if err != nil {
		g.Print(err.Error())
		g.AutoTarget = InvalidPos
		return false
	}
	if valid(g.AutoTarget) && g.Player.P == g.AutoTarget {
		g.AutoTarget = InvalidPos
	}
	return !again
}

func (g *game) WaitTurn() {
	g.Stats.Waits++
}

func (g *game) MonsterCount() (count int) {
	for _, mons := range g.Monsters {
		if mons.Exists() {
			count++
		}
	}
	return count
}

func (g *game) Rest() error {
	if terrain(g.Dungeon.Cell(g.Player.P)) != BarrelCell {
		return fmt.Errorf("This place is not safe for sleeping.")
	}
	if cld, ok := g.Clouds[g.Player.P]; ok && cld == CloudFire {
		return errors.New("You cannot rest on flames.")
	}
	if !g.NeedsRegenRest() && !g.StatusRest() {
		return errors.New("You do not need to rest.")
	}
	if g.Player.Bananas <= 0 {
		return errors.New("You cannot sleep without eating for dinner a banana first.")
	}
	// TODO: animation
	//g.ui.DrawMessage("Resting...")
	g.Resting = true
	g.RestingTurns = RandInt(5) // you do not wake up when you want
	g.Player.Bananas--
	return nil
}

func (g *game) StatusRest() bool {
	for st, q := range g.Player.Statuses {
		if st.Info() {
			continue
		}
		if q > 0 {
			return true
		}
	}
	return false
}

func (g *game) NeedsRegenRest() bool {
	return g.Player.HP < g.Player.HPMax() || g.Player.MP < g.Player.MPMax()
}

func (g *game) Teleportation() {
	var p gruid.Point
	i := 0
	count := 0
	for {
		count++
		if count > 1000 {
			panic("Teleportation")
		}
		p = g.FreePassableCell()
		if distance(p, g.Player.P) < 15 && i < 1000 {
			i++
			continue
		}
		break

	}
	if valid(p) {
		// should always happen
		opos := g.Player.P
		g.Print("You teleport away.")
		g.md.TeleportAnimation(opos, p, true)
		g.PlacePlayerAt(p)
	} else {
		// should not happen
		g.Print("Something went wrong with the teleportation.")
	}
}

const MaxBananas = 4

func (g *game) CollectGround() {
	p := g.Player.P
	c := g.Dungeon.Cell(p)
	if c.IsNotable() {
		g.AutoexploreMapRebuild = true
	switchcell:
		switch terrain(c) {
		case BarrelCell:
			// TODO: move here message
		case BananaCell:
			if g.Player.Bananas >= MaxBananas {
				g.Print("There is a banana, but your pack is already full.")
			} else {
				g.Print("You take a banana.")
				g.Player.Bananas++
				g.StoryPrintf("Found banana (bananas: %d)", g.Player.Bananas)
				g.Dungeon.SetCell(p, GroundCell)
				delete(g.Objects.Bananas, p)
				if g.Player.Bananas == MaxBananas {
					AchBananaCollector.Get(g)
				}
			}
		case MagaraCell:
			for i, mag := range g.Player.Magaras {
				if mag.Kind != NoMagara {
					continue
				}
				g.Player.Magaras[i] = g.Objects.Magaras[p]
				delete(g.Objects.Magaras, p)
				g.Dungeon.SetCell(p, GroundCell)
				g.Printf("You take the %s.", g.Player.Magaras[i])
				g.StoryPrintf("Took %s", g.Player.Magaras[i])
				break switchcell
			}
			g.Printf("You stand over %s.", Indefinite(g.Objects.Magaras[p].String(), false))
		case FakeStairCell:
			g.Dungeon.SetCell(p, GroundCell)
			g.PrintStyled("You stand over fake stairs.", logSpecial)
			g.PrintStyled("Harmonic illusions!", logSpecial)
			g.StoryPrint("Found harmonic fake stairs!")
			g.md.FoundFakeStairsAnimation()
		case PotionCell:
			g.DrinkPotion(p)
		default:
			g.Printf("You are standing over %s.", c.ShortDesc(g, p))
		}
	} else if terrain(c) == DoorCell {
		g.Print("You stand at the door.")
	}
	if terrain(c).ReachNotable() {
		g.Reach(p)
	}
}

func (g *game) FallAbyss(style descendstyle) {
	if g.Player.HasStatus(StatusLevitation) {
		return
	}
	g.Player.HP -= 2
	if g.Player.HP <= 0 {
		g.Player.HP = 1
	}
	g.Player.MP -= 2
	if g.Player.MP < 0 {
		g.Player.MP = 0
	}
	if g.Player.Bananas > 0 {
		g.Player.Bananas--
	}
	if style == DescendFall && g.Depth == MaxDepth || g.Depth == WinDepth {
		g.Player.HP = 0
		return
	}
	g.Descend(style)
}

func (g *game) AbyssJumpConfirmation() {
	g.Print("Do you really want to jump into the abyss? (DANGEROUS) [y/N]")
	g.md.mode = modeJumpConfirmation
	// TODO confirmation prompt abyss
	//return g.ui.PromptConfirmation()
}

func (g *game) DeepChasmDepth() bool {
	return g.Depth == WinDepth || g.Depth == MaxDepth
}

func (g *game) AbyssJump() error {
	if g.DeepChasmDepth() {
		return errors.New("You cannot jump into deep chasm.")
	}
	g.AbyssJumpConfirmation()
	return nil
}

func (g *game) PlayerBump(p gruid.Point) (again bool, err error) {
	if !valid(p) {
		return again, errors.New("You cannot move there.")
	}
	c := g.Dungeon.Cell(p)
	switch {
	case terrain(c) == BarrierCell && !g.Player.HasStatus(StatusLevitation):
		return again, errors.New("You cannot move into a magical barrier.")
	case terrain(c) == WindowCell && !g.Player.HasStatus(StatusDig):
		return again, errors.New("You cannot pass through the closed window.")
	case terrain(c) == BarrelCell && g.MonsterLOS[g.Player.P]:
		return again, errors.New("You cannot enter a barrel while seen.")
	}
	mons := g.MonsterAt(p)
	if c.IsJumpPropulsion() && !g.Player.HasStatus(StatusDig) {
		err := g.WallJump(p)
		if err != nil {
			return again, err
		}
	} else if !mons.Exists() {
		if g.Player.HasStatus(StatusLignification) {
			return again, errors.New("You cannot move while lignified.")
		}
		if terrain(c) == ChasmCell && !g.Player.HasStatus(StatusLevitation) {
			again = true
			return again, g.AbyssJump()
		}
		if terrain(c) == BarrelCell {
			g.Print("You hide yourself inside the barrel.")
		} else if terrain(c) == TableCell {
			g.Print("You hide yourself under the table.")
		} else if terrain(c) == TreeCell {
			g.Print("You climb to the top.")
		} else if terrain(c) == HoledWallCell {
			g.Print("You crawl under the wall.")
		}
		if c.IsDiggable() && terrain(c) != HoledWallCell {
			g.Dungeon.SetCell(p, RubbleCell)
			g.MakeNoise(WallNoise, p)
			g.Print(g.CrackSound())
			g.Fog(p, 1)
			g.Stats.Digs++
			g.Stats.DestructionUse++
			if g.Stats.DestructionUse == 20 {
				AchDestructorNovice.Get(g)
			}
			if g.Stats.DestructionUse == 40 {
				AchDestructorInitiate.Get(g)
			}
			if g.Stats.DestructionUse == 60 {
				AchDestructorMaster.Get(g)
			}
		}
		if g.Player.Inventory.Body == CloakSmoke {
			_, ok := g.Clouds[g.Player.P]
			if !ok && g.Dungeon.Cell(g.Player.P).AllowsFog() {
				g.Clouds[g.Player.P] = CloudFog
				g.PushEventD(&posEvent{P: g.Player.P, Action: CloudEnd}, DurationSmokingCloakFog)
			}
		}
		//}
		g.Stats.Moves++
		g.PlacePlayerAt(p)
	} else if again, err = g.Jump(mons); err != nil {
		return again, err
	}
	if g.Player.HasStatus(StatusSwift) {
		again = true
		g.Player.Statuses[StatusSwift]--
		if !g.Player.HasStatus(StatusSwift) {
			g.Print("You no longer feel swift.")
		}
		g.md.updateMapInfo()
		return again, nil
	}
	return again, nil
}

func (g *game) PushPlayerTurn() {
	g.PushEventD(&playerEvent{Action: PlayerTurn}, DurationTurn)
}

func (g *game) SwiftFog() {
	dij := &noisePath{state: g}
	nodes := g.PR.DijkstraMap(dij, []gruid.Point{g.Player.P}, 2)
	for _, n := range nodes {
		_, ok := g.Clouds[n.P]
		if !ok && g.Dungeon.Cell(n.P).AllowsFog() {
			g.Clouds[n.P] = CloudFog
			g.PushEvent(&posEvent{P: n.P, Action: CloudEnd}, g.Turn+DurationFog+RandInt(DurationFog/2))
		}
	}
	g.PutStatus(StatusSwift, DurationShortSwiftness)
	g.ComputeLOS()
	g.Print("You feel an energy burst and smoke comes out from you.")
}

func (g *game) Confusion() {
	if g.PutStatus(StatusConfusion, DurationConfusionPlayer) {
		g.Print("You feel confused.")
	}
}

// PlacePlayerAt moves the player to a given position, swapping positions with
// a monster if necessary, and handles LOS and monsters awareness update,
// ground collecting, and footsteps noise.
func (g *game) PlacePlayerAt(p gruid.Point) {
	if p == g.Player.P {
		return
	}
	g.Player.Dir = dirnorm(g.Player.P, p)
	m := g.MonsterAt(p)
	ppos := g.Player.P
	g.Player.P = p
	if m.Exists() {
		m.MoveTo(g, ppos)
		m.Swapped = true
	}
	if terrain(g.Dungeon.Cell(g.Player.P)) == QueenRockCell && !g.Player.HasStatus(StatusLevitation) {
		g.MakeNoise(QueenRockFootstepNoise, g.Player.P)
		g.Print("Tap-tap.")
	}
	g.CollectGround()
	g.ComputeLOS()
	g.MakeMonstersAware()
}

const LignificationHPbonus = 4

func (g *game) EnterLignification() {
	if g.PutStatus(StatusLignification, DurationLignificationPlayer) {
		g.Print("You feel rooted to the ground.")
		g.Player.HPbonus += LignificationHPbonus
	}
}

func (g *game) ExtinguishFire() error {
	g.Dungeon.SetCell(g.Player.P, ExtinguishedLightCell)
	g.Objects.Lights[g.Player.P] = false
	g.Stats.Extinguishments++
	if g.Stats.Extinguishments >= 15 {
		AchExtinguisher.Get(g)
	}
	g.Print("You extinguish the fire.")
	return nil
}

func (g *game) PutStatus(st status, duration int) bool {
	if g.Player.Statuses[st] != 0 {
		return false
	}
	g.Player.Statuses[st] += duration
	g.PushEventD(&statusEvent{Status: st}, DurationStatusStep)
	g.Stats.Statuses[st]++
	if st.Good() {
		g.Player.Expire[st] = g.Turn + duration
	}
	return true
}

func (g *game) PutFakeStatus(st status, duration int) bool {
	if g.Player.Statuses[st] != 0 {
		return false
	}
	g.Player.Statuses[st] += duration
	g.Stats.Statuses[st]++
	if st.Good() {
		g.Player.Expire[st] = g.Turn + duration
	}
	return true
}

func (g *game) UpdateKnowledge(p gruid.Point, c cell) {
	if g.Player.Sees(p) {
		return
	}
	_, ok := g.TerrainKnowledge[p]
	if !ok {
		g.TerrainKnowledge[p] = c
	}
}

func (g *game) PlayerCanPass(p gruid.Point) bool {
	if !valid(p) {
		return false
	}
	c := g.Dungeon.Cell(p)
	return c.IsPlayerPassable() ||
		g.Player.HasStatus(StatusLevitation) && (terrain(c) == BarrierCell || c.IsLevitatePassable()) ||
		g.Player.HasStatus(StatusDig) && c.IsDiggable()
}

func (g *game) PlayerCanJumpPass(p gruid.Point) bool {
	if !valid(p) {
		return false
	}
	c := g.Dungeon.Cell(p)
	return c.IsJumpPassable() ||
		g.Player.HasStatus(StatusLevitation) && terrain(c) == BarrierCell
}
