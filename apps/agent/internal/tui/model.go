package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

const (
	inputHeight      = 3
	logTick          = 200 * time.Millisecond
	mouseScrollLines = 3
)

type model struct {
	ctx              context.Context
	session          *session.Session
	logs             *LogBuffer
	width            int
	height           int
	input            textinput.Model
	running          bool
	status           string
	streamingIndex   int
	initialRun       bool
	mouseEnabled     bool
	conversation     []conversationItem
	conversationView viewport.Model
	lastRunUsage     usage.TokenUsage
	sessionUsage     usage.TokenUsage
	logLines         []string
	logView          viewport.Model
	err              error
}

type conversationRole string

const (
	conversationRoleUser      conversationRole = "user"
	conversationRoleAssistant conversationRole = "assistant"
	conversationRoleActivity  conversationRole = "activity"
)

type conversationItem struct {
	role         conversationRole
	activityKind activityKind
	text         string
}

type activityKind string

const (
	activityKindRound    activityKind = "round"
	activityKindThinking activityKind = "thinking"
	activityKindModel    activityKind = "model"
	activityKindTool     activityKind = "tool"
	activityKindObserve  activityKind = "observe"
	activityKindError    activityKind = "error"
	activityKindDone     activityKind = "done"
	activityKindPlan     activityKind = "plan"
	activityKindStep     activityKind = "step"
)

// agentEventMsg carries one streamed agent event back into the TUI update loop.
type agentEventMsg struct {
	stream <-chan agent.Event
	event  agent.Event
	ok     bool
}

// logRefreshMsg carries the latest log snapshot into the TUI update loop.
type logRefreshMsg struct {
	lines []string
}

// newModel creates the Bubble Tea model and initializes viewports and input.
func newModel(ctx context.Context, session *session.Session, logs *LogBuffer, initialPrompt string) model {
	initialPrompt = strings.TrimSpace(initialPrompt)
	conversationView := viewport.New(1, 1)
	conversationView.MouseWheelDelta = mouseScrollLines
	conversationView.Style = paneStyle
	logView := viewport.New(1, 1)
	logView.MouseWheelDelta = mouseScrollLines
	logView.Style = paneStyle
	input := textinput.New()
	input.Prompt = "> "
	input.PromptStyle = inputPromptStyle
	input.TextStyle = inputTextStyle
	input.Cursor.Style = inputCursorStyle
	input.SetValue(initialPrompt)
	input.Focus()

	return model{
		ctx:              ctx,
		session:          session,
		logs:             logs,
		input:            input,
		streamingIndex:   -1,
		initialRun:       initialPrompt != "",
		mouseEnabled:     false,
		conversation:     []conversationItem{},
		conversationView: conversationView,
		logLines:         []string{},
		logView:          logView,
	}
}

// Init starts periodic log refresh and the text-input cursor blink.
func (m model) Init() tea.Cmd {
	if m.initialRun {
		return tea.Batch(tickLogs(m.logs), textinput.Blink, submitPrompt(m.input.Value()))
	}
	return tea.Batch(tickLogs(m.logs), textinput.Blink)
}

// Update handles Bubble Tea messages and returns the next model state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		m.resizeViewports()
		return m, nil
	case tea.KeyMsg:
		return m.updateKey(typed)
	case tea.MouseMsg:
		if !m.mouseEnabled {
			return m, nil
		}
		return m.updateMouse(typed)
	case agentEventMsg:
		if !typed.ok {
			m.running = false
			return m, nil
		}
		m = m.handleAgentEvent(typed.event)
		return m, waitAgentEvent(typed.stream)
	case logRefreshMsg:
		wasAtBottom := m.logView.AtBottom()
		m.logLines = typed.lines
		m.syncLogs(wasAtBottom)
		return m, tickLogs(m.logs)
	case submitPromptMsg:
		return m.submitPrompt(typed.prompt)
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

