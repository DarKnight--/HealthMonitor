// Package webui implements the web interface for the OWTF monitor module
package webui

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/valyala/fasthttp"

	"health_monitor/api"
	"health_monitor/setup"
	"health_monitor/utils"
)

var (
	staticRoot   = path.Join("webui", "static", "%s")
	templateRoot = path.Join("webui", "templates", "%s")
	logFile      *os.File
)

// RunServer will serve the webui content
func RunServer(port string) {
	var err error
	logFileName := path.Join(setup.ConfigVars.HomeDir, "webui.log")
	logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		utils.PLogError(err)
	}
	defer logFile.Close()
	if err = fasthttp.ListenAndServe(":"+port, requestHandler); err != nil {
		utils.ModuleError(logFile, "Unable to run the server", err.Error())
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	tempPath := strings.SplitN(string(ctx.Path()), "/", 3)
	if len(tempPath) == 2 {
		render(ctx, "index.html")
		return
	}
	switch tempPath[1] {
	case "static":
		staticHandler(ctx, tempPath[2])
	case "settings": // Serves the json data of the module's config.
		configHandler(ctx, tempPath[2])
	case "preferences":
		render(ctx, "settings.html")
	case "description": //Serves the page for serving modal
		if strings.HasSuffix(tempPath[2], "html") {
			render(ctx, tempPath[2])
		} else {
			render(ctx, tempPath[2]+"-setting")
		}
	case "moduleStatus":
		moduleStatusHandler(ctx, tempPath[2])
	default:
		if api.ModuleStatus(tempPath[2]) || tempPath[2] == "main" {
			switch tempPath[1] {
			case "module": // Serves the json data of the module's status.
				statusHandler(ctx, tempPath[2])
			case "template": // Serves the template for short description
				templateHandler(ctx, tempPath[2])
			default:
				ctx.Error("not found", fasthttp.StatusNotFound)
			}
		} else {
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
}

func render(ctx *fasthttp.RequestCtx, tmpl string) {
	tmpl = fmt.Sprintf(templateRoot, tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		utils.ModuleError(logFile, "template parsing error ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	err = t.Execute(ctx, "")
	if err != nil {
		utils.ModuleError(logFile, "template executing error: ", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.SetContentType("text/html; charset=utf-8")
}

func staticHandler(ctx *fasthttp.RequestCtx, filePath string) {
	filePath = fmt.Sprintf(staticRoot, filePath)
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		fasthttp.ServeFile(ctx, filePath)
		utils.ModuleLogs(logFile, fmt.Sprintf("[200] File: %s", ctx.Path()))
		return
	}
	utils.ModuleLogs(logFile, fmt.Sprintf("[404] Unable to find the requested static file: %s",
		ctx.Path()))
	ctx.NotFound()
}

func statusHandler(ctx *fasthttp.RequestCtx, module string) {
	if _, ok := api.StatusFunc[module]; ok {
		ctx.SetContentType("application/json")
		ctx.SetBody(api.GetStatusJSON(module))
		return
	}
	utils.ModuleLogs(logFile, fmt.Sprintf("[404] Unable to find the requested json: %s",
		ctx.Path()))
	ctx.NotFound()
}

func configHandler(ctx *fasthttp.RequestCtx, module string) {
	if ctx.IsPost() {
		if _, ok := api.ConfSaveFunc[module]; ok {
			err := api.SaveConfig(module, ctx.PostBody())
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				utils.ModuleError(logFile, fmt.Sprintf("[404] Unable to save data: %s",
					ctx.Path()), err.Error())
			}
			return
		}
		utils.ModuleLogs(logFile, fmt.Sprintf("[404] Unable to find the requested module: %s",
			module))
	}
	if _, ok := api.ConfFunc[module]; ok {
		ctx.SetContentType("application/json")
		ctx.SetBody(api.GetConfJSON(module))
		return
	}
	utils.ModuleLogs(logFile, fmt.Sprintf("[404] Unable to find the requested json: %s",
		ctx.Path()))
	ctx.NotFound()
}

func templateHandler(ctx *fasthttp.RequestCtx, tmpl string) {
	switch tmpl {
	case "live":
		liveTemplateHandler(ctx)
	case "disk":
		diskTemplateHandler(ctx)
	case "inode":
		inodeTemplateHandler(ctx)
	case "ram":
		ramTemplateHandler(ctx)
	case "cpu":
		cpuTemplateHandler(ctx)
	case "target":
		targetTemplateHandler(ctx)
	default:
		utils.ModuleLogs(logFile, fmt.Sprintf("[404] Unable to find the requested template: %s",
			ctx.Path()))
		ctx.NotFound()
	}
}

func moduleStatusHandler(ctx *fasthttp.RequestCtx, module string) {
	if string(ctx.PostBody()) == "1" {
		api.ChangeModuleStatus(module, true)
		utils.ModuleLogs(logFile, fmt.Sprintf("Turning on %s module",
			module))
	} else if string(ctx.PostBody()) == "0" {
		api.ChangeModuleStatus(module, false)
		utils.ModuleLogs(logFile, fmt.Sprintf("Turning off %s module",
			module))
	} else {
		ctx.NotFound()
		return
	}
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}
