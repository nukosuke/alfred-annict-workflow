package main

import (
  "fmt"
  "github.com/urfave/cli"
  "github.com/pascalw/go-alfred"
  "github.com/spf13/viper"
  "github.com/nukosuke/go-annict/annict"
  "golang.org/x/oauth2"
  "os"
)

const (
  bundleID = "tech.shibuya.alfred-annict-workflow"
  version  = "0.0.0"
)

func newAnnictClient() (*annict.Client, error) {
  var client *annict.Client
  return client, nil
}

func main() {
  app := cli.NewApp()
  app.Name    = "alfred-annict-workflow"
  app.Version = version
  app.Commands = []cli.Command{
    {
      Name: "works",
      Action: func(c *cli.Context) error {
        viper.SetConfigName("config")
        viper.AddConfigPath(os.Getenv("HOME")+"/.annict")
        viper.SetConfigType("json")
        err := viper.ReadInConfig()
        if err != nil {
          return err
        }

        access_token := viper.GetString("access_token")

        query := c.Args()

        ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: access_token})
        tc := oauth2.NewClient(oauth2.NoContext, ts)
        client := annict.NewClient(tc)

        res, _, err := client.Works.List(&annict.WorksListOptions{
          FilterTitle: query[0],
        })

        if err == nil {
          list := alfred.NewResponse()
          for _, item := range res.Works {
            list.AddItem(&alfred.AlfredResponseItem{
              Valid: true,
              Uid: fmt.Sprintf("%d", item.Id),
              Title: item.Title,
              Arg: fmt.Sprintf("https://annict.com/works/%d", item.Id),
              Subtitle: item.MediaText+" "+item.SeasonNameText,
            })
          }

          list.Print()
        }
        return nil
      },
    },
  }

  app.Run(os.Args)
}
