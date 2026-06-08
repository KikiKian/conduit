package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	accentColor = lipgloss.Color("#00FF9C")
	subtleColor = lipgloss.Color("#3C3C3C")
	whiteColor  = lipgloss.Color("#EDEDED")
	redColor    = lipgloss.Color("#FF5E5B")
	yellowColor = lipgloss.Color("#FFD166")
)

// Styles
var (
	whiteStyle = lipgloss.NewStyle().
			Foreground(whiteColor)

	logoStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(whiteColor).
			Bold(true).
			MarginBottom(1)

	dividerStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	labelStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Width(10)

	dimStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	successStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(redColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(yellowColor)

	selectedStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	unselectedStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor).
			Padding(1, 2).
			MarginTop(1)

	promptStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)
)

const logo = `
  ██████╗ ██████╗ ███╗   ██╗██████╗ ██╗   ██╗██╗████████╗
  ██╔════╝██╔═══██╗████╗  ██║██╔══██╗██║   ██║██║╚══██╔══╝
  ██║     ██║   ██║██╔██╗ ██║██║  ██║██║   ██║██║   ██║
  ██║     ██║   ██║██║╚██╗██║██║  ██║██║   ██║██║   ██║
  ╚██████╗╚██████╔╝██║ ╚████║██████╔╝╚██████╔╝██║   ██║
   ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚═╝   ╚═╝`

var divider = dividerStyle.Render(strings.Repeat("─", 53))

// Style accessors
func DimStyle() lipgloss.Style     { return dimStyle }
func ErrorStyle() lipgloss.Style   { return errorStyle }
func SuccessStyle() lipgloss.Style { return successStyle }

func PrintHelp() {
	fmt.Println(dimStyle.Render("  usage: conduit <command>"))
	fmt.Println(dimStyle.Render("  commands: init, add, watch"))
}

type initStep int

const (
	stepSelectLangs initStep = iota
	stepFilePaths
	stepDone
)

type initModel struct {
	step      initStep
	langs     []string
	selected  map[string]bool
	cursor    int
	inputs    map[string]textinput.Model
	inputKeys []string
	inputIdx  int
	targets   []WatchTarget
	err       string
}

func newInitModel() initModel {
	langs := []string{"typescript", "python", "go", "java"}
	return initModel{
		step:     stepSelectLangs,
		langs:    langs,
		selected: make(map[string]bool),
		inputs:   make(map[string]textinput.Model),
	}
}

func (m initModel) Init() tea.Cmd {
	return nil
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {

		case stepSelectLangs:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.langs)-1 {
					m.cursor++
				}
			case " ":
				lang := m.langs[m.cursor]
				m.selected[lang] = !m.selected[lang]
			case "enter":
				if len(m.selected) == 0 {
					m.err = "select at least one language"
					return m, nil
				}
				m.err = ""
				for _, lang := range m.langs {
					if m.selected[lang] {
						m.inputKeys = append(m.inputKeys, lang)
						ti := textinput.New()
						ti.Placeholder = "./path/to/file"
						ti.Width = 40
						m.inputs[lang] = ti
					}
				}
				first := m.inputs[m.inputKeys[0]]
				first.Focus()
				m.inputs[m.inputKeys[0]] = first
				m.step = stepFilePaths
			}

		case stepFilePaths:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				current := m.inputKeys[m.inputIdx]
				val := strings.TrimSpace(m.inputs[current].Value())
				if val == "" {
					m.err = "file path cannot be empty"
					return m, nil
				}
				m.err = ""
				m.inputIdx++
				if m.inputIdx >= len(m.inputKeys) {
					for _, lang := range m.inputKeys {
						m.targets = append(m.targets, WatchTarget{
							Lang:     lang,
							FilePath: m.inputs[lang].Value(),
						})
					}
					if err := Save(m.targets); err != nil {
						m.err = err.Error()
						return m, nil
					}
					m.step = stepDone
					return m, tea.Quit
				}
				next := m.inputs[m.inputKeys[m.inputIdx]]
				next.Focus()
				m.inputs[m.inputKeys[m.inputIdx]] = next
			default:
				current := m.inputKeys[m.inputIdx]
				ti, cmd := m.inputs[current].Update(msg)
				m.inputs[current] = ti
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m initModel) View() string {
	var b strings.Builder

	b.WriteString(logoStyle.Render(logo) + "\n")
	b.WriteString(dimStyle.Render("  share variables across languages in real-time.") + "\n")
	b.WriteString("  " + divider + "\n\n")

	switch m.step {

	case stepSelectLangs:
		b.WriteString(promptStyle.Render("  ? ") + titleStyle.Render("select languages") + "\n")
		b.WriteString(dimStyle.Render("    space to toggle · enter to confirm") + "\n\n")

		for i, lang := range m.langs {
			cursor := "  "
			if i == m.cursor {
				cursor = promptStyle.Render("❯ ")
			}
			checkbox := "○"
			style := unselectedStyle
			if m.selected[lang] {
				checkbox = "◉"
				style = selectedStyle
			}
			b.WriteString(fmt.Sprintf("  %s%s\n", cursor, style.Render(checkbox+" "+lang)))
		}

		if m.err != "" {
			b.WriteString("\n  " + errorStyle.Render("✗ "+m.err) + "\n")
		}

	case stepFilePaths:
		b.WriteString(promptStyle.Render("  ? ") + titleStyle.Render("set file paths") + "\n\n")

		for i, lang := range m.inputKeys {
			if i < m.inputIdx {
				val := m.inputs[lang].Value()
				b.WriteString(fmt.Sprintf("  %s %s %s\n",
					successStyle.Render("✓"),
					labelStyle.Render(lang),
					dimStyle.Render("→ "+val),
				))
			} else if i == m.inputIdx {
				b.WriteString(fmt.Sprintf("  %s %s %s\n",
					promptStyle.Render("?"),
					labelStyle.Render(lang),
					m.inputs[lang].View(),
				))
			} else {
				b.WriteString(fmt.Sprintf("  %s %s\n",
					dimStyle.Render("·"),
					dimStyle.Render(lang),
				))
			}
		}

		if m.err != "" {
			b.WriteString("\n  " + errorStyle.Render("✗ "+m.err) + "\n")
		}
	}

	return b.String()
}

