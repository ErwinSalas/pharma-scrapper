package types

import "encoding/xml"

type URLSet struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"` // Importante: Namespace
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc string `xml:"loc"`
}

type Product struct {
	Name     string
	Price    string
	Lab      string
	Active   string
	Img      string
	Pharmacy string
}

type Pharmacy struct {
	ID               string
	Domain           string
	ProductSelectors Product
}
