package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"encoding/json"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type responseMsg struct {
	status int
	data   Co2DataDto
	err    error
}

func listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * 100)
			sub <- struct{}{}
		}
	}
}

func waitForActivity(sub chan struct{}, m model) tea.Cmd {
	return func() tea.Msg {
		if m.response.status == 0 {
			return checkServer()
		}

		time.Sleep(time.Second * 60)
		<-sub
		return checkServer()
	}
}

func checkServer() tea.Msg {
	apiKey := os.Getenv("X_API_KEY")
	apiUrl := os.Getenv("API_URL")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Set("X-API-KEY", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return responseMsg{
			status: 0,
			data:   Co2DataDto{},
			err:    fmt.Errorf("Error: %s", err),
		}
	}
	defer client.CloseIdleConnections()

	if resp.StatusCode == http.StatusOK {
		var data Co2DataDto
		err := json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return responseMsg{
				status: http.StatusInternalServerError,
				data:   Co2DataDto{},
				err:    fmt.Errorf("Error: %s", err),
			}
		}
		return responseMsg{
			status: resp.StatusCode,
			data:   data,
			err:    nil,
		}
	}

	return nil
}

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

type model struct {
	sub      chan struct{}
	response responseMsg
	spinner  spinner.Model
	quitting bool
}

type Co2DataDto struct {
	CreatedAt time.Time `json:"created_at"`
	CO2       int       `json:"co2"`
	Temp      float32   `json:"temp"`
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m.sub),
		waitForActivity(m.sub, m),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case responseMsg:
		m.response = msg.(responseMsg)
		return m, waitForActivity(m.sub, m) // wait for next event
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	var s string
	if m.response.err != nil {
		s = fmt.Sprintf("\n %s Error: %s; %s", m.spinner.View(), m.response.err.Error(), time.Now().Format("2006-01-02 15:04:05"))
	} else {
		s = fmt.Sprintf("\n %s CO2: %d; Temp: %.1f; %s", m.spinner.View(), m.response.data.CO2, m.response.data.Temp, m.response.data.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\n\nPress any key to exit\n")
	if m.quitting {
		s += "\n"
	}
	return s
}

func main() {
	loadEnvVariables()

	m := model{}
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Globe
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00b8f5"))
	m.sub = make(chan struct{})

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
