// Package webui implements the web interface for the OWTF monitor module
package webui

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/valyala/fasthttp"
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
