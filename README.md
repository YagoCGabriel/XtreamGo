# XtreamGo

> Cliente IPTV para listas Xtream Code no terminal. Assiste TV ao Vivo, Filmes e Series
> com uma TUI split-panel no MPV, VLC ou KMPlayer.

<img width="952" height="552" alt="Interface" src="https://github.com/user-attachments/assets/923fb96d-f09e-45fe-9f26-27af12be9d0f" />


---

## Funcionalidades

- Layout split com painel de detalhes ao lado da lista (toggle com Tab)
- TV ao Vivo, Filmes (VOD) e Series com navegacao por categorias
- Busca global em todos os Filmes ou Series sem precisar escolher categoria
- Filtro fuzzy por nome em qualquer lista (tecla /)
- Episodios marcados com estrela (automatico ao reproduzir, ou manual com W)
- Progress bar de temporada: quantos episodios ja foram assistidos
- Proximos episodios nao assistidos exibidos no painel de detalhes
- Interface em portugues (pt-BR) ou ingles (en-US), alternavel com tecla L
- Header adaptativo: trunca informacoes conforme a largura do terminal
- Multiplos servidores Xtream Code com troca rapida
- Suporte a MPV, VLC e KMPlayer (incluindo versoes portable)
- Hardware decoding configuravel por servidor (vaapi, nvdec, videotoolbox...)
- Configuracoes salvas imediatamente no config.json
- Historico de episodios assistidos salvo em watched.json

---

## Requisitos

- Go 1.21+              https://go.dev/dl/
- MPV                   https://mpv.io/installation/
- VLC (opcional)        https://www.videolan.org/vlc/
- KMPlayer (opcional)   https://www.kmplayer.com/

---

## Instalacao

### Baixar binario pre-compilado (recomendado)

Acesse a pagina de releases e baixe o executavel para o seu sistema:

    https://github.com/user/xtreamgo/releases

Plataformas disponíveis:

    xtreamgo-windows-amd64.exe
    xtreamgo-linux-amd64
    xtreamgo-darwin-amd64
    xtreamgo-darwin-arm64

### Compilar a partir do codigo-fonte

    git clone https://github.com/user/xtreamgo.git
    cd xtreamgo
    go mod tidy
    go build -o xtreamgo .

### Instalar globalmente (opcional)

    # Linux / macOS
    sudo mv xtreamgo /usr/local/bin/

    # Windows: mova xtreamgo.exe para uma pasta no PATH

### Players Portable (opcional)

Coloque os arquivos dos players portable ao lado do executavel:

    xtreamgo.exe
    players/
      vlc/
        vlc-portable.exe
        (demais arquivos do VLC portable)
      kmplayer/
        KMPlayerPortable.exe
        (demais arquivos do KMPlayer portable)

#### Ordem de deteccao - VLC

  1. players/vlc/vlc-portable.exe       (relativo ao executavel)
  2. players/vlc/vlc.exe
  3. players/vlc/vlc                    (Linux/macOS)
  4. vlc no PATH do sistema
  5. C:\Program Files\VideoLAN\VLC\vlc.exe
  6. /Applications/VLC.app/Contents/MacOS/VLC

#### Ordem de deteccao - KMPlayer

  1. players/kmplayer/KMPlayerPortable.exe   (relativo ao executavel)
  2. players/kmplayer/KMPlayer.exe
  3. players/kmplayer/kmplayer.exe
  4. KMPlayer no PATH do sistema
  5. C:\Program Files\KMPlayer\KMPlayer.exe
  6. C:\Program Files (x86)\KMPlayer\KMPlayer.exe

---

## Uso

### Adicionar servidor

    xtreamgo add

    Name / Nome: MeuIPTV
    URL: http://meuservidor.com:8080
    Username / Usuario: meuusuario
    Password / Senha: minhasenha
    Player (mpv/vlc/kmplayer) [mpv]: mpv
    Testing connection... OK! meuusuario | Active | expires: 2026-12-31
    Server 'MeuIPTV' added (player: mpv).

### Abrir a TUI

    xtreamgo

### Outros comandos

    xtreamgo list         listar servidores (* = ativo) e idioma atual
    xtreamgo use 1        ativar servidor pelo indice
    xtreamgo remove 0     remover servidor pelo indice

---

## Interface (TUI)

Acima de 110 colunas de largura, a interface usa um layout split:

    ┌─────────────────────────────────────────────────────────────────────┐
    │  XtreamGo v1.0                         MeuIPTV  ·  mpv  ·  pt-BR   │
    ├─────────────────────────────────────────────────────────────────────┤
    │  Home › Séries › Drama › Dark                         tab ⊠        │
    ├───────────────────────────────────────┬─────────────────────────────┤
    │                                       │                             │
    │    Season 1        8 ep. ★ 5/8        │  Temporada 1               │
    │  ▸ Season 2        10 ep. ★ 2/10      │  ────────────────          │
    │    Season 3        8 ep.              │  Episodios  10             │
    │                                       │  Assistidos  2 / 10        │
    │                                       │                             │
    │                                       │  ████████░░░░░░░░░░░░  20% │
    │                                       │                             │
    │                                       │  Proximos nao assistidos:  │
    │                                       │    E03  Adam und Eva        │
    ├───────────────────────────────────────┴─────────────────────────────┤
    │  enter selecionar  ·  esc voltar  ·  / filtrar  ·  q sair          │
    └─────────────────────────────────────────────────────────────────────┘

