/*
	Copyright 2017-2018 OneLedger

	Cli to init a node (server)
*/
package main

import (
	"fmt"
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"os"
	"runtime/debug"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize node (server)",
	Run:   InitNode,
}

type InitCmdArguments struct {
	Password string
	NewPassword string
}

var initCmdArguments *InitCmdArguments = &InitCmdArguments{}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&initCmdArguments.Password, "password", "", "existing node password")
	initCmd.Flags().StringVar(&initCmdArguments.NewPassword, "newpassword", "", "new node password")
}

func InitNode(cmd *cobra.Command, args []string) {

	// Catch any underlying panics, for now just print out the details properly and stop
	defer func() {
		if r := recover(); r != nil {
			log.Error("Fullnode Fatal Panic, shutting down", "r", r)
			debug.PrintStack()
			if service != nil {
				service.Stop()
			}
			os.Exit(-1)
		}
	}()

	log.Debug("Initializing", "appAddress", global.Current.AppAddress, "on", global.Current.NodeName)

	shouldReplacePassword := false

	newPlainPassword := initCmdArguments.NewPassword
	currentPlainPassword := initCmdArguments.Password

	node := app.NewApplication()

	adminPassword := node.GetPassword()

	if adminPassword == nil {
		shouldReplacePassword = true
	}

	if adminPassword != nil {
		if currentPlainPassword == "" {
			tty := shared.Tty{}

			currentPlainPassword = tty.Password("Enter a password:")
		}

		err := bcrypt.CompareHashAndPassword(adminPassword.([]byte), []byte(currentPlainPassword))

		if err != nil {
			log.Fatal("Wrong password", "error", err)
			return
		}

		// TODO were already initialized, nothing to do now?
		return
	}

	if shouldReplacePassword {
		if newPlainPassword == "" {
			tty := shared.Tty{}

			for true {
				newPlainPassword = tty.Password("Enter a new password:")

				// @TODO need some actual password policy rules here or maybe in another place
				if len(newPlainPassword) < 4 {
					fmt.Println("Password should be longer than 4 characters")
					continue
				}

				passwordConfirm := tty.Password("Confirm a new password:")

				if newPlainPassword != passwordConfirm {
					fmt.Println("Passwords don't match")
					continue
				}

				break
			}
		} else {
			// @TODO need some actual password policy rules here or maybe in another place
			if len(newPlainPassword) < 4 {
				log.Fatal("Password should be longer than 4 characters")
				return
			}
		}

		passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(newPlainPassword), bcrypt.DefaultCost)

		if err != nil {
			log.Fatal("Can't encrypt password", "error", err)
		}

		session := node.Admin.Begin()
		session.Set(data.DatabaseKey("Password"), passwordEncrypted)
		session.Commit()
	}
}
