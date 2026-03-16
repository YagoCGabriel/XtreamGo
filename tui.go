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
// PALETA E ESTILOS
// =============================================================================

var (
	colorAccent  = lipgloss.Color("#C084FC")
	colorAccent2 = lipgloss.Color("#818CF8")
	colorCyan    = lipgloss.Color("#22D3EE")
	colorGreen   = lipgloss.Color("#4ADE80")
	colorRed     = lipgloss.Color("#F87171")
	colorAmber   = lipgloss.Color("#FBBF24")
	colorText    = lipgloss.Color("#F1F5F9")
	colorSub     = lipgloss.Color("#94A3B8")
	colorMuted   = lipgloss.Color("#475569")
	colorFaint   = lipgloss.Color("#334155")
	colorBgMid   = lipgloss.Color("#1E293B")
	colorBgLight = lipgloss.Color("#253347")
	colorBorder  = lipgloss.Color("#334155")

	subStyle     = lipgloss.NewStyle().Foreground(colorSub)
	mutedStyle   = lipgloss.NewStyle().Foreground(colorMuted)
	faintStyle   = lipgloss.NewStyle().Foreground(colorFaint)
	accentStyle  = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	accent2Style = lipgloss.NewStyle().Foreground(colorAccent2)
	cyanStyle    = lipgloss.NewStyle().Foreground(colorCyan)
	amberStyle   = lipgloss.NewStyle().Foreground(colorAmber)

	successStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(colorRed).Bold(true)

	headerLogoStyle    = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	headerVersionStyle = lipgloss.NewStyle().Foreground(colorMuted)

	bcSepStyle = lipgloss.NewStyle().Foreground(colorFaint)
	bcStyle    = lipgloss.NewStyle().Foreground(colorMuted)
	bcCurStyle = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)

	footerKeyStyle  = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	footerDescStyle = lipgloss.NewStyle().Foreground(colorMuted)
	footerDotStyle  = lipgloss.NewStyle().Foreground(colorFaint)

	detailTitleStyle = lipgloss.NewStyle().Foreground(colorText).Bold(true)
	detailLabelStyle = lipgloss.NewStyle().Foreground(colorMuted).Width(16)
	detailValueStyle = lipgloss.NewStyle().Foreground(colorCyan)
	detailPlotStyle  = lipgloss.NewStyle().Foreground(colorSub).Italic(true)

	detailBadgeLive = lipgloss.NewStyle().
			Foreground(colorBgMid).Background(colorRed).Bold(true).Padding(0, 1)
	detailBadgeWatched = lipgloss.NewStyle().
				Foreground(colorBgMid).Background(colorAmber).Bold(true).Padding(0, 1)
	detailBadgeNew = lipgloss.NewStyle().
			Foreground(colorBgMid).Background(colorGreen).Bold(true).Padding(0, 1)

	settingLabelStyle   = lipgloss.NewStyle().Foreground(colorSub).Width(20)
	settingValueStyle   = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
	settingHintStyle    = lipgloss.NewStyle().Foreground(colorMuted).Italic(true)
	settingKeyStyle     = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	settingSectionStyle = lipgloss.NewStyle().
				Foreground(colorAccent2).Bold(true).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(colorFaint)
)

// =============================================================================
// DELEGATE
// =============================================================================

type customDelegate struct{}

func newDelegate() customDelegate { return customDelegate{} }

func (d customDelegate) Height() int                               { return 2 }
func (d customDelegate) Spacing() int                             { return 0 }
func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d customDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(menuItem)
	if !ok {
		return
	}
	if i.id == "separator" {
		fmt.Fprintf(w, " %s\n ", faintStyle.Render(strings.Repeat("─", m.Width()-4)))
		return
	}

	sel := index == m.Index()
	width := m.Width() - 2
	if width < 10 {
		width = 10
	}
	icon := itemIcon(i.id, i.extra)
	watched := i.desc == "★ assistido" || i.desc == "★ watched"

	if i.id == "search_all" {
		if sel {
			fmt.Fprintf(w, "%s\n%s",
				lipgloss.NewStyle().Foreground(colorAmber).Bold(true).Background(colorBgLight).
					Width(width).PaddingLeft(1).Render("▸ "+icon+i.title),
				lipgloss.NewStyle().Foreground(colorSub).Background(colorBgLight).
					Width(width).PaddingLeft(3).Render(i.desc),
			)
		} else {
			fmt.Fprintf(w, "%s\n%s",
				amberStyle.PaddingLeft(2).Width(width).Render(icon+i.title),
				faintStyle.PaddingLeft(4).Width(width).Render(i.desc),
			)
		}
		return
	}

	if sel {
		titleFg := colorText
		extra := ""
		if watched {
			titleFg = colorAmber
			extra = "  ★"
		}
		descText := i.desc
		if watched {
			descText = "✓"
		}
		fmt.Fprintf(w, "%s\n%s",
			lipgloss.NewStyle().Foreground(titleFg).Bold(true).Background(colorBgLight).
				Width(width).PaddingLeft(1).Render("▸ "+icon+i.title+extra),
			lipgloss.NewStyle().Foreground(colorSub).Background(colorBgLight).
				Width(width).PaddingLeft(3).Render(descText),
		)
	} else {
		titleFg := colorSub
		extra := ""
		if watched {
			titleFg = colorFaint
			extra = " ★"
		}
		fmt.Fprintf(w, "%s\n%s",
			lipgloss.NewStyle().Foreground(titleFg).Width(width).PaddingLeft(2).Render(icon+i.title+extra),
			faintStyle.Width(width).PaddingLeft(4).Render(i.desc),
		)
	}
}

