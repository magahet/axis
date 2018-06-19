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

		i := 0
		for _, r := range recordings.Recording {
			exclude, err := r.Filter(maxLength, daytime, sunrise, sunset)
			if err != nil {
				return errors.Wrap(err, "problem checking filters")
			}
			if exclude {
				continue
			}
			fmt.Println(r)
			// Max download count
			i += 1
			if count > 0 && i >= count {
				break
			}
		}
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
	lsCmd.Flags().IntVarP(&count, "count", "c", 0, "Max number of videos to download (default: unlimited)")
	lsCmd.Flags().StringVarP(&maxLength, "max-length", "m", "10m", "Exlude recordings longer than MAXLENGTH")
	lsCmd.Flags().BoolVarP(&daytime, "daytime", "d", false, "Exclude recordings that occur at night")
	lsCmd.Flags().StringVar(&sunrise, "sunrise", "7:00AM", "Sunrise")
	lsCmd.Flags().StringVar(&sunset, "sunset", "6:00PM", "Sunset")
}