func RunInit() {
	m := newInitModel()
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		fmt.Println(errorStyle.Render("✗ error: " + err.Error()))
		os.Exit(1)
	}

	final := result.(initModel)
	if final.step == stepDone {
		fmt.Println()
		fmt.Println(boxStyle.Render(
			successStyle.Render("  conduit.config.json created!\n\n") +
				dimStyle.Render("  next steps:\n") +
				warningStyle.Render("    1. ") + whiteStyle.Render("conduit add") + dimStyle.Render(" — add your first variable\n") +
				warningStyle.Render("    2. ") + whiteStyle.Render("conduit watch") + dimStyle.Render(" — start watching for changes\n") +
				warningStyle.Render("    3. ") + dimStyle.Render("use ") + whiteStyle.Render("# conduit:import <name>") + dimStyle.Render(" in your files"),
		))
		fmt.Println()
	}
}

type addField int

const (
	fieldName addField = iota
	fieldType
	fieldValue
	fieldConfirm
)

var supportedTypes = []string{"string", "int", "bool"}

type addModel struct {
	field      addField
	nameInput  textinput.Model
	valInput   textinput.Model
	typeCursor int
	err        string
	done       bool
}

func newAddModel() addModel {
	name := textinput.New()
	name.Placeholder = "e.g. max_retries"
	name.Focus()
	name.Width = 30

	val := textinput.New()
	val.Placeholder = "e.g. 3"
	val.Width = 30

	return addModel{
		nameInput: name,
		valInput:  val,
	}
}

func (m addModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m addModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			switch m.field {
			case fieldName:
				if strings.TrimSpace(m.nameInput.Value()) == "" {
					m.err = "name cannot be empty"
					return m, nil
				}
				m.err = ""
				m.field = fieldType

			case fieldType:
				m.err = ""
				m.field = fieldValue
				m.valInput.Focus()
				return m, textinput.Blink

			case fieldValue:
				if strings.TrimSpace(m.valInput.Value()) == "" {
					m.err = "value cannot be empty"
					return m, nil
				}
				m.err = ""
				m.field = fieldConfirm

			case fieldConfirm:
				varType := supportedTypes[m.typeCursor]
				typed, err := CastValue(varType, m.valInput.Value())
				if err != nil {
					m.err = err.Error()
					m.field = fieldValue
					return m, nil
				}
				if err := AddEntry(m.nameInput.Value(), varType, typed); err != nil {
					m.err = err.Error()
					return m, nil
				}
				m.done = true
				return m, tea.Quit
			}

		case "up", "k":
			if m.field == fieldType && m.typeCursor > 0 {
				m.typeCursor--
			}
		case "down", "j":
			if m.field == fieldType && m.typeCursor < len(supportedTypes)-1 {
				m.typeCursor++
			}
		}

		var cmd tea.Cmd
		switch m.field {
		case fieldName:
			m.nameInput, cmd = m.nameInput.Update(msg)
		case fieldValue, fieldConfirm:
			m.valInput, cmd = m.valInput.Update(msg)
		}
		return m, cmd
	}
	return m, nil
}

