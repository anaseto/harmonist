package main

import (
	"errors"
	"fmt"
)

type magara struct {
	Kind    magaraKind
	Charges int
}

type magaraKind int

const (
	NoMagara magaraKind = iota
	BlinkMagara
	DigMagara
	TeleportMagara
	SwiftnessMagara
	LevitationMagara
	FireMagara
	FogMagara
	ShadowsMagara
	NoiseMagara
	ConfusionMagara
	SleepingMagara
	TeleportOtherMagara
	SwappingMagara
	SlowingMagara
	ObstructionMagara
	LignificationMagara
	//BarrierMagara
)

const NumMagaras = int(LignificationMagara)

func (mag magara) Harmonic() bool {
	switch mag.Kind {
	case FogMagara, ShadowsMagara, NoiseMagara, ConfusionMagara, SleepingMagara, SlowingMagara:
		return true
	default:
		return false
	}
}

func (mag magara) Oric() bool {
	switch mag.Kind {
	case BlinkMagara, DigMagara, TeleportMagara, LevitationMagara, TeleportOtherMagara, SwappingMagara, ObstructionMagara:
		return true
	default:
		return false
	}
}

func (m magaraKind) DefaultCharges() int {
	switch m {
	case FireMagara:
		return 10
	case BlinkMagara, LevitationMagara, FogMagara, NoiseMagara:
		return 9
	case SlowingMagara, ObstructionMagara, SwiftnessMagara:
		return 8
	default:
		return 7
	}
}

func (g *game) RandomMagara() magara {
	var mag magaraKind
loop:
	for {
		mag = magaraKind(1 + RandInt(NumMagaras))
		for _, m := range g.GeneratedMagaras {
			if m == mag {
				continue loop
			}
		}
		break
	}
	return magara{Kind: mag, Charges: mag.DefaultCharges()}
}

func (g *game) EquipMagara(i int, ev event) (err error) {
	omagara := g.Player.Magaras[i]
	g.Player.Magaras[i] = g.Objects.Magaras[g.Player.Pos]
	g.Objects.Magaras[g.Player.Pos] = omagara
	g.Printf("You take the %s.", g.Player.Magaras[i])
	g.Printf("You leave the %s.", omagara)
	g.StoryPrintf("Took %s (%d), left %s (%d)", g.Player.Magaras[i], g.Player.Magaras[i].Charges, omagara, omagara.Charges)
	ev.Renew(g, 5)
	return nil
}

