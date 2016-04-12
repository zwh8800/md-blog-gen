package sitemap

import (
	"encoding/xml"
	"io"
	"time"
)

type Item struct {
	XMLName xml.Name  `xml:"url"`
	Link    string    `xml:"loc"`
	Updated time.Time `xml:"lastmod"`

	// Updated  time.Time `xml:"-"`
	// UpdatedS string    `xml:"lastmod"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`

	Items []*Item
}

func (p *Sitemap) Add(item *Item) {
	// item.UpdatedS = item.Updated.Format(time.RFC3339)
	p.Items = append(p.Items, item)
}

func New() *Sitemap {
	return &Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Items: make([]*Item, 0),
	}
}

func Xml(s *Sitemap, w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	en := xml.NewEncoder(w)
	//en.Indent("", "  ")
	return en.Encode(s)
}
