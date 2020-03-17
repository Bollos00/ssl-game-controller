package engine

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/config"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/statemachine"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_Engine(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal("Could not create a temporary directory: ", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("Could not cleanup tmpDir %v: %v", tmpDir, err)
		}
	}()

	gameConfig := config.DefaultControllerConfig().Game
	sm := statemachine.NewStateMachine(gameConfig, "/tmp/foo")
	engine := NewEngine(sm, 0, tmpDir+"/store.json.stream")
	hook := make(chan statemachine.StateChange)
	engine.RegisterHook(hook)
	if err := engine.Start(); err != nil {
		t.Fatal("Could not start engine")
	}
	engine.Enqueue(statemachine.Change{
		ChangeType: statemachine.ChangeTypeNewCommand,
		NewCommand: &statemachine.NewCommand{
			Command:    state.CommandDirect,
			CommandFor: state.Team_BLUE,
		},
	})
	gameEventTypeGoalLine := state.GameEventType_BALL_LEFT_FIELD_GOAL_LINE
	byTeam := state.Team_YELLOW
	engine.Enqueue(statemachine.Change{
		ChangeType: statemachine.ChangeTypeAddGameEvent,
		AddGameEvent: &statemachine.AddGameEvent{
			GameEvent: state.GameEvent{
				Type: &gameEventTypeGoalLine,
				Event: &state.GameEvent_BallLeftFieldGoalLine{
					BallLeftFieldGoalLine: &state.GameEvent_BallLeftField{
						ByTeam: &byTeam,
					},
				},
			},
		},
	})
	// wait for the changes to be processed
	<-hook
	<-hook

	engine.Stop()

	wantNewState := &state.State{
		Command:                    state.CommandDirect,
		CommandFor:                 state.Team_BLUE,
		CurrentActionTimeRemaining: gameConfig.FreeKickTime[config.DivA],
		NextCommand:                state.CommandDirect,
		NextCommandFor:             byTeam.Opposite(),
	}

	if gotNewState := engine.LatestStateInStore(); !reflect.DeepEqual(gotNewState, wantNewState) {
		t.Errorf("State mismatch:\n%v\n%v", gotNewState, wantNewState)
	}
}
