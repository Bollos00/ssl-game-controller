package controller

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/config"
	"github.com/RoboCup-SSL/ssl-game-controller/pkg/refproto"
	"github.com/RoboCup-SSL/ssl-game-controller/pkg/timer"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

type GameControllerState struct {
	Division               config.Division                     `json:"division" yaml:"division"`
	AutoContinue           bool                                `json:"autoContinue" yaml:"autoContinue"`
	FirstKickoffTeam       Team                                `json:"firstKickoffTeam" yaml:"firstKickoffTeam"`
	GameEventBehavior      map[GameEventType]GameEventBehavior `json:"gameEventBehavior" yaml:"gameEventBehavior"`
	GameEventProposals     []*GameEventProposal                `json:"gameEventProposals" yaml:"gameEventProposals"`
	AutoRefsConnected      []string                            `json:"autoRefsConnected" yaml:"autoRefsConnected"`
	TeamConnected          map[Team]bool                       `json:"teamConnected" yaml:"teamConnected"`
	TeamConnectionVerified map[Team]bool                       `json:"teamConnectionVerified" yaml:"teamConnectionVerified"`
	MatchState             *State                              `json:"matchState" yaml:"matchState"`
}

func NewGameControllerState() (s *GameControllerState) {
	s = new(GameControllerState)

	s.Division = config.DivA
	s.AutoContinue = true
	s.FirstKickoffTeam = TeamYellow

	s.GameEventBehavior = map[GameEventType]GameEventBehavior{}
	for _, event := range AllGameEvents() {
		s.GameEventBehavior[event] = GameEventBehaviorOn
	}

	s.GameEventProposals = nil

	s.AutoRefsConnected = []string{}
	s.TeamConnected = map[Team]bool{}
	s.TeamConnected[TeamYellow] = false
	s.TeamConnected[TeamBlue] = false
	s.TeamConnectionVerified = map[Team]bool{}
	s.TeamConnectionVerified[TeamYellow] = false
	s.TeamConnectionVerified[TeamBlue] = false

	s.MatchState = NewState()
	return
}

func (s GameControllerState) DeepCopy() (c GameControllerState) {
	c = s
	if len(s.GameEventProposals) > 0 {
		c.GameEventProposals = make([]*GameEventProposal, len(s.GameEventProposals))
		copy(c.GameEventProposals, s.GameEventProposals)
	} else {
		c.GameEventProposals = nil
	}
	c.GameEventBehavior = make(map[GameEventType]GameEventBehavior)
	for k, v := range s.GameEventBehavior {
		c.GameEventBehavior[k] = v
	}
	c.MatchState = s.MatchState.DeepCopy()
	return
}

type Engine struct {
	State           *State
	GcState         *GameControllerState
	StageTimes      map[Stage]time.Duration
	config          config.Game
	TimeProvider    timer.TimeProvider
	LastTimeUpdate  time.Time
	PersistentState *PersistentState
	Geometry        config.Geometry
	Rand            *rand.Rand
}

func NewEngine(config config.Game, seed int64) (e Engine) {
	e.config = config
	e.GcState = NewGameControllerState()
	e.GcState.Division = e.config.DefaultDivision
	e.Geometry = *e.config.DefaultGeometry[e.GcState.Division]
	e.PersistentState = new(PersistentState)
	e.loadStages()
	e.ResetGame()
	e.TimeProvider = func() time.Time { return time.Now() }
	e.LastTimeUpdate = e.TimeProvider()
	e.Rand = rand.New(rand.NewSource(seed))
	return
}

func (e *Engine) loadStages() {
	e.StageTimes = map[Stage]time.Duration{}
	for _, stage := range Stages {
		e.StageTimes[stage] = 0
	}
	e.StageTimes[StageFirstHalf] = e.config.Normal.HalfDuration
	e.StageTimes[StageHalfTime] = e.config.Normal.HalfTimeDuration
	e.StageTimes[StageSecondHalf] = e.config.Normal.HalfDuration
	e.StageTimes[StageOvertimeBreak] = e.config.Normal.BreakAfter
	e.StageTimes[StageOvertimeFirstHalf] = e.config.Overtime.HalfDuration
	e.StageTimes[StageOvertimeHalfTime] = e.config.Overtime.HalfTimeDuration
	e.StageTimes[StageOvertimeSecondHalf] = e.config.Overtime.HalfDuration
	e.StageTimes[StageShootoutBreak] = e.config.Overtime.BreakAfter
}

func (e *Engine) ResetGame() {
	e.State = NewState()
	e.GcState.MatchState = e.State
	e.PersistentState.CurrentState = e.GcState
	e.PersistentState.Protocol = []*ProtocolEntry{}

	for _, team := range []Team{TeamBlue, TeamYellow} {
		e.State.TeamState[team].TimeoutTimeLeft = e.config.Normal.TimeoutDuration
		e.State.TeamState[team].TimeoutsLeft = e.config.Normal.Timeouts
		e.State.TeamState[team].MaxAllowedBots = e.config.MaxBots[e.GcState.Division]
	}

	e.updateNextCommandForStage()
}

// Update updates times and triggers events if needed
func (e *Engine) Update() (timeChanged bool, eventTriggered bool) {
	currentTime := e.TimeProvider()
	delta := currentTime.Sub(e.LastTimeUpdate)
	e.LastTimeUpdate = currentTime

	timeChanged = e.updateTimes(currentTime, delta)
	eventTriggered = e.triggerTimedEvents(delta)

	return
}

