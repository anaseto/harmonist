package main

import (
	"errors"
	"fmt"

	"github.com/anaseto/gruid"
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
	ParalysisMagara
	ObstructionMagara
	LignificationMagara
	EnergyMagara
	TransparencyMagara
	DisguiseMagara
	DelayedNoiseMagara
	DispersalMagara
	DelayedOricExplosionMagara
	//BarrierMagara
)

func (mag magara) Harmonic() bool {
	switch mag.Kind {
	case FogMagara, ShadowsMagara, NoiseMagara, ConfusionMagara,
		SleepingMagara, ParalysisMagara, TransparencyMagara, DisguiseMagara,
		DelayedNoiseMagara:
		return true
	default:
		return false
	}
}

func (mag magara) Oric() bool {
	switch mag.Kind {
	case BlinkMagara, DigMagara, TeleportMagara, LevitationMagara,
		TeleportOtherMagara, SwappingMagara, ObstructionMagara,
		DispersalMagara, DelayedOricExplosionMagara:
		return true
	default:
		return false
	}
}

func (m magaraKind) DefaultCharges() int {
	switch m {
	case LevitationMagara, FogMagara, NoiseMagara, FireMagara, DelayedNoiseMagara:
		return 6
	case ParalysisMagara, ShadowsMagara, ObstructionMagara, TransparencyMagara, DispersalMagara, DelayedOricExplosionMagara:
		return 5
	case EnergyMagara:
		return 1
	default:
		return 4
	}
}

func (g *game) RandomStartingMagara() magara {
	mags := []magaraKind{BlinkMagara, DigMagara, TeleportMagara,
		SwiftnessMagara, LevitationMagara, FireMagara, FogMagara,
		ShadowsMagara, NoiseMagara, ConfusionMagara, SleepingMagara,
		TeleportOtherMagara, SwappingMagara, ParalysisMagara,
		ObstructionMagara, LignificationMagara, TransparencyMagara,
		DelayedNoiseMagara, DisguiseMagara}
	var mag magaraKind
loop:
	for {
		mag = mags[RandInt(len(mags))]
		for _, m := range g.GeneratedMagaras {
			if m == mag {
				continue loop
			}
		}
		break
	}
	return magara{Kind: mag, Charges: mag.DefaultCharges()}
}

func (g *game) RandomMagara() magara {
	mags := []magaraKind{BlinkMagara, DigMagara, TeleportMagara,
		SwiftnessMagara, LevitationMagara, FireMagara, FogMagara,
		ShadowsMagara, NoiseMagara, ConfusionMagara, SleepingMagara,
		TeleportOtherMagara, SwappingMagara, ParalysisMagara,
		ObstructionMagara, LignificationMagara, TransparencyMagara,
		DelayedNoiseMagara, DisguiseMagara, EnergyMagara}
	var mag magaraKind
loop:
	for {
		mag = mags[RandInt(len(mags))]
		for _, m := range g.GeneratedMagaras {
			if m == mag {
				continue loop
			}
		}
		break
	}
	return magara{Kind: mag, Charges: mag.DefaultCharges()}
}

func (g *game) EquipMagara(i int) (err error) {
	omagara := g.Player.Magaras[i]
	g.Player.Magaras[i] = g.Objects.Magaras[g.Player.P]
	g.Objects.Magaras[g.Player.P] = omagara
	g.Printf("You take the %s.", g.Player.Magaras[i])
	g.Printf("You leave the %s.", omagara)
	g.StoryPrintf("Took %s (%d), left %s (%d)", g.Player.Magaras[i], g.Player.Magaras[i].Charges, omagara, omagara.Charges)
	return nil
}

