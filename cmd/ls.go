package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/magahet/axis/types"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List available recordings",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("ls called")
		host := viper.GetString("host")
		recordings, err := types.GetRecordings(host)
		if err != nil {
			fmt.Println(err)
			return errors.Wrap(err, "could not get recordings")
		}

		fmt.Println(recordings)
		return nil
	},
}

func init() {
	recordCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
