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
	} else {
		switch tempPath[1] {
		case "static":
			staticHandler(ctx, tempPath[2])
		case "module": // Serves the json data of the module's status.
			statusHandler(ctx, tempPath[2])
		case "template": // Serves the template for short description
			templateHandler(ctx, tempPath[2])
		case "description": //Serves the page for serving modal
			render(ctx, tempPath[2])
		case "settings": // Serves the json data of the module's config.
			configHandler(ctx, tempPath[2])
		case "preferences":
			render(ctx, "settings.html")
		default:
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
	if status, ok := api.StatusFunc[module]; ok {
		ctx.SetContentType("application/json")
		ctx.SetBody(status())
		return
	}
	utils.ModuleLogs(logFile, fmt.Sprintf("[404] Unable to find the requested json: %s",
		ctx.Path()))
	ctx.NotFound()
}

func configHandler(ctx *fasthttp.RequestCtx, module string) {
	if ctx.IsPost() {
		if status, ok := api.ConfSaveFunc[module]; ok {
			err := status(ctx.PostBody())
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
	if status, ok := api.ConfFunc[module]; ok {
		ctx.SetContentType("application/json")
		ctx.SetBody(status())
		return
	}
	utils.ModuleLogs(logFile, fmt.Sprintf("[404] Unable to find the requested json: %s",
		ctx.Path()))
	ctx.NotFound()
}

func templateHandler(ctx *fasthttp.RequestCtx, tmpl string) {
	switch tmpl {
	case "disk-status":
		diskTemplateHandler(ctx, tmpl)
	case "inode-status":
		inodeTemplateHandler(ctx, tmpl)
	default:
		utils.ModuleLogs(logFile, fmt.Sprintf("[404] Unable to find the requested template: %s",
			ctx.Path()))
		ctx.NotFound()
	}
}
