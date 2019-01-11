package {{Name}}

import (
	"html/template"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/masterminds/sprig"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"github.com/unrolled/render"
)

func NewRenderer(config *Config) *render.Render {
	templates := packr.New("templates", "./templates")
	return render.New(render.Options{
		IsDevelopment: config.IsDevelopment,
		Extensions:    []string{".html"},
		Layout:        "layout/main",
		Directory:     "templates",
		Asset: func(name string) ([]byte, error) {
			return templates.Find(strings.TrimPrefix(name, "templates/"))
		},
		AssetNames: func() []string {
			result := []string{}
			for _, i := range templates.List() {
				result = append(result, strings.Join([]string{"templates", i}, "/"))
			}
			return result
		},
		Funcs: []template.FuncMap{
			sprig.FuncMap(),
			template.FuncMap{
				"markdown": func(content string) template.HTML {
					return template.HTML(
						bluemonday.UGCPolicy().SanitizeBytes(
							blackfriday.Run(
								[]byte(content),
								blackfriday.WithRenderer(
									blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
										HeadingLevelOffset: 1,
									}),
								),
								blackfriday.WithExtensions(blackfriday.AutoHeadingIDs),
							),
						),
					)
				},
			},
		},
	})
}
