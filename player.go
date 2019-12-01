package main

import (
	"errors"
	"fmt"
)

type player struct {
	HP      int
	HPbonus int
	MP      int
	Bananas int
	Magaras []magara
	Dir     direction
	//Aptitudes map[aptitude]bool
	Statuses  map[status]int
	Expire    map[status]int
	Pos       position
	Target    position
	LOS       map[position]bool
	Rays      rayMap
	Inventory inventory
}

type inventory struct {
	Body item
	Neck item
	Misc item
}

const DefaultHealth = 5

func (p *player) HPMax() int {
	hpmax := DefaultHealth
	if p.Inventory.Body == CloakVitality {
		hpmax += 2
	}
	if hpmax < 2 {
		hpmax = 2
	}
	return hpmax
}

const DefaultMPmax = 6

func (p *player) MPMax() int {
	mpmax := DefaultMPmax
	if p.Inventory.Body == CloakMagic {
		mpmax += 2
	}
	return mpmax
}

func (p *player) HasStatus(st status) bool {
	return p.Statuses[st] > 0
}

func (g *game) AutoToDir(ev event) bool {
	if g.MonsterInLOS() == nil {
		err := g.MovePlayer(g.Player.Pos.To(g.AutoDir), ev)
		if err != nil {
			g.Print(err.Error())
			g.AutoDir = NoDir
			return false
		}
		return true
	}
	g.AutoDir = NoDir
	return false
}

func (g *game) GoToDir(dir direction, ev event) error {
	if g.MonsterInLOS() != nil {
		g.AutoDir = NoDir
		return errors.New("You cannot travel while there are monsters in view.")
	}
	err := g.MovePlayer(g.Player.Pos.To(dir), ev)
	if err != nil {
		return err
	}
	g.AutoDir = dir
	return nil
}

func (g *game) MoveToTarget(ev event) bool {
	if !g.AutoTarget.valid() {
		return false
	}
	path := g.PlayerPath(g.Player.Pos, g.AutoTarget)
	if g.MonsterInLOS() != nil {
		g.AutoTarget = InvalidPos
	}
	if len(path) < 1 {
		g.AutoTarget = InvalidPos
		return false
	}
	var err error
	if len(path) > 1 {
		err = g.MovePlayer(path[len(path)-2], ev)
		if g.ExclusionsMap[path[len(path)-2]] {
			g.AutoTarget = InvalidPos
		}
	} else {
		g.WaitTurn(ev)
	}
	if err != nil {
		g.Print(err.Error())
		g.AutoTarget = InvalidPos
		return false
	}
	if g.AutoTarget.valid() && g.Player.Pos == g.AutoTarget {
		g.AutoTarget = InvalidPos
	}
	return true
}

func (g *game) WaitTurn(ev event) {
	ev.Renew(g, DurationTurn)
}

func (g *game) MonsterCount() (count int) {
	for _, mons := range g.Monsters {
		if mons.Exists() {
			count++
		}
	}
	return count
}