func itemIcon(id string, extra any) string {
	switch id {
	case "live":
		return "  "
	case "vod":
		return "  "
	case "series":
		return "  "
	case "settings":
		return "  "
	case "search_all":
		return "  "
	}
	switch extra.(type) {
	case LiveStream:
		return "  "
	case VODStream:
		return "  "
	case Series:
		return "  "
	case Episode:
		return "  "
	}
	if strings.HasPrefix(id, "T") {
		return "  "
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

func searchAllItem(t Lang, desc string) menuItem {
	return menuItem{title: t.SearchAll, desc: desc, id: "search_all"}
}
func separatorItem() menuItem { return menuItem{id: "separator"} }

func prependSearchAll(t Lang, desc string, items []list.Item) []list.Item {
	out := make([]list.Item, 0, len(items)+2)
	out = append(out, searchAllItem(t, desc), separatorItem())
	return append(out, items...)
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
type watchedSavedMsg struct{}

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
	selectedSeason   string
	seriesInfo       *SeriesInfo
	watched          *WatchedDB
	width, height    int
	breadcrumb       []string
	showDetail       bool
}

const minSplitWidth = 110

func newAppModel(client *Client, cfg *Config, playOpts PlayOpts) appModel {
	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = lipgloss.NewStyle().Foreground(colorAccent)

	l := list.New(nil, newDelegate(), 80, 20)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(colorAccent)
	l.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(colorText)
	l.Title = ""
	l.Styles.Title = lipgloss.NewStyle()

	watched, err := loadWatched()
	if err != nil {
		watched = &WatchedDB{Episodes: make(map[string]WatchedEntry)}
	}

	m := appModel{
		client: client, cfg: cfg, screen: screenMain,
		list: l, spinner: sp, playOpts: playOpts,
		watched: watched, showDetail: true,
		breadcrumb: []string{T(cfg).Home},
	}
	l.SetItems(m.mainMenuItems())
	m.list = l
	return m
}

func (m appModel) t() Lang { return T(m.cfg) }

func (m appModel) mainMenuItems() []list.Item {
	t := m.t()
	return []list.Item{
		menuItem{title: t.MenuLive, desc: t.MenuLiveDesc, id: "live"},
		menuItem{title: t.MenuVOD, desc: t.MenuVODDesc, id: "vod"},
		menuItem{title: t.MenuSeries, desc: t.MenuSerDesc, id: "series"},
		menuItem{title: t.MenuSettings, desc: t.MenuSetDesc, id: "settings"},
	}
}

func (m appModel) Init() tea.Cmd { return m.spinner.Tick }

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
		case "w":
			if m.screen == screenEpisodes {
				return m.toggleWatched()
			}
		case "tab", "i":
			m.showDetail = !m.showDetail
			m.resizeList()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeList()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case itemsLoadedMsg:
		m.loading = false
		m.list.SetItems(msg.items)
		m.status, m.errMsg = "", ""

	case loadErrMsg:
		m.loading = false
		m.errMsg = msg.err.Error()
		m.status = ""

	case configSavedMsg:
		m.status = m.t().StatusSaved

	case watchedSavedMsg:
		// silencioso

	case seriesInfoLoadedMsg:
		m.loading = false
		m.seriesInfo = msg.info
		m.list.SetItems(m.buildSeasonItems(msg.info))
		m.status = ""
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, tea.Batch(cmd, m.spinner.Tick)
}

func (m *appModel) resizeList() {
	w, h := m.width, m.height
	if w == 0 { w = 80 }
	if h == 0 { h = 24 }
	listH := h - 9
	if listH < 5 { listH = 5 }
	m.list.SetSize(m.listWidth()-2, listH)
}

func (m appModel) useSplit() bool {
	return m.showDetail && m.width >= minSplitWidth &&
		m.screen != screenSettings && m.screen != screenServerSelect
}

func (m appModel) listWidth() int {
	if !m.useSplit() {
		return m.width - 2
	}
	lw := int(float64(m.width) * 0.42)
	if lw < 38 { lw = 38 }
	if lw > 62 { lw = 62 }
	return lw
}

func (m appModel) detailWidth() int { return m.width - m.listWidth() - 3 }

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
		srv.Player = cycleNext([]string{"mpv", "vlc", "kmplayer"}, srv.Player, "mpv")
		m.playOpts.Player = srv.Player
		return m, m.saveConfig()
	case "h":
		srv.HWDec = cycleNext([]string{"no", "auto", "vaapi", "vdpau", "nvdec", "videotoolbox"}, srv.HWDec, "no")
		m.playOpts.HWDec = srv.HWDec
		return m, m.saveConfig()
	case "f":
		m.playOpts.Fullscreen = !m.playOpts.Fullscreen
		srv.Fullscreen = m.playOpts.Fullscreen
		return m, m.saveConfig()
	case "l":
		if m.cfg.Language == "pt-BR" {
			m.cfg.Language = "en-US"
		} else {
			m.cfg.Language = "pt-BR"
		}
		m.list.SetItems(m.mainMenuItems())
		return m, m.saveConfig()
	case "s":
		t := m.t()
		m.screen = screenServerSelect
		m.breadcrumb = []string{t.Home, t.Settings, t.Servers}
		m.list.SetFilteringEnabled(false)
		items := make([]list.Item, len(m.cfg.Servers))
		for i, s := range m.cfg.Servers {
			desc := s.URL
			if i == m.cfg.Current {
				desc = "● active  ·  " + s.URL
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
		if !ok { break }
		var idx int
		fmt.Sscanf(sel.id, "%d", &idx)
		m.cfg.Current = idx
		srv := m.cfg.Servers[idx]
		m.client = NewClient(srv.URL, srv.Username, srv.Password)
		m.playOpts = PlayOpts{Player: srv.Player, HWDec: srv.HWDec, Fullscreen: srv.Fullscreen}
		m.status = m.t().StatusServer
		return m, m.saveConfig()
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m appModel) openSettings() (tea.Model, tea.Cmd) {
	t := m.t()
	m.screen = screenSettings
	m.list.SetFilteringEnabled(false)
	m.list.SetItems(nil)
	m.errMsg = ""
	m.breadcrumb = []string{t.Home, t.Settings}
	return m, nil
}

func (m appModel) saveConfig() tea.Cmd {
	cfg := m.cfg
	return func() tea.Msg {
		if err := saveConfig(cfg); err != nil { return loadErrMsg{err} }
		return configSavedMsg{}
	}
}

// =============================================================================
// WATCHED
// =============================================================================

func (m appModel) toggleWatched() (tea.Model, tea.Cmd) {
	sel, ok := m.list.SelectedItem().(menuItem)
	if !ok { return m, nil }
	ep, ok := sel.extra.(Episode)
	if !ok { return m, nil }

	marked := m.watched.Toggle(ep.ID, sel.title)
	db := m.watched
	m = m.rebuildEpisodeList()
	t := m.t()
	if marked {
		m.status = t.StatusWatched
	} else {
		m.status = t.StatusUnwatched
	}
	return m, func() tea.Msg { saveWatched(db); return watchedSavedMsg{} }
}

func (m appModel) watchedDesc() string {
	if m.cfg.Language == "en-US" {
		return "★ watched"
	}
	return "★ assistido"
}

func (m appModel) rebuildEpisodeList() appModel {
	if m.seriesInfo == nil { return m }
	episodes, ok := m.seriesInfo.Seasons[m.selectedSeason]
	if !ok { return m }
	m.list.SetItems(m.buildEpisodeItems(episodes))
	return m
}

func (m appModel) buildEpisodeItems(episodes []Episode) []list.Item {
	sorted := make([]Episode, len(episodes))
	copy(sorted, episodes)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].EpisodeNum.Int() < sorted[j].EpisodeNum.Int() })
	wd := m.watchedDesc()
	items := make([]list.Item, len(sorted))
	for i, ep := range sorted {
		label := fmt.Sprintf("E%02d", ep.EpisodeNum.Int())
		if ep.Title != "" { label += "  " + ep.Title }
		desc := ""
		if m.watched.IsWatched(ep.ID) { desc = wd }
		items[i] = menuItem{title: label, id: ep.ID, extra: ep, desc: desc}
	}
	return items
}

