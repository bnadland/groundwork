package {{Name}}_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bnadland/{{Name}}"
	"github.com/gocolly/colly"
)

func NewTestServer(tb testing.TB) (*httptest.Server, *url.URL, *colly.Collector) {
	if testing.Verbose() == false {
		{{Name}}.DisableLog()
	}
	s := httptest.NewServer({{Name}}.NewHandler(NewTestConfig()))
	u, err := url.Parse(s.URL)
	if err != nil {
		tb.Fatal(err)
	}
	c := colly.NewCollector(
		colly.AllowedDomains(u.Host),
	)
	return s, u, c
}
