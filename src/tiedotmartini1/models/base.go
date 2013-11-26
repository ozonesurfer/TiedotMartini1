// base
package models

import (
	"encoding/json"
	"fmt"
	"tiedotmartini1"
	//	"loveoneanother.at/tiedot/db"
	"github.com/HouzuoGuo/tiedot/db"
	"strconv"
)

type MyDoc map[string]interface{}

func ToObjectId(in string) uint64 {
	id, err := strconv.ParseUint(in, 10, 64)
	if err != nil {
		panic(err)
	}
	return id
}

func GetDB() *db.DB {
	myDB, err := db.OpenDB(tiedotmartini1.DATABASE_DIR)
	if err != nil {
		panic(err)
	}
	return myDB
}

type DocWithID struct {
	DocKey uint64
	Value  map[string]interface{}
}

func GetAll(docType string) []DocWithID {
	database := GetDB()
	defer database.Close()
	collection := database.Use(docType)
	var query interface{}
	result := make(map[uint64]struct{})
	json.Unmarshal([]byte(`"all"`), &query)
	db.EvalQueryV2(query, collection, &result)
	var docs []DocWithID
	for id := range result {
		var doc MyDoc
		collection.Read(id, &doc)
		docObj := DocWithID{DocKey: id, Value: doc}
		docs = append(docs, docObj)
	}
	return docs
}

func AddDoc(doc MyDoc, docType string) (uint64, error) {
	database := GetDB()
	collection := database.Use(docType)
	newId, err := collection.Insert(doc)
	return newId, err
}

func GetDoc(id uint64, docType string) DocWithID {
	database := GetDB()
	defer database.Close()
	collection := database.Use(docType)
	var value MyDoc
	collection.Read(id, &value)
	doc := DocWithID{DocKey: id, Value: value}
	return doc
}

type Album struct {
	Name    string `json:"album_name"`
	Year    int    `json:"year"`
	GenreId uint64 `json:"genre_id"`
}

type Band struct {
	Name       string  `json:"name"`
	LocationId uint64  `json:"location_id"`
	Albums     []Album `json:"albums"`
}

type Genre struct {
	Name string `json:"name"`
}

type Location struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

func (this *Album) GetGenreName() string {
	rawDoc := GetDoc(this.GenreId, tiedotmartini1.GENRE_COL)
	return rawDoc.Value["name"].(string)

}

func GetGenreName(id uint64) string {
	rawDoc := GetDoc(id, tiedotmartini1.GENRE_COL)
	//	genre := rawDoc.Value.(map[string]interface{})
	return rawDoc.Value["name"].(string)
}

func GetBandsByGenre(id uint64) []DocWithID {
	database := GetDB()
	defer database.Close()
	collection := database.Use(tiedotmartini1.BAND_COL)
	var query interface{}
	q := `{"eq": ` + strconv.FormatUint(id, 10) + `, "in": ["albums", "genre_id"]}`
	//	json.Unmarshal([]byte(`"all"`), &query)
	json.Unmarshal([]byte(q), &query)
	result := make(map[uint64]struct{})
	db.EvalQueryV2(query, collection, &result)
	fmt.Println("Ran query")
	var docs []DocWithID
	for id2, value := range result {
		fmt.Println("Found", id2, value)
		var readback map[string]interface{}
		collection.Read(id2, &readback)
		doc := DocWithID{DocKey: id2, Value: readback}
		docs = append(docs, doc)
	}
	if docs == nil {
		fmt.Println("Returning empty value")
	}
	return docs
}

func (this *DocWithID) LocToString() string {
	location := this.Value
	city := location["city"].(string)
	state := location["state"].(string)
	country := location["country"].(string)
	var strCity, strState string
	if city == "" {
		strCity = "(city)"
	} else {
		strCity = city
	}
	if state == "" {
		strState = "(state/province)"
	} else {
		strState = state
	}
	result := strCity + ", " + strState + " " + country
	return result
}

func (this DocWithID) GetLocation() string {
	//	original := this.Value.(Band)
	original := this.Value
	id := original["location_id"].(float64)
	fmt.Println("id =", id)
	rawDoc := GetDoc(uint64(id), tiedotmartini1.LOCATION_COL)
	result := rawDoc.LocToString()
	//	result := this
	return result
}

func (this *DocWithID) GetName() string {
	original := this.Value
	name := original["name"].(string)
	return name
}
func (this *DocWithID) AddAlbum(album Album) error {
	database := GetDB()
	defer database.Close()
	collection := database.Use(tiedotmartini1.BAND_COL)
	original := this.Value
	x := original["location_id"].(float64)
	band := Band{Name: original["name"].(string),
		LocationId: uint64(x)}
	band.Albums = []Album{}
	if original["albums"] != nil {
		for _, a := range original["albums"].([]interface{}) {
			x := a.(map[string]interface{})
			z := x["genre_id"].(float64)
			y := x["year"].(float64)
			q := Album{Name: x["album_name"].(string), Year: int(y),
				GenreId: uint64(z)}
			band.Albums = append(band.Albums, q)
		}
	}
	band.Albums = append(band.Albums, album)
	newKey, err := collection.Update(this.DocKey, band)
	this.DocKey = newKey
	return err
}

func (this DocWithID) GetAlbums() []Album {
	original := this.Value
	var cds []Album
	for _, a := range original["albums"].([]interface{}) {
		x := a.(map[string]interface{})
		q := Album{Name: x["album_name"].(string), Year: int(x["year"].(float64)),
			GenreId: uint64(x["genre_id"].(float64))}
		cds = append(cds, q)
	}
	return cds
}
