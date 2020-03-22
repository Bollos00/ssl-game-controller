package statemachine

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/config"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
	"github.com/go-test/deep"
	"testing"
)

func Test_Statemachine(t *testing.T) {

	gameConfig := config.DefaultControllerConfig().Game
	type args struct {
		initState func(*state.State)
		change    Change
	}
	tests := []struct {
		name         string
		args         args
		wantNewState func(*state.State)
	}{
		{
			name: "Command",
			args: args{
				initState: func(s *state.State) {
				},
				change: Change{
					ChangeType: ChangeTypeNewCommand,
					NewCommand: &NewCommand{
						Command:    state.CommandDirect,
						CommandFor: state.Team_BLUE,
					},
				}},
			wantNewState: func(s *state.State) {
				s.Command = state.CommandDirect
				s.CommandFor = state.Team_BLUE
				s.CurrentActionTimeRemaining = gameConfig.FreeKickTime[config.DivA]
			},
		},
		{
			name: "Stage",
			args: args{
				initState: func(s *state.State) {
					s.Stage = state.Referee_NORMAL_FIRST_HALF_PRE
				},
				change: Change{
					ChangeType: ChangeTypeChangeStage,
					ChangeStage: &ChangeStage{
						NewStage: state.Referee_NORMAL_FIRST_HALF,
					},
				},
			},
			wantNewState: func(s *state.State) {
				s.Stage = state.Referee_NORMAL_FIRST_HALF
				s.StageTimeLeft = gameConfig.Normal.HalfDuration
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewStateMachine(gameConfig)
			currentState := state.NewState()
			tt.args.initState(&currentState)
			newState := state.NewState()
			tt.wantNewState(&newState)

			gotNewState, _ := sm.Process(currentState, tt.args.change)
			diffs := deep.Equal(gotNewState, newState)
			if len(diffs) > 0 {
				t.Error("States differ: ", diffs)
			}
		})
	}
}
