package main

import (
	"fmt"
)

type monsterState int

const (
	Resting monsterState = iota
	Hunting
	Wandering
	Watching
)

func (m monsterState) String() string {
	var st string
	switch m {
	case Resting:
		st = "resting"
	case Wandering:
		st = "wandering"
	case Hunting:
		st = "hunting"
	case Watching:
		st = "watching"
	}
	return st
}

type monsterStatus int

const (
	MonsConfused monsterStatus = iota
	MonsExhausted
	MonsParalysed
	MonsSatiated
	MonsLignified
)

const NMonsStatus = int(MonsLignified) + 1

func (st monsterStatus) String() (text string) {
	switch st {
	case MonsConfused:
		text = "confused"
	case MonsExhausted:
		text = "exhausted"
	case MonsParalysed:
		text = "paralysed"
	case MonsSatiated:
		text = "satiated"
	case MonsLignified:
		text = "lignified"
	}
	return text
}

type monsterKind int

const (
	MonsGuard monsterKind = iota
	MonsYack
	MonsSatowalgaPlant
	MonsMadNixe
	MonsBlinkingFrog
	MonsWorm
	MonsMirrorSpecter
	MonsTinyHarpy
	//MonsOgre
	MonsOricCelmist
	MonsHarmonicCelmist
	//MonsBrizzia
	MonsDog
	//MonsGiantBee
	MonsHighGuard
	//MonsHydra
	//MonsSkeletonWarrior
	MonsSpider
	MonsWingedMilfid
	//MonsLich
	MonsEarthDragon
	MonsAcidMound
	MonsExplosiveNadre
	//MonsMindCelmist
	MonsVampire
	MonsTreeMushroom
	//MonsMarevorHelith
	MonsButterfly
	MonsCrazyImp
	MonsHazeCat
)

func (mk monsterKind) String() string {
	return MonsData[mk].name
}

func (mk monsterKind) Letter() rune {
	return MonsData[mk].letter
}

func (mk monsterKind) BaseAttack() int {
	return 1
}

func (mk monsterKind) Dangerousness() int {
	return MonsData[mk].dangerousness
}