func (g *game) Rest(ev event) error {
	if g.Dungeon.Cell(g.Player.Pos).T != BarrelCell {
		return fmt.Errorf("This place is not safe for sleeping.")
	}
	if cld, ok := g.Clouds[g.Player.Pos]; ok && cld == CloudFire {
		return errors.New("You cannot rest on flames.")
	}
	if !g.NeedsRegenRest() && !g.StatusRest() {
		return errors.New("You do not need to rest.")
	}
	if g.Player.Bananas <= 0 {
		return errors.New("You cannot sleep without eating for dinner a banana first.")
	}
	g.ui.DrawMessage("Resting...")
	g.WaitTurn(ev)
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

func (g *game) Teleportation(ev event) {
	// XXX ev is not used
	var pos position
	i := 0
	count := 0
	for {
		count++
		if count > 1000 {
			panic("Teleportation")
		}
		pos = g.FreePassableCell()
		if pos.Distance(g.Player.Pos) < 15 && i < 1000 {
			i++
			continue
		}
		break

	}
	if pos.valid() {
		// should always happen
		opos := g.Player.Pos
		g.Print("You teleport away.")
		g.ui.TeleportAnimation(opos, pos, true)
		g.PlacePlayerAt(pos)
	} else {
		// should not happen
		g.Print("Something went wrong with the teleportation.")
	}
}

const MaxBananas = 4

func (g *game) CollectGround() {
	pos := g.Player.Pos
	c := g.Dungeon.Cell(pos)
	if c.IsNotable() {
		g.DijkstraMapRebuild = true
	switchcell:
		switch c.T {
		case BarrelCell:
			// TODO: move here message
		case BananaCell:
			if g.Player.Bananas >= MaxBananas {
				g.Print("There is a banana, but your pack is already full.")
			} else {
				g.Print("You take a banana.")
				g.Player.Bananas++
				g.StoryPrintf("Found banana (bananas: %d)", g.Player.Bananas)
				g.Dungeon.SetCell(pos, GroundCell)
				delete(g.Objects.Bananas, pos)
				if g.Player.Bananas == MaxBananas {
					AchBananaCollector.Get(g)
				}
			}
		case MagaraCell:
			for i, mag := range g.Player.Magaras {
				if mag.Kind != NoMagara {
					continue
				}
				g.Player.Magaras[i] = g.Objects.Magaras[pos]
				delete(g.Objects.Magaras, pos)
				g.Dungeon.SetCell(pos, GroundCell)
				g.Printf("You take the %s.", g.Player.Magaras[i])
				g.StoryPrintf("Took %s", g.Player.Magaras[i])
				break switchcell
			}
			g.Printf("You stand over %s.", Indefinite(g.Objects.Magaras[pos].String(), false))
		case FakeStairCell:
			g.Dungeon.SetCell(pos, GroundCell)
			g.PrintStyled("You stand over fake stairs.", logSpecial)
			g.PrintStyled("Harmonic illusions!", logSpecial)
			g.StoryPrint("Found harmonic fake stairs!")
			g.ui.FoundFakeStairsAnimation()
		case PotionCell:
			g.DrinkPotion(pos)
		default:
			g.Printf("You are standing over %s.", c.ShortDesc(g, pos))
		}
	} else if c.T == DoorCell {
		g.Print("You stand at the door.")
	}
	if c.T.ReachNotable() {
		pos.Reach(g)
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

func (g *game) AbyssJump() error {
	if g.Depth == WinDepth || g.Depth == MaxDepth {
		return errors.New("You cannot jump into deep chasm.")
	}
	g.Print("Do you really want to jump into the abyss? (DANGEROUS) [y/N]")
	g.ui.DrawDungeonView(NoFlushMode)
	g.ui.Flush()
	jump := g.ui.PromptConfirmation()
	if jump {
		g.FallAbyss(DescendJump)
		return nil
	}
	return errors.New(DoNothing)
}
func (g *game) MovePlayer(pos position, ev event) error {
	if !pos.valid() {
		return errors.New("You cannot move there.")
	}
	c := g.Dungeon.Cell(pos)
	if c.T == WallCell && !g.Player.HasStatus(StatusDig) {
		return errors.New("You cannot move into a wall.")
	} else if c.T == BarrierCell && !g.Player.HasStatus(StatusLevitation) {
		return errors.New("You cannot move into a magical barrier.")
	} else if c.T == WindowCell && !g.Player.HasStatus(StatusDig) {
		return errors.New("You cannot pass through the closed window.")
	} else if c.T == BarrelCell && g.MonsterLOS[g.Player.Pos] {
		return errors.New("You cannot enter a barrel while seen.")
	}
	mons := g.MonsterAt(pos)
	if !mons.Exists() {
		if g.Player.HasStatus(StatusLignification) {
			return errors.New("You cannot move while lignified.")
		}
		if c.T == ChasmCell && !g.Player.HasStatus(StatusLevitation) {
			return g.AbyssJump()
		}
		if c.T == BarrelCell {
			g.Print("You hide yourself inside the barrel.")
		} else if c.T == TableCell {
			g.Print("You hide yourself under the table.")
		} else if c.T == TreeCell {
			g.Print("You climb to the top.")
		} else if c.T == HoledWallCell {
			g.Print("You crawl under the wall.")
		}
		if c.T.IsDiggable() && c.T != HoledWallCell {
			g.Dungeon.SetCell(pos, RubbleCell)
			g.MakeNoise(WallNoise, pos)
			g.Print(g.CrackSound())
			g.Fog(pos, 1, ev)
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
			_, ok := g.Clouds[g.Player.Pos]
			if !ok && g.Dungeon.Cell(g.Player.Pos).AllowsFog() {
				g.Clouds[g.Player.Pos] = CloudFog
				g.PushEvent(&posEvent{ERank: ev.Rank() + DurationSmokingCloakFog, EAction: CloudEnd, Pos: g.Player.Pos})
			}
		}
		//}
		g.Stats.Moves++
		g.PlacePlayerAt(pos)
	} else if err := g.Jump(mons, ev); err != nil {
		return err
	}
	if g.Player.HasStatus(StatusSwift) {
		ev.Renew(g, 0)
		g.Player.Statuses[StatusSwift]--
		if !g.Player.HasStatus(StatusSwift) {
			g.Print("You no longer feel swift.")
		}
		return nil
	}
	ev.Renew(g, DurationTurn)
	return nil
}

func (g *game) SwiftFog(ev event) {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{g.Player.Pos}, 2)
	nm.iter(g.Player.Pos, func(n *node) {
		pos := n.Pos
		_, ok := g.Clouds[pos]
		if !ok && g.Dungeon.Cell(pos).AllowsFog() {
			g.Clouds[pos] = CloudFog
			g.PushEvent(&posEvent{ERank: ev.Rank() + DurationFog + RandInt(DurationFog/2), EAction: CloudEnd, Pos: pos})
		}
	})
	g.PutStatus(StatusSwift, DurationShortSwiftness)
	g.ComputeLOS()
	g.Print("You feel an energy burst and smoke comes out from you.")
}

func (g *game) Confusion(ev event) {
	if g.PutStatus(StatusConfusion, DurationConfusionPlayer) {
		g.Print("You feel confused.")
	}
}

func (g *game) PlacePlayerAt(pos position) {
	if pos == g.Player.Pos {
		return
	}
	g.Player.Dir = pos.Dir(g.Player.Pos)
	switch g.Player.Dir {
	case ENE, ESE:
		g.Player.Dir = E
	case NNE, NNW:
		g.Player.Dir = N
	case WNW, WSW:
		g.Player.Dir = W
	case SSW, SSE:
		g.Player.Dir = S
	}
	g.Player.Pos = pos
	if g.Dungeon.Cell(pos).T == QueenRockCell {
		g.MakeNoise(QueenRockFootstepNoise, pos)
		g.Print("Tap-tap.")
	}
	g.CollectGround()
	g.ComputeLOS()
	g.MakeMonstersAware()
}

const LignificationHPbonus = 4

func (g *game) EnterLignification(ev event) {
	if g.PutStatus(StatusLignification, DurationLignificationPlayer) {
		g.Print("You feel rooted to the ground.")
		g.Player.HPbonus += LignificationHPbonus
	}
}

func (g *game) ExtinguishFire() error {
	g.Dungeon.SetCell(g.Player.Pos, ExtinguishedLightCell)
	g.Objects.Lights[g.Player.Pos] = false
	g.Stats.Extinguishments++
	if g.Stats.Extinguishments >= 15 {
		AchExtinguisher.Get(g)
	}
	g.Print("You extinguish the fire.")
	g.Ev.Renew(g, DurationTurn)
	return nil
}

func (g *game) PutStatus(st status, duration int) bool {
	if g.Player.Statuses[st] != 0 {
		return false
	}
	g.Player.Statuses[st] += duration
	g.PushEvent(&simpleEvent{ERank: g.Ev.Rank() + DurationStatusStep, EAction: StatusEndActions[st]})
	g.Stats.Statuses[st]++
	if st.Good() {
		g.Player.Expire[st] = g.Ev.Rank() + duration
	}
	return true
}

func (g *game) UpdateKnowledge(pos position, t terrain) {
	if g.Player.Sees(pos) {
		return
	}
	_, ok := g.TerrainKnowledge[pos]
	if !ok {
		g.TerrainKnowledge[pos] = t
	}
}
