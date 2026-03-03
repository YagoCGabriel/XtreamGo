package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// =============================================================================
// TEMA
// =============================================================================

var (
	colorPurple  = lipgloss.Color("#7B61FF")
	colorCyan    = lipgloss.Color("#00D4FF")
	colorGreen   = lipgloss.Color("#00F5A0")
	colorRed     = lipgloss.Color("#FF4D6A")
	colorYellow  = lipgloss.Color("#FFD93D")
	colorDim     = lipgloss.Color("#4A4A6A")
	colorMuted   = lipgloss.Color("#6E6E9A")
	colorText    = lipgloss.Color("#E2E2F0")
	colorBgPanel = lipgloss.Color("#13131F")
	colorBorder  = lipgloss.Color("#2A2A45")

	successStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(colorRed).Bold(true)
	dimStyle     = lipgloss.NewStyle().Foreground(colorDim)
	mutedStyle   = lipgloss.NewStyle().Foreground(colorMuted)
	cyanStyle    = lipgloss.NewStyle().Foreground(colorCyan)
	purpleStyle  = lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	yellowStyle  = lipgloss.NewStyle().Foreground(colorYellow)

	logoStyle = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)

	versionStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Italic(true)

	breadcrumbSepStyle = lipgloss.NewStyle().
				Foreground(colorBorder)

	breadcrumbActiveStyle = lipgloss.NewStyle().
				Foreground(colorCyan).
				Bold(true)

	footerStyle = lipgloss.NewStyle().
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			Foreground(colorMuted).
			Padding(0, 2)

	keyStyle = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)

	listStyle = lipgloss.NewStyle().Padding(0, 1)

	settingLabelStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Width(22)

	settingValueStyle = lipgloss.NewStyle().
				Foreground(colorCyan).
				Bold(true)

	settingHintStyle = lipgloss.NewStyle().
				Foreground(colorDim).
				Italic(true)

	settingKeyBindStyle = lipgloss.NewStyle().
				Foreground(colorPurple).
				Bold(true)

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(colorPurple).
				Bold(true).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(colorBorder)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)
)

// =============================================================================
// DELEGATE CUSTOMIZADO
// =============================================================================

type customDelegate struct{}

func newDelegate() customDelegate { return customDelegate{} }

func (d customDelegate) Height() int                              { return 2 }
func (d customDelegate) Spacing() int                            { return 0 }
func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d customDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(menuItem)
	if !ok {
		return
	}

	// separador visual
	if i.id == "separator" {
		fmt.Fprintf(w, "\n%s", dimStyle.Render("  "+i.title))
		return
	}

	isSelected := index == m.Index()
	width := m.Width() - 4
	if width < 20 {
		width = 20
	}

	// item de busca global — estilo amarelo destacado
	if i.id == "search_all" {
		if isSelected {
			fmt.Fprintf(w, "%s\n%s",
				lipgloss.NewStyle().
					Foreground(colorYellow).Bold(true).
					Background(lipgloss.Color("#1A1A2E")).
					BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(colorYellow).
					PaddingRight(1).
					Render("▶  󰍉  "+i.title),
				lipgloss.NewStyle().
					Foreground(colorMuted).PaddingLeft(5).
					Background(lipgloss.Color("#1A1A2E")).
					Render(i.desc),
			)
		} else {
			fmt.Fprintf(w, "%s\n%s",
				yellowStyle.PaddingLeft(3).Render("󰍉  "+i.title),
				dimStyle.PaddingLeft(5).Render(i.desc),
			)
		}
		return
	}

	icon := itemIcon(i.id, i.extra)

	if isSelected {
		indicator := lipgloss.NewStyle().Foreground(colorPurple).Bold(true).Render("▶ ")
		title := lipgloss.NewStyle().
			Foreground(colorText).Bold(true).
			Width(width - 4).
			Render(icon + i.title)
		desc := lipgloss.NewStyle().
			Foreground(colorMuted).
			Width(width - 2).
			PaddingLeft(4).
			Render(i.desc)
		sel := lipgloss.NewStyle().
			Background(lipgloss.Color("#1A1A2E")).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorPurple).
			PaddingRight(1)
		fmt.Fprintf(w, "%s\n%s",
			sel.Render(indicator+title),
			lipgloss.NewStyle().Background(lipgloss.Color("#1A1A2E")).PaddingLeft(1).Render(desc),
		)
	} else {
		title := lipgloss.NewStyle().
			Foreground(colorMuted).Width(width).PaddingLeft(3).
			Render(icon + i.title)
		desc := lipgloss.NewStyle().
			Foreground(colorDim).Width(width).PaddingLeft(5).
			Render(i.desc)
		fmt.Fprintf(w, "%s\n%s", title, desc)
	}
}

