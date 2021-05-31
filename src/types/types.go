package types

type StringStream struct {
	Read  chan string
	Write chan string
}

func StringStreamCreate() StringStream {
	return StringStream{
		Read:  make(chan string),
		Write: make(chan string)}
}

/*
type TriggerStream struct {
	Flag chan struct{}
}

func TriggerStreamCreate() TriggerStream {
	return TriggerStream {
		Flag: make(chan struct{}) }
}
*/
