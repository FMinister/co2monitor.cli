package main

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"encoding/json"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type ResponseMsg struct {
	Status int
	Data   Co2DataDto
	Err    error
}

func listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * 100) // nolint:gosec
			sub <- struct{}{}
		}
	}
}

func waitForActivity(sub chan struct{}, m model) tea.Cmd {
	return func() tea.Msg {
		if m.response.Status == 0 {
			fmt.Println("No response yet")
			return checkServer()
		}

		time.Sleep(time.Second * 10)
		<-sub
		return checkServer()
	}
}

func checkServer() tea.Msg {
	apiKey := os.Getenv("X_API_KEY")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://co2.leyrer.io/api/co2data/1/latest", nil)
	req.Header.Set("X-API-KEY", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return ResponseMsg{
			Status: 0,
			Data:   Co2DataDto{},
			Err:    fmt.Errorf("Error: %s", err),
		}
	}
	defer client.CloseIdleConnections()

	if resp.StatusCode == http.StatusOK {
		var data Co2DataDto
		err := json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return ResponseMsg{
				Status: http.StatusInternalServerError,
				Data:   Co2DataDto{},
				Err:    fmt.Errorf("Error: %s", err),
			}
		}
		return ResponseMsg{
			Status: resp.StatusCode,
			Data:   data,
			Err:    nil,
		}
	}

	return nil
}

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

type model struct {
	sub      chan struct{}
	response ResponseMsg
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
	case ResponseMsg:
		m.response = msg.(ResponseMsg)
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
	s := fmt.Sprintf("\n %s CO2: %d; Temp: %.1f; %s \n\n Press any key to exit\n", m.spinner.View(), m.response.Data.CO2, m.response.Data.Temp, m.response.Data.CreatedAt.Format("2006-01-02 15:04:05"))
	if m.quitting {
		s += "\n"
	}
	return s
}

func main() {
	LoadEnvVariables()
	p := tea.NewProgram(model{
		sub:     make(chan struct{}),
		spinner: spinner.New(),
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
