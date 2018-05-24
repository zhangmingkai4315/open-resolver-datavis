package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zhangmingkai4315/open-resolver-datavis/config"
)

func init() {
	resolverCmd.PersistentFlags().String("path:src", "", "raw ip address from zmap(53)")
	resolverCmd.PersistentFlags().String("path:dest", "open-resolver.ip", "raw ip address from zmap(53)")
	rootCmd.AddCommand(resolverCmd)
}

var resolverCmd = &cobra.Command{
	Use:   "resolver",
	Short: "resolver will remove any ip from file not responese to dns requery or just refuse",
	Run: func(cmd *cobra.Command, args []string) {
		rawIPFile := cmd.Flag("path:src").Value.String()
		if rawIPFile == "" {
			log.Printf("File %s not set, please use --path:src to load file first", rawIPFile)
			os.Exit(1)
		}
		output := cmd.Flag("path:dest").Value.String()
		log.Printf("Try send dns data to candidte open reslover and save to %s", output)
		err := collectOpenDNSResolver(rawIPFile, output)
		if err != nil {
			log.Printf("Error:%s", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func queryWorker(id int, jobs <-chan string, results chan<- string) {
	for ip := range jobs {
		// do query job
		results <- ip
	}
}

func collectOpenDNSResolver(srcFile string, destFile string) error {
	// open src file read million lines
	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// create dest file and write results to this file
	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	// create worker to query dns and save result to channes
	scanner := bufio.NewScanner(f)
	cfg := config.GetAppConfig()
	worker := cfg.GlobalConfig.Worker
	jobs := make(chan string, 100)
	results := make(chan string, 100)

	for w := 1; w <= worker; w++ {
		go queryWorker(w, jobs, results)
	}

	go func(dest *os.File, results chan string) {
		w := bufio.NewWriter(dest)
		for r := range results {
			fmt.Fprintln(w, r)
		}
		w.Flush()
	}(dest, results)

	for scanner.Scan() {
		line := scanner.Text()
		jobs <- line
	}

	close(jobs)
	close(results)
	return nil
}
