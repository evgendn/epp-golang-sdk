package eppoclient

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CallbackMock struct {
	mock.Mock
}

func (m *CallbackMock) CallbackFn() {
	m.Called()
}

func Test_PollerPoll_InvokesCallbackUntilStoped(t *testing.T) {
	expected := 5

	callbackMock := CallbackMock{}
	callbackMock.On("CallbackFn").Return()

	poller := newPoller(10*time.Millisecond, callbackMock.CallbackFn)
	poller.Start()
	time.Sleep(55 * time.Millisecond)
	poller.Stop()

	callbackMock.AssertNumberOfCalls(t, "CallbackFn", expected)
}

func Test_PollerPoll_StopsOnError(t *testing.T) {
	var callCount int32
	var expected int32 = 3

	poller := newPoller(10*time.Millisecond, func() {
		if atomic.AddInt32(&callCount, 1) == 3 {
			panic("some_error")
		}
	})
	poller.Start()

	time.Sleep(55 * time.Millisecond)
	assert.Equal(t, expected, atomic.LoadInt32(&callCount))
}

func Test_PollerPoll_ManualStop(t *testing.T) {
	expected := 3

	callbackMock := CallbackMock{}
	callbackMock.On("CallbackFn").Return()

	var poller = newPoller(10*time.Millisecond, callbackMock.CallbackFn)
	poller.Start()

	time.Sleep(35 * time.Millisecond)

	poller.Stop()

	time.Sleep(20 * time.Millisecond)
	callbackMock.AssertNumberOfCalls(t, "CallbackFn", expected)

	poller.Stop()
}
