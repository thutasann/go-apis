package advanced

import "fmt"

type Logger interface {
	Log(message string)
}

type CustomLogger struct{}

func (c CustomLogger) Log(mesasge string) {
	fmt.Println("Log: ", mesasge)
}

type Service struct {
	Logger
}

func InterfaceComposition() {
	svc := Service{Logger: CustomLogger{}}
	svc.Log("Service started")
}
