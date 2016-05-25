package opensearch

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/util"
)

type OpenSearchXmlImage struct {
	XMLName xml.Name `xml:"Image"`

	Height int    `xml:"height,attr"`
	Width  int    `xml:"width,attr"`
	Type   string `xml:"type,attr"`
	Value  string `xml:",chardata"`
}

type OpenSearchXmlUrl struct {
	XMLName xml.Name `xml:"Url"`

	Type     string `xml:"type,attr"`
	Template string `xml:"template,attr"`
}

type OpenSearchXml struct {
	XMLName       xml.Name `xml:"OpenSearchDescription"`
	XMLNameSpace  string   `xml:"xmlns,attr"`
	InputEncoding string   `xml:"InputEncoding"`
	ShortName     string   `xml:"ShortName"`
	Description   string   `xml:"Description"`

	Image OpenSearchXmlImage

	Url OpenSearchXmlUrl
}

func OpenSearch(c *gin.Context) {
	c.XML(http.StatusOK, &OpenSearchXml{
		XMLNameSpace:  "http://a9.com/-/spec/opensearch/1.1/",
		InputEncoding: "utf-8",
		ShortName:     conf.Conf.Site.Name,
		Description:   conf.Conf.Site.Description,

		Image: OpenSearchXmlImage{
			Height: 16,
			Width:  16,
			Type:   "image/x-icon",
			Value:  util.UrlJoin("favicon.ico"),
		},
		Url: OpenSearchXmlUrl{
			Type:     "text/html",
			Template: util.GetSearchUrl("") + "/{searchTerms}",
		},
	})
}
