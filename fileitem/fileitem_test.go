package fileitem

import (
	"fmt"
	"testing"
	"time"
//        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)





func Test01AddItem(t *testing.T) {
  collection :=GetCollection("fileitems")
  var tags = []string{"1-1","3-2"}
  var vs = []VisitInfo{VisitInfo{Time:time.Now(),Visitor:"tom"}}
  temp := &FileItem{
	  Title:"demo",
	  CreatedTime:time.Now(),
	  Owner:"alex",
	  Tags : tags,
	  Uses:  vs,
  }

  showerr(collection.Insert(temp));

}


func Test02UpdateItem(t *testing.T) {
  collection :=GetCollection("fileitems")
   var tags = []string{"5","8"}
   collection.Update(bson.M{"title": "demo"},bson.M{"$addToSet": bson.M{"tags":bson.M{"$each": tags}}})//,bson.M{"$slice":-2}
  //err = collection.Update(bson.M{"name": "ddd"}, bson.M{"$set": bson.M{"phone": "12345678"}})

  //_, err = collection.RemoveAll(bson.M{"name": "Ale"})
 
}
func test06RemoveTags(t *testing.T) {
  collection :=GetCollection("fileitems")
  var tags = []string{}
  collection.Update(bson.M{"title": "demo"},bson.M{ "$set": bson.M{"tags":tags} })
 
}

func Test91Findtem(t *testing.T) {
 collection :=GetCollection("fileitems")
  result := &FileItem{}
  collection.Find(bson.M{"title": "demo"}).One(result)
  fmt.Printf("%v",result)

}
func Test92Count(t *testing.T) {
  collection :=GetCollection("fileitems")
  countNum, _ := collection.Count()
  fmt.Println("count: ", countNum)
}

func Test99FindAll(t *testing.T) {
  collection :=GetCollection("fileitems")
  iter := collection.Find(nil).Iter()
  result := &FileItem{}
  for iter.Next(result) {
    fmt.Printf("Result: %v\n", result)
  }
}