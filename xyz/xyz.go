package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Simulate a process that sends events at an irregular interval in real time.
func listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * time.Duration(500)) // Sleep for 500 milliseconds
			sub <- struct{}{}
		}
	}
}

// A command that makes an HTTP request to leyrer.io to get the status.
func getStatus(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		// Make an HTTP GET request to leyrer.io
		resp, err := http.Get("https://leyrer.io")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		// If the status code is 200, we consider it a success
		if resp.StatusCode == http.StatusOK {
			sub <- struct{}{}
		}

		return nil
	}
}

type model struct {
	sub      chan struct{} // where we'll receive activity notifications
	spinner  spinner.Model
	quitting bool
	status   string
}

type responseMsg struct {
	status int
	err    error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m.sub), // generate activity
		getStatus(m.sub),         // check the status
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case string:
		m.status = "Status: OK"
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {
	s := fmt.Sprintf("\n %s %s\n\n Press any key to exit\n", m.spinner.View(), m.status)
	if m.quitting {
		s += "\n"
	}
	return s
}

func main() {
	p := tea.NewProgram(model{
		sub:     make(chan struct{}),
		spinner: spinner.New(),
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
