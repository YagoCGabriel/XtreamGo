# XtreamGo

CLI interativo para assistir listas Xtream Code (IPTV) no MPV, VLC ou KMPlayer,
com TUI navegavel no terminal.

<img width="952" height="552" alt="image" src="https://github.com/user-attachments/assets/5a4ccebb-b304-4fcc-84ee-fa4a698b2bd3" />


## Funcionalidades

- TV ao Vivo, Filmes (VOD) e Series com navegacao por categorias
- Busca global em todos os Filmes ou Series de uma vez (sem escolher categoria)
- Filtro por nome em qualquer lista (tecla /)
- Breadcrumb de navegacao sempre visivel
- Spinner animado durante carregamentos
- Multiplos servidores Xtream Code com troca rapida
- Suporte a MPV, VLC e KMPlayer (incluindo versoes portable)
- Hardware decoding configuravel por servidor
- Configuracoes alteraveis dentro da TUI sem sair do programa
- Config persistente em ~/.config/xtream-mpv/config.json

## Requisitos

- Go 1.21+              https://go.dev/dl/
- MPV                   https://mpv.io/installation/
- VLC (opcional)        https://www.videolan.org/vlc/
- KMPlayer (opcional)   https://www.kmplayer.com/

## Instalacao

### Compilar

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

## Uso

### Adicionar servidor

    xtreamgo add

    Nome (ex: MeuIPTV): MeuIPTV
    URL (ex: http://server.com:8080): http://meuservidor.com:8080
    Usuario: meuusuario
    Senha: minhasenha
    Player (mpv/vlc/kmplayer) [mpv]: mpv
    Testando conexao... OK! meuusuario | Active | expira: 2026-12-31
    Servidor 'MeuIPTV' adicionado.

### Abrir a TUI

    xtreamgo

### Outros comandos

    xtreamgo list         listar servidores (* = ativo)
    xtreamgo use 1        ativar servidor pelo indice
    xtreamgo remove 0     remover servidor pelo indice

## Interface

A TUI conta com:

- Header fixo com logo, versao e servidor/player ativo
- Breadcrumb de navegacao com separador › e item atual em ciano
- Spinner animado durante carregamentos
- Item selecionado com borda esquerda roxa e fundo destacado
- Footer dinamico que muda os atalhos conforme a tela atual

## Navegacao na TUI

    Setas / Enter   Navegar e selecionar
    Esc             Voltar
    /               Filtrar por nome na lista atual
    q               Sair

### Fluxo

    Menu principal
     TV ao Vivo  ->  Categorias               ->  Canais      ->  Player
     Filmes      ->  [Buscar em tudo]          ->  Todos       ->  Player
                  ->  Categoria                ->  Filmes      ->  Player
     Series      ->  [Buscar em tudo]          ->  Todas       ->  Player
                  ->  Categoria  ->  Series
                                  ->  Temporadas
                                  ->  Episodios ->  Player
     Configuracoes

### Busca global em Filmes e Series

Nas telas de categorias de Filmes e Series, o primeiro item da lista e sempre
"Buscar em tudo". Ao seleciona-lo, o XtreamGo carrega todos os itens de todas
as categorias de uma vez e abre a lista com o filtro ativo. Use / para digitar
e encontrar qualquer titulo independente da categoria.

    Filmes
     󰍉  Buscar em tudo   <- carrega TODOS os filmes, use / para filtrar
     ──────────────────
     Acao
     Comedia
     Drama
     ...

    Series
     󰍉  Buscar em tudo   <- carrega TODAS as series, use / para filtrar
     ──────────────────
     Animacao
     Drama
     ...

## Tela de Configuracoes

Acesse pelo menu principal -> Configuracoes.

    p   Player          mpv -> vlc -> kmplayer -> mpv (ciclo)
    h   HW Decoding     no -> auto -> vaapi -> vdpau -> nvdec -> videotoolbox
    f   Fullscreen      sim <-> nao
    s   Servidor ativo  lista de servidores salvos
    Esc Voltar

Cada tecla salva imediatamente no config.json.

## Players

### MPV (padrao recomendado)

Player leve e poderoso, ideal para streams IPTV. Melhor compatibilidade
com os parametros de correcao de GOP e cache configurados pelo XtreamGo.

    https://mpv.io/installation/

### VLC

Amplamente compativel, boa opcao de fallback. Suporta versao portable.

    https://www.videolan.org/vlc/
    Portable: players/vlc/vlc-portable.exe

### KMPlayer

Player popular no Windows com interface propria. Suporta versao portable.

    https://www.kmplayer.com/
    Portable: players/kmplayer/KMPlayerPortable.exe

## Hardware Decoding

Disponivel apenas no MPV. Configuravel por servidor.

    no             Software (padrao). Maxima compatibilidade para streams IPTV.
    auto           Detecta automaticamente.
    vaapi          Intel/AMD no Linux.
    vdpau          NVIDIA no Linux (legado).
    nvdec          NVIDIA moderno (RTX/GTX recentes).
    videotoolbox   macOS.

Dica: se aparecerem artefatos ou blocos na imagem, use "no" (software).
Hardware decoding pode causar corrupcao em streams com bitrate instavel.

## Estrutura do projeto

    xtreamgo/
    ├── main.go         entrypoint + comandos CLI
    ├── api.go          cliente da API Xtream Codes
    ├── config.go       configuracao e servidores
    ├── player.go       launcher MPV, VLC e KMPlayer
    ├── tui.go          interface TUI (Bubbletea)
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

## Dependencias

    github.com/charmbracelet/bubbletea   Framework da TUI
    github.com/charmbracelet/bubbles     Componente de lista com filtro
    github.com/charmbracelet/lipgloss    Estilos e cores no terminal

## Config JSON

Salvo em ~/.config/xtream-mpv/config.json:

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
      "current": 0
    }

## Licenca

MIT License - veja o arquivo LICENSE para detalhes.
