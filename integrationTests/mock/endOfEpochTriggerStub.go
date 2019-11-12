package mock

import "github.com/ElrondNetwork/elrond-go/data"

type EndOfEpochTriggerStub struct {
	ForceEndOfEpochCalled func(round int64) error
	IsEndOfEpochCalled    func() bool
	EpochCalled           func() uint32
	ReceivedHeaderCalled  func(handler data.HeaderHandler)
	UpdateCalled          func(round int64)
	ProcessedCalled       func()
}

func (e *EndOfEpochTriggerStub) Update(round int64) {
	if e.UpdateCalled != nil {
		e.UpdateCalled(round)
	}
}

func (e *EndOfEpochTriggerStub) Processed() {
	if e.ProcessedCalled != nil {
		e.ProcessedCalled()
	}
}

func (e *EndOfEpochTriggerStub) ForceEndOfEpoch(round int64) error {
	if e.ForceEndOfEpochCalled != nil {
		return e.ForceEndOfEpochCalled(round)
	}
	return nil
}

func (e *EndOfEpochTriggerStub) IsEndOfEpoch() bool {
	if e.IsEndOfEpochCalled != nil {
		return e.IsEndOfEpochCalled()
	}
	return false
}

func (e *EndOfEpochTriggerStub) Epoch() uint32 {
	if e.EpochCalled != nil {
		return e.EpochCalled()
	}
	return 0
}

func (e *EndOfEpochTriggerStub) ReceivedHeader(header data.HeaderHandler) {
	if e.ReceivedHeaderCalled != nil {
		e.ReceivedHeaderCalled(header)
	}
}

func (e *EndOfEpochTriggerStub) IsInterfaceNil() bool {
	return e == nil
}