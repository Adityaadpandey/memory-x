package types

import "time"

type User struct {
	Id         string      `json:"id"`
	Name       string      `json:"name" validate:"required"`
	Email      string      `json:"email" validate:"required,email"`
	Password   string      `json:"password,omitempty"` // optional
	Workspaces []Workspace `json:"workspaces,omitempty"`
	Stories    []Story     `json:"stories,omitempty"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

type Workspace struct {
	Id          string    `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	Caption     string    `json:"caption,omitempty"`
	Users       []User    `json:"users,omitempty"`
	Images      []Image   `json:"images,omitempty"`
	Stories     []Story   `json:"stories,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Image struct {
	Id          string     `json:"id"`
	URL         string     `json:"url" validate:"required,url"`
	StoryID     *string    `json:"storyId,omitempty"`
	WorkspaceID *string    `json:"workspaceId,omitempty"`
	Story       *Story     `json:"story,omitempty"`
	Workspace   *Workspace `json:"workspace,omitempty"`
	Captions    []Caption  `json:"captions,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type Caption struct {
	Id        string    `json:"id"`
	Text      string    `json:"text" validate:"required"`
	ImageID   string    `json:"imageId"`
	Image     *Image    `json:"image,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Story struct {
	Id          string     `json:"id"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description,omitempty"`
	WorkspaceID string     `json:"workspaceId"`
	UserID      string     `json:"userId"`
	Caption     string     `json:"caption,omitempty"`
	Images      []Image    `json:"images,omitempty"`
	User        *User      `json:"user,omitempty"`
	Workspace   *Workspace `json:"workspace,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