// triggerTimedEvents checks for elapsed stage time, timeouts and cards
func (e *Engine) triggerTimedEvents(delta time.Duration) (eventTriggered bool) {
	eventTriggered = false
	if e.countStageTime() {
		if e.State.StageTimeLeft+delta > 0 && e.State.StageTimeLeft <= 0 {
			e.LogTime("Stage time elapsed", "")
			eventTriggered = true
		}
	}

	if e.countCardTime() {
		for team, teamState := range e.State.TeamState {
			eventTriggered = e.removeElapsedYellowCards(team, teamState) || eventTriggered
		}
		e.updateMaxBots()
	}

	if e.State.GameState() == GameStateTimeout && e.State.CommandFor.Known() {
		timeLeft := e.State.TeamState[e.State.CommandFor].TimeoutTimeLeft
		if timeLeft+delta > 0 && timeLeft <= 0 {
			e.LogTime("Timeout time elapsed", e.State.CommandFor)
			eventTriggered = true
		}
	}
	return
}

// updateTimes updates the times of the state
func (e *Engine) updateTimes(currentTime time.Time, delta time.Duration) (timeChanged bool) {

	if e.countStageTime() {
		timeChanged = timeChanged || isNewFullSecond(e.State.StageTimeElapsed, delta)

		e.State.StageTimeElapsed += delta
		e.State.StageTimeLeft -= delta
	}

	if e.countCardTime() {
		for _, teamState := range e.State.TeamState {
			timeChanged = reduceYellowCardTimes(teamState, delta) || timeChanged
		}
	}

	if e.State.GameState() == GameStateTimeout && e.State.CommandFor.Known() {
		timeChanged = timeChanged || isNewFullSecond(e.State.TeamState[e.State.CommandFor].TimeoutTimeLeft, delta)
		e.State.TeamState[e.State.CommandFor].TimeoutTimeLeft -= delta
	}

	if e.State.MatchTimeStart.After(time.Unix(0, 0)) {
		newMatchDuration := currentTime.Sub(e.State.MatchTimeStart)
		timeChanged = timeChanged || isNewFullSecond(e.State.MatchDuration, newMatchDuration-e.State.MatchDuration)
		e.State.MatchDuration = newMatchDuration
	}

	if e.State.CurrentActionDeadline.After(time.Unix(0, 0)) {
		newCurrentActionTimeRemaining := e.State.CurrentActionDeadline.Sub(currentTime)
		timeChanged = timeChanged || isNewFullSecond(e.State.CurrentActionTimeRemaining, newCurrentActionTimeRemaining-e.State.CurrentActionTimeRemaining)
		e.State.CurrentActionTimeRemaining = newCurrentActionTimeRemaining
	}

	return
}

func (e *Engine) countStageTime() bool {
	return e.State.Stage.IsPausedStage() || e.State.GameState() == GameStateRunning
}

func (e *Engine) countCardTime() bool {
	return e.State.GameState() == GameStateRunning
}

func (e *Engine) SendCommand(command RefCommand, forTeam Team) {
	e.State.PrevCommands = append(e.State.PrevCommands, e.State.Command)
	e.State.PrevCommandsFor = append(e.State.PrevCommandsFor, e.State.CommandFor)
	if len(e.State.PrevCommands) > 10 {
		e.State.PrevCommands = e.State.PrevCommands[1:]
		e.State.PrevCommandsFor = e.State.PrevCommandsFor[1:]
	}

	e.State.Command = command
	e.State.CommandFor = forTeam

	switch command {
	case CommandBallPlacement:
		e.setCurrentActionTimeout(e.config.BallPlacementTime)
	case CommandHalt, CommandStop, CommandTimeout, CommandUnknown:
		// nothing to do
	default:
		// reset placement pos
		e.State.PlacementPos = nil
	}

	e.LogCommand()

	if command == CommandTimeout {
		e.State.TeamState[forTeam].TimeoutsLeft--
	} else if command == CommandNormalStart {
		e.updatePreStages()
	}

	if command.ContinuesGame() {
		// reset game events
		if len(e.State.GameEvents) > 0 {
			e.State.GameEvents = []*GameEvent{}
		}
		// reset game event proposals
		if len(e.GcState.GameEventProposals) > 0 {
			e.GcState.GameEventProposals = nil
		}
		// reset ball placement pos and follow ups
		e.State.PlacementPos = nil
		e.State.NextCommand = CommandUnknown
		e.State.NextCommandFor = TeamUnknown

		// update current action timeout
		if command == CommandIndirect || command == CommandDirect {
			e.setCurrentActionTimeout(e.config.FreeKickTime[e.GcState.Division])
		} else if command != CommandKickoff && command != CommandPenalty {
			e.setCurrentActionTimeout(e.config.GeneralTime)
		}
	}
}

func (e *Engine) setCurrentActionTimeout(timeout time.Duration) {
	e.State.CurrentActionTimeRemaining = timeout
	e.State.CurrentActionDeadline = e.TimeProvider().Add(e.State.CurrentActionTimeRemaining)
}

func (e *Engine) AddGameEvent(gameEvent GameEvent) {
	e.LogGameEvent(gameEvent, e.State.DeepCopy())
	e.State.GameEvents = append(e.State.GameEvents, &gameEvent)
}

