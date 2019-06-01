// combat utility functions

package main

import "errors"

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
	g.ui.WoundedAnimation()
	if oldHP > max && g.Player.HP <= max {
		g.StoryPrintf("Critical: hit by %s (HP: %d)", m.Kind.Indefinite(false), g.Player.HP)
		g.ui.CriticalHPWarning()
	} else {
		g.StoryPrintf("Hit by %s (HP: %d)", m.Kind.Indefinite(false), g.Player.HP)
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

func (g *game) MakeNoise(noise int, at position) {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{at}, noise)
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
		if m.State == Resting && d > noise/2 || m.Status(MonsExhausted) && m.State == Resting && d > noise/3 {
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

func (g *game) Jump(mons *monster, ev event) error {
	if mons.Peaceful(g) && mons.Kind != MonsEarthDragon {
		ompos := mons.Pos
		mons.MoveTo(g, g.Player.Pos)
		mons.Swapped = true
		g.PlacePlayerAt(ompos)
		return nil
	}
	dir := mons.Pos.Dir(g.Player.Pos)
	pos := g.Player.Pos
	for {
		pos = pos.To(dir)
		if !pos.valid() || !g.Dungeon.Cell(pos).T.IsPlayerPassable() {
			break
		}
		m := g.MonsterAt(pos)
		if !m.Exists() {
			break
		}
	}
	if !pos.valid() || !g.Dungeon.Cell(pos).T.IsPlayerPassable() {
		return errors.New("You cannot jump in that direction.")
	}
	if g.Player.HasStatus(StatusSlow) {
		return errors.New("You cannot jump while slowed.")
	}
	if g.Player.HasStatus(StatusExhausted) {
		return errors.New("You cannot jump while exhausted.")
	}
	if !g.Player.HasStatus(StatusSwift) && g.Player.Inventory.Body != CloakAcrobat {
		g.PutStatus(StatusExhausted, 50)
	}
	if mons.Kind == MonsEarthDragon {
		g.Confusion(ev)
	}
	g.PlacePlayerAt(pos)
	g.Stats.Jumps++
	if g.Stats.Jumps == 10 {
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
	DmgExtra  = 2
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
	if mons.Kind.Dangerousness() > 10 {
		g.StoryPrintf("%s died.", mons.Kind.Indefinite(true))
	}
}

const (
	WallNoise           = 9
	TemporalWallNoise   = 5
	ExplosionNoise      = 12
	MagicHitNoise       = 5
	BarkNoise           = 9
	MagicExplosionNoise = 12
	MagicCastNoise      = 5
	BaseHitNoise        = 2
	QueenStoneNoise     = 9
)

func (g *game) ClangMsg() (sclang string) {
	if RandInt(2) == 0 {
		sclang = " Clang!"
	} else {
		sclang = " Smash!"
	}
	return sclang
}
