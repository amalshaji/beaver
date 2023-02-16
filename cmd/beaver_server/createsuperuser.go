package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/amalshaji/beaver/internal/server/admin"
	"github.com/amalshaji/beaver/internal/server/db"
	"github.com/amalshaji/beaver/internal/utils"
	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var superUserCmd = &cobra.Command{
	Use:   "createsuperuser",
	Short: "Create a new super user",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var email string
		var password []byte

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Email: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSuffix(email, "\n")

		if err := utils.ValidateEmail(email); err != nil {
			fmt.Println(color.Red("Enter a valid email address"))
			os.Exit(1)
		}

		fmt.Print("Password: ")
		password, _ = term.ReadPassword(0)
		fmt.Println(strings.Repeat("*", len(password)))

		if err := utils.ValidatePassword(string(password)); err != nil {
			fmt.Println(color.Red(err.Error()))
			os.Exit(1)
		}

		// Create new user service
		db := db.NewStore()
		ctx := context.Background()
		user := admin.NewUserService(db)

		_, err = user.CreateSuperUser(ctx, email, string(password))
		if err != nil {
			fmt.Println(color.Red(err.Error()))
		} else {
			fmt.Println(color.Green("Superuser created 🎉"))
		}
	},
}

func init() {
	rootCmd.AddCommand(superUserCmd)
}
