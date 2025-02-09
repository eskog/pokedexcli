package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	cachedData, exists := conf.cache.Get(conf.endpoint)
	if !exists {
		req, err := http.NewRequest("GET", conf.endpoint, nil)
		if err != nil {
			return err
		}
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		conf.cache.Add(conf.endpoint, body)
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}

	} else {
		err := json.Unmarshal(cachedData, &data)
		if err != nil {
			return err
		}
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
	req, err := http.NewRequest("GET", conf.previous, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	data := result{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	conf.endpoint = data.Next
	conf.previous = data.Previous
	//conf.previous = data.Previous

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(conf *config) error {
	result := locationarea{}
	completeEndpoint := conf.baseEndpoint + conf.subcommands[0]
	if data, exist := conf.cache.Get(completeEndpoint); !exist {
		req, err := http.NewRequest("GET", completeEndpoint, nil)
		fmt.Println("connecting to ", completeEndpoint)
		if err != nil {
			return err
		}
		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		conf.cache.Add(completeEndpoint, body)
		err = json.Unmarshal(body, &result)
		if err != nil {
			return err
		}
	} else {
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
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
	if data, exist := conf.cache.Get(apiEndpoint); !exist {
		req, err := http.NewRequest("GET", apiEndpoint, nil)
		if err != nil {
			return err
		}
		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		conf.cache.Add(apiEndpoint, body)
		err = json.Unmarshal(body, &result)
		if err != nil {
			return err
		}
	} else {
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}
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
