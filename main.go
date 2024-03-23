package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/sprpgmr/papermc-fetch/files"
	paperapi "github.com/sprpgmr/papermc-fetch/paper-api"
)

type programArgs struct {
	Experimental bool   `long:"experimental" description:"check for experimental builds"`
	Filename     string `short:"f" long:"file" description:"file to output to" value-name:"FILE" default:"paper.jar"`
	SkipDownload bool   `long:"skip-download" description:"skip downloading files"`
	Prefix       string `short:"p" long:"prefix" description:"only look for builds containing this version prefix"`
}

func main() {
	service := paperapi.GetPaperAPIService()
	fileService := files.GetFileService()

	err := runMainProgram(service, fileService, os.Args)
	if err != nil && !flags.WroteHelp(err) {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func runMainProgram(paperAPIService paperapi.Service, fileService files.Service, args []string) error {
	opts := &programArgs{}
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		return err
	}

	fmt.Printf("Checking for latest version of paper...\n")

	buildInfo, err := paperAPIService.GetLatestBuild(opts.Experimental, opts.Prefix)
	if err != nil {
		return err
	}

	if buildInfo == nil {
		return errors.New("no builds found")
	}

	msg := fmt.Sprintf("Latest paper version is %s - build #%d", buildInfo.Version, buildInfo.Build)
	if buildInfo.Channel != "default" {
		msg += " EXPERIMENTAL"
	}

	fmt.Println(msg)

	exists, err := paperAPIService.DownloadExists(opts.Filename, buildInfo)
	if err != nil {
		return err
	}

	if exists {
		fmt.Println("You already have this version of paper.")
		return nil
	}

	if !opts.SkipDownload {
		err = fileService.DeleteIfExists(opts.Filename)
		if err != nil {
			return err
		}
	}

	if opts.SkipDownload {
		return nil
	}

	fmt.Println("Downloading...")

	err = paperAPIService.DownloadJar(buildInfo, opts.Filename)
	if err != nil {
		return err
	}

	fmt.Println("Finished downloading.")

	fmt.Println("Verifying file integrity...")
	valid, err := paperAPIService.IsValidDownload(opts.Filename, buildInfo.Downloads.Application.Sha256)
	if err != nil {
		return err
	}

	if !valid {
		fmt.Println("Download is invalid!!")
		return errors.New("download invalid")
	}

	fmt.Println("Download verified.")

	return nil
}
