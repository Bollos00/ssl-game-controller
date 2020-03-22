package engine

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/config"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/statemachine"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/store"
	"github.com/RoboCup-SSL/ssl-game-controller/pkg/timer"
	"github.com/pkg/errors"
	"log"
	"time"
)

const changeOriginEngine = "Engine"

type Engine struct {
	gameConfig      config.Game
	stateStore      *store.Store
	currentState    state.State
	stateMachine    *statemachine.StateMachine
	queue           chan statemachine.Change
	quit            chan int
	hooks           []chan statemachine.StateChange
	timeProvider    timer.TimeProvider
	lastTimeUpdate  time.Time
	readyToContinue *bool
	ballPlaced      *bool
}

func NewEngine(gameConfig config.Game) (s *Engine) {
	s = new(Engine)
	s.stateStore = store.NewStore(gameConfig.StateStoreFile)
	s.stateMachine = statemachine.NewStateMachine(gameConfig)
	s.queue = make(chan statemachine.Change, 100)
	s.quit = make(chan int)
	s.hooks = []chan statemachine.StateChange{}
	s.timeProvider = func() time.Time { return time.Now() }
	s.lastTimeUpdate = s.timeProvider()
	return
}

func (e *Engine) Enqueue(change statemachine.Change) {
	e.queue <- change
}

func (e *Engine) RegisterHook(hook chan statemachine.StateChange) {
	e.hooks = append(e.hooks, hook)
}

func (e *Engine) UnregisterHook(hook chan statemachine.StateChange) bool {
	for i, h := range e.hooks {
		if h == hook {
			e.hooks = append(e.hooks[:i], e.hooks[i+1:]...)
			select {
			case <-hook:
			case <-time.After(10 * time.Millisecond):
			}
			return true
		}
	}
	return false
}

func (e *Engine) Start() error {
	if err := e.stateStore.Open(); err != nil {
		return errors.Wrap(err, "Could not open state store")
	}
	if err := e.stateStore.Load(); err != nil {
		return errors.Wrap(err, "Could not load state store")
	}
	e.currentState = e.initialStateFromStore()
	e.stateMachine.UpdateGeometry(e.gameConfig.DefaultGeometry[e.currentState.Division])
	go e.processChanges()
	return nil
}

func (e *Engine) Stop() {
	e.quit <- 0
}

func (e *Engine) CurrentState() state.State {
	return e.currentState
}

func (e *Engine) processChanges() {
	for {
		select {
		case <-e.quit:
			return
		case change := <-e.queue:
			entry := statemachine.StateChange{}
			entry.Change = change
			var newChanges []statemachine.Change
			entry.State, newChanges = e.stateMachine.Process(e.currentState, change)
			e.currentState = entry.State

			e.postProcessChange(change)

			for _, newChange := range newChanges {
				e.queue <- newChange
			}

			// do not save state for ticks
			if err := e.stateStore.Add(store.Entry(entry)); err != nil {
				log.Println("Could not add new state to store: ", err)
			}
			for _, hook := range e.hooks {
				hook <- entry
			}
		case <-time.After(10 * time.Millisecond):
			e.Tick()
		}
	}
}

// currentState gets the current state or returns an empty default state
func (e *Engine) initialStateFromStore() state.State {
	latestEntry := e.stateStore.LatestEntry()
	if latestEntry == nil {
		return state.NewState()
	}
	return latestEntry.State
}

func (e *Engine) postProcessChange(change statemachine.Change) {
	switch change.ChangeType {
	case statemachine.ChangeTypeChangeStage:
		if change.ChangeStage.NewStage == state.StageFirstHalf {
			e.currentState.MatchTimeStart = e.timeProvider()
		}
	}
}
