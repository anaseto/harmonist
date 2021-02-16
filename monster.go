package main

import (
	"fmt"
	"log"

	"github.com/anaseto/gruid"
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

func (mk monsterKind) GoodFlair() bool {
	switch mk {
	case MonsDog, MonsHazeCat, MonsVampire, MonsExplosiveNadre:
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

func (mk monsterKind) ShallowSleep() bool {
	switch mk {
	case MonsCrazyImp, MonsHazeCat:
		return true
	default:
		return false
	}
}

func (mk monsterKind) ResistsLignification() bool {
	switch mk {
	case MonsSatowalgaPlant, MonsTreeMushroom:
		return true
	default:
		return false
	}
}

func (mk monsterKind) ReflectsTeleport() bool {
	switch mk {
	case MonsBlinkingFrog:
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
	MonsWorm:            "Farmer worms are ugly creeping creatures. They furrow as they move, helping new foliage to grow.",
	//MonsBrizzia:         "Brizzias are big slow moving biped creatures. They are quite hardy, and when hurt they can cause nausea, impeding the use of potions.",
	MonsAcidMound: "Acid mounds are acidic creatures. They can corrode your magaras, reducing their number of charges.",
	MonsDog:       "Dogs are carnivore quadrupeds. They can bark, and smell you from up to 5 tiles away when hunting or watching for you.",
	MonsYack:      "Yacks are quite large herbivorous quadrupeds. They tend to eat grass peacefully, but upon seing you they may attack, pushing you up to 5 cells away.",
	//MonsGiantBee:        "Giant bees are fragile but extremely fast moving creatures. Their bite can sometimes enrage you.",
	MonsHighGuard: "High guards watch over a particular location. They can throw javelins.",
	//MonsHydra:           "Hydras are enormous creatures with four heads that can hit you each at once.",
	//MonsSkeletonWarrior: "Skeleton warriors are good fighters, clad in chain mail.",
	MonsSpider:       "Spiders are small creatures, with panoramic vision and whose bite can confuse you.",
	MonsWingedMilfid: "Winged milfids are  humanoids that can fly over you and make you swap positions. They tend to be very agressive creatures.",
	MonsBlinkingFrog: "Blinking frogs are big frog-like creatures, whose bite can make you blink away. The science behind their attack is not clear, but many think it relies on some kind of oric deviation magic. They can jump to attack from below.",
	//MonsLich:           "Liches are non-living mages wearing a leather armour. They can throw a bolt of torment at you, halving your HP.",
	MonsEarthDragon:    "Earth dragons are big creatures from a dragon species that wander in the Underground. They are peaceful creatures, but they may hurt you inadvertently, pushing you up to 6 tiles away (3 if confused). They naturally emit powerful oric energies, allowing them to eat rocks and dig tunnels. Their oric energies can confuse you if you're close enough, for example if they hurt you or you jump over them.",
	MonsMirrorSpecter:  "Mirror specters are very insubstantial creatures, which can absorb your mana.",
	MonsExplosiveNadre: "Nadres are dragon-like biped creatures that are famous for exploding upon dying. Explosive nadres are a tiny nadre race that explodes upon attacking. The explosion confuses any adjacent creatures and occasionally destroys walls.",
	MonsSatowalgaPlant: "Satowalga Plants are immobile bushes that throw viscous acidic projectiles at you, destroying some of your magara charges. They attack at half normal speed.",
	MonsMadNixe:        "Nixes are magical humanoids. Usually, they specialize in illusion harmonic magic, but the so called mad nixes are a perverted variant who learned the oric arts to create a spell that can attract their foes to them, so that they can kill them without pursuing them.",
	//MonsMindCelmist:     "Mind celmists are mages that use magical smitting mind attacks that bypass armour. They can occasionally confuse or slow you. They try to avoid melee.",
	MonsVampire:      "Vampires are humanoids that drink blood to survive. Their nauseous spitting can cause confusion, impeding the use of magaras for a few turns.",
	MonsTreeMushroom: "Tree mushrooms are big clunky creatures. They can throw lignifying spores at you, leaving you unable to move for a few turns, though the spores will also provide some protection against harm.",
	//MonsMarevorHelith: "Marevor Helith is an ancient undead nakrus very fond of teleporting people away. He is a well-known expert in the field of magaras - items that many people simply call magical objects. His current research focus is monolith creation. Marevor, a repentant necromancer, is now searching for his old disciple Jaixel in the Underground to help him overcome the past.",
	MonsButterfly: "Underground's butterflies, called kerejats, wander peacefully around, illuminating their surroundings.",
	MonsCrazyImp:  "Crazy Imp is a crazy creature that likes to sing with its small guitar. It seems to be fond of monkeys and quite capable at finding them by flair. While singing it may attract unwanted attention.",
	MonsHazeCat:   "Haze cats are a special variety of cats found in the Underground. They have very good night vision and are always alert.",
}

type bandInfo struct {
	Path []gruid.Point
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
	PairSpider
	PairHazeCat
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
	PairSpider:                 {Band: true, Distribution: map[monsterKind]int{MonsSpider: 2}},
	PairHazeCat:                {Band: true, Distribution: map[monsterKind]int{MonsHazeCat: 2}},
	PairSatowalga:              {Band: true, Distribution: map[monsterKind]int{MonsSatowalgaPlant: 2}},
	PairWorm:                   {Band: true, Distribution: map[monsterKind]int{MonsWorm: 2}},
	PairVampire:                {Band: true, Distribution: map[monsterKind]int{MonsVampire: 2}},
	PairOricCelmist:            {Band: true, Distribution: map[monsterKind]int{MonsOricCelmist: 2}},
	PairHarmonicCelmist:        {Band: true, Distribution: map[monsterKind]int{MonsHarmonicCelmist: 2}},
	PairNixe:                   {Band: true, Distribution: map[monsterKind]int{MonsMadNixe: 2}},
	PairExplosiveNadre:         {Band: true, Distribution: map[monsterKind]int{MonsExplosiveNadre: 2}},
	PairWingedMilfid:           {Band: true, Distribution: map[monsterKind]int{MonsWingedMilfid: 2}},
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
	P             gruid.Point
	LastKnownPos  gruid.Point
	Target        gruid.Point
	Path          []gruid.Point // cache
	FireReady     bool
	Seen          bool
	LOS           map[gruid.Point]bool
	LastSeenState monsterState
	Swapped       bool
	Watching      int
	Left          bool
	Search        gruid.Point
	Alerted       bool
	Waiting       int
}

func (m *monster) Init() {
	m.Attack = m.Kind.BaseAttack()
	m.P = InvalidPos
	m.LOS = make(map[gruid.Point]bool)
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

func (m *monster) CanPass(g *game, p gruid.Point) bool {
	if !valid(p) {
		return false
	}
	c := g.Dungeon.Cell(p)
	return c.IsPassable() ||
		c.IsDoorPassable() && m.Kind.CanOpenDoors() ||
		c.IsLevitatePassable() && m.Kind.CanFly() ||
		c.IsSwimPassable() && (m.Kind.CanSwim() || m.Kind.CanFly()) ||
		terrain(c) == HoledWallCell && m.Kind.Size() == MonsSmall
}

func (m *monster) CanPassDestruct(g *game, p gruid.Point) bool {
	if !valid(p) {
		return false
	}
	c := g.Dungeon.Cell(p)
	destruct := false
	if m.Kind == MonsEarthDragon {
		destruct = true
	}
	return m.CanPass(g, p) || c.IsDestructible() && destruct
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
	p := m.P
	i := 0
	count := 0
	for {
		count++
		if count > 1000 {
			panic("TeleportOther")
		}
		p = g.FreePassableCell()
		if Distance(p, m.P) < 15 && i < 1000 {
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
		m.Target = m.P
	}
	if g.Player.Sees(m.P) {
		g.Printf("%s teleports away.", m.Kind.Definite(true))
	}
	opos := m.P
	m.MoveTo(g, p)
	if g.Player.Sees(opos) {
		g.md.TeleportAnimation(opos, p, false)
	}
}

func (m *monster) MoveTo(g *game, p gruid.Point) {
	if g.Player.Sees(p) {
		m.UpdateKnowledge(g, p)
	} else if g.Player.Sees(m.P) {
		if Distance(m.P, p) == 1 {
			// You know the direction, so you know where the
			// monster should be.
			m.UpdateKnowledge(g, p)
		} else {
			delete(g.LastMonsterKnownAt, m.P)
		}
		m.LastKnownPos = InvalidPos
	}
	if !g.Player.Sees(m.P) && g.Player.Sees(p) {
		if !m.Seen {
			m.Seen = true
			g.Printf("%s (%v) comes into view.", m.Kind.Indefinite(true), m.State)
		}
		g.StopAuto()
	}
	recomputeLOS := g.Player.Sees(m.P) && terrain(g.Dungeon.Cell(m.P)) == DoorCell ||
		g.Player.Sees(p) && terrain(g.Dungeon.Cell(p)) == DoorCell
	m.PlaceAt(g, p)
	if recomputeLOS {
		g.ComputeLOS()
	}
	c := g.Dungeon.Cell(p)
	if terrain(c) == ChasmCell && !m.Kind.CanFly() || terrain(c) == WaterCell && !m.Kind.CanSwim() && !m.Kind.CanFly() {
		m.Dead = true
		if g.Player.Sees(m.P) {
			g.HandleKill(m)
			switch terrain(c) {
			case ChasmCell:
				g.Printf("%s falls into the abyss.", m.Kind.Definite(true))
			case WaterCell:
				g.Printf("%s drowns.", m.Kind.Definite(true))
			}
		}
	}
}

func (m *monster) PlaceAtStart(g *game, p gruid.Point) {
	i := idx(p)
	if g.MonstersPosCache[i] > 0 {
		log.Printf("used monster starting position: %v", p)
		// should not happen
		return
	}
	m.P = p
	g.MonstersPosCache[i] = m.Index + 1
	npos := m.RandomFreeNeighbor(g)
	if npos != m.P {
		m.Dir = Dir(m.P, npos)
	} else {
		m.Dir = E
	}
}

func (m *monster) PlaceAt(g *game, p gruid.Point) {
	if !valid(m.P) {
		log.Printf("monster place at: bad position %v", m.P)
		// should not happen
		return
	}
	if p == InvalidPos {
		log.Printf("monster new place at invalid position")
		// should not happen
		return
	}
	if p == m.P {
		// should not happen
		log.Printf("monster place at: same position %v", m.P)
		return
	}
	if p == g.Player.P {
		log.Printf("monster place at: player position %v", p)
		// should not happen
		return
	}
	i := idx(p)
	j := idx(m.P)
	m.Dir = Dir(m.P, p)
	m.CorrectDir()
	mons := g.MonsterAt(p)
	g.MonstersPosCache[i], g.MonstersPosCache[j] = g.MonstersPosCache[j], g.MonstersPosCache[i]
	if mons.Exists() {
		m.P, mons.P = mons.P, m.P
		mons.Swapped = true
	} else {
		m.P = p
	}
	m.Waiting = 0
}

func (g *game) MonsterAt(p gruid.Point) *monster {
	pi := idx(p)
	if pi < 0 || pi >= len(g.MonstersPosCache) {
		log.Printf("monster at: bad index %v for pos %v", pi, p)
		// should not happen
		return nil
	}
	i := g.MonstersPosCache[pi]
	if i <= 0 {
		return nil
	}
	m := g.Monsters[i-1]
	//if m.P != p {
	//log.Printf("monster position mismatch: %v vs %v", m.P, p)
	//}
	return m
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

func (m *monster) EndTurn(g *game) {
	g.PushEventD(&monsterTurnEvent{Index: m.Index}, 1)
}

func (m *monster) AttackAction(g *game) bool {
	m.Dir = Dir(m.P, g.Player.P)
	m.CorrectDir()
	switch m.Kind {
	case MonsExplosiveNadre:
		g.StoryPrint("Nadre explosion")
		m.Explode(g)
		return false
	default:
		if g.Player.HasStatus(StatusDispersal) {
			m.Blink(g)
		} else {
			m.HitPlayer(g)
		}
	}
	return true
}

func (m *monster) Explode(g *game) {
	m.Dead = true
	neighbors := ValidCardinalNeighbors(m.P)
	g.Printf("%s %s explodes with a loud boom.", g.ExplosionSound(), m.Kind.Definite(true))
	g.md.ExplosionAnimation(FireExplosion, m.P)
	g.MakeNoise(ExplosionNoise, m.P)
	for _, p := range append(neighbors, m.P) {
		c := g.Dungeon.Cell(p)
		if c.Flammable() {
			g.Burn(p)
		}
		mons := g.MonsterAt(p)
		if mons.Exists() && !mons.Status(MonsConfused) {
			mons.EnterConfusion(g)
			if mons.State != Hunting && mons.State != Watching {
				mons.StartWatching()
			}
		} else if g.Player.P == p {
			m.InflictDamage(g, 1, 1)
		} else if c.IsDestructible() && RandInt(3) > 0 {
			if c.IsDiggable() {
				g.Dungeon.SetCell(p, RubbleCell)
			} else {
				g.Dungeon.SetCell(p, GroundCell)
			}
			if terrain(c) == BarrelCell {
				delete(g.Objects.Barrels, p)
			}
			g.Stats.Digs++
			g.UpdateKnowledge(p, terrain(c))
			if g.Player.Sees(p) {
				g.md.WallExplosionAnimation(p)
			}
			g.MakeNoise(WallNoise, p)
			g.Fog(p, 1)
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

func (m *monster) RandomFreeNeighbor(g *game) gruid.Point {
	p := m.P
	neighbors := [4]gruid.Point{p.Add(gruid.Point{1, 0}), p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, -1}), p.Add(gruid.Point{0, 1})}
	fnb := []gruid.Point{}
	for _, nbpos := range neighbors {
		if !valid(nbpos) {
			continue
		}
		c := g.Dungeon.Cell(nbpos)
		if c.IsPassable() {
			fnb = append(fnb, nbpos)
		}
	}
	if len(fnb) == 0 {
		return m.P
	}
	samedir := fnb[RandInt(len(fnb))]
	for _, p := range fnb {
		// invariant: pos != m.Pos
		if m.Dir.InViewCone(m.P, To(Dir(m.P, p), p)) {
			samedir = p
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

var SearchAroundCache []gruid.Point

func (m *monster) SearchAround(g *game, p gruid.Point, radius int) gruid.Point {
	dij := &monPath{state: g, monster: m}
	nodes := g.PR.DijkstraMap(dij, []gruid.Point{p}, radius)
	SearchAroundCache = SearchAroundCache[:0]
	for _, n := range nodes {
		SearchAroundCache = append(SearchAroundCache, n.P)
	}
	if len(SearchAroundCache) > 0 {
		p := SearchAroundCache[RandInt(len(SearchAroundCache))]
		return p
	}
	return InvalidPos
}

func (m *monster) NextTarget(g *game) (p gruid.Point) {
	band := g.Bands[m.Band]
	p = m.P
	switch band.Beh {
	case BehWander:
		if Distance(m.P, band.Path[0]) < 8+RandInt(8) {
			p = m.SearchAround(g, m.P, 4)
			if p != InvalidPos {
				break
			}
		}
		if m.Search != InvalidPos && RandInt(2) == 0 {
			p = m.SearchAround(g, m.Search, 7)
			if p != InvalidPos {
				break
			}
		}
		p = m.SearchAround(g, band.Path[0], 7)
		if p != InvalidPos {
			break
		}
		p = band.Path[0]
	case BehExplore:
		if m.Kind.CanOpenDoors() {
			if m.Search != InvalidPos && RandInt(4) == 0 {
				p = m.SearchAround(g, m.Search, 7)
			} else {
				p = m.SearchAround(g, p, 5)
			}
			if p != InvalidPos {
				break
			}
		}
		p = band.Path[RandInt(len(band.Path))]
	case BehGuard:
		if m.Search != InvalidPos && Distance(m.Search, m.P) < 5 && RandInt(2) == 0 {
			p = m.SearchAround(g, m.Search, 3)
			if p != InvalidPos {
				break
			}
		}
		p = band.Path[0]
	case BehPatrol:
		if m.Search != InvalidPos && RandInt(4) > 0 {
			p = m.SearchAround(g, m.Search, 7)
			if p != m.P && p != InvalidPos {
				break
			}
		}
		if band.Path[0] == m.Target {
			p = band.Path[1]
		} else if band.Path[1] == m.Target {
			p = band.Path[0]
		} else if Distance(band.Path[0], m.P) < Distance(band.Path[1], m.P) {
			p = band.Path[0]
			if RandInt(4) == 0 {
				p = band.Path[1]
			}
		} else {
			p = band.Path[1]
			if RandInt(4) == 0 {
				p = band.Path[0]
			}
		}
	case BehCrazyImp:
		path := m.APath(g, m.P, g.Player.P)
		if len(path) == 0 {
			p = m.SearchAround(g, m.P, 3)
		} else {
			p = g.Player.P
		}
	}
	return p
}

func (m *monster) HandleMonsSpecifics(g *game) (done bool) {
	switch m.Kind {
	case MonsSatowalgaPlant:
		switch m.State {
		case Hunting:
			if m.Target != InvalidPos && m.Target != m.P {
				m.Dir = Dir(m.P, m.Target)
			}
			if !m.SeesPlayer(g) {
				m.StartWatching()
			}
		default:
			if RandInt(4) > 0 {
				m.Alternate()
			}
		}
		// oklob plants are static ranged-only
		return true
	case MonsGuard, MonsHighGuard:
		if m.State != Wandering && m.State != Watching {
			break
		}
		for p, on := range g.Objects.Lights {
			if !on && p == m.P {
				g.Dungeon.SetCell(m.P, LightCell)
				g.Objects.Lights[m.P] = true
				if g.Player.Sees(m.P) {
					g.Printf("%s makes a new fire.", m.Kind.Definite(true))
				} else {
					g.UpdateKnowledge(m.P, ExtinguishedLightCell)
				}
				return true
			} else if !on && m.SeesLight(g, p) {
				m.Target = p
			}
		}
	case MonsCrazyImp:
		if g.Player.Sees(m.P) && RandInt(2) == 0 && !m.Status(MonsConfused) && !m.Status(MonsExhausted) {
			g.PrintStyled("Crazy Imp: “♫ larilon, larila ♫ ♪”", logSpecial)
			g.MakeNoise(SingingNoise, m.P)
			//g.ui.MusicAnimation(m.Pos)
			m.Exhaust(g)
		}
	}
	return false
}

const DogFlairDist = 5
const DisguiseFlairDist = 3

func (m *monster) HandleWatching(g *game) {
	turns := 4
	if m.Kind == MonsHazeCat {
		turns = 3
	}
	if m.Watching+RandInt(2) < turns {
		m.Alternate()
		m.Watching++
		if m.Kind == MonsDog {
			dij := &monPath{state: g, monster: m}
			g.PR.DijkstraMap(dij, []gruid.Point{m.P}, DogFlairDist)
			if c := g.PR.DijkstraMapAt(g.Player.P); c <= DogFlairDist {
				m.Target = g.Player.P
				m.MakeWander()
			}
		}
	} else {
		// pick a random cell: more escape strategies for the player
		m.Target = m.NextTarget(g)
		switch g.Bands[m.Band].Beh {
		case BehGuard:
			m.Alternate()
			if m.P != m.Target {
				m.MakeWander()
				m.GatherBand(g)
			}
		default:
			m.MakeWander()
			m.GatherBand(g)
		}
	}
}

func (m *monster) ComputePath(g *game) {

	if !(len(m.Path) > 0 && m.Path[len(m.Path)-1] == m.Target && m.Path[0] == m.P) {
		m.Path = m.APath(g, m.P, m.Target)
		if len(m.Path) == 0 && !m.Status(MonsConfused) {
			// if target is not accessible, try free neighbor cells
			for _, npos := range g.Dungeon.FreeNeighbors(m.Target) {
				m.Path = m.APath(g, m.P, npos)
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
	if m.State != Hunting && g.Player.HasStatus(StatusDisguised) && (!m.Kind.GoodFlair() || Distance(m.P, g.Player.P) > DisguiseFlairDist) {
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
	if len(m.Path) == 0 && m.Search != InvalidPos && Distance(m.Search, m.Target) < 10 && m.P != m.Target {
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
}

func (m *monster) MakeWanderAt(target gruid.Point) {
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

func (m *monster) HandleMove(g *game) bool {
	target := m.Path[1]
	mons := g.MonsterAt(target)
	monstarget := InvalidPos
	if mons.Exists() && len(mons.Path) >= 2 {
		monstarget = mons.Path[1]
	}
	c := g.Dungeon.Cell(target)
	switch {
	case m.Peaceful(g) && target == g.Player.P:
		switch m.Kind {
		case MonsEarthDragon:
			return m.AttackAction(g)
		default:
			m.Path = m.APath(g, m.P, m.Target)
		}
	case !mons.Exists():
		if m.Kind == MonsEarthDragon && c.IsDestructible() {
			g.Dungeon.SetCell(target, RubbleCell)
			if terrain(c) == BarrelCell {
				delete(g.Objects.Barrels, target)
			}
			g.Stats.Digs++
			g.UpdateKnowledge(target, terrain(c))
			g.MakeNoise(WallNoise, m.P)
			g.Fog(m.P, 1)
			if Distance(g.Player.P, target) < 12 {
				// XXX use dijkstra distance ?
				if c.IsWall() {
					g.Printf("%s You hear an earth-splitting noise.", g.CrackSound())
				} else if terrain(c) == BarrelCell || terrain(c) == DoorCell || terrain(c) == TableCell {
					g.Printf("%s You hear an wood-splitting noise.", g.CrackSound())
				}
				g.StopAuto()
			}
			m.MoveTo(g, target)
			m.Path = m.Path[1:]
		} else if !m.CanPass(g, target) {
			m.Path = m.APath(g, m.P, m.Target)
		} else {
			m.InvertFoliage(g)
			m.MoveTo(g, target)
			if (m.Kind.Ranged() || m.Kind.Smiting()) && !m.FireReady && m.SeesPlayer(g) {
				m.FireReady = true
			}
			m.Path = m.Path[1:]
		}
	case (mons.P == target && m.P == monstarget || m.Waiting > 5+RandInt(2)) && !mons.Status(MonsLignified):
		target := mons.P
		m.MoveTo(g, target)
		m.Path = m.Path[1:]
		if len(mons.Path) > 0 {
			mons.Path = mons.Path[1:]
		}
	case m.State == Hunting && mons.State != Hunting:
		if m.Waiting > 2+RandInt(3) {
			if mons.Peaceful(g) {
				mons.MakeWander()
			} else {
				mons.MakeWanderAt(m.Target)
				mons.GatherBand(g)
			}
		} else {
			m.Path = m.APath(g, m.P, m.Target)
		}
		m.Waiting++
	case !mons.SeesPlayer(g) && mons.State != Hunting:
		if m.Waiting > 1+RandInt(2) && mons.Kind != MonsSatowalgaPlant {
			mons.MakeWanderAt(mons.RandomFreeNeighbor(g))
		} else {
			m.Path = m.APath(g, m.P, m.Target)
		}
		m.Waiting++
	default:
		m.Path = m.APath(g, m.P, m.Target)
		m.Waiting++
	}
	return true
}

func (m *monster) HandleTurn(g *game) {
	m.handleTurn(g)
	if m.Exists() {
		m.EndTurn(g)
	}
}

func (m *monster) handleTurn(g *game) {
	if m.Status(MonsParalysed) {
		return
	}
	if m.Swapped {
		m.Swapped = false
		return
	}
	ppos := g.Player.P
	mpos := m.P
	switch m.Kind {
	case MonsGuard, MonsHighGuard:
		// they have to put lights on, could be optimized (TODO)
		m.ComputeLOS(g)
	}
	m.MakeAware(g)
	if m.State == Resting {
		if RandInt(3000) == 0 || m.Kind.ShallowSleep() && RandInt(10) == 0 {
			m.NaturalAwake(g)
		}
		return
	}
	if m.State == Hunting && m.RangedAttack(g) {
		return
	}
	if m.State == Hunting && m.SmitingAttack(g) {
		return
	}
	if m.HandleMonsSpecifics(g) {
		return
	}
	if Distance(mpos, ppos) == 1 && terrain(g.Dungeon.Cell(ppos)) != BarrelCell && !m.Peaceful(g) {
		if m.Status(MonsConfused) {
			g.Printf("%s appears too confused to attack.", m.Kind.Definite(true))
			return
		}
		if terrain(g.Dungeon.Cell(ppos)) == TreeCell && !m.Kind.CanAttackOnTree() {
			g.Printf("%s watches you from below.", m.Kind.Definite(true))
			return
		}
		m.AttackAction(g)
		return
	}
	if m.Status(MonsLignified) {
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
	c := g.Dungeon.Cell(m.P)
	if terrain(c) == CavernCell {
		g.Dungeon.SetCell(m.P, FoliageCell)
		invert = true
	} else if terrain(c) == FoliageCell {
		g.Dungeon.SetCell(m.P, CavernCell)
		invert = true
	}
	if invert {
		if g.Player.Sees(m.P) {
			g.ComputeLOS()
		} else {
			g.UpdateKnowledge(m.P, terrain(c))
		}
	}
}

func (m *monster) Exhaust(g *game) {
	m.ExhaustTime(g, DurationExhaustionMonster+RandInt(DurationExhaustionMonster/2))
}

func (m *monster) ExhaustTime(g *game, t int) {
	m.PutStatus(g, MonsExhausted, t)
}

func (m *monster) HitPlayer(g *game) {
	if g.Player.HP <= 0 || Distance(g.Player.P, m.P) > 1 {
		return
	}
	dmg := m.Attack
	clang := RandInt(4) == 0
	noise := g.HitNoise(clang)
	g.MakeNoise(noise, g.Player.P)
	var sclang string
	if clang {
		sclang = g.ClangMsg()
	}
	g.PrintfStyled("%s hits you (%d dmg).%s", logMonsterHit, m.Kind.Definite(true), dmg, sclang)
	m.InflictDamage(g, dmg, m.Attack)
	if g.Player.HP <= 0 {
		return
	}
	m.HitSideEffects(g)
	const HeavyWoundHP = 2
	if g.Player.HP >= HeavyWoundHP {
		return
	}
	switch g.Player.Inventory.Neck {
	case AmuletConfusion:
		m.EnterConfusion(g)
		// TODO: maybe affect all monsters in sight?
		g.Printf("Your amulet releases confusing harmonies.")
	case AmuletFog:
		g.Print("Your amulet feels warm.")
		g.SwiftFog()
	case AmuletObstruction:
		opos := m.P
		p := m.LeaveRoomForPlayer(g)
		if p == InvalidPos {
			m.Blink(g)
		} else {
			m.MoveTo(g, p)
		}
		if opos != m.P {
			g.MagicalBarrierAt(opos)
			g.Print("The amulet releases an oric wind.")
			m.Exhaust(g)
		}
	case AmuletTeleport:
		g.Print("Your amulet shines.")
		m.TeleportAway(g)
		if m.Kind.ReflectsTeleport() {
			g.Printf("The %s reflected back some energies.", m.Kind)
			g.Teleportation()
		}
	case AmuletLignification:
		g.Print("Your amulet glows.")
		if !m.Kind.ResistsLignification() {
			m.EnterLignification(g)
		}
	}
}

func (m *monster) PutStatus(g *game, st monsterStatus, duration int) bool {
	if m.Status(st) {
		return false
	}
	m.Statuses[st] += duration
	g.PushEventD(&monsterStatusEvent{
		Index:  m.Index,
		Status: st}, DurationStatusStep)

	return true
}

func (m *monster) EnterConfusion(g *game) bool {
	if m.PutStatus(g, MonsConfused, DurationConfusionMonster) {
		m.Path = m.Path[:0]
		if g.Player.Sees(m.P) {
			g.Printf("%s looks confused.", m.Kind.Definite(true))
		}
		return true
	}
	return false
}

func (m *monster) EnterLignification(g *game) {
	if m.PutStatus(g, MonsLignified, DurationLignificationMonster) {
		m.Path = m.Path[:0]
		if g.Player.Sees(m.P) {
			g.Printf("%s is rooted to the ground.", m.Kind.Definite(true))
		}
	}
}

func (m *monster) HitSideEffects(g *game) {
	switch m.Kind {
	case MonsEarthDragon:
		if m.Status(MonsConfused) {
			m.PushPlayer(g, 3)
		} else {
			m.PushPlayer(g, 6)
		}
		g.Confusion()
	case MonsBlinkingFrog:
		if g.Blink() {
			g.StoryPrintf("Blinked away by %s", m.Kind)
			g.Stats.TimesBlinked++
		}
	case MonsYack:
		m.PushPlayer(g, 5)
	case MonsAcidMound:
		m.Corrode(g)
	case MonsSpider:
		g.Confusion()
	case MonsWingedMilfid:
		if m.Status(MonsExhausted) || g.Player.HasStatus(StatusLignification) {
			break
		}
		c := g.Dungeon.Cell(g.Player.P)
		if !(c.IsPassable() || c.IsSwimPassable() || c.IsDoorPassable() || c.IsLevitatePassable()) {
			break
		}
		g.PlacePlayerAt(m.P)
		g.Print("The flying milfid makes you swap positions.")
		g.StoryPrintf("Position swap by %s", m.Kind)
		m.ExhaustTime(g, 5+RandInt(5))
		if terrain(g.Dungeon.Cell(g.Player.P)) == ChasmCell {
			g.PushEventFirst(&playerEvent{Action: AbyssFall}, g.Turn)
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
	dir := Dir(m.P, g.Player.P)
	p := g.Player.P
	q := p
	path := []gruid.Point{p}
	i := 0
	for {
		i++
		q = To(dir, q)
		path = append(path, q)
		if !valid(q) || g.Dungeon.Cell(q).BlocksRange() {
			path = path[:len(path)-1]
			break
		}
		mons := g.MonsterAt(q)
		if mons.Exists() {
			continue
		}
		p = q
		if i >= dist {
			break
		}
	}
	if p == g.Player.P {
		// TODO: do more interesting things, perhaps?
		return
	}
	g.Stats.TimesPushed++
	c := g.Dungeon.Cell(p)
	var cs string
	if m.Kind == MonsEarthDragon {
		cs = " inadvertently"
		if m.Status(MonsConfused) {
			cs = " out of confusion"
		}
	}
	g.md.PushAnimation(path)
	g.PlacePlayerAt(p)
	g.Printf("%s pushes you%s.", m.Kind.Definite(true), cs)
	g.StoryPrintf("Pushed by %s%s", m.Kind.Definite(true), cs)
	if terrain(c) == ChasmCell {
		g.PushEventFirst(&playerEvent{Action: AbyssFall}, g.Turn)
	}
}

func (m *monster) RangedAttack(g *game) bool {
	if !m.Kind.Ranged() {
		return false
	}
	if m.Status(MonsConfused) {
		g.Printf("%s appears too confused to attack.", m.Kind.Definite(true))
		return false
	}
	if Distance(m.P, g.Player.P) <= 1 && m.Kind != MonsSatowalgaPlant {
		return false
	}
	if !m.SeesPlayer(g) {
		m.FireReady = false
		return false
	}
	if !m.FireReady {
		m.FireReady = true
		return Distance(m.P, g.Player.P) <= 3
	}
	if m.Status(MonsExhausted) {
		return false
	}
	switch m.Kind {
	//case MonsLich:
	//return m.TormentBolt(g, ev)
	case MonsHighGuard:
		return m.ThrowJavelin(g)
	case MonsSatowalgaPlant:
		return m.ThrowAcid(g)
	case MonsMadNixe:
		return m.NixeAttraction(g)
	case MonsVampire:
		return m.VampireSpit(g)
	case MonsTreeMushroom:
		return m.ThrowSpores(g)
	}
	return false
}

func (m *monster) RangeBlocked(g *game) bool {
	ray := g.Ray(m.P)
	if len(ray) == 1 {
		return false
	}
	if len(ray) == 0 {
		// should not happen
		return true
	}
	for _, p := range ray[1:] {
		c := g.Dungeon.Cell(p)
		if c.BlocksRange() {
			return true
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			return true
		}
	}
	return false
}

func (g *game) BarrierCandidates(p gruid.Point, todir direction) []gruid.Point {
	candidates := ValidCardinalNeighbors(p)
	bestpos := To(todir, p)
	if Distance(bestpos, p) > 1 {
		j := 0
		for i := 0; i < len(candidates); i++ {
			if Distance(candidates[i], bestpos) == 1 {
				candidates[j] = candidates[i]
				j++
			}
		}
		if len(candidates) > 2 {
			candidates = candidates[0:2]
		}
		return candidates
	}
	worstpos := To(Dir(bestpos, p), p)
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

func (m *monster) CreateBarrier(g *game) bool {
	// TODO: add noise?
	dir := Dir(m.P, g.Player.P)
	candidates := g.BarrierCandidates(g.Player.P, dir)
	done := false
	for _, p := range candidates {
		c := g.Dungeon.Cell(p)
		mons := g.MonsterAt(p)
		if mons.Exists() || c.IsWall() {
			continue
		}
		g.MagicalBarrierAt(p)
		done = true
		g.Print("The oric celmist creates a magical barrier.")
		g.StoryPrintf("Blocked by %s barrier", m.Kind)
		g.Stats.TimesBlocked++
		break
	}
	if !done {
		return false
	}
	m.Exhaust(g)
	return true
}

func (m *monster) Illuminate(g *game) bool {
	if g.PutStatus(StatusIlluminated, DurationIlluminated) {
		g.Print("The harmonic celmist casts magical harmonies on you.")
		g.StoryPrintf("Illuminated by %s", m.Kind)
		g.MakeNoise(HarmonicNoise, g.Player.P)
		m.Exhaust(g)
		return true
	}
	return false
}

func (m *monster) VampireSpit(g *game) bool {
	blocked := m.RangeBlocked(g)
	if blocked || g.Player.HasStatus(StatusConfusion) {
		return false
	}
	g.Print("The vampire spits at you.")
	g.Print("A vampire spitted at you.")
	g.Confusion()
	m.Exhaust(g)
	return true
}

func (m *monster) ThrowSpores(g *game) bool {
	blocked := m.RangeBlocked(g)
	if blocked || g.Player.HasStatus(StatusLignification) {
		return false
	}
	g.Print("The tree mushroom releases spores.")
	g.StoryPrintf("Lignified by %s", m.Kind)
	g.EnterLignification()
	m.Exhaust(g)
	return true
}

func (m *monster) ThrowJavelin(g *game) bool {
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
	g.md.MonsterJavelinAnimation(g.Ray(m.P), true)
	g.MakeNoise(noise, g.Player.P)
	m.InflictDamage(g, dmg, dmg)
	m.ExhaustTime(g, 10+RandInt(5))
	return true
}

func (m *monster) Corrode(g *game) {
	count := 0
	for i := range g.Player.Magaras {
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

func (m *monster) ThrowAcid(g *game) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	dmg := DmgNormal
	noise := g.HitNoise(false) // no clang with acid projectiles
	g.Printf("%s throws acid at you (%d dmg).", m.Kind.Definite(true), dmg)
	g.md.MonsterProjectileAnimation(g.Ray(m.P), '*', ColorGreen)
	g.MakeNoise(noise, g.Player.P)
	m.InflictDamage(g, dmg, dmg)
	m.Corrode(g)
	m.ExhaustTime(g, 2)
	return true
}

func (m *monster) NixeAttraction(g *game) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	g.MakeNoise(MagicCastNoise, m.P)
	g.PrintfStyled("%s lures you to her.", logMonsterHit, m.Kind.Definite(true))
	g.StoryPrintf("Lured by %s", m.Kind)
	ray := g.Ray(m.P)
	g.md.MonsterProjectileAnimation(ray, '*', ColorCyan)
	if len(ray) > 1 {
		// should always be the case
		g.md.TeleportAnimation(g.Player.P, ray[1], true)
		g.PlacePlayerAt(ray[1])
	}
	m.Exhaust(g)
	return true
}

func (m *monster) SmitingAttack(g *game) bool {
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
		return Distance(m.P, g.Player.P) <= 3
	}
	if m.Status(MonsExhausted) {
		return false
	}
	switch m.Kind {
	case MonsMirrorSpecter:
		return m.AbsorbMana(g)
	case MonsOricCelmist:
		return m.CreateBarrier(g)
		//case MonsMindCelmist:
		//return m.MindAttack(g, ev)
	case MonsHarmonicCelmist:
		return m.Illuminate(g)
	}
	return false
}

func (m *monster) AbsorbMana(g *game) bool {
	if g.Player.MP == 0 {
		return false
	}
	g.Player.MP -= 1
	g.Printf("%s absorbs your mana.", m.Kind.Definite(true))
	g.StoryPrintf("Mana absorbed by %s (MP: %d)", m.Kind, g.Player.MP)
	m.ExhaustTime(g, 1+RandInt(2))
	return true
}

func (m *monster) Blink(g *game) {
	npos := g.BlinkPos(true)
	if !valid(npos) || npos == g.Player.P || npos == m.P {
		return
	}
	opos := m.P
	g.Printf("The %s blinks away.", m.Kind)
	g.md.TeleportAnimation(opos, npos, true)
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
	m.Search = g.Player.P
	m.Target = g.Player.P
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
			g.MakeNoise(BarkNoise, m.P)
		}
	}
}

func (m *monster) MakeAware(g *game) {
	if m.Peaceful(g) || m.Status(MonsSatiated) {
		if m.State == Resting && Distance(m.P, g.Player.P) == 1 {
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
		g.MakeNoise(BarkNoise, m.P)
	}
}

func (m *monster) GatherBand(g *game) {
	if !MonsBands[g.Bands[m.Band].Kind].Band {
		return
	}
	dij := &noisePath{state: g}
	g.PR.BreadthFirstMap(dij, []gruid.Point{m.P}, 4)
	for _, mons := range g.Monsters {
		if mons.Band == m.Band {
			if mons.State == Hunting && m.State != Hunting {
				continue
			}
			c := g.PR.BreadthFirstMapAt(mons.P)
			if c > 4 || mons.State == Resting && mons.Status(MonsExhausted) {
				continue
			}
			mons.Target = m.Target
			if mons.State == Resting {
				mons.MakeWander()
			}
		}
	}
}

func (g *game) MonsterInLOS() *monster {
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.P) {
			return mons
		}
	}
	return nil
}
