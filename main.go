package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	port = flag.String("p", "80", "port number")
)

func main() {
	flag.Parse()
	address := fmt.Sprintf(":%s", *port)
	fp := getFilePath()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.HideBanner = true
	e.HidePort = true
	static(e, "/", fp)
	e.Logger.Fatal(e.Start(address))
}

func getFilePath() string {
	length := len(os.Args)
	if length == 0 {
		dir, err := os.Getwd()
		if err != nil {
			return "./"
		}
		return dir
	}
	return os.Args[length-1]
}

func static(e *echo.Echo, prefix, root string) {
	h := func(c echo.Context) error {
		p, err := url.PathUnescape(c.Param("*"))
		if err != nil {
			return err
		}
		name := filepath.Join(root, path.Clean("/"+p)) // "/"+ for security
		return c.File(name)
	}
	e.Any(prefix, h)
	if prefix == "/" {
		e.Any(prefix+"*", h)
		return
	}

	e.Any(prefix+"/*", h)
	return
}