// updateKey handles keyboard shortcuts and prompt editing.
func (m model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "ctrl+y":
		m.mouseEnabled = !m.mouseEnabled
		if m.mouseEnabled {
			return m, tea.EnableMouseCellMotion
		}
		return m, tea.DisableMouse
	case "enter":
		return m.submitPrompt(m.input.Value())
	case "q":
		if m.input.Value() == "" && !m.running {
			return m, tea.Quit
		}
	case "pgup":
		m.conversationView.PageUp()
		m.logView.PageUp()
		return m, nil
	case "pgdown":
		m.conversationView.PageDown()
		m.logView.PageDown()
		return m, nil
	case "home":
		m.conversationView.GotoTop()
		m.logView.GotoTop()
		return m, nil
	case "end":
		m.conversationView.GotoBottom()
		m.logView.GotoBottom()
		return m, nil
	}

	if m.running {
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// updateMouse routes mouse wheel events to the pane under the cursor.
func (m model) updateMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Y >= m.mainHeight() {
		return m, nil
	}

	if msg.X < m.width/2 {
		var cmd tea.Cmd
		m.conversationView, cmd = m.conversationView.Update(msg)
		return m, cmd
	}
	var cmd tea.Cmd
	m.logView, cmd = m.logView.Update(msg)
	return m, cmd
}

// submitPromptMsg asks the update loop to submit a prompt after initialization.
type submitPromptMsg struct {
	prompt string
}

// submitPrompt creates a command that submits a prompt on the next update tick.
func submitPrompt(prompt string) tea.Cmd {
	return func() tea.Msg {
		return submitPromptMsg{prompt: prompt}
	}
}

// submitPrompt appends a user message to the pane and starts the agent run.
func (m model) submitPrompt(input string) (tea.Model, tea.Cmd) {
	prompt := strings.TrimSpace(input)
	if prompt == "" || m.running {
		return m, nil
	}
	m.input.SetValue("")
	m.initialRun = false
	m.running = true
	m.err = nil
	m.streamingIndex = -1
	m.conversation = append(m.conversation, conversationItem{
		role: conversationRoleUser,
		text: prompt,
	})
	m.syncConversation(true)
	return m, runAgent(m.ctx, m.session, prompt)
}

// handleAgentEvent applies a streamed agent event to the TUI state.
func (m model) handleAgentEvent(event agent.Event) model {
	switch typed := event.(type) {
	case agent.RoundStartedEvent:
		m.status = fmt.Sprintf("round %d", typed.Round)
		m.streamingIndex = -1
		m.appendActivity(activityKindRound, fmt.Sprintf("round %d", typed.Round))
	case agent.ModelCallStartedEvent:
		m.status = "thinking"
		m.appendActivity(activityKindThinking, "thinking")
	case agent.ModelCallCompletedEvent:
		m.status = "model response"
		m.appendActivity(activityKindModel, modelCompletedText(typed))
	case agent.ModelTextDeltaEvent:
		m.status = "streaming"
		m.appendAssistantDelta(typed.Delta)
	case agent.PlanCreatedEvent:
		m.status = "plan"
		m.streamingIndex = -1
		m.appendActivity(activityKindPlan, fmt.Sprintf("plan: %s", typed.Summary))
	case agent.StepStartedEvent:
		m.status = "step " + typed.StepID
		m.streamingIndex = -1
		m.appendActivity(activityKindStep, fmt.Sprintf("%s: %s", typed.StepID, typed.Description))
	case agent.StepCompletedEvent:
		m.status = "step completed " + typed.StepID
		m.appendActivity(activityKindObserve, fmt.Sprintf("%s completed: %s", typed.StepID, typed.Summary))
	case agent.StepFailedEvent:
		m.status = "step failed " + typed.StepID
		m.appendActivity(activityKindError, fmt.Sprintf("%s failed: %v", typed.StepID, typed.Err))
	case agent.ReplanStartedEvent:
		m.status = "replanning"
		m.appendActivity(activityKindPlan, "replanning after "+typed.StepID)
	case agent.PlanUpdatedEvent:
		m.status = "plan updated"
		m.appendActivity(activityKindPlan, "plan updated: "+typed.Summary)
	case agent.FinalizingEvent:
		m.status = "finalizing"
		m.appendActivity(activityKindThinking, "finalizing")
	case agent.ToolCallStartedEvent:
		m.status = "tool " + typed.Call.Name
		m.streamingIndex = -1
		m.appendActivity(activityKindTool, toolStartedText(typed))
	case agent.ToolCallCompletedEvent:
		m.status = "observe " + typed.Result.Name
		m.appendActivity(activityKindObserve, toolCompletedText(typed))
	case agent.ToolCallFailedEvent:
		m.status = "tool failed " + typed.Call.Name
		m.appendActivity(activityKindError, "tool failed "+typed.Call.Name+": "+typed.Observation)
	case agent.FinalAnswerEvent:
		m.status = "answer"
		m.setFinalAnswer(conversation.Text(typed.Message))
		m.syncConversation(true)
	case agent.RunFailedEvent:
		m.status = ""
		m.streamingIndex = -1
		m.err = typed.Err
		m.updateUsageStats(typed.Stats.TokenUsage)
		m.conversation = append(m.conversation, conversationItem{
			role: conversationRoleAssistant,
			text: "agent error: " + typed.Err.Error(),
		})
		m.syncConversation(true)
	case agent.RunCompletedEvent:
		m.status = ""
		m.streamingIndex = -1
		m.err = nil
		m.updateUsageStats(typed.Stats.TokenUsage)
		m.appendActivity(activityKindDone, "completed: "+string(typed.Reason))
	}
	return m
}

