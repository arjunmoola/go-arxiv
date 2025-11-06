package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/arjunmoola/go-arxiv/client"
	"time"
	"io"
	"fmt"
	"strings"
)

type App struct {
	arxivClient *client.Client

	searchBar textinput.Model
	searchResults list.Model
	showResults bool
}

type resultEntry struct {
	title string
	author string
	summary viewport.Model
	submittedDate time.Time
	updatedDate time.Time
	pdfLink string
}

func newResultEntry(entry client.Entry) resultEntry {
	vp := viewport.New(5, 10)
	vp.SetContent(entry.Summary)

	return resultEntry{
		title: entry.Title,
		summary: vp,
	}
}

func (e resultEntry) render(w io.Writer, isCursor bool, styles resultEntryDelegateStyles) {
	var builder strings.Builder
	builder.WriteString(styles.title.Render(fmt.Sprintf("Title: %s\n", e.title)))
	//builder.WriteString(styles.author.Render(fmt.Sprintf("Author: %s\n", e.author)))
	builder.WriteString(styles.summary.Render(e.summary.View()))
	view := builder.String()

	if isCursor {
		view = styles.cursor.Render(view)
	}

	io.WriteString(w, view)
}

func (e resultEntry) FilterValue() string {
	return ""
}

func (e resultEntryDelegate) Height() int { return e.height }
func (e resultEntryDelegate) Spacing() int { return 0 }
func (e resultEntryDelegate) Update(msg tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (e resultEntryDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	entry, ok := item.(resultEntry)

	if !ok {
		return
	}

	isCursor := m.Index() == index

	entry.render(w, isCursor, e.styles)
}

type resultEntryDelegate struct {
	styles resultEntryDelegateStyles
	height int
}

type resultEntryDelegateStyles struct {
	title lipgloss.Style
	summary lipgloss.Style
	author lipgloss.Style
	normal lipgloss.Style
	cursor lipgloss.Style
	height int
	width int
}

func newResultEntryDelegateStyles() resultEntryDelegateStyles {
	return resultEntryDelegateStyles{
		title: lipgloss.NewStyle(),
		summary: lipgloss.NewStyle(),
		author: lipgloss.NewStyle(),
		normal: lipgloss.NewStyle(),
		cursor: lipgloss.NewStyle(),
	}
}

func newResultEntryDelegate() resultEntryDelegate {
	return resultEntryDelegate{
		styles: newResultEntryDelegateStyles(),
		height: 10,
	}
}

func newSearchResultList() list.Model {
	model := list.New(nil, newResultEntryDelegate(), 30, 30)
	model.SetShowPagination(false)
	model.SetShowHelp(false)
	model.SetShowStatusBar(false)

	return model
}

func New() *App {
	input := textinput.New()
	input.Width = 20
	client := client.New()

	searchRes := newSearchResultList()

	return &App{
		arxivClient: client,
		searchBar: input,
		searchResults: searchRes,
		
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (a *App) View() string {
	var builder strings.Builder

	builder.WriteString(a.searchBar.View())
	builder.WriteRune('\n')

	if a.showResults {
	}

}
