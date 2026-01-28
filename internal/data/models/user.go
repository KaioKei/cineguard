package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users, alias:u"`
	ID            int64     `bun:"id,pk,type:bigserial,notnull"`
	Name          string    `bun:"name,type:text,notnull"`
	Email         string    `bun:"email,unique,type:text,notnull"`
	CreatedAt     time.Time `bun:"createdAt,type:time,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updatedAt,type:time,notnull,default:current_timestamp"`
}

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type Profile struct {
	bun.BaseModel `bun:"table:profiles, alias:p"`
	ID            int64  `bun:"id,pk,type:bigserial,notnull"`
	UserID        int64  `bun:"user_id,type:bigserial,notnull"`
	Bio           string `bun:"bio,type:text,notnull"`
	Avatar        string

	// Relationship
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
