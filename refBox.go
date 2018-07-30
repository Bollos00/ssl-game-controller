package main

import "time"

type RefBox struct {
	State           *RefBoxState
	timer           Timer
	MatchTimeStart  time.Time
	newEventChannel chan RefBoxEvent
}

func NewRefBox() (refBox *RefBox) {
	refBox = new(RefBox)
	refBox.State = NewRefBoxState()
	refBox.timer = NewTimer()
	refBox.newEventChannel = make(chan RefBoxEvent)
	refBox.MatchTimeStart = time.Unix(0, 0)

	return
}

func (r *RefBox) Run() (err error) {
	r.timer.Start()

	go func() {
		for {
			r.timer.WaitTillNextFullSecond()
			r.Tick()
			r.newEventChannel <- RefBoxEvent{}
		}
	}()
	return nil
}

func (r *RefBox) Tick() {
	delta := r.timer.Delta()
	updateTimes(r, delta)

	if r.MatchTimeStart.After(time.Unix(0, 0)) {
		r.State.MatchDuration = time.Now().Sub(r.MatchTimeStart)
	}
}

func updateTimes(r *RefBox, delta time.Duration) {
	if r.State.GameState == GameStateRunning {
		r.State.GameTimeElapsed += delta
		r.State.GameTimeLeft -= delta

		for _, teamState := range r.State.TeamState {
			reduceYellowCardTimes(teamState, delta)
			removeElapsedYellowCards(teamState)
		}
	}

	if r.State.GameState == GameStateTimeout && r.State.GameStateFor != nil {
		r.State.TeamState[*r.State.GameStateFor].TimeoutTimeLeft -= delta
	}
}

func reduceYellowCardTimes(teamState *RefBoxTeamState, delta time.Duration) {
	for i := range teamState.YellowCardTimes {
		teamState.YellowCardTimes[i] -= delta
	}
}

func removeElapsedYellowCards(teamState *RefBoxTeamState) {
	b := teamState.YellowCardTimes[:0]
	for _, x := range teamState.YellowCardTimes {
		if x > 0 {
			b = append(b, x)
		}
	}
	teamState.YellowCardTimes = b
}
