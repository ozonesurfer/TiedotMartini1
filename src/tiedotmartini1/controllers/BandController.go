// BandController
package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"tiedotmartini1"
	"tiedotmartini1/models"
)

/*func main() {
	fmt.Println("Hello World!")
}
*/

func BandAdd(r http.ResponseWriter, rq *http.Request) {
	locations := models.GetAll(tiedotmartini1.LOCATION_COL)
	t, err := template.ParseFiles("src/tiedotmartini1/views/band/add.html")
	if err != nil {
		panic(err)
	} else {
		t.Execute(r, struct {
			Title     string
			Locations []models.DocWithID
		}{Title: "Adding A Band", Locations: locations})
	}
}

func BandVerify(r http.ResponseWriter, rq *http.Request) {
	name := rq.FormValue("name")
	locType := rq.FormValue("loctype")
	var locationId uint64
	errorString := "no errors"
	switch locType {
	case "existing":
		locationIdString := rq.FormValue("location_id")
		if locationIdString == "" {
			errorString = "No location was selected"
		} else {
			locationId = models.ToObjectId(locationIdString)
		}
		break
	case "new":
		if rq.FormValue("country") == "" {
			errorString = "Country is required"
		} else {
			location := map[string]interface{}{"city": rq.FormValue("city"),
				"state":   rq.FormValue("state"),
				"country": rq.FormValue("country")}
			id, err := models.AddDoc(location, tiedotmartini1.LOCATION_COL)
			if err != nil {
				errorString = fmt.Sprintf("Add location error: %s", err.Error())
				locationId = 0
			} else {
				locationId = id
			}
		}
		break
	default:
		errorString = "You need to select an option"
		locationId = 0
	}
	if errorString == "no errors" {
		band := map[string]interface{}{"name": name,
			"location_id": locationId,
			"albums":      []models.Album{}}
		_, err := models.AddDoc(band, tiedotmartini1.BAND_COL)
		if err != nil {
			errorString = fmt.Sprintf("Add band error: %s", err.Error())
		}
	}
	t, err := template.ParseFiles("src/tiedotmartini1/views/band/verify.html")
	if err != nil {
		panic(err)
	} else {
		t.Execute(r, struct {
			Title   string
			Message string
		}{Title: "Verifying Band",
			Message: errorString})
	}
}
