package event

type Event string

const (
	MediaUploaded    Event = "MediaUploaded"
	CreateThumbnail  Event = "CreateThumbnail"
	ThumbnailCreated Event = "ThumbnailCreated"
	RemoveObject     Event = "RemoveObject"
)
