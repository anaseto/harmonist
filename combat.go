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
	g.md.WoundedAnimation()
	if oldHP > max && g.Player.HP <= max {
		g.StoryPrintf("Critical hit by %s (HP: %d)", m.Kind, g.Player.HP)
		g.md.WoundedAnimation() // twice
		g.md.criticalHPWarning()
	} else if g.Player.HP > 0 {
		g.StoryPrintf("Hit by %s (HP: %d)", m.Kind, g.Player.HP)
	} else {
		g.StoryPrintf("Killed by %s", m.Kind)
	}
	if g.Player.HP > 0 && g.Player.Inventory.Body == CloakConversion && g.Player.MP < g.Player.MPMax() {
		g.Player.MP++
	}
}

func (md *model) criticalHPWarning() {
	md.mode = modeHPCritical
	md.g.PrintStyled("*** CRITICAL HP WARNING ***", logCritic)
	md.g.PrintStyled("[(x) to continue]", logConfirm)
}

func (g *game) MakeMonstersAware() {
	for _, m := range g.Monsters {
		if m.Dead {
			continue
		}
		if g.Player.Sees(m.P) {
			m.MakeAware(g)
			if m.State != Resting {
				m.GatherBand(g)
			}
		}
	}
}

func (g *game) MakeNoise(noise int, at gruid.Point) {
	dij := &noisePath{g: g}
	g.PR.BreadthFirstMap(dij, []gruid.Point{at}, noise)
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
		d := g.PR.BreadthFirstMapAt(m.P)
		if d > noise {
			continue
		}
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
	dij := &monPath{g: g, monster: m}
	nodes := g.PR.DijkstraMap(dij, []gruid.Point{m.P}, 10)
	free := invalidPos
	dist := unreachable
	for _, n := range nodes {
		if !m.CanPass(g, n.P) {
			continue
		}
		if n.P == g.Player.P || n.P == m.P {
			continue
		}
		mons := g.MonsterAt(n.P)
		if mons.Exists() {
			continue
		}
		if distance(n.P, m.P) < dist {
			free = n.P
			dist = distance(n.P, m.P)
		}
	}
	// free should be valid except in really rare cases
	return free
}

func (g *game) FindJumpTarget(m *monster) gruid.Point {
	dij := &jumpPath{g: g}
	nodes := g.PR.DijkstraMap(dij, []gruid.Point{m.P}, 10)
	free := invalidPos
	dist := unreachable
	for _, n := range nodes {
		if !g.PlayerCanPass(n.P) {
			continue
		}
		if n.P == g.Player.P || n.P == m.P {
			continue
		}
		mons := g.MonsterAt(n.P)
		if mons.Exists() {
			continue
		}
		if distance(n.P, m.P) < dist {
			free = n.P
			dist = distance(n.P, m.P)
		}
	}
	// free should be valid except in really rare cases
	return free
}

func (g *game) Jump(mons *monster) (bool, error) {
	if mons.Peaceful(g) && mons.Kind != MonsEarthDragon {
		op := mons.P
		if terrain(g.Dungeon.Cell(op)) == ChasmCell && !g.Player.HasStatus(StatusLevitation) {
			return true, g.AbyssJump()
		}
		if !mons.CanPass(g, g.Player.P) {
			p := mons.LeaveRoomForPlayer(g)
			if p != invalidPos {
				mons.MoveTo(g, p)
				mons.Swapped = true
				g.PlacePlayerAt(op)
				return false, nil
			}
			// otherwise (which should not happen in practice), swap anyways
		}
		g.PlacePlayerAt(mons.P)
		return false, nil
	}
	if g.Player.HasStatus(StatusExhausted) {
		return false, errors.New("You cannot jump while exhausted.")
	}
	dir := dirnorm(g.Player.P, mons.P)
	p := g.Player.P
	for {
		p = p.Add(dir)
		if !g.PlayerCanPass(p) {
			break
		}
		m := g.MonsterAt(p)
		if !m.Exists() {
			break
		}
	}
	if !g.PlayerCanPass(p) {
		p = g.FindJumpTarget(mons)
		if !g.PlayerCanPass(p) {
			// should not happen in practice, but better safe than sorry
			g.Teleportation()
			return false, nil
		}
	}
	if !g.Player.HasStatus(StatusSwift) && g.Player.Inventory.Body != CloakAcrobat {
		g.PutStatus(StatusExhausted, DurationExhaustion)
	}
	if mons.Kind == MonsEarthDragon {
		g.Confusion()
	}
	g.PlacePlayerAt(p)
	g.Stats.Jumps++
	g.Printf("You jump over %s", mons.Kind.Definite(false))
	g.StoryPrintf("Jumped over %s", mons.Kind)
	if g.Stats.Jumps+g.Stats.WallJumps == 15 {
		AchAcrobat.Get(g)
	}
	return false, nil
}

func (g *game) WallJump(p gruid.Point) error {
	c := g.Dungeon.Cell(g.Player.P)
	if c.IsEnclosing() {
		return fmt.Errorf("You cannot jump from %s.", c.ShortDesc(g, g.Player.P))
	}
	if g.Player.HasStatus(StatusExhausted) {
		return errors.New("You cannot jump while exhausted.")
	}
	dir := dirnorm(p, g.Player.P)
	p = g.Player.P
	q := p
	count := 0
	path := []gruid.Point{q}
	for count < 4 {
		p = p.Add(dir)
		if !g.PlayerCanJumpPass(p) {
			break
		}
		count++
		q = p
		path = append(path, q)
		m := g.MonsterAt(p)
		if !m.Exists() && count == 3 {
			break
		}
	}
	m := g.MonsterAt(q)
	if m.Exists() {
		return errors.New("There's not enough room to jump.")
	}
	if count == 3 && !g.PlayerCanPass(q) {
		q = path[len(path)-2]
		m := g.MonsterAt(q)
		if m.Exists() {
			return errors.New("There's not enough room to jump.")
		}
		path = path[:len(path)-1]
	}
	if !g.PlayerCanPass(q) || (count != 3 && count != 2) {
		return errors.New("There's not enough room to jump.")
	}
	if !g.Player.HasStatus(StatusSwift) && g.Player.Inventory.Body != CloakAcrobat {
		g.PutStatus(StatusExhausted, DurationExhaustion)
	}
	g.md.PushAnimation(path)
	g.PlacePlayerAt(q)
	g.Stats.WallJumps++
	g.Print("You jump by propulsing yourself against the wall.")
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
	if g.Player.Sees(mons.P) {
		AchAssassin.Get(g)
	}
	if terrain(g.Dungeon.Cell(mons.P)) == DoorCell {
		g.ComputeLOS()
	}
	g.StoryPrintf("Death of %s", mons.Kind.Indefinite(false))
}

const (
	WallNoise              = 10
	ExplosionNoise         = 14
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
