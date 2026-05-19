package tui

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

const (
	inputHeight      = 3
	logTick          = 200 * time.Millisecond
	mouseScrollLines = 3
)

type model struct {
	ctx              context.Context
	session          *runtime.Session
	agent            agent.Agent
	logs             *LogBuffer
	width            int
	height           int
	input            textinput.Model
	running          bool
	initialRun       bool
	mouseEnabled     bool
	conversation     []string
	conversationView viewport.Model
	logLines         []string
	logView          viewport.Model
	err              error
}

// agentResultMsg carries the completed agent response back into the TUI update
// loop.
type agentResultMsg struct {
	text string
	err  error
}

// logRefreshMsg carries the latest log snapshot into the TUI update loop.
type logRefreshMsg struct {
	lines []string
}

// newModel creates the Bubble Tea model and initializes viewports and input.
func newModel(ctx context.Context, session *runtime.Session, runtimeAgent agent.Agent, logs *LogBuffer, initialPrompt string) model {
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
		agent:            runtimeAgent,
		logs:             logs,
		input:            input,
		initialRun:       initialPrompt != "",
		mouseEnabled:     false,
		conversation:     []string{},
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
	case agentResultMsg:
		m.running = false
		if typed.err != nil {
			m.err = typed.err
			m.conversation = append(m.conversation, "assistant: agent error: "+typed.err.Error())
			m.syncConversation(true)
			return m, nil
		}
		m.conversation = append(m.conversation, "assistant: "+typed.text)
		m.syncConversation(true)
		return m, nil
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
	m.conversation = append(m.conversation, "user: "+prompt)
	m.syncConversation(true)
	return m, runAgent(m.ctx, m.agent, m.session, prompt)
}

// View renders the full TUI frame.
func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	conversationPane := m.conversationView.View()
	logPane := m.logView.View()
	main := lipgloss.JoinHorizontal(lipgloss.Top, conversationPane, logPane)

	return lipgloss.JoinVertical(lipgloss.Left, main, m.renderInput())
}

// renderInput renders the bottom prompt box and transient status text.
func (m model) renderInput() string {
	status := ""
	if m.running {
		status = " running"
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
	m.conversationView.Height = mainHeight
	m.logView.Width = rightWidth
	m.logView.Height = mainHeight
	m.syncConversation(false)
	m.syncLogs(false)
}

// syncConversation refreshes the conversation viewport content.
func (m *model) syncConversation(gotoBottom bool) {
	lines := append([]string{titleStyle.Render("Conversation"), ""}, m.conversation...)
	m.conversationView.SetContent(paneContent(lines, m.conversationView.Width))
	if gotoBottom || m.conversationView.PastBottom() {
		m.conversationView.GotoBottom()
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

// runAgent starts an agent run as a Bubble Tea command.
func runAgent(ctx context.Context, runtimeAgent agent.Agent, session *runtime.Session, prompt string) tea.Cmd {
	return func() tea.Msg {
		message, err := runtimeAgent.Run(ctx, session, prompt)
		if err != nil {
			return agentResultMsg{err: err}
		}
		return agentResultMsg{text: conversation.Text(message)}
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
