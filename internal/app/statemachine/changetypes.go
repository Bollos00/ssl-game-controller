package statemachine

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/config"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
)

type ChangeType string
type ChangeOrigin string

const (
	ChangeTypeNewCommand               ChangeType = "NewCommand"
	ChangeTypeAddGameEvent             ChangeType = "AddGameEvent"
	ChangeTypeAddProposedGameEvent     ChangeType = "AddProposedGameEvent"
	ChangeTypeChangeStage              ChangeType = "ChangeStage"
	ChangeTypeSetBallPlacementPos      ChangeType = "SetBallPlacementPos"
	ChangeTypeAddYellowCard            ChangeType = "AddYellowCard"
	ChangeTypeAddRedCard               ChangeType = "AddRedCard"
	ChangeTypeYellowCardOver           ChangeType = "ChangeTypeYellowCardOver"
	ChangeTypeUpdateStateMachineConfig ChangeType = "UpdateStateMachineConfig"
	ChangeTypeUpdateTeamState          ChangeType = "UpdateTeamState"
	ChangeTypeSwitchColor              ChangeType = "SwitchColor"
	ChangeTypeRevert                   ChangeType = "Revert"
)

type Change struct {
	ChangeType               ChangeType                `json:"changeType,omitempty" yaml:"changeType"`
	ChangeOrigin             ChangeOrigin              `json:"changeOrigin,omitempty" yaml:"changeOrigin"`
	NewCommand               *NewCommand               `json:"newCommand,omitempty" yaml:"newCommand"`
	AddGameEvent             *AddGameEvent             `json:"addGameEvent,omitempty" yaml:"addGameEvent"`
	AddProposedGameEvent     *AddProposedGameEvent     `json:"addProposedGameEvent,omitempty" yaml:"addProposedGameEvent"`
	ChangeStage              *ChangeStage              `json:"changeStage,omitempty" yaml:"changeStage"`
	SetBallPlacementPos      *SetBallPlacementPos      `json:"setBallPlacementPos,omitempty" yaml:"setBallPlacementPos"`
	AddYellowCard            *AddYellowCard            `json:"addYellowCard,omitempty" yaml:"addYellowCard"`
	AddRedCard               *AddRedCard               `json:"addRedCard,omitempty" yaml:"addRedCard"`
	UpdateStateMachineConfig *UpdateStateMachineConfig `json:"updateStateMachineConfig,omitempty" yaml:"updateStateMachineConfig"`
	UpdateTeamState          *UpdateTeamState          `json:"updateTeamState,omitempty" yaml:"updateTeamState"`
	Revert                   *Revert                   `json:"revert,omitempty" yaml:"revert"`
}

type NewCommand struct {
	Command    state.RefCommand `json:"command" yaml:"command"`
	CommandFor state.Team       `json:"commandFor" yaml:"commandFor"`
}

type SetBallPlacementPos struct {
	Pos state.Location `json:"pos" yaml:"pos"`
}

type AddGameEvent struct {
	GameEvent state.GameEvent `json:"gameEvent" yaml:"gameEvent"`
}

type AddProposedGameEvent struct {
	GameEvent state.ProposedGameEvent `json:"gameEvent" yaml:"gameEvent"`
}

type ChangeStage struct {
	NewStage state.Stage `json:"newStage" yaml:"newStage"`
}

type AddYellowCard struct {
	ForTeam           state.Team       `json:"forTeam" yaml:"forTeam"`
	CausedByGameEvent *state.GameEvent `json:"causedByGameEvent,omitempty" yaml:"causedByGameEvent"`
}

type AddRedCard struct {
	ForTeam           state.Team       `json:"forTeam" yaml:"forTeam"`
	CausedByGameEvent *state.GameEvent `json:"causedByGameEvent,omitempty" yaml:"causedByGameEvent"`
}

type UpdateStateMachineConfig struct {
	Division          *config.Division                          `json:"division,omitempty" yaml:"division"`
	FirstKickoffTeam  *state.Team                               `json:"firstKickoffTeam,omitempty" yaml:"firstKickoffTeam"`
	AutoContinue      *bool                                     `json:"autoContinue,omitempty" yaml:"autoContinue"`
	GameEventBehavior map[state.GameEventType]GameEventBehavior `json:"gameEventBehavior,omitempty" yaml:"gameEventBehavior"`
}

type UpdateTeamState struct {
	ForTeam state.Team `json:"forTeam"`

	Goals                 *int                     `json:"goals,omitempty" yaml:"goals"`
	Goalkeeper            *int                     `json:"goalkeeper,omitempty" yaml:"goalkeeper"`
	YellowCards           map[int]state.YellowCard `json:"yellowCards,omitempty" yaml:"yellowCards"`
	YellowCardsRemove     *int                     `json:"yellowCardsRemove,omitempty" yaml:"yellowCardsRemove"`
	RedCards              map[int]state.RedCard    `json:"redCards,omitempty" yaml:"redCards"`
	RedCardsRemove        *int                     `json:"yellowCardsRemove,omitempty" yaml:"yellowCardsRemove"`
	TimeoutsLeft          *int                     `json:"timeoutsLeft,omitempty" yaml:"timeoutsLeft"`
	TimeoutTimeLeft       *string                  `json:"timeoutTimeLeft,omitempty" yaml:"timeoutTimeLeft"`
	OnPositiveHalf        *bool                    `json:"onPositiveHalf,omitempty" yaml:"onPositiveHalf"`
	TeamName              *string                  `json:"teamName,omitempty" yaml:"teamName"`
	Fouls                 map[int]state.Foul       `json:"fouls,omitempty" yaml:"fouls"`
	FoulsRemove           *int                     `json:"foulsRemove,omitempty" yaml:"foulsRemove"`
	BallPlacementFailures *int                     `json:"ballPlacementFailures,omitempty" yaml:"ballPlacementFailures"`
	CanPlaceBall          *bool                    `json:"canPlaceBall,omitempty" yaml:"canPlaceBall"`
	BotSubstitutionIntend *bool                    `json:"botSubstitutionIntend,omitempty" yaml:"botSubstitutionIntend"`
}

type Revert struct {
	StateId int `json:"stateId" yaml:"stateId"`
}
