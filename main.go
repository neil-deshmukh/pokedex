package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type Pokemon struct {
	name      string
	height    float64
	weight    int
	stats     map[string]int
	poketypes []string
}

func main() {
	fmt.Println(`Heres a list of commands i built: 
	help - for more information
	map - gives you a batch of 20 locations
	mapb - gives you the previous batch of 20 location
	explore (location) - takes a argument which is a location and return all the pokemon found in that area
	catch (pokemon) - Tries to catch the given pokemon and returns if you caught it or not
	inspect (pokemon) - inspects and gives you information about the given pokemon
	pokedex - Gives you a list of all the pokemon you have caught`)
	iOfArea := 0
	pokedex := make(map[string]Pokemon)
	for {
		pokemonRreader := bufio.NewReader(os.Stdin)
		pokemon, err := pokemonRreader.ReadString('\n')
		if err != nil {
			fmt.Println("Sorry there was a error: ", err)
			continue
		}
		pokemon = strings.TrimSpace(pokemon)
		if pokemon == "help" {
			fmt.Println(`This is details about each command:
			map - gives you the next batch of 20 locations
			mapb - gives you the previous batch of 20 location. It will return an error if we are on the first batch
			explore - Takes a location as input and returns all pokemon found in that area as output
			catch - Takes a pokemon as input and tells you if you caught it or not. The stronger the pokemon lower the chance
			inspect - Takes a pokemon as input and returns you the information about the pokemon. It will return error if you have not caught the given pokemon
			pokedex - Returns a list of all the pokemon you have`)
			continue
		} else if pokemon == "exit" {
			break
		} else if pokemon == "map" {
			if iOfArea != 0 {
				iOfArea += 20
			}
			urlstr := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%v", iOfArea)
			resp, gerr := http.Get(urlstr)
			if gerr != nil {
				fmt.Println("Opps There was a error while fetching: ", gerr)
				continue
			}

			body, rerr := io.ReadAll(resp.Body)
			if rerr != nil {
				fmt.Println("Opps, There was a promlem while reading: ", rerr)
				continue
			}

			areajson := string(body)
			areaobj := make(map[string][]map[string]string)
			json.Unmarshal([]byte(areajson), &areaobj)
			areas := areaobj["results"]
			for _, area := range areas {
				fmt.Println(" - ", area["name"])
			}
		} else if pokemon == "mapb" {
			if iOfArea == 0 {
				fmt.Println("You have to go forward before you can come back")
				continue
			}
			iOfArea -= 20
			urlstr := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%v", iOfArea)
			resp, gerr := http.Get(urlstr)
			if gerr != nil {
				fmt.Println("Opps There was a error while fetching: ", gerr)
				continue
			}

			body, rerr := io.ReadAll(resp.Body)
			if rerr != nil {
				fmt.Println("Opps, There was a promlem while reading: ", rerr)
				continue
			}

			areajson := string(body)
			areaobj := make(map[string][]map[string]string)
			json.Unmarshal([]byte(areajson), &areaobj)
			areas := areaobj["results"]
			for _, area := range areas {
				fmt.Println(" - ", area["name"])
			}
		} else if strings.Contains(pokemon, "explore") {
			inis := strings.Split(pokemon, " ")
			fmt.Println("Exploring " + inis[1] + ".....")
			urlstr := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", inis[1])
			resp, err := http.Get(urlstr)
			if err != nil {
				fmt.Println("Opps, There was a error: ", err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Opps, There was a error: ", err)
			}

			areaDataJSON := string(body)
			areaData := make(map[string][]map[string]map[string]string)
			json.Unmarshal([]byte(areaDataJSON), &areaData)
			pokemonobjs := areaData["pokemon_encounters"]
			for _, pokeobj := range pokemonobjs {
				fmt.Println(" - ", pokeobj["pokemon"]["name"])
			}
		} else if strings.Contains(pokemon, "catch") {
			poke := strings.Split(pokemon, " ")
			urlstr := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", poke[1])
			resp, gerr := http.Get(urlstr)
			if gerr != nil {
				fmt.Println("Opps, There was a error: ", gerr)
			}

			body, rerr := io.ReadAll(resp.Body)
			if rerr != nil {
				fmt.Println("Opps, there was a error: ", rerr)
			}

			pokejson := string(body)
			pokeobjs := make(map[string]string)
			pokeobji := make(map[string]int)
			pokeobjsi := make(map[string][]map[string]int)
			pokeobjso := make(map[string][]map[string]map[string]string)
			json.Unmarshal([]byte(pokejson), &pokeobjs)
			json.Unmarshal([]byte(pokejson), &pokeobji)
			json.Unmarshal([]byte(pokejson), &pokeobjsi)
			json.Unmarshal([]byte(pokejson), &pokeobjso)
			chance := rand.Intn(pokeobji["base_experience"])
			fmt.Println("Throwing a pokeball at " + pokeobjs["name"] + ".....")
			if chance <= 100 {
				fmt.Println("You caught " + pokeobjs["name"] + "!!!")
				orgstats := make(map[string]int)
				keys := make([]string, 0)
				values := make([]int, 0)
				for _, statobj := range pokeobjso["stats"] {
					keys = append(keys, statobj["stat"]["name"])
				}
				for _, statobj := range pokeobjsi["stats"] {
					values = append(values, statobj["base_stat"])
				}
				for i := range keys {
					orgstats[keys[i]] = values[i]
				}
				orgtypes := make([]string, 0)
				for _, typeobj := range pokeobjso["types"] {
					orgtypes = append(orgtypes, typeobj["type"]["name"])
				}
				pokedex[pokeobjs["name"]] = Pokemon{
					name:      pokeobjs["name"],
					height:    float64(pokeobji["height"]),
					weight:    pokeobji["weight"],
					stats:     orgstats,
					poketypes: orgtypes,
				}
				continue
			}
			fmt.Println(pokeobjs["name"] + " escaped!")
		} else if strings.Contains(pokemon, "inspect") {
			inis := strings.Split(pokemon, " ")
			pokeobj, ok := pokedex[inis[1]]
			if !ok {
				fmt.Println("You do not have this pokemon")
				continue
			}
			fmt.Println("Name: ", pokeobj.name)
			fmt.Println("Height: ", pokeobj.height/10)
			fmt.Println("Weight: ", pokeobj.weight/10)
			fmt.Println("Stats: ")
			for key, value := range pokeobj.stats {
				fmt.Println(" -", key, ": ", value)
			}
			fmt.Println("Types: ")
			for _, typ := range pokeobj.poketypes {
				fmt.Println(" -", typ)
			}
		} else if pokemon == "pokedex" {
			fmt.Println("Your Pokedex: ")
			for key := range pokedex {
				fmt.Println(" - " + key)
			}
		}
	}
}