func (m *model) updateUsageStats(lastRun usage.TokenUsage) {
	m.lastRunUsage = lastRun
	m.sessionUsage.Add(lastRun)
}

func (m *model) appendActivity(kind activityKind, text string) {
	m.conversation = append(m.conversation, conversationItem{
		role:         conversationRoleActivity,
		activityKind: kind,
		text:         text,
	})
	m.syncConversation(true)
}

func (m *model) appendAssistantDelta(delta string) {
	if delta == "" {
		return
	}
	if m.streamingIndex < 0 || m.streamingIndex >= len(m.conversation) {
		m.conversation = append(m.conversation, conversationItem{
			role: conversationRoleAssistant,
			text: delta,
		})
		m.streamingIndex = len(m.conversation) - 1
		m.syncConversation(true)
		return
	}
	m.conversation[m.streamingIndex].text += delta
	m.syncConversation(true)
}

func (m *model) setFinalAnswer(text string) {
	if m.streamingIndex >= 0 && m.streamingIndex < len(m.conversation) {
		m.conversation[m.streamingIndex].text = text
		m.streamingIndex = -1
		return
	}
	m.conversation = append(m.conversation, conversationItem{
		role: conversationRoleAssistant,
		text: text,
	})
}

func modelCompletedText(event agent.ModelCallCompletedEvent) string {
	messageText := strings.TrimSpace(event.MessageText)
	if event.TextStreamed {
		messageText = ""
	}
	if len(event.ToolCallNames) == 0 {
		if messageText != "" {
			return fmt.Sprintf(
				"model response: %s, input=%d output=%d",
				activityText(messageText),
				event.Usage.InputTokens,
				event.Usage.OutputTokens,
			)
		}
		return fmt.Sprintf(
			"model response: no tool calls, input=%d output=%d",
			event.Usage.InputTokens,
			event.Usage.OutputTokens,
		)
	}
	if messageText != "" {
		return fmt.Sprintf(
			"model response: %s; tool calls %s, input=%d output=%d",
			activityText(messageText),
			strings.Join(event.ToolCallNames, ", "),
			event.Usage.InputTokens,
			event.Usage.OutputTokens,
		)
	}
	return fmt.Sprintf(
		"model response: tool calls %s, input=%d output=%d",
		strings.Join(event.ToolCallNames, ", "),
		event.Usage.InputTokens,
		event.Usage.OutputTokens,
	)
}

func toolCompletedText(event agent.ToolCallCompletedEvent) string {
	preview := strings.TrimSpace(event.Result.Content)
	if preview == "" {
		return "observe " + event.Result.Name
	}
	return fmt.Sprintf("observe %s: %s", event.Result.Name, activityText(preview))
}

func toolStartedText(event agent.ToolCallStartedEvent) string {
	args := previewToolArguments(event.Call.Arguments)
	if args == "" {
		return "tool " + event.Call.Name
	}
	return fmt.Sprintf("tool %s %s", event.Call.Name, args)
}