func (g *game) UseMagara(n int) (err error) {
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
		err = g.EvokeBlink()
	case DigMagara:
		err = g.EvokeDig()
	case TeleportMagara:
		err = g.EvokeTeleport()
	case SwiftnessMagara:
		err = g.EvokeSwiftness()
	case LevitationMagara:
		err = g.EvokeLevitation()
	case FireMagara:
		err = g.EvokeFire()
	case FogMagara:
		err = g.EvokeFog()
	case ShadowsMagara:
		err = g.EvokeShadows()
	case NoiseMagara:
		err = g.EvokeNoise()
	case ConfusionMagara:
		err = g.EvokeConfusion()
	case ParalysisMagara:
		err = g.EvokeParalysis()
	case SleepingMagara:
		err = g.EvokeSleeping()
	case TeleportOtherMagara:
		err = g.EvokeTeleportOther()
	case SwappingMagara:
		err = g.EvokeSwapping()
	case ObstructionMagara:
		err = g.EvokeObstruction()
	case LignificationMagara:
		err = g.EvokeLignification()
	case EnergyMagara:
		err = g.EvokeEnergyMagara()
	case TransparencyMagara:
		err = g.EvokeTransparencyMagara()
	case DisguiseMagara:
		err = g.EvokeDisguiseMagara()
	case DelayedNoiseMagara:
		err = g.EvokeDelayedNoiseMagara()
	case DispersalMagara:
		err = g.EvokeDispersalMagara()
	case DelayedOricExplosionMagara:
		err = g.EvokeOricExplosionMagara()
	}
	if err != nil {
		return err
	}
	g.Stats.MagarasUsed++
	g.Stats.UsedMagaras[mag.Kind]++
	g.Stats.DMagaraUses[g.Depth]++
	g.Player.MP -= mag.MPCost(g)
	g.Player.Magaras[n].Charges--
	g.StoryPrintf("Evoked %s (MP: %d, Charges: %d)", mag, g.Player.MP, g.Player.Magaras[n].Charges)
	if mag.Harmonic() {
		g.Stats.HarmonicMagUse++
		if g.Stats.HarmonicMagUse == 6 {
			AchHarmonistNovice.Get(g)
		}
		if g.Stats.HarmonicMagUse == 11 {
			AchHarmonistInitiate.Get(g)
		}
		if g.Stats.HarmonicMagUse == 16 {
			AchHarmonistMaster.Get(g)
		}
	} else if mag.Oric() {
		g.Stats.OricMagUse++
		if g.Stats.OricMagUse == 6 {
			AchNoviceOricCelmist.Get(g)
		}
		if g.Stats.OricMagUse == 11 {
			AchInitiateOricCelmist.Get(g)
		}
		if g.Stats.OricMagUse == 16 {
			AchMasterOricCelmist.Get(g)
		}
	} else if mag.Kind == FireMagara {
		g.Stats.FireUse++
		if g.Stats.FireUse == 2 {
			AchPyromancerNovice.Get(g)
		}
		if g.Stats.FireUse == 4 {
			AchPyromancerInitiate.Get(g)
		}
		if g.Stats.FireUse == 6 {
			AchPyromancerMaster.Get(g)
		}
	}
	switch mag.Kind {
	case TeleportMagara, TeleportOtherMagara, BlinkMagara, SwappingMagara, DispersalMagara:
		g.Stats.OricTelUse++
		if g.Stats.OricTelUse == 14 {
			AchTeleport.Get(g)
		}
	}
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
	case ParalysisMagara:
		desc = "magara of paralysis"
	case ObstructionMagara:
		desc = "magara of obstruction"
	case LignificationMagara:
		desc = "magara of lignification"
	case EnergyMagara:
		desc = "magara of energy"
	case TransparencyMagara:
		desc = "magara of transparency"
	case DisguiseMagara:
		desc = "magara of disguise"
	case DelayedNoiseMagara:
		desc = "magara of delayed noise"
	case DispersalMagara:
		desc = "magara of dispersal"
	case DelayedOricExplosionMagara:
		desc = "magara of oric explosion"
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
		desc = "makes you able to move several times in a row for free."
	case LevitationMagara:
		desc = "makes you levitate with oric energies, allowing you to move over chasms, as well as through oric barriers."
	case FireMagara:
		desc = "throws small magical sparks at flammable terrain adjacent to you. Flammable terrain is first consumed by magical flames that are by themselves harmless to creatures. Then smoke will produce night clouds inducing sleep and confusion in monsters. As a gawalt monkey, you resist sleepiness, but you will still feel confused. The fire does often expand to other adjacent flammable terrain."
	case FogMagara:
		desc = "creates a dense fog in a 2-range radius using harmonic energies. The fog will dissipate with time."
	case ShadowsMagara:
		desc = "surrounds you by harmonic shadows. While standing on a dark cell, only adjacent monsters will be able to see you. It does not affect your visibility on lighted cells."
	case NoiseMagara:
		desc = "tricks monsters in a 12-range area with harmonic magical sounds, making them go away from you for a few turns. The possible monster destinations will be marked with noise symbols in the map. It only works on monsters that are not already seeing you."
	case ConfusionMagara:
		desc = "confuses monsters in sight with harmonic light and sounds, leaving them unable to attack you."
	case ParalysisMagara:
		desc = "makes monsters in sight unable to act by disturbing their senses with sound and light illusions."
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
	case EnergyMagara:
		desc = "replenishes your MP and HP."
	case TransparencyMagara:
		desc = "feeds surrounding light to harmonic magic to make you transparent. When standing on a lighted cell, only adjacent monsters will be able to see you. It does not affect your visibility on dark cells."
	case DisguiseMagara:
		desc = "surrounds you with harmonic illusions that make you look like a guard. As a result, most monsters will ignore you. Monsters with good flair may see through the illusions at less than 3 tiles away. Monsters that are already hunting you will continue doing so."
	case DelayedNoiseMagara:
		desc = "will produce a thunderous harmonic noise in your current cell. The noise will happen after a delay."
	case DispersalMagara:
		desc = "will make monsters that attempt to hit you blink away."
	case DelayedOricExplosionMagara:
		desc = "will produce a big rock-destroying oric explosion at your current cell. The explosion will happen after a delay. It destroys only walls."
	}
	duration := 0
	switch mag.Kind {
	case ConfusionMagara:
		duration = DurationConfusionMonster
	case ParalysisMagara:
		duration = DurationParalysisMonster
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
	case TransparencyMagara:
		duration = DurationTransparency
	case DisguiseMagara:
		duration = DurationDisguise
	case DelayedNoiseMagara:
		duration = DurationHarmonicNoiseDelay
	case DelayedOricExplosionMagara:
		duration = DurationOricExplosionDelay
	}
	if duration > 0 {
		switch mag.Kind {
		case DelayedNoiseMagara, DelayedOricExplosionMagara:
			desc += fmt.Sprintf(" Delay lasts for %d turns.", duration)
		default:
			desc += fmt.Sprintf(" Effect lasts for %d turns.", duration)
		}
	}
	desc += fmt.Sprintf("\n\nIt currently has %d charges.", mag.Charges)
	return fmt.Sprintf("The %s %s", mag, desc)
}