func (g *game) UseMagara(n int, ev event) (err error) {
	if g.Player.HasStatus(StatusNausea) {
		return errors.New("You cannot use magaras while sick.")
	}
	if g.Player.HasStatus(StatusConfusion) {
		return errors.New("You cannot use magaras while confused.")
	}
	mag := g.Player.Magaras[n]
	if mag.Kind == NoMagara {
		return errors.New("You cannot evoke an empty slot!")
	}
	if mag.MPCost(g) > g.Player.MP {
		return errors.New("Not enough magic points for using this magara.")
	}
	if mag.Charges <= 0 {
		return errors.New("Not enough charges for using this magara.")
	}
	switch mag.Kind {
	case BlinkMagara:
		err = g.EvokeBlink(ev)
	case DigMagara:
		err = g.EvokeDig(ev)
	case TeleportMagara:
		err = g.EvokeTeleport(ev)
	case SwiftnessMagara:
		err = g.EvokeSwiftness(ev)
	case LevitationMagara:
		err = g.EvokeLevitation(ev)
	case FireMagara:
		err = g.EvokeFire(ev)
	case FogMagara:
		err = g.EvokeFog(ev)
	case ShadowsMagara:
		err = g.EvokeShadows(ev)
	case NoiseMagara:
		err = g.EvokeNoise(ev)
	case ConfusionMagara:
		err = g.EvokeConfusion(ev)
	case SlowingMagara:
		err = g.EvokeSlowing(ev)
	case SleepingMagara:
		err = g.EvokeSleeping(ev)
	case TeleportOtherMagara:
		err = g.EvokeTeleportOther(ev)
	case SwappingMagara:
		err = g.EvokeSwapping(ev)
	case ObstructionMagara:
		err = g.EvokeObstruction(ev)
	case LignificationMagara:
		err = g.EvokeLignification(ev)
	}
	if err != nil {
		return err
	}
	g.Stats.MagarasUsed++
	g.Stats.UsedMagaras[mag]++
	g.Stats.DMagaraUses[g.Depth]++
	g.Player.MP -= mag.MPCost(g)
	g.Player.Magaras[n].Charges--
	g.StoryPrintf("Evoked %s (MP: %d, Charges: %d)", mag, g.Player.MP, g.Player.Magaras[n].Charges)
	if mag.Harmonic() {
		g.Stats.HarmonicMagUse++
		if g.Stats.HarmonicMagUse == 7 {
			AchHarmonistNovice.Get(g)
		}
		if g.Stats.HarmonicMagUse == 14 {
			AchHarmonistInitiate.Get(g)
		}
		if g.Stats.HarmonicMagUse == 21 {
			AchHarmonistMaster.Get(g)
		}
	} else if mag.Oric() {
		g.Stats.OricMagUse++
		if g.Stats.OricMagUse == 7 {
			AchNoviceOricCelmist.Get(g)
		}
		if g.Stats.OricMagUse == 14 {
			AchInitiateOricCelmist.Get(g)
		}
		if g.Stats.OricMagUse == 21 {
			AchMasterOricCelmist.Get(g)
		}
	} else if mag.Kind == FireMagara {
		g.Stats.FireUse++
		if g.Stats.FireUse == 3 {
			AchPyromancerNovice.Get(g)
		}
		if g.Stats.FireUse == 6 {
			AchPyromancerInitiate.Get(g)
		}
		if g.Stats.FireUse == 10 {
			AchPyromancerMaster.Get(g)
		}
	}
	switch mag.Kind {
	case TeleportMagara, TeleportOtherMagara, BlinkMagara, SwappingMagara:
		g.Stats.OricTelUse++
		if g.Stats.OricTelUse == 15 {
			AchTeleport.Get(g)
		}
	}
	ev.Renew(g, 5)
	return nil
}

func (mag magara) String() (desc string) {
	switch mag.Kind {
	case NoMagara:
		desc = "empty slot"
	case BlinkMagara:
		desc = "magara of blinking"
	case DigMagara:
		desc = "magara of digging"
	case TeleportMagara:
		desc = "magara of teleportation"
	case SwiftnessMagara:
		desc = "magara of swiftness"
	case LevitationMagara:
		desc = "magara of levitation"
	case FireMagara:
		desc = "magara of fire"
	case FogMagara:
		desc = "magara of fog"
	case ShadowsMagara:
		desc = "magara of shadows"
	case NoiseMagara:
		desc = "magara of noise"
	case ConfusionMagara:
		desc = "magara of confusion"
	case SleepingMagara:
		desc = "magara of sleeping"
	case TeleportOtherMagara:
		desc = "magara of teleport other"
	case SwappingMagara:
		desc = "magara of swapping"
	case SlowingMagara:
		desc = "magara of slowing"
	case ObstructionMagara:
		desc = "magara of obstruction"
	case LignificationMagara:
		desc = "magara of lignification"
	}
	return desc
}

func (mag magara) ShortDesc() string {
	return fmt.Sprintf("%s (%d)", mag.String(), mag.Charges)
}

func (mag magara) Desc(g *game) (desc string) {
	switch mag.Kind {
	case NoMagara:
		desc = "can be used for a new magara."
	case BlinkMagara:
		desc = "makes you blink away within your line of sight by using an oric energy disturbance. The magara is more susceptible to send you to the cells that are most far from you."
	case DigMagara:
		desc = "makes you dig walls by walking into them like an earth dragon thanks to destructive oric magic."
	case TeleportMagara:
		desc = "creates an oric energy disturbance, making you teleport far away on the same level."
	case SwiftnessMagara:
		desc = "makes you move faster for a short time by giving a boost to your jaipu energies."
	case LevitationMagara:
		desc = "makes you levitate with oric energies, allowing you to move over chasms, as well as through oric barriers."
	case FireMagara:
		desc = "throws small magical sparks at flammable terrain adjacent to you. Flammable terrain is first consumed by magical flames that are at first harmless to creatures. Then smoke will produce night clouds inducing sleep in monsters. As a gawalt monkey, you resist sleepiness, but you will still feel slowed."
	case FogMagara:
		desc = "creates a dense fog in a 2-range radius using harmonic energies. The fog will dissipate with time."
	case ShadowsMagara:
		desc = "surrounds you by harmonic shadows, making you detectable only by adjacent monsters when you're not in an lighted cell."
	case NoiseMagara:
		desc = "tricks monsters in a 12-range area with harmonic magical sounds, making them go away from you for a few turns. It only works on monsters that are not already seeing you."
	case ConfusionMagara:
		desc = "confuses monsters in sight with harmonic light and sounds, leaving them unable to attack you."
	case SlowingMagara:
		desc = "induces slow movement and attack for monsters in sight by disturbing their senses with sound and light illusions."
	case SleepingMagara:
		desc = "induces deep sleeping and exhaustion for up to three random monsters you see in cardinal directions using hypnotic illusions."
	case TeleportOtherMagara:
		desc = "creates oric energy disturbances, teleporting up to two random monsters you see in cardinal directions."
	case SwappingMagara:
		desc = "makes you swap positions with the farthest monster in sight. If there is more than one at the same distance, it will be chosen randomly."
	case ObstructionMagara:
		desc = "creates temporal barriers with oric energy between you and monsters in sight."
	case LignificationMagara:
		desc = "liberates magical spores that lignify up to 2 monsters in view, so that they cannot move. The monsters can still fight."
	}
	duration := 0
	switch mag.Kind {
	case ConfusionMagara:
		duration = DurationConfusionMonster
	case SlowingMagara:
		duration = DurationSlowMonster
	case ObstructionMagara:
		duration = DurationMagicalBarrier
	case LignificationMagara:
		duration = DurationLignificationMonster
	case ShadowsMagara:
		duration = DurationShadows
	case DigMagara:
		duration = DurationDigging
	case SwiftnessMagara:
		duration = DurationSwiftness
	case LevitationMagara:
		duration = DurationLevitation
	}
	if duration > 0 {
		desc += fmt.Sprintf(" Effect lasts for %d turns.", duration/10)
	}
	desc += fmt.Sprintf("\n\nIt currently has %d charges.", mag.Charges)
	return fmt.Sprintf("The %s %s", mag, desc)
}

func (mag magara) MPCost(g *game) int {
	if mag.Kind == NoMagara {
		return 0
	}
	cost := 1
	return cost
}

func (g *game) EvokeBlink(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	g.Blink(ev)
	return nil
}

func (g *game) Blink(ev event) bool {
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

func (g *game) EvokeTeleport(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot teleport while lignified.")
	}
	g.Teleportation(ev)
	g.Print("You teleported away.")
	return nil
}

