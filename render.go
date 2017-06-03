package main

import (
	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/unrolled/render"
)

func newRender() *render.Render {
	renderOptions := render.Options{
		IsDevelopment:             config.IsDevelopment,
		DisableHTTPErrorRendering: !config.IsDevelopment,
		Extensions:                []string{".html"},
		Layout:                    "layout",
		IndentJSON:                true,
		Funcs:                     []template.FuncMap{sprig.FuncMap()},
	}
	if !config.IsDevelopment {
		renderOptions.Asset = Asset
		renderOptions.AssetNames = AssetNames
	}
	r := render.New(renderOptions)
	return r
}
