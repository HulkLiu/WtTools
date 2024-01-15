package model

type Profile struct {
	Title         string
	LanguageType  string
	Source        string
	Description   string
	Catalog       string
	ResourcesLink string
	Bk            string
	CreateAt      string
	LastTime      string
	Status        string
}

type SearchResult struct {
	Hits     int64
	Start    int
	Query    string
	PrevFrom int
	NextFrom int
	Items    []interface{}
}
type Item struct {
	Url     string
	Type    string
	Id      string
	Payload interface{}
}
