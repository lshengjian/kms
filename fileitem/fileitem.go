package fileitem

import (
	"time"
)

type VisitInfo struct {
   Time time.Time `bson:"time"`
   Visitor string  `bson:"userName"`
}
type FileItem struct {
  Title  string  `bson:"title"`
  CreatedTime time.Time   `bson:"time"`
  Owner  string  `bson:"owner"`
  Tags  []string `bson:"tags"`
  Uses  []VisitInfo `bson:"uses"`
}

