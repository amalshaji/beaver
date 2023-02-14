package main

import (
	"context"
	"fmt"
	"log"

	"github.com/amalshaji/beaver/internal/server/admin"
	"github.com/amalshaji/beaver/internal/server/db"
	"github.com/amalshaji/beaver/internal/utils"
	"github.com/labstack/gommon/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func handlePromptRenderError(err error) {
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
}

var superUserCmd = &cobra.Command{
	Use:   "createsuperuser",
	Short: "Create a new super user",
	Run: func(cmd *cobra.Command, args []string) {
		emailPrompt := promptui.Prompt{
			Label:    "Email",
			Validate: utils.ValidateEmail,
		}
		email, err := emailPrompt.Run()
		handlePromptRenderError(err)

		passwordPrompt := promptui.Prompt{
			Label:    "Password",
			Validate: utils.ValidatePassword,
			Mask:     '*',
		}
		password, err := passwordPrompt.Run()
		handlePromptRenderError(err)

		// Create new user service
		db := db.NewStore()
		ctx := context.Background()
		user := admin.NewUserService(db)

		err = user.CreateSuperUser(ctx, email, password)
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
