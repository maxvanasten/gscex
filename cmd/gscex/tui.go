package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"gscex/pkg/config"
	"gscex/pkg/index"
	"gscex/pkg/search"
)

var (
	accentColor   = lipgloss.Color("#7C3AED")
	textColor     = lipgloss.Color("#E2E8F0")
	dimColor      = lipgloss.Color("#64748B")
	selectedColor = lipgloss.Color("#A78BFA")
	borderColor   = lipgloss.Color("#475569")
	bgColor       = lipgloss.Color("#0F172A")
)

type mode int

const (
	modeSearch mode = iota
	modeResults
	modePreview
	modeHelp
)

type resultItem struct {
	result search.Result
}

func (i resultItem) Title() string { return fmt.Sprintf("%s:%d", i.result.File, i.result.Line) }
func (i resultItem) Description() string {
	ctx := i.result.Context

	// If no context and Line is set, this is a file-only result - show line count
	if len(ctx) == 0 && i.result.Line > 0 {
		return fmt.Sprintf("%d lines", i.result.Line)
	}

	// If no context but Line is 0, show content if available
	if len(ctx) == 0 {
		return i.result.Content
	}

	// Need at least 6 lines from 7-line context (3 before + matched + 3 after)
	// With default ContextLines=3, matched is at index 3, we want indices 1-5
	if len(ctx) < 6 {
		return i.result.Content
	}

	// Extract 5 lines: indices 1-5 from context [0, 1, 2, 3(MATCHED), 4, 5, 6]
	before2 := ctx[1] // 2 lines before matched
	before1 := ctx[2] // 1 line before matched
	matched := ctx[3] // The matched line (center)
	after1 := ctx[4]  // 1 line after matched
	after2 := ctx[5]  // 2 lines after matched

	// Format with matched line highlighted using >> prefix
	return fmt.Sprintf("%s\n%s\n>> %s\n%s\n%s",
		before2, before1, matched, after1, after2)
}
func (i resultItem) FilterValue() string { return i.result.File + i.result.Content }

type fileItem struct {
	path string
}

func (i fileItem) Title() string       { return i.path }
func (i fileItem) Description() string { return "" }
func (i fileItem) FilterValue() string { return i.path }

type model struct {
	mode       mode
	width      int
	height     int
	searchType string // "text", "func", "method", "files"

	// Search input
	searchInput textinput.Model

	// File filter input (shows file filter when activated)
	fileFilterInput textinput.Model
	showFileFilter  bool

	// Results list
	resultsList     list.Model
	contextDelegate list.DefaultDelegate
	compactDelegate list.DefaultDelegate

	// File list (for browsing)
	fileList list.Model

	// Currently selected file content
	currentFile     string
	currentLine     int
	fileContent     []string
	scrollOffset    int
	previewViewport int

	// Help
	showHelp bool

	// Game selection
	currentGame    string
	availableGames []string
	gameFilter     string // empty means load all games

	// Data (support multiple games)
	indices map[string]*index.Index
	engines map[string]*search.Engine

	// Config reference for accessing script paths
	cfg *config.Config

	// Search state
	isSearching  bool
	spinnerFrame int
}

func (m *model) isInputFocused() bool {
	return m.searchInput.Focused() || m.fileFilterInput.Focused()
}

func (m *model) setResultsDelegate() {
	if m.searchType == "files" {
		m.resultsList.SetDelegate(m.compactDelegate)
	} else {
		m.resultsList.SetDelegate(m.contextDelegate)
	}
}

