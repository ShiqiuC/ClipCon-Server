package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ClipboardItem 存储剪贴板内容的结构
type ClipboardItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"createdAt"`
}