func itemIcon(id string, extra any) string {
	switch id {
	case "live":
		return "󰋙 "
	case "vod":
		return "󰿎 "
	case "series":
		return "󰎁 "
	case "settings":
		return "󰒓 "
	case "search_all":
		return "󰍉 "
	}
	switch extra.(type) {
	case LiveStream:
		return "󰋙 "
	case VODStream:
		return "󰿎 "
	case Series:
		return "󰎁 "
	case Episode:
		return "󰏃 "
	}
	if strings.HasPrefix(id, "T") {
		return "󰄛 "
	}
	return "  "
}

// =============================================================================
// MENU ITEM
// =============================================================================

type menuItem struct {
	title, desc string
	id          string
	intID       int
	extra       any
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

func searchAllItem(desc string) menuItem {
	return menuItem{title: "Buscar em tudo", desc: desc, id: "search_all"}
}

func separatorItem() menuItem {
	return menuItem{title: "───────────────────────────", id: "separator"}
}

func prependSearchAll(desc string, items []list.Item) []list.Item {
	result := make([]list.Item, 0, len(items)+2)
	result = append(result, searchAllItem(desc), separatorItem())
	result = append(result, items...)
	return result
}

// =============================================================================
// SCREENS
// =============================================================================

type screen int

const (
	screenMain screen = iota
	screenLiveCats
	screenLiveStreams
	screenVODCats
	screenVODStreams
	screenSeriesCats
	screenSeriesList
	screenSeasons
	screenEpisodes
	screenSettings
	screenServerSelect
)

// =============================================================================
// MESSAGES
// =============================================================================

type itemsLoadedMsg struct{ items []list.Item }
type loadErrMsg struct{ err error }
type seriesInfoLoadedMsg struct{ info *SeriesInfo }
type configSavedMsg struct{}

// =============================================================================
// MODEL
// =============================================================================

type appModel struct {
	client           *Client
	cfg              *Config
	screen           screen
	list             list.Model
	spinner          spinner.Model
	loading          bool
	status           string
	errMsg           string
	playOpts         PlayOpts
	selectedCatID    string
	selectedSeriesID int
	seriesInfo       *SeriesInfo
	width, height    int
	breadcrumb       []string
}

func newAppModel(client *Client, cfg *Config, playOpts PlayOpts) appModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(colorPurple)

	l := list.New(nil, newDelegate(), 80, 20)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(colorPurple)
	l.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(colorText)
	l.Title = ""
	l.Styles.Title = lipgloss.NewStyle()
	l.SetItems(mainMenuItems())

	return appModel{
		client:     client,
		cfg:        cfg,
		screen:     screenMain,
		list:       l,
		spinner:    sp,
		playOpts:   playOpts,
		breadcrumb: []string{"Início"},
	}
}

func mainMenuItems() []list.Item {
	return []list.Item{
		menuItem{title: "TV ao Vivo", desc: "Canais em tempo real", id: "live"},
		menuItem{title: "Filmes (VOD)", desc: "Video on demand", id: "vod"},
		menuItem{title: "Series", desc: "Episodios por temporada", id: "series"},
		menuItem{title: "Configuracoes", desc: "Player, servidor, hwdec", id: "settings"},
	}
}

