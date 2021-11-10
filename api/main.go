package main

import (
	"fmt"
	"net"

	"net/http"

	"Himawari/controllers"
	"Himawari/models/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
)

func init() {
	fmt.Println(" _   _ _                                    _ ")
	fmt.Println("| | | (_)_ __ ___   __ ___      ____ _ _ __(_)")
	fmt.Println("| |_| | | '_ ` _ \\ / _` \\ \\ /\\ / / _` | '__| |")
	fmt.Println("|  _  | | | | | | | (_| |\\ V  V / (_| | |  | |")
	fmt.Println("|_| |_|_|_| |_| |_|\\__,_| \\_/\\_/ \\__,_|_|  |_|")
	fmt.Println("")
	fmt.Println("MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM")
	fmt.Println("MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM")
	fmt.Println("MMMMMMMMMMMMMMMMMM^_?MMMMMMM9\"TMMMMMMMMMMMMMMM")
	fmt.Println("MMMMMMMMMMMMMMMM#\\    ?MM@'   .MMMMMMMMMMMMMMM")
	fmt.Println("MMMMMMMMMMMHMMMM#      jD      MMMMMMMMMMMMMMM")
	fmt.Println("MMMMMMMMMb     ?T,  t  .: .'  .H\"!~`~?MMMMMMMM")
	fmt.Println("MMMMMMMMM#.      ,-JZVOOOV0-..^      .MMMMMMMM")
	fmt.Println("MMMMMMMMMMp.  _&vOllttOOOtllOn./!   .MMMMMMMMM")
	fmt.Println("MMMMMMMB\"7771.,OllOrrrrvrrrOllZ, ..7TMMMMMMMMM")
	fmt.Println("MMMMNr     . .Iltrrvrvrrvrvrwlld}      7MMMMMM")
	fmt.Println("MMMMMNo.     (OlOvrvrvvrrvrvvZltf7?!    .MMMMM")
	fmt.Println("MMMMMMMNm-...jZlOrvrrrvrvrrvrZlln.....jNMMMMMM")
	fmt.Println("MMMMMMMMMY`  .XllOvrvrrvrvrrOlld:  .7MMMMMMMMM")
	fmt.Println("MMMMMMM@!  .^  4llOrvvrrvrOOtlO^ ?`   TMMMMMMM")
	fmt.Println("MMMMMM#>     ..=?wtlltttlllOO=?(.......MMMMMMM")
	fmt.Println("MMMMMMN&JJ+g#=  (`.1C7O774>`.l  ,MMMMMMMMMMMMM")
	fmt.Println("MMMMMMBBYYWMt     .Z  j  .L      vHBYHMMMMMMMM")
	fmt.Println("MMMR>>>>>>>j\\   .gMb     .Mm..   (>>>>>>>dMMMM")
	fmt.Println("MMMNn>;>;>++6vC<<XMN,   .MMB1<+111<>;>>;jMMMMM")
	fmt.Println("MMMMNy>;>>;?T&+>>>zMN,.dMM6;>>+&vC>>;>>jMMMMMM")
	fmt.Println("MMMMMNm+>>>>>><?Tl>OMb~jMD>1vC<>>>>;>+uMMMMMMM")
	fmt.Println("MMMMMMMNmx+>;>>;++ugMR.JMme+>>;>>;+&dNMMMMMMMM")
	fmt.Println("MMMMMMMMMMMNNNNNMMMMMR.jMMMMNNNgNNMMMMMMMMMMMM")
	fmt.Println("MMMMMMMMMMMMMMMMMMMMMR((MMMMMMMMMMMMMMMMMMMMMM")
	fmt.Println("MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM")
	fmt.Println("MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM")
	logger.LoggingSetting()
}

func isRunSunflower() string {
	sunflower := "localhost:18080"
	_, err := net.Dial("tcp", sunflower)
	if err != nil {
		return "http://localhost:8080/"
	} else {
		return "http://localhost:8080/Himawari/?url=http%3A%2F%2Flocalhost%3A18080%2F"
	}
}

func openBrowser(target string) {
	err := browser.OpenURL(target)
	if logger.ErrHandle(err) {
		//		panic(err)
	}
}

func main() {
	//target := isRunSunflower()
	go openBrowser("http://localhost:3000/")
	go openBrowser(isRunSunflower())
	router := gin.Default()
	router.Static("/views", "./views")

	router.StaticFS("/Himawari", http.Dir("./views/static"))

	api := router.Group("/api")
	{
		api.GET("/sitemap", controllers.ReadSitemap)
		api.POST("/crawl", controllers.Crawl)
		api.GET("/outoforigin", controllers.ExportOutOfOrigin)
		api.GET("/sort", controllers.Sort)
		api.GET("/reset", controllers.Reset)
		api.GET("/scan", controllers.Scan)
	}
	//	router.POST("/api/deletePath", controller.DeletePath)
	router.GET("/download", controllers.DownloadSitemap)
	router.POST("/upload", controllers.UploadSitemap)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/Himawari")
	})
	router.Run(":8080")
}
