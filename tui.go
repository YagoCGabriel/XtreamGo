package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7B61FF")).Padding(0, 1)
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	docStyle     = lipgloss.NewStyle().Margin(1, 2)

	settingKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#79c0ff")).Width(20)
	settingValStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	settingDimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
)

// ─── List item ────────────────────────────────────────────────────────────────

type menuItem struct {
	title, desc string
	id          string
	intID       int
	extra       any
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

// ─── Screens ──────────────────────────────────────────────────────────────────

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

// ─── Messages ─────────────────────────────────────────────────────────────────

type itemsLoadedMsg struct{ items []list.Item }
type loadErrMsg struct{ err error }
type seriesInfoLoadedMsg struct{ info *SeriesInfo }
type configSavedMsg struct{}

// ─── Model ────────────────────────────────────────────────────────────────────

type appModel struct {
	client           *Client
	cfg              *Config
	screen           screen
	list             list.Model
	status, errMsg   string
	playOpts         PlayOpts
	selectedCatID    string
	selectedSeriesID int
	seriesInfo       *SeriesInfo
}

func newAppModel(client *Client, cfg *Config, playOpts PlayOpts) appModel {
	del := list.NewDefaultDelegate()
	del.Styles.SelectedTitle = del.Styles.SelectedTitle.Foreground(lipgloss.Color("#7B61FF"))
	del.Styles.SelectedDesc = del.Styles.SelectedDesc.Foreground(lipgloss.Color("#9D8FFF"))
	l := list.New(nil, del, 80, 20)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Title = "XtreamGO"
	l.SetItems(mainMenuItems())
	return appModel{
		client:   client,
		cfg:      cfg,
		screen:   screenMain,
		list:     l,
		playOpts: playOpts,
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

func (m appModel) Init() tea.Cmd { return nil }

func (m appModel) navigate(s screen, title string, loader tea.Cmd) (tea.Model, tea.Cmd) {
	m.screen = s
	m.list.Title = title
	m.status = "Carregando..."
	m.list.SetItems(nil)
	return m, loader
}

// ─── Update ───────────────────────────────────────────────────────────────────

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		// settings tem teclas próprias
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
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-4)

	case itemsLoadedMsg:
		m.list.SetItems(msg.items)
		m.status, m.errMsg = "", ""

	case loadErrMsg:
		m.errMsg = msg.err.Error()
		m.status = ""

	case configSavedMsg:
		m.status = successStyle.Render("Configuracoes salvas!")

	case seriesInfoLoadedMsg:
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
				desc:  fmt.Sprintf("%d ep.", len(msg.info.Seasons[s])),
				id:    s,
			}
		}
		m.list.SetItems(items)
		m.status = ""
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// ─── Settings screen ──────────────────────────────────────────────────────────

