package tui

import (
	"context"
	"strings"
	"time"

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
	input            string
	running          bool
	initialRun       bool
	conversation     []string
	conversationView viewport.Model
	logLines         []string
	logView          viewport.Model
	err              error
}

type agentResultMsg struct {
	text string
	err  error
}

type logRefreshMsg struct {
	lines []string
}

func newModel(ctx context.Context, session *runtime.Session, runtimeAgent agent.Agent, logs *LogBuffer, initialPrompt string) model {
	initialPrompt = strings.TrimSpace(initialPrompt)
	conversationView := viewport.New(1, 1)
	conversationView.MouseWheelDelta = mouseScrollLines
	conversationView.Style = paneStyle
	logView := viewport.New(1, 1)
	logView.MouseWheelDelta = mouseScrollLines
	logView.Style = paneStyle

	return model{
		ctx:              ctx,
		session:          session,
		agent:            runtimeAgent,
		logs:             logs,
		input:            initialPrompt,
		initialRun:       initialPrompt != "",
		conversation:     []string{},
		conversationView: conversationView,
		logLines:         []string{},
		logView:          logView,
	}
}

func (m model) Init() tea.Cmd {
	if m.initialRun {
		return tea.Batch(tickLogs(m.logs), submitPrompt(m.input))
	}
	return tickLogs(m.logs)
}

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
		return m, nil
	}
}

func (m model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "enter":
		return m.submitPrompt(m.input)
	case "backspace":
		if m.input != "" {
			runes := []rune(m.input)
			m.input = string(runes[:len(runes)-1])
		}
		return m, nil
	case "q":
		if m.input == "" && !m.running {
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

	if len(msg.Runes) > 0 && !m.running {
		m.input += string(msg.Runes)
	}
	return m, nil
}

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

type submitPromptMsg struct {
	prompt string
}

func submitPrompt(prompt string) tea.Cmd {
	return func() tea.Msg {
		return submitPromptMsg{prompt: prompt}
	}
}

func (m model) submitPrompt(input string) (tea.Model, tea.Cmd) {
	prompt := strings.TrimSpace(input)
	if prompt == "" || m.running {
		return m, nil
	}
	m.input = ""
	m.initialRun = false
	m.running = true
	m.err = nil
	m.conversation = append(m.conversation, "user: "+prompt)
	m.syncConversation(true)
	return m, runAgent(m.ctx, m.agent, m.session, prompt)
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	conversationPane := m.conversationView.View()
	logPane := m.logView.View()
	main := lipgloss.JoinHorizontal(lipgloss.Top, conversationPane, logPane)

	return lipgloss.JoinVertical(lipgloss.Left, main, m.renderInput())
}

func (m model) renderInput() string {
	status := ""
	if m.running {
		status = " running"
	}
	if m.err != nil {
		status = " error"
	}

	text := "> " + m.input
	if m.running {
		text = "> waiting for assistant..."
	}
	contentWidth := m.width - 4
	if contentWidth < 1 {
		contentWidth = 1
	}
	return inputStyle.Width(contentWidth).Height(inputHeight - 2).Render(text + status)
}

func (m model) mainHeight() int {
	mainHeight := m.height - inputHeight
	if mainHeight < 1 {
		return 1
	}
	return mainHeight
}

func (m *model) resizeViewports() {
	mainHeight := m.mainHeight()
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	m.conversationView.Width = leftWidth
	m.conversationView.Height = mainHeight
	m.logView.Width = rightWidth
	m.logView.Height = mainHeight
	m.syncConversation(false)
	m.syncLogs(false)
}

func (m *model) syncConversation(gotoBottom bool) {
	lines := append([]string{"Conversation", ""}, m.conversation...)
	m.conversationView.SetContent(paneContent(lines, m.conversationView.Width))
	if gotoBottom || m.conversationView.PastBottom() {
		m.conversationView.GotoBottom()
	}
}

func (m *model) syncLogs(gotoBottom bool) {
	lines := append([]string{"Logs", ""}, m.logLines...)
	m.logView.SetContent(paneContent(lines, m.logView.Width))
	if gotoBottom || m.logView.PastBottom() {
		m.logView.GotoBottom()
	}
}

func paneContentWidth(width int) int {
	contentWidth := width - paneStyle.GetHorizontalFrameSize()
	if contentWidth < 1 {
		contentWidth = 1
	}
	return contentWidth
}

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

func runAgent(ctx context.Context, runtimeAgent agent.Agent, session *runtime.Session, prompt string) tea.Cmd {
	return func() tea.Msg {
		message, err := runtimeAgent.Run(ctx, session, prompt)
		if err != nil {
			return agentResultMsg{err: err}
		}
		return agentResultMsg{text: conversation.Text(message)}
	}
}

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
