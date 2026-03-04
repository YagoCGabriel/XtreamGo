package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    args := os.Args[1:]
    if len(args) == 0 {
        runTUI()
        return
    }
    switch args[0] {
    case "add":
        cmdAdd()
    case "list":
        cmdList()
    case "use":
        if len(args) < 2 {
            fmt.Fprintln(os.Stderr, "Uso: xtreamgo use <index>")
            os.Exit(1)
        }
        cmdUse(args[1])
    case "remove":
        if len(args) < 2 {
            fmt.Fprintln(os.Stderr, "Uso: xtreamgo remove <index>")
            os.Exit(1)
        }
        cmdRemove(args[1])
    default:
        fmt.Println("Uso: xtreamgo [add|list|use <n>|remove <n>]")
    }
}

func cmdAdd() {
	cfg, _ := loadConfig()
	var name, url, user, pass, player string
	fmt.Print("Nome (ex: MeuIPTV): ")
	fmt.Scanln(&name)
	fmt.Print("URL (ex: http://server.com:8080): ")
	fmt.Scanln(&url)
	fmt.Print("Usuario: ")
	fmt.Scanln(&user)
	fmt.Print("Senha: ")
	fmt.Scanln(&pass)
	fmt.Print("Player (mpv/vlc) [mpv]: ")
	fmt.Scanln(&player)
	if player != "vlc" {
		player = "mpv"
	}

	url = strings.TrimRight(url, "/")
	fmt.Print("Testando conexao... ")
	client := NewClient(url, user, pass)
	info, err := client.Authenticate()
	if err != nil {
		fmt.Println("FALHOU:", err)
		fmt.Print("Salvar mesmo assim? (s/N): ")
		var ans string
		fmt.Scanln(&ans)
		if strings.ToLower(ans) != "s" {
			return
		}
	} else {
		fmt.Printf("OK! %s | %s | expira: %s\n", info.Username, info.Status, info.ExpDate)
	}

	cfg.Servers = append(cfg.Servers, Server{
		Name: name, URL: url,
		Username: user, Password: pass,
		Player: player,
	})
	cfg.Current = len(cfg.Servers) - 1
	saveConfig(cfg)
	fmt.Printf("Servidor '%s' adicionado (player: %s).\n", name, player)
}

func cmdList() {
    cfg, _ := loadConfig()
    if len(cfg.Servers) == 0 {
        fmt.Println("Nenhum servidor. Use: xtream-mpv add")
        return
    }
    for i, s := range cfg.Servers {
        mark := "  "
        if i == cfg.Current {
            mark = "* "
        }
        fmt.Printf("%s[%d] %s (%s)\n", mark, i, s.Name, s.URL)
    }
}

func cmdUse(idx string) {
    cfg, _ := loadConfig()
    var n int
    fmt.Sscanf(idx, "%d", &n)
    if n < 0 || n >= len(cfg.Servers) {
        fmt.Fprintln(os.Stderr, "Indice invalido")
        os.Exit(1)
    }
    cfg.Current = n
    saveConfig(cfg)
    fmt.Println("Ativo:", cfg.Servers[n].Name)
}

func cmdRemove(idx string) {
    cfg, _ := loadConfig()
    var n int
    fmt.Sscanf(idx, "%d", &n)
    if n < 0 || n >= len(cfg.Servers) {
        fmt.Fprintln(os.Stderr, "Indice invalido")
        os.Exit(1)
    }
    name := cfg.Servers[n].Name
    cfg.Servers = append(cfg.Servers[:n], cfg.Servers[n+1:]...)
    if cfg.Current >= len(cfg.Servers) && len(cfg.Servers) > 0 {
        cfg.Current = len(cfg.Servers) - 1
    }
    saveConfig(cfg)
    fmt.Printf("Removido: %s\n", name)
}