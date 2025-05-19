package component

type UpdateChannel struct {
	Topic   string
	Channel <-chan string
}

type CommandChannel struct {
	Topic   string
	Channel chan<- string
}

type Updatable interface {
	UpdateChannels() []UpdateChannel
}

type Commandable interface {
	CommandChannels() []CommandChannel
}
