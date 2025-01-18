package main

import "github.com/eskog/pokedexcli/internal/pokecache"

type config struct {
	endpoint    string
	next        string
	previous    string
	cache       *pokecache.Cache
	subcommands []string
}

type NamedAPIResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationarea struct {
	Id                     int                    `json:"id"`
	Name                   string                 `json:"name"`
	Game_index             int                    `json:"game_index"`
	Encounter_method_rates []encountermethodrates `json:"encounter_method_rates"`
	Location               NamedAPIResource       `json:"location"`
	Names                  []name                 `json:"names"`
	Pokemon_encounters     []pokemonEncounter     `json:"pokemon_encounters"`
}

type encountermethodrates struct {
	Encounter_method NamedAPIResource          `json:"encounter_method"`
	Version_details  []encounterVersionDetails `json:"version_details"`
}

type pokemonEncounter struct {
	Pokemon         NamedAPIResource         `json:"pokemon"`
	Version_details []versionEncounterDetail `json:"version_details"`
}

type encounterVersionDetails struct {
	Rate    int              `json:"rate"`
	Version NamedAPIResource `json:"version"`
}

type versionEncounterDetail struct {
	Version           NamedAPIResource `json:"version"`
	Max_chance        int              `json:"max_chance"`
	Encounter_details []encounter      `json:"encounter_details"`
}

type encounter struct {
	Min_level        int                `json:"min_level"`
	Max_level        int                `json:"max_level"`
	Condition_values []NamedAPIResource `json:"condition_values"`
	Chance           int                `json:"chance"`
	Method           NamedAPIResource   `json:"method"`
}

type name struct {
	Name     string           `json:"name"`
	Language NamedAPIResource `json:"language"`
}

type result struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}
