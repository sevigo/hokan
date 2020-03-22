package core

type ServerSideEventCreater interface {
	PublishMessage(string)
}