func (m appModel) buildSeasonItems(info *SeriesInfo) []list.Item {
	t := m.t()
	seasons := make([]string, 0, len(info.Seasons))
	for k := range info.Seasons { seasons = append(seasons, k) }
	sort.Strings(seasons)
	items := make([]list.Item, len(seasons))
	for i, s := range seasons {
		eps := info.Seasons[s]
		wc := 0
		for _, ep := range eps {
			if m.watched.IsWatched(ep.ID) { wc++ }
		}
		desc := fmt.Sprintf(t.EpCount, len(eps))
		if wc > 0 { desc += "  " + fmt.Sprintf(t.WatchedCount, wc, len(eps)) }
		items[i] = menuItem{title: t.Season + " " + s, desc: desc, id: "T" + s, extra: s}
	}
	return items
}

// =============================================================================
// VIEW
// =============================================================================

func (m appModel) View() string {
	w, h := m.width, m.height
	if w == 0 { w = 80 }
	if h == 0 { h = 24 }

	header    := m.viewHeader(w)
	breadcrumb := m.viewBreadcrumb(w)
	footer    := m.viewFooter(w)
	bodyH     := h - lipgloss.Height(header) - lipgloss.Height(breadcrumb) - lipgloss.Height(footer)
	if bodyH < 3 { bodyH = 3 }

	var body string
	switch {
	case m.screen == screenSettings:
		body = m.viewSettings(w, bodyH)
	case m.useSplit():
		body = m.viewSplit(w, bodyH)
	default:
		body = lipgloss.NewStyle().Height(bodyH).Render(m.list.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, breadcrumb, body, footer)
}

func (m appModel) viewSplit(w, h int) string {
	lw := m.listWidth()
	dw := m.detailWidth()

	leftBox := lipgloss.NewStyle().Width(lw).Height(h).MaxHeight(h).Render(m.list.View())

	divLines := make([]string, h)
	for i := range divLines { divLines[i] = faintStyle.Render("│") }
	divider := strings.Join(divLines, "\n")

	rightBox := lipgloss.NewStyle().Width(dw).Height(h).MaxHeight(h).PaddingLeft(1).
		Render(m.viewDetail(dw-3, h))

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, divider, rightBox)
}

// =============================================================================
// PAINEL DE DETALHES
// =============================================================================