func (m appModel) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	srv := &m.cfg.Servers[m.cfg.Current]
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc", "backspace":
		m.screen = screenMain
		m.list.Title = "XtreamGO"
		m.list.SetItems(mainMenuItems())
		m.status, m.errMsg = "", ""
	// player
	case "p":
    options := []string{"mpv", "vlc", "kmplayer"}
    current := srv.Player
    if current == "" {
        current = "mpv"
    }
    next := "mpv"
    for i, o := range options {
        if o == current {
            next = options[(i+1)%len(options)]
            break
        }
    }
    srv.Player = next
    m.playOpts.Player = next
    return m, m.saveConfig()
	// hwdec
	case "h":
		options := []string{"no", "auto", "vaapi", "vdpau", "nvdec", "videotoolbox"}
		current := srv.HWDec
		if current == "" {
			current = "no"
		}
		next := "no"
		for i, o := range options {
			if o == current {
				next = options[(i+1)%len(options)]
				break
			}
		}
		srv.HWDec = next
		m.playOpts.HWDec = next
		return m, m.saveConfig()
	// fullscreen
	case "f":
		m.playOpts.Fullscreen = !m.playOpts.Fullscreen
		srv.Fullscreen = m.playOpts.Fullscreen
		return m, m.saveConfig()
	// trocar servidor
	case "s":
		m.screen = screenServerSelect
		m.list.Title = "Configuracoes > Servidor"
		m.list.SetFilteringEnabled(false)
		items := make([]list.Item, len(m.cfg.Servers))
		for i, s := range m.cfg.Servers {
			mark := ""
			if i == m.cfg.Current {
				mark = " (ativo)"
			}
			items[i] = menuItem{
				title: s.Name + mark,
				desc:  s.URL,
				id:    fmt.Sprintf("%d", i),
			}
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
		return m, m.saveConfig()
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m appModel) openSettings() (tea.Model, tea.Cmd) {
	m.screen = screenSettings
	m.list.SetFilteringEnabled(false)
	m.list.Title = "Configuracoes"
	m.list.SetItems(nil) // settings tem view própria
	m.errMsg = ""
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

// ─── View ─────────────────────────────────────────────────────────────────────

func (m appModel) View() string {
	if m.screen == screenSettings {
		return m.viewSettings()
	}

	var sb strings.Builder
	sb.WriteString(docStyle.Render(m.list.View()))
	sb.WriteString("\n")
	if m.status != "" {
		sb.WriteString("  " + m.status + "\n")
	}
	if m.errMsg != "" {
		sb.WriteString("  " + errorStyle.Render("Erro: "+m.errMsg) + "\n")
	}
	footer := "  enter=selecionar  esc=voltar  /=filtrar  q=sair"
	if m.screen == screenServerSelect {
		footer = "  enter=selecionar  esc=voltar  q=sair"
	}
	sb.WriteString(statusStyle.Render(footer))
	return sb.String()
}

func (m appModel) viewSettings() string {
	srv := m.cfg.Servers[m.cfg.Current]

	player := srv.Player
	if player == "" {
		player = "mpv"
	}
	hwdec := srv.HWDec
	if hwdec == "" {
		hwdec = "no"
	}
	fs := "nao"
	if m.playOpts.Fullscreen {
		fs = "sim"
	}

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Configuracoes") + "\n\n")

	row := func(key, val, hint string) {
		sb.WriteString("  ")
		sb.WriteString(settingKeyStyle.Render(key))
		sb.WriteString(settingValStyle.Render(val))
		sb.WriteString("  " + settingDimStyle.Render(hint))
		sb.WriteString("\n")
	}

	sep := func(label string) {
		sb.WriteString("\n  " + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B61FF")).Bold(true).
			Render("── "+label) + "\n\n")
	}

	sep("Servidor ativo")
	row("Nome", srv.Name, "")
	row("URL", srv.URL, "")
	row("Usuario", srv.Username, "")

	sep("Reproducao")
	row("Player        [p]", player, "alternar mpv / vlc / kmplayer")
	row("HW Decode     [h]", hwdec, "no / auto / vaapi / vdpau / nvdec / videotoolbox")
	row("Fullscreen    [f]", fs, "abrir em tela cheia")

	sep("Servidores")
	row("Trocar        [s]", fmt.Sprintf("%d configurado(s)", len(m.cfg.Servers)), "selecionar servidor ativo")

	sb.WriteString("\n")
	if m.status != "" {
		sb.WriteString("  " + m.status + "\n")
	}
	if m.errMsg != "" {
		sb.WriteString("  " + errorStyle.Render("Erro: "+m.errMsg) + "\n")
	}
	sb.WriteString(statusStyle.Render("  p=player  h=hwdec  f=fullscreen  s=servidores  esc=voltar"))
	return sb.String()
}

// ─── handleEnter ──────────────────────────────────────────────────────────────

func (m appModel) handleEnter() (tea.Model, tea.Cmd) {
	sel, ok := m.list.SelectedItem().(menuItem)
	if !ok {
		return m, nil
	}
	m.errMsg = ""

	switch m.screen {
	case screenMain:
		switch sel.id {
		case "live":
			return m.navigate(screenLiveCats, "TV ao Vivo > Categorias", func() tea.Msg {
				cats, err := m.client.GetLiveCategories()
				if err != nil {
					return loadErrMsg{err}
				}
				return itemsLoadedMsg{catsToItems(cats)}
			})
		case "vod":
			return m.navigate(screenVODCats, "Filmes > Categorias", func() tea.Msg {
				cats, err := m.client.GetVODCategories()
				if err != nil {
					return loadErrMsg{err}
				}
				return itemsLoadedMsg{catsToItems(cats)}
			})
		case "series":
			return m.navigate(screenSeriesCats, "Series > Categorias", func() tea.Msg {
				cats, err := m.client.GetSeriesCategories()
				if err != nil {
					return loadErrMsg{err}
				}
				return itemsLoadedMsg{catsToItems(cats)}
			})
		case "settings":
			return m.openSettings()
		}

	case screenLiveCats:
		m.selectedCatID = sel.id
		catID := sel.id
		return m.navigate(screenLiveStreams, "TV ao Vivo > "+sel.title, func() tea.Msg {
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

	case screenLiveStreams:
		s := sel.extra.(LiveStream)
		opts := m.playOpts
		opts.Title = s.Name
		go Play(m.client.LiveStreamURL(s.ID, s.ContainerExtension), opts)
		m.status = successStyle.Render("Abrindo: " + s.Name)

	case screenVODCats:
		m.selectedCatID = sel.id
		catID := sel.id
		return m.navigate(screenVODStreams, "Filmes > "+sel.title, func() tea.Msg {
			streams, err := m.client.GetVODStreams(catID)
			if err != nil {
				return loadErrMsg{err}
			}
			items := make([]list.Item, len(streams))
			for i, s := range streams {
				desc := ""
				if s.Rating.String() != "" {
					desc = "Nota: " + s.Rating.String()
				}
				items[i] = menuItem{title: s.Name, desc: desc, id: fmt.Sprintf("%d", s.ID), extra: s}
			}
			return itemsLoadedMsg{items}
		})

	case screenVODStreams:
		v := sel.extra.(VODStream)
		opts := m.playOpts
		opts.Title = v.Name
		go Play(m.client.VODStreamURL(v.ID, v.ContainerExtension), opts)
		m.status = successStyle.Render("Abrindo: " + v.Name)

	case screenSeriesCats:
		m.selectedCatID = sel.id
		catID := sel.id
		return m.navigate(screenSeriesList, "Series > "+sel.title, func() tea.Msg {
			series, err := m.client.GetSeries(catID)
			if err != nil {
				return loadErrMsg{err}
			}
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
			return itemsLoadedMsg{items}
		})

	case screenSeriesList:
		s := sel.extra.(Series)
		m.selectedSeriesID = s.ID
		sid := s.ID
		return m.navigate(screenSeasons, "Series > "+s.Name+" > Temporadas", func() tea.Msg {
			info, err := m.client.GetSeriesInfo(sid)
			if err != nil {
				return loadErrMsg{err}
			}
			return seriesInfoLoadedMsg{info}
		})

	case screenSeasons:
		m.screen = screenEpisodes
		parts := strings.SplitN(m.list.Title, " > Temporadas", 2)
		m.list.Title = parts[0] + " > " + sel.title
		episodes := m.seriesInfo.Seasons[sel.id]
		sort.Slice(episodes, func(i, j int) bool {
			return episodes[i].EpisodeNum.Int() < episodes[j].EpisodeNum.Int()
		})
		items := make([]list.Item, len(episodes))
		for i, ep := range episodes {
			t := fmt.Sprintf("E%02d", ep.EpisodeNum.Int())
			if ep.Title != "" {
				t += " - " + ep.Title
			}
			items[i] = menuItem{title: t, id: ep.ID, extra: ep}
		}
		m.list.SetItems(items)

	case screenEpisodes:
		ep := sel.extra.(Episode)
		opts := m.playOpts
		opts.Title = sel.title
		go Play(m.client.SeriesStreamURL(ep.ID, ep.ContainerExtension), opts)
		m.status = successStyle.Render("Abrindo: " + sel.title)
	}

	return m, nil
}

// ─── goBack ───────────────────────────────────────────────────────────────────

func (m appModel) goBack() (tea.Model, tea.Cmd) {
	m.errMsg, m.status = "", ""
	m.list.SetFilteringEnabled(true)
	switch m.screen {
	case screenMain:
		return m, tea.Quit
	case screenLiveCats, screenVODCats, screenSeriesCats, screenSettings:
		m.screen = screenMain
		m.list.Title = "XtreamGO"
		m.list.SetItems(mainMenuItems())
	case screenLiveStreams:
		return m.navigate(screenLiveCats, "TV ao Vivo > Categorias", func() tea.Msg {
			cats, err := m.client.GetLiveCategories()
			if err != nil {
				return loadErrMsg{err}
			}
			return itemsLoadedMsg{catsToItems(cats)}
		})
	case screenVODStreams:
		return m.navigate(screenVODCats, "Filmes > Categorias", func() tea.Msg {
			cats, err := m.client.GetVODCategories()
			if err != nil {
				return loadErrMsg{err}
			}
			return itemsLoadedMsg{catsToItems(cats)}
		})
	case screenSeriesList:
		return m.navigate(screenSeriesCats, "Series > Categorias", func() tea.Msg {
			cats, err := m.client.GetSeriesCategories()
			if err != nil {
				return loadErrMsg{err}
			}
			return itemsLoadedMsg{catsToItems(cats)}
		})
	case screenSeasons:
		catID := m.selectedCatID
		return m.navigate(screenSeriesList, "Series", func() tea.Msg {
			series, err := m.client.GetSeries(catID)
			if err != nil {
				return loadErrMsg{err}
			}
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
			return itemsLoadedMsg{items}
		})
	case screenEpisodes:
		sid := m.selectedSeriesID
		return m.navigate(screenSeasons, "Temporadas", func() tea.Msg {
			info, err := m.client.GetSeriesInfo(sid)
			if err != nil {
				return loadErrMsg{err}
			}
			return seriesInfoLoadedMsg{info}
		})
	}
	return m, nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func catsToItems(cats []Category) []list.Item {
	items := make([]list.Item, len(cats))
	for i, c := range cats {
		items[i] = menuItem{title: c.Name, id: c.ID}
	}
	return items
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

// ─── Entry point ──────────────────────────────────────────────────────────────

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