func (e *Engine) Continue() {
	log.Println("Continue")
	substitutionIntend := e.State.BotSubstitutionIntend()
	if substitutionIntend != TeamUnknown {
		teamProto := substitutionIntend.toProto()
		event := GameEvent{
			Type: GameEventBotSubstitution,
			Details: GameEventDetails{
				BotSubstitution: &refproto.GameEvent_BotSubstitution{ByTeam: &teamProto}}}
		e.AddGameEvent(event)
		e.LogHint("botSubstitution", "game halted for bot substitution", substitutionIntend)
		e.State.TeamState[TeamBlue].BotSubstitutionIntend = false
		e.State.TeamState[TeamYellow].BotSubstitutionIntend = false
		e.SendCommand(CommandHalt, "")
	} else if e.State.Command == CommandHalt {
		// continue with a stop after halt to let the robots prepare
		e.SendCommand(CommandStop, "")
	} else if e.State.NextCommand != CommandUnknown {
		log.Printf("Let game continue with next command")
		e.SendCommand(e.State.NextCommand, e.State.NextCommandFor)
	} else {
		if e.State.Command != CommandStop {
			// halt the game, if not in STOP.
			// Rational: After ball placement and no next command, halt the game to indicate that manual action is required
			// If in STOP, that was most likely triggered manually already and a suddenly halted game might be confusing and not intended
			e.LogHint("missingFollowUp", "No next command available", TeamUnknown)
			e.SendCommand(CommandHalt, "")
		}
	}
}

func (e *Engine) CommandForEvent(event *GameEvent) (command RefCommand, forTeam Team, err error) {
	command = CommandUnknown
	forTeam = TeamUnknown

	primaryGameEvent := e.State.PrimaryGameEvent()
	if (event.IsSecondary() && event.Type != GameEventPlacementFailed) || primaryGameEvent == nil {
		return
	}

	if e.State.bothTeamsCanPlaceBall() &&
		e.State.ballPlacementFailedBefore() &&
		!e.allTeamsFailedPlacement() {
		// A failed placement will result in an indirect free kick for the opposing team.
		command = CommandIndirect
		forTeam = event.ByTeam()
	} else {
		forTeam = event.ByTeam().Opposite()
		switch event.Type {
		case
			GameEventBallLeftFieldTouchLine,
			GameEventAimlessKick,
			GameEventBotKickedBallTooFast,
			GameEventBotDribbledBallTooFar,
			GameEventAttackerDoubleTouchedBall,
			GameEventAttackerInDefenseArea,
			GameEventAttackerTouchedKeeper,
			GameEventKickTimeout,
			GameEventKeeperHeldBall:
			command = CommandIndirect
		case
			GameEventBallLeftFieldGoalLine,
			GameEventIndirectGoal,
			GameEventPossibleGoal,
			GameEventChippedGoal,
			GameEventDefenderInDefenseAreaPartially,
			GameEventAttackerTooCloseToDefenseArea,
			GameEventBotTippedOver,
			GameEventBotCrashUnique,
			GameEventBotPushedBot,
			GameEventBotHeldBallDeliberately:
			command = CommandDirect
		case
			GameEventGoal:
			command = CommandKickoff
		case
			GameEventBotCrashDrawn,
			GameEventNoProgressInGame:
			command = CommandForceStart
		case
			GameEventDefenderInDefenseArea,
			GameEventMultipleCards:
			command = CommandPenalty
		default:
			err = errors.Errorf("Unhandled game event: %v", e.State.GameEvents)
		}

		if forTeam.Known() &&
			e.GcState.Division == config.DivA && // For division A
			!e.State.TeamState[forTeam].CanPlaceBall && // If team in favor can not place the ball
			e.State.TeamState[forTeam.Opposite()].CanPlaceBall && // If opponent team can place the ball
			primaryGameEvent.Type.resultsFromBallLeavingField() { // event is caused by the ball leaving the field
			// All free kicks that were a result of the ball leaving the field, are awarded to the opposing team.
			forTeam = forTeam.Opposite()
		}
	}

	if command.NeedsTeam() && forTeam.Unknown() {
		// if the command needs a team and there is no unique team available, reset command
		command = CommandUnknown
		forTeam = TeamUnknown
	}

	return
}

func (e *Engine) randomTeam() Team {
	if e.Rand.Intn(2) == 0 {
		return TeamYellow
	}
	return TeamBlue
}

func (g GameEventType) resultsFromBallLeavingField() bool {
	switch g {
	case
		GameEventBallLeftFieldTouchLine,
		GameEventBallLeftFieldGoalLine,
		GameEventAimlessKick,
		GameEventIndirectGoal,
		GameEventPossibleGoal,
		GameEventChippedGoal:
		return true
	}
	return false
}

func (s *State) bothTeamsCanPlaceBall() bool {
	return s.TeamState[TeamYellow].CanPlaceBall && s.TeamState[TeamBlue].CanPlaceBall
}

func (s *State) noTeamCanPlaceBall() bool {
	return !s.TeamState[TeamYellow].CanPlaceBall && !s.TeamState[TeamBlue].CanPlaceBall
}

func (s *State) ballPlacementFailedBefore() bool {
	for _, gameEvent := range s.GameEvents {
		if gameEvent.Type == GameEventPlacementFailed {
			return true
		}
	}
	return false
}

func (e *Engine) Process(event Event) error {
	err := e.processEvent(event)
	if err != nil {
		return err
	}
	// do not execute this for timestamp updates
	e.updateMaxBots()
	e.updateNextCommand()
	return nil
}

