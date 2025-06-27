package main

import "fmt"

type customEvent struct {
	Name string
	Data any
}

type EventBusService struct {
	eventChannel chan customEvent
}

func NewEventBusService(eventChannel chan customEvent) *EventBusService {
	return &EventBusService{
		eventChannel: eventChannel,
	}
}

func (e *EventBusService) EmitEvent(name string, data any) {
	fmt.Println("Emitting event:", name, data)
	e.eventChannel <- customEvent{
		Name: name,
		Data: data,
	}
}
