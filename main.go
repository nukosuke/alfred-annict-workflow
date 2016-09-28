package main

import (
	"errors"
	"fmt"
	"github.com/nukosuke/go-annict/annict"
	"github.com/nukosuke/go-alfred"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	"os"
)

const (
	bundleID = "tech.shibuya.alfred-annict-workflow"
	version  = "0.0.0"
)

func newAnnictClient() (*annict.Client, error) {
	var client *annict.Client

	viper.SetConfigName("config")
	viper.AddConfigPath(os.Getenv("HOME") + "/.annict")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	access_token := viper.GetString("access_token")
	if access_token == "" {
		client = annict.NewClient(nil)
		return client, errors.New("access_token was not found")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: access_token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = annict.NewClient(tc)

	return client, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "alfred-annict-workflow"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name: "works",
			Action: func(c *cli.Context) error {
				query := c.Args()
				client, err := newAnnictClient()

				if err != nil {
					return err
				}

				res, _, err := client.Works.List(&annict.WorksListOptions{
					Fields: []string{"id", "title", "media_text", "season_name_text"},
					FilterTitle: query[0],
				})

				if err != nil {
					return err
				}

				list := alfred.NewResponse()
				for _, item := range res.Works {
					list.AddItem(&alfred.AlfredResponseItem{
						Valid:    true,
						Uid:      fmt.Sprintf("%d", item.Id),
						Title:    item.Title,
						Arg:      fmt.Sprintf("https://annict.com/works/%d", item.Id),
						Subtitle: item.MediaText + " " + item.SeasonNameText,
					})
				}

				list.Print()
				return nil
			},
		},
		{
			Name: "watching",
			Action: func(c *cli.Context) error {
				client, err := newAnnictClient()

				res, _, err := client.Me.Works.List(&annict.MeWorksListOptions{
					Fields: []string{"id", "title", "media_text", "season_name_text"},
					FilterStatus: "watching",
				})

				if err != nil {
					return err
				}

				list := alfred.NewResponse()
				for _, item := range res.Works {
					list.AddItem(&alfred.AlfredResponseItem{
						Valid:    true,
						Uid:      fmt.Sprintf("%d", item.Id),
						Title:    item.Title,
						Arg:      fmt.Sprintf("https://annict.com/works/%d", item.Id),
						Subtitle: item.MediaText + " " + item.SeasonNameText,
					})
				}

				list.Print()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