func (e *Engine) updateMaxBots() {
	for _, team := range []Team{TeamBlue, TeamYellow} {
		max := e.config.MaxBots[e.GcState.Division]
		yellowCards := len(e.State.TeamState[team].YellowCardTimes)
		redCards := e.State.TeamState[team].RedCards
		e.State.TeamState[team].MaxAllowedBots = max - yellowCards - redCards
	}
}

func (e *Engine) updateNextCommand() {
	if e.State.Command == CommandPenalty || e.State.Command == CommandKickoff {
		e.State.NextCommand = CommandNormalStart
		e.State.NextCommandFor = ""
		return
	}
	primaryEvent := e.State.PrimaryGameEvent()
	if primaryEvent == nil {
		return
	}

	if primaryEvent.Type == GameEventDefenderTooCloseToKickPoint {
		e.State.NextCommand, e.State.NextCommandFor =
			e.lastCommand([]RefCommand{CommandIndirect, CommandDirect, CommandKickoff, CommandPenalty})
	} else {
		lastCommand, lastCommandFor :=
			e.lastCommand([]RefCommand{CommandKickoff, CommandPenalty, CommandNormalStart})
		if lastCommand.IsPrepare() {
			// if the last command was a kickoff or penalty, restart with that
			e.State.NextCommand = lastCommand
			e.State.NextCommandFor = lastCommandFor
		} else {
			command, forTeam, err := e.CommandForEvent(primaryEvent)
			if err != nil {
				log.Print("Warn: ", err)
				return
			}
			e.State.NextCommand = command
			e.State.NextCommandFor = forTeam
		}
	}
}

func (e *Engine) lastCommand(commands []RefCommand) (RefCommand, Team) {
	for i := len(e.State.PrevCommands) - 1; i >= 0; i-- {
		for _, cmd := range commands {
			if e.State.PrevCommands[i] == cmd {
				return e.State.PrevCommands[i], e.State.PrevCommandsFor[i]
			}
		}
	}
	return CommandUnknown, TeamUnknown
}

func (e *Engine) processEvent(event Event) error {
	if event.Command != nil {
		return e.processCommand(event.Command)
	} else if event.Modify != nil {
		return e.processModify(event.Modify)
	} else if event.Stage != nil {
		return e.processStage(event.Stage)
	} else if event.Card != nil {
		return e.processCard(event.Card)
	} else if event.Trigger != nil {
		return e.processTrigger(event.Trigger)
	} else if event.GameEvent != nil {
		return e.processGameEvent(event.GameEvent)
	} else if event.RevertProtocolEntry != nil {
		return e.processRevertProtocolEntry(*event.RevertProtocolEntry)
	}
	return errors.New("unknown event")
}

func (e *Engine) processRevertProtocolEntry(id string) error {
	err := e.PersistentState.RevertProtocolEntry(id)
	if err != nil {
		return err
	}
	e.GcState = e.PersistentState.CurrentState
	e.State = e.GcState.MatchState

	if e.State.Command != CommandHalt {
		e.SendCommand(CommandHalt, "")
	}

	log.Printf("Revert protocol entry %v", id)
	return nil
}

func (e *Engine) processCommand(c *EventCommand) error {

	if c.Location != nil {
		e.State.PlacementPos = c.Location
	}

	switch c.Type {
	case CommandDirect, CommandIndirect, CommandKickoff, CommandPenalty, CommandTimeout, CommandBallPlacement:
		if c.ForTeam == nil {
			return errors.Errorf("Team required for %v", c.Type)
		}
		e.SendCommand(c.Type, *c.ForTeam)
	case CommandHalt, CommandStop, CommandForceStart, CommandNormalStart:
		e.SendCommand(c.Type, "")
	default:
		return errors.Errorf("Unknown command: %v", c)
	}

	log.Printf("Processed command %v", *c)
	return nil
}

func (e *Engine) processModify(m *EventModifyValue) error {
	// process team-independent modifies
	if m.Division != nil {
		if *m.Division == config.DivA || *m.Division == config.DivB {
			e.GcState.Division = *m.Division
		} else {
			return errors.Errorf("Invalid division: %v", *m.Division)
		}
	} else if m.AutoContinue != nil {
		e.GcState.AutoContinue = *m.AutoContinue
	} else if m.FirstKickoffTeam != nil {
		e.GcState.FirstKickoffTeam = Team(*m.FirstKickoffTeam)
		e.updateNextCommandForStage()
	} else if m.GameEventBehavior != nil {
		if m.GameEventBehavior.GameEventType == GameEventAll {
			for i := range e.GcState.GameEventBehavior {
				e.GcState.GameEventBehavior[i] = m.GameEventBehavior.GameEventBehavior
			}
		} else {
			e.GcState.GameEventBehavior[m.GameEventBehavior.GameEventType] = m.GameEventBehavior.GameEventBehavior
		}
	} else if m.RemoveGameEvent != nil {
		i := *m.RemoveGameEvent
		if i >= len(e.State.GameEvents) {
			return errors.Errorf("Game event id %d is too large.", i)
		}
		e.State.GameEvents = append(e.State.GameEvents[:i], e.State.GameEvents[i+1:]...)
	} else if err := e.processTeamModify(m); err != nil {
		return err
	}

	// do not log timestamp modification - it will be send very often, if sent
	e.LogModify(*m)
	log.Printf("Processed %v", m)
	return nil
}

