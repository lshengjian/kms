package fileitem

import (
//	"time"
        "gopkg.in/mgo.v2"
//        "gopkg.in/mgo.v2/bson"
)
const  (
  DB_URL = "192.168.1.178:27017"
)
var (
  session *mgo.Session 
)

func init(){
	//print("init")
	se, err := mgo.Dial(DB_URL)
        showerr(err)
	session=se
	session.SetMode(mgo.Monotonic, true)

}
func closeSession(){
	session.Close()
}
func showerr(err error){
  if err != nil {
    panic(err)
  }
}


func GetCollection(name string)(*mgo.Collection){
  db := session.DB("km_db")
  return db.C(name)
}