package main

import (
	"log"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

type room struct {
	p       gruid.Point
	w       int
	h       int
	entries []rentry
	places  []place
	vault   *rl.Vault
	special bool
	tunnels int
}

type places struct {
	Shaedra  gruid.Point
	Monolith gruid.Point
	Marevor  gruid.Point
	Artifact gruid.Point
}

type rentry struct {
	p       gruid.Point
	used    bool
	virtual bool
}

type placeKind int

const (
	PlaceDoor placeKind = iota
	PlacePatrol
	PlaceStatic
	PlaceSpecialStatic
	PlaceItem
	PlaceStory
	PlacePatrolSpecial
)

type place struct {
	p    gruid.Point
	kind placeKind
	used bool
}

func roomDistance(r1, r2 *room) int {
	// TODO: use the center?
	return abs(r1.p.X+r1.w/2-r2.p.X-r2.w/2) + abs(r1.p.Y+r1.h/2-r2.p.Y-r2.h/2)
}

func (r *room) HasSpace(dg *dgen) bool {
	if DungeonWidth-r.p.X < r.w || DungeonHeight-r.p.Y < r.h {
		return false
	}
	for i := r.p.X - 1; i <= r.p.X+r.w; i++ {
		for j := r.p.Y - 1; j <= r.p.Y+r.h; j++ {
			p := gruid.Point{i, j}
			if valid(p) && dg.room[p] {
				return false
			}
		}
	}
	return true
}

func (r *room) Dig(dg *dgen) {
	runedraw := func(p gruid.Point, c rune) {
		q := gruid.Point{X: r.p.X + p.X, Y: r.p.Y + p.Y}
		if valid(q) && c != '?' {
			dg.room[q] = true
		}
		switch c {
		case '.', '>', '!', 'P', '_', '|', 'G', '-':
			if valid(q) {
				dg.d.SetCell(q, GroundCell)
			}
		case 'B':
			// obstacle
			t := WallCell
			switch RandInt(9) {
			case 0, 6:
				t = TreeCell
			case 1:
				if RandInt(2) == 0 {
					t = QueenRockCell
				} else {
					t = LightCell
				}
			case 2:
				if RandInt(2) == 0 {
					t = ChasmCell
				} else {
					t = TableCell
				}
			case 3:
				t = TableCell
			case 4, 5:
				t = GroundCell
			}
			if valid(q) {
				dg.d.SetCell(q, t)
			}
		case '#', '+':
			if valid(q) {
				dg.d.SetCell(q, WallCell)
			}
		case 'T':
			if valid(q) {
				dg.d.SetCell(q, TreeCell)
			}
		case 'π':
			if valid(q) {
				dg.d.SetCell(q, TableCell)
			}
		case 'l':
			if valid(q) {
				dg.d.SetCell(q, LightCell)
			}
		case 'W':
			if valid(q) {
				dg.d.SetCell(q, WindowCell)
			}
		case '"',
			'?',
			',',
			'~',
			'c',
			'q',
			'S',
			'M',
			'Δ',
			'A':
		default:
			log.Fatalf("Invalid terrain: %c for room w:%d h:%d pos:%+v\n%s", c, r.w, r.h, r.p, r.vault.Content())
		}
		switch c {
		case '>':
			r.places = append(r.places, place{p: q, kind: PlaceSpecialStatic})
		case '!':
			r.places = append(r.places, place{p: q, kind: PlaceItem})
		case 'P':
			r.places = append(r.places, place{p: q, kind: PlacePatrol})
		case 'G':
			r.places = append(r.places, place{p: q, kind: PlacePatrolSpecial})
		case '_':
			r.places = append(r.places, place{p: q, kind: PlaceStatic})
		case '|':
			r.places = append(r.places, place{p: q, kind: PlaceDoor})
		case '+', '-':
			if q.X == 0 || q.X == DungeonWidth-1 || q.Y == 0 || q.Y == DungeonHeight-1 {
				break
			}
			e := rentry{}
			e.p = q
			if c == '-' {
				e.virtual = true
			}
			r.entries = append(r.entries, e)
		case '"':
			if valid(q) {
				dg.d.SetCell(q, FoliageCell)
			}
		case ',':
			if valid(q) {
				dg.d.SetCell(q, CavernCell)
			}
		case '~':
			if valid(q) {
				dg.d.SetCell(q, WaterCell)
			}
		case 'c':
			if valid(q) {
				dg.d.SetCell(q, ChasmCell)
			}
		case 'q':
			if valid(q) {
				dg.d.SetCell(q, QueenRockCell)
			}
		case 'S':
			r.places = append(r.places, place{p: q, kind: PlaceStory})
			dg.spl.Shaedra = q
			dg.d.SetCell(q, StoryCell)
		case 'M':
			r.places = append(r.places, place{p: q, kind: PlaceStory})
			dg.spl.Marevor = q
			dg.d.SetCell(q, StoryCell)
		case 'Δ':
			r.places = append(r.places, place{p: q, kind: PlaceStory})
			dg.spl.Monolith = q
			dg.d.SetCell(q, StoryCell)
		case 'A':
			r.places = append(r.places, place{p: q, kind: PlaceStory})
			dg.spl.Artifact = q
			dg.d.SetCell(q, StoryCell)
		}
	}
	r.vault.Iter(runedraw)
}

// UnusedEntry returns an unused entry, if possible, or a random entry
// otherwise.
func (r *room) UnusedEntry() int {
	ens := []int{}
	for i, e := range r.entries {
		if !e.used {
			ens = append(ens, i)
		}
	}
	if len(ens) == 0 {
		return RandInt(len(r.entries))
	}
	return ens[RandInt(len(ens))]
}

func (r *room) RandomPlace(kind placeKind) gruid.Point {
	var p []int
	for i, pl := range r.places {
		if pl.kind == kind && !pl.used {
			p = append(p, i)
		}
	}
	if len(p) == 0 {
		return InvalidPos
	}
	j := p[RandInt(len(p))]
	r.places[j].used = true
	return r.places[j].p
}

var PlaceSpecialOrStatic = []placeKind{PlaceSpecialStatic, PlaceStatic}

func (r *room) RandomPlaces(kinds []placeKind) gruid.Point {
	p := InvalidPos
	for _, kind := range kinds {
		p = r.RandomPlace(kind)
		if p != InvalidPos {
			break
		}
	}
	return p
}

const (
	RoomAlmostSquare = `
?###+##??
#_....!#?
+..PBP.!#
#!....._#
?###+###?`
	RoomSquareBis = `
?##?#+##?
#_!#!.._#
#..|.P.P+
#_!#!.._#
?##?#+##?`
	RoomRoundSimple = `
??#+#??
?#!.!#?
#_.P._#
+.PBP.+
#_.P._#
?#!.!#?
??#+#??`
	RoomLittle = `
?#+#?
#_._#
+.P.+
#_.!#
?#+#?`
	RoomLittleDiamond = `
??#+#??
##!._##
+..P..+
##_._##
??#+#??`
	RoomLittleColumnDiamond = `
?##+#??
#_..!##
+.PBP.+
##!.._#
??#+##?`
	RoomLittleTreeDiamond = `
??#+##?
##!.B_#
+.P.P.#
#_B.!##
?##+#??`
	RoomRound = `
???##+##???
??#".P."#??
##".!B_."##
+.P.B#B.P.+
##"._B_."##
??#".P."#??
???##+##???`
	RoomRoundTree = `
???##+##???
??#".P."#??
##"._._."##
+.P..B..P.+
##"._.!."##
??#".P."#??
???##+##???`
)

var roomNormalTemplates = []string{RoomAlmostSquare, RoomSquareBis, RoomRoundSimple, RoomLittle, RoomLittleDiamond, RoomLittleColumnDiamond, RoomRound, RoomLittleTreeDiamond, RoomRoundTree}

const (
	RoomBigColumns = `
?####?#++#?####?
#!.._##..##!..>#
##.P........P.##
+...._####._...+
##.P........P.##
#!..>##..##>..!#
?####?#++#?####?`
	RoomBigGarden = `
?####?#++#?####?
#""""##..##...!#
#"T""".!P...~~.#
#""""">P_>.~~~.+
#"T""".!P...~~"#
#""""##..##..""#
?####?#++#?####?`
	RoomBigRooms = `
?####?#++#?###??
#_..!##..##!..#?
#"""..#..#..π._#
#"""P.|.P|.P..>#
#"""..#..#..π.!#
#>..!##..##...#?
?####?#++#?###??`
	RoomColumns = `
###+##+###
+..P..P..+
#.#>B#!#.#
#.#!##>B.#
+..P..P..+
###+##+###`
	RoomTriangles = `
?####?######?
#...>#_.....#
+.P.#...#.P.+
#..#_...!#..#
#!....P...#>#
?#..!#|#..>#?
??###.-.###??`
	RoomHome1 = `
?########+#?
#>.P..|..P!#
#...B.#....#
#!...!#_...#
###|######|#
#...P...|..#
#>.....!#P.#
?########+#?`
	RoomHome2 = `
##########??
+...#....>#?
#.P.|..P...#
##|###..._.#
#...>#!...!#
#....##|####
#!P..|..P..+
?########+##`
	RoomHome3 = `
?#############?
#>P.|.........#
#..!##|##!.P..+
####!...!#_...#
+..|.P>._###|##
####!...!#!...#
#!.>##|##..P..+
#.P.....|...B.#
?#......#....#?
??###########??`
	RoomHome4 = `
?############?
#..._.#.....!#
#.P...|..B...#
##|####....P.+
#...._#>.##|##
+..P..####>..#
#!....|....P!#
?###+########?`
	RoomHome5 = `
????##########?
??###!...>#..!#
?###>.....!#..#
#..#...P.B.#.P+
+P.#_......#..+
#..####|####..#
#!.....P....._#
?######+######?`
	RoomHome6 = `
?###+#???
#>.#q##??
#.P|P_##?
#.!#.#>.#
#_##.##|#
##!.....#
+.PBPB..#
#......_#
?###+###?`
	RoomCaban = `
???????-??????
?????""""?????
???""""""""???
??"""###."""??
?"""#>!|PP"""?
-""""###.""""-
??""""""""""??
????"""""?????
???????-??????`
	RoomDolmen = `
???????-?????
????.....????
????.!#!.????
???..#>#..???
??....P....??
?...B_._B...?
-.P...P...P.-
??...#>#...??
????.!#!.????
????.....????
???????-?????`
	RoomSmallTemple = `
?????????????
?????...?????
????..#..????
???..#>#..???
??..#>._#..??
??..#...#..??
?...#!P!#...?
-.P..#|#..P.-
?..#.....#..?
??..#.P_#..??
???.......???
??????-??????`
	RoomTemple = `
????.....????
??...###...??
?..##!>!##..?
..#_....._#..
.#....B....#.
.#.##.P.##.#.
.#.>#...#>.#.
-#..P...P..#.
..###|||###..
?.....P.....-
??."T...T".??
???"".-.""???`
	RoomSchool = `
????#######???
???#l..!..l#??
??#!..Pπ..!##?
?#.........|.#
.W.π.π..π..#>#
.#.........l##
.W.π.π..π.##>#
.#........#..#
-W.π.π..π.#P.#
.#l...P...|._#
..####||#####-
?.....--P.....`
	RoomShop = `
??????###???
?????#>P>#??
??#####|##??
?#l...P..l#?
.W.π.π!.π.#?
.#.π.π..π.#?
-W.π.π..π.#?
.#...π.P..l#
.#l!..qq_!#?
..####||##.-
?.....--P..?`
	RoomTavern = `
??##########??
?#l.π...|..P#?
.#..!...###>#?
.W.π.π...!###?
.#....P.π.#>.#
.W.π.π..π.#P.#
.#l...qqπ.|._#
-.####||###|#-
?.....--P.....`
	RoomDoctor = `
????#####??
???#l...!#?
??##.Pππ.#?
?#.|.....#?
.W>#_...l#?
.####|###??
.W...P#>.#?
.#.π..##|#.
.W.!.....W-
.#l..Pq._#.
-.####|##..
?.....-.P.?`
	RoomRuins = `
????????-????
????....P????
???..###.????
???.##>#..???
-.....P....??
?..##_"!##..?
??.#""""#.P.-
??.""B###..??
????"#>#.????
????..P..????
??????-??????`
	RoomPillars = `
???????-?????
???.......???
?...#!#P#...?
?.#!......#.?
-.P.B>B>B.P.-
?.#......!#.?
?...#P#!#...?
???.......???
?????-???????`
	RoomRoundColumns = `
???##+##???
??#!...!#??
##_.BP#._##
+...P>P...+
##_.#P#._##
??#!...!#??
???##+##???`
	RoomTriangle = `
?????#?????
????#>#????
???#!.!#???
??#_..._#??
##!.#.B.!##
+..P...P..+
##_..P.._##
??###+###??`
	RoomSpiraling = `
?#####+#
#!.P...+
#.>#####
#!.P..!#
#####>.#
+...P.!#
#+#####?`
	RoomSpiralingCircle = `
??##+#?##??
##.P.!#.>#?
+.P.##....#
#._#...#_.#
#....##!P.+
?#.>#.P..##
??###+###??`
	RoomCircleDouble = `
???####+#???
??#""..P#???
?#""..#|#???
#"""!#!P!#??
#.._#.....#?
#..#..>#>..#
+.P|P.###.P#
#..#.._#>..#
#"._#.....##
#"""!#!P!#!#
?#""..#|#..#
??#"...P..#?
???####+##??`
	RoomGardenHome = `
???#######???
??#.#>#>#.#??
?#!...P...!#?
#...#._.#...#
#_........._#
######|######
+.....P..#_>#
+.......#!..#
#######|#...#
#""""..P|.P.#
?#"T"...#...+
??#""".!#!..#
???#######+#?`
	RoomAltar = `
#+#??#######??#+#
+P_##>..!..>##_P+
#...#..#>#..#...#
?##.#!..P..!#.##?
??#..P.....P..#??
???#####+#####???`
	RoomRoundGarden = `
???##+##???
??#_.P.>#??
##!.""".!##
+.P""B""P.+
##!.""".!##
??#>.P.>#??
???##+##???`
	RoomLongHall = `
#################
+.P...........P.+
#..!#!>.B.>!#!..#
+.P...........P.+
#################`
	RoomGardenHall = `
?#############?
#"""".>!>.""""#
+..P...P...P..+
#""""._!>.""""#
?#############?`
	RoomRoundHall = `
????###?###+#?
?###!_##>#.P.#
#....###|#...#
+.P_..P....P.+
#....#!#.....#
?###.>#?#.π.##
????##???###??`
	RoomToilets = `
??######
?#!.P..+
-Wπ..###
.#!..|>#
.WπP.##?
.#!..|>#
-|.P.###
?####???`
	RoomPicnic = `
????""""????
??""""T"""??
?"T.!P..T""?
?"".πlπ!.."?
-..P.>.P...-
?""!π.π..."?
?"T..P!.T""?
??""".."""??
??????-?????`
	RoomSnake = `
?????#####???###
?#?##!..._#?##.+
#B#........#.|.#
#...######.P.###
#..##....###..#?
-P.#...#.!#>#.#?
-..#..##>P..#.#?
#π#"...!####..#?
?##"""T.P....#??
???#"""....!#???
????########????`
)

var roomBigTemplates = []string{RoomBigColumns, RoomBigGarden, RoomColumns,
	RoomRoundColumns, RoomRoundGarden, RoomLongHall, RoomGardenHall,
	RoomTriangles, RoomHome1, RoomHome2, RoomHome3, RoomHome4, RoomHome5, RoomHome6,
	RoomTriangle, RoomSpiraling, RoomSpiralingCircle, RoomAltar,
	RoomCircleDouble, RoomGardenHome, RoomBigRooms, RoomCaban, RoomDolmen,
	RoomSmallTemple, RoomTemple, RoomSchool, RoomTavern, RoomShop,
	RoomDoctor, RoomRuins, RoomPillars, RoomRoundHall, RoomToilets, RoomPicnic, RoomSnake}

const (
	CellShaedra = `
?#?#?#?#?
#########
#SMΔ#_!_#
##|###|##
+.G.P.G.+
##|###|##
#_!_#_!_#
#########
?#?#?#?#?
`
	CellShaedra2 = `
ccccccc-?
c#####c..
c#SMΔ#ccP
c##|####.
cWG..G.|-
c##|##.##
c#_!_#!π#
c#####|#?
cccccc-P?
`
	CellShaedra3 = `
ccccccc-?
c#####c..
c#_!_#ccP
c##|####.
cWG..G.|-
c##|##.##
c#SMΔ#!π#
c#####|#?
cccccc-P?
`
	CellShaedra4 = `
?????c?????
????ccc????
??ccccccc??
?cc##W##cc?
cc##!G!##cc
c#_#...#S#c
c#!|...|M#c
c#_#.G.#Δ#c
c###_._###c
-.###|###.-
?.P..-..P.?
`
)

// TODO: add indestructible walls?

var roomCellTemplates = []string{CellShaedra, CellShaedra2, CellShaedra3, CellShaedra4}

const (
	RoomArtifact = `
????#????
???#c#???
??#cAc#??
?##.MΔ##?
####|####
#!_#.#_!#
+P.G.G.P+
#>..P..>#
?###+###?`
	RoomArtifact2 = `
??cccc??
?cc#####
cc#cc!P+
c##M#G.#
c#ΔA|..#
c##.#G.#
c###π.!#
-|q..P>#
P######?`
	RoomArtifact3 = `
?#####+#?
#B..Gcq.#
#.Gc.cP.#
#.ccMccq#
#.c.AΔcP#
#.>ccc!.#
#...P...#
#!.#q_.B#
?###+###?`
)

var roomArtifactTemplates = []string{RoomArtifact, RoomArtifact2, RoomArtifact3}

const (
	RoomSpecialNixes = `
?#########???
##.#~>~##>##?
#>.G~.~G...!#
##.......B..#
#!.#qqq#..._#
##.G...G.#..#
#...._....P.+
+P.......#!##
#....P...!#??
?####+####???`
	RoomSpecialVampires = `
?#####+#####?
#>.G.#.#...>#
?#.......G.#?
??#!_.P._.#!#
???##....#>G#
??#~_#P._!#|#
?#...|......#
#>.G.#.#..._#
?#####+#####?`
	RoomSpecialVampires2 = `
?~~~~~~~~~~-.
~~##W##W##~~.
~#>._#.π.##~P
~W.G.|...|.~-
~|q.!#.G.|.G-
~####>...|.~-
~#>G|....##~P
~~###.π.!#~~.
?~~~##W###~~.
???~~~~~~~~-.`
	RoomSpecialFrogs = `
?~~-~~?~??
~~.P~~~~~?
~..._~B.~?
?~.G.>~.~~
~~..!.G.B~
-P..B>~.~~
~~.G>~~.~?
?~!..BG.~~
?~~~P..~~?
???~-~~???`
	RoomSpecialMilfids = `
?????????-???
???......P.??
??..!?G._?..?
?.?.?#>c.?.??
-P.G.>#>G...?
?.?._c>c.??.?
??.......?..-
????!.G...P??
????????-????`
	RoomSpecialMilfids2 = `
?????????-??
-P.ccc?c.Gc?
?c.G..cc..cc
??ccc._G..>c
?ccc>....!cc
??c!....>ccc
??cc.G.cccc?
???c.....???
?????P-.????`
	RoomSpecialTreeMushrooms = `
????"--.???????
??"""..G."""???
??""?....!"""??
??....T>T..G..?
?.G.T..!..T..P-
-..!..T>T._"".-
-P???..G..""""?
.-??????.-?????`
	RoomSpecialTreeMushrooms2 = `
???""-???????
??""?.."""???
??""..G..""??
?""?....!..??
?".T.T.T.T..?
?.G..>!>.G.P-
-.!π.T.T....-
-.?......_""?
?P??".TG.""??
.-?????.-????`
	RoomSpecialHarpies = `
?-????##??????
?P???#..#?????
?.???#G.>####?
?.??#...##.>#?
?.?#.._....G.#
-G........!..#
??.#.G._._#>#?
??P?#.>###?#??
??-??##???????`
	RoomSpecialHarpies2 = `
?-#????##..P
-P.####...#-
?#.G.#..G.!#
?##..._....#
?#.G.....!##
#>.#_.G.##??
?##?###.>#??
???????##???`
	RoomSpecialCelmists = `
?############+##?
#>#_......>#.P._#
#...G!#!G.##....#
#....###...|...P+
#_...###...|...P+
##..G!#!G.##....#
#>........>#.P._#
?############+##?`
	RoomSpecialCelmists2 = `
?##?##########
#_.#....G....+
#..#>#W#|#W#!#
#..###>.._.###
+P.|.G.....G>#
#..###._.!.###
#..#!#W#|#W#>#
#_#.....G....#
?##+#########?`
	RoomSpecialCelmists3 = `
?###???.-.??????
#...#..P...?????
+.P.|....P.....-
##|##W#|||#W####
?#....G...G....#
??#!._#._.#_.!#?
???#..G...G..#??
????#!.....!#???
?????#>#>#>#????
??????#?#?#?????`
	RoomSpecialMirrorSpecters = `
########-########
-P.....W.W..q..P-
##W##_.#q#._##W##
-P......G......P-
##W##..##W##.#W##
#>.!W..W>!>W.W!>#
#.G.#..#.G.#.#G.#
#...............#
#################`
	RoomSpecialMirrorSpecters2 = `
######--#####
-P...W.P...P-
#.#####B.##.#
#.W>.G...!W.#
#.###..#.##G#
#.....#..q..#
#W#.B#>.!##W#
#>#!.....#.>#
#..#..B.#..#?
?#G..#!.G.#??
??###?####???`
)

type specialRoom int

const (
	noSpecialRoom specialRoom = iota
	roomMilfids
	roomFrogs
	roomNixes
	roomVampires
	roomCelmists
	roomHarpies
	roomTreeMushrooms
	roomMirrorSpecters
	roomShaedra
	roomArtifact
)

func (sr specialRoom) Templates() (tpl []string) {
	switch sr {
	case roomMilfids:
		tpl = append(tpl, RoomSpecialMilfids, RoomSpecialMilfids2)
	case roomFrogs:
		tpl = append(tpl, RoomSpecialFrogs)
	case roomVampires:
		tpl = append(tpl, RoomSpecialVampires, RoomSpecialVampires2)
	case roomCelmists:
		tpl = append(tpl, RoomSpecialCelmists, RoomSpecialCelmists2, RoomSpecialCelmists3)
	case roomNixes:
		tpl = append(tpl, RoomSpecialNixes)
	case roomHarpies:
		tpl = append(tpl, RoomSpecialHarpies, RoomSpecialHarpies2)
	case roomTreeMushrooms:
		tpl = append(tpl, RoomSpecialTreeMushrooms, RoomSpecialTreeMushrooms2)
	case roomMirrorSpecters:
		tpl = append(tpl, RoomSpecialMirrorSpecters, RoomSpecialMirrorSpecters2)
	case roomShaedra:
		tpl = roomCellTemplates
	case roomArtifact:
		tpl = roomArtifactTemplates
	}
	return tpl
}
