package statemachine

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/config"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
	"reflect"
	"testing"
)

func Test_Statemachine(t *testing.T) {

	gameConfig := config.DefaultControllerConfig().Game
	type args struct {
		currentState *state.State
		change       Change
	}
	tests := []struct {
		name         string
		args         args
		wantNewState *state.State
	}{
		{
			name: "Command",
			args: args{
				currentState: &state.State{},
				change: Change{
					ChangeType: ChangeTypeNewCommand,
					NewCommand: &NewCommand{
						Command:    state.CommandDirect,
						CommandFor: state.Team_BLUE,
					},
				}},
			wantNewState: &state.State{
				Command:                    state.CommandDirect,
				CommandFor:                 state.Team_BLUE,
				CurrentActionTimeRemaining: gameConfig.FreeKickTime[config.DivA],
			},
		},
		{
			name: "Stage",
			args: args{
				currentState: &state.State{
					Stage: state.StagePreGame,
				},
				change: Change{
					ChangeType: ChangeTypeChangeStage,
					ChangeStage: &ChangeStage{
						NewStage: state.StageFirstHalf,
					},
				},
			},
			wantNewState: &state.State{
				Stage: state.StageFirstHalf,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewStateMachine(gameConfig, "/tmp/foo")
			if gotNewState, _ := sm.Process(tt.args.currentState, tt.args.change); !reflect.DeepEqual(gotNewState, tt.wantNewState) {
				t.Errorf("Process() != want:\n%v\n%v", gotNewState, tt.wantNewState)
			}
		})
	}
}
