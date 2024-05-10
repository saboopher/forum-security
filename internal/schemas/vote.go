package schemas

import "github.com/gofrs/uuid"

type ShowVote struct {
	VoteID    uuid.UUID
	UserID    uuid.UUID
	PostID    uuid.UUID
	CommentID uuid.UUID
	Binary    int
}

type CreateVote struct {
	ShowVote
}
