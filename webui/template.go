package webui

import (
	"fmt"
	"html/template"

	"health_monitor/disk"
	"health_monitor/utils"

	"github.com/valyala/fasthttp"
)

func percent(a int, b int) int {
	return (a * 100) / b
}

func diskTemplateHandler(ctx *fasthttp.RequestCtx, tmpl string) {
	tmpl = fmt.Sprintf(templateRoot, "disk-status")
	funcMap := template.FuncMap{"percent": percent}
	t, err := template.New("disk-status").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = t.Execute(ctx, disk.GetStatus())
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func inodeTemplateHandler(ctx *fasthttp.RequestCtx, tmpl string) {
	tmpl = fmt.Sprintf(templateRoot, "inode-status")
	funcMap := template.FuncMap{"percent": percent}
	t, err := template.New("inode-status").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = t.Execute(ctx, disk.GetStatus())
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}
