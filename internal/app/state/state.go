package state

import (
	"encoding/json"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/config"
	"time"
)

// State of the game
type State struct {
	Stage                      Stage                               `json:"stage" yaml:"stage"`
	Command                    RefCommand                          `json:"command" yaml:"command"`
	CommandFor                 Team                                `json:"commandForTeam" yaml:"commandForTeam"`
	StageTimeElapsed           time.Duration                       `json:"stageTimeElapsed" yaml:"stageTimeElapsed"`
	StageTimeLeft              time.Duration                       `json:"stageTimeLeft" yaml:"stageTimeLeft"`
	MatchTimeStart             time.Time                           `json:"matchTimeStart" yaml:"matchTimeStart"`
	MatchDuration              time.Duration                       `json:"matchDuration" yaml:"matchDuration"` // MatchDuration contains the updated match duration based on MatchTimeStart for the UI
	TeamState                  map[Team]*TeamInfo                  `json:"teamState" yaml:"teamState"`
	PlacementPos               *Location                           `json:"placementPos" yaml:"placementPos"`
	NextCommand                RefCommand                          `json:"nextCommand" yaml:"nextCommand"`
	NextCommandFor             Team                                `json:"nextCommandFor" yaml:"nextCommandFor"`
	PrevCommands               []RefCommand                        `json:"prevCommands" yaml:"prevCommands"`
	PrevCommandsFor            []Team                              `json:"prevCommandsFor" yaml:"prevCommandsFor"`
	CurrentActionTimeRemaining time.Duration                       `json:"currentActionTimeRemaining" yaml:"currentActionTimeRemaining"` // CurrentActionTimeRemaining contains the updated remaining lack of progress time for the UI
	GameEvents                 []GameEvent                         `json:"gameEvents" yaml:"gameEvents"`
	ProposedGameEvents         []ProposedGameEvent                 `json:"gameEventProposals" yaml:"gameEventProposals"`
	Division                   config.Division                     `json:"division" yaml:"division"`
	AutoContinue               bool                                `json:"autoContinue" yaml:"autoContinue"`
	FirstKickoffTeam           Team                                `json:"firstKickoffTeam" yaml:"firstKickoffTeam"`
	GameEventBehavior          map[GameEventType]GameEventBehavior `json:"gameEventBehavior" yaml:"gameEventBehavior"`
}

// NewState creates a new state, initialized for the start of a new game
func NewState() (s State) {
	s.Stage = StagePreGame
	s.Command = CommandHalt

	s.StageTimeLeft = 0
	s.StageTimeElapsed = 0
	s.MatchDuration = 0
	s.MatchTimeStart = time.Unix(0, 0)

	s.TeamState = map[Team]*TeamInfo{}
	s.TeamState[Team_YELLOW] = newTeamInfo()
	s.TeamState[Team_BLUE] = newTeamInfo()
	s.TeamState[Team_BLUE].OnPositiveHalf = !s.TeamState[Team_YELLOW].OnPositiveHalf

	s.GameEvents = []GameEvent{}
	s.ProposedGameEvents = []ProposedGameEvent{}

	s.Division = config.DivA
	s.FirstKickoffTeam = Team_YELLOW
	s.AutoContinue = true

	s.GameEventBehavior = map[GameEventType]GameEventBehavior{}
	for _, event := range AllGameEvents() {
		s.GameEventBehavior[event] = GameEventBehaviorOn
	}

	return
}

func (s State) GameState() GameState {
	return s.Command.GameState()
}

func (s State) String() string {
	bytes, e := json.Marshal(s)
	if e != nil {
		return e.Error()
	}
	return string(bytes)
}