func (e *Engine) processTeamModify(m *EventModifyValue) error {
	if m.ForTeam.Unknown() {
		return errors.Errorf("Unknown team: %v", m.ForTeam)
	}
	teamState := e.State.TeamState[m.ForTeam]
	if m.Goals != nil {
		teamState.Goals = *m.Goals
	} else if m.Goalkeeper != nil {
		teamState.Goalkeeper = *m.Goalkeeper
	} else if m.YellowCards != nil {
		teamState.YellowCards = *m.YellowCards
	} else if m.RedCards != nil {
		teamState.RedCards = *m.RedCards
	} else if m.TimeoutsLeft != nil {
		teamState.TimeoutsLeft = *m.TimeoutsLeft
	} else if m.TeamName != nil {
		teamState.Name = *m.TeamName
	} else if m.OnPositiveHalf != nil {
		teamState.OnPositiveHalf = *m.OnPositiveHalf
		e.State.TeamState[m.ForTeam.Opposite()].OnPositiveHalf = !*m.OnPositiveHalf
	} else if m.FoulCounter != nil {
		incremented := *m.FoulCounter > teamState.FoulCounter
		teamState.FoulCounter = *m.FoulCounter
		if incremented {
			if aEvent := e.FoulCounterIncremented(m.ForTeam); aEvent != nil {
				err := e.processGameEvent(aEvent)
				if err != nil {
					return err
				}
			}
		}
	} else if m.BallPlacementFailures != nil {
		incremented := *m.BallPlacementFailures > teamState.BallPlacementFailures
		teamState.BallPlacementFailures = *m.BallPlacementFailures
		if incremented {
			if aEvent := e.PlacementFailuresIncremented(m.ForTeam); aEvent != nil {
				err := e.processGameEvent(aEvent)
				if err != nil {
					return err
				}
			}
		}
	} else if m.YellowCardTime != nil {
		cardId := m.YellowCardTime.CardID
		if cardId < 0 || cardId >= len(teamState.YellowCardTimes) {
			return errors.Errorf("Invalid card index: %v", cardId)
		}
		if duration, err := strToDuration(m.YellowCardTime.Duration); err == nil {
			teamState.YellowCardTimes[cardId] = duration
			sort.SliceStable(teamState.YellowCardTimes, func(i, j int) bool {
				return teamState.YellowCardTimes[i] < teamState.YellowCardTimes[j]
			})
		}
	} else if m.TimeoutTimeLeft != nil {
		if duration, err := strToDuration(*m.TimeoutTimeLeft); err == nil {
			teamState.TimeoutTimeLeft = duration
		} else {
			return err
		}
	} else if m.CanPlaceBall != nil {
		teamState.CanPlaceBall = *m.CanPlaceBall
	} else if m.BotSubstitutionIntend != nil {
		teamState.BotSubstitutionIntend = *m.BotSubstitutionIntend
	} else {
		return errors.Errorf("Unknown modify: %v", m)
	}
	return nil
}

func (e *Engine) processStage(s *EventStage) error {
	if e.State.GameState() != GameStateHalted && e.State.GameState() != GameStateStopped {
		return errors.New("The game state must be halted or stopped to change the stage")
	}

	if s.StageOperation == StageNext {
		e.updateStage(e.State.Stage.Next())
		e.State.TeamState[TeamYellow].BallPlacementFailures = 0
		e.State.TeamState[TeamBlue].BallPlacementFailures = 0
	} else if s.StageOperation == StagePrevious {
		e.updateStage(e.State.Stage.Previous())
	} else if s.StageOperation == StageEndGame {
		e.updateStage(StagePostGame)
	} else {
		return errors.Errorf("Unknown stage operation: %v", s.StageOperation)
	}

	log.Printf("Processed stage %v", s.StageOperation)

	return nil
}

func (e *Engine) updateStage(stage Stage) {
	e.LogStage(stage, e.State.DeepCopy())

	e.State.StageTimeLeft = e.StageTimes[stage]
	e.State.StageTimeElapsed = 0

	if !e.State.Stage.IsPreStage() {
		e.SendCommand(CommandHalt, "")
	}

	if stage == StageFirstHalf {
		e.State.MatchTimeStart = e.TimeProvider()
	}

	if stage.IsPausedStage() {
		e.State.GameEvents = []*GameEvent{}
	}

	if stage == StageOvertimeFirstHalfPre {
		e.State.TeamState[TeamYellow].TimeoutsLeft = e.config.Overtime.Timeouts
		e.State.TeamState[TeamYellow].TimeoutTimeLeft = e.config.Overtime.TimeoutDuration
		e.State.TeamState[TeamBlue].TimeoutsLeft = e.config.Overtime.Timeouts
		e.State.TeamState[TeamBlue].TimeoutTimeLeft = e.config.Overtime.TimeoutDuration
	}

	e.State.Stage = stage

	e.updateNextCommandForStage()
}

func (e *Engine) updateNextCommandForStage() {
	switch e.State.Stage {
	case StagePreGame, StageOvertimeFirstHalfPre:
		e.State.NextCommand = CommandKickoff
		e.State.NextCommandFor = e.GcState.FirstKickoffTeam

	case StageSecondHalfPre, StageOvertimeSecondHalfPre:
		e.State.NextCommand = CommandKickoff
		e.State.NextCommandFor = e.GcState.FirstKickoffTeam.Opposite()
	default:
		e.State.NextCommand = ""
		e.State.NextCommandFor = ""
	}
}

func (e *Engine) updatePreStages() {

	switch e.State.Stage {
	case StagePreGame:
		e.updateStage(StageFirstHalf)
	case StageSecondHalfPre:
		e.updateStage(StageSecondHalf)
	case StageOvertimeFirstHalfPre:
		e.updateStage(StageOvertimeFirstHalf)
	case StageOvertimeSecondHalfPre:
		e.updateStage(StageOvertimeSecondHalf)
	}
}

func (e *Engine) processCard(card *EventCard) (err error) {
	prevState := e.State.DeepCopy()
	if card.ForTeam != TeamYellow && card.ForTeam != TeamBlue {
		return errors.Errorf("Unknown team: %v", card.ForTeam)
	}
	if card.Type != CardTypeYellow && card.Type != CardTypeRed {
		return errors.Errorf("Unknown card type: %v", card.Type)
	}
	teamState := e.State.TeamState[card.ForTeam]
	if card.Operation == CardOperationAdd {
		event := e.addCard(card, card.ForTeam, e.config.YellowCardDuration)
		if event != nil {
			if err := e.processGameEvent(event); err != nil {
				log.Printf("Could not process new additional game event: %v", err)
			}
		}
	} else if card.Operation == CardOperationRevoke {
		err = revokeCard(card, teamState)
	} else if card.Operation == CardOperationModify {
		return modifyCard(card, teamState)
	} else {
		return errors.Errorf("Unknown operation: %v", card.Operation)
	}

	if err == nil {
		e.LogCard(card, prevState)
	}
	return
}

func modifyCard(card *EventCard, teamState *TeamInfo) error {
	if card.Type == CardTypeRed {
		return errors.New("Red cards can not be modified")
	}
	nCardTimes := len(teamState.YellowCardTimes)
	if card.Modification.CardID >= nCardTimes {
		return errors.Errorf("Invalid card id %v. Only %v card times available", card.Modification.CardID, nCardTimes)
	}
	teamState.YellowCardTimes[card.Modification.CardID] = card.Modification.TimeLeft
	return nil
}

func (e *Engine) addCard(card *EventCard, team Team, duration time.Duration) *GameEvent {
	if card.Type == CardTypeYellow {
		log.Printf("Add yellow card for team %v", card.ForTeam)
		e.State.TeamState[team].YellowCards++
		e.State.TeamState[team].YellowCardTimes = append(e.State.TeamState[team].YellowCardTimes, duration)
	} else if card.Type == CardTypeRed {
		log.Printf("Add red card for team %v", card.ForTeam)
		e.State.TeamState[team].RedCards++
	}
	return e.CardNumberIncremented(team)
}

func (e *Engine) processTrigger(t *EventTrigger) (err error) {
	if t.Type == TriggerResetMatch {
		e.ResetGame()
	} else if t.Type == TriggerSwitchColor {
		yellow := e.State.TeamState[TeamYellow]
		e.State.TeamState[TeamYellow] = e.State.TeamState[TeamBlue]
		e.State.TeamState[TeamBlue] = yellow
		e.LogHint("switchColors", "", TeamUnknown)
	} else if t.Type == TriggerSwitchSides {
		yellowOnPositiveHalf := e.State.TeamState[TeamYellow].OnPositiveHalf
		e.State.TeamState[TeamYellow].OnPositiveHalf = !yellowOnPositiveHalf
		e.State.TeamState[TeamBlue].OnPositiveHalf = yellowOnPositiveHalf
		e.LogHint("switchSides", "", TeamUnknown)
	} else if t.Type == TriggerContinue {
		e.Continue()
		e.LogHint("continue", "", TeamUnknown)
	} else {
		return errors.Errorf("Unknown trigger: %v", t.Type)
	}
	log.Printf("Processed trigger %v", t.Type)
	return err
}

func (e *Engine) disabledGameEvent(event GameEventType) bool {
	for eventType, behavior := range e.GcState.GameEventBehavior {
		if event == eventType {
			return behavior == GameEventBehaviorOff
		}
	}
	return false
}

