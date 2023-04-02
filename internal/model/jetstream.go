package model

type PublisherUsecase interface {
	CreateStream() error
}

type ConsumerUsecase interface {
	ConsumeEvent() error
}

type JSDeleteObjectPayload struct {
	ObjectID string `json:"objectID"`
}
