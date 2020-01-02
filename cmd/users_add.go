package cmd

import (
	"github.com/filebrowser/filebrowser/v2/users"
	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(usersAddCmd)
	addUserFlags(usersAddCmd.Flags())
}

var usersAddCmd = &cobra.Command{
	Use:   "add <username> <password> <email> <spacen>",
	Short: "Create a new user",
	Long:  `Create a new user and add it to the database.`,
	Args:  cobra.ExactArgs(4),
	Run: python(func(cmd *cobra.Command, args []string, d pythonData) {
		s, err := d.store.Settings.Get()
		checkErr(err)
		getUserDefaults(cmd.Flags(), &s.Defaults, false)

		password, err := users.HashPwd(args[1])
		checkErr(err)

		user := &users.User{
			Username:     args[0],
			Password:     password,
			Email:        args[2],
			Space:        args[3],
			LockPassword: mustGetBool(cmd.Flags(), "lockPassword"),
		}

		s.Defaults.Apply(user)

		servSettings, err := d.store.Settings.GetServer()
		checkErr(err)
		//since getUserDefaults() polluted s.Defaults.Scope
		//which makes the Scope not the one saved in the db
		//we need the right s.Defaults.Scope here
		s2, err := d.store.Settings.Get()
		checkErr(err)

		userHome, err := s2.MakeUserDir(user.Username, user.Scope, servSettings.Root)
		checkErr(err)

		userSpace, err := s2.MakeSpaceDir(user.Space, servSettings.Root)
		checkErr(err)
		if user.Space == "" {
			user.Scope = userHome
		} else {
			user.Scope = userSpace
		}

		err = d.store.Users.Save(user)
		checkErr(err)
		printUsers([]*users.User{user})
	}, pythonConfig{}),
}
