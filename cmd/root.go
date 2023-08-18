package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/pkg/entity"
)

var rootCmd = &cobra.Command{
	Use:     "nali",
	Short:   "An offline tool for querying IP geographic information",
	Long:    `An offline tool for querying IP geographic information. ...`,
	Version: constant.Version,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		gbk, _ := cmd.Flags().GetBool("gbk")
		filePath, _ := cmd.Flags().GetString("f")

		// Read IP addresses from the specified file
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var ipList, locationList []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if gbk {
				line, _, _ = transform.String(simplifiedchinese.GBK.NewDecoder(), line)
			}
			if line := strings.TrimSpace(line); line == "quit" || line == "exit" {
				return
			}
			result := entity.ParseLine(line).ColorString()
			fmt.Println(result)
			ipList = append(ipList, line)
			locationList = append(locationList, result)
		}

		createCSVFile(ipList, locationList)
	},
}

func createCSVFile(ipList []string, locationList []string) {
	csvFilePath := "IP_attribution.csv" // Change this to the desired CSV file name
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer csvFile.Close()

	writer := bufio.NewWriter(csvFile)
	defer writer.Flush()

	// Write header
	_, _ = fmt.Fprintf(writer, "IP地址,IP归属地\n")

	// Write data
	for i, ip := range ipList {
		_, _ = fmt.Fprintf(writer, "%s,%s\n", ip, locationList[i])
	}
}

// Execute parse subcommand and run
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	rootCmd.Flags().Bool("gbk", false, "Use GBK decoder")
	rootCmd.Flags().StringP("f", "f", "1.txt", "Path to the file containing IP addresses")
}