func (m appModel) viewDetail(w, h int) string {
	sel := m.list.SelectedItem()
	var lines []string
	switch m.screen {
	case screenMain:
		lines = m.detailMain(w)
	case screenLiveCats, screenVODCats, screenSeriesCats:
		lines = m.detailCategory(w, sel)
	case screenLiveStreams:
		lines = m.detailLiveStream(w, sel)
	case screenVODStreams:
		lines = m.detailVOD(w, sel)
	case screenSeriesList:
		lines = m.detailSeries(w, sel)
	case screenSeasons:
		lines = m.detailSeason(w, sel)
	case screenEpisodes:
		lines = m.detailEpisode(w, sel)
	default:
		lines = []string{"", mutedStyle.Render("  " + m.t().NoItems)}
	}
	return lipgloss.NewStyle().Width(w).Height(h).MaxHeight(h).Render(strings.Join(lines, "\n"))
}

func (m appModel) detailMain(w int) []string {
	t := m.t()
	lines := []string{"", accentStyle.Render("  XtreamGo"), faintStyle.Render("  "+strings.Repeat("─", minInt(w-4, 22))), ""}
	if len(m.cfg.Servers) > 0 {
		srv := m.cfg.Servers[m.cfg.Current]
		player := srv.Player
		if player == "" { player = "mpv" }
		lines = append(lines,
			detailRow(t.Server, srv.Name, w),
			detailRow(t.URL, ellipsis(srv.URL, w-20), w),
			detailRow(t.User, srv.Username, w),
			detailRow(t.Player, player, w),
		)
	}
	lines = append(lines, "", faintStyle.Render("  "+strings.Repeat("─", minInt(w-4, 22))), "")
	total := len(m.watched.Episodes)
	lines = append(lines, "  "+amberStyle.Render(fmt.Sprintf(t.SetWatched, total)))
	lines = append(lines, "", "  "+faintStyle.Render(t.TogglePanel))
	return lines
}

func (m appModel) detailCategory(w int, sel list.Item) []string {
	i, ok := sel.(menuItem)
	if !ok || i.id == "separator" || i.id == "search_all" {
		return []string{"", mutedStyle.Render("  " + m.t().SelectCat)}
	}
	return []string{"", "  " + detailTitleStyle.Render(i.title),
		"  " + faintStyle.Render(strings.Repeat("─", minInt(len(i.title)+2, w-4)))}
}

func (m appModel) detailLiveStream(w int, sel list.Item) []string {
	i, ok := sel.(menuItem)
	if !ok { return nil }
	s, ok := i.extra.(LiveStream)
	if !ok { return nil }
	t := m.t()
	player := m.playOpts.Player
	if player == "" { player = "mpv" }
	lines := []string{"", "  " + detailBadgeLive.Render(t.LiveBadge), "",
		"  " + detailTitleStyle.Render(wordWrapFirst(s.Name, w-4)),
		"  " + faintStyle.Render(strings.Repeat("─", minInt(w-4, 30))), "",
		detailRow("ID", fmt.Sprintf("%d", s.ID), w),
	}
	if s.ContainerExtension != "" {
		lines = append(lines, detailRow(t.Format, s.ContainerExtension, w))
	}
	lines = append(lines, "", "  "+mutedStyle.Render(fmt.Sprintf(t.OpenIn, player)))
	return lines
}

func (m appModel) detailVOD(w int, sel list.Item) []string {
	i, ok := sel.(menuItem)
	if !ok { return nil }
	v, ok := i.extra.(VODStream)
	if !ok { return nil }
	t := m.t()
	lines := []string{"", "  " + detailTitleStyle.Render(wordWrapFirst(v.Name, w-4)),
		"  " + faintStyle.Render(strings.Repeat("─", minInt(w-4, 36))), ""}
	if r := v.Rating.String(); r != "" && r != "0" {
		lines = append(lines, "  "+amberStyle.Render(ratingStars(r))+"  "+subStyle.Render(r), "")
	}
	if v.Plot != "" {
		for j, line := range wordWrap(v.Plot, w-4) {
			if j >= 6 { lines = append(lines, "  "+faintStyle.Render("…")); break }
			lines = append(lines, "  "+detailPlotStyle.Render(line))
		}
		lines = append(lines, "")
	}
	if v.ContainerExtension != "" {
		lines = append(lines, detailRow(t.Format, v.ContainerExtension, w))
	}
	lines = append(lines, "", "  "+mutedStyle.Render("Enter → "+m.playOpts.Player))
	return lines
}

func (m appModel) detailSeries(w int, sel list.Item) []string {
	i, ok := sel.(menuItem)
	if !ok || i.id == "separator" || i.id == "search_all" {
		return []string{"", mutedStyle.Render("  " + m.t().SelectSeries)}
	}
	s, ok := i.extra.(Series)
	if !ok { return nil }
	lines := []string{"", "  " + detailTitleStyle.Render(wordWrapFirst(s.Name, w-4)),
		"  " + faintStyle.Render(strings.Repeat("─", minInt(w-4, 36))), ""}
	if s.Rating != "" && s.Rating != "0" {
		lines = append(lines, "  "+amberStyle.Render(ratingStars(s.Rating))+"  "+subStyle.Render(s.Rating), "")
	}
	if s.Plot != "" {
		for j, line := range wordWrap(s.Plot, w-4) {
			if j >= 5 { lines = append(lines, "  "+faintStyle.Render("…")); break }
			lines = append(lines, "  "+detailPlotStyle.Render(line))
		}
		lines = append(lines, "")
	}
	lines = append(lines, "  "+mutedStyle.Render("Enter → "+m.t().Seasons))
	return lines
}

