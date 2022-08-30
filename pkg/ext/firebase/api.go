package firebase

import (
	"context"
	"fmt"
	"time"
	core "videosmover/pkg"
	"videosmover/pkg/web"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/lcserny/goutils"
	"google.golang.org/api/option"
)

func NewCloudDatabase(cfg *core.ProxyConfig) web.CloudDatabase {
	return &firebaseCloudDatabase{
		databaseUrl:    cfg.CloudDBUrl,
		accountKeyFile: cfg.CloudDBAccountKeyfile,
		serverName:     cfg.ServerName,
	}
}

type firebaseCloudDatabase struct {
	databaseUrl    string
	accountKeyFile string
	serverName     string
}

func (f firebaseCloudDatabase) Init() {
	client, err := f.createClient()
	if err != nil {
		goutils.LogError(fmt.Errorf("error creating Firebase client: %v", err))
		return
	}

	serverRef := client.NewRef("servers").Child(f.serverName)
	if err := serverRef.Child("actionsPending").Set(context.Background(), []string{}); err != nil {
		goutils.LogError(err)
	}

	go func(serverRef *db.Ref) {
		for range time.NewTicker(time.Second * 10).C {
			if err := serverRef.Child("lastPingDate").Set(context.Background(), time.Now().UnixMilli()); err != nil {
				goutils.LogError(err)
				return
			}
	
			commands := make([]string, len(core.AvailableCommands))
			i := 0
			for command := range core.AvailableCommands {
				commands[i] = command
				i++
			}
			if err := serverRef.Child("actionsAvailable").Set(context.Background(), commands); err != nil {
				goutils.LogError(err)
				return
			}
	
			var actionsPending []string
			if err := serverRef.Child("actionsPending").Get(context.Background(), &actionsPending); err != nil {
				goutils.LogError(err)
				return
			}
	
			if len(actionsPending) > 0 {
				action := actionsPending[0]
				if err := serverRef.Child("actionsPending").Set(context.Background(), actionsPending[1:]); err != nil {
					goutils.LogError(err)
					return
				}

				if len(action) > 0 {
					if f, ok := core.AvailableCommands[action]; ok {
						f()
					}
				}
			}
		}
	}(serverRef)
}

func (f firebaseCloudDatabase) createClient() (*db.Client, error) {
	opt := option.WithCredentialsFile(f.accountKeyFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.DatabaseWithURL(context.Background(), f.databaseUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}