func initialModel(gameFilter string) *model {
	m := &model{
		mode:       modeSearch,
		searchType: "text",
		gameFilter: gameFilter,
	}

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Search for text, functions, or methods..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	m.searchInput = ti

	// Initialize file filter input
	fi := textinput.New()
	fi.Placeholder = "Filter by filename (e.g., zm_tomb)..."
	fi.CharLimit = 100
	fi.Width = 40
	m.fileFilterInput = fi

	// Create delegates for different search types
	contextDelegate := list.NewDefaultDelegate()
	contextDelegate.SetHeight(6) // 1 line for title + 5 lines for description with context

	compactDelegate := list.NewDefaultDelegate()
	compactDelegate.SetHeight(2) // 1 line for title + 1 line for file info

	m.resultsList = list.New([]list.Item{}, contextDelegate, 0, 0)
	m.resultsList.Title = "Search Results"
	m.resultsList.SetShowStatusBar(false)
	m.resultsList.SetFilteringEnabled(false)

	// Store delegates for switching
	m.contextDelegate = contextDelegate
	m.compactDelegate = compactDelegate

	// Initialize file list
	m.fileList = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.fileList.Title = "Files"
	m.fileList.SetShowStatusBar(false)
	m.fileList.SetFilteringEnabled(false)

	return m
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			// Load config
			cfg, err := config.Load()
			if err != nil {
				return errMsg{err: fmt.Errorf("failed to load config: %w", err)}
			}

			// Load indices based on game filter
			var games []string
			if m.gameFilter != "" && m.gameFilter != "all" {
				games = []string{m.gameFilter}
			} else {
				games = []string{"t5", "t6"}
			}
			var loadedGames []string
			var allFiles []string

			for _, game := range games {
				idxPath := cfg.IndexPath(game)
				if _, err := os.Stat(idxPath); os.IsNotExist(err) {
					continue
				}

				idx, err := index.Load(idxPath)
				if err != nil {
					continue
				}

				if m.indices == nil {
					m.indices = make(map[string]*index.Index)
				}
				if m.engines == nil {
					m.engines = make(map[string]*search.Engine)
				}

				m.indices[game] = idx
				m.engines[game] = search.New(idx)
				loadedGames = append(loadedGames, game)

				// Collect files with game prefix if multiple games loaded
				for _, f := range idx.Files {
					if len(games) > 1 {
						allFiles = append(allFiles, fmt.Sprintf("[%s] %s", game, f))
					} else {
						allFiles = append(allFiles, f)
					}
				}
			}

			if len(loadedGames) == 0 {
				return errMsg{err: fmt.Errorf("no indices found. Run 'gscex init' first")}
			}

			// Set default game
			if m.currentGame == "" && len(loadedGames) > 0 {
				m.currentGame = loadedGames[0]
			}

			// Populate file list
			var items []list.Item
			sort.Strings(allFiles)
			for _, f := range allFiles {
				items = append(items, fileItem{path: f})
			}
			m.fileList.SetItems(items)

			return loadedMsg{games: loadedGames, cfg: cfg}
		},
	)
}

type errMsg struct {
	err error
}

type loadedMsg struct {
	games []string
	cfg   *config.Config
}
type searchMsg struct {
	results []search.Result
}
type editorFinishedMsg struct {
	err error
}

type spinnerMsg struct{}

func (m *model) spinnerTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerMsg{}
	})
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update list dimensions
		m.resultsList.SetSize(msg.Width-4, msg.Height-8)
		m.fileList.SetSize(msg.Width/3-4, msg.Height-8)

		return m, nil

	case tea.KeyMsg:
		// Global shortcuts
		switch msg.String() {
		case "ctrl+c":
			// Always allow Ctrl+C to quit
			return m, tea.Quit

		case "q":
			// Block 'q' when typing in an input field
			if m.isInputFocused() {
				return m, nil
			}
			if m.mode != modeSearch {
				m.mode = modeSearch
				m.searchInput.Focus()
				return m, nil
			}
			return m, tea.Quit

		case "?":
			// Block '?' when typing in an input field
			if m.isInputFocused() {
				return m, nil
			}
			m.showHelp = !m.showHelp
			return m, nil

		case "tab":
			// Cycle search types
			switch m.searchType {
			case "text":
				m.searchType = "func"
				m.searchInput.Placeholder = "Search for function..."
			case "func":
				m.searchType = "method"
				m.searchInput.Placeholder = "entity method (e.g., player give_weapon)"
			case "method":
				m.searchType = "files"
				m.searchInput.Placeholder = "Search for files..."
			case "files":
				m.searchType = "text"
				m.searchInput.Placeholder = "Search for text, functions, or methods..."
			}
			m.setResultsDelegate()
			return m, nil

		case "esc":
			if m.isSearching {
				// Cancel search
				m.isSearching = false
				return m, nil
			}
			if m.mode != modeSearch {
				m.mode = modeSearch
				m.searchInput.Focus()
			}
			return m, nil

		case "enter":
			if m.mode == modeSearch {
				// Prevent duplicate searches while one is in progress
				if m.isSearching {
					return m, nil
				}

				if m.fileFilterInput.Focused() {
					// Blur filter and focus search, then execute search
					m.fileFilterInput.Blur()
					m.searchInput.Focus()
				}

				// Start search with loading state
				m.isSearching = true
				m.spinnerFrame = 0
				return m, tea.Batch(m.performSearch(), m.spinnerTick())
			} else if m.mode == modeResults {
				if item, ok := m.resultsList.SelectedItem().(resultItem); ok {
					m.currentFile = item.result.File
					m.loadFileContent(item.result.File, item.result.Line)
					m.mode = modePreview
				}
			}
			return m, nil

		case "ctrl+n":
			if m.mode == modeResults {
				m.resultsList.CursorDown()
			} else if m.mode == modePreview {
				m.scrollDown()
			}
			return m, nil

		case "ctrl+p":
			if m.mode == modeResults {
				m.resultsList.CursorUp()
			} else if m.mode == modePreview {
				m.scrollUp()
			}
			return m, nil

		case "ctrl+f":
			// Toggle file filter visibility and focus
			if !m.showFileFilter {
				// Show filter and focus it
				m.showFileFilter = true
				m.mode = modeSearch
				m.fileFilterInput.Focus()
				m.searchInput.Blur()
			} else if m.fileFilterInput.Focused() {
				// Switch focus back to search input
				m.fileFilterInput.Blur()
				m.searchInput.Focus()
				// Hide filter if empty
				if m.fileFilterInput.Value() == "" {
					m.showFileFilter = false
				}
			} else {
				// Focus filter input
				m.fileFilterInput.Focus()
				m.searchInput.Blur()
			}
			return m, nil

		case "ctrl+b":
			m.mode = modeResults
			return m, nil

		}

	case errMsg:
		m.isSearching = false
		return m, tea.Quit

	case spinnerMsg:
		if m.isSearching {
			m.spinnerFrame = (m.spinnerFrame + 1) % 4
			return m, m.spinnerTick()
		}
		return m, nil

	case searchMsg:
		m.isSearching = false
		var items []list.Item
		for _, r := range msg.results {
			items = append(items, resultItem{result: r})
		}
		m.resultsList.SetItems(items)
		if len(msg.results) > 0 {
			m.mode = modeResults
		}
		return m, nil

	case loadedMsg:
		m.cfg = msg.cfg
		return m, nil

	case editorFinishedMsg:
		// Editor closed, just resume
		return m, nil
	}

	// Update sub-components based on mode
	var cmd tea.Cmd
	if m.mode == modeSearch {
		// Only update the focused input
		if m.fileFilterInput.Focused() {
			m.fileFilterInput, cmd = m.fileFilterInput.Update(msg)
		} else {
			m.searchInput, cmd = m.searchInput.Update(msg)
		}
	} else if m.mode == modeResults {
		m.resultsList, cmd = m.resultsList.Update(msg)
	} else if m.mode == modePreview {
		// Handle preview scrolling
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "down", "j":
				m.scrollDown()
			case "up", "k":
				m.scrollUp()
			case "pgdown":
				for i := 0; i < 10; i++ {
					m.scrollDown()
				}
			case "pgup":
				for i := 0; i < 10; i++ {
					m.scrollUp()
				}
			case "v":
				if m.currentFile != "" {
					return m, m.openEditor()
				}
			}
		}
	}

	return m, cmd
}

func (m *model) performSearch() tea.Cmd {
	return func() tea.Msg {
		query := m.searchInput.Value()
		if query == "" {
			return searchMsg{results: []search.Result{}}
		}

		opts := search.Options{
			MaxResults:   100,
			ContextLines: 3,
		}

		var allResults []search.Result

		for game, eng := range m.engines {
			var results []search.Result

			switch m.searchType {
			case "text":
				results = eng.SearchText(query, opts)
				// Add game prefix if multiple games
				if len(m.engines) > 1 {
					for i := range results {
						results[i].File = fmt.Sprintf("[%s] %s", game, results[i].File)
					}
				}
			case "func":
				// Try exact match first
				fn, usages, ok := eng.SearchFunction(query)
				if ok {
					// Add definition
					if len(m.engines) > 1 {
						results = append(results, search.Result{
							File:    fmt.Sprintf("[%s] %s", game, fn.File),
							Line:    fn.Line,
							Content: fn.Signature,
							Context: fn.Context,
						})
					} else {
						results = append(results, search.Result{
							File:    fn.File,
							Line:    fn.Line,
							Content: fn.Signature,
							Context: fn.Context,
						})
					}
					// Add usages
					for _, u := range usages {
						if len(m.engines) > 1 {
							results = append(results, search.Result{
								File:    fmt.Sprintf("[%s] %s", game, u.File),
								Line:    u.Line,
								Content: u.Content,
								Context: u.Context,
							})
						} else {
							results = append(results, search.Result{
								File:    u.File,
								Line:    u.Line,
								Content: u.Content,
								Context: u.Context,
							})
						}
					}
				}

				// Always do fuzzy search to get ALL matching functions
				fuzzyResults := eng.SearchFunctionsFuzzy(query, 100)
				for _, fn := range fuzzyResults {
					// Skip if already added (exact match case)
					if len(results) > 0 && results[0].File == fn.File && results[0].Line == fn.Line {
						continue
					}
					if len(m.engines) > 1 {
						results = append(results, search.Result{
							File:    fmt.Sprintf("[%s] %s", game, fn.File),
							Line:    fn.Line,
							Content: fn.Signature,
							Context: fn.Context,
						})
					} else {
						results = append(results, search.Result{
							File:    fn.File,
							Line:    fn.Line,
							Content: fn.Signature,
							Context: fn.Context,
						})
					}
				}
			case "method":
				// Parse "entity method" format
				parts := strings.Fields(query)
				if len(parts) >= 2 {
					entries := eng.SearchMethod(parts[0], parts[1])
					for _, e := range entries {
						if len(m.engines) > 1 {
							results = append(results, search.Result{
								File:    fmt.Sprintf("[%s] %s", game, e.File),
								Line:    e.Line,
								Content: e.Content,
								Context: e.Context,
							})
						} else {
							results = append(results, search.Result{
								File:    e.File,
								Line:    e.Line,
								Content: e.Content,
								Context: e.Context,
							})
						}
					}
				}
			case "files":
				files := eng.ListFiles(query)
				for _, f := range files {
					// Get line count from index
					lineCount := 0
					if idx, ok := m.indices[game]; ok {
						if raw, ok := idx.Raw[f]; ok {
							lineCount = strings.Count(raw, "\n") + 1
						}
					}
					if len(m.engines) > 1 {
						results = append(results, search.Result{
							File: fmt.Sprintf("[%s] %s", game, f),
							Line: lineCount,
						})
					} else {
						results = append(results, search.Result{
							File: f,
							Line: lineCount,
						})
					}
				}
			}

			allResults = append(allResults, results...)
		}

		// Sort results
		sort.Slice(allResults, func(i, j int) bool {
			if allResults[i].File == allResults[j].File {
				return allResults[i].Line < allResults[j].Line
			}
			return allResults[i].File < allResults[j].File
		})

		// Apply file filter if active
		fileFilter := m.fileFilterInput.Value()
		if fileFilter != "" {
			fileFilterLower := strings.ToLower(fileFilter)
			var filtered []search.Result
			for _, r := range allResults {
				// Extract actual filename (remove game prefix if present)
				fileName := r.File
				if strings.HasPrefix(fileName, "[") {
					idx := strings.Index(fileName, "]")
					if idx > 0 {
						fileName = strings.TrimSpace(fileName[idx+1:])
					}
				}
				if strings.Contains(strings.ToLower(fileName), fileFilterLower) {
					filtered = append(filtered, r)
				}
			}
			allResults = filtered
		}

		return searchMsg{results: allResults}
	}
}

