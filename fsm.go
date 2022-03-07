package easyfsm

import (
	"github.com/wuqinqiang/easyfsm/log"
)

// FSM is the finite state machine
type FSM struct {
	// current state
	state State
	// current business
	businessName BusinessName
}

// NewFSM creates a new fsm
func NewFSM(businessName BusinessName, initState State) (fsm *FSM) {
	fsm = new(FSM)
	fsm.state = initState
	fsm.businessName = businessName
	return
}

// Call call the state's event func
func (f *FSM) Call(eventName EventName, opts ...ParamOption) (State, error) {
	businessMap, ok := stateMachineMap[f.businessName]
	if !ok {
		return f.getState(), UnKnownBusinessError{businessName: f.businessName}
	}
	// TODO 获取当前state转移到下一state的所有event处理器
	events, ok := businessMap[f.getState()]
	if !ok || events == nil {
		// TODO 当前状态下没有任何状态转移Handler，报错
		return f.getState(), UnKnownStateError{businessName: f.businessName, state: f.getState()}
	}

	opt := new(Param)
	for _, fn := range opts {
		fn(opt)
	}

	// TODO 获取当前state下处理eventName对应的事件的状态转移Handler
	eventEntity, ok := events[eventName]
	if !ok || eventEntity == nil {
		// TODO 当前状态下没有根据eventName对应的事件转移的下一状态的状态转移Handler, 报错
		return f.getState(), UnKnownEventError{businessName: f.businessName, state: f.getState(), event: eventName}
	}

	// call eventName func
	// TODO 调用状态转移Handler执行业务逻辑并将状态转移到下一状态
	state, err := eventEntity.Execute(opt)
	if err != nil {
		return f.getState(), err
	}
	oldState := f.getState()
	// TODO 更新状态
	f.setState(state)
	log.DefaultLogger.Log(log.LevelInfo, "eventName:", eventName,
		"beforeState:", oldState, "afterState:", f.getState())
	return f.getState(), nil
}

// getState get the state
func (f *FSM) getState() State {
	return f.state
}

// setState set the state
func (f *FSM) setState(newState State) {
	f.state = newState
}