func previewToolArguments(arguments map[string]any) string {
	if len(arguments) == 0 {
		return ""
	}
	data, err := json.Marshal(arguments)
	if err != nil {
		return activityText(fmt.Sprintf("%v", arguments))
	}
	return activityText(string(data))
}

func activityText(text string) string {
	return text
}

// View renders the full TUI frame.
func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	conversationPane := m.conversationView.View()
	conversationPane = lipgloss.JoinVertical(lipgloss.Left, conversationPane, m.renderUsageLine())
	logPane := m.logView.View()
	main := lipgloss.JoinHorizontal(lipgloss.Top, conversationPane, logPane)

	return lipgloss.JoinVertical(lipgloss.Left, main, m.renderInput())
}

// renderInput renders the bottom prompt box and transient status text.
func (m model) renderInput() string {
	status := ""
	if m.running {
		status = " " + m.agentStatus()
	}
	if m.err != nil {
		status += " error"
	}
	if m.mouseEnabled {
		status += " mouse-scroll"
	}

	text := m.input.View()
	if m.running {
		text = "> waiting for assistant..."
	}
	contentWidth := m.width - 4
	if contentWidth < 1 {
		contentWidth = 1
	}
	return inputStyle.Width(contentWidth).Height(inputHeight - 2).Render(text + status)
}

func (m model) agentStatus() string {
	if m.status != "" {
		return m.status
	}
	return "running"
}

// mainHeight returns the height available to the conversation and log panes.
func (m model) mainHeight() int {
	mainHeight := m.height - inputHeight
	if mainHeight < 1 {
		return 1
	}
	return mainHeight
}

// resizeViewports recalculates pane sizes after a terminal resize.
func (m *model) resizeViewports() {
	mainHeight := m.mainHeight()
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	inputWidth := m.width - inputStyle.GetHorizontalFrameSize()
	if inputWidth < 1 {
		inputWidth = 1
	}
	m.input.Width = inputWidth
	m.conversationView.Width = leftWidth
	m.conversationView.Height = mainHeight - 1
	if m.conversationView.Height < 1 {
		m.conversationView.Height = 1
	}
	m.logView.Width = rightWidth
	m.logView.Height = mainHeight
	m.syncConversation(false)
	m.syncLogs(false)
}

func (m model) renderUsageLine() string {
	content := fmt.Sprintf(
		"last i/o %d/%d | session i/o %d/%d",
		m.lastRunUsage.InputTokens,
		m.lastRunUsage.OutputTokens,
		m.sessionUsage.InputTokens,
		m.sessionUsage.OutputTokens,
	)
	width := m.conversationView.Width - usageLineStyle.GetHorizontalFrameSize()
	if width < 1 {
		width = 1
	}
	return usageLineStyle.Width(width).Render(content)
}

// syncConversation refreshes the conversation viewport content.
func (m *model) syncConversation(gotoBottom bool) {
	lines := []string{titleStyle.Render("Conversation"), ""}
	for _, item := range m.conversation {
		lines = append(lines, renderConversationItem(item), "")
	}
	m.conversationView.SetContent(paneContent(lines, m.conversationView.Width))
	if gotoBottom || m.conversationView.PastBottom() {
		m.conversationView.GotoBottom()
	}
}

func renderConversationItem(item conversationItem) string {
	switch item.role {
	case conversationRoleUser:
		return userMessageStyle.Render("> " + item.text)
	case conversationRoleAssistant:
		return assistantMessageStyle.Render(item.text)
	case conversationRoleActivity:
		return activityMessageStyle(item.activityKind).Render("- " + item.text)
	default:
		return item.text
	}
}

// syncLogs refreshes the log viewport content.
func (m *model) syncLogs(gotoBottom bool) {
	lines := append([]string{titleStyle.Render("Logs"), ""}, colorLogLines(m.logLines)...)
	m.logView.SetContent(paneContent(lines, m.logView.Width))
	if gotoBottom || m.logView.PastBottom() {
		m.logView.GotoBottom()
	}
}

