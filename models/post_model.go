package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	ID            uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID        uuid.UUID          `gorm:"type:uuid;not null" json:"userID"`
	User          User               `gorm:"" json:"user"`
	Content       string             `gorm:"type:text;not null" json:"content"`
	CreatedAt     time.Time          `gorm:"default:current_timestamp" json:"postedAt"`
	Images        pq.StringArray     `gorm:"type:text[]" json:"images"`
	Hashes        pq.StringArray     `gorm:"type:text[]" json:"hashes"`
	NoShares      int                `gorm:"default:0" json:"noShares"`
	NoLikes       int                `gorm:"default:0" json:"noLikes"`
	NoComments    int                `gorm:"default:0" json:"noComments"`
	Tags          pq.StringArray     `gorm:"type:text[]" json:"tags"`
	Edited        bool               `gorm:"default:false" json:"edited"`
	Comments      []Comment          `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments"`
	Notifications []Notification     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	Messages      []Message          `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	Likes         []UserPostLike     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	BookMarkItems []PostBookmarkItem `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	CommentLikes  []UserCommentLike  `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	UsersTagged   []UserPostTag      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserPostLike struct { //TODO create a single liked model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"likedByID"`
	User      User      `gorm:"" json:"likedBy"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"postID"`
	Post      Post      `gorm:"" json:"post"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"likedAt"`
}

type UserPostTag struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User   User      `gorm:"" json:"user"`
	PostID uuid.UUID `gorm:"type:uuid;not null" json:"postID"`
}
