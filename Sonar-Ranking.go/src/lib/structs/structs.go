package structs

type Project struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Projects struct {
	Components []Project `json:"components"`
}

type Rule struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Lang     string `json:"lang"`
	LangName string `json:"langName"`
}

type FacetValue struct {
	Val   string `json:"val"`
	Count int    `json:"count"`
}

type Facet struct {
	Property string       `json:"property"`
	Values   []FacetValue `json:"values"`
}

type Issues struct {
	Rules  []Rule  `json:"rules"`
	Facets []Facet `json:"facets"`
}

type OwaspCategory struct {
	Category        string `json:"category"`
	Vulnerabilities int    `json:"vulnerabilities"`
}

type SecurityReport struct {
	Categories []OwaspCategory `json:"categories"`
}

type Top10Rule struct {
	Key   string
	Name  string
	Lang  string
	Count int
}
