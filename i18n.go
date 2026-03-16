package main

// =============================================================================
// I18N — INTERNACIONALIZAÇÃO
// =============================================================================

type Lang struct {
	// Menu principal
	MenuLive     string
	MenuLiveDesc string
	MenuVOD      string
	MenuVODDesc  string
	MenuSeries   string
	MenuSerDesc  string
	MenuSettings string
	MenuSetDesc  string

	// Navegação / labels
	Home          string
	LiveTV        string
	Movies        string
	Series        string
	Settings      string
	Categories    string
	Seasons       string
	Episodes      string
	Servers       string
	All           string
	AllMovies     string
	AllSeries     string
	SearchAll     string
	SearchMovies  string
	SearchSeries  string

	// Detail panel
	Server     string
	URL        string
	User       string
	Player     string
	Format     string
	Rating     string
	Season     string
	Episode    string
	Watched    string
	Unwatched  string
	Episodes2  string // "episódios" (noun)
	WatchedOf  string // "%d / %d assistidos"
	Completed  string // "%d%% concluído"
	NextUnwatched string
	NoItems    string
	SelectCat  string
	SelectSeries string
	OpenIn     string // "Abrir no %s"
	PlayMark   string
	ToggleMark string
	LiveBadge  string
	TogglePanel string

	// Settings fields
	SetServerActive string
	SetPlayback     string
	SetServers      string
	SetHistory      string
	SetWatched      string // "★ %d episódio(s) assistido(s)"
	SetPlayerCycle  string
	SetHWCycle      string
	SetFullscreen   string
	SetFullYes      string
	SetFullNo       string
	SetSwitch       string
	SetConfigured   string // "%d configurado(s)"

	// Status / feedback
	StatusSaved    string
	StatusServer   string
	StatusWatched  string
	StatusUnwatched string

	// Footer keys
	KeySelect    string
	KeyWatch     string
	KeyMark      string
	KeyBack      string
	KeyFilter    string
	KeyQuit      string
	KeyPlayer    string
	KeyHWDec     string
	KeyFullscreen string
	KeyServers   string

	// Misc
	NoneConfigured string
	TabHint        string
	EpCount        string   // "%d ep."
	WatchedCount   string   // "★ %d/%d"
}

