// Package iam provides structures and functionality related to AWS Identity and Access Management (IAM) groups.
package iam

import "time"

// HTTP error messages related to IAM groups.
const (
	HttpErrGroupNotFound = "group not found"
	HttpErrUpdatingGroup = "error updating group"
	HttpErrCreatingGroup = "error creating group"
)

// Group represents an IAM group.
type Group struct {
	Arn        string    `json:"arn"`
	CreateDate time.Time `json:"created_date"`
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
}

// GroupDetailResponse represents a response detailing an IAM group.
type GroupDetailResponse struct {
	GroupName  string `json:"name"`
	GroupID    string `json:"id"`
	CreateDate string `json:"created_date"`
}
