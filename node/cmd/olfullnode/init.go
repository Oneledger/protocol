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
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	"golang.org/x/crypto/bcrypt"
	"os"
	"path/filepath"
	"runtime/debug"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize node (server)",
	RunE:  initNode,
}

type InitCmdArguments struct {
	password    string
	newPassword string
	genesis     string
	folder      string
}

var initCmdArguments *InitCmdArguments = &InitCmdArguments{}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&initCmdArguments.password, "password", "", "existing node password")
	initCmd.Flags().StringVar(&initCmdArguments.newPassword, "newpassword", "", "new node password")
	initCmd.Flags().StringVar(&initCmdArguments.genesis, "genesis", "", "Gensis file to use to generate new node key file")
	initCmd.Flags().StringVar(&initCmdArguments.folder, "dir", "./", "Directory to store initialization files for the node, default current folder")
}

func initNode(cmd *cobra.Command, _ []string) error {
	args := initCmdArguments
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

	genesisdoc, err := types.GenesisDocFromFile(filepath.Join(args.folder, args.genesis))
	if err != nil {
		return err
	}
	dir := filepath.Join(args.folder, "consensus", "config")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = genesisdoc.SaveAs(filepath.Join(dir, "genesis.json"))
	if err != nil {
		return err
	}
	// Make node key
	_, err = p2p.LoadOrGenNodeKey(filepath.Join(dir, "node_key.json"))
	if err != nil {
		return err
	}
	// Make private validator file
	pvFile := privval.GenFilePV(filepath.Join(dir, "priv_validator_key.json"), filepath.Join(dir, "../data/priv_validator_state.json"))
	pvFile.Save()

	return nil
}

func setupPasswod() {

	log.Debug("Setup Password")
	shouldReplacePassword := false

	newPlainPassword := initCmdArguments.newPassword
	currentPlainPassword := initCmdArguments.password

	node := app.NewApplication()
	node.Initialize()

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
			}
		}

		passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(newPlainPassword), bcrypt.DefaultCost)

		if err != nil {
			log.Fatal("Can't encrypt password", "error", err)
			return
		}

		session := node.Admin.Begin()
		session.Set(data.DatabaseKey("Password"), passwordEncrypted)
		session.Commit()
	}
}
