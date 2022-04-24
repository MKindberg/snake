package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"flag"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pair struct {
	x int
	y int
}

func contains(haystack []pair, needle pair) bool {
	for _, p := range haystack {
		if p.x == needle.x && p.y == needle.y {
			return true
		}
	}
	return false
}

type model struct {
	width  int
	height int
	snake  []pair
	food   pair
	dir    string
	score  int
	lost   bool
}

func (m model) Init() tea.Cmd {
	return nil
}

type tickMsg struct {
	idx int
}

var idx int

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.dir = msg.String()
	case tickMsg:
		if msg.idx < idx { // event expired
			return m, nil
		}
	}
	x, y := m.snake[0].x, m.snake[0].y
	switch m.dir {

	// These keys should exit the program.
	case "r":
		return initialModel(m.width, m.height), nil
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k", "w":
		if y > 0 && !contains(m.snake, pair{x, y - 1}) {
			y--
		} else {
			m.lost = true
		}

	case "down", "j", "s":
		if y < m.height-1 && !contains(m.snake, pair{x, y + 1}) {
			y++
		} else {
			m.lost = true
		}

	case "left", "h", "a":
		if x > 0 && !contains(m.snake, pair{x - 1, y}) {
			x--
		} else {
			m.lost = true
		}

	case "right", "l", "d":
		if x < m.width-1 && !contains(m.snake, pair{x + 1, y}) {
			x++
		} else {
			m.lost = true
		}
	}

	if !m.lost {
		if m.food.x == x && m.food.y == y {
			m.snake = append([]pair{{x, y}}, m.snake...)

			fx, fy := rand.Intn(m.width), rand.Intn(m.height)
			for contains(m.snake, pair{fx, fy}) {
				fx, fy = rand.Intn(m.width), rand.Intn(m.height)
			}
			m.food.x, m.food.y = fx, fy
			m.score++
		} else {
			m.snake = append([]pair{{x, y}}, m.snake[:len(m.snake)-1]...)
		}

		idx++
		curIdx := idx
		sleep := time.Duration(50) * time.Millisecond
		if m.score < 45 {
			sleep = time.Duration(500-m.score*10) * time.Millisecond
		}
		cmd = tea.Tick(sleep, func(t time.Time) tea.Msg {
			return tickMsg{curIdx}
		})
	}

	return m, cmd
}

func (m model) View() string {
	const BODY = 1
	const HEAD = 2
	const FOOD = 3

	var bodyChar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Background(lipgloss.Color("46")).
		Render("X")
	var headChar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Render("O")
	var foodChar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("228")).
		Render("o")
	var lossStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("1"))

	board := make([][]int, m.height)
	for r := 0; r < len(board); r++ {
		board[r] = make([]int, m.width)
	}
	board[m.food.y][m.food.x] = FOOD
	for i := 1; i < len(m.snake); i++ {
		board[m.snake[i].y][m.snake[i].x] = BODY
	}
	board[m.snake[0].y][m.snake[0].x] = HEAD

	s := ""
	if m.lost {
		s = fmt.Sprintf("Final score: %d\n\n", m.score)
	} else {
		s = fmt.Sprintf("Score: %d\n\n", m.score)
	}

	for i := 0; i < m.width+2; i++ {
		if m.lost {
			s += lossStyle.Render("-")
		} else {
			s += "-"
		}
	}
	s += "\n"
	for _, row := range board {
		if m.lost {
			s += lossStyle.Render("|")
		} else {
			s += "|"
		}
		for _, elem := range row {
			switch elem {
			case HEAD:
				s += headChar
			case BODY:
				s += bodyChar
			case FOOD:
				s += foodChar
			default:
				s += " "
			}
		}
		if m.lost {
			s += lossStyle.Render("|")
		} else {
			s += "|"
		}
		s += "\n"
	}
	for i := 0; i < m.width+2; i++ {
		if m.lost {
			s += lossStyle.Render("-")
		} else {
			s += "-"
		}
	}

	s += "\nPress r to restart.\n"
	s += "Press q to quit.\n"

	return s
}

func initialModel(width int, height int) model {
	return model{
		width:  width,
		height: height,
		snake:  []pair{{width / 2, height / 2}},
		food:   pair{rand.Intn(width), rand.Intn(height)},
	}
}

func main() {
	var width, height int
	flag.IntVar(&width, "width", 100, "Width of gameboard")
	flag.IntVar(&height, "height", 20, "Height of gameboard")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(width, height))
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
