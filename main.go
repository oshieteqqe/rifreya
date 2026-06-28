package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Task struct {
	title, desc string
	isCompleted bool
}

func (t Task) Title() string       {
	if t.isCompleted {
		return "[X] " + t.title	
	}
	return "[ ] " + t.title }

	func (t Task) Description() string { return t.desc }
	func (t Task) FilterValue() string { return t.title }
	///
	type listKeyMap struct {
		keyToggle key.Binding
	}

	func newListKeyMap() *listKeyMap {
		return &listKeyMap{
			keyToggle: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "complete task"),
			),
		}
	}

	func newTaskDelegate(keys *listKeyMap) list.DefaultDelegate {
		d := list.NewDefaultDelegate()

		d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
			task, ok := m.SelectedItem().(Task)
			if !ok {
				return nil
			}

			if kp, ok := msg.(tea.KeyPressMsg); ok {
				if key.Matches(kp, keys.keyToggle) {

					idx := m.Index()
					items := m.Items()
					items[idx] = Task{
						title:       task.title,
						desc:        task.desc,
						isCompleted: !task.isCompleted,
					}
					m.SetItems(items)

					status := string(task.title)+"\tIncomplete"
					if !task.isCompleted {
						status =string(task.title)+"\tCompleted"
					}
					return m.NewStatusMessage(status)
				}
			}
			return nil
		}

		help := []key.Binding{keys.keyToggle}
		d.ShortHelpFunc = func() []key.Binding { return help }
		d.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{help} }

		return d
	}
	///
	type model struct {
		taskList list.Model
	}


	func (m model) Init() tea.Cmd {
		return nil
	}

	func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if m.taskList.FilterState() == list.Filtering {
				break
			}
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit 

			}
		case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.taskList.SetSize(msg.Width-h, msg.Height-v)

		}
		var cmd tea.Cmd
		m.taskList, cmd = m.taskList.Update(msg)
		return m, cmd
	}



	func (m model) View() tea.View {
		v := tea.NewView(docStyle.Render(m.taskList.View()))
		v.AltScreen = true
		return v
	}

	func main() {
		tasks := []list.Item{ 
			Task{title: "Cleaning ", desc: "F", isCompleted: true},
			Task{title: "Cleaniff2f", desc: "F", isCompleted: false},
		}
		keys := newListKeyMap()
		delegate := newTaskDelegate(keys)

		l := list.New(tasks, delegate, 0, 0)
		l.Title = "To Do for U"

		m := model{taskList: l}

		p := tea.NewProgram(m)
		if _,err := p.Run(); err != nil {
			fmt.Printf("oh my error: %v", err)
			os.Exit(1)
		}
	} 
