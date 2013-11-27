// AlbumController
package controllers

import (
	"fmt"
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
	//	id, _ := strconv.ParseUint(rawId, 10, 64)
	id := models.ToObjectId(rawId)
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
	//	id, _ := strconv.ParseUint(rawId, 10, 64)
	id := models.ToObjectId(rawId)
	genres := models.GetAll(tiedotmartini1.GENRE_COL)
	title := "Add Album"
	t, err := template.ParseFiles("src/tiedotmartini1/views/album/add.html")
	if err != nil {
		panic(err)
	}
	t.Execute(r, struct {
		Title  string
		Genres []models.DocWithID
		Id     uint64
	}{Title: title, Genres: genres, Id: id})
}

func AlbumVerify(params martini.Params, r http.ResponseWriter, rq *http.Request) {
	rawId := params["id"]
	//	id, _ := strconv.ParseUint(rawId, 10, 64)
	id := models.ToObjectId(rawId)
	name := rq.FormValue("name")
	yearString := rq.FormValue("year")
	year, _ := strconv.Atoi(yearString)
	genreType := rq.FormValue("genretype")
	var genreId uint64
	var err error
	errString := "no errors"
	switch genreType {
	case "existing":
		rawGenreId := rq.FormValue("genre_id")
		if rawGenreId == "" {
			errString = "You need to select a genre"
			genreId = 0
		} else {
			//			genreId, _ = strconv.ParseUint(rawGenreId, 10, 64)
			genreId = models.ToObjectId(rawGenreId)
		}
		break
	case "new":
		genreName := rq.FormValue("genre_name")
		fmt.Println("Genre name:", genreName)
		if genreName == "" {
			errString = "You need to enter a name"
			//			genreId = 0
		} else {
			genre := map[string]interface{}{"name": genreName}
			genreId, err = models.AddDoc(genre, tiedotmartini1.GENRE_COL)
			fmt.Println("Attempted to create genre #", genreId)
			if err != nil {
				errString = fmt.Sprintf("Error on genre creation: %s", err.Error())
				fmt.Println(errString)
				fmt.Println("Error on genre creation:", err)
				//				genreId = 0
			}
		}
	default:
		errString = "You need to select an option"
		//		genreId = 0
	}
	if errString == "no errors" {
		band := models.GetDoc(id, tiedotmartini1.BAND_COL)
		//		yearString := rq.FormValue("year")
		//		year, _ := strconv.Atoi(yearString)
		album := models.Album{Name: name, Year: year, GenreId: genreId}
		err := band.AddAlbum(album)
		if err != nil {
			errString = fmt.Sprintf("Error on album addition: %s", err.Error())
		} else {
			id = band.DocKey
		}
	}
	title := "Verifying Album"
	t, err := template.ParseFiles("src/tiedotmartini1/views/album/verify.html")
	if err != nil {
		panic(err)
	}
	t.Execute(r, struct {
		Title   string
		Message string
		Id      uint64
	}{Title: title, Message: errString, Id: id})

}
