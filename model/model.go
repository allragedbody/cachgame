package model

type Organization struct {
	Name   string
	Income float64
	Retio  float64
}

type Country struct {
	Name   string
	Income float64
	Retio  float64
}

type Gamer struct {
	Name         string
	Gamertype    string
	Total        float64
	CashTotal    float64
	SharesTotals float64
	Shares       map[string]Share
}

type Share struct {
	ShareID     string
	Number      float64
	SharesTotal float64
	AvgPrice    float64
}

type ShareInfo struct {
	ShareID   string
	ShareName string
	Price     float64
}

type ShareGame struct {
	Organization Organization
	Country      Country
	Gamers       map[string]Gamer
}
