package engine

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/config"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/statemachine"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/store"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"time"
)

type Engine struct {
	stateStore   *store.Store
	currentState *state.State
	cfg          *Config
	cfgFilename  string
	gameConfig   config.Game
	geometry     config.Geometry
	stageTimes   map[state.Stage]time.Duration
	queue        chan statemachine.Change
	quit         chan int
	hooks        []chan statemachine.StateChange
	rand         *rand.Rand
}

func NewEngine(gameConfig config.Game, seed int64, storeFilename, cfgFilename string) (s *Engine) {
	s = new(Engine)
	s.stateStore = store.NewStore(storeFilename)
	s.currentState = &state.State{}
	s.cfg = DefaultConfig()
	s.cfgFilename = cfgFilename
	s.cfg.Division = gameConfig.DefaultDivision
	s.gameConfig = gameConfig
	s.geometry = *gameConfig.DefaultGeometry[s.cfg.Division]
	s.stageTimes = loadStageTimes(gameConfig)
	s.queue = make(chan statemachine.Change, 100)
	s.quit = make(chan int)
	s.hooks = []chan statemachine.StateChange{}
	s.rand = rand.New(rand.NewSource(seed))
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
	e.cfg.LoadFrom(e.cfgFilename)
	go e.processChanges()
	return nil
}

func (e *Engine) Stop() {
	e.quit <- 0
}

func (e *Engine) LatestStateInStore() *state.State {
	entry := e.stateStore.LatestEntry()
	if entry != nil {
		return entry.State
	}
	return nil
}

func (e *Engine) processChanges() {
	for {
		select {
		case <-e.quit:
			if err := e.cfg.SaveTo(e.cfgFilename); err != nil {
				log.Printf("Could not save engine config to %v: %v", e.cfgFilename, err)
			}
			return
		case change := <-e.queue:
			newState, newChanges := statemachine.Process(e.currentState, change)
			entry := statemachine.StateChange{
				State:  newState,
				Change: change,
			}
			e.currentState = newState

			for _, newChange := range newChanges {
				e.queue <- newChange
			}

			// do not save state for ticks
			if change.ChangeType != statemachine.ChangeTypeTick {
				if err := e.stateStore.Add(store.Entry(entry)); err != nil {
					log.Println("Could not add new state to store: ", err)
				}
			}
			for _, hook := range e.hooks {
				hook <- entry
			}
		}
	}
}

// currentState gets the current state or returns an empty default state
func (e *Engine) initialStateFromStore() *state.State {
	latestEntry := e.stateStore.LatestEntry()
	var currentState *state.State
	if latestEntry == nil {
		currentState = &state.State{}
	} else {
		currentState = latestEntry.State.DeepCopy()
	}
	return currentState
}