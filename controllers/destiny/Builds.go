package destiny

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Builds struct {
	Warlock []struct {
		Name       string   `json:"name"`
		Preference []string `json:"preference"`
		Subclass   struct {
			Item      string        `json:"item"`
			Aspects   []interface{} `json:"aspects"`
			Fragments []interface{} `json:"fragments"`
		} `json:"subclass"`
		Helmet struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"helmet"`
		Gauntlets struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"gauntlets"`
		ChestArmor struct {
			Item           string        `json:"item"`
			RecomendedMods []string      `json:"recomended_mods"`
			OptionalMods   []interface{} `json:"optional_mods"`
		} `json:"chest_armor"`
		LegArmor struct {
			Item           string        `json:"item"`
			RecomendedMods []string      `json:"recomended_mods"`
			OptionalMods   []interface{} `json:"optional_mods"`
		} `json:"leg_armor"`
		Bond struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"bond"`
	} `json:"warlock"`
	Hunter []struct {
		Name       string   `json:"name"`
		Preference []string `json:"preference"`
		Subclass   struct {
			Item      string        `json:"item"`
			Aspects   []interface{} `json:"aspects"`
			Fragments []interface{} `json:"fragments"`
		} `json:"subclass"`
		Helmet struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"helmet"`
		Gauntlets struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"gauntlets"`
		ChestArmor struct {
			Item           string        `json:"item"`
			RecomendedMods []string      `json:"recomended_mods"`
			OptionalMods   []interface{} `json:"optional_mods"`
		} `json:"chest_armor"`
		LegArmor struct {
			Item           string        `json:"item"`
			RecomendedMods []string      `json:"recomended_mods"`
			OptionalMods   []interface{} `json:"optional_mods"`
		} `json:"leg_armor"`
		Bond struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"cloak"`
	} `json:"hunter"`
	Titan []struct{
		Name       string   `json:"name"`
		Preference []string `json:"preference"`
		Subclass   struct {
			Item      string        `json:"item"`
			Aspects   []interface{} `json:"aspects"`
			Fragments []interface{} `json:"fragments"`
		} `json:"subclass"`
		Helmet struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"helmet"`
		Gauntlets struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"gauntlets"`
		ChestArmor struct {
			Item           string        `json:"item"`
			RecomendedMods []string      `json:"recomended_mods"`
			OptionalMods   []interface{} `json:"optional_mods"`
		} `json:"chest_armor"`
		LegArmor struct {
			Item           string        `json:"item"`
			RecomendedMods []string      `json:"recomended_mods"`
			OptionalMods   []interface{} `json:"optional_mods"`
		} `json:"leg_armor"`
		Bond struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"mark"`
	} `json:"titan"`
}


func GetBuilds(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonFile, err := os.Open("./resources/builds.json")
	if err != nil {
		m := Message{Response: "Unable to open file"}
		json.NewEncoder(w).Encode(m)
		return
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		m := Message{Response: "Unable to read file"}
		json.NewEncoder(w).Encode(m)
		return
	}
	var builds Builds
	err2 := json.Unmarshal(jsonData, &builds)
	if err2 != nil {
		m := Message{Response: "Unable to read json data"}
		json.NewEncoder(w).Encode(m)
		return
	}
	
	var buildsData map[string]interface{}
	buildsData["warlock"] = perChar("warlock", builds)
	//buildsData["hunter"] = perChar("hunter", builds)
	//buildsData["titan"] = perChar("titan", builds)

	json.NewEncoder(w).Encode(buildsData)
	


}

func perChar(classname string, builds Builds) ([]interface{}){

	var buildsData []interface{}
	for _, build := range builds.Warlock{
		var buildData map[string]interface{}

		//General
		buildData["name"] = build.Name
		buildData["preference"] = build.Preference

	


		//Gauntlet
		var manifestGauntletData map[string]interface{}
		_, data := DestinyManifestQuery(build.Gauntlets.Item, "DestinyInventoryItemDefinition")
		manifestGauntletData["item"] = data

		//Gauntlet recomended mods
		var recomended_mods []string
		for id := range build.Gauntlets.RecomendedMods{
			_, data = DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestGauntletData["recomended_mods"] = recomended_mods

		//Gauntlet optional mods
		var optional_mods []string
		for id := range build.Gauntlets.OptionalMods{
			_, data = DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestGauntletData["optional_mods"] = optional_mods
		
		//Adding it to object
		buildData["gauntlets"] = manifestGauntletData
		buildsData = append(buildsData, buildData)
	}
		
	return buildsData
}