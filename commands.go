package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"math/rand"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandMap(conf *config) error {
	data := result{}

	body, err := makeAPICall(conf, conf.endpoint)
	if err != nil {
		fmt.Printf("Error calling api or cache: %s", err)
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	conf.endpoint = data.Next
	conf.previous = data.Previous

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(conf *config) error {
	if conf.previous == "" {
		fmt.Println("You are on the first page")
		return nil
	}
	data := result{}
	body, err := makeAPICall(conf, conf.previous)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	conf.endpoint = data.Next
	conf.previous = data.Previous

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(conf *config) error {
	result := locationarea{}
	completeEndpoint := conf.baseEndpoint + conf.subcommands[0]
	body, err := makeAPICall(conf, completeEndpoint)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	for _, pokemon := range result.Pokemon_encounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(conf *config) error {
	result := pokemon{}
	fmt.Printf("Throwing a Pokeball at %s...\n", conf.subcommands[0])
	apiEndpoint := conf.pokemonEndpoint + conf.subcommands[0]
	body, err := makeAPICall(conf, apiEndpoint)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	catchRoll := rand.Intn(100)
	for {
		if result.Base_experience > 100 {
			result.Base_experience -= result.Base_experience / 10
		} else {
			break
		}
	}
	if catchRoll >= result.Base_experience {
		fmt.Printf("You caught %s\n", conf.subcommands[0])
		conf.pokemon[conf.subcommands[0]] = result
	} else {
		fmt.Printf("%s escaped\n", conf.subcommands[0])
	}
	return nil
}

func commandInspect(conf *config) error {
	pokemon, exists := conf.pokemon[conf.subcommands[0]]
	if !exists {
		return errors.New("pokemon not caught")
	}
	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(conf *config) error {
	if len(conf.pokemon) < 1 {
		return errors.New("you have not caught any pokemons")
	}
	fmt.Println("Your pokedex:")
	for _, pokemon := range conf.pokemon {
		fmt.Println("-", pokemon.Name)
	}
	return nil
}
