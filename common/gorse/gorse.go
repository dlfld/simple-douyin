package gorse

var Client *GorseClient

func init() {
	Client = NewGorseClient("http://gorse:8088", "")
}
