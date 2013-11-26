// HomeController
package controllers

import (
	//	"fmt"
	"html/template"
	"net/http"
	"tiedotmartini1"
	"tiedotmartini1/models"
)

/* func main() {
	fmt.Println("Hello World!")
} */

func HomeIndex(r http.ResponseWriter, rw *http.Request) {
	bands := models.GetAll(tiedotmartini1.BAND_COL)
	t, err := template.ParseFiles("src/tiedotmartini1/views/home/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(r, struct {
		Title string
		Bands []models.DocWithID
	}{Title: "My CD Catalog", Bands: bands})
}
