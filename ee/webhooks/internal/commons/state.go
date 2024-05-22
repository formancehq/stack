package commons

import(
	"github.com/formancehq/stack/libs/go-libs/sync/shared"	
	"time"

)


type State struct {

	HooksById 			*MapSharedHook
	ActiveHooksByEvent 	*MapSharedHooks

	AttemptsById 	*MapSharedAttempt
	WaitingAttempts *SharedAttempts


	ToSaveAttempts *SharedAttempts

}

func (s *State) RoutineEvent(stopChan chan struct{}, eventChan chan Event, handleEvent func(e Event)){
	for{
		select {
		case <- stopChan:
			return 
		case ev := <-eventChan:
			handleEvent(ev)
		}
	}
}

func (s *State) RoutineTime(stopChan chan struct{}, ticker *time.Ticker, handleTime func()){
	for {
		select {
		case <- stopChan:
			ticker.Stop()
			return 
		case <-ticker.C:
			handleTime()
		}	 
	}
}

func (s *State) LoadHooks(hooks *[]*Hook){
	var sHooks *SharedHooks = (&SharedHooks{}).From(hooks)
	for _, sH := range *sHooks.Val {
		s.HooksById.Add(sH.Val.ID, sH)
		if(sH.Val.IsActive()){
			s.ActiveHooksByEvent.Adds(sH.Val.Events, sH)
		}
		
	}
}

func(s *State) AddNewHook(hook *Hook){
	sHook := shared.NewShared(hook)
	s.HooksById.Add(sHook.Val.ID, &sHook)
	if(sHook.Val.IsActive()){
		s.ActiveHooksByEvent.Adds(sHook.Val.Events, &sHook)
	}
}

func (s *State) DeleteHook(id string) *SharedHook{
	sHook := s.HooksById.Get(id)
	if(sHook == nil){return nil}

	sHook.WLock()
	sHook.Val.Delete()
	sHook.WUnlock()

	s.ActiveHooksByEvent.Removes(sHook.Val.Events, sHook)
	s.HooksById.Remove(sHook.Val.ID)

	
	return sHook
}

func (s *State) ActivateHook(id string) *SharedHook{
	sHook := s.HooksById.Get(id)
	if(sHook == nil) {return nil}
	
	sHook.WLock()
	sHook.Val.Enable()
	sHook.WUnlock()

	s.ActiveHooksByEvent.Adds(sHook.Val.Events, sHook)

	return sHook
}

func (s *State) DisableHook(id string) *SharedHook{
	sHook := s.HooksById.Get(id)
	if(sHook == nil) {return nil}

	sHook.WLock()
	sHook.Val.Disable()
	sHook.WUnlock()

	s.ActiveHooksByEvent.Removes(sHook.Val.Events, sHook)

	return sHook
}

func (s *State) UpdateHookEndpoint(id string, endpoint string) *SharedHook{
	sHook := s.HooksById.Get(id)
	if(sHook == nil) {return nil}

	sHook.WLock()
	sHook.Val.Endpoint = endpoint 
	sHook.WUnlock()

	return sHook
}

func (s *State) UpdateHookSecret(id string, secret string) *SharedHook{
	sHook := s.HooksById.Get(id)
	if(sHook == nil) {return nil}

	sHook.WLock()
	sHook.Val.Secret = secret 
	sHook.WUnlock()

	return sHook
}

func (s *State) UpdateHookRetry(id string, retry bool) *SharedHook {
	sHook := s.HooksById.Get(id)
	if(sHook == nil) {return nil}

	sHook.WLock()
	sHook.Val.Retry = retry 
	sHook.WUnlock()

	return sHook
}

func (s *State) LoadWaitingAttempts(attempts *[]*Attempt){
	var sAttempts *SharedAttempts = (&SharedAttempts{}).From(attempts)
	for _, sA := range *sAttempts.Val {
		s.WaitingAttempts.Add(sA)
		s.AttemptsById.Add(sA.Val.ID, sA)
	}
}

func (s *State) AddNewAttempt(attempt *Attempt){
	sA := shared.NewShared(attempt)

	s.WaitingAttempts.Add(&sA)
	s.AttemptsById.Add(sA.Val.ID, &sA)
}

func (s *State) FlushAttempt(id string) *SharedAttempt{
	sAttempt := s.AttemptsById.Get(id)
	if(sAttempt == nil) {return nil}
	
	sAttempt.WLock()
	sAttempt.Val.NextTry = time.Now()
	sAttempt.WUnlock()

	return sAttempt
}

func (s *State) FlushAttempts() {
	for _, sA := range *s.WaitingAttempts.Val{
		sA.WLock()
		sA.Val.NextTry = time.Now()
		sA.WUnlock()
	}
}

func (s *State) AbortAttempt(id string) *SharedAttempt{
	sAttempt := s.AttemptsById.Get(id)
	if(sAttempt == nil) {return nil}

	s.AttemptsById.Remove(sAttempt.Val.ID)
	s.WaitingAttempts.Remove(sAttempt)

	sAttempt.WLock()
	sAttempt.Val.Status = AbortStatus
	sAttempt.WUnlock()

	return sAttempt

}

func NewState() *State {
	return &State{
		HooksById: NewMapSharedHook(),
		ActiveHooksByEvent: NewMapSharedHooks(),
		AttemptsById: NewMapSharedAttempt(),
		WaitingAttempts: NewSharedAttempts(),
		ToSaveAttempts: NewSharedAttempts(),
	}
}