package cli

import (
	"context"
	"io"
	"os"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/howeyc/gopass"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type userFromFile struct {
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	Email       string   `yaml:"email"`
	Permissions []string `yaml:"permissions"`
}

func AddUserFromFile(client api.UserServiceClient, ctx context.Context, fn string) error {
	fc, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer fc.Close()

	var user userFromFile
	decoder := yaml.NewDecoder(fc)
	for {
		if err := decoder.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if _, err := client.Add(ctx, &api.AddUserRequest{Username: user.Username, Password: user.Password, Email: user.Email}); err != nil {
			return err
		}

		permissions := make([]api.Permission, len(user.Permissions))
		for idx, perm := range user.Permissions {
			p := types.Permission(perm)
			permissions[idx] = api.ConvertPermission(p)
		}
		if _, err := client.Grant(ctx, &api.GrantPermissionsRequest{Username: user.Username, Permission: permissions}); err != nil {
			return err
		}

		log.WithField("name", user.Username).WithField("permissions", permissions).Debug("Added user")
	}

	return nil
}

func GetPassword(cmd *cobra.Command) (string, error) {
	viper.BindPFlag("password", cmd.Flags().Lookup("password"))
	viper.BindEnv("password", "RURUKU_PASSWORD")
	password := viper.GetString("password")
	if password == "" {
		os.Stderr.WriteString("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return "", err
		}
		password = string(pass)
		os.Stderr.WriteString("\n")
	}
	return password, nil
}
