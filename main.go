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
			fmt.Fprintln(os.Stderr, "Usage: xtreamgo use <index>")
			os.Exit(1)
		}
		cmdUse(args[1])
	case "remove":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: xtreamgo remove <index>")
			os.Exit(1)
		}
		cmdRemove(args[1])
	default:
		fmt.Println("Usage: xtreamgo [add|list|use <n>|remove <n>]")
	}
}

func cmdAdd() {
	cfg, _ := loadConfig()
	if cfg == nil {
		cfg = &Config{Language: "pt-BR"}
	}

	var name, url, user, pass, player string

	fmt.Print("Name / Nome: ")
	fmt.Scanln(&name)
	fmt.Print("URL: ")
	fmt.Scanln(&url)
	fmt.Print("Username / Usuario: ")
	fmt.Scanln(&user)
	fmt.Print("Password / Senha: ")
	fmt.Scanln(&pass)
	fmt.Print("Player (mpv/vlc/kmplayer) [mpv]: ")
	fmt.Scanln(&player)
	if player != "vlc" && player != "kmplayer" {
		player = "mpv"
	}

	url = strings.TrimRight(url, "/")

	fmt.Print("Testing connection / Testando conexao... ")
	client := NewClient(url, user, pass)
	info, err := client.Authenticate()
	if err != nil {
		fmt.Println("FAILED / FALHOU:", err)
		fmt.Print("Save anyway? / Salvar mesmo assim? (y/s/N): ")
		var ans string
		fmt.Scanln(&ans)
		ans = strings.ToLower(ans)
		if ans != "y" && ans != "s" {
			return
		}
	} else {
		fmt.Printf("OK! %s | %s | expires: %s\n", info.Username, info.Status, info.ExpDate)
	}

	cfg.Servers = append(cfg.Servers, Server{
		Name:     name,
		URL:      url,
		Username: user,
		Password: pass,
		Player:   player,
	})
	cfg.Current = len(cfg.Servers) - 1
	saveConfig(cfg)
	fmt.Printf("Server '%s' added (player: %s).\n", name, player)
}

func cmdList() {
	cfg, _ := loadConfig()
	if cfg == nil || len(cfg.Servers) == 0 {
		fmt.Println("No servers. Run: xtreamgo add")
		return
	}
	for i, s := range cfg.Servers {
		mark := "  "
		if i == cfg.Current {
			mark = "* "
		}
		player := s.Player
		if player == "" {
			player = "mpv"
		}
		fmt.Printf("%s[%d] %s (%s) [%s]\n", mark, i, s.Name, s.URL, player)
	}
	fmt.Printf("\nLanguage / Idioma: %s\n", cfg.Language)
}

func cmdUse(idx string) {
	cfg, _ := loadConfig()
	var n int
	fmt.Sscanf(idx, "%d", &n)
	if n < 0 || n >= len(cfg.Servers) {
		fmt.Fprintln(os.Stderr, "Invalid index / Indice invalido")
		os.Exit(1)
	}
	cfg.Current = n
	saveConfig(cfg)
	fmt.Println("Active:", cfg.Servers[n].Name)
}

func cmdRemove(idx string) {
	cfg, _ := loadConfig()
	var n int
	fmt.Sscanf(idx, "%d", &n)
	if n < 0 || n >= len(cfg.Servers) {
		fmt.Fprintln(os.Stderr, "Invalid index / Indice invalido")
		os.Exit(1)
	}
	name := cfg.Servers[n].Name
	cfg.Servers = append(cfg.Servers[:n], cfg.Servers[n+1:]...)
	if cfg.Current >= len(cfg.Servers) && len(cfg.Servers) > 0 {
		cfg.Current = len(cfg.Servers) - 1
	}
	saveConfig(cfg)
	fmt.Printf("Removed / Removido: %s\n", name)
}