Abaixo de 110 colunas a lista ocupa a tela inteira.

### Teclas globais

    Setas / Enter   Navegar e selecionar
    Esc             Voltar
    /               Filtrar por nome
    Tab             Mostrar / ocultar painel de detalhes
    q               Sair

### Teclas na tela de episodios

    Enter           Reproduzir e marcar como assistido automaticamente
    W               Marcar / desmarcar manualmente como assistido

### Teclas nas Configuracoes

    p   Player          ciclo: mpv -> vlc -> kmplayer
    h   HW Decoding     ciclo: no -> auto -> vaapi -> vdpau -> nvdec -> videotoolbox
    f   Fullscreen      sim <-> nao
    l   Idioma          pt-BR <-> en-US
    s   Servidor        lista de servidores para trocar
    Esc Voltar

Cada tecla salva imediatamente.

### Busca global

Nas telas de Filmes e Series o primeiro item e sempre "Buscar em tudo".
Ao seleciona-lo, todos os itens de todas as categorias sao carregados.
Use / para filtrar por nome em toda a base.

---

## Episodios Assistidos

O XtreamGo marca automaticamente um episodio como assistido ao reproduzi-lo.
Use W para marcar ou desmarcar manualmente sem reproduzir.

O estado e salvo em:

    Windows:   %APPDATA%\xtream-mpv\watched.json
    Linux:     ~/.config/xtream-mpv/watched.json
    macOS:     ~/Library/Application Support/xtream-mpv/watched.json

Na tela de temporadas a desc mostra quantos episodios foram assistidos:

    Season 1   8 ep.  ★ 8/8    <- temporada concluida
    Season 2   10 ep. ★ 3/10   <- em andamento
    Season 3   6 ep.            <- nao iniciada

---

## Idioma

O idioma padrao e pt-BR. Para alternar vá em Configuracoes e pressione L.
A mudanca e imediata e persistente. Toda a interface muda: menus, labels,
breadcrumbs, badges e teclas do footer.

    pt-BR   Portugues brasileiro (padrao)
    en-US   English

---

## Players

### MPV (padrao recomendado)

Player leve e poderoso, ideal para streams IPTV. O XtreamGo configura
automaticamente flags de cache, correcao de GOP e tolerancia a erros.

    https://mpv.io/installation/

### VLC

Amplamente compativel, boa opcao de fallback. Suporta versao portable.

    https://www.videolan.org/vlc/
    Portable: players/vlc/vlc-portable.exe

### KMPlayer

Player popular no Windows com interface propria. Suporta versao portable.

    https://www.kmplayer.com/
    Portable: players/kmplayer/KMPlayerPortable.exe

---

## Hardware Decoding (MPV)

    no             Software (padrao). Maxima compatibilidade para streams IPTV.
    auto           Detecta automaticamente.
    vaapi          Intel/AMD no Linux.
    vdpau          NVIDIA no Linux (legado).
    nvdec          NVIDIA moderno (RTX/GTX recentes).
    videotoolbox   macOS (qualquer Mac).

Dica: se aparecerem artefatos ou blocos na imagem, use "no".
Hardware decoding pode causar corrupcao em streams com bitrate instavel.

---

## Estrutura do projeto

    xtreamgo/
    ├── main.go         entrypoint + comandos CLI
    ├── api.go          cliente da API Xtream Codes
    ├── config.go       configuracao e servidores
    ├── player.go       launcher MPV, VLC e KMPlayer
    ├── tui.go          interface TUI com layout split
    ├── i18n.go         strings pt-BR e en-US
    ├── watched.go      historico de episodios assistidos
    ├── flexjson.go     tipos flexiveis para JSON inconsistente da API
    ├── go.mod
    ├── go.sum
    ├── README.md
    ├── LICENSE
    └── players/
        ├── vlc/
        │   └── vlc-portable.exe
        └── kmplayer/
            └── KMPlayerPortable.exe

---

## Arquivos de configuracao

    Windows   %APPDATA%\xtream-mpv\
    Linux     ~/.config/xtream-mpv/
    macOS     ~/Library/Application Support/xtream-mpv/

    config.json    servidores, player, idioma, fullscreen, hwdec
    watched.json   historico de episodios assistidos

### config.json

    {
      "servers": [
        {
          "name": "MeuIPTV",
          "url": "http://meuservidor.com:8080",
          "username": "usuario",
          "password": "senha",
          "player": "mpv",
          "hwdec": "no",
          "fullscreen": false
        }
      ],
      "current": 0,
      "language": "pt-BR"
    }

---

## Dependencias

    github.com/charmbracelet/bubbletea   Framework da TUI
    github.com/charmbracelet/bubbles     Lista com filtro e spinner
    github.com/charmbracelet/lipgloss    Estilos e cores no terminal

---

## Licenca

MIT License - veja o arquivo LICENSE para detalhes.