// paneContentWidth returns the usable text width inside a bordered pane.
func paneContentWidth(width int) int {
	contentWidth := width - paneStyle.GetHorizontalFrameSize()
	if contentWidth < 1 {
		contentWidth = 1
	}
	return contentWidth
}

// paneContent wraps pane lines to the current viewport width.
func paneContent(lines []string, width int) string {
	width = paneContentWidth(width)
	wrapped := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			wrapped = append(wrapped, "")
			continue
		}
		wrapped = append(wrapped, strings.Split(ansi.Wrap(line, width, ""), "\n")...)
	}
	return strings.Join(wrapped, "\n")
}

// colorLogLines applies syntax coloring to structured slog text lines.
func colorLogLines(lines []string) []string {
	colored := make([]string, 0, len(lines))
	for _, line := range lines {
		colored = append(colored, colorLogLine(line))
	}
	return colored
}

// colorLogLine colors key=value fields in one structured log line.
func colorLogLine(line string) string {
	fields := splitLogFields(line)
	if len(fields) == 0 {
		return line
	}

	colored := make([]string, 0, len(fields))
	for _, field := range fields {
		key, value, ok := strings.Cut(field, "=")
		if !ok || key == "" {
			colored = append(colored, field)
			continue
		}
		colored = append(colored, logKeyStyle.Render(key)+logEqualsStyle.Render("=")+logValueStyle(key, value).Render(value))
	}
	return strings.Join(colored, " ")
}

// splitLogFields splits a slog text line while preserving quoted values.
func splitLogFields(line string) []string {
	fields := []string{}
	var builder strings.Builder
	inQuote := false

	for _, r := range line {
		switch r {
		case '"':
			inQuote = !inQuote
			builder.WriteRune(r)
		case ' ':
			if inQuote {
				builder.WriteRune(r)
				continue
			}
			if builder.Len() > 0 {
				fields = append(fields, builder.String())
				builder.Reset()
			}
		default:
			builder.WriteRune(r)
		}
	}
	if builder.Len() > 0 {
		fields = append(fields, builder.String())
	}
	return fields
}

// logValueStyle chooses a color style for a log value based on its key.
func logValueStyle(key string, value string) lipgloss.Style {
	if key == "msg" {
		return logMessageStyle
	}
	if key != "level" {
		return logValueStyleDefault
	}

	switch strings.Trim(value, `"`) {
	case "DEBUG":
		return logDebugStyle
	case "INFO":
		return logInfoStyle
	case "WARN":
		return logWarnStyle
	case "ERROR":
		return logErrorStyle
	default:
		return logValueStyleDefault
	}
}

// runAgent starts an agent event stream as a Bubble Tea command.
func runAgent(ctx context.Context, session *session.Session, prompt string) tea.Cmd {
	stream := session.Run(ctx, prompt)
	return waitAgentEvent(stream)
}

// waitAgentEvent waits for the next event from an agent stream.
func waitAgentEvent(stream <-chan agent.Event) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-stream
		return agentEventMsg{stream: stream, event: event, ok: ok}
	}
}

// tickLogs schedules the next periodic log-buffer snapshot.
func tickLogs(logs *LogBuffer) tea.Cmd {
	return tea.Tick(logTick, func(time.Time) tea.Msg {
		if logs == nil {
			return logRefreshMsg{}
		}
		return logRefreshMsg{lines: logs.Lines()}
	})
}

var paneStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(0, 1)

var inputStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(0, 1)

var inputPromptStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250"))

var inputTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250"))

var inputCursorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("212"))

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("39"))

var userMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("220"))

var assistantMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250"))

var usageLineStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("244")).
	Padding(0, 1)

func activityMessageStyle(kind activityKind) lipgloss.Style {
	switch kind {
	case activityKindRound:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
	case activityKindThinking:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("183"))
	case activityKindModel:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	case activityKindTool:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	case activityKindObserve:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	case activityKindError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	case activityKindDone:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("120"))
	case activityKindPlan:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	case activityKindStep:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	}
}

var logKeyStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("75"))

var logEqualsStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("244"))

var logValueStyleDefault = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250"))

var logMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("229"))

var logDebugStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("244"))

var logInfoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("82"))

var logWarnStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("214"))

var logErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("203"))