var Langs = map[string]Lang{
	"pt-BR": {
		MenuLive:     "TV ao Vivo",
		MenuLiveDesc: "Canais em tempo real",
		MenuVOD:      "Filmes",
		MenuVODDesc:  "Video on demand",
		MenuSeries:   "Séries",
		MenuSerDesc:  "Episódios por temporada",
		MenuSettings: "Configurações",
		MenuSetDesc:  "Player, servidor, idioma",

		Home:          "Início",
		LiveTV:        "TV ao Vivo",
		Movies:        "Filmes",
		Series:        "Séries",
		Settings:      "Configurações",
		Categories:    "Categorias",
		Seasons:       "Temporadas",
		Episodes:      "Episódios",
		Servers:       "Servidores",
		All:           "Todos",
		AllMovies:     "Todos os Filmes",
		AllSeries:     "Todas as Séries",
		SearchAll:     "Buscar em tudo",
		SearchMovies:  "Pesquisar em todos os filmes",
		SearchSeries:  "Pesquisar em todas as séries",

		Server:        "Servidor",
		URL:           "URL",
		User:          "Usuário",
		Player:        "Player",
		Format:        "Formato",
		Rating:        "Nota",
		Season:        "Temporada",
		Episode:       "Episódio %02d",
		Watched:       "★ ASSISTIDO",
		Unwatched:     "▶ NÃO ASSISTIDO",
		Episodes2:     "episódios",
		WatchedOf:     "%d / %d assistidos",
		Completed:     "%d%% concluído",
		NextUnwatched: "Próximos não assistidos:",
		NoItems:       "Nenhum item selecionado.",
		SelectCat:     "Selecione uma categoria.",
		SelectSeries:  "Selecione uma série.",
		OpenIn:        "Enter para abrir no %s",
		PlayMark:      "Enter  reproduzir + marcar",
		ToggleMark:    "W      marcar / desmarcar",
		LiveBadge:     "● AO VIVO",
		TogglePanel:   "tab  mostrar/ocultar detalhes",

		SetServerActive: "Servidor Ativo",
		SetPlayback:     "Reprodução",
		SetServers:      "Servidores",
		SetHistory:      "Histórico",
		SetWatched:      "★  %d episódio(s) assistido(s)",
		SetPlayerCycle:  "mpv → vlc → kmplayer",
		SetHWCycle:      "no → auto → vaapi → vdpau → nvdec → videotoolbox",
		SetFullscreen:   "abrir em tela cheia",
		SetFullYes:      "sim",
		SetFullNo:       "não",
		SetSwitch:       "selecionar servidor ativo",
		SetConfigured:   "%d configurado(s)",

		StatusSaved:     "configuração salva",
		StatusServer:    "servidor alterado",
		StatusWatched:   "★ marcado como assistido",
		StatusUnwatched: "☆ desmarcado",

		KeySelect:    "selecionar",
		KeyWatch:     "assistir",
		KeyMark:      "marcar",
		KeyBack:      "voltar",
		KeyFilter:    "filtrar",
		KeyQuit:      "sair",
		KeyPlayer:    "player",
		KeyHWDec:     "hwdec",
		KeyFullscreen: "fullscreen",
		KeyServers:   "servidores",

		NoneConfigured: "Nenhum servidor configurado.",
		TabHint:        "tab ⊠",
		EpCount:        "%d ep.",
		WatchedCount:   "★ %d/%d",
	},

	"en-US": {
		MenuLive:     "Live TV",
		MenuLiveDesc: "Real-time channels",
		MenuVOD:      "Movies",
		MenuVODDesc:  "Video on demand",
		MenuSeries:   "Series",
		MenuSerDesc:  "Episodes by season",
		MenuSettings: "Settings",
		MenuSetDesc:  "Player, server, language",

		Home:          "Home",
		LiveTV:        "Live TV",
		Movies:        "Movies",
		Series:        "Series",
		Settings:      "Settings",
		Categories:    "Categories",
		Seasons:       "Seasons",
		Episodes:      "Episodes",
		Servers:       "Servers",
		All:           "All",
		AllMovies:     "All Movies",
		AllSeries:     "All Series",
		SearchAll:     "Search all",
		SearchMovies:  "Search all movies",
		SearchSeries:  "Search all series",

		Server:        "Server",
		URL:           "URL",
		User:          "Username",
		Player:        "Player",
		Format:        "Format",
		Rating:        "Rating",
		Season:        "Season",
		Episode:       "Episode %02d",
		Watched:       "★ WATCHED",
		Unwatched:     "▶ NOT WATCHED",
		Episodes2:     "episodes",
		WatchedOf:     "%d / %d watched",
		Completed:     "%d%% completed",
		NextUnwatched: "Next unwatched:",
		NoItems:       "No item selected.",
		SelectCat:     "Select a category.",
		SelectSeries:  "Select a series.",
		OpenIn:        "Enter to open in %s",
		PlayMark:      "Enter  play + mark watched",
		ToggleMark:    "W      mark / unmark",
		LiveBadge:     "● LIVE",
		TogglePanel:   "tab  show/hide details",

		SetServerActive: "Active Server",
		SetPlayback:     "Playback",
		SetServers:      "Servers",
		SetHistory:      "History",
		SetWatched:      "★  %d episode(s) watched",
		SetPlayerCycle:  "mpv → vlc → kmplayer",
		SetHWCycle:      "no → auto → vaapi → vdpau → nvdec → videotoolbox",
		SetFullscreen:   "open fullscreen",
		SetFullYes:      "yes",
		SetFullNo:       "no",
		SetSwitch:       "select active server",
		SetConfigured:   "%d configured",

		StatusSaved:     "settings saved",
		StatusServer:    "server changed",
		StatusWatched:   "★ marked as watched",
		StatusUnwatched: "☆ unmarked",

		KeySelect:    "select",
		KeyWatch:     "watch",
		KeyMark:      "mark",
		KeyBack:      "back",
		KeyFilter:    "filter",
		KeyQuit:      "quit",
		KeyPlayer:    "player",
		KeyHWDec:     "hwdec",
		KeyFullscreen: "fullscreen",
		KeyServers:   "servers",

		NoneConfigured: "No server configured.",
		TabHint:        "tab ⊠",
		EpCount:        "%d ep.",
		WatchedCount:   "★ %d/%d",
	},
}

func T(cfg *Config) Lang {
	lang := cfg.Language
	if lang == "" {
		lang = "pt-BR"
	}
	l, ok := Langs[lang]
	if !ok {
		return Langs["pt-BR"]
	}
	return l
}
