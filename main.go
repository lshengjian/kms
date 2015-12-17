package main

import (
	"fmt"
	"io/ioutil"

	"io"
	"net/http"
	"os"
	"labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)
type Person struct {
  NAME  string
  PHONE string
}

type Men struct {
  Persons []Person
}

const  (
  URL = "192.168.1.178:27017"
)

func upload(c *echo.Context) error {
	mr, err := c.Request().MultipartReader()
	if err != nil {
		return err
	}

	// Read form field `name`
	part, err := mr.NextPart()
	if err != nil {
		return err
	}
	defer part.Close()
	b, err := ioutil.ReadAll(part)
	if err != nil {
		return err
	}
	name := string(b)

	// Read form field `email`
	part, err = mr.NextPart()
	if err != nil {
		return err
	}
	defer part.Close()
	b, err = ioutil.ReadAll(part)
	if err != nil {
		return err
	}
	email := string(b)

	// Read files
	i := 0
	for {
		part, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		defer part.Close()

		file, err := os.Create("uploads/"+part.FileName())
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(file, part); err != nil {
			return err
		}
		i++
	}
	return c.String(http.StatusOK, fmt.Sprintf("Thank You! %s <%s>, %d files uploaded successfully.",name, email, i))
}

func DbDemo(){
  session, err := mgo.Dial(URL)
  if err != nil {
    panic(err)
  }
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("mydb")
  collection := db.C("person")



  countNum, err := collection.Count()
  if err != nil {
    panic(err)
  }
  fmt.Println("Things objects count: ", countNum)


  temp := &Person{PHONE: "18811577546", NAME:  "zhangzheHero"}
  err = collection.Insert(&Person{"Ale", "+55 53 8116 9639"}, temp)
  if err != nil {
    panic(err)
  }

  result := Person{}
  err = collection.Find(bson.M{"phone": "456"}).One(&result)
  fmt.Println("Phone:", result.NAME, result.PHONE)

  var personAll Men
  iter := collection.Find(nil).Iter()
  for iter.Next(&result) {
    fmt.Printf("Result: %v\n", result.NAME)
    personAll.Persons = append(personAll.Persons, result)
  }

  err = collection.Update(bson.M{"name": "ccc"}, bson.M{"$set": bson.M{"name": "ddd"}})
  //err = collection.Update(bson.M{"name": "ddd"}, bson.M{"$set": bson.M{"phone": "12345678"}})
  //err = collection.Update(bson.M{"name": "aaa"}, bson.M{"phone": "1245", "name": "bbb"})
  _, err = collection.RemoveAll(bson.M{"name": "Ale"})

}

func main() {
        DbDemo()
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Static("/", "public")
	e.Post("/upload", upload)
	print("start at port:3000!");

	e.Run(":3000")
}