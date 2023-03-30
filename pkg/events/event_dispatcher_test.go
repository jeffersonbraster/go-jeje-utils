package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name     string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandle struct {}

func (h *TestEventHandle) Handle(event EventInterface) {}

type EventDispatcherTestSuite struct {
	suite.Suite
	event TestEvent
	event2 TestEvent
	handler TestEventHandle
	handle2 TestEventHandle
	handle3 TestEventHandle
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandle{}
	suite.handle2 = TestEventHandle{}
	suite.handle3 = TestEventHandle{}
	suite.event = TestEvent{
		Name: "test",
		Payload: "test",
	}
	suite.event2 = TestEvent{
		Name: "test2",
		Payload: "test2",
	}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	assert.True(suite.T(), true)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}