package main

import (
	"errors"
)

func (g *game) EvokeBlink() error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	g.Blink()
	return nil
}

func (g *game) Blink() bool {
	if g.Player.HasStatus(StatusLignification) {
		return false
	}
	npos := g.BlinkPos(false)
	if !npos.valid() {
		// should not happen
		g.Print("You could not blink.")
		return false
	}
	opos := g.Player.Pos
	if npos == opos {
		g.Print("You blink in-place.")
	} else {
		g.Print("You blink away.")
	}
	g.ui.TeleportAnimation(opos, npos, true)
	g.PlacePlayerAt(npos)
	return true
}

func (g *game) BlinkPos(mpassable bool) position {
	losPos := []position{}
	for pos, b := range g.Player.LOS {
		// XXX: skip if not seen by monster?
		if !b {
			continue
		}
		c := g.Dungeon.Cell(pos)
		if !c.T.IsPlayerPassable() || mpassable && !c.IsPassable() || c.T == StoryCell {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		losPos = append(losPos, pos)
	}
	if len(losPos) == 0 {
		return InvalidPos
	}
	npos := losPos[RandInt(len(losPos))]
	for i := 0; i < 4; i++ {
		pos := losPos[RandInt(len(losPos))]
		if npos.Distance(g.Player.Pos) < pos.Distance(g.Player.Pos) {
			npos = pos
		}
	}
	return npos
}

func (g *game) EvokeTeleport() error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot teleport while lignified.")
	}
	g.Teleportation()
	return nil
}

func (g *game) EvokeDig() error {
	if !g.PutStatus(StatusDig, DurationDigging) {
		return errors.New("You are already digging.")
	}
	g.Print("You feel like an earth dragon.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) MonstersInLOS() []*monster {
	ms := []*monster{}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.Pos) {
			ms = append(ms, mons)
		}
	}
	// shuffle before, because the order could be unnaturally predicted
	for i := 0; i < len(ms); i++ {
		j := i + RandInt(len(ms)-i)
		ms[i], ms[j] = ms[j], ms[i]
	}
	return ms
}

func (g *game) MonstersInCardinalLOS() []*monster {
	ms := []*monster{}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.Pos) && (mons.Pos.X == g.Player.Pos.X || mons.Pos.Y == g.Player.Pos.Y) {
			ms = append(ms, mons)
		}
	}
	// shuffle before, because the order could be unnaturally predicted
	for i := 0; i < len(ms); i++ {
		j := i + RandInt(len(ms)-i)
		ms[i], ms[j] = ms[j], ms[i]
	}
	return ms
}

func (g *game) EvokeTeleportOther() error {
	ms := g.MonstersInCardinalLOS()
	if len(ms) == 0 {
		return errors.New("There are no targetable monsters.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	for i := 0; i < max; i++ {
		if ms[i].Search == InvalidPos {
			ms[i].Search = ms[i].Pos
		}
		ms[i].TeleportAway(g)
		if ms[i].Kind.ReflectsTeleport() {
			g.Printf("The %s reflected back some energies.", ms[i].Kind)
			g.Teleportation()
			break
		}
	}

	return nil
}

func (g *game) EvokeHealWounds() error {
	g.Player.HP = g.Player.HPMax()
	g.Print("Your feel healthy again.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeRefillMagic() error {
	g.Player.MP = g.Player.MPMax()
	g.Print("Your magic forces return.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeSwiftness() error {
	if g.Player.HasStatus(StatusSwift) {
		return errors.New("You are already swift.")
	}
	g.Player.Statuses[StatusSwift] += DurationSwiftness
	g.Printf("You feel swift.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeLevitation() error {
	if !g.PutStatus(StatusLevitation, DurationLevitation) {
		return errors.New("You are already levitating.")
	}
	g.Printf("You feel light.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeSwapping() error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot use this magara while lignified.")
	}
	ms := g.MonstersInLOS()
	var mons *monster
	best := 0
	for _, m := range ms {
		if m.Status(MonsLignified) {
			continue
		}
		if m.Pos.Distance(g.Player.Pos) > best {
			best = m.Pos.Distance(g.Player.Pos)
			mons = m
		}
	}
	if !mons.Exists() {
		return errors.New("No monsters suitable for swapping in view.")
	}
	if mons.Kind.CanOpenDoors() {
		// only intelligent monsters understand swapping
		mons.Search = mons.Pos
	}
	g.SwapWithMonster(mons)
	return nil
}

func (g *game) SwapWithMonster(mons *monster) {
	ompos := mons.Pos
	g.Printf("You swap positions with the %s.", mons.Kind)
	g.ui.SwappingAnimation(mons.Pos, g.Player.Pos)
	mons.MoveTo(g, g.Player.Pos)
	g.PlacePlayerAt(ompos)
	mons.MakeAware(g)
	if g.Dungeon.Cell(g.Player.Pos).T == ChasmCell {
		g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: AbyssFall})
	}
}

type cloud int

const (
	CloudFog cloud = iota
	CloudFire
	CloudNight
)

func (g *game) EvokeFog() error {
	g.Fog(g.Player.Pos, 3)
	g.Print("You are surrounded by a dense fog.")
	return nil
}

func (g *game) Fog(at position, radius int) {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{at}, radius)
	nm.iter(at, func(n *node) {
		pos := n.Pos
		_, ok := g.Clouds[pos]
		if !ok && g.Dungeon.Cell(pos).AllowsFog() {
			g.Clouds[pos] = CloudFog
			g.PushEvent(&posEvent{ERank: g.Ev.Rank() + DurationFog + RandInt(DurationFog/2), EAction: CloudEnd, Pos: pos})
		}
	})
	g.ComputeLOS()
}

