package main

import (
	"fmt"
	"io"
	"os"

	"log"
	"net/http"
	"time"

	"Himawari/controllers"

	"github.com/gin-gonic/gin"
)

func loggingSetting() {
	layout := "2006-01-02_15:04:05"
	dirName := "log"
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.Mkdir(dirName, 0666)
	}
	t := time.Now()
	fileName := "log/" + t.Format(layout) + ".log"
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetFlags(log.Flags() &^ log.LstdFlags)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetPrefix("======================================================\n")
}

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
	loggingSetting()
}

func main() {
	/*
		proxyUrl, err := url.Parse("http://172.16.82.190:8001")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	*/

	router := gin.Default()
	router.Static("/views", "./views")

	router.StaticFS("/Himawari", http.Dir("./views/static"))

	api := router.Group("/api")
	{
		api.GET("/sitemap", controllers.ReadSitemap)
		api.POST("/crawl", controllers.Crawl)
		api.GET("/found", controllers.FoundItem)
		api.GET("/sort", controllers.Sort)
	}
	//	router.POST("/api/deletePath", controller.DeletePath)
	router.GET("/download", controllers.DownloadSitemap)
	router.POST("/upload", controllers.UploadSitemap)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/Himawari")
	})
	router.Run(":8080")
}
