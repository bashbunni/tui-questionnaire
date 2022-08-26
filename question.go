package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

const (
	// this will affect the:
	// - choices
	// - list title
	// so each time something is submitted, I would need to reset state of the
	// model and increment the index until you're at the end
	Q1 int = iota
	Q2
)

var (
	currentQ  = Q1
	questions []question
)

func initQuestions() {
	questions = []question{}
	questions = append(questions, NewQuestion("What do you want for dinner?",
		[]list.Item{
			item("Ramen"),
			item("Tomato Soup"),
			item("Hamburgers"),
			item("Cheeseburgers"),
			item("Currywurst"),
			item("Okonomiyaki"),
			item("Pasta"),
			item("Fillet Mignon"),
			item("Caviar"),
			item("Just Wine"),
		}))

	questions = append(questions, NewQuestion(
		"What is your favourite lipgloss?",
		[]list.Item{
			item("Laneige"),
			item("Lancome"),
			item("Burt's Bees"),
		}))
}

func handleQ1(answer string) question {
	switch answer {
	case "Ramen", "Pasta":
		return NewQuestion("Do we have noodles?", []list.Item{item("Yes"), item("No")})
	default:
		currentQ = Q2
		return questions[currentQ]
	}
}

func NewQuestion(title string, choices []list.Item) question {
	const defaultWidth = 20
	l := list.New(choices, itemDelegate{}, defaultWidth, listHeight)
	l.Title = title
	// set styles
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return question{list: l}
}

/* TODO
- have an array of question structs, let the enum be based on their indices
- you could change what question gets returned based on an answer to a particular question
-> switch on Q#, have a handler func for that question based on A?
*/

type question struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m question) Init() tea.Cmd {
	return nil
}

func (m question) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
				// TODO: write to file (or wherever you want to store answers to questionnaire)
			}
			if currentQ == len(questions)-1 {
				return m, tea.Quit
			}
			switch currentQ {
			case Q1:
				return handleQ1(string(i)), nil
				// TODO: handle others
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m question) View() string {
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}
