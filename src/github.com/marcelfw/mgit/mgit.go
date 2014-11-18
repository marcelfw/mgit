// Copyright (c) 2014 Marcel Wouters

// Package main glues everything together :-)
package main

import (
	"fmt"
	"github.com/marcelfw/mgit/config"
	"github.com/marcelfw/mgit/engine"
	"github.com/marcelfw/mgit/repository"
	"log"
	"os"
)

// current version
var version = "0.0.1"

func main() {
	/*
		f, err := os.Create("mgit.pprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	//*/

	filterDefs := config.GetFilterDefs()
	commands := config.GetCommands()

	textCommand, flagInteractive, args, filter, ok := config.ParseCommandline(os.Args[1:], filterDefs)
	if ok == false {
		return
	}

	var curCommand repository.Command
	if curCommand, ok = commands[textCommand]; ok == false {
		log.Printf("No such command \"%s\"", textCommand)
		config.Usage(commands)
		return
	}

	if flagInteractive {
		if iCommand, ok := curCommand.(repository.InteractiveCommand); ok {
			iCommand.ForceInteractive()
		}
	}

	// Let the command initialize itself with the arguments.
	initResult := curCommand.Init(args)
	// @note no pointer receiver so for now we do this
	if newCommand, ok := initResult.(repository.Command); ok == true {
		curCommand = newCommand
	}

	if repositoryCommand, ok := curCommand.(repository.RepositoryCommand); ok {
		// Run the actual command.
		engine.RunCommand(repositoryCommand, filter)
	} else if infoCommand, ok := curCommand.(repository.InfoCommand); ok {
		fmt.Fprintln(os.Stdout, infoCommand.Output(commands, version))
	} else {
		fmt.Fprintln(os.Stdout, curCommand.Help())
	}
}
