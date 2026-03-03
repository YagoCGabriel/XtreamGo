# XtreamGo

CLI interativo para assistir listas Xtream Code (IPTV) no MPV ou VLC, com TUI navegavel no terminal.

## Funcionalidades

- TV ao Vivo, Filmes (VOD) e Series com navegacao por categorias
- Filtro por nome em qualquer lista (tecla /)
- Multiplos servidores Xtream Code com troca rapida
- Suporte a MPV e VLC (incluindo VLC portable)
- Hardware decoding configuravel por servidor
- Configuracoes alteraveis dentro da TUI sem sair do programa
- Config persistente em ~/.config/xtream-mpv/config.json

## Requisitos

- Go 1.21+        https://go.dev/dl/
- MPV             https://mpv.io/installation/
- VLC (opcional)  https://www.videolan.org/vlc/

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

### VLC Portable (opcional)

Coloque os arquivos do VLC portable nesta estrutura ao lado do executavel:

    xtreamgo.exe
    players/
      vlc/
        vlc-portable.exe
        (demais arquivos do VLC portable)

Ordem de deteccao automatica do VLC:
  1. players/vlc/vlc-portable.exe  (relativo ao executavel)
  2. players/vlc/vlc.exe
  3. players/vlc/vlc               (Linux/macOS)
  4. vlc no PATH do sistema
  5. C:\Program Files\VideoLAN\VLC\vlc.exe
  6. /Applications/VLC.app/Contents/MacOS/VLC

## Uso

### Adicionar servidor

    xtreamgo add

    Nome (ex: MeuIPTV): MeuIPTV
    URL (ex: http://server.com:8080): http://meuservidor.com:8080
    Usuario: meuusuario
    Senha: minhasenha
    Player (mpv/vlc) [mpv]: mpv
    Testando conexao... OK! meuusuario | Active | expira: 2026-12-31
    Servidor 'MeuIPTV' adicionado.

### Abrir a TUI

    xtreamgo

### Outros comandos

    xtreamgo list         listar servidores (* = ativo)
    xtreamgo use 1        ativar servidor pelo indice
    xtreamgo remove 0     remover servidor pelo indice

## Navegacao na TUI

    Setas / Enter   Navegar e selecionar
    Esc             Voltar
    /               Filtrar por nome
    q               Sair

### Fluxo

    Menu principal
     TV ao Vivo  ->  Categorias  ->  Canais      ->  Player
     Filmes      ->  Categorias  ->  Filmes       ->  Player
     Series      ->  Categorias  ->  Series
                                  ->  Temporadas
                                  ->  Episodios   ->  Player
     Configuracoes

## Tela de Configuracoes

Acesse pelo menu principal -> Configuracoes.

    p   Player          mpv <-> vlc
    h   HW Decoding     no -> auto -> vaapi -> vdpau -> nvdec -> videotoolbox
    f   Fullscreen      sim <-> nao
    s   Servidor ativo  lista de servidores salvos
    Esc Voltar

Cada tecla salva imediatamente no config.json.

## Hardware Decoding

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
    ├── main.go       entrypoint + comandos CLI
    ├── api.go        cliente da API Xtream Codes
    ├── config.go     configuracao e servidores
    ├── player.go     launcher MPV e VLC
    ├── tui.go        interface TUI (Bubbletea)
    ├── flexjson.go   tipos flexiveis para JSON inconsistente da API
    ├── go.mod
    ├── go.sum
    ├── README.md
    ├── LICENSE
    └── players/
        └── vlc/
            └── vlc-portable.exe

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