func (m *model) loadFileContent(file string, highlightLine int) {
	m.currentFile = file
	m.currentLine = highlightLine

	// Extract game prefix if present (e.g., "[t6] some/file.gsc")
	game := ""
	actualFile := file
	if strings.HasPrefix(file, "[") {
		idx := strings.Index(file, "]")
		if idx > 0 {
			game = file[1:idx]
			actualFile = strings.TrimSpace(file[idx+1:])
		}
	}

	// Get the raw content from the appropriate index
	if game != "" {
		if idx, ok := m.indices[game]; ok {
			if raw, ok := idx.Raw[actualFile]; ok {
				m.fileContent = strings.Split(raw, "\n")
			}
		}
	} else if len(m.indices) > 0 {
		// Use the first available index if no game prefix
		for _, idx := range m.indices {
			if raw, ok := idx.Raw[actualFile]; ok {
				m.fileContent = strings.Split(raw, "\n")
				break
			}
		}
	}

	// Center the view on the selected line
	m.calculateScrollOffset()
}

func (m *model) calculateScrollOffset() {
	if m.currentLine == 0 {
		m.scrollOffset = 0
		return
	}
	// Show about 10 lines above the current line
	m.scrollOffset = m.currentLine - 10
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
	// Adjust if we're near the end
	visibleLines := m.previewViewport - 5
	maxScroll := len(m.fileContent) - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
	}
}

func (m *model) scrollDown() {
	maxScroll := len(m.fileContent) - m.previewViewport + 5
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollOffset < maxScroll {
		m.scrollOffset++
	}
}

func (m *model) scrollUp() {
	if m.scrollOffset > 0 {
		m.scrollOffset--
	}
}

func (m *model) openEditor() tea.Cmd {
	return func() tea.Msg {
		if m.cfg == nil {
			return editorFinishedMsg{err: fmt.Errorf("config not loaded")}
		}

		// Parse game prefix from currentFile (e.g., "[t6] some/file.gsc")
		game := ""
		actualFile := m.currentFile
		if strings.HasPrefix(m.currentFile, "[") {
			idx := strings.Index(m.currentFile, "]")
			if idx > 0 {
				game = m.currentFile[1:idx]
				actualFile = strings.TrimSpace(m.currentFile[idx+1:])
			}
		}

		// Determine game if not prefixed
		if game == "" {
			for g := range m.indices {
				game = g
				break
			}
		}

		if game == "" {
			return editorFinishedMsg{err: fmt.Errorf("no game found")}
		}

		// Build source path
		sourcePath := filepath.Join(m.cfg.ScriptsPath(game), actualFile)

		// Create temp directory with timestamp
		tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("gscex-%d", time.Now().Unix()))
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return editorFinishedMsg{err: fmt.Errorf("failed to create temp dir: %w", err)}
		}

		// Copy file to temp location
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return editorFinishedMsg{err: fmt.Errorf("failed to read source file: %w", err)}
		}

		tempPath := filepath.Join(tempDir, filepath.Base(actualFile))
		if err := os.WriteFile(tempPath, content, 0644); err != nil {
			return editorFinishedMsg{err: fmt.Errorf("failed to write temp file: %w", err)}
		}

		// Build neovim command with line number
		cmd := exec.Command("nvim", fmt.Sprintf("+%d", m.currentLine), tempPath)

		// Execute neovim and wait for it to finish
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			return editorFinishedMsg{err: err}
		})()
	}
}

func (m *model) View() string {
	if m.showHelp {
		return m.helpView()
	}

	switch m.mode {
	case modeSearch:
		return m.searchView()
	case modeResults:
		return m.resultsView()
	case modePreview:
		return m.previewView()
	default:
		return m.searchView()
	}
}

func (m *model) searchView() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Show spinner when searching
	if m.isSearching {
		var sb strings.Builder

		headerStyle := lipgloss.NewStyle().
			Background(accentColor).
			Foreground(bgColor).
			Bold(true).
			Padding(0, 2).
			Width(m.width)

		sb.WriteString(headerStyle.Render("GSC Reference Browser - Press ? for help"))
		sb.WriteString("\n\n")

		// Spinner frames
		spinnerFrames := []string{"◐", "◓", "◑", "◒"}
		spinner := spinnerFrames[m.spinnerFrame]

		loadingStyle := lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			MarginTop(2).
			MarginLeft(2)

		sb.WriteString(loadingStyle.Render(fmt.Sprintf("%s Searching...", spinner)))

		// Show what we're searching for
		query := m.searchInput.Value()
		if query != "" {
			queryStyle := lipgloss.NewStyle().
				Foreground(textColor).
				MarginLeft(2)
			sb.WriteString("\n" + queryStyle.Render(fmt.Sprintf("Looking for: %s", query)))
		}

		// Show file filter if active
		if m.fileFilterInput.Value() != "" {
			filterStyle := lipgloss.NewStyle().
				Foreground(dimColor).
				MarginLeft(2)
			sb.WriteString("\n" + filterStyle.Render(fmt.Sprintf("Filter: %s", m.fileFilterInput.Value())))
		}

		// Help text
		helpStyle := lipgloss.NewStyle().
			Foreground(dimColor).
			MarginTop(2).
			MarginLeft(2)
		sb.WriteString("\n" + helpStyle.Render("Press Esc to cancel"))

		return sb.String()
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Background(accentColor).
		Foreground(bgColor).
		Bold(true).
		Padding(0, 2).
		Width(m.width)

	sb.WriteString(headerStyle.Render("GSC Reference Browser - Press ? for help"))
	sb.WriteString("\n")

	// Search type indicator
	typeStyle := lipgloss.NewStyle().
		Foreground(selectedColor).
		Bold(true)

	searchTypes := []string{"text", "func", "method", "files"}
	var typeIndicators []string
	for _, t := range searchTypes {
		if t == m.searchType {
			typeIndicators = append(typeIndicators, typeStyle.Render("["+t+"]"))
		} else {
			typeIndicators = append(typeIndicators, lipgloss.NewStyle().Foreground(dimColor).Render(t))
		}
	}

	sb.WriteString("\n  Mode: " + strings.Join(typeIndicators, " | ") + " (Tab to switch)\n\n")

	// Search input
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1).
		Width(m.width - 4)

	sb.WriteString(inputStyle.Render(m.searchInput.View()))

	// File filter input (if visible)
	if m.showFileFilter {
		filterStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1).
			Width(m.width - 4).
			MarginTop(1)
		if m.fileFilterInput.Focused() {
			filterStyle = filterStyle.BorderForeground(selectedColor)
		}
		sb.WriteString("\n" + filterStyle.Render(m.fileFilterInput.View()))
	}

	// Quick stats
	if len(m.engines) > 0 {
		statsStyle := lipgloss.NewStyle().
			Foreground(dimColor).
			MarginTop(2)

		var statsLines []string
		for game, eng := range m.engines {
			files, funcs, methods := eng.Stats()
			statsLines = append(statsLines, fmt.Sprintf("[%s] %d files, %d functions, %d methods", game, files, funcs, methods))
		}
		sort.Strings(statsLines)
		stats := "Indexed: " + strings.Join(statsLines, " | ")
		sb.WriteString("\n\n  " + statsStyle.Render(stats))
	}

	return sb.String()
}

