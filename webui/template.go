package webui

import (
	"fmt"
	"html/template"

	"github.com/owtf/health_monitor/api"
	"github.com/owtf/health_monitor/utils"

	"github.com/valyala/fasthttp"
)

func percent(a int, b int) int {
	if b == 0 {
		return 0
	}
	return (a * 100) / b
}

func diskTemplateHandler(ctx *fasthttp.RequestCtx) {
	tmpl := fmt.Sprintf(templateRoot, "disk-status")
	funcMap := template.FuncMap{"percent": percent}
	t, err := template.New("disk-status").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = t.Execute(ctx, api.DiskStatus())
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func inodeTemplateHandler(ctx *fasthttp.RequestCtx) {
	tmpl := fmt.Sprintf(templateRoot, "inode-status")
	funcMap := template.FuncMap{"percent": percent}
	t, err := template.New("inode-status").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = t.Execute(ctx, api.DiskStatus())
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func liveTemplateHandler(ctx *fasthttp.RequestCtx) {
	tmpl := fmt.Sprintf(templateRoot, "live-status")
	t, err := template.New("live-status").ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = t.Execute(ctx, api.LiveStatus().Normal)
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func cpuTemplateHandler(ctx *fasthttp.RequestCtx) {
	tmpl := fmt.Sprintf(templateRoot, "cpu-status")
	funcMap := template.FuncMap{"percent": percent}
	t, err := template.New("cpu-status").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = t.Execute(ctx, api.CPUStatus())
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func ramTemplateHandler(ctx *fasthttp.RequestCtx) {
	tmpl := fmt.Sprintf(templateRoot, "ram-status")
	funcMap := template.FuncMap{"percent": percent}
	t, err := template.New("ram-status").Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = t.Execute(ctx, api.RAMStatus())
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func targetTemplateHandler(ctx *fasthttp.RequestCtx) {
	tmpl := fmt.Sprintf(templateRoot, "target-status")
	t, err := template.New("target-status").ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = t.Execute(ctx, api.TargetStatus())
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func settingProfileHandler(ctx *fasthttp.RequestCtx) {
	tmpl := fmt.Sprintf(templateRoot, "main-setting")
	t, err := template.New("main-setting").ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	type vars struct {
		CurrentProfile string
		AllProfiles    []string
	}

	err = t.Execute(ctx, vars{CurrentProfile: api.GetActiveProfile(), AllProfiles: api.GetAllProfiles()})
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}