func (mk monsterKind) Ranged() bool {
	switch mk {
	//case MonsLich, MonsCyclop, MonsHighGuard, MonsSatowalgaPlant, MonsMadNixe, MonsVampire, MonsTreeMushroom:
	case MonsHighGuard, MonsSatowalgaPlant, MonsMadNixe, MonsVampire, MonsTreeMushroom:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Smiting() bool {
	switch mk {
	//case MonsMirrorSpecter, MonsMindCelmist:
	case MonsMirrorSpecter, MonsOricCelmist, MonsHarmonicCelmist:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Peaceful() bool {
	switch mk {
	case MonsButterfly, MonsEarthDragon, MonsCrazyImp:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Notable() bool {
	switch mk {
	case MonsCrazyImp, MonsEarthDragon, MonsHazeCat:
		return true
	default:
		return false
	}
}

func (mk monsterKind) CanOpenDoors() bool {
	switch mk {
	case MonsGuard, MonsHighGuard, MonsMadNixe, MonsOricCelmist, MonsHarmonicCelmist, MonsVampire, MonsWingedMilfid:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Patrolling() bool {
	switch mk {
	case MonsGuard, MonsHighGuard, MonsMadNixe, MonsOricCelmist, MonsHarmonicCelmist:
		return true
	default:
		return false
	}
}

func (mk monsterKind) CanFly() bool {
	switch mk {
	case MonsWingedMilfid, MonsMirrorSpecter, MonsButterfly, MonsTinyHarpy:
		return true
	default:
		return false
	}
}

func (mk monsterKind) CanSwim() bool {
	switch mk {
	case MonsBlinkingFrog, MonsVampire, MonsDog:
		return true
	default:
		return false
	}
}

func (mk monsterKind) CanAttackOnTree() bool {
	switch {
	case mk.Size() == MonsLarge:
		return true
	case mk.CanFly():
		return true
	case mk == MonsBlinkingFrog || mk == MonsHazeCat:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Desc() string {
	return MonsDesc[mk]
}

func (mk monsterKind) Indefinite(capital bool) (text string) {
	switch mk {
	//case MonsMarevorHelith:
	//text = mk.String()
	default:
		text = Indefinite(mk.String(), capital)
	}
	return text
}

func (mk monsterKind) Definite(capital bool) (text string) {
	switch mk {
	//case MonsMarevorHelith:
	//text = mk.String()
	default:
		if capital {
			text = fmt.Sprintf("The %s", mk.String())
		} else {
			text = fmt.Sprintf("the %s", mk.String())
		}
	}
	return text
}

func (mk monsterKind) Size() monsize {
	return MonsData[mk].size
}

type monsize int

const (
	MonsSmall monsize = iota
	MonsMedium
	MonsLarge
)

func (ms monsize) String() (text string) {
	switch ms {
	case MonsSmall:
		text = "small"
	case MonsMedium:
		text = "average"
	case MonsLarge:
		text = "large"
	}
	return text
}

type monsterData struct {
	size          monsize
	letter        rune
	name          string
	dangerousness int
}

var MonsData = []monsterData{
	MonsGuard:     {MonsMedium, 'g', "guard", 3},
	MonsTinyHarpy: {MonsSmall, 't', "tiny harpy", 3},
	//MonsOgre:            {10, 2, 20, 3, 'O', "ogre", 7},
	MonsOricCelmist:     {MonsMedium, 'o', "oric celmist", 9},
	MonsHarmonicCelmist: {MonsMedium, 'h', "harmonic celmist", 9},
	MonsWorm:            {MonsSmall, 'w', "farmer worm", 4},
	//MonsBrizzia:         {15, 1, 10, 3, 'z', "brizzia", 6},
	MonsAcidMound: {MonsSmall, 'a', "acid mound", 4},
	MonsDog:       {MonsMedium, 'd', "dog", 5},
	MonsYack:      {MonsMedium, 'y', "yack", 5},
	//MonsGiantBee:        {5, 1, 10, 1, 'B', "giant bee", 6},
	MonsHighGuard: {MonsMedium, 'G', "high guard", 5},
	//MonsHydra:           {10, 1, 10, 4, 'H', "hydra", 10},
	//MonsSkeletonWarrior: {10, 1, 10, 3, 'S', "skeleton warrior", 6},
	MonsSpider:       {MonsSmall, 's', "spider", 15},
	MonsWingedMilfid: {MonsMedium, 'W', "winged milfid", 6},
	MonsBlinkingFrog: {MonsMedium, 'F', "blinking frog", 6},
	//MonsLich:           {10, MonsMedium, 'L', "lich", 15},
	MonsEarthDragon:    {MonsLarge, 'D', "earth dragon", 18},
	MonsMirrorSpecter:  {MonsMedium, 'm', "mirror specter", 11},
	MonsExplosiveNadre: {MonsMedium, 'n', "explosive nadre", 8},
	MonsSatowalgaPlant: {MonsLarge, 'P', "satowalga plant", 7},
	MonsMadNixe:        {MonsMedium, 'N', "mad nixe", 14},
	//MonsMindCelmist:     {10, 1, 20, 2, 'c', "mind celmist", 12},
	MonsVampire:      {MonsMedium, 'V', "vampire", 13},
	MonsTreeMushroom: {MonsLarge, 'T', "tree mushroom", 17},
	//MonsMarevorHelith: {10, MonsMedium, 'M', "Marevor Helith", 18},
	MonsButterfly: {MonsSmall, 'b', "kerejat", 2},
	MonsCrazyImp:  {MonsSmall, 'i', "Crazy Imp", 19},
	MonsHazeCat:   {MonsSmall, 'c', "haze cat", 16},
}

var MonsDesc = []string{
	MonsGuard:     "Guards are low rank soldiers who patrol between Dayoriah Clan's buildings.",
	MonsTinyHarpy: "Tiny harpies are little humanoid flying creatures. They are aggressive when hungry, but peaceful when satiated. This Underground harpy species eats fruits (including bananas) and other vegetables.",
	//MonsOgre:            "Ogres are big clunky humanoids that can hit really hard.",
	MonsOricCelmist:     "Oric celmists are mages that can create magical barriers in cells adjacent to you, complicating your escape.\n\nDayoriah Clan's oric celmists are famous for their knowledge of oric magic force manipulations. They are the ones who instigated the steal of Marevor's Gem Portal Artifact. According to Marevor, they plan on doing some dangerous oric experiments with the Artifact, though that's all you can say about it, because his boring explanations were a bit over your head.",
	MonsHarmonicCelmist: "Harmonic celmists are mages specialized in manipulation of sound and light. They can illuminate you with harmonic light, making it more difficult to hide from them. They also use alert harmonic sounds around you.\n\nHarmonies are usually mainly used for sneaking around in the shadows, but they can also be used to reveal ennemies, sadly for you. Although harmonies are often considered as less prestigious magic energies than oric energies, the Dayoriah Clan knows how to make good use of them, as they clearly showed when they stole Marevor's Gem Portal Artifact.",
	MonsWorm:            "Farmer worms are ugly slow moving creatures. They furrow as they move, helping new foliage to grow.",
	//MonsBrizzia:         "Brizzias are big slow moving biped creatures. They are quite hardy, and when hurt they can cause nausea, impeding the use of potions.",
	MonsAcidMound: "Acid mounds are acidic creatures. They can corrode your magaras, reducing their number of charges.",
	MonsDog:       "Dogs are fast moving carnivore quadrupeds. They can bark, and smell you.",
	MonsYack:      "Yacks are quite large herbivorous quadrupeds. They tend to eat grass peacefully, but upon seing you they may attack, pushing you up to 5 cells away.",
	//MonsGiantBee:        "Giant bees are fragile but extremely fast moving creatures. Their bite can sometimes enrage you.",
	MonsHighGuard: "High guards watch over a particular location. They can throw javelins.",
	//MonsHydra:           "Hydras are enormous creatures with four heads that can hit you each at once.",
	//MonsSkeletonWarrior: "Skeleton warriors are good fighters, clad in chain mail.",
	MonsSpider:       "Spiders are small creatures, with panoramic vision and whose bite can confuse you.",
	MonsWingedMilfid: "Winged milfids are fast moving humanoids that can fly over you and make you swap positions. They tend to be very agressive creatures.",
	MonsBlinkingFrog: "Blinking frogs are big frog-like creatures, whose bite can make you blink away. The science behind their attack is not clear, but many think it relies on some kind of oric deviation magic. They can jump to attack from below.",
	//MonsLich:           "Liches are non-living mages wearing a leather armour. They can throw a bolt of torment at you, halving your HP.",
	MonsEarthDragon:    "Earth dragons are big creatures from a dragon species that wander in the Underground. They are peaceful creatures, but they may hurt you inadvertently, pushing you up to 6 tiles away (3 if confused). They naturally emit powerful oric energies, allowing them to eat rocks and dig tunnels. Their oric energies can confuse you if you're close enough, for example if they hurt you or you jump over them.",
	MonsMirrorSpecter:  "Mirror specters are very insubstantial creatures, which can absorb your mana.",
	MonsExplosiveNadre: "Nadres are dragon-like biped creatures that are famous for exploding upon dying. Explosive nadres are a tiny nadre race that explodes upon attacking. The explosion confuses any adjacent creatures and occasionally destroys walls.",
	MonsSatowalgaPlant: "Satowalga Plants are immobile bushes that throw slowing viscous acidic projectiles at you, destroying some of your magara charges. They attack at half normal speed.",
	MonsMadNixe:        "Nixes are magical humanoids. Usually, they specialize in illusion harmonic magic, but the so called mad nixes are a perverted variant who learned the oric arts to create a spell that can attract their foes to them, so that they can kill them without pursuing them.",
	//MonsMindCelmist:     "Mind celmists are mages that use magical smitting mind attacks that bypass armour. They can occasionally confuse or slow you. They try to avoid melee.",
	MonsVampire:      "Vampires are humanoids that drink blood to survive. Their nauseous spitting can cause confusion, impeding the use of magaras for a few turns.",
	MonsTreeMushroom: "Tree mushrooms are big clunky slow-moving creatures. They can throw lignifying spores at you, leaving you unable to move for a few turns, though the spores will also provide some protection against harm.",
	//MonsMarevorHelith: "Marevor Helith is an ancient undead nakrus very fond of teleporting people away. He is a well-known expert in the field of magaras - items that many people simply call magical objects. His current research focus is monolith creation. Marevor, a repentant necromancer, is now searching for his old disciple Jaixel in the Underground to help him overcome the past.",
	MonsButterfly: "Underground's butterflies, called kerejats, wander peacefully around, illuminating their surroundings.",
	MonsCrazyImp:  "Crazy Imp is a crazy creature that likes to sing with its small guitar. It seems to be fond of monkeys and quite capable at finding them by flair. While singing it may attract unwanted attention.",
	MonsHazeCat:   "Haze cats are a special variety of cats found in the Underground. They have very good night vision and are always alert.",
}

type bandInfo struct {
	Path []position
	I    int
	Kind monsterBand
	Beh  mbehaviour
}

type monsterBand int

const (
	LoneGuard monsterBand = iota
	LoneHighGuard
	LoneYack
	LoneOricCelmist
	LoneHarmonicCelmist
	LoneSatowalgaPlant
	LoneBlinkingFrog
	LoneWorm
	LoneMirrorSpecter
	LoneDog
	LoneExplosiveNadre
	LoneWingedMilfid
	LoneMadNixe
	LoneTreeMushroom
	LoneEarthDragon
	LoneButterfly
	LoneVampire
	LoneHarpy
	LoneHazeCat
	LoneAcidMound
	LoneSpider
	PairGuard
	PairYack
	PairFrog
	PairDog
	PairTreeMushroom
	PairSatowalga
	PairWorm
	PairOricCelmist
	PairHarmonicCelmist
	PairVampire
	PairNixe
	PairExplosiveNadre
	PairWingedMilfid
	SpecialLoneVampire
	SpecialLoneNixe
	SpecialLoneMilfid
	SpecialLoneOricCelmist
	SpecialArtifactBand
	SpecialLoneHarmonicCelmist
	SpecialLoneHighGuard
	SpecialLoneHarpy
	SpecialLoneTreeMushroom
	SpecialLoneMirrorSpecter
	SpecialLoneAcidMound
	SpecialLoneHazeCat
	SpecialLoneSpider
	SpecialLoneBlinkingFrog
	SpecialLoneExplosiveNadre
	SpecialLoneYack
	SpecialLoneDog
	UniqueCrazyImp
)

type monsterBandData struct {
	Distribution map[monsterKind]int
	Band         bool
	Monster      monsterKind
}

var MonsBands = []monsterBandData{
	LoneGuard:                  {Monster: MonsGuard},
	LoneHighGuard:              {Monster: MonsHighGuard},
	LoneYack:                   {Monster: MonsYack},
	LoneOricCelmist:            {Monster: MonsOricCelmist},
	LoneHarmonicCelmist:        {Monster: MonsHarmonicCelmist},
	LoneSatowalgaPlant:         {Monster: MonsSatowalgaPlant},
	LoneBlinkingFrog:           {Monster: MonsBlinkingFrog},
	LoneWorm:                   {Monster: MonsWorm},
	LoneMirrorSpecter:          {Monster: MonsMirrorSpecter},
	LoneDog:                    {Monster: MonsDog},
	LoneExplosiveNadre:         {Monster: MonsExplosiveNadre},
	LoneWingedMilfid:           {Monster: MonsWingedMilfid},
	LoneMadNixe:                {Monster: MonsMadNixe},
	LoneTreeMushroom:           {Monster: MonsTreeMushroom},
	LoneEarthDragon:            {Monster: MonsEarthDragon},
	LoneButterfly:              {Monster: MonsButterfly},
	LoneVampire:                {Monster: MonsVampire},
	LoneHarpy:                  {Monster: MonsTinyHarpy},
	LoneHazeCat:                {Monster: MonsHazeCat},
	LoneAcidMound:              {Monster: MonsAcidMound},
	LoneSpider:                 {Monster: MonsSpider},
	PairGuard:                  {Band: true, Distribution: map[monsterKind]int{MonsGuard: 2}},
	PairYack:                   {Band: true, Distribution: map[monsterKind]int{MonsYack: 2}},
	PairFrog:                   {Band: true, Distribution: map[monsterKind]int{MonsBlinkingFrog: 2}},
	PairDog:                    {Band: true, Distribution: map[monsterKind]int{MonsDog: 2}},
	PairTreeMushroom:           {Band: true, Distribution: map[monsterKind]int{MonsTreeMushroom: 2}},
	PairSatowalga:              {Band: true, Distribution: map[monsterKind]int{MonsSatowalgaPlant: 2}},
	PairWorm:                   {Band: true, Distribution: map[monsterKind]int{MonsWorm: 2}},
	PairVampire:                {Band: true, Distribution: map[monsterKind]int{MonsVampire: 2}},
	PairOricCelmist:            {Band: true, Distribution: map[monsterKind]int{MonsOricCelmist: 2}},
	PairHarmonicCelmist:        {Band: true, Distribution: map[monsterKind]int{MonsHarmonicCelmist: 2}},
	PairNixe:                   {Band: true, Distribution: map[monsterKind]int{MonsMadNixe: 2}},
	PairExplosiveNadre:         {Band: true, Distribution: map[monsterKind]int{MonsExplosiveNadre: 2}},
	PairWingedMilfid:           {Band: true, Distribution: map[monsterKind]int{MonsWingedMilfid: 2}},
	SpecialArtifactBand:        {Band: true, Distribution: map[monsterKind]int{MonsOricCelmist: 1, MonsHighGuard: 1}},
	SpecialLoneVampire:         {Monster: MonsVampire},
	SpecialLoneNixe:            {Monster: MonsMadNixe},
	SpecialLoneMilfid:          {Monster: MonsWingedMilfid},
	SpecialLoneOricCelmist:     {Monster: MonsOricCelmist},
	SpecialLoneHarmonicCelmist: {Monster: MonsHarmonicCelmist},
	SpecialLoneHighGuard:       {Monster: MonsHighGuard},
	SpecialLoneHarpy:           {Monster: MonsTinyHarpy},
	SpecialLoneTreeMushroom:    {Monster: MonsTreeMushroom},
	SpecialLoneMirrorSpecter:   {Monster: MonsMirrorSpecter},
	SpecialLoneAcidMound:       {Monster: MonsAcidMound},
	SpecialLoneHazeCat:         {Monster: MonsHazeCat},
	SpecialLoneSpider:          {Monster: MonsSpider},
	SpecialLoneBlinkingFrog:    {Monster: MonsBlinkingFrog},
	SpecialLoneDog:             {Monster: MonsDog},
	SpecialLoneExplosiveNadre:  {Monster: MonsExplosiveNadre},
	SpecialLoneYack:            {Monster: MonsYack},
	UniqueCrazyImp:             {Monster: MonsCrazyImp},
}

type monster struct {
	Kind          monsterKind
	Band          int
	Index         int
	Dir           direction
	Attack        int
	Dead          bool
	State         monsterState
	Statuses      [NMonsStatus]int
	Pos           position
	LastKnownPos  position
	Target        position
	Path          []position // cache
	FireReady     bool
	Seen          bool
	LOS           map[position]bool
	LastSeenState monsterState
	Swapped       bool
	Watching      int
	Left          bool
	Search        position
	Alerted       bool
	Waiting       int
}

func (m *monster) Init() {
	m.Attack = m.Kind.BaseAttack()
	m.Pos = InvalidPos
	m.LOS = map[position]bool{}
	m.LastKnownPos = InvalidPos
	m.Search = InvalidPos
	if RandInt(2) == 0 {
		m.Left = true
	}
	switch m.Kind {
	case MonsButterfly:
		m.MakeWander()
	case MonsSatowalgaPlant:
		m.StartWatching()
	}
}

func (m *monster) StartWatching() {
	m.State = Watching
	m.Watching = 0
}

func (m *monster) Status(st monsterStatus) bool {
	return m.Statuses[st] > 0
}

func (m *monster) Exists() bool {
	return m != nil && !m.Dead
}

func (m *monster) Alternate() {
	if m.Left {
		if RandInt(4) > 0 {
			m.Dir = m.Dir.Left()
		} else {
			m.Dir = m.Dir.Right()
			m.Left = false
		}
	} else {
		if RandInt(3) > 0 {
			m.Dir = m.Dir.Right()
		} else {
			m.Dir = m.Dir.Left()
			m.Left = true
		}
	}
}

func (m *monster) TeleportAway(g *game) {
	pos := m.Pos
	i := 0
	count := 0
	for {
		count++
		if count > 1000 {
			panic("TeleportOther")
		}
		pos = g.FreePassableCell()
		if pos.Distance(m.Pos) < 15 && i < 1000 {
			i++
			continue
		}
		break
	}

	switch m.State {
	case Hunting:
		// TODO: change the target or state?
	case Resting, Wandering:
		m.MakeWander()
		m.Target = m.Pos
	}
	if g.Player.Sees(m.Pos) {
		g.Printf("%s teleports away.", m.Kind.Definite(true))
	}
	opos := m.Pos
	m.MoveTo(g, pos)
	if g.Player.Sees(opos) {
		g.ui.TeleportAnimation(opos, pos, false)
	}
}

func (m *monster) MoveTo(g *game, pos position) {
	if g.Player.Sees(pos) {
		m.UpdateKnowledge(g, pos)
	} else if g.Player.Sees(m.Pos) {
		delete(g.LastMonsterKnownAt, m.Pos)
		m.LastKnownPos = InvalidPos
	}
	if !g.Player.Sees(m.Pos) && g.Player.Sees(pos) {
		if !m.Seen {
			m.Seen = true
			g.Printf("%s (%v) comes into view.", m.Kind.Indefinite(true), m.State)
		}
		g.StopAuto()
	}
	recomputeLOS := g.Player.Sees(m.Pos) && g.Dungeon.Cell(m.Pos).T == DoorCell ||
		g.Player.Sees(pos) && g.Dungeon.Cell(pos).T == DoorCell
	m.PlaceAt(g, pos)
	if recomputeLOS {
		g.ComputeLOS()
	}
	c := g.Dungeon.Cell(pos)
	if c.T == ChasmCell && !m.Kind.CanFly() || c.T == WaterCell && !m.Kind.CanSwim() && !m.Kind.CanFly() {
		m.Dead = true
		if g.Player.Sees(m.Pos) {
			g.HandleKill(m)
			switch c.T {
			case ChasmCell:
				g.Printf("%s falls into the abyss.", m.Kind.Definite(true))
			case WaterCell:
				g.Printf("%s drowns.", m.Kind.Definite(true))
			}
		}
	}
}

func (m *monster) PlaceAt(g *game, pos position) {
	if !m.Pos.valid() {
		m.Pos = pos
		g.MonstersPosCache[m.Pos.idx()] = m.Index + 1
		npos := m.RandomFreeNeighbor(g)
		if npos != m.Pos {
			m.Dir = npos.Dir(m.Pos)
		} else {
			m.Dir = E
		}
		return
	}
	if pos == m.Pos {
		// should not happen
		return
	}
	m.Dir = pos.Dir(m.Pos)
	m.CorrectDir()
	g.MonstersPosCache[m.Pos.idx()] = 0
	m.Pos = pos
	g.MonstersPosCache[m.Pos.idx()] = m.Index + 1
	m.Waiting = 0
}

func (m *monster) CorrectDir() {
	switch m.Dir {
	case ENE, ESE:
		m.Dir = E
	case NNE, NNW:
		m.Dir = N
	case WNW, WSW:
		m.Dir = W
	case SSW, SSE:
		m.Dir = S
	}
}

func (m *monster) AttackAction(g *game, ev event) {
	m.Dir = g.Player.Pos.Dir(m.Pos)
	m.CorrectDir()
	switch m.Kind {
	case MonsExplosiveNadre:
		g.StoryPrint("Nadre explosion")
		m.Explode(g, ev)
		return
	default:
		m.HitPlayer(g, ev)
	}
	ev.Renew(g, DurationTurn)
}

func (m *monster) Explode(g *game, ev event) {
	m.Dead = true
	neighbors := m.Pos.ValidCardinalNeighbors()
	g.Printf("%s %s explodes with a loud boom.", g.ExplosionSound(), m.Kind.Definite(true))
	g.ui.ExplosionAnimation(FireExplosion, m.Pos)
	g.MakeNoise(ExplosionNoise, m.Pos)
	for _, pos := range append(neighbors, m.Pos) {
		c := g.Dungeon.Cell(pos)
		if c.Flammable() {
			g.Burn(pos, ev)
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() && !mons.Status(MonsConfused) {
			mons.EnterConfusion(g, ev)
			if mons.State != Hunting && mons.State != Watching {
				mons.StartWatching()
			}
		} else if g.Player.Pos == pos {
			m.InflictDamage(g, 1, 1)
		} else if c.IsDestructible() && RandInt(3) > 0 {
			if c.T.IsDiggable() {
				g.Dungeon.SetCell(pos, RubbleCell)
			} else {
				g.Dungeon.SetCell(pos, GroundCell)
			}
			if c.T == BarrelCell {
				delete(g.Objects.Barrels, pos)
			}
			g.Stats.Digs++
			g.UpdateKnowledge(pos, c.T)
			if g.Player.Sees(pos) {
				g.ui.WallExplosionAnimation(pos)
			}
			g.MakeNoise(WallNoise, pos)
			g.Fog(pos, 1, ev)
		}
	}
}

func (m *monster) NaturalAwake(g *game) {
	m.Target = m.NextTarget(g)
	switch g.Bands[m.Band].Beh {
	case BehGuard:
		m.StartWatching()
	default:
		m.MakeWander()
	}
	m.GatherBand(g)
}

func (m *monster) CanPass(g *game, pos position) bool {
	if !pos.valid() {
		return false
	}
	c := g.Dungeon.Cell(pos)
	return c.IsPassable() || c.IsDoorPassable() && m.Kind.CanOpenDoors() ||
		c.IsLevitatePassable() && m.Kind.CanFly() ||
		c.IsSwimPassable() && (m.Kind.CanSwim() || m.Kind.CanFly()) ||
		c.T == HoledWallCell && m.Kind.Size() == MonsSmall
}

func (m *monster) RandomFreeNeighbor(g *game) position {
	pos := m.Pos
	neighbors := [4]position{pos.E(), pos.W(), pos.N(), pos.S()}
	fnb := []position{}
	for _, nbpos := range neighbors {
		if !nbpos.valid() {
			continue
		}
		c := g.Dungeon.Cell(nbpos)
		if c.IsPassable() {
			fnb = append(fnb, nbpos)
		}
	}
	if len(fnb) == 0 {
		return m.Pos
	}
	samedir := fnb[RandInt(len(fnb))]
	for _, pos := range fnb {
		if m.Dir.InViewCone(m.Pos, pos.To(pos.Dir(m.Pos))) {
			samedir = pos
			break
		}
	}
	if RandInt(4) > 0 {
		return samedir
	}
	return fnb[RandInt(len(fnb))]
}

type mbehaviour int

const (
	BehPatrol mbehaviour = iota
	BehGuard
	BehWander
	BehExplore
	BehCrazyImp
)

var SearchAroundCache []position

func (m *monster) SearchAround(g *game, pos position, radius int) position {
	dij := &monPath{game: g, monster: m}
	nm := Dijkstra(dij, []position{pos}, radius)
	SearchAroundCache = SearchAroundCache[:0]
	nm.iter(pos, func(n *node) {
		SearchAroundCache = append(SearchAroundCache, n.Pos)
	})
	if len(SearchAroundCache) > 0 {
		p := SearchAroundCache[RandInt(len(SearchAroundCache))]
		return p
	}
	return InvalidPos
}

func (m *monster) NextTarget(g *game) (pos position) {
	band := g.Bands[m.Band]
	pos = m.Pos
	switch band.Beh {
	case BehWander:
		if m.Pos.Distance(band.Path[0]) < 8+RandInt(8) {
			pos = m.SearchAround(g, m.Pos, 4)
			if pos != InvalidPos {
				break
			}
		}
		if m.Search != InvalidPos && RandInt(2) == 0 {
			pos = m.SearchAround(g, m.Search, 7)
			if pos != InvalidPos {
				break
			}
		}
		pos = m.SearchAround(g, band.Path[0], 7)
		if pos != InvalidPos {
			break
		}
		pos = band.Path[0]
	case BehExplore:
		if m.Kind.CanOpenDoors() {
			if m.Search != InvalidPos && RandInt(4) == 0 {
				pos = m.SearchAround(g, m.Search, 7)
			} else {
				pos = m.SearchAround(g, pos, 5)
			}
			if pos != InvalidPos {
				break
			}
		}
		pos = band.Path[RandInt(len(band.Path))]
	case BehGuard:
		if m.Search != InvalidPos && m.Search.Distance(m.Pos) < 4 && RandInt(2) == 0 {
			pos = m.SearchAround(g, m.Search, 3)
			if pos != InvalidPos {
				break
			}
		}
		pos = band.Path[0]
	case BehPatrol:
		if m.Search != InvalidPos && RandInt(4) > 0 {
			pos = m.SearchAround(g, m.Search, 7)
			if pos != m.Pos && pos != InvalidPos {
				break
			}
		}
		if band.Path[0] == m.Target {
			pos = band.Path[1]
		} else if band.Path[1] == m.Target {
			pos = band.Path[0]
		} else if band.Path[0].Distance(m.Pos) < band.Path[1].Distance(m.Pos) && RandInt(4) > 0 {
			pos = band.Path[0]
		} else {
			pos = band.Path[1]
		}
	case BehCrazyImp:
		path := m.APath(g, m.Pos, g.Player.Pos)
		if len(path) == 0 {
			pos = m.SearchAround(g, m.Pos, 3)
		} else {
			pos = g.Player.Pos
		}
	}
	return pos
}

func (m *monster) HandleMonsSpecifics(g *game) (done bool) {
	switch m.Kind {
	case MonsSatowalgaPlant:
		switch m.State {
		case Hunting:
			if !m.SeesPlayer(g) {
				m.Alternate()
				if RandInt(5) == 0 {
					m.StartWatching()
				}
			}
		default:
			if RandInt(4) > 0 {
				m.Alternate()
			}
		}
		// oklob plants are static ranged-only
		g.Ev.Renew(g, DurationTurn)
		return true
	case MonsGuard, MonsHighGuard:
		if m.State != Wandering && m.State != Watching {
			break
		}
		for pos, on := range g.Objects.Lights {
			if !on && pos == m.Pos {
				g.Dungeon.SetCell(m.Pos, LightCell)
				g.Objects.Lights[m.Pos] = true
				g.Ev.Renew(g, DurationTurn)
				if g.Player.Sees(m.Pos) {
					g.Printf("%s makes a new fire.", m.Kind.Definite(true))
				} else {
					g.UpdateKnowledge(m.Pos, ExtinguishedLightCell)
				}
				return true
			} else if !on && m.SeesLight(g, pos) {
				m.Target = pos
			}
		}
	case MonsCrazyImp:
		if g.Player.Sees(m.Pos) && RandInt(2) == 0 && !m.Status(MonsConfused) && !m.Status(MonsExhausted) {
			g.PrintStyled("Crazy Imp: “♫ larilon, larila ♫ ♪”", logSpecial)
			g.MakeNoise(SingingNoise, m.Pos)
			//g.ui.MusicAnimation(m.Pos)
			m.Exhaust(g)
		}
	}
	return false
}

func (m *monster) HandleWatching(g *game) {
	turns := 4
	if m.Kind == MonsHazeCat {
		turns = 3
	}
	if m.Watching+RandInt(2) < turns {
		m.Alternate()
		m.Watching++
		if m.Kind == MonsDog {
			dij := &monPath{game: g, monster: m}
			nm := Dijkstra(dij, []position{m.Pos}, 5)
			if _, ok := nm.at(g.Player.Pos); ok {
				m.Target = g.Player.Pos
				m.MakeWander()
			}
		}
	} else {
		// pick a random cell: more escape strategies for the player
		m.Target = m.NextTarget(g)
		switch g.Bands[m.Band].Beh {
		case BehGuard:
			m.Alternate()
			if m.Pos != m.Target {
				m.MakeWander()
				m.GatherBand(g)
			}
		default:
			m.MakeWander()
			m.GatherBand(g)
		}
	}
	g.Ev.Renew(g, DurationTurn)
}

func (m *monster) ComputePath(g *game) {

	if !(len(m.Path) > 0 && m.Path[0] == m.Target && m.Path[len(m.Path)-1] == m.Pos) {
		m.Path = m.APath(g, m.Pos, m.Target)
		if len(m.Path) == 0 && !m.Status(MonsConfused) {
			// if target is not accessible, try free neighbor cells
			for _, npos := range g.Dungeon.FreeNeighbors(m.Target) {
				m.Path = m.APath(g, m.Pos, npos)
				if len(m.Path) > 0 {
					m.Target = npos
					break
				}
			}
		}
	}
}

func (m *monster) Peaceful(g *game) bool {
	if m.Kind.Peaceful() {
		return true
	}
	switch m.Kind {
	case MonsTinyHarpy:
		if m.Status(MonsSatiated) || g.Player.Bananas == 0 {
			return true
		}
	}
	return false
}

func (m *monster) HandleEndPath(g *game) {
	if len(m.Path) == 0 && m.Search != InvalidPos && m.Search.Distance(m.Target) < 10 && m.Pos != m.Target {
		// the cell where the player was last noticed may not be recheable for the monster
		m.Search = InvalidPos
	}
	switch m.State {
	case Wandering, Hunting:
		if !m.Peaceful(g) {
			if !m.SeesPlayer(g) {
				m.StartWatching()
				m.Alternate()
			}
		} else {
			m.Target = m.NextTarget(g)
		}
	}
	g.Ev.Renew(g, DurationTurn)
}

func (m *monster) MakeWanderAt(target position) {
	m.Target = target
	if m.Kind == MonsSatowalgaPlant {
		m.State = Hunting
	} else {
		m.State = Wandering
	}
}

func (m *monster) MakeWander() {
	if m.Kind == MonsSatowalgaPlant {
		m.State = Watching
	} else {
		m.State = Wandering
	}
}

func (m *monster) HandleMove(g *game) {
	target := m.Path[len(m.Path)-2]
	mons := g.MonsterAt(target)
	monstarget := InvalidPos
	if mons.Exists() && len(mons.Path) >= 2 {
		monstarget = mons.Path[len(mons.Path)-2]
	}
	c := g.Dungeon.Cell(target)
	switch {
	case m.Peaceful(g) && target == g.Player.Pos:
		switch m.Kind {
		case MonsEarthDragon:
			m.AttackAction(g, g.Ev)
			return
		default:
			m.Path = m.APath(g, m.Pos, m.Target)
		}
	case !mons.Exists():
		if m.Kind == MonsEarthDragon && c.IsDestructible() {
			g.Dungeon.SetCell(target, RubbleCell)
			if c.T == BarrelCell {
				delete(g.Objects.Barrels, target)
			}
			g.Stats.Digs++
			g.UpdateKnowledge(target, c.T)
			g.MakeNoise(WallNoise, m.Pos)
			g.Fog(m.Pos, 1, g.Ev)
			if g.Player.Pos.Distance(target) < 12 {
				// XXX use dijkstra distance ?
				if c.IsWall() {
					g.Printf("%s You hear an earth-splitting noise.", g.CrackSound())
				} else if c.T == BarrelCell || c.T == DoorCell || c.T == TableCell {
					g.Printf("%s You hear an wood-splitting noise.", g.CrackSound())
				}
				g.StopAuto()
			}
			m.MoveTo(g, target)
			m.Path = m.Path[:len(m.Path)-1]
		} else if !m.CanPass(g, target) {
			m.Path = m.APath(g, m.Pos, m.Target)
		} else {
			m.InvertFoliage(g)
			m.MoveTo(g, target)
			if (m.Kind.Ranged() || m.Kind.Smiting()) && !m.FireReady && m.SeesPlayer(g) {
				m.FireReady = true
			}
			m.Path = m.Path[:len(m.Path)-1]
		}
	case (mons.Pos == target && m.Pos == monstarget || m.Waiting > 5+RandInt(2)) && !mons.Status(MonsLignified):
		target := mons.Pos
		monstarget := m.Pos
		m.MoveTo(g, target)
		m.Path = m.Path[:len(m.Path)-1]
		mons.MoveTo(g, monstarget)
		if len(mons.Path) > 0 {
			mons.Path = mons.Path[:len(mons.Path)-1]
		}
		g.MonstersPosCache[m.Pos.idx()] = m.Index + 1
		mons.Swapped = true
	case m.State == Hunting && mons.State != Hunting:
		if m.Waiting > 2+RandInt(3) {
			if mons.Peaceful(g) {
				mons.MakeWander()
			} else {
				mons.MakeWanderAt(m.Target)
				mons.GatherBand(g)
			}
		} else {
			m.Path = m.APath(g, m.Pos, m.Target)
		}
		m.Waiting++
	case !mons.SeesPlayer(g) && mons.State != Hunting:
		if m.Waiting > 1+RandInt(2) && mons.Kind != MonsSatowalgaPlant {
			mons.MakeWanderAt(mons.RandomFreeNeighbor(g))
		} else {
			m.Path = m.APath(g, m.Pos, m.Target)
		}
		m.Waiting++
	default:
		m.Path = m.APath(g, m.Pos, m.Target)
		m.Waiting++
	}
	g.Ev.Renew(g, DurationTurn)
}

func (m *monster) HandleTurn(g *game, ev event) {
	if m.Status(MonsParalysed) {
		ev.Renew(g, DurationTurn)
		return
	}
	if m.Swapped {
		m.Swapped = false
		ev.Renew(g, DurationTurn)
		return
	}
	ppos := g.Player.Pos
	mpos := m.Pos
	switch m.Kind {
	case MonsGuard, MonsHighGuard:
		// they have to put lights on, could be optimized (TODO)
		m.ComputeLOS(g)
	}
	m.MakeAware(g)
	if m.State == Resting {
		if RandInt(3000) == 0 || (m.Kind == MonsCrazyImp || m.Kind == MonsHazeCat) && RandInt(20) == 0 {
			m.NaturalAwake(g)
		}
		ev.Renew(g, DurationTurn)
		return
	}
	if m.State == Hunting && m.RangedAttack(g, ev) {
		return
	}
	if m.State == Hunting && m.SmitingAttack(g, ev) {
		return
	}
	if m.HandleMonsSpecifics(g) {
		return
	}
	if mpos.Distance(ppos) == 1 && g.Dungeon.Cell(ppos).T != BarrelCell && !m.Peaceful(g) {
		if m.Status(MonsConfused) {
			g.Printf("%s appears too confused to attack.", m.Kind.Definite(true))
			ev.Renew(g, DurationTurn) // wait
			return
		}
		if g.Dungeon.Cell(ppos).T == TreeCell && !m.Kind.CanAttackOnTree() {
			g.Printf("%s watches you from below.", m.Kind.Definite(true))
			ev.Renew(g, DurationTurn) // wait
			return
		}
		m.AttackAction(g, ev)
		return
	}
	if m.Status(MonsLignified) {
		ev.Renew(g, DurationTurn) // wait
		return
	}
	switch m.State {
	case Watching:
		m.HandleWatching(g)
		return
	}
	m.ComputePath(g)
	if len(m.Path) < 2 {
		m.HandleEndPath(g)
		return
	}
	m.HandleMove(g)
}

func (m *monster) InvertFoliage(g *game) {
	if m.Kind != MonsWorm {
		return
	}
	invert := false
	c := g.Dungeon.Cell(m.Pos)
	if c.T == CavernCell {
		g.Dungeon.SetCell(m.Pos, FoliageCell)
		invert = true
	} else if c.T == FoliageCell {
		g.Dungeon.SetCell(m.Pos, CavernCell)
		invert = true
	}
	if invert {
		if g.Player.Sees(m.Pos) {
			g.ComputeLOS()
		} else {
			g.UpdateKnowledge(m.Pos, c.T)
		}
	}
}

func (m *monster) Exhaust(g *game) {
	m.ExhaustTime(g, DurationExhaustionMonster+RandInt(DurationExhaustionMonster/2))
}

func (m *monster) ExhaustTime(g *game, t int) {
	m.PutStatus(g, MonsExhausted, t)
}

func (m *monster) HitPlayer(g *game, ev event) {
	if g.Player.HP <= 0 || g.Player.Pos.Distance(m.Pos) > 1 {
		return
	}
	dmg := m.Attack
	clang := RandInt(4) == 0
	noise := g.HitNoise(clang)
	g.MakeNoise(noise, g.Player.Pos)
	var sclang string
	if clang {
		sclang = g.ClangMsg()
	}
	g.PrintfStyled("%s hits you (%d dmg).%s", logMonsterHit, m.Kind.Definite(true), dmg, sclang)
	m.InflictDamage(g, dmg, m.Attack)
	if g.Player.HP <= 0 {
		return
	}
	m.HitSideEffects(g, ev)
	const HeavyWoundHP = 2
	if g.Player.HP >= HeavyWoundHP {
		return
	}
	switch g.Player.Inventory.Neck {
	case AmuletConfusion:
		m.EnterConfusion(g, ev)
		// TODO: maybe affect all monsters in sight?
		g.Printf("Your amulet releases confusing harmonies.")
	case AmuletFog:
		g.Print("Your amulet feels warm.")
		g.SwiftFog(ev)
	case AmuletObstruction:
		opos := m.Pos
		m.Blink(g)
		if opos != m.Pos {
			g.MagicalBarrierAt(opos, ev)
			g.Print("The amulet releases an oric wind.")
			m.Exhaust(g)
		}
	case AmuletTeleport:
		g.Print("Your amulet shines.")
		m.TeleportAway(g)
	case AmuletLignification:
		g.Print("Your amulet glows.")
		m.EnterLignification(g, ev)
	}
}

func (m *monster) PutStatus(g *game, st monsterStatus, duration int) bool {
	if m.Status(st) {
		return false
	}
	m.Statuses[st] += duration
	g.PushEvent(&monsterEvent{
		ERank:   g.Ev.Rank() + DurationStatusStep,
		NMons:   m.Index,
		EAction: MonsStatusEndActions[st]})
	return true
}

func (m *monster) EnterConfusion(g *game, ev event) {
	if m.PutStatus(g, MonsConfused, DurationConfusionMonster) {
		m.Path = m.Path[:0]
		if g.Player.Sees(m.Pos) {
			g.Printf("%s looks confused.", m.Kind.Definite(true))
		}
	}
}

func (m *monster) EnterLignification(g *game, ev event) {
	if m.PutStatus(g, MonsLignified, DurationLignificationMonster) {
		m.Path = m.Path[:0]
		if g.Player.Sees(m.Pos) {
			g.Printf("%s is rooted to the ground.", m.Kind.Definite(true))
		}
	}
}

func (m *monster) HitSideEffects(g *game, ev event) {
	switch m.Kind {
	case MonsEarthDragon:
		if m.Status(MonsConfused) {
			m.PushPlayer(g, 3)
		} else {
			m.PushPlayer(g, 6)
		}
		g.Confusion(ev)
	case MonsBlinkingFrog:
		if g.Blink(ev) {
			g.StoryPrintf("Blinked away by %s", m.Kind)
			g.Stats.TimesBlinked++
		}
	case MonsYack:
		m.PushPlayer(g, 5)
	case MonsAcidMound:
		m.Corrode(g)
	case MonsSpider:
		g.Confusion(g.Ev)
	case MonsWingedMilfid:
		if m.Status(MonsExhausted) || g.Player.HasStatus(StatusLignification) {
			break
		}
		ompos := m.Pos
		m.MoveTo(g, g.Player.Pos)
		g.PlacePlayerAt(ompos)
		g.Print("The flying milfid makes you swap positions.")
		g.StoryPrintf("Position swap by %s", m.Kind)
		m.ExhaustTime(g, 5+RandInt(5))
		if g.Dungeon.Cell(g.Player.Pos).T == ChasmCell {
			g.PushAgainEvent(&simpleEvent{ERank: ev.Rank(), EAction: AbyssFall})
		}
	case MonsTinyHarpy:
		if m.Status(MonsSatiated) {
			return
		}
		g.Player.Bananas--
		if g.Player.Bananas < 0 {
			g.Player.Bananas = 0
		} else {
			m.PutStatus(g, MonsSatiated, DurationSatiationMonster)
			g.Print("The tiny harpy steals a banana from you.")
			g.StoryPrintf("Banana stolen by %s (bananas: %d)", m.Kind, g.Player.Bananas)
			g.Stats.StolenBananas++
			m.Target = m.NextTarget(g)
			m.MakeWander()
		}
	}
}

func (m *monster) PushPlayer(g *game, dist int) {
	if g.Player.HasStatus(StatusLignification) {
		return
	}
	dir := g.Player.Pos.Dir(m.Pos)
	pos := g.Player.Pos
	npos := pos
	path := []position{pos}
	i := 0
	for {
		i++
		npos = npos.To(dir)
		path = append(path, npos)
		if !npos.valid() || g.Dungeon.Cell(npos).BlocksRange() {
			path = path[:len(path)-1]
			break
		}
		mons := g.MonsterAt(npos)
		if mons.Exists() {
			continue
		}
		pos = npos
		if i >= dist {
			break
		}
	}
	if pos == g.Player.Pos {
		// TODO: do more interesting things, perhaps?
		return
	}
	g.Stats.TimesPushed++
	c := g.Dungeon.Cell(pos)
	var cs string
	if m.Kind == MonsEarthDragon {
		cs = " inadvertently"
		if m.Status(MonsConfused) {
			cs = " out of confusion"
		}
	}
	g.PlacePlayerAt(pos)
	g.Printf("%s pushes you%s.", m.Kind.Definite(true), cs)
	g.StoryPrintf("Pushed by %s%s", m.Kind.Definite(true), cs)
	g.ui.PushAnimation(path)
	if c.T == ChasmCell {
		g.PushAgainEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: AbyssFall})
	}
}

func (m *monster) RangedAttack(g *game, ev event) bool {
	if !m.Kind.Ranged() {
		return false
	}
	if m.Status(MonsConfused) {
		g.Printf("%s appears too confused to attack.", m.Kind.Definite(true))
		return false
	}
	if m.Pos.Distance(g.Player.Pos) <= 1 && m.Kind != MonsSatowalgaPlant {
		return false
	}
	if !m.SeesPlayer(g) {
		m.FireReady = false
		return false
	}
	if !m.FireReady {
		m.FireReady = true
		if m.Pos.Distance(g.Player.Pos) <= 3 {
			ev.Renew(g, DurationTurn)
			return true
		} else {
			return false
		}
	}
	if m.Status(MonsExhausted) {
		return false
	}
	switch m.Kind {
	//case MonsLich:
	//return m.TormentBolt(g, ev)
	case MonsHighGuard:
		return m.ThrowJavelin(g, ev)
	case MonsSatowalgaPlant:
		return m.ThrowAcid(g, ev)
	case MonsMadNixe:
		return m.NixeAttraction(g, ev)
	case MonsVampire:
		return m.VampireSpit(g, ev)
	case MonsTreeMushroom:
		return m.ThrowSpores(g, ev)
	}
	return false
}

func (m *monster) RangeBlocked(g *game) bool {
	ray := g.Ray(m.Pos)
	if len(ray) == 1 {
		return false
	}
	if len(ray) == 0 {
		// should not happen
		return true
	}
	for _, pos := range ray[1:] {
		c := g.Dungeon.Cell(pos)
		if c.BlocksRange() {
			return true
		}
		mons := g.MonsterAt(pos)
		if mons != nil {
			return true
		}
	}
	return false
}

func (g *game) BarrierCandidates(pos position, todir direction) []position {
	candidates := pos.ValidCardinalNeighbors()
	bestpos := pos.To(todir)
	if bestpos.Distance(pos) > 1 {
		j := 0
		for i := 0; i < len(candidates); i++ {
			if candidates[i].Distance(bestpos) == 1 {
				candidates[j] = candidates[i]
				j++
			}
		}
		if len(candidates) > 2 {
			candidates = candidates[0:2]
		}
		return candidates
	}
	worstpos := pos.To(pos.Dir(bestpos))
	for i := 1; i < len(candidates); i++ {
		if candidates[i] == bestpos {
			candidates[0], candidates[i] = candidates[i], candidates[0]
		}
	}
	for i := 1; i < len(candidates)-1; i++ {
		if candidates[i] == worstpos {
			candidates[len(candidates)-1], candidates[i] = candidates[i], candidates[len(candidates)-1]
		}
	}
	if len(candidates) == 4 && RandInt(2) == 0 {
		candidates[1], candidates[2] = candidates[2], candidates[1]
	}
	if len(candidates) == 4 {
		candidates = candidates[0:3]
	}
	return candidates
}

func (m *monster) CreateBarrier(g *game, ev event) bool {
	// TODO: add noise?
	dir := g.Player.Pos.Dir(m.Pos)
	candidates := g.BarrierCandidates(g.Player.Pos, dir)
	done := false
	for _, pos := range candidates {
		c := g.Dungeon.Cell(pos)
		mons := g.MonsterAt(pos)
		if mons.Exists() || c.IsWall() {
			continue
		}
		g.MagicalBarrierAt(pos, ev)
		done = true
		g.Print("The oric celmist creates a magical barrier.")
		g.StoryPrintf("Blocked by %s barrier", m.Kind)
		g.Stats.TimesBlocked++
		break
	}
	if !done {
		return false
	}
	ev.Renew(g, DurationTurn)
	m.Exhaust(g)
	return true
}

func (m *monster) Illuminate(g *game, ev event) bool {
	if g.PutStatus(StatusIlluminated, DurationIlluminated) {
		g.Print("The harmonic celmist casts magical harmonies on you.")
		g.StoryPrintf("Illuminated by %s", m.Kind)
		g.MakeNoise(HarmonicNoise, g.Player.Pos)
		ev.Renew(g, DurationTurn)
		m.Exhaust(g)
		return true
	}
	return false
}

func (m *monster) VampireSpit(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked || g.Player.HasStatus(StatusConfusion) {
		return false
	}
	g.Print("The vampire spits at you.")
	g.Print("A vampire spitted at you.")
	g.Confusion(ev)
	m.Exhaust(g)
	ev.Renew(g, DurationTurn)
	return true
}

func (m *monster) ThrowSpores(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked || g.Player.HasStatus(StatusLignification) {
		return false
	}
	g.Print("The tree mushroom releases spores.")
	g.StoryPrintf("Lignified by %s", m.Kind)
	g.EnterLignification(ev)
	m.Exhaust(g)
	ev.Renew(g, DurationTurn)
	return true
}

func (m *monster) ThrowJavelin(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	dmg := DmgNormal
	clang := RandInt(4) == 0
	noise := g.HitNoise(clang)
	var sclang string
	if clang {
		sclang = g.ClangMsg()
	}
	g.Printf("%s throws a javelin at you (%d dmg).%s", m.Kind.Definite(true), dmg, sclang)
	g.StoryPrintf("Targeted by %s javelin", m.Kind)
	g.ui.MonsterJavelinAnimation(g.Ray(m.Pos), true)
	g.MakeNoise(noise, g.Player.Pos)
	m.InflictDamage(g, dmg, dmg)
	m.ExhaustTime(g, 10+RandInt(5))
	ev.Renew(g, DurationTurn)
	return true
}

func (m *monster) Corrode(g *game) {
	count := 0
	for i, _ := range g.Player.Magaras {
		n := RandInt(2)
		g.Player.Magaras[i].Charges -= n
		if g.Player.Magaras[i].Charges < 0 {
			g.Player.Magaras[i].Charges = 0
		} else {
			count += n
		}
	}
	if count > 0 {
		g.Printf("You lose %d magara charges by corrosion.", count)
		g.StoryPrintf("Corroded by %s (lost %d magara charges)", m.Kind, count)
	}
}

func (m *monster) ThrowAcid(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	dmg := DmgNormal
	noise := g.HitNoise(false) // no clang with acid projectiles
	g.Printf("%s throws acid at you (%d dmg).", m.Kind.Definite(true), dmg)
	g.ui.MonsterProjectileAnimation(g.Ray(m.Pos), '*', ColorGreen)
	g.MakeNoise(noise, g.Player.Pos)
	m.InflictDamage(g, dmg, dmg)
	m.Corrode(g)
	m.ExhaustTime(g, 2)
	ev.Renew(g, DurationTurn)
	return true
}

func (m *monster) NixeAttraction(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	g.MakeNoise(MagicCastNoise, m.Pos)
	g.PrintfStyled("%s lures you to her.", logMonsterHit, m.Kind.Definite(true))
	g.StoryPrintf("Lured by %s", m.Kind)
	ray := g.Ray(m.Pos)
	g.ui.MonsterProjectileAnimation(ray, '*', ColorCyan)
	if len(ray) > 1 {
		// should always be the case
		g.ui.TeleportAnimation(g.Player.Pos, ray[1], true)
		g.PlacePlayerAt(ray[1])
	}
	m.Exhaust(g)
	ev.Renew(g, DurationTurn)
	return true
}

func (m *monster) SmitingAttack(g *game, ev event) bool {
	if !m.Kind.Smiting() {
		return false
	}
	if m.Status(MonsConfused) {
		g.Printf("%s appears too confused to attack.", m.Kind.Definite(true))
		return false
	}
	if !m.SeesPlayer(g) {
		m.FireReady = false
		return false
	}
	if !m.FireReady {
		m.FireReady = true
		if m.Pos.Distance(g.Player.Pos) <= 3 {
			ev.Renew(g, DurationTurn)
			return true
		} else {
			return false
		}
	}
	if m.Status(MonsExhausted) {
		return false
	}
	switch m.Kind {
	case MonsMirrorSpecter:
		return m.AbsorbMana(g, ev)
	case MonsOricCelmist:
		return m.CreateBarrier(g, ev)
		//case MonsMindCelmist:
		//return m.MindAttack(g, ev)
	case MonsHarmonicCelmist:
		return m.Illuminate(g, ev)
	}
	return false
}

func (m *monster) AbsorbMana(g *game, ev event) bool {
	if g.Player.MP == 0 {
		return false
	}
	g.Player.MP -= 1
	g.Printf("%s absorbs your mana.", m.Kind.Definite(true))
	g.StoryPrintf("Mana absorbed by %s (MP: %d)", m.Kind, g.Player.MP)
	m.ExhaustTime(g, 1+RandInt(2))
	ev.Renew(g, DurationTurn)
	return true
}

func (m *monster) Blink(g *game) {
	npos := g.BlinkPos(true)
	if !npos.valid() || npos == g.Player.Pos || npos == m.Pos {
		return
	}
	opos := m.Pos
	g.Printf("The %s blinks away.", m.Kind)
	g.ui.TeleportAnimation(opos, npos, true)
	m.MoveTo(g, npos)
}

func (m *monster) MakeHunt(g *game) (noticed bool) {
	if m.State != Hunting {
		m.State = Hunting
		g.Stats.NSpotted++
		g.Stats.DSpotted[g.Depth]++
		if !m.Alerted {
			g.Stats.NUSpotted++
			g.Stats.DUSpotted[g.Depth]++
			if g.Stats.DUSpotted[g.Depth] >= 90*len(g.Monsters)/100 {
				AchUnstealthy.Get(g)
			}
		}
		m.Alerted = true
		noticed = true
	}
	m.Search = g.Player.Pos
	m.Target = g.Player.Pos
	return noticed
}

func (m *monster) MakeWatchIfHurt(g *game) {
	// TODO: not used now.
	if m.Exists() && m.State != Hunting {
		m.MakeHunt(g)
		if m.State == Resting {
			g.Printf("%s awakens.", m.Kind.Definite(true))
		}
		if m.Kind == MonsDog {
			g.Printf("%s barks.", m.Kind.Definite(true))
			g.MakeNoise(BarkNoise, m.Pos)
		}
	}
}

func (m *monster) MakeAware(g *game) {
	if m.Peaceful(g) || m.Status(MonsSatiated) {
		if m.State == Resting && m.Pos.Distance(g.Player.Pos) == 1 {
			g.Printf("%s awakens.", m.Kind.Definite(true))
			m.MakeWander()
		}
		return
	}
	if !m.SeesPlayer(g) {
		return
	}
	if m.State == Resting {
		g.Printf("%s awakens.", m.Kind.Definite(true))
	} else if m.State == Wandering || m.State == Watching {
		g.Printf("%s notices you.", m.Kind.Definite(true))
	}
	noticed := m.MakeHunt(g)
	if noticed && m.Kind == MonsDog {
		g.Printf("%s barks.", m.Kind.Definite(true))
		g.StoryPrintf("Barked at by %s", m.Kind)
		g.MakeNoise(BarkNoise, m.Pos)
	}
}

func (m *monster) GatherBand(g *game) {
	if !MonsBands[g.Bands[m.Band].Kind].Band {
		return
	}
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{m.Pos}, 4)
	for _, mons := range g.Monsters {
		if mons.Band == m.Band {
			if mons.State == Hunting && m.State != Hunting {
				continue
			}
			n, ok := nm.at(mons.Pos)
			if !ok || n.Cost > 4 || mons.State == Resting && mons.Status(MonsExhausted) {
				continue
			}
			mons.Target = m.Target
			if mons.State == Resting {
				mons.MakeWander()
			}
		}
	}
}

func (g *game) MonsterAt(pos position) *monster {
	if !pos.valid() {
		return nil
	}
	i := g.MonstersPosCache[pos.idx()]
	if i <= 0 {
		return nil
	}
	return g.Monsters[i-1]
}

func (g *game) MonsterInLOS() *monster {
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.Pos) {
			return mons
		}
	}
	return nil
}