func (m addModel) View() string {
	var b strings.Builder

	b.WriteString(logoStyle.Render(logo) + "\n")
	b.WriteString(dimStyle.Render("  share variables across languages in real-time.") + "\n")
	b.WriteString("  " + divider + "\n\n")
	b.WriteString(promptStyle.Render("  ? ") + titleStyle.Render("add a variable") + "\n\n")

	varType := supportedTypes[m.typeCursor]

	if m.field == fieldName {
		b.WriteString(fmt.Sprintf("  %s %s %s\n", promptStyle.Render("?"), labelStyle.Render("name"), m.nameInput.View()))
	} else {
		b.WriteString(fmt.Sprintf("  %s %s %s\n", successStyle.Render("✓"), labelStyle.Render("name"), dimStyle.Render(m.nameInput.Value())))
	}

	if m.field == fieldType {
		b.WriteString(fmt.Sprintf("\n  %s %s\n", promptStyle.Render("?"), labelStyle.Render("type")))
		b.WriteString(dimStyle.Render("    up/down to select · enter to confirm\n\n"))
		for i, t := range supportedTypes {
			cursor := "  "
			if i == m.typeCursor {
				cursor = promptStyle.Render("❯ ")
				b.WriteString(fmt.Sprintf("  %s%s\n", cursor, selectedStyle.Render(t)))
			} else {
				b.WriteString(fmt.Sprintf("  %s%s\n", cursor, unselectedStyle.Render(t)))
			}
		}
	} else if m.field > fieldType {
		b.WriteString(fmt.Sprintf("  %s %s %s\n", successStyle.Render("✓"), labelStyle.Render("type"), dimStyle.Render(varType)))
	}

	if m.field == fieldValue {
		b.WriteString(fmt.Sprintf("\n  %s %s %s\n", promptStyle.Render("?"), labelStyle.Render("value"), m.valInput.View()))
	} else if m.field == fieldConfirm {
		b.WriteString(fmt.Sprintf("  %s %s %s\n", successStyle.Render("✓"), labelStyle.Render("value"), dimStyle.Render(m.valInput.Value())))
	}

	if m.field == fieldConfirm {
		b.WriteString("\n" + boxStyle.Render(
			dimStyle.Render("adding: ")+
				successStyle.Render(m.nameInput.Value())+
				dimStyle.Render(" (")+warningStyle.Render(varType)+dimStyle.Render(")")+
				dimStyle.Render(" = ")+whiteStyle.Render(m.valInput.Value()),
		) + "\n")
		b.WriteString(dimStyle.Render("  press enter to confirm") + "\n")
	}

	if m.err != "" {
		b.WriteString("\n  " + errorStyle.Render("✗ "+m.err) + "\n")
	}

	return b.String()
}

func RunAdd() {
	m := newAddModel()
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		fmt.Println(errorStyle.Render("✗ " + err.Error()))
		os.Exit(1)
	}

	final := result.(addModel)
	if final.done {
		fmt.Println()
		fmt.Println(successStyle.Render("  ✓ ") +
			whiteStyle.Render(final.nameInput.Value()) +
			dimStyle.Render(" ("+supportedTypes[final.typeCursor]+")") +
			dimStyle.Render(" = ") +
			whiteStyle.Render(final.valInput.Value()) +
			dimStyle.Render(" added to .conduit"),
		)
		fmt.Println()
	}
}

func PrintWatchHeader(targets []WatchTarget) {
	var b strings.Builder

	b.WriteString(logoStyle.Render(logo) + "\n")
	b.WriteString(dimStyle.Render("  share variables across languages in real-time.") + "\n")
	b.WriteString("  " + divider + "\n\n")

	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("watching"), dimStyle.Render(".conduit")))

	for i, t := range targets {
		prefix := "           "
		if i == 0 {
			prefix = "  " + labelStyle.Render("targets") + "  "
		}
		b.WriteString(fmt.Sprintf("%s%s %s\n",
			prefix,
			warningStyle.Render(t.Lang),
			dimStyle.Render("→ "+t.FilePath),
		))
	}

	b.WriteString("  " + divider + "\n")
	b.WriteString(successStyle.Render("  ✓ ready") + "\n\n")

	fmt.Print(b.String())
}

func PrintWatchEvent(timestamp, filePath, lang string) {
	fmt.Printf("  %s %s %s %s\n",
		dimStyle.Render("["+timestamp+"]"),
		successStyle.Render("✓"),
		whiteStyle.Render(filePath),
		dimStyle.Render("("+lang+")"),
	)
}

func PrintWatchError(timestamp, filePath string, err error) {
	fmt.Printf("  %s %s %s %s\n",
		dimStyle.Render("["+timestamp+"]"),
		errorStyle.Render("✗"),
		whiteStyle.Render(filePath),
		errorStyle.Render(err.Error()),
	)
}
