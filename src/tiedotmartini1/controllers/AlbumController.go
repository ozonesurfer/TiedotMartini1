// AlbumController
package controllers

import (
	//	"fmt"
	"github.com/codegangsta/martini"
	"html/template"
	"net/http"
	"strconv"
	"tiedotmartini1"
	"tiedotmartini1/models"
)

/* func main() {
	fmt.Println("Hello World!")
}
*/

func AlbumIndex(params martini.Params, r http.ResponseWriter) {
	rawId := params["id"]
	id, _ := strconv.ParseUint(rawId, 10, 64)
	band := models.GetDoc(id, tiedotmartini1.BAND_COL)
	title := "Albums by " + band.Value["name"].(string)
	t, err := template.ParseFiles("src/tiedotmartini1/views/album/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(r, struct {
		Title string
		Band  models.DocWithID
		Id    uint64
	}{Title: title, Band: band, Id: id})
}

func AlbumAdd(params martini.Params, r http.ResponseWriter) {
	rawId := params["id"]
	id, _ := strconv.ParseUint(rawId, 10, 64)
    genres := models.GetAll(tiedotmartini1.GENRE_COL)
	title := "Add Album"
	t, err := template.ParseFiles("src/tiedotmartini1/views/album/add.html")
	if err != nil {
		         panic(err)
	}
	t.Execute(r, struct{
			Title string
			Genres []models.DocWithID
			Id uint64
		}{Title: title, Genres: genres, Id: id})
}

 func AlbumVerify(params martini.Params, r http.ResponseWriter, rq *http.Request) {
	 rawId := params["id"]
	 id, _ := strconv.ParseUint(rawId, 10, 64)
	 name := rq.FormValue("name")
	 genreType := rq.FormValue("genretype")
	 var (
		 genreId uint64
	 )
	 errString := "no errors"
	 switch genreType {
	 case "existing":
		 rawGenreId := rq.FormValue("genre_id")
		 if rawGenreId == "" {
			 errString = "You need to select a genre"
			 genreId = 0
		 } else {
			 genreId, _ = strconv.ParseUint(rawGenreId, 10, 64)
		 }
	 case "new":
		 genreName := rq.FormValue("genre_name")
		 if genreName == "" {
			 errString = "You need to enter a name"
			 genreId = 0
		 } else {
			 genre := map[string]interface {}{"name": genreName}
			 i, err := models.AddDoc(genre, tiedotmartini1.GENRE_COL)
			 if err != nil {
				 errString = "Error on genre creation: " + err.Error()
				 genreId = 0
			 } else {
				 genreId = i
			 }
		 }
	 default:
		 errString = "You need to select an option"
		 genreId = 0
	 }
	 if errString == "no errors" {
		 band := models.GetDoc(id, tiedotmartini1.BAND_COL)
		 yearString := rq.FormValue("year")
		 year, _ := strconv.Atoi(yearString)
		 album := models.Album{Name: name, GenreId: genreId, Year: year}
		 err := band.AddAlbum(album)
		 if err != nil {
			 errString = "Error of album addition: " + err.Error()
		 } else {
			 id = band.DocKey             // changing the date might change the band's key value
		 }
	 }
	 title:= "Verifying Album"
	 t, err := template.ParseFiles("src/tiedotmartini1/views/album/verify.html")
	 if err != nil {
		 panic(err)
	 }
	 t.Execute(r, struct {
			 Title string
			 Message string
			 Id uint64
		 }{Title: title, Message: errString, Id: id})

 }
