package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/iciwhite/gitplus/internal/ai"
	"github.com/iciwhite/gitplus/internal/auth"
	"github.com/iciwhite/gitplus/internal/config"
	"github.com/iciwhite/gitplus/internal/github"
)

type model struct {
	list        list.Model
	input       textinput.Model
	ghClient    *github.Client
	aiAssistant *ai.Assistant
	cfg         *config.Config
	auth        *auth.OAuthService
	repos       []item
	quitting    bool
}

type item struct {
	name, desc string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.name }

func NewTUI(cfg *config.Config, authService *auth.OAuthService) *tea.Program {
	m := &model{
		cfg:         cfg,
		auth:        authService,
		ghClient:    github.NewClient(authService.GetClient()),
		aiAssistant: ai.NewAssistant(cfg),
		input:       textinput.New(),
	}
	m.input.Placeholder = "Enter command (list, pr, ai, quit)"
	m.input.Focus()

	items := []list.Item{}
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "GitPlus - GitHub Manager"

	return tea.NewProgram(m)
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.fetchReposCmd())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			cmd := strings.TrimSpace(m.input.Value())
			m.input.Reset()
			return m, m.handleCommand(cmd)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		m.list.View(),
		m.input.View(),
	)
}

func (m *model) fetchReposCmd() tea.Cmd {
	return func() tea.Msg {
		repos, err := m.ghClient.ListRepos(m.auth.GetClient().Transport.(*oauth2.Transport).Context(), "")
		if err != nil {
			return err
		}
		items := make([]list.Item, len(repos))
		for i, r := range repos {
			items[i] = item{name: *r.FullName, desc: *r.Description}
		}
		return items
	}
}

func (m *model) handleCommand(cmd string) tea.Cmd {
	switch {
	case cmd == "list":
		return m.fetchReposCmd()
	case strings.HasPrefix(cmd, "pr"):
		return m.createPRCmd(cmd)
	case strings.HasPrefix(cmd, "ai"):
		return m.aiSuggestionCmd(cmd)
	default:
		return nil
	}
}

func (m *model) createPRCmd(cmd string) tea.Cmd {
	return func() tea.Msg {
		return "PR creation not implemented in demo"
	}
}

func (m *model) aiSuggestionCmd(cmd string) tea.Cmd {
	return func() tea.Msg {
		return "AI suggestion not implemented in demo"
	}
}