package model

const (
	TaskProductUpdateThumbnail = "product:updateThumbnail"
)

type TaskUpdateThumbnailPayload struct {
	OldObjectID string `json:"oldObjectID"`
	NewObjectID string `json:"newObjectID"`
}