func (e *Engine) processGameEvent(event *GameEvent) error {

	if event.Details.EventType() == GameEventNone {
		return errors.Errorf("Incomplete game event: %v", event)
	}

	e.applyGameEventFilters(event)

	event.Origins = append(event.Origins, collectAllOrigins(e.collectAllMatchingProposals(event))...)
	e.GcState.GameEventProposals = e.collectNonMatchingProposals(event)

	if e.disabledGameEvent(event.Type) {
		e.LogIgnoredGameEvent(*event)
		return nil
	}
	event.SetOccurred(e.TimeProvider())

	var additionalEvents []*GameEvent

	if event.IncrementsFoulCounter() {
		team := event.ByTeam()
		if team.Unknown() {
			e.State.TeamState[TeamYellow].FoulCounter++
			if aEvent := e.FoulCounterIncremented(TeamYellow); aEvent != nil {
				additionalEvents = append(additionalEvents, aEvent)
			}
			e.State.TeamState[TeamBlue].FoulCounter++
			if aEvent := e.FoulCounterIncremented(TeamBlue); aEvent != nil {
				additionalEvents = append(additionalEvents, aEvent)
			}
		} else {
			e.State.TeamState[team].FoulCounter++
			if aEvent := e.FoulCounterIncremented(team); aEvent != nil {
				additionalEvents = append(additionalEvents, aEvent)
			}
		}
	}

	e.AddGameEvent(*event)

	if event.AddsYellowCard() {
		team := event.ByTeam()
		if team.Unknown() {
			return errors.New("Missing team in game event")
		}
		aEvent := e.addCard(&EventCard{Type: CardTypeYellow, ForTeam: team, Operation: CardOperationAdd}, team, e.config.YellowCardDuration)
		if aEvent != nil {
			additionalEvents = append(additionalEvents, aEvent)
		}
	}

	if event.AddsRedCard() {
		team := event.ByTeam()
		if team.Unknown() {
			return errors.New("Missing team in game event")
		}
		aEvent := e.addCard(&EventCard{Type: CardTypeRed, ForTeam: team, Operation: CardOperationAdd}, team, 0)
		if aEvent != nil {
			additionalEvents = append(additionalEvents, aEvent)
		}
	}

	if event.Type == GameEventPlacementFailed {
		team := event.ByTeam()
		if team.Unknown() {
			return errors.New("Missing team in game event")
		}
		e.State.TeamState[team].BallPlacementFailures++
		aEvent := e.PlacementFailuresIncremented(team)
		if aEvent != nil {
			additionalEvents = append(additionalEvents, aEvent)
		}
	}

	if event.Type == GameEventPlacementSucceeded {
		team := event.ByTeam()
		if team.Unknown() {
			return errors.New("Missing team in game event")
		}
		e.State.TeamState[team].BallPlacementFailures = 0
	}

	if event.Type == GameEventGoal {
		team := event.ByTeam()
		if team.Unknown() {
			return errors.New("Missing team in game event")
		}
		e.State.TeamState[team].Goals++
	}

	if event.Type == GameEventBotInterferedPlacement {
		// reset ball placement timer
		e.setCurrentActionTimeout(e.config.BallPlacementTime)
	}

	e.State.PlacementPos = e.BallPlacementPos()

	if e.State.GameState() == GameStateHalted {
		log.Printf("Received a game event while halted: %v", event)
	} else if event.Type == GameEventDefenderTooCloseToKickPoint {
		// stop the game and let bots move away from the ball first. The autoRef will continue the game afterwards
		e.SendCommand(CommandStop, "")
	} else if event.Type == GameEventPossibleGoal {
		// halt the game to let the human referee decide if this was a valid goal
		e.SendCommand(CommandHalt, "")
	} else if !event.IsSkipped() && !event.IsSecondary() && !e.State.Command.IsPrepare() {
		e.placeBall(event)
	} else if e.GcState.AutoContinue && event.IsContinueGame() {
		e.Continue()
	} else {
		log.Printf("No change in game with event %v", event)
	}

	log.Printf("Processed game event %v", event)

	for _, aEvent := range additionalEvents {
		if err := e.processGameEvent(aEvent); err != nil {
			log.Print("Could not process new additional game event: ", err)
		}
	}

	return nil
}

func (e *Engine) placeBall(event *GameEvent) {
	teamInFavor := event.ByTeam().Opposite()

	if teamInFavor.Unknown() {
		// select a team by 50% chance
		teamInFavor = e.randomTeam()
		log.Printf("No team in favor. Choosing randomly: %v", teamInFavor)
	}

	if e.State.PlacementPos == nil || e.State.noTeamCanPlaceBall() {
		// placement not possible, human ref must help out
		e.SendCommand(CommandHalt, "")
		e.LogHint("manualPlacement", "placement pos unknown", TeamUnknown)
		return
	}

	nextCommand, _, _ := e.CommandForEvent(event)
	if nextCommand == CommandPenalty || nextCommand == CommandKickoff {
		e.SendCommand(CommandHalt, "")
		e.LogHint("manualPlacement", "manual placement required for kickoff and penalty", teamInFavor)
	} else if e.GcState.Division == config.DivA && // For division A
		!e.State.TeamState[teamInFavor].CanPlaceBall && // If team in favor can not place the ball
		e.State.TeamState[teamInFavor.Opposite()].CanPlaceBall && // If opponent team can place the ball
		event.Type.resultsFromBallLeavingField() {
		// Rule: All free kicks that were a result of the ball leaving the field, are awarded to the opposing team.
		e.SendCommand(CommandBallPlacement, teamInFavor.Opposite())
		log.Printf("Let opponent place the ball in DivA")
	} else if e.GcState.Division == config.DivB && // For division B
		!e.State.TeamState[teamInFavor].CanPlaceBall && // If team in favor can not place the ball
		e.State.TeamState[teamInFavor.Opposite()].CanPlaceBall && // If opponent team can place the ball
		event.Type != GameEventPlacementFailed { // opponent team has not failed recently
		// Rule: [...] the team is allowed to bring the ball into play, after the ball was placed by the opposing team.
		e.SendCommand(CommandBallPlacement, teamInFavor.Opposite())
		log.Printf("Let opponent place the ball in DivB")
	} else if e.State.TeamState[teamInFavor].CanPlaceBall &&
		!e.allTeamsFailedPlacement() {
		// If team can place ball, let it place
		e.SendCommand(CommandBallPlacement, teamInFavor)
		log.Printf("Let team in favor place ball")
	} else {
		// If team can not place ball, human ref has to help out
		e.SendCommand(CommandHalt, "")
		e.LogHint("manualPlacement", "team can not place ball", teamInFavor)
		return
	}
}

