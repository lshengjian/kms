package main

import (
   	"strconv"
    "net/http"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/lshengjian/kms/util"
)

var (
	web WebOptions
)

type WebOptions struct {
	filerServer *string
    masterServer *string
    port *int
}

func init() {
	cmdWeb.ExecuteFunc = runWeb // break init cycle
	web.filerServer = cmdWeb.Flag.String("filerServer", "localhost:8888", "SeaweedFS filer location")
    web.masterServer = cmdWeb.Flag.String("masterServer", "localhost:9333", "SeaweedFS master location")
    web.port = cmdWeb.Flag.Int("port", 3000, "kms server http listen port")
}

var cmdWeb = &Command{
	UsageLine: "web -filerServer=localhost:8888",
	Short:     "show files in filer server",
	Long: `show files in filer server.

  first start seaweedfs filer server.

  `,
}

func runWeb(cmd *Command, args []string) bool {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Static("/", "public")
	e.Index("public/index.html")

	e.Favicon("public/favicon.ico")

	e.Post("/upload", upload)

	e.Get("/hello", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello!")
	})

	util.Info("start at port:",*web.port);
    
    
	e.Run(":"+strconv.Itoa(*web.port))
	return true
}

