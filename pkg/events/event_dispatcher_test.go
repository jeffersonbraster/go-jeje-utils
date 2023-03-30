package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type TestEventHandle struct {
	ID int
}

func (h *TestEventHandle) Handle(event EventInterface) {}

type EventDispatcherTestSuite struct {
	suite.Suite
	event TestEvent
	event2 TestEvent
	handler TestEventHandle
	handler2 TestEventHandle
	handle3 TestEventHandle
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandle{
		ID: 1,
	}
	suite.handler2 = TestEventHandle{
		ID: 2,
	}
	suite.handle3 = TestEventHandle{
		ID: 3,
	}
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
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	suite.Nil(err)

	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.handlers[suite.event.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Equal(ErrHandleAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	//Event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	//Event 2
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handle3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	//Event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handle3))
}

type MockHandle struct {
	mock.Mock
}

func (m *MockHandle) Handle(event EventInterface) {
	m.Called(event)
}

func(suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eh := &MockHandle{}
	eh.On("Handle", &suite.event)
	suite.eventDispatcher.Register(suite.event.GetName(), eh)
	suite.eventDispatcher.Dispatch(&suite.event)
	eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	//Event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	//Event 2
	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][0])

	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler2)
	assert.Equal(suite.T(), 0, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
}


func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}