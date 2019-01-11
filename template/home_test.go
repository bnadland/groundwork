package {{Name}}_test

import (
	"testing"

	"github.com/gocolly/colly"
)

func TestHomeHandler(t *testing.T) {
	s, u, c := NewTestServer(t)
	defer s.Close()

	c.OnResponse(func(resp *colly.Response) {
		if resp.StatusCode != 200 {
			t.Errorf("unexpected status code: %v", resp.StatusCode)
		}
	})

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		if e.Text != "hello, world" {
			t.Errorf("unexpected headline: %s", e.Text)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		t.Error(err)
	})

	u.Path = "/"
	c.Visit(u.String())
	c.Wait()
}

func BenchmarkHomeHandler(b *testing.B) {
	b.ReportAllocs()
	s, u, c := NewTestServer(b)
	defer s.Close()

	u.Path = "/"
	for n := 0; n < b.N; n++ {
		c.Visit(u.String())
		c.Wait()
	}
}