func (g *game) EvokeShadows() error {
	if g.Player.HasStatus(StatusIlluminated) {
		return errors.New("You cannot surround yourself by shadows while illuminated.")
	}
	if !g.PutStatus(StatusShadows, DurationShadows) {
		return errors.New("You are already surrounded by shadows.")
	}
	g.Print("You are surrounded by shadows.")
	return nil
}

func (g *game) EvokeParalysis() error {
	ms := g.MonstersInLOS()
	count := 0
	for _, mons := range ms {
		if mons.PutStatus(g, MonsParalysed, DurationParalysisMonster) {
			count++
			if mons.Search == InvalidPos {
				mons.Search = mons.Pos
			}
		}
	}
	if count == 0 {
		return errors.New("No suitable targets in view.")
	}
	g.Print("Whoosh! A slowing luminous wave emerges.")
	g.ui.LOSWavesAnimation(DefaultLOSRange, WaveSlowing, g.Player.Pos)
	return nil
}

func (g *game) EvokeSleeping() error {
	ms := g.MonstersInCardinalLOS()
	if len(ms) == 0 {
		return errors.New("There are no targetable monsters.")
	}
	max := 3
	if max > len(ms) {
		max = len(ms)
	}
	targets := []position{}
	// XXX: maybe use noise distance instead of LOS?
	for i := 0; i < max; i++ {
		mons := ms[i]
		if mons.State != Resting {
			g.Printf("%s falls asleep.", mons.Kind.Definite(true))
		} else {
			continue
		}
		mons.State = Resting
		mons.Dir = NoDir
		mons.ExhaustTime(g, 4+RandInt(2))
		targets = append(targets, g.Ray(mons.Pos)...)
	}
	if len(targets) == 0 {
		return errors.New("There are no suitable targets.")
	}
	if max == 1 {
		g.Print("A beam of sleeping emerges.")
	} else {
		g.Print("Two beams of sleeping emerge.")
	}
	g.ui.BeamsAnimation(targets, BeamSleeping)

	return nil
}

func (g *game) EvokeLignification() error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	targets := []position{}
	for i := 0; i < max; i++ {
		mons := ms[i]
		if mons.Status(MonsLignified) || mons.Kind.ResistsLignification() {
			continue
		}
		mons.EnterLignification(g)
		if mons.Search == InvalidPos {
			mons.Search = mons.Pos
		}
		targets = append(targets, g.Ray(mons.Pos)...)
	}
	if len(targets) == 0 {
		return errors.New("There are no suitable targets.")
	}
	if max == 1 {
		g.Print("A beam of lignification emerges.")
	} else {
		g.Print("Two beams of lignification emerge.")
	}
	g.ui.BeamsAnimation(targets, BeamLignification)
	return nil
}

func (g *game) EvokeNoise() error {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{g.Player.Pos}, 23)
	noises := []position{}
	g.NoiseIllusion = map[position]bool{}
	for _, mons := range g.Monsters {
		if !mons.Exists() {
			continue
		}
		n, ok := nm.at(mons.Pos)
		if !ok || n.Cost > DefaultLOSRange {
			continue
		}
		if mons.SeesPlayer(g) {
			continue
		}
		mp := &monPath{game: g, monster: mons}
		target := mons.Pos
		best := n.Cost
		for {
			ncost := best
			for _, pos := range mp.Neighbors(target) {
				node, ok := nm.at(pos)
				if !ok {
					continue
				}
				ncost := node.Cost
				if ncost > best {
					target = pos
					best = ncost
				}
			}
			if ncost == best {
				break
			}
		}
		if mons.Kind == MonsSatowalgaPlant {
			mons.State = Hunting
		} else if mons.State != Hunting {
			mons.State = Wandering
		}
		mons.Target = target
		noises = append(noises, target)
		g.NoiseIllusion[target] = true
	}
	g.ui.NoiseAnimation(noises)
	g.Print("Monsters are tricked by magical sounds.")
	return nil
}

