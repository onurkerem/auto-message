package cli

import (
	"fmt"
	"os"

	"github.com/onurkerem/auto-message/packages/cli/internal/config"
	"github.com/onurkerem/auto-message/packages/cli/internal/telegram"
	"github.com/spf13/cobra"
)

var appVersion = "dev"

func SetVersion(v string) {
	appVersion = v
}

func Execute() error {
	rootCmd := &cobra.Command{
		Use:     "auto-message",
		Short:   "Send Telegram notifications from the terminal",
		Long:    "A CLI tool for sending Telegram bot notifications from terminal commands and automation scripts.",
		Version: appVersion,
	}

	rootCmd.AddCommand(sendCmd())
	rootCmd.AddCommand(configCmd())

	return rootCmd.Execute()
}

func sendCmd() *cobra.Command {
	var cfgName string

	cmd := &cobra.Command{
		Use:   "send <message>",
		Short: "Send a message via Telegram",
		Long:  "Send a message via Telegram bot. Uses a named config profile, the default profile, or AUTO_MESSAGE_TOKEN and AUTO_MESSAGE_CHAT_ID environment variables.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			text := args[0]
			var token, chatID string

			if cfgName != "" {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("cannot load config: %w", err)
				}
				p, err := config.GetProfile(cfg, cfgName)
				if err != nil {
					return err
				}
				token = p.Token
				chatID = p.ChatID
			} else {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("cannot load config: %w", err)
				}
				p, err := config.GetDefault(cfg)
				if err == nil {
					token = p.Token
					chatID = p.ChatID
				}
			}

			if token == "" {
				token = os.Getenv("AUTO_MESSAGE_TOKEN")
			}
			if chatID == "" {
				chatID = os.Getenv("AUTO_MESSAGE_CHAT_ID")
			}

			if token == "" || chatID == "" {
				return fmt.Errorf("no bot token or chat ID found. Run 'auto-message config add' to create a profile, or set AUTO_MESSAGE_TOKEN and AUTO_MESSAGE_CHAT_ID environment variables")
			}

			client := telegram.NewClient()
			if err := client.Send(telegram.SendParams{
				Token:  token,
				ChatID: chatID,
				Text:   text,
			}); err != nil {
				return err
			}

			fmt.Println("Message sent successfully.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&cfgName, "config", "c", "", "named config profile to use")

	return cmd
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage bot configuration profiles",
		Long:  "Manage named Telegram bot configuration profiles stored at ~/.config/auto-message.",
	}

	cmd.AddCommand(configAddCmd())
	cmd.AddCommand(configListCmd())
	cmd.AddCommand(configDefaultCmd())
	cmd.AddCommand(configRemoveCmd())

	return cmd
}

func configAddCmd() *cobra.Command {
	var token, chatID string
	var setDefault bool

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a new bot configuration profile",
		Long:  "Add a named Telegram bot configuration with a bot token and chat ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if token == "" {
				return fmt.Errorf("bot token is required. Use --token flag")
			}
			if chatID == "" {
				return fmt.Errorf("chat ID is required. Use --chat-id flag")
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("cannot load config: %w", err)
			}

			p := config.Profile{
				Name:    name,
				Token:   token,
				ChatID:  chatID,
				Default: setDefault,
			}

			if err := config.AddProfile(cfg, p); err != nil {
				return err
			}

			fmt.Printf("Profile '%s' added successfully.\n", name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "Telegram bot token")
	cmd.Flags().StringVar(&chatID, "chat-id", "", "Telegram chat ID")
	cmd.Flags().BoolVarP(&setDefault, "default", "d", false, "set as default profile")

	return cmd
}

func configListCmd() *cobra.Command {
	var showTokens bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all bot configuration profiles",
		Long:  "List all saved Telegram bot configuration profiles. Tokens are masked by default.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("cannot load config: %w", err)
			}

			if len(cfg.Profiles) == 0 {
				fmt.Println("No profiles configured. Run 'auto-message config add' to create one.")
				return nil
			}

			for _, p := range cfg.Profiles {
				token := config.MaskToken(p.Token)
				if showTokens {
					token = p.Token
				}
				defaultMark := ""
				if p.Default {
					defaultMark = " (default)"
				}
				fmt.Printf("  %s%s\n", p.Name, defaultMark)
				fmt.Printf("    token:   %s\n", token)
				fmt.Printf("    chat-id: %s\n", p.ChatID)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&showTokens, "show-tokens", false, "show full tokens instead of masked values")

	return cmd
}

func configDefaultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default <name>",
		Short: "Set a profile as the default",
		Long:  "Set the specified profile as the default for the send command.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("cannot load config: %w", err)
			}

			if err := config.SetDefault(cfg, name); err != nil {
				return err
			}

			fmt.Printf("Profile '%s' set as default.\n", name)
			return nil
		},
	}

	return cmd
}

func configRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a bot configuration profile",
		Long:  "Remove the specified profile from the configuration.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("cannot load config: %w", err)
			}

			if err := config.RemoveProfile(cfg, name); err != nil {
				return err
			}

			fmt.Printf("Profile '%s' removed.\n", name)
			return nil
		},
	}

	return cmd
}
