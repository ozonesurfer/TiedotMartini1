// main
package main

import (
	//	"fmt"
	//	"github.com/HouzuoGuo/tiedot/db"
	"github.com/codegangsta/martini"
	"tiedotmartini1"
	"tiedotmartini1/controllers"
	"tiedotmartini1/models"
)

func main() {
	//	fmt.Println("Hello World!")
	database := models.GetDB()
	//	database.Drop(tiedotmartini1.BAND_COL)
	//	database.Drop(tiedotmartini1.LOCATION_COL)
	database.Create(tiedotmartini1.BAND_COL)
	database.Create(tiedotmartini1.LOCATION_COL)
	database.Create(tiedotmartini1.GENRE_COL)
	database.Close()
	m := martini.Classic()
	m.Get("/", controllers.HomeIndex)
	m.Get("/home/index", controllers.HomeIndex)
	m.Get("/band/add", controllers.BandAdd)
	m.Post("/band/verify", controllers.BandVerify)
	m.Get("/album/index/:id", controllers.AlbumIndex)
	m.Get("/album/add/:id", controllers.AlbumAdd)
	m.Post("/album/verify/:id", controllers.AlbumVerify)
	m.Use(martini.Static("assets"))
	m.Run()
}