func (g *game) EvokeConfusion() error {
	ms := g.MonstersInLOS()
	count := 0
	for _, mons := range ms {
		if mons.EnterConfusion(g) {
			count++
			if mons.Search == InvalidPos {
				mons.Search = mons.Pos
			}
		}
	}
	if count == 0 {
		return errors.New("No suitable targets in view.")
	}
	g.Print("Whoosh! A confusing luminous wave emerges.")
	g.ui.LOSWavesAnimation(DefaultLOSRange, WaveConfusion, g.Player.Pos)
	return nil
}

func (g *game) EvokeFire() error {
	burnpos := g.Dungeon.CardinalFlammableNeighbors(g.Player.Pos)
	if len(burnpos) == 0 {
		return errors.New("You are not surrounded by any flammable terrain.")
	}
	g.Print("Sparks emanate from the magara.")
	for _, pos := range burnpos {
		g.Burn(pos)
	}
	return nil
}

func (g *game) EvokeObstruction() error {
	targets := []position{}
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		ray := g.Ray(mons.Pos)
		for i, pos := range ray[1:] {
			if pos == g.Player.Pos {
				break
			}
			mons := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			g.MagicalBarrierAt(pos)
			if len(ray) == 0 {
				break
			}
			ray = ray[i+1:]
			targets = append(targets, ray...)
			break
		}
	}
	if len(targets) == 0 {
		return errors.New("No targetable monsters in view.")
	}
	g.Print("Magical barriers emerged.")
	g.ui.BeamsAnimation(targets, BeamObstruction)
	return nil
}

func (g *game) MagicalBarrierAt(pos position) {
	if g.Dungeon.Cell(pos).T == WallCell || g.Dungeon.Cell(pos).T == BarrierCell {
		return
	}
	g.UpdateKnowledge(pos, g.Dungeon.Cell(pos).T)
	g.CreateMagicalBarrierAt(pos)
	g.ComputeLOS()
}

func (g *game) CreateMagicalBarrierAt(pos position) {
	t := g.Dungeon.Cell(pos).T
	g.Dungeon.SetCell(pos, BarrierCell)
	delete(g.Clouds, pos)
	g.MagicalBarriers[pos] = t
	g.PushEvent(&posEvent{ERank: g.Ev.Rank() + DurationMagicalBarrier + RandInt(DurationMagicalBarrier/2), Pos: pos, EAction: ObstructionEnd})
}

func (g *game) EvokeEnergyMagara() error {
	if g.Player.MP == g.Player.MPMax() && g.Player.HP == g.Player.HPMax() {
		return errors.New("You are already full of energy.")
	}
	g.Print("The magara glows.")
	g.ui.PlayerGoodEffectAnimation()
	g.Player.MP = g.Player.MPMax()
	g.Player.HP = g.Player.HPMax()
	return nil
}

func (g *game) EvokeTransparencyMagara() error {
	if !g.PutStatus(StatusTransparent, DurationTransparency) {
		return errors.New("You are already transparent.")
	}
	g.Print("Light makes you diaphanous.")
	return nil
}

func (g *game) EvokeDisguiseMagara() error {
	if !g.PutStatus(StatusDisguised, DurationDisguise) {
		return errors.New("You are already disguised.")
	}
	g.Print("You look now like a normal guard.")
	return nil
}

func (g *game) EvokeDelayedNoiseMagara() error {
	if !g.PutFakeStatus(StatusDelay, DurationHarmonicNoiseDelay) {
		return errors.New("You are already using delayed magic.")
	}
	g.PushEvent(&posEvent{ERank: g.Ev.Rank() + DurationTurn,
		Pos: g.Player.Pos, EAction: DelayedHarmonicNoiseEvent,
		Timer: DurationHarmonicNoiseDelay})
	g.Print("Timer activated.")
	return nil
}

func (g *game) EvokeDispersalMagara() error {
	if !g.PutStatus(StatusDispersal, DurationDispersal) {
		return errors.New("You are already dispersing.")
	}
	g.Print("You feel unstable.")
	return nil
}

func (g *game) EvokeOricExplosionMagara() error {
	if !g.PutFakeStatus(StatusDelay, DurationHarmonicNoiseDelay) {
		return errors.New("You are already using delayed magic.")
	}
	g.PushEvent(&posEvent{ERank: g.Ev.Rank() + DurationTurn,
		Pos: g.Player.Pos, EAction: DelayedOricExplosionEvent,
		Timer: DurationOricExplosionDelay})
	g.Print("Timer activated.")
	return nil
}
