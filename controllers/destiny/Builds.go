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
		ClassArmor struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"class_armor"`
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
		ClassArmor struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"class_armor"`
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
		ClassArmor struct {
			Item           string   `json:"item"`
			RecomendedMods []string `json:"recomended_mods"`
			OptionalMods   []string `json:"optional_mods"`
		} `json:"class_armor"`
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
	
	buildsData := make(map[string]interface{})
	buildsData["warlock"] = perChar("warlock", builds)
	//buildsData["hunter"] = perChar("hunter", builds)
	//buildsData["titan"] = perChar("titan", builds)

	json.NewEncoder(w).Encode(buildsData)
	


}

func perChar(classname string, builds Builds) ([]interface{}){
	var buildsData []interface{}
	for _, build := range builds.Warlock{
		buildData :=  make(map[string]interface{})

		//General
		buildData["name"] = build.Name
		buildData["preference"] = build.Preference


		//subclass
		manifestItemData := make(map[string]interface{})

		if build.Subclass.Item != "" {
			_, data := DestinyManifestQuery(build.Subclass.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Subclass aspects
		var aspects []string
		for id := range build.Subclass.Aspects{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			aspects = append(aspects, data)
		}
		manifestItemData["aspects"] = aspects

		//Subclass fragments
		var fragments []string
		for id := range build.Subclass.Fragments{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			fragments = append(fragments, data)
		}
		manifestItemData["fragments"] = fragments

		//Adding it to object
		buildData["subclass"] = manifestItemData

		//Helmet
		if build.Helmet.Item != "" {
			manifestItemData := make(map[string]interface{})
			_, data := DestinyManifestQuery(build.Helmet.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Helmet recomended mods
		var recomended_mods []string
		for id := range build.Helmet.RecomendedMods{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Helmet optional mods
		var optional_mods []string
		for id := range build.Helmet.OptionalMods{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods

		//Adding it to object
		buildData["helmet"] = manifestItemData


		//Gauntlet
		if build.Gauntlets.Item != ""{
			manifestItemData := make(map[string]interface{})
			_, data := DestinyManifestQuery(build.Gauntlets.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Gauntlet recomended mods
		recomended_mods = make([]string, 0)
		for _, id := range build.Gauntlets.RecomendedMods{
			_, data := DestinyManifestQuery(id, "DestinySandboxPerkDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Gauntlet optional mods
		optional_mods = make([]string, 0)
		for _, id := range build.Gauntlets.OptionalMods{
			_, data := DestinyManifestQuery(id, "DestinySandboxPerkDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods
		
		//Adding it to object
		buildData["gauntlets"] = manifestItemData


		//Chest armor
		if build.ChestArmor.Item != ""{
			manifestItemData := make(map[string]interface{})
			_, data := DestinyManifestQuery(build.ChestArmor.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Chest armor recomended mods
		recomended_mods = make([]string, 0)
		for id := range build.ChestArmor.RecomendedMods{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Chest armor optional mods
		optional_mods = make([]string, 0)
		for id := range build.ChestArmor.OptionalMods{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods

		//Adding it to object
		buildData["chest_armor"] = manifestItemData


		//Leg armor
		if build.LegArmor.Item != "" {
			manifestItemData := make(map[string]interface{})
			_, data := DestinyManifestQuery(build.LegArmor.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Leg armor recomended mods
		recomended_mods = make([]string, 0)
		for id := range build.LegArmor.RecomendedMods{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Leg armor optional mods
		optional_mods = make([]string, 0)
		for id := range build.LegArmor.OptionalMods{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods

		//Adding it to object
		buildData["leg_armor"] = manifestItemData

		//Class armor
		if build.ClassArmor.Item != "" {
			manifestItemData := make(map[string]interface{})
			_, data := DestinyManifestQuery(build.ClassArmor.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Class armor recomended mods
		recomended_mods = make([]string, 0)
		for id := range build.ClassArmor.RecomendedMods{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Class armor optional mods
		optional_mods = make([]string, 0)
		for id := range build.ClassArmor.OptionalMods{
			_, data := DestinyManifestQuery(strconv.Itoa(id), "DestinySandboxPerkDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods

		//Adding it to object
		buildData["class_armor"] = manifestItemData


		//Putting it all together in an object
		buildsData = append(buildsData, buildData)
	}
		
	return buildsData
}