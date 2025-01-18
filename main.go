package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/eskog/pokedexcli/internal/pokecache"
)

func main() {
	Commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "We need alot of help running this",
			callback:    commandHelp,
		},
		"map": {
			name:        "help",
			description: "print next page of map areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "print previous page of map areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a region",
			callback:    commandExplore,
		},
	}
	conf := config{
		endpoint: "https://pokeapi.co/api/v2/location-area/",
		next:     "",
		previous: "",
		cache:    pokecache.NewCache(time.Second * 30),
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		input := cleanInput(scanner.Text())
		if len(input) > 1 {
			conf.subcommands = append(conf.subcommands, input[1])
		}
		command, exist := Commands[input[0]]
		if exist {
			if err := command.callback(&conf); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
		fmt.Print("Pokedex > ")
	}
}