func (m appModel) detailSeason(w int, sel list.Item) []string {
	i, ok := sel.(menuItem)
	if !ok { return nil }
	seasonKey, ok := i.extra.(string)
	if !ok || m.seriesInfo == nil { return nil }
	t := m.t()
	eps := m.seriesInfo.Seasons[seasonKey]
	total := len(eps)
	wc := 0
	for _, ep := range eps {
		if m.watched.IsWatched(ep.ID) { wc++ }
	}
	lines := []string{"", "  " + detailTitleStyle.Render(i.title),
		"  " + faintStyle.Render(strings.Repeat("─", minInt(w-4, 28))), "",
		detailRow(t.Episodes2, fmt.Sprintf("%d", total), w),
		detailRow(t.Watched, fmt.Sprintf(t.WatchedOf, wc, total), w), "",
	}
	if total > 0 {
		pct := int(float64(wc) / float64(total) * 100)
		lines = append(lines,
			"  "+progressBar(wc, total, w-6),
			"  "+mutedStyle.Render(fmt.Sprintf(t.Completed, pct)),
			"",
		)
	}
	next := nextUnwatched(eps, m.watched, 3)
	if len(next) > 0 {
		lines = append(lines, "  "+faintStyle.Render(t.NextUnwatched))
		for _, ep := range next {
			label := fmt.Sprintf("  E%02d", ep.EpisodeNum.Int())
			if ep.Title != "" { label += "  " + ep.Title }
			lines = append(lines, "  "+subStyle.Render(label))
		}
	}
	return lines
}

func (m appModel) detailEpisode(w int, sel list.Item) []string {
	i, ok := sel.(menuItem)
	if !ok { return nil }
	ep, ok := i.extra.(Episode)
	if !ok { return nil }
	t := m.t()
	watched := m.watched.IsWatched(ep.ID)

	badge := detailBadgeNew.Render(t.Unwatched)
	if watched { badge = detailBadgeWatched.Render(t.Watched) }

	lines := []string{"", "  " + badge, "",
		"  " + faintStyle.Render(fmt.Sprintf(t.Episode, ep.EpisodeNum.Int())),
	}
	if ep.Title != "" {
		lines = append(lines, "  "+detailTitleStyle.Render(wordWrapFirst(ep.Title, w-4)))
	}
	lines = append(lines, "  "+faintStyle.Render(strings.Repeat("─", minInt(w-4, 30))), "")
	if ep.ContainerExtension != "" {
		lines = append(lines, detailRow(t.Format, ep.ContainerExtension, w), "")
	}
	if m.seriesInfo != nil {
		eps := m.seriesInfo.Seasons[m.selectedSeason]
		total := len(eps)
		wc := 0
		for _, e := range eps {
			if m.watched.IsWatched(e.ID) { wc++ }
		}
		if total > 0 {
			lines = append(lines,
				"  "+faintStyle.Render(t.Season+fmt.Sprintf("  %d/%d", wc, total)),
				"  "+progressBar(wc, total, w-6), "",
			)
		}
	}
	lines = append(lines, "  "+mutedStyle.Render(t.PlayMark), "  "+mutedStyle.Render(t.ToggleMark))
	return lines
}

// =============================================================================
// HEADER / BREADCRUMB / FOOTER
// =============================================================================

