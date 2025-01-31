package dtos

import (
	"github.com/google/uuid"
	"time"
)

type Review struct {
	Id          uuid.UUID     `json:"id"`
	Body        string        `json:"body"`
	CreatedTime time.Time     `json:"createdTime"`
	UpdatedTime time.Time     `json:"updatedTime"`
	User        UserForReview `json:"userForReview"`
}
