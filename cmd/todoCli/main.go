package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zGraund/TodoCli/internal/cli"
	"github.com/zGraund/TodoCli/internal/db"
	"github.com/zGraund/TodoCli/internal/models"
)

func main() {
	if err := db.Get().AutoMigrate(&models.Todo{}); err != nil {
		log.Fatalln("failed to auto migrate database:", err)
	}

	if len(os.Args) > 1 && (os.Args[1] == "-directory" || os.Args[1] == "-d") {
		fmt.Println("The database directory is:\n" + db.Directory())
		os.Exit(0)
	}

	// fmt.Println(models.NewStats())
	p := tea.NewProgram(cli.InitMainModel())
	if _, err := p.Run(); err != nil {
		log.Fatalln("Alas, there's been an error:", err)
	}
}