func (g *game) EvokeDig(ev event) error {
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

func (g *game) EvokeTeleportOther(ev event) error {
	ms := g.MonstersInCardinalLOS()
	if len(ms) == 0 {
		return errors.New("There are no targetable monsters.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	for i := 0; i < max; i++ {
		ms[i].TeleportAway(g)
	}

	return nil
}

func (g *game) EvokeHealWounds(ev event) error {
	g.Player.HP = g.Player.HPMax()
	g.Print("Your feel healthy again.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeRefillMagic(ev event) error {
	g.Player.MP = g.Player.MPMax()
	g.Print("Your magic forces return.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeSwiftness(ev event) error {
	if !g.PutStatus(StatusSwift, DurationSwiftness) {
		return errors.New("You are already swift.")
	}
	g.Printf("You feel speedy.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeLevitation(ev event) error {
	if !g.PutStatus(StatusLevitation, DurationLevitation) {
		return errors.New("You are already levitating.")
	}
	g.Printf("You feel light.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeSwapping(ev event) error {
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
}

type cloud int

const (
	CloudFog cloud = iota
	CloudFire
	CloudNight
)

func (g *game) EvokeFog(ev event) error {
	g.Fog(g.Player.Pos, 3, ev)
	g.Print("You are surrounded by a dense fog.")
	return nil
}

func (g *game) Fog(at position, radius int, ev event) {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{at}, radius)
	nm.iter(at, func(n *node) {
		pos := n.Pos
		_, ok := g.Clouds[pos]
		if !ok && g.Dungeon.Cell(pos).AllowsFog() {
			g.Clouds[pos] = CloudFog
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationFog + RandInt(DurationFog/2), EAction: CloudEnd, Pos: pos})
		}
	})
	g.ComputeLOS()
}

func (g *game) EvokeShadows(ev event) error {
	if g.Player.HasStatus(StatusIlluminated) {
		return errors.New("You cannot surround yourself by shadows while illuminated.")
	}
	if !g.PutStatus(StatusShadows, DurationShadows) {
		return errors.New("You are already surrounded by shadows.")
	}
	g.Print("You are surrounded by shadows.")
	return nil
}

func (g *game) EvokeSlowing(ev event) error {
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		mons.PutStatus(g, MonsSlow, DurationSlowMonster)
	}
	g.Print("Whoosh! A slowing luminous wave emerges.")
	g.ui.LOSWavesAnimation(DefaultLOSRange, WaveSlowing, g.Player.Pos)

	return nil
}

func (g *game) EvokeSleeping(ev event) error {
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
		}
		mons.State = Resting
		mons.Dir = NoDir
		mons.ExhaustTime(g, 40+RandInt(10))
		targets = append(targets, g.Ray(mons.Pos)...)
	}
	if max == 1 {
		g.Print("A beam of sleeping emerges.")
	} else {
		g.Print("Two beams of sleeping emerge.")
	}
	g.ui.BeamsAnimation(targets, BeamSleeping)

	return nil
}

func (g *game) EvokeLignification(ev event) error {
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
		if mons.Status(MonsLignified) {
			continue
		}
		mons.EnterLignification(g, ev)
		targets = append(targets, g.Ray(mons.Pos)...)
	}
	if max == 1 {
		g.Print("A beam of lignification emerges.")
	} else {
		g.Print("Two beams of lignification emerge.")
	}
	g.ui.BeamsAnimation(targets, BeamLignification)
	return nil
}

func (g *game) EvokeNoise(ev event) error {
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
		if mons.SeesPlayer(g) || mons.Kind == MonsSatowalgaPlant {
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
		if mons.State != Hunting {
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

func (g *game) EvokeConfusion(ev event) error {
	g.Print("Whoosh! A confusing luminous wave emerges.")
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		mons.EnterConfusion(g, ev)
	}
	g.ui.LOSWavesAnimation(DefaultLOSRange, WaveConfusion, g.Player.Pos)
	return nil
}

func (g *game) EvokeFire(ev event) error {
	burnpos := g.Dungeon.CardinalFlammableNeighbors(g.Player.Pos)
	if len(burnpos) == 0 {
		return errors.New("You are not surrounded by any flammable terrain.")
	}
	g.Print("Sparks emanate from the magara.")
	for _, pos := range burnpos {
		g.Burn(pos, ev)
	}
	return nil
}

func (g *game) EvokeObstruction(ev event) error {
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
			g.MagicalBarrierAt(pos, ev)
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

func (g *game) MagicalBarrierAt(pos position, ev event) {
	if g.Dungeon.Cell(pos).T == WallCell || g.Dungeon.Cell(pos).T == BarrierCell {
		return
	}
	g.UpdateKnowledge(pos, g.Dungeon.Cell(pos).T)
	g.CreateMagicalBarrierAt(pos, ev)
	g.ComputeLOS()
}

func (g *game) CreateMagicalBarrierAt(pos position, ev event) {
	t := g.Dungeon.Cell(pos).T
	g.Dungeon.SetCell(pos, BarrierCell)
	delete(g.Clouds, pos)
	g.MagicalBarriers[pos] = t
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationMagicalBarrier + RandInt(DurationMagicalBarrier/2), Pos: pos, EAction: ObstructionEnd})
}
