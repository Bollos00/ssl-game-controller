package publish

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/statemachine"
	"time"
)

type MessageGenerator struct {
	goal             map[state.Team]bool
	commandCounter   uint32
	commandTimestamp uint64
}

// NewPublisher creates a new publisher that publishes referee messages via UDP to the teams
func NewMessageGenerator() (m *MessageGenerator) {
	m = new(MessageGenerator)
	m.goal = map[state.Team]bool{}
	m.goal[state.Team_BLUE] = false
	m.goal[state.Team_YELLOW] = false

	return
}

// GenerateRefereeMessages generates a list of referee messages that result from the given state change
func (p *MessageGenerator) GenerateRefereeMessages(change statemachine.StateChange) (rs []*state.Referee) {
	// send the GOAL command based on the team score for compatibility with old behavior
	if change.Change.ChangeType == statemachine.ChangeTypeAddGameEvent &&
		*change.Change.AddGameEvent.GameEvent.Type == state.GameEventType_GOAL {
		p.updateCommand()
		refereeMsg := p.stateToRefereeMessage(change.State)
		if change.Change.AddGameEvent.GameEvent.ByTeam() == state.Team_YELLOW {
			*refereeMsg.Command = state.Referee_GOAL_YELLOW
		} else {
			*refereeMsg.Command = state.Referee_GOAL_BLUE
		}
		rs = append(rs, refereeMsg)
	}

	if change.Change.ChangeType == statemachine.ChangeTypeNewCommand {
		p.updateCommand()
	}

	rs = append(rs, p.stateToRefereeMessage(change.State))

	return
}

func (p *MessageGenerator) updateCommand() {
	p.commandCounter++
	p.commandTimestamp = uint64(time.Now().UnixNano() / 1000)
}

func (p *MessageGenerator) stateToRefereeMessage(matchState *state.State) (r *state.Referee) {
	r = newRefereeMessage()
	r.DesignatedPosition = mapLocation(matchState.PlacementPos)
	r.ProposedGameEvents = mapProposedGameEvents(matchState.ProposedGameEvents)

	*r.Command = mapCommand(matchState.Command, matchState.CommandFor)
	*r.CommandCounter = p.commandCounter
	*r.CommandTimestamp = p.commandTimestamp
	*r.PacketTimestamp = uint64(time.Now().UnixNano() / 1000)
	*r.Stage = mapStage(matchState.Stage)
	*r.StageTimeLeft = int32(matchState.StageTimeLeft.Nanoseconds() / 1000)
	*r.BlueTeamOnPositiveHalf = matchState.TeamState[state.Team_BLUE].OnPositiveHalf
	*r.NextCommand = mapCommand(matchState.NextCommand, matchState.NextCommandFor)
	*r.CurrentActionTimeRemaining = int32(matchState.CurrentActionTimeRemaining.Nanoseconds() / 1000)

	updateTeam(r.Yellow, matchState.TeamState[state.Team_YELLOW])
	updateTeam(r.Blue, matchState.TeamState[state.Team_BLUE])
	return
}

func updateTeam(teamInfo *state.Referee_TeamInfo, teamState *state.TeamInfo) {
	*teamInfo.Name = teamState.Name
	*teamInfo.Score = uint32(teamState.Goals)
	*teamInfo.RedCards = uint32(len(teamState.RedCards))
	teamInfo.YellowCardTimes = mapYellowCardTimes(teamState.YellowCards)
	*teamInfo.YellowCards = uint32(len(teamState.YellowCards))
	*teamInfo.Timeouts = uint32(teamState.TimeoutsLeft)
	*teamInfo.Goalkeeper = uint32(teamState.Goalkeeper)
	*teamInfo.FoulCounter = uint32(len(teamState.Fouls))
	*teamInfo.BallPlacementFailures = uint32(teamState.BallPlacementFailures)
	*teamInfo.BallPlacementFailuresReached = teamState.BallPlacementFailuresReached
	*teamInfo.CanPlaceBall = teamState.CanPlaceBall
	*teamInfo.MaxAllowedBots = uint32(teamState.MaxAllowedBots)
	*teamInfo.BotSubstitutionIntent = teamState.BotSubstitutionIntend
	*teamInfo.TimeoutTime = mapTime(teamState.TimeoutTimeLeft)
}

func newRefereeMessage() (m *state.Referee) {
	m = new(state.Referee)
	m.PacketTimestamp = new(uint64)
	m.Stage = new(state.Referee_Stage)
	m.StageTimeLeft = new(int32)
	m.Command = new(state.Referee_Command)
	m.CommandCounter = new(uint32)
	m.CommandTimestamp = new(uint64)
	m.Yellow = newTeamInfo()
	m.Blue = newTeamInfo()
	m.BlueTeamOnPositiveHalf = new(bool)
	m.NextCommand = new(state.Referee_Command)
	m.CurrentActionTimeRemaining = new(int32)
	return
}

func newTeamInfo() (t *state.Referee_TeamInfo) {
	t = new(state.Referee_TeamInfo)
	t.Name = new(string)
	t.Score = new(uint32)
	t.RedCards = new(uint32)
	t.YellowCards = new(uint32)
	t.Timeouts = new(uint32)
	t.TimeoutTime = new(uint32)
	t.Goalkeeper = new(uint32)
	t.FoulCounter = new(uint32)
	t.BallPlacementFailures = new(uint32)
	t.BallPlacementFailuresReached = new(bool)
	t.CanPlaceBall = new(bool)
	t.MaxAllowedBots = new(uint32)
	t.BotSubstitutionIntent = new(bool)
	return
}