func (e *Engine) allTeamsFailedPlacement() bool {
	possibleFailures := 0
	if e.State.TeamState[TeamYellow].CanPlaceBall {
		possibleFailures++
	}
	if e.State.TeamState[TeamBlue].CanPlaceBall {
		possibleFailures++
	}
	failures := 0
	for _, e := range e.State.GameEvents {
		if e.Type == GameEventPlacementFailed {
			failures++
			if failures >= possibleFailures {
				return true
			}
		}
	}
	return false
}

func (e *Engine) filterAimlessKickForDivA(gameEvent *GameEvent) {
	details := gameEvent.Details
	if e.GcState.Division == config.DivA && gameEvent.Type == GameEventAimlessKick {
		// there is no aimless kick in division A. Map it to a ball left field event
		gameEvent.Type = GameEventBallLeftFieldGoalLine
		aimlessKickDetails := details.AimlessKick
		gameEvent.Details = GameEventDetails{
			BallLeftFieldGoalLine: &refproto.GameEvent_BallLeftField{
				ByBot:    aimlessKickDetails.ByBot,
				ByTeam:   aimlessKickDetails.ByTeam,
				Location: aimlessKickDetails.Location,
			}}
	}
}

func (e *Engine) applyGameEventFilters(gameEvent *GameEvent) {
	e.filterAimlessKickForDivA(gameEvent)
}

func (e *Engine) FoulCounterIncremented(team Team) *GameEvent {
	if e.State.TeamState[team].FoulCounter%e.config.MultipleFoulStep == 0 {
		teamProto := team.toProto()
		event := GameEvent{Type: GameEventMultipleFouls,
			Details: GameEventDetails{MultipleFouls: &refproto.GameEvent_MultipleFouls{ByTeam: &teamProto}}}
		return &event
	}
	return nil
}

func (e *Engine) CardNumberIncremented(team Team) *GameEvent {
	cards := e.State.TeamState[team].YellowCards + e.State.TeamState[team].RedCards
	if cards%e.config.MultipleCardStep == 0 {
		teamProto := team.toProto()
		event := GameEvent{Type: GameEventMultipleCards,
			Details: GameEventDetails{MultipleCards: &refproto.GameEvent_MultipleCards{ByTeam: &teamProto}}}
		return &event
	}
	return nil
}

func (e *Engine) PlacementFailuresIncremented(team Team) *GameEvent {
	if e.State.TeamState[team].BallPlacementFailures == e.config.MultiplePlacementFailures {
		teamProto := team.toProto()
		e.State.TeamState[team].CanPlaceBall = false
		event := GameEvent{Type: GameEventMultiplePlacementFailures,
			Details: GameEventDetails{MultiplePlacementFailures: &refproto.GameEvent_MultiplePlacementFailures{ByTeam: &teamProto}}}
		return &event
	}
	return nil
}

func revokeCard(card *EventCard, teamState *TeamInfo) error {
	if card.Type == CardTypeYellow {
		if teamState.YellowCards > 0 {
			log.Printf("Revoke yellow card for team %v", card.ForTeam)
			teamState.YellowCards--
			nCardTimes := len(teamState.YellowCardTimes)
			if nCardTimes > 0 {
				log.Printf("Revoke yellow card time for team %v", card.ForTeam)
				teamState.YellowCardTimes = teamState.YellowCardTimes[:nCardTimes-1]
			}
		} else {
			return errors.Errorf("No yellow cards left to revoke for team %v", card.ForTeam)
		}
	} else if card.Type == CardTypeRed {
		if teamState.RedCards > 0 {
			log.Printf("Revoke red card for team %v", card.ForTeam)
			teamState.RedCards--
		} else {
			return errors.Errorf("No red cards left to revoke for team %v", card.ForTeam)
		}
	}
	return nil
}

func strToDuration(s string) (duration time.Duration, err error) {
	duration = 0
	err = nil

	parts := strings.Split(s, ":")
	if len(parts) > 2 {
		err = errors.Errorf("Invalid duration format: %v", s)
		return
	}

	var secondsIndex int
	var minutes int
	if len(parts) == 1 {
		secondsIndex = 0
		minutes = 0
	} else {
		secondsIndex = 1
		minutes, err = strconv.Atoi(parts[0])
		if err != nil {
			return
		}
	}
	seconds, err := strconv.Atoi(parts[secondsIndex])
	if err != nil {
		return
	}

	duration += time.Minute * time.Duration(minutes)
	duration += time.Second * time.Duration(seconds)
	return
}

func reduceYellowCardTimes(teamState *TeamInfo, delta time.Duration) (newFullSecond bool) {
	newFullSecond = false
	for i := range teamState.YellowCardTimes {
		newFullSecond = newFullSecond || isNewFullSecond(teamState.YellowCardTimes[i], delta)
		teamState.YellowCardTimes[i] -= delta
	}
	return
}

func isNewFullSecond(t time.Duration, delta time.Duration) bool {
	return time.Duration(t+delta).Truncate(time.Second) > t.Truncate(time.Second)
}

func (e *Engine) removeElapsedYellowCards(team Team, teamState *TeamInfo) (removed bool) {
	removed = false
	b := teamState.YellowCardTimes[:0]
	for _, x := range teamState.YellowCardTimes {
		if x > 0 {
			b = append(b, x)
		} else {
			e.LogTime("Yellow card time elapsed", team)
			removed = true
		}
	}
	teamState.YellowCardTimes = b
	return
}