func (m *model) resultsView() string {
	if m.width == 0 {
		return "Loading..."
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Background(accentColor).
		Foreground(bgColor).
		Bold(true).
		Padding(0, 2).
		Width(m.width)

	sb.WriteString(headerStyle.Render(fmt.Sprintf("Results for: %s", m.searchInput.Value())))
	sb.WriteString("\n\n")

	// Results list
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Height(m.height - 6).
		Width(m.width - 4)

	sb.WriteString(listStyle.Render(m.resultsList.View()))

	// Footer help
	helpStyle := lipgloss.NewStyle().
		Foreground(dimColor).
		MarginTop(1)

	sb.WriteString("\n  " + helpStyle.Render("Enter: view | Ctrl+F: file filter | Ctrl+N/P: navigate | Q: back | ?: help"))

	return sb.String()
}

func (m *model) previewView() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Calculate viewport height
	m.previewViewport = m.height - 10
	if m.previewViewport < 10 {
		m.previewViewport = 10
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Background(accentColor).
		Foreground(bgColor).
		Bold(true).
		Padding(0, 2).
		Width(m.width)

	lineInfo := ""
	if m.currentLine > 0 {
		lineInfo = fmt.Sprintf(" (Line %d)", m.currentLine)
	}
	sb.WriteString(headerStyle.Render(fmt.Sprintf("File: %s%s", m.currentFile, lineInfo)))
	sb.WriteString("\n\n")

	// File content
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1).
		Height(m.previewViewport).
		Width(m.width - 4)

	var content strings.Builder

	// Calculate visible range
	start := m.scrollOffset
	end := start + m.previewViewport - 2
	if end > len(m.fileContent) {
		end = len(m.fileContent)
	}
	if start < 0 {
		start = 0
	}

	// Show visible lines
	for i := start; i < end && i < len(m.fileContent); i++ {
		line := m.fileContent[i]
		lineNum := fmt.Sprintf("%4d | ", i+1)

		// Highlight current line
		if i+1 == m.currentLine {
			lineNum = lipgloss.NewStyle().
				Background(selectedColor).
				Foreground(bgColor).
				Bold(true).
				Render(fmt.Sprintf("%4d > ", i+1))
			line = lipgloss.NewStyle().
				Background(lipgloss.Color("#1E1B4B")).
				Render(line)
		} else {
			lineNum = lipgloss.NewStyle().
				Foreground(dimColor).
				Render(lineNum)
		}

		content.WriteString(lineNum)
		content.WriteString(line)
		content.WriteString("\n")
	}

	sb.WriteString(contentStyle.Render(content.String()))

	// Footer help with scroll position
	helpStyle := lipgloss.NewStyle().
		Foreground(dimColor).
		MarginTop(1)

	scrollInfo := fmt.Sprintf("Line %d/%d", m.scrollOffset+1, len(m.fileContent))
	help := fmt.Sprintf("↑/↓ or j/k: scroll | PgUp/PgDn: page | V: open in nvim | Q: back | ?: help | %s", scrollInfo)
	sb.WriteString("\n  " + helpStyle.Render(help))

	return sb.String()
}

func (m *model) helpView() string {
	if m.width == 0 {
		return "Loading..."
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Background(accentColor).
		Foreground(bgColor).
		Bold(true).
		Padding(0, 2).
		Width(m.width)

	sb.WriteString(headerStyle.Render("Keyboard Shortcuts"))
	sb.WriteString("\n\n")

	help := `
Global:
  ?          Toggle this help
  q          Quit (or go back)
  ctrl+c     Quit

Search Mode:
  tab        Switch search type (text/func/method/files)
  enter      Execute search
  ctrl+f     Toggle file filter input (show/hide)
  esc        Cancel search (when loading)

Search Loading:
  Shows spinner while searching
  Enter blocked during search to prevent duplicates
  Press Esc to cancel an in-progress search

File Filter:
  Filter results to only show files matching substring
  e.g., "zm_tomb" filters to tomb map files
  Persists across searches until cleared

Results Mode:
  ctrl+n     Next result
  ctrl+p     Previous result
  enter      Open file at result

Preview Mode:
  ↑/↓, j/k   Scroll up/down
  pgup/pgdn  Page up/down
  v          Open in neovim (temp copy)
  q          Back to results

Search Types:
  text       Search for any text in all files
  func       Search for function definitions
  method     Search for method calls (format: "entity method")
  files      Search for filenames
`

	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(2).
		Width(m.width - 4)

	sb.WriteString(helpStyle.Render(help))

	sb.WriteString("\n  Press ? or q to close help")

	return sb.String()
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive TUI browser",
	Long:  `Launch an interactive terminal UI for browsing T6 GSC scripts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(initialModel(gameFlag), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to run TUI: %w", err)
		}
		return nil
	},
}
