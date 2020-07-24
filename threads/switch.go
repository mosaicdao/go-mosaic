package threads

import (
	"fmt"

	"github.com/mosaicdao/go-mosaic/libs/service"
)

// Switch defines a switch struct.
type Switch struct {
	service.BaseService

	topics            []*Topic
	reactors          map[string]Reactor
	reactorsByTopicID map[TopicID]Reactor
}

func (sw *Switch) OnStart() error {
	// starts reactors
	for _, reactor := range sw.reactors {
		err := reactor.Start()
		if err != nil {
			return fmt.Errorf("failed to start %v: %w", reactor, err)
		}
	}

	return nil
}

func (sw *Switch) OnStop() {
	// stops reactors
	for _, reactor := range sw.reactors {
		reactor.Stop()
	}
}

// AddReactor adds a reactor to the switch.
// The function updates a mapping from a topic id to a reactor based on the reactor's topics.
// The function updates a mapping from the given reactor name to the reactor.
// The function requires that no two reactors can share the same topic.
// The function sets the current object as a switch to the given reactor and returns it.
func (sw *Switch) AddReactor(name string, reactor Reactor) Reactor {
	for _, topic := range reactor.GetTopics() {
		topicID := topic.ID

		// No two reactors can share the same topic.
		if sw.reactorsByTopicID[topicID] != nil {
			panic(
				fmt.Sprintf(
					"There is already a reactor (%v) registered for the topic %X",
					sw.reactorsByTopicID[topicID],
					topicID,
				),
			)
		}

		sw.topics = append(sw.topics, topic)
		sw.reactorsByTopicID[topicID] = reactor
	}

	sw.reactors[name] = reactor
	reactor.SetSwitch(sw)

	return reactor
}

func (sw *Switch) RemoveReactor(name string, reactor Reactor) {

}

// func (sw *Switch) Reactors() map[string]Reactor
// func (sw *Switch) Reactor(name string) Reactor
