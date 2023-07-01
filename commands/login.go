package commands

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your account",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		token, err := config.ReadToken()
		client := api.NewClient(token)
		reader := bufio.NewReader(os.Stdin)

		if err != nil {
			log.Println("Error read token: ", err, "; would do login")
			token, email, err := doLogin(ctx, reader)
			if err != nil {
				log.Fatal("login error: ", err)
			}
			config.UpdateToken(token)
			log.Printf("Logged in as %s successfully!", email)
		} else {
			user, err := client.Me()
			if err == nil {
				log.Println("loaded existing token. User: ", user.Email)
				fmt.Printf("Login as another user?(Y/%s)", wrapUnderline("N"))

				option, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal("read input failed with err: ", err)
				}
				option = strings.TrimSpace(option)
				if option != "Y" {
					return
				}
			} else {
				log.Println("login using existing token failed: ", err.Error(), "; would proceed to login")
			}

			token, email, err := doLogin(ctx, reader)
			if err != nil {
				log.Fatal("login error: ", err)
			}
			config.UpdateToken(token)
			log.Printf("Logged in as %s successfully!", email)
		}
	},
}

func wrapUnderline(s string) string {
	return fmt.Sprintf("\x1b[4m\x1b[1m%s\x1b[0m", s)
}

func doLogin(ctx context.Context, reader *bufio.Reader) (token, email string, err error) {
	fmt.Print("Enter Email: ")
	email, err = reader.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Print("Enter Password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	resp, err := api.Login(ctx, email, password)
	if err != nil {
		log.Fatal(err)
	}
	token = string(resp.Token)
	email = string(resp.User.Email)
	return
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