func (m appModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// =============================================================================
// UPDATE
// =============================================================================

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		if m.screen == screenSettings {
			return m.updateSettings(msg)
		}
		if m.screen == screenServerSelect {
			return m.updateServerSelect(msg)
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "backspace":
			return m.goBack()
		case "enter", " ":
			return m.handleEnter()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-2, msg.Height-9)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case itemsLoadedMsg:
		m.loading = false
		m.list.SetItems(msg.items)
		m.status = ""
		m.errMsg = ""

	case loadErrMsg:
		m.loading = false
		m.errMsg = msg.err.Error()
		m.status = ""

	case configSavedMsg:
		m.status = "salvo"

	case seriesInfoLoadedMsg:
		m.loading = false
		m.seriesInfo = msg.info
		seasons := make([]string, 0, len(msg.info.Seasons))
		for k := range msg.info.Seasons {
			seasons = append(seasons, k)
		}
		sort.Strings(seasons)
		items := make([]list.Item, len(seasons))
		for i, s := range seasons {
			items[i] = menuItem{
				title: "Temporada " + s,
				desc:  fmt.Sprintf("%d episodio(s)", len(msg.info.Seasons[s])),
				id:    "T" + s,
				extra: s,
			}
		}
		m.list.SetItems(items)
		m.status = ""
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, tea.Batch(cmd, m.spinner.Tick)
}

// =============================================================================
// SETTINGS
// =============================================================================

func (m appModel) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	srv := &m.cfg.Servers[m.cfg.Current]
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc", "backspace":
		return m.goHome()
	case "p":
		options := []string{"mpv", "vlc", "kmplayer"}
		cur := srv.Player
		if cur == "" {
			cur = "mpv"
		}
		next := "mpv"
		for i, o := range options {
			if o == cur {
				next = options[(i+1)%len(options)]
				break
			}
		}
		srv.Player = next
		m.playOpts.Player = next
		return m, m.saveConfig()
	case "h":
		options := []string{"no", "auto", "vaapi", "vdpau", "nvdec", "videotoolbox"}
		cur := srv.HWDec
		if cur == "" {
			cur = "no"
		}
		next := "no"
		for i, o := range options {
			if o == cur {
				next = options[(i+1)%len(options)]
				break
			}
		}
		srv.HWDec = next
		m.playOpts.HWDec = next
		return m, m.saveConfig()
	case "f":
		m.playOpts.Fullscreen = !m.playOpts.Fullscreen
		srv.Fullscreen = m.playOpts.Fullscreen
		return m, m.saveConfig()
	case "s":
		m.screen = screenServerSelect
		m.breadcrumb = []string{"Início", "Configuracoes", "Servidores"}
		m.list.SetFilteringEnabled(false)
		items := make([]list.Item, len(m.cfg.Servers))
		for i, s := range m.cfg.Servers {
			desc := s.URL
			if i == m.cfg.Current {
				desc = "● ativo  |  " + s.URL
			}
			items[i] = menuItem{title: s.Name, desc: desc, id: fmt.Sprintf("%d", i)}
		}
		m.list.SetItems(items)
	}
	return m, nil
}

func (m appModel) updateServerSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc", "backspace":
		return m.openSettings()
	case "enter", " ":
		sel, ok := m.list.SelectedItem().(menuItem)
		if !ok {
			break
		}
		var idx int
		fmt.Sscanf(sel.id, "%d", &idx)
		m.cfg.Current = idx
		srv := m.cfg.Servers[idx]
		m.client = NewClient(srv.URL, srv.Username, srv.Password)
		m.playOpts = PlayOpts{Player: srv.Player, HWDec: srv.HWDec, Fullscreen: srv.Fullscreen}
		m.status = "servidor alterado"
		return m, m.saveConfig()
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m appModel) openSettings() (tea.Model, tea.Cmd) {
	m.screen = screenSettings
	m.list.SetFilteringEnabled(false)
	m.list.SetItems(nil)
	m.errMsg = ""
	m.breadcrumb = []string{"Início", "Configuracoes"}
	return m, nil
}

