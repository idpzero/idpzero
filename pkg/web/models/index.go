package models

type UrlInfo struct {
	Description string
	Url         string
}

type IndexModel struct {
	Urls []UrlInfo
}
