package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"rishik.com/enums"
)

type Account struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Provider     enums.ProviderEnum `bson:"provider"`
	ProviderID   string             `bson:"providerId"`
	UserID       primitive.ObjectID `bson:"userId"`
	RefreshToken *string            `bson:"refreshToken,omitempty"`
	TokenExpiry  *time.Time         `bson:"tokenExpiry,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}
type Member struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"userId"`
	WorkspaceID primitive.ObjectID `bson:"workspaceId"`
	RoleID      primitive.ObjectID `bson:"role"`
	JoinedAt    time.Time          `bson:"joinedAt"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}
type Project struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description *string            `bson:"description,omitempty"`
	Emoji       string             `bson:"emoji"`
	WorkspaceID primitive.ObjectID `bson:"workspace"`
	CreatedBy   primitive.ObjectID `bson:"createdBy"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}
type Role struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty"`
	Name        enums.RoleType         `bson:"name"`
	Permissions []enums.PermissionType `bson:"permissions"`
	CreatedAt   time.Time              `bson:"createdAt"`
	UpdatedAt   time.Time              `bson:"updatedAt"`
}
type Task struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty"`
	TaskCode    string                 `bson:"taskCode"`
	Title       string                 `bson:"title"`
	Description *string                `bson:"description,omitempty"`
	ProjectID   primitive.ObjectID     `bson:"project"`
	WorkspaceID primitive.ObjectID     `bson:"workspace"`
	Status      enums.TaskSTatusEnum   `bson:"status"`
	Priority    enums.TaskPriorityEnum `bson:"priority"`
	AssignedTo  *primitive.ObjectID    `bson:"assignedTo,omitempty"`
	CreatedBy   primitive.ObjectID     `bson:"createdBy"`
	DueDate     *time.Time             `bson:"dueDate,omitempty"`
	CreatedAt   time.Time              `bson:"createdAt"`
	UpdatedAt   time.Time              `bson:"updatedAt"`
}
type User struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty"`
	Name             *string             `bson:"name,omitempty"`
	Email            string              `bson:"email"`
	Password         *string             `bson:"password,omitempty"`
	ProfilePicture   *string             `bson:"profilePicture,omitempty"`
	CurrentWorkspace *primitive.ObjectID `bson:"currentWorkspace,omitempty"`
	IsActive         bool                `bson:"isActive"`
	LastLogin        *time.Time          `bson:"lastLogin,omitempty"`
	CreatedAt        time.Time           `bson:"createdAt"`
	UpdatedAt        time.Time           `bson:"updatedAt"`
}
type Workspace struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description *string            `bson:"description,omitempty"`
	OwnerID     primitive.ObjectID `bson:"owner"`
	InviteCode  string             `bson:"inviteCode"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

// heper Function
// omiting the Password and returning a user witjut the password field
func (u *User) OmitPassword() *User {
	return &User{
		ID:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		ProfilePicture:   u.ProfilePicture,
		CurrentWorkspace: u.CurrentWorkspace,
		IsActive:         u.IsActive,
		LastLogin:        u.LastLogin,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}

// hash password
func (u *User) HashPassword() error {
	if u.Password == nil {
		return nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedStr := string(hashed)
	u.Password = &hashedStr
	return nil
}

// compasre hash password
func (u *User) ComparePassword(password string) error {
	if u.Password == nil {
		return errors.New("no password set")
	}
	return bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
}

// BeforeInsert is a MongoDB hook to hash password before insert
func (u *User) BeforeInsert(ctx context.Context) error {
	return u.HashPassword()
}

// BeforeUpdate is a MongoDB hook to hash password if it's being updated
func (u *User) BeforeUpdate(ctx context.Context) error {
	if u.Password != nil {
		return u.HashPassword()
	}
	return nil
}

// ResetInvideCode geberates a new invide cofe dor the worksapce
func (w *Workspace) ResetInviteCode(generateCode func() string) {
	w.InviteCode = generateCode()
}