func (mag magara) MPCost(g *game) int {
	if mag.Kind == NoMagara || mag.Kind == EnergyMagara {
		return 0
	}
	cost := 1
	return cost
}

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
	if !valid(npos) {
		// should not happen
		g.Print("You could not blink.")
		return false
	}
	opos := g.Player.P
	if npos == opos {
		g.Print("You blink in-place.")
	} else {
		g.Print("You blink away.")
	}
	g.md.TeleportAnimation(opos, npos, true)
	g.PlacePlayerAt(npos)
	return true
}

func (g *game) BlinkPos(mpassable bool) gruid.Point {
	losPos := []gruid.Point{}
	for p, b := range g.Player.LOS {
		// XXX: skip if not seen by monster?
		if !b {
			continue
		}
		c := g.Dungeon.Cell(p)
		if !c.IsPlayerPassable() || mpassable && !c.IsPassable() || terrain(c) == StoryCell {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		losPos = append(losPos, p)
	}
	if len(losPos) == 0 {
		return invalidPos
	}
	q := losPos[RandInt(len(losPos))]
	for i := 0; i < 4; i++ {
		p := losPos[RandInt(len(losPos))]
		if distance(q, g.Player.P) < distance(p, g.Player.P) {
			q = p
		}
	}
	return q
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
	g.md.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) MonstersInLOS() []*monster {
	ms := []*monster{}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.P) {
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
		if mons.Exists() && g.Player.Sees(mons.P) && (mons.P.X == g.Player.P.X || mons.P.Y == g.Player.P.Y) {
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
		if ms[i].Search == invalidPos {
			ms[i].Search = ms[i].P
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
	g.md.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeRefillMagic() error {
	g.Player.MP = g.Player.MPMax()
	g.Print("Your magic forces return.")
	g.md.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeSwiftness() error {
	if g.Player.HasStatus(StatusSwift) {
		return errors.New("You are already swift.")
	}
	g.Player.Statuses[StatusSwift] += DurationSwiftness
	g.Printf("You feel swift.")
	g.md.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeLevitation() error {
	if !g.PutStatus(StatusLevitation, DurationLevitation) {
		return errors.New("You are already levitating.")
	}
	g.Printf("You feel light.")
	g.md.PlayerGoodEffectAnimation()
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
		if distance(m.P, g.Player.P) > best {
			best = distance(m.P, g.Player.P)
			mons = m
		}
	}
	if !mons.Exists() {
		return errors.New("No monsters suitable for swapping in view.")
	}
	if mons.Kind.CanOpenDoors() {
		// only intelligent monsters understand swapping
		mons.Search = mons.P
	}
	g.SwapWithMonster(mons)
	return nil
}

func (g *game) SwapWithMonster(mons *monster) {
	g.Printf("You swap positions with the %s.", mons.Kind)
	g.md.SwappingAnimation(mons.P, g.Player.P)
	g.PlacePlayerAt(mons.P)
	if terrain(g.Dungeon.Cell(g.Player.P)) == ChasmCell {
		g.PushEventFirst(&playerEvent{Action: AbyssFall}, g.Turn)
	}
}

type cloud int

const (
	CloudFog cloud = iota
	CloudFire
	CloudNight
)

func (cl cloud) String() (s string) {
	switch cl {
	case CloudFog:
		s = "fog"
	case CloudFire:
		s = "fire"
	case CloudNight:
		s = "purple fog"
	}
	return s
}

func (g *game) EvokeFog() error {
	g.Fog(g.Player.P, 3)
	g.Print("You are surrounded by a dense fog.")
	g.md.EffectAtPPAnimation()
	return nil
}

func (g *game) Fog(at gruid.Point, radius int) {
	dij := &noisePath{g: g}
	nodes := g.PR.DijkstraMap(dij, []gruid.Point{at}, radius)
	for _, n := range nodes {
		_, ok := g.Clouds[n.P]
		if !ok && g.Dungeon.Cell(n.P).AllowsFog() {
			g.Clouds[n.P] = CloudFog
			g.PushEvent(&posEvent{P: n.P, Action: CloudEnd}, g.Turn+DurationFog+RandInt(DurationFog/2))
		}
	}
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
	g.md.EffectAtPPAnimation()
	return nil
}

func (g *game) EvokeParalysis() error {
	ms := g.MonstersInLOS()
	count := 0
	for _, mons := range ms {
		if mons.PutStatus(g, MonsParalysed, DurationParalysisMonster) {
			count++
			if mons.Search == invalidPos {
				mons.Search = mons.P
			}
		}
	}
	if count == 0 {
		return errors.New("No suitable targets in view.")
	}
	g.Print("Whoosh! A slowing luminous wave emerges.")
	g.md.LOSWavesAnimation(DefaultLOSRange, WaveSlowing, g.Player.P)
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
	targets := []gruid.Point{}
	// XXX: maybe use noise distance instead of LOS?
	for i := 0; i < max; i++ {
		mons := ms[i]
		if mons.State != Resting {
			g.Printf("%s falls asleep.", mons.Kind.Definite(true))
		} else {
			continue
		}
		mons.State = Resting
		mons.Dir = ZP
		mons.ExhaustTime(g, 4+RandInt(2))
		targets = append(targets, g.Ray(mons.P)...)
	}
	if len(targets) == 0 {
		return errors.New("There are no suitable targets.")
	}
	if max == 1 {
		g.Print("A beam of sleeping emerges.")
	} else {
		g.Print("Two beams of sleeping emerge.")
	}
	g.md.BeamsAnimation(targets, BeamSleeping)

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
	targets := []gruid.Point{}
	for i := 0; i < max; i++ {
		mons := ms[i]
		if mons.Status(MonsLignified) || mons.Kind.ResistsLignification() {
			continue
		}
		mons.EnterLignification(g)
		if mons.Search == invalidPos {
			mons.Search = mons.P
		}
		targets = append(targets, g.Ray(mons.P)...)
	}
	if len(targets) == 0 {
		return errors.New("There are no suitable targets.")
	}
	if max == 1 {
		g.Print("A beam of lignification emerges.")
	} else {
		g.Print("Two beams of lignification emerge.")
	}
	g.md.BeamsAnimation(targets, BeamLignification)
	return nil
}

func (g *game) EvokeNoise() error {
	dij := &noisePath{g: g}
	const noiseDist = 23
	g.PR.BreadthFirstMap(dij, []gruid.Point{g.Player.P}, noiseDist)
	noises := []gruid.Point{}
	g.NoiseIllusion = map[gruid.Point]bool{}
	for _, mons := range g.Monsters {
		if !mons.Exists() {
			continue
		}
		c := g.PR.BreadthFirstMapAt(mons.P)
		if c > DefaultLOSRange {
			continue
		}
		if mons.SeesPlayer(g) {
			continue
		}
		mp := &monPath{g: g, monster: mons}
		target := mons.P
		best := c
		for {
			ncost := best
			for _, p := range mp.Neighbors(target) {
				c := g.PR.BreadthFirstMapAt(p)
				if c > noiseDist {
					continue
				}
				ncost := c
				if ncost > best {
					target = p
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
	g.md.NoiseAnimation(noises)
	g.Print("Monsters are tricked by magical sounds.")
	return nil
}

func (g *game) EvokeConfusion() error {
	ms := g.MonstersInLOS()
	count := 0
	for _, mons := range ms {
		if mons.EnterConfusion(g) {
			count++
			if mons.Search == invalidPos {
				mons.Search = mons.P
			}
		}
	}
	if count == 0 {
		return errors.New("No suitable targets in view.")
	}
	g.Print("Whoosh! A confusing luminous wave emerges.")
	g.md.LOSWavesAnimation(DefaultLOSRange, WaveConfusion, g.Player.P)
	return nil
}

func (g *game) EvokeFire() error {
	burnpos := g.flammableNeighbors(g.Player.P)
	if len(burnpos) == 0 {
		return errors.New("You are not surrounded by any flammable terrain.")
	}
	g.Print("Sparks emanate from the magara.")
	g.md.EffectAtPPAnimation()
	for _, p := range burnpos {
		g.Burn(p)
	}
	return nil
}

func (g *game) EvokeObstruction() error {
	targets := []gruid.Point{}
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.P) {
			continue
		}
		ray := g.Ray(mons.P)
		for i, p := range ray[1:] {
			if p == g.Player.P {
				break
			}
			mons := g.MonsterAt(p)
			if mons.Exists() {
				continue
			}
			g.MagicalBarrierAt(p)
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
	g.md.BeamsAnimation(targets, BeamObstruction)
	return nil
}

func (g *game) MagicalBarrierAt(p gruid.Point) {
	if terrain(g.Dungeon.Cell(p)) == WallCell || terrain(g.Dungeon.Cell(p)) == BarrierCell {
		return
	}
	g.UpdateKnowledge(p, terrain(g.Dungeon.Cell(p)))
	g.CreateMagicalBarrierAt(p)
	g.ComputeLOS()
}

func (g *game) CreateMagicalBarrierAt(p gruid.Point) {
	t := terrain(g.Dungeon.Cell(p))
	g.Dungeon.SetCell(p, BarrierCell)
	delete(g.Clouds, p)
	g.MagicalBarriers[p] = t
	g.PushEvent(&posEvent{P: p, Action: ObstructionEnd}, g.Turn+DurationMagicalBarrier+RandInt(DurationMagicalBarrier/2))
}

func (g *game) EvokeEnergyMagara() error {
	if g.Player.MP == g.Player.MPMax() && g.Player.HP == g.Player.HPMax() {
		return errors.New("You are already full of energy.")
	}
	g.Print("The magara glows.")
	g.md.PlayerGoodEffectAnimation()
	g.Player.MP = g.Player.MPMax()
	g.Player.HP = g.Player.HPMax()
	return nil
}

func (g *game) EvokeTransparencyMagara() error {
	if !g.PutStatus(StatusTransparent, DurationTransparency) {
		return errors.New("You are already transparent.")
	}
	g.Print("Light makes you diaphanous.")
	g.md.EffectAtPPAnimation()
	return nil
}

func (g *game) EvokeDisguiseMagara() error {
	if !g.PutStatus(StatusDisguised, DurationDisguise) {
		return errors.New("You are already disguised.")
	}
	g.Print("You look now like a normal guard.")
	g.md.EffectAtPPAnimation()
	return nil
}

func (g *game) EvokeDelayedNoiseMagara() error {
	if !g.PutFakeStatus(StatusDelay, DurationHarmonicNoiseDelay) {
		return errors.New("You are already using delayed magic.")
	}
	g.PushEventD(&posEvent{P: g.Player.P, Action: DelayedHarmonicNoiseEvent,
		Timer: DurationHarmonicNoiseDelay}, DurationTurn)
	g.Print("Timer activated.")
	g.md.EffectAtPPAnimation()
	return nil
}

func (g *game) EvokeDispersalMagara() error {
	if !g.PutStatus(StatusDispersal, DurationDispersal) {
		return errors.New("You are already dispersing.")
	}
	g.Print("You feel unstable.")
	g.md.EffectAtPPAnimation()
	return nil
}

func (g *game) EvokeOricExplosionMagara() error {
	if !g.PutFakeStatus(StatusDelay, DurationHarmonicNoiseDelay) {
		return errors.New("You are already using delayed magic.")
	}
	g.PushEventD(&posEvent{P: g.Player.P, Action: DelayedOricExplosionEvent,
		Timer: DurationOricExplosionDelay}, DurationTurn)
	g.Print("Timer activated.")
	g.md.EffectAtPPAnimation()
	return nil
}
