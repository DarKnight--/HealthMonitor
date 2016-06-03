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
)

var (
	staticRoot   = path.Join("webui", "static", "%s")
	templateRoot = path.Join("webui", "templates", "%s")
)

// RunServer will serve the webui content
func RunServer(port string) {
	if err := fasthttp.ListenAndServe(":"+port, requestHandler); err != nil {
		fmt.Println("error in running server")
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
		case "module":
			statusHandler(ctx, tempPath[2])
		case "template":
			templateHandler(ctx, tempPath[2])
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
}

func render(ctx *fasthttp.RequestCtx, tmpl string) {
	tmpl = fmt.Sprintf(templateRoot, tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("template parsing error: ", err)
	}
	err = t.Execute(ctx, "")
	if err != nil {
		fmt.Print("template executing error: ", err)
	}
	ctx.Response.Header.SetContentType("text/html; charset=utf-8")
}

func staticHandler(ctx *fasthttp.RequestCtx, filePath string) {
	filePath = fmt.Sprintf(staticRoot, filePath)
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		fasthttp.ServeFile(ctx, filePath)
		return
	}
	ctx.NotFound()
}

func statusHandler(ctx *fasthttp.RequestCtx, module string) {
	if status, ok := api.StatusFunc[module]; ok {
		ctx.SetContentType("application/json")
		ctx.SetBody(status())
		return
	}
	ctx.NotFound()
}

func templateHandler(ctx *fasthttp.RequestCtx, tmpl string) {
	switch tmpl {
	case "disk-status":
		diskTemplateHandler(ctx, tmpl)
	case "inode-status":
		inodeTemplateHandler(ctx, tmpl)
	default:
		ctx.NotFound()
	}
	ctx.Response.Header.Add("module", tmpl)
}
