package webui

import (
	"fmt"
	"html/template"

	"health_monitor/disk"

	"github.com/valyala/fasthttp"
)

func percent(a int, b int) int {
	return (a * 100) / b
}

func diskTemplateHandler(ctx *fasthttp.RequestCtx) {
	tmpl := fmt.Sprintf(templateRoot, "disk-status")
	funcMap := template.FuncMap{"percent": percent}
	t, err := template.New("disk-status").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		fmt.Print("template parsing error: ", err)
		ctx.NotFound()
		return
	}

	err = t.Execute(ctx, disk.GetStatus())
	if err != nil {
		fmt.Print("template executing error: ", err)
	}
}
