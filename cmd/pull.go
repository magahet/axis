package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/magahet/axis/types"
)

const (
	ExportRecordingURI = "axis-cgi/record/export/exportrecording.cgi"
	SchemaVersion      = "1"
	ExportFormat       = "matroska"
	DiskID             = "SD_DISK"
)

var (
	count     int
	maxLength int
	daytime   bool
	sunrise   string
	sunset    string
)

func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func filterRecording(r *types.Recording) (bool, error) {
	startLocal, _ := types.ToLocalTime(r.StartTime)

	// No stop time. Still recording.
	if r.StopTime == "" {
		if viper.GetBool("verbose") {
			fmt.Printf("%s has no stop time\n", startLocal)
		}
		return true, nil
	}

	// Max length
	if maxLength > 0 {
		start, err := time.Parse(time.RFC3339, r.StartTime)
		if err != nil {
			return true, errors.Wrap(err, "could not parse start time")
		}
		stop, err := time.Parse(time.RFC3339, r.StopTime)
		if err != nil {
			return true, errors.Wrap(err, "could not parse stop time")
		}

		duration := stop.Sub(start)
		if viper.GetBool("verbose") {
			fmt.Printf("%s checking duration: %s\n", startLocal, duration)
		}
		if time.Duration(maxLength) < duration {
			if viper.GetBool("verbose") {
				fmt.Printf("%s too long. duration: %s\n", startLocal, duration)
			}
			return true, nil
		}
	}

	// Daytime only
	if daytime {
		start, err := time.Parse(time.RFC3339, r.StartTime)
		if err != nil {
			return true, errors.Wrap(err, "could not parse start time")
		}
		after, err := time.Parse(time.Kitchen, sunrise)
		if err != nil {
			return true, errors.Wrap(err, "could not parse sunrise")
		}
		before, err := time.Parse(time.Kitchen, sunset)
		if err != nil {
			return true, errors.Wrap(err, "could not parse sunset")
		}

		if clockMinutes(start) < clockMinutes(after) || clockMinutes(start) > clockMinutes(before) {
			if viper.GetBool("verbose") {
				fmt.Printf("%s outside of daylight hours\n", startLocal)
			}
			return true, nil
		}
	}

	return false, nil
}

func clockMinutes(t time.Time) int {
	h, m, _ := t.Clock()
	return h*60 + m
}

func pull(cmd *cobra.Command, args []string) error {
	fmt.Println("pull called")
	host := viper.GetString("host")
	saveDir := args[0]
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		return errors.Wrap(err, "output directory does not exist")
	}

	fmt.Print("getting list of available recordings... ")
	recordings, err := types.GetRecordings(host)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return errors.Wrap(err, "could not get recordings")
	}

	fmt.Println("done")

	if viper.GetBool("verbose") {
		fmt.Println(recordings)
	}

	for i, r := range recordings.Recording {
		// Max download count
		if count > 0 && i >= count {
			break
		}

		exclude, err := filterRecording(r)
		if err != nil {
			return errors.Wrap(err, "problem checking filters")
		}
		if exclude {
			continue
		}

		startTime, _ := types.ToLocalTime(r.StartTime)
		fileName := fmt.Sprintf("%s.mkv", startTime)
		path := filepath.Join(saveDir, fileName)
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("%s already downloaded\n", fileName)
			continue
		}

		fmt.Printf("downloading %s... ", fileName)

		v := url.Values{}
		// Defaults
		v.Set("schemaversion", SchemaVersion)
		v.Set("exportformat", ExportFormat)
		v.Set("diskid", DiskID)

		v.Set("recordingid", r.RecordingID)

		url := fmt.Sprintf("http://%s/%s?%s", host, ExportRecordingURI, v.Encode())
		if err := DownloadFile(path, url); err != nil {
			fmt.Printf("error: %s\n", err)
		}

		fmt.Println("done")
	}

	return nil
}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull DOWNLOAD_PATH",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	RunE:  pull,
}

func init() {
	recordCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pullCmd.Flags().IntVarP(&count, "count", "c", 0, "Max number of videos to download (default: unlimited)")
	pullCmd.Flags().IntVarP(&maxLength, "max-length", "m", 0, "Exlude recordings longer than MAXLENGTH")
	pullCmd.Flags().BoolVarP(&daytime, "daytime", "d", false, "Exclude recordings that occur at night")
	pullCmd.Flags().StringVar(&sunrise, "sunrise", "7:00AM", "Sunrise")
	pullCmd.Flags().StringVar(&sunset, "sunset", "6:00PM", "Sunset")
}