func (m appModel) viewHeader(w int) string {
	logo    := headerLogoStyle.Render("  XtreamGo")
	version := headerVersionStyle.Render(" v1.0")
	left    := logo + version

	right := ""
	if len(m.cfg.Servers) > 0 {
		srv    := m.cfg.Servers[m.cfg.Current]
		player := srv.Player
		if player == "" { player = "mpv" }
		lang   := m.cfg.Language

		// espaço disponível para o lado direito
		maxRight := w - lipgloss.Width(left) - 6
		if maxRight < 0 { maxRight = 0 }

		// completo: nome · player · lang
		full := faintStyle.Render(srv.Name) + faintStyle.Render("  ·  ") +
			accent2Style.Render(player) + faintStyle.Render("  ·  ") +
			faintStyle.Render(lang)

		if lipgloss.Width(full) <= maxRight {
			right = full
		} else {
			// médio: player · lang
			mid := accent2Style.Render(player) + faintStyle.Render("  ·  ") + faintStyle.Render(lang)
			if lipgloss.Width(mid) <= maxRight {
				right = mid
			} else {
				right = faintStyle.Render(lang)
			}
		}
	}

	gap := w - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 { gap = 1 }

	return lipgloss.NewStyle().
		Background(colorBgMid).
		BorderBottom(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(colorBorder).
		Width(w).Padding(0, 2).
		Render(left + strings.Repeat(" ", gap) + right)
}

func (m appModel) viewBreadcrumb(w int) string {
	parts := make([]string, len(m.breadcrumb))
	for i, b := range m.breadcrumb {
		if i == len(m.breadcrumb)-1 { parts[i] = bcCurStyle.Render(b) } else { parts[i] = bcStyle.Render(b) }
	}
	crumb := strings.Join(parts, bcSepStyle.Render("  ›  "))
	statusStr := ""
	if m.loading {
		statusStr = "  " + m.spinner.View() + " " + faintStyle.Render("…")
	} else if m.status != "" {
		statusStr = "  " + successStyle.Render("✓ "+m.status)
	} else if m.errMsg != "" {
		statusStr = "  " + errorStyle.Render("✗ "+m.errMsg)
	}
	splitHint := ""
	if m.width >= minSplitWidth && m.screen != screenSettings {
		icon := "⊠"
		if !m.showDetail { icon = "⊡" }
		splitHint = faintStyle.Render("tab " + icon)
	}
	left := crumb + statusStr
	gap := w - lipgloss.Width(left) - lipgloss.Width(splitHint) - 4
	if gap < 1 { gap = 1 }
	return lipgloss.NewStyle().Width(w).Padding(0, 2).Render(left + strings.Repeat(" ", gap) + splitHint)
}

func (m appModel) viewFooter(w int) string {
	t := m.t()
	type kp struct{ k, v string }
	var pairs []kp
	switch m.screen {
	case screenSettings:
		pairs = []kp{{"p", t.KeyPlayer}, {"h", t.KeyHWDec}, {"f", t.KeyFullscreen},
			{"l", "lang"}, {"s", t.KeyServers}, {"esc", t.KeyBack}}
	case screenServerSelect:
		pairs = []kp{{"enter", t.KeySelect}, {"esc", t.KeyBack}, {"q", t.KeyQuit}}
	case screenEpisodes:
		pairs = []kp{{"enter", t.KeyWatch}, {"w", t.KeyMark}, {"esc", t.KeyBack},
			{"/", t.KeyFilter}, {"q", t.KeyQuit}}
	default:
		pairs = []kp{{"enter", t.KeySelect}, {"esc", t.KeyBack},
			{"/", t.KeyFilter}, {"q", t.KeyQuit}}
	}
	parts := make([]string, len(pairs))
	for i, p := range pairs {
		parts[i] = footerKeyStyle.Render(p.k) + footerDescStyle.Render(" "+p.v)
	}
	return lipgloss.NewStyle().
		BorderTop(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(colorBorder).
		Width(w).Padding(0, 2).
		Render(strings.Join(parts, footerDotStyle.Render("  ·  ")))
}

// =============================================================================
// VIEW SETTINGS
// =============================================================================

func (m appModel) viewSettings(w, h int) string {
	t := m.t()
	if len(m.cfg.Servers) == 0 {
		return mutedStyle.Render("  " + t.NoneConfigured)
	}
	srv := m.cfg.Servers[m.cfg.Current]
	player := srv.Player
	if player == "" { player = "mpv" }
	hwdec := srv.HWDec
	if hwdec == "" { hwdec = "no" }
	fs := t.SetFullNo
	if m.playOpts.Fullscreen { fs = t.SetFullYes }
	lang := m.cfg.Language
	if lang == "" { lang = "pt-BR" }
	sw := w - 8
	if sw < 40 { sw = 40 }

	row := func(label, bind, value, hint string) string {
		return "   " + settingLabelStyle.Render(label) + settingKeyStyle.Render("["+bind+"]") +
			"  " + settingValueStyle.Render(value) + "   " + settingHintStyle.Render(hint)
	}
	section := func(title string) string {
		return "\n   " + settingSectionStyle.Width(sw).Render(title) + "\n"
	}

	var sb strings.Builder
	sb.WriteString(section(t.SetServerActive))
	sb.WriteString("\n")
	sb.WriteString("   " + settingLabelStyle.Render(t.Server) + cyanStyle.Render(srv.Name) + "\n")
	sb.WriteString("   " + settingLabelStyle.Render(t.URL) + subStyle.Render(ellipsis(srv.URL, w-30)) + "\n")
	sb.WriteString("   " + settingLabelStyle.Render(t.User) + subStyle.Render(srv.Username) + "\n")
	sb.WriteString(section(t.SetPlayback))
	sb.WriteString("\n")
	sb.WriteString(row("Player       ", "p", player, t.SetPlayerCycle) + "\n")
	sb.WriteString(row("HW Decode    ", "h", hwdec, t.SetHWCycle) + "\n")
	sb.WriteString(row("Fullscreen   ", "f", fs, t.SetFullscreen) + "\n")
	sb.WriteString(section(t.SetServers))
	sb.WriteString("\n")
	sb.WriteString(row("Trocar/Switch", "s", fmt.Sprintf(t.SetConfigured, len(m.cfg.Servers)), t.SetSwitch) + "\n")
	sb.WriteString(section("Language / Idioma"))
	sb.WriteString("\n")
	sb.WriteString(row("Language     ", "l", lang, "toggle: pt-BR ↔ en-US") + "\n")
	sb.WriteString(section(t.SetHistory))
	sb.WriteString("\n")
	sb.WriteString("   " + settingLabelStyle.Render(t.Watched) +
		amberStyle.Render(fmt.Sprintf(t.SetWatched, len(m.watched.Episodes))) + "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).BorderForeground(colorBorder).
		Width(w - 2).Padding(0, 1).
		Render(sb.String())
}

// =============================================================================
// HANDLE ENTER
// =============================================================================

func (m appModel) handleEnter() (tea.Model, tea.Cmd) {
	sel, ok := m.list.SelectedItem().(menuItem)
	if !ok || sel.id == "separator" { return m, nil }
	m.errMsg = ""
	t := m.t()

	switch m.screen {
	case screenMain:
		switch sel.id {
		case "live":
			return m.navigate(screenLiveCats, []string{t.Home, t.LiveTV}, func() tea.Msg {
				cats, err := m.client.GetLiveCategories()
				if err != nil { return loadErrMsg{err} }
				return itemsLoadedMsg{catsToItems(cats)}
			})
		case "vod":
			return m.navigate(screenVODCats, []string{t.Home, t.Movies}, func() tea.Msg {
				cats, err := m.client.GetVODCategories()
				if err != nil { return loadErrMsg{err} }
				return itemsLoadedMsg{prependSearchAll(t, t.SearchMovies, catsToItems(cats))}
			})
		case "series":
			return m.navigate(screenSeriesCats, []string{t.Home, t.Series}, func() tea.Msg {
				cats, err := m.client.GetSeriesCategories()
				if err != nil { return loadErrMsg{err} }
				return itemsLoadedMsg{prependSearchAll(t, t.SearchSeries, catsToItems(cats))}
			})
		case "settings":
			return m.openSettings()
		}

	case screenLiveCats:
		m.selectedCatID = sel.id
		catID := sel.id
		return m.navigate(screenLiveStreams, append(m.breadcrumb, sel.title), func() tea.Msg {
			streams, err := m.client.GetLiveStreams(catID)
			if err != nil { return loadErrMsg{err} }
			items := make([]list.Item, len(streams))
			for i, s := range streams {
				items[i] = menuItem{title: s.Name, id: fmt.Sprintf("%d", s.ID), extra: s}
			}
			return itemsLoadedMsg{items}
		})

	case screenLiveStreams:
		s := sel.extra.(LiveStream)
		opts := m.playOpts; opts.Title = s.Name
		go Play(m.client.LiveStreamURL(s.ID, s.ContainerExtension), opts)
		m.status = s.Name

	case screenVODCats:
		if sel.id == "search_all" {
			return m.navigate(screenVODStreams, []string{t.Home, t.Movies, t.All}, func() tea.Msg {
				streams, err := m.client.GetVODStreams("")
				if err != nil { return loadErrMsg{err} }
				return itemsLoadedMsg{vodStreamsToItems(streams)}
			})
		}
		m.selectedCatID = sel.id; catID := sel.id
		return m.navigate(screenVODStreams, append(m.breadcrumb, sel.title), func() tea.Msg {
			streams, err := m.client.GetVODStreams(catID)
			if err != nil { return loadErrMsg{err} }
			return itemsLoadedMsg{vodStreamsToItems(streams)}
		})

	case screenVODStreams:
		v := sel.extra.(VODStream)
		opts := m.playOpts; opts.Title = v.Name
		go Play(m.client.VODStreamURL(v.ID, v.ContainerExtension), opts)
		m.status = v.Name

	case screenSeriesCats:
		if sel.id == "search_all" {
			return m.navigate(screenSeriesList, []string{t.Home, t.Series, t.All}, func() tea.Msg {
				series, err := m.client.GetSeries("")
				if err != nil { return loadErrMsg{err} }
				return itemsLoadedMsg{seriesToItems(series)}
			})
		}
		m.selectedCatID = sel.id; catID := sel.id
		return m.navigate(screenSeriesList, append(m.breadcrumb, sel.title), func() tea.Msg {
			series, err := m.client.GetSeries(catID)
			if err != nil { return loadErrMsg{err} }
			return itemsLoadedMsg{seriesToItems(series)}
		})

	case screenSeriesList:
		s := sel.extra.(Series)
		m.selectedSeriesID = s.ID; sid := s.ID
		return m.navigate(screenSeasons, append(m.breadcrumb, s.Name), func() tea.Msg {
			info, err := m.client.GetSeriesInfo(sid)
			if err != nil { return loadErrMsg{err} }
			return seriesInfoLoadedMsg{info}
		})

	case screenSeasons:
		seasonKey := strings.TrimPrefix(sel.id, "T")
		m.selectedSeason = seasonKey
		m.screen = screenEpisodes
		m.breadcrumb = append(m.breadcrumb, sel.title)
		m.list.SetItems(m.buildEpisodeItems(m.seriesInfo.Seasons[seasonKey]))

	case screenEpisodes:
		ep := sel.extra.(Episode)
		opts := m.playOpts; opts.Title = sel.title
		go Play(m.client.SeriesStreamURL(ep.ID, ep.ContainerExtension), opts)
		m.watched.Mark(ep.ID, sel.title)
		db := m.watched
		m = m.rebuildEpisodeList()
		m.status = sel.title
		return m, func() tea.Msg { saveWatched(db); return watchedSavedMsg{} }
	}
	return m, nil
}

// =============================================================================
// NAVIGATE / GO BACK
// =============================================================================

func (m appModel) navigate(s screen, crumb []string, loader tea.Cmd) (tea.Model, tea.Cmd) {
	m.screen = s
	m.loading, m.status, m.errMsg = true, "", ""
	m.breadcrumb = crumb
	m.list.SetFilteringEnabled(true)
	m.list.SetItems(nil)
	return m, tea.Batch(loader, m.spinner.Tick)
}

func (m appModel) goHome() (tea.Model, tea.Cmd) {
	t := m.t()
	m.screen = screenMain
	m.breadcrumb = []string{t.Home}
	m.errMsg, m.status, m.loading = "", "", false
	m.list.SetFilteringEnabled(true)
	m.list.SetItems(m.mainMenuItems())
	return m, nil
}

func (m appModel) goBack() (tea.Model, tea.Cmd) {
	t := m.t()
	m.errMsg, m.status = "", ""
	m.list.SetFilteringEnabled(true)
	switch m.screen {
	case screenMain:
		return m, tea.Quit
	case screenLiveCats, screenVODCats, screenSeriesCats, screenSettings:
		return m.goHome()
	case screenLiveStreams:
		return m.navigate(screenLiveCats, []string{t.Home, t.LiveTV}, func() tea.Msg {
			cats, err := m.client.GetLiveCategories()
			if err != nil { return loadErrMsg{err} }
			return itemsLoadedMsg{catsToItems(cats)}
		})
	case screenVODStreams:
		return m.navigate(screenVODCats, []string{t.Home, t.Movies}, func() tea.Msg {
			cats, err := m.client.GetVODCategories()
			if err != nil { return loadErrMsg{err} }
			return itemsLoadedMsg{prependSearchAll(t, t.SearchMovies, catsToItems(cats))}
		})
	case screenSeriesList:
		return m.navigate(screenSeriesCats, []string{t.Home, t.Series}, func() tea.Msg {
			cats, err := m.client.GetSeriesCategories()
			if err != nil { return loadErrMsg{err} }
			return itemsLoadedMsg{prependSearchAll(t, t.SearchSeries, catsToItems(cats))}
		})
	case screenSeasons:
		catID := m.selectedCatID
		prev := breadcrumbPop(m.breadcrumb)
		return m.navigate(screenSeriesList, prev, func() tea.Msg {
			series, err := m.client.GetSeries(catID)
			if err != nil { return loadErrMsg{err} }
			return itemsLoadedMsg{seriesToItems(series)}
		})
	case screenEpisodes:
		sid := m.selectedSeriesID
		prev := breadcrumbPop(m.breadcrumb)
		return m.navigate(screenSeasons, prev, func() tea.Msg {
			info, err := m.client.GetSeriesInfo(sid)
			if err != nil { return loadErrMsg{err} }
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
	for i, c := range cats { items[i] = menuItem{title: c.Name, id: c.ID} }
	return items
}

func vodStreamsToItems(streams []VODStream) []list.Item {
	items := make([]list.Item, len(streams))
	for i, s := range streams {
		desc := ""
		if r := s.Rating.String(); r != "" && r != "0" { desc = "★ " + r }
		items[i] = menuItem{title: s.Name, desc: desc, id: fmt.Sprintf("%d", s.ID), extra: s}
	}
	return items
}

func seriesToItems(series []Series) []list.Item {
	items := make([]list.Item, len(series))
	for i, s := range series {
		desc := ""
		if s.Rating != "" && s.Rating != "0" { desc = "★ " + s.Rating }
		items[i] = menuItem{title: s.Name, desc: desc, id: fmt.Sprintf("%d", s.ID), intID: s.ID, extra: s}
	}
	return items
}

func ellipsis(s string, n int) string {
	if n <= 0 || len(s) <= n { return s }
	if n < 4 { return s[:n] }
	return s[:n-3] + "…"
}

func breadcrumbPop(crumb []string) []string {
	if len(crumb) <= 1 { return crumb }
	out := make([]string, len(crumb)-1)
	copy(out, crumb)
	return out
}

func cycleNext(opts []string, current, def string) string {
	if current == "" { current = def }
	for i, o := range opts {
		if o == current { return opts[(i+1)%len(opts)] }
	}
	return def
}

func progressBar(done, total, width int) string {
	if total == 0 || width <= 0 { return "" }
	if width > 38 { width = 38 }
	filled := int(float64(done) / float64(total) * float64(width))
	if filled > width { filled = width }
	return amberStyle.Render(strings.Repeat("█", filled)) +
		faintStyle.Render(strings.Repeat("░", width-filled))
}

func ratingStars(r string) string {
	var val float64
	fmt.Sscanf(r, "%f", &val)
	if val > 10 { val = 10 }
	stars := int(val / 2)
	return strings.Repeat("★", stars) + strings.Repeat("☆", 5-stars)
}

func wordWrapFirst(s string, w int) string {
	if w <= 0 || len(s) <= w { return s }
	for i := w; i > 0; i-- {
		if s[i] == ' ' { return s[:i] + "…" }
	}
	return s[:w] + "…"
}

func wordWrap(s string, w int) []string {
	if w <= 0 { return []string{s} }
	var lines []string
	current := ""
	for _, word := range strings.Fields(s) {
		if current == "" {
			current = word
		} else if len(current)+1+len(word) <= w {
			current += " " + word
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" { lines = append(lines, current) }
	return lines
}

func detailRow(label, value string, w int) string {
	labelW := lipgloss.Width(detailLabelStyle.Render(label))
	maxVal := w - labelW - 4
	if maxVal < 4 { maxVal = 4 }
	return "  " + detailLabelStyle.Render(label) + detailValueStyle.Render(ellipsis(value, maxVal))
}

func nextUnwatched(episodes []Episode, db *WatchedDB, limit int) []Episode {
	sorted := make([]Episode, len(episodes))
	copy(sorted, episodes)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].EpisodeNum.Int() < sorted[j].EpisodeNum.Int() })
	var result []Episode
	for _, ep := range sorted {
		if !db.IsWatched(ep.ID) {
			result = append(result, ep)
			if len(result) >= limit { break }
		}
	}
	return result
}

func minInt(a, b int) int {
	if a < b { return a }
	return b
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
	p := tea.NewProgram(
		newAppModel(
			NewClient(srv.URL, srv.Username, srv.Password),
			cfg,
			PlayOpts{Player: srv.Player, HWDec: srv.HWDec, Fullscreen: srv.Fullscreen},
		),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Erro:", err)
		os.Exit(1)
	}
}