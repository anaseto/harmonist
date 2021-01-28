// combat utility functions

package main

import (
	"errors"
	"fmt"

	"github.com/anaseto/gruid"
)

func (g *game) DamagePlayer(damage int) {
	g.Stats.Damage += damage
	g.Stats.DDamage[g.Depth] += damage
	g.Player.HPbonus -= damage
	if g.Player.HPbonus < 0 {
		g.Player.HP += g.Player.HPbonus
		g.Player.HPbonus = 0
	}
}

func (m *monster) InflictDamage(g *game, damage, max int) {
	g.Stats.ReceivedHits++
	oldHP := g.Player.HP
	g.DamagePlayer(damage)
	// TODO: animation
	//g.ui.WoundedAnimation()
	if oldHP > max && g.Player.HP <= max {
		g.StoryPrintf("Critical hit by %s (HP: %d)", m.Kind, g.Player.HP)
		// TODO: animation
		g.ui.criticalHPWarning()
	} else if g.Player.HP > 0 {
		g.StoryPrintf("Hit by %s (HP: %d)", m.Kind, g.Player.HP)
	} else {
		g.StoryPrintf("Killed by %s", m.Kind)
	}
	if g.Player.HP > 0 && g.Player.Inventory.Body == CloakConversion && g.Player.MP < g.Player.MPMax() {
		g.Player.MP++
	}
}

func (g *game) MakeMonstersAware() {
	for _, m := range g.Monsters {
		if m.Dead {
			continue
		}
		if g.Player.LOS[m.Pos] {
			m.MakeAware(g)
			if m.State != Resting {
				m.GatherBand(g)
			}
		}
	}
}

func (g *game) MakeNoise(noise int, at gruid.Point) {
	dij := &noisePath{state: g}
	nm := Dijkstra(dij, []gruid.Point{at}, noise)
	//if at.Distance(g.Player.Pos)-noise < DefaultLOSRange && noise > 4 {
	//g.ui.LOSWavesAnimation(noise, WaveNoise, at)
	//}
	for _, m := range g.Monsters {
		if !m.Exists() {
			continue
		}
		if m.State == Hunting {
			continue
		}
		n, ok := nm.at(m.Pos)
		if !ok {
			continue
		}
		d := n.Cost
		if m.State == Resting && 3*d > 2*noise || m.Status(MonsExhausted) && m.State == Resting && 3*d > noise {
			continue
		}
		if m.SeesPlayer(g) {
			m.MakeAware(g)
		} else {
			m.MakeWanderAt(at)
		}
		m.GatherBand(g)
	}
}

func (m *monster) LeaveRoomForPlayer(g *game) gruid.Point {
	dij := &monPath{state: g, monster: m}
	nm := Dijkstra(dij, []gruid.Point{m.Pos}, 10)
	free := InvalidPos
	dist := unreachable
	nm.iter(m.Pos, func(n *node) {
		if !m.CanPass(g, n.Pos) {
			return
		}
		if n.Pos == g.Player.Pos || n.Pos == m.Pos {
			return
		}
		mons := g.MonsterAt(n.Pos)
		if mons.Exists() {
			return
		}
		if Distance(n.Pos, m.Pos) < dist {
			free = n.Pos
			dist = Distance(n.Pos, m.Pos)
		}
	})
	// free should be valid except in really rare cases
	return free
}

func (g *game) FindJumpTarget(m *monster) gruid.Point {
	dij := &jumpPath{state: g}
	nm := Dijkstra(dij, []gruid.Point{m.Pos}, 10)
	free := InvalidPos
	dist := unreachable
	nm.iter(m.Pos, func(n *node) {
		if !g.PlayerCanPass(n.Pos) {
			return
		}
		if n.Pos == g.Player.Pos || n.Pos == m.Pos {
			return
		}
		mons := g.MonsterAt(n.Pos)
		if mons.Exists() {
			return
		}
		if Distance(n.Pos, m.Pos) < dist {
			free = n.Pos
			dist = Distance(n.Pos, m.Pos)
		}
	})
	// free should be valid except in really rare cases
	return free
}

func (g *game) Jump(mons *monster) error {
	if mons.Peaceful(g) && mons.Kind != MonsEarthDragon {
		ompos := mons.Pos
		if g.Dungeon.Cell(ompos).T == ChasmCell && !g.Player.HasStatus(StatusLevitation) {
			if g.DeepChasmDepth() {
				return errors.New("You cannot jump into deep chasm.")
			}
			if !g.AbyssJumpConfirmation() {
				return errors.New(doNothing)
			}
			g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: AbyssFall})
		}
		if !mons.CanPass(g, g.Player.Pos) {
			pos := mons.LeaveRoomForPlayer(g)
			if pos != InvalidPos {
				mons.MoveTo(g, pos)
				mons.Swapped = true
				g.PlacePlayerAt(ompos)
				return nil
			}
			// otherwise (which should not happen in practice), swap anyways
		}
		mons.MoveTo(g, g.Player.Pos)
		mons.Swapped = true
		g.PlacePlayerAt(ompos)
		return nil
	}
	if g.Player.HasStatus(StatusExhausted) {
		return errors.New("You cannot jump while exhausted.")
	}
	dir := Dir(g.Player.Pos, mons.Pos)
	pos := g.Player.Pos
	for {
		pos = To(dir, pos)
		if !g.PlayerCanPass(pos) {
			break
		}
		m := g.MonsterAt(pos)
		if !m.Exists() {
			break
		}
	}
	if !g.PlayerCanPass(pos) {
		pos = g.FindJumpTarget(mons)
		if !g.PlayerCanPass(pos) {
			// should not happen in practice, but better safe than sorry
			g.Teleportation()
			return nil
		}
	}
	if !g.Player.HasStatus(StatusSwift) && g.Player.Inventory.Body != CloakAcrobat {
		g.PutStatus(StatusExhausted, 5)
	}
	if mons.Kind == MonsEarthDragon {
		g.Confusion()
	}
	g.PlacePlayerAt(pos)
	g.Stats.Jumps++
	g.Printf("You jump over %s", mons.Kind.Definite(false))
	g.StoryPrintf("Jumped over %s", mons.Kind)
	if g.Stats.Jumps+g.Stats.WallJumps == 15 {
		AchAcrobat.Get(g)
	}
	return nil
}

func (g *game) WallJump(pos gruid.Point) error {
	c := g.Dungeon.Cell(g.Player.Pos)
	if c.IsEnclosing() {
		return fmt.Errorf("You cannot jump from %s.", c.ShortDesc(g, g.Player.Pos))
	}
	if g.Player.HasStatus(StatusExhausted) {
		return errors.New("You cannot jump while exhausted.")
	}
	dir := Dir(pos, g.Player.Pos)
	pos = g.Player.Pos
	tpos := pos
	count := 0
	path := []gruid.Point{tpos}
	for count < 4 {
		pos = To(dir, pos)
		if !g.PlayerCanJumpPass(pos) {
			break
		}
		count++
		tpos = pos
		path = append(path, tpos)
		m := g.MonsterAt(pos)
		if !m.Exists() && count == 3 {
			break
		}
	}
	m := g.MonsterAt(tpos)
	if m.Exists() {
		return errors.New("There's not enough room to jump.")
	}
	if count == 3 && !g.PlayerCanPass(tpos) {
		tpos = path[len(path)-2]
		m := g.MonsterAt(tpos)
		if m.Exists() {
			return errors.New("There's not enough room to jump.")
		}
		path = path[:len(path)-1]
	}
	if !g.PlayerCanPass(tpos) || (count != 3 && count != 2) {
		return errors.New("There's not enough room to jump.")
	}
	if !g.Player.HasStatus(StatusSwift) && g.Player.Inventory.Body != CloakAcrobat {
		g.PutStatus(StatusExhausted, 5)
	}
	g.PlacePlayerAt(tpos)
	g.Stats.WallJumps++
	g.Print("You jump by propulsing yourself against the wall.")
	// TODO: animation
	//g.ui.PushAnimation(path)
	if g.Stats.Jumps+g.Stats.WallJumps == 15 {
		AchAcrobat.Get(g)
	}
	return nil
}

func (g *game) HitNoise(clang bool) int {
	noise := BaseHitNoise
	if clang {
		noise += 5
	}
	return noise
}

const (
	DmgNormal = 1
)

func (g *game) HandleKill(mons *monster) {
	g.Stats.Killed++
	g.Stats.KilledMons[mons.Kind]++
	if g.Player.Sees(mons.Pos) {
		AchAssassin.Get(g)
	}
	if g.Dungeon.Cell(mons.Pos).T == DoorCell {
		g.ComputeLOS()
	}
	g.StoryPrintf("Death of %s", mons.Kind.Indefinite(false))
}

const (
	WallNoise              = 10
	TemporalWallNoise      = 5
	ExplosionNoise         = 14
	MagicHitNoise          = 5
	BarkNoise              = 10
	MagicCastNoise         = 5
	HarmonicNoise          = 7
	BaseHitNoise           = 2
	QueenStoneNoise        = 10
	SingingNoise           = 12
	EarthquakeNoise        = 35
	QueenRockFootstepNoise = 7
	DelayedHarmonicNoise   = 25
	OricExplosionNoise     = 20
)

func (g *game) ClangMsg() (sclang string) {
	return " Smash!"
}
