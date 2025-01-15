package main

import "github.com/eskog/pokedexcli/internal/pokecache"

type config struct {
	endpoint string
	next     string
	previous string
	cache    *pokecache.Cache
}

type NamedAPIResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type result struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}
