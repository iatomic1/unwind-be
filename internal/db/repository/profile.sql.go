// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: profile.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const getProfileById = `-- name: GetProfileById :one
SELECT id, user_id, profile_pic, name, username, cover_pic, created_at, updated_at FROM profile WHERE id = $1
`

func (q *Queries) GetProfileById(ctx context.Context, id uuid.UUID) (*Profile, error) {
	row := q.db.QueryRow(ctx, getProfileById, id)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProfilePic,
		&i.Name,
		&i.Username,
		&i.CoverPic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getProfileByUserId = `-- name: GetProfileByUserId :one
SELECT id, user_id, profile_pic, name, username, cover_pic, created_at, updated_at FROM profile WHERE user_id = $1
`

func (q *Queries) GetProfileByUserId(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	row := q.db.QueryRow(ctx, getProfileByUserId, userID)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProfilePic,
		&i.Name,
		&i.Username,
		&i.CoverPic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const insertProfile = `-- name: InsertProfile :one
INSERT INTO profile (
  id,
  username,
  user_id,
  profile_pic,
  name,
  cover_pic
) VALUES (
  uuid_generate_v4(),
  $1,               -- username (required)
  $2,               -- user_id (required)
  $3,               -- profile_pic (optional)
  $4,               -- name (optional)
  $5                -- cover_pic (optional)
)
RETURNING id, user_id, profile_pic, name, username, cover_pic, created_at, updated_at
`

type InsertProfileParams struct {
	Username   string    `binding:"required,min=8" example:"Slimmm Shaddy" json:"username"`
	UserID     uuid.UUID `binding:"required,uuid" json:"userId"`
	ProfilePic *string   `json:"profilePic"`
	Name       *string   `json:"name"`
	CoverPic   *string   `json:"coverPic"`
}

func (q *Queries) InsertProfile(ctx context.Context, arg InsertProfileParams) (*Profile, error) {
	row := q.db.QueryRow(ctx, insertProfile,
		arg.Username,
		arg.UserID,
		arg.ProfilePic,
		arg.Name,
		arg.CoverPic,
	)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProfilePic,
		&i.Name,
		&i.Username,
		&i.CoverPic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const updateProfile = `-- name: UpdateProfile :one
UPDATE profile
SET 
    username = COALESCE($1, username),
    name = COALESCE($2, name),
    cover_pic = COALESCE($3, cover_pic),
    profile_pic = COALESCE($4, profile_pic),
    updated_at = now()
WHERE id = $5
RETURNING id, user_id, profile_pic, name, username, cover_pic, created_at, updated_at
`

type UpdateProfileParams struct {
	Username   string    `binding:"required,min=8" example:"Slimmm Shaddy" json:"username"`
	Name       *string   `json:"name"`
	CoverPic   *string   `json:"coverPic"`
	ProfilePic *string   `json:"profilePic"`
	ID         uuid.UUID `json:"id"`
}

func (q *Queries) UpdateProfile(ctx context.Context, arg UpdateProfileParams) (*Profile, error) {
	row := q.db.QueryRow(ctx, updateProfile,
		arg.Username,
		arg.Name,
		arg.CoverPic,
		arg.ProfilePic,
		arg.ID,
	)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProfilePic,
		&i.Name,
		&i.Username,
		&i.CoverPic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