func (m appModel) saveConfig() tea.Cmd {
	cfg := m.cfg
	return func() tea.Msg {
		if err := saveConfig(cfg); err != nil {
			return loadErrMsg{err}
		}
		return configSavedMsg{}
	}
}

// =============================================================================
// VIEW
// =============================================================================

func (m appModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	header := m.viewHeader(w)
	breadcrumb := m.viewBreadcrumb(w)
	footer := m.viewFooter(w)

	var body string
	if m.screen == screenSettings {
		body = m.viewSettings(w)
	} else {
		body = listStyle.Render(m.list.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, breadcrumb, body, footer)
}

func (m appModel) viewHeader(w int) string {
	logo := logoStyle.Render("󰿎  XtreamGo")
	version := versionStyle.Render("v1.0")

	srv := ""
	if len(m.cfg.Servers) > 0 {
		s := m.cfg.Servers[m.cfg.Current]
		player := s.Player
		if player == "" {
			player = "mpv"
		}
		srv = dimStyle.Render("  ") +
			mutedStyle.Render(s.Name) +
			dimStyle.Render("  ") +
			purpleStyle.Render(player)
	}

	left := lipgloss.JoinHorizontal(lipgloss.Center, logo, "  ", version)
	gap := w - lipgloss.Width(left) - lipgloss.Width(srv) - 4
	if gap < 1 {
		gap = 1
	}

	line := left + strings.Repeat(" ", gap) + srv
	return lipgloss.NewStyle().
		Background(colorBgPanel).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Width(w).
		Padding(0, 2).
		Render(line)
}

func (m appModel) viewBreadcrumb(w int) string {
	parts := make([]string, len(m.breadcrumb))
	for i, b := range m.breadcrumb {
		if i == len(m.breadcrumb)-1 {
			parts[i] = breadcrumbActiveStyle.Render(b)
		} else {
			parts[i] = mutedStyle.Render(b)
		}
	}
	sep := breadcrumbSepStyle.Render("  ›  ")
	crumb := strings.Join(parts, sep)

	status := ""
	if m.loading {
		status = "  " + m.spinner.View() + " " + dimStyle.Render("carregando...")
	} else if m.status != "" {
		status = "  " + successStyle.Render("✓ "+m.status)
	} else if m.errMsg != "" {
		status = "  " + errorStyle.Render("✗ "+m.errMsg)
	}

	return lipgloss.NewStyle().
		Width(w).
		Padding(0, 2).
		Render(crumb + status)
}

func (m appModel) viewSettings(w int) string {
	if len(m.cfg.Servers) == 0 {
		return dimStyle.Render("  Nenhum servidor configurado.")
	}

	srv := m.cfg.Servers[m.cfg.Current]
	player := srv.Player
	if player == "" {
		player = "mpv"
	}
	hwdec := srv.HWDec
	if hwdec == "" {
		hwdec = "no"
	}
	fs := "não"
	if m.playOpts.Fullscreen {
		fs = "sim"
	}

	innerW := w - 8
	if innerW < 40 {
		innerW = 40
	}

	row := func(key, bind, val, hint string) string {
		k := settingLabelStyle.Render(key)
		b := settingKeyBindStyle.Render("[" + bind + "]")
		v := settingValueStyle.Render(val)
		h := settingHintStyle.Render(hint)
		return "  " + k + b + "  " + v + "  " + h
	}

	section := func(title string) string {
		return "\n  " + sectionTitleStyle.Width(innerW).Render(title) + "\n\n"
	}

	var sb strings.Builder
	sb.WriteString(section("Servidor Ativo"))
	sb.WriteString("  " + settingLabelStyle.Render("Nome") + cyanStyle.Render(srv.Name) + "\n")
	sb.WriteString("  " + settingLabelStyle.Render("URL") + mutedStyle.Render(srv.URL) + "\n")
	sb.WriteString("  " + settingLabelStyle.Render("Usuário") + mutedStyle.Render(srv.Username) + "\n")

	sb.WriteString(section("Reprodução"))
	sb.WriteString(row("Player       ", "p", player, "mpv → vlc → kmplayer") + "\n")
	sb.WriteString(row("HW Decode    ", "h", hwdec, "no → auto → vaapi → vdpau → nvdec → videotoolbox") + "\n")
	sb.WriteString(row("Fullscreen   ", "f", fs, "abrir em tela cheia") + "\n")

	sb.WriteString(section("Servidores"))
	sb.WriteString(row("Trocar       ", "s", fmt.Sprintf("%d configurado(s)", len(m.cfg.Servers)), "selecionar servidor ativo") + "\n")

	return panelStyle.Width(w - 4).Render(sb.String())
}

func (m appModel) viewFooter(w int) string {
	var keys []string
	switch m.screen {
	case screenSettings:
		keys = []string{"p player", "h hwdec", "f fullscreen", "s servidores", "esc voltar"}
	case screenServerSelect:
		keys = []string{"enter selecionar", "esc voltar", "q sair"}
	default:
		keys = []string{"enter selecionar", "esc voltar", "/ filtrar", "q sair"}
	}

	parts := make([]string, len(keys))
	for i, k := range keys {
		fields := strings.SplitN(k, " ", 2)
		if len(fields) == 2 {
			parts[i] = keyStyle.Render(fields[0]) + dimStyle.Render(" "+fields[1])
		} else {
			parts[i] = dimStyle.Render(k)
		}
	}
	return footerStyle.Width(w).Render(strings.Join(parts, dimStyle.Render("  ·  ")))
}

// =============================================================================
// HANDLE ENTER
// =============================================================================

func (m appModel) handleEnter() (tea.Model, tea.Cmd) {
	sel, ok := m.list.SelectedItem().(menuItem)
	if !ok {
		return m, nil
	}
	// ignorar clique no separador
	if sel.id == "separator" {
		return m, nil
	}
	m.errMsg = ""

	switch m.screen {

	// ── Menu principal ────────────────────────────────────────────────────────
	case screenMain:
		switch sel.id {
		case "live":
			return m.navigate(screenLiveCats, []string{"Início", "TV ao Vivo"}, func() tea.Msg {
				cats, err := m.client.GetLiveCategories()
				if err != nil {
					return loadErrMsg{err}
				}
				return itemsLoadedMsg{catsToItems(cats)}
			})
		case "vod":
			return m.navigate(screenVODCats, []string{"Início", "Filmes"}, func() tea.Msg {
				cats, err := m.client.GetVODCategories()
				if err != nil {
					return loadErrMsg{err}
				}
				items := prependSearchAll("Pesquisar em todos os filmes", catsToItems(cats))
				return itemsLoadedMsg{items}
			})
		case "series":
			return m.navigate(screenSeriesCats, []string{"Início", "Series"}, func() tea.Msg {
				cats, err := m.client.GetSeriesCategories()
				if err != nil {
					return loadErrMsg{err}
				}
				items := prependSearchAll("Pesquisar em todas as series", catsToItems(cats))
				return itemsLoadedMsg{items}
			})
		case "settings":
			return m.openSettings()
		}

	// ── TV ao Vivo — categorias ───────────────────────────────────────────────
	case screenLiveCats:
		m.selectedCatID = sel.id
		catID := sel.id
		return m.navigate(screenLiveStreams, append(m.breadcrumb, sel.title), func() tea.Msg {
			streams, err := m.client.GetLiveStreams(catID)
			if err != nil {
				return loadErrMsg{err}
			}
			items := make([]list.Item, len(streams))
			for i, s := range streams {
				items[i] = menuItem{title: s.Name, id: fmt.Sprintf("%d", s.ID), extra: s}
			}
			return itemsLoadedMsg{items}
		})

	// ── TV ao Vivo — streams ──────────────────────────────────────────────────
	case screenLiveStreams:
		s := sel.extra.(LiveStream)
		opts := m.playOpts
		opts.Title = s.Name
		go Play(m.client.LiveStreamURL(s.ID, s.ContainerExtension), opts)
		m.status = s.Name

	// ── Filmes — categorias ───────────────────────────────────────────────────
	case screenVODCats:
		if sel.id == "search_all" {
			return m.navigate(screenVODStreams, []string{"Início", "Filmes", "Todos"}, func() tea.Msg {
				streams, err := m.client.GetVODStreams("")
				if err != nil {
					return loadErrMsg{err}
				}
				return itemsLoadedMsg{vodStreamsToItems(streams)}
			})
		}
		m.selectedCatID = sel.id
		catID := sel.id
		return m.navigate(screenVODStreams, append(m.breadcrumb, sel.title), func() tea.Msg {
			streams, err := m.client.GetVODStreams(catID)
			if err != nil {
				return loadErrMsg{err}
			}
			return itemsLoadedMsg{vodStreamsToItems(streams)}
		})

	// ── Filmes — streams ──────────────────────────────────────────────────────
	case screenVODStreams:
		v := sel.extra.(VODStream)
		opts := m.playOpts
		opts.Title = v.Name
		go Play(m.client.VODStreamURL(v.ID, v.ContainerExtension), opts)
		m.status = v.Name

	// ── Series — categorias ───────────────────────────────────────────────────
	case screenSeriesCats:
		if sel.id == "search_all" {
			return m.navigate(screenSeriesList, []string{"Início", "Series", "Todas"}, func() tea.Msg {
				series, err := m.client.GetSeries("")
				if err != nil {
					return loadErrMsg{err}
				}
				return itemsLoadedMsg{seriesToItems(series)}
			})
		}
		m.selectedCatID = sel.id
		catID := sel.id
		return m.navigate(screenSeriesList, append(m.breadcrumb, sel.title), func() tea.Msg {
			series, err := m.client.GetSeries(catID)
			if err != nil {
				return loadErrMsg{err}
			}
			return itemsLoadedMsg{seriesToItems(series)}
		})

	// ── Series — lista ────────────────────────────────────────────────────────
	case screenSeriesList:
		s := sel.extra.(Series)
		m.selectedSeriesID = s.ID
		sid := s.ID
		return m.navigate(screenSeasons, append(m.breadcrumb, s.Name), func() tea.Msg {
			info, err := m.client.GetSeriesInfo(sid)
			if err != nil {
				return loadErrMsg{err}
			}
			return seriesInfoLoadedMsg{info}
		})

	// ── Series — temporadas ───────────────────────────────────────────────────
	case screenSeasons:
		seasonKey := strings.TrimPrefix(sel.id, "T")
		m.screen = screenEpisodes
		m.breadcrumb = append(m.breadcrumb, sel.title)
		episodes := m.seriesInfo.Seasons[seasonKey]
		sort.Slice(episodes, func(i, j int) bool {
			return episodes[i].EpisodeNum.Int() < episodes[j].EpisodeNum.Int()
		})
		items := make([]list.Item, len(episodes))
		for i, ep := range episodes {
			t := fmt.Sprintf("E%02d", ep.EpisodeNum.Int())
			if ep.Title != "" {
				t += "  " + ep.Title
			}
			items[i] = menuItem{title: t, id: ep.ID, extra: ep}
		}
		m.list.SetItems(items)

	// ── Episodios ─────────────────────────────────────────────────────────────
	case screenEpisodes:
		ep := sel.extra.(Episode)
		opts := m.playOpts
		opts.Title = sel.title
		go Play(m.client.SeriesStreamURL(ep.ID, ep.ContainerExtension), opts)
		m.status = sel.title
	}

	return m, nil
}

// =============================================================================
// NAVIGATE / GO BACK
// =============================================================================

func (m appModel) navigate(s screen, crumb []string, loader tea.Cmd) (tea.Model, tea.Cmd) {
	m.screen = s
	m.loading = true
	m.status = ""
	m.errMsg = ""
	m.breadcrumb = crumb
	m.list.SetFilteringEnabled(true)
	m.list.SetItems(nil)
	return m, tea.Batch(loader, m.spinner.Tick)
}

func (m appModel) goHome() (tea.Model, tea.Cmd) {
	m.screen = screenMain
	m.breadcrumb = []string{"Início"}
	m.errMsg, m.status = "", ""
	m.loading = false
	m.list.SetFilteringEnabled(true)
	m.list.SetItems(mainMenuItems())
	return m, nil
}

func (m appModel) goBack() (tea.Model, tea.Cmd) {
	m.errMsg, m.status = "", ""
	m.list.SetFilteringEnabled(true)

	switch m.screen {
	case screenMain:
		return m, tea.Quit

	case screenLiveCats, screenVODCats, screenSeriesCats, screenSettings:
		return m.goHome()

	case screenLiveStreams:
		return m.navigate(screenLiveCats, []string{"Início", "TV ao Vivo"}, func() tea.Msg {
			cats, err := m.client.GetLiveCategories()
			if err != nil {
				return loadErrMsg{err}
			}
			return itemsLoadedMsg{catsToItems(cats)}
		})

	case screenVODStreams:
		return m.navigate(screenVODCats, []string{"Início", "Filmes"}, func() tea.Msg {
			cats, err := m.client.GetVODCategories()
			if err != nil {
				return loadErrMsg{err}
			}
			items := prependSearchAll("Pesquisar em todos os filmes", catsToItems(cats))
			return itemsLoadedMsg{items}
		})

	case screenSeriesList:
		return m.navigate(screenSeriesCats, []string{"Início", "Series"}, func() tea.Msg {
			cats, err := m.client.GetSeriesCategories()
			if err != nil {
				return loadErrMsg{err}
			}
			items := prependSearchAll("Pesquisar em todas as series", catsToItems(cats))
			return itemsLoadedMsg{items}
		})

	case screenSeasons:
		catID := m.selectedCatID
		prev := breadcrumbPop(m.breadcrumb)
		return m.navigate(screenSeriesList, prev, func() tea.Msg {
			series, err := m.client.GetSeries(catID)
			if err != nil {
				return loadErrMsg{err}
			}
			return itemsLoadedMsg{seriesToItems(series)}
		})

	case screenEpisodes:
		sid := m.selectedSeriesID
		prev := breadcrumbPop(m.breadcrumb)
		return m.navigate(screenSeasons, prev, func() tea.Msg {
			info, err := m.client.GetSeriesInfo(sid)
			if err != nil {
				return loadErrMsg{err}
			}
			return seriesInfoLoadedMsg{info}
		})
	}

	return m, nil
}

// =============================================================================
// HELPERS
// =============================================================================

func catsToItems(cats []Category) []list.Item {
	items := make([]list.Item, len(cats))
	for i, c := range cats {
		items[i] = menuItem{title: c.Name, id: c.ID}
	}
	return items
}

func vodStreamsToItems(streams []VODStream) []list.Item {
	items := make([]list.Item, len(streams))
	for i, s := range streams {
		desc := ""
		if s.Rating.String() != "" {
			desc = "⭐ " + s.Rating.String()
		}
		items[i] = menuItem{title: s.Name, desc: desc, id: fmt.Sprintf("%d", s.ID), extra: s}
	}
	return items
}

func seriesToItems(series []Series) []list.Item {
	items := make([]list.Item, len(series))
	for i, s := range series {
		items[i] = menuItem{
			title: s.Name,
			desc:  truncate(s.Plot, 80),
			id:    fmt.Sprintf("%d", s.ID),
			intID: s.ID,
			extra: s,
		}
	}
	return items
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func breadcrumbPop(crumb []string) []string {
	if len(crumb) <= 1 {
		return crumb
	}
	out := make([]string, len(crumb)-1)
	copy(out, crumb)
	return out
}

// =============================================================================
// ENTRY POINT
// =============================================================================

func runTUI() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro:", err)
		os.Exit(1)
	}
	srv, err := cfg.activeServer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	playOpts := PlayOpts{
		Player:     srv.Player,
		HWDec:      srv.HWDec,
		Fullscreen: srv.Fullscreen,
	}
	p := tea.NewProgram(
		newAppModel(NewClient(srv.URL, srv.Username, srv.Password), cfg, playOpts),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro:", err)
		os.Exit(1)
	}
}