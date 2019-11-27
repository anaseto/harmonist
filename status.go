package main

type status int

const (
	StatusSlow status = iota
	StatusExhausted
	StatusSwift
	StatusLignification
	StatusConfusion
	StatusNausea
	StatusFlames // fake status
	StatusHidden
	StatusUnhidden
	StatusLight
	StatusDig
	StatusLevitation
	StatusShadows
	StatusIlluminated
	StatusTransparent
)

func (st status) Info() bool {
	switch st {
	case StatusFlames, StatusHidden, StatusUnhidden, StatusLight:
		return true
	}
	return false
}

func (st status) Good() bool {
	switch st {
	case StatusSwift, StatusDig, StatusHidden, StatusLevitation, StatusShadows, StatusTransparent:
		return true
	default:
		return false
	}
}

func (st status) Bad() bool {
	switch st {
	case StatusSlow, StatusConfusion, StatusNausea, StatusUnhidden, StatusIlluminated:
		return true
	default:
		return false
	}
}

func (st status) String() string {
	switch st {
	case StatusSlow:
		return "Slow"
	case StatusExhausted:
		return "Exhausted"
	case StatusSwift:
		return "Swift"
	case StatusLignification:
		return "Lignified"
	case StatusConfusion:
		return "Confused"
	case StatusNausea:
		return "Nausea"
	case StatusFlames:
		return "Flames"
	case StatusHidden:
		return "Hidden"
	case StatusUnhidden:
		return "Unhidden"
	case StatusDig:
		return "Dig"
	case StatusLight:
		return "Light"
	case StatusLevitation:
		return "Levitating"
	case StatusShadows:
		return "Shadows"
	case StatusIlluminated:
		return "Illuminated"
	case StatusTransparent:
		return "Transparent"
	default:
		// should not happen
		return "unknown"
	}
}

func (st status) Short() string {
	switch st {
	case StatusSlow:
		return "Sl"
	case StatusExhausted:
		return "Ex"
	case StatusSwift:
		return "Sw"
	case StatusLignification:
		return "Lg"
	case StatusConfusion:
		return "Co"
	case StatusNausea:
		return "Na"
	case StatusFlames:
		return "Fl"
	case StatusHidden:
		return "H+"
	case StatusUnhidden:
		return "H-"
	case StatusDig:
		return "Di"
	case StatusLight:
		return "Li"
	case StatusLevitation:
		return "Le"
	case StatusShadows:
		return "Sh"
	case StatusTransparent:
		return "Tr"
	default:
		// should not happen
		return "?"
	}
}
