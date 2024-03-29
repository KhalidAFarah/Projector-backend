package destiny

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	//"strconv"
)

type Class struct {
	Name       string   `json:"name"`
	Preference []string `json:"preference"`
	Subclass   struct {
		Item      string        `json:"item"`
		Aspects   []interface{} `json:"aspects"`
		Fragments []interface{} `json:"fragments"`
	} `json:"subclass"`
	Kinetic struct {
		Item            string        `json:"item"`
		RecomendedPerks []interface{} `json:"recomended_perks"`
	} `json:"kinetic"`
	Energy struct {
		Item            string        `json:"item"`
		RecomendedPerks []interface{} `json:"recomended_perks"`
	} `json:"energy"`
	Heavy struct {
		Item            string        `json:"item"`
		RecomendedPerks []interface{} `json:"recomended_perks"`
	} `json:"heavy"`
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
	var builds map[string][]Class
	err2 := json.Unmarshal(jsonData, &builds)
	if err2 != nil {
		m := Message{Response: "Unable to read json data"}
		json.NewEncoder(w).Encode(m)
		return
	}

	buildsData := make(map[string]interface{})
	buildsData["warlock"] = perChar(builds["warlock"])
	//buildsData["hunter"] = perChar(builds["hunter"])
	//buildsData["titan"] = perChar(builds["titan"])

	json.NewEncoder(w).Encode(buildsData)

}

func perChar(builds []Class) []interface{} {
	var buildsData []interface{}
	for _, build := range builds {
		buildData := make(map[string]interface{})

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
		aspects := make([]Item, 0)
		for _, id := range build.Subclass.Aspects {
			_, data := DestinyManifestQuery(id.(string), "DestinyInventoryItemDefinition")
			aspects = append(aspects, data)
		}
		manifestItemData["aspects"] = aspects

		//Subclass fragments
		fragments := make([]Item, 0)
		for _, id := range build.Subclass.Fragments {
			_, data := DestinyManifestQuery(id.(string), "DestinyInventoryItemDefinition")
			fragments = append(fragments, data)
		}
		manifestItemData["fragments"] = fragments

		//Adding it to object
		buildData["subclass"] = manifestItemData

		//Primary
		manifestItemData = make(map[string]interface{})
		if build.Kinetic.Item != "" {
			_, data := DestinyManifestQuery(build.Kinetic.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Recomended perks for kinetic
		recomended_perks := make([]Item, 0)
		for _, id := range build.Kinetic.RecomendedPerks {
			_, data := DestinyManifestQuery(id.(string), "DestinyInventoryItemDefinition")
			recomended_perks = append(recomended_perks, data)
		}
		manifestItemData["recomended_perks"] = recomended_perks

		//Adding it to object
		buildData["kinetic"] = manifestItemData

		//Energy
		manifestItemData = make(map[string]interface{})
		if build.Energy.Item != "" {
			_, data := DestinyManifestQuery(build.Energy.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Recomended perks for energy
		recomended_perks = make([]Item, 0)
		for _, id := range build.Energy.RecomendedPerks {
			_, data := DestinyManifestQuery(id.(string), "DestinyInventoryItemDefinition")
			recomended_perks = append(recomended_perks, data)
		}
		manifestItemData["recomended_perks"] = recomended_perks

		//Adding it to object
		buildData["energy"] = manifestItemData

		//Heavy
		manifestItemData = make(map[string]interface{})
		if build.Heavy.Item != "" {
			_, data := DestinyManifestQuery(build.Heavy.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Recomended perks for energy
		recomended_perks = make([]Item, 0)
		for _, id := range build.Heavy.RecomendedPerks {
			_, data := DestinyManifestQuery(id.(string), "DestinyInventoryItemDefinition")
			recomended_perks = append(recomended_perks, data)
		}
		manifestItemData["recomended_perks"] = recomended_perks

		//Adding it to object
		buildData["heavy"] = manifestItemData

		//Helmet
		manifestItemData = make(map[string]interface{})
		if build.Helmet.Item != "" {
			_, data := DestinyManifestQuery(build.Helmet.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Helmet recomended mods
		var recomended_mods []Item
		for _, id := range build.Helmet.RecomendedMods {
			_, data := DestinyManifestQuery(id, "DestinyInventoryItemDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Helmet optional mods
		var optional_mods []Item
		for _, id := range build.Helmet.OptionalMods {
			_, data := DestinyManifestQuery(id, "DestinyInventoryItemDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods

		//Adding it to object
		buildData["helmet"] = manifestItemData

		//Gauntlet
		manifestItemData = make(map[string]interface{})
		if build.Gauntlets.Item != "" {

			_, data := DestinyManifestQuery(build.Gauntlets.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Gauntlet recomended mods
		recomended_mods = make([]Item, 0)
		for _, id := range build.Gauntlets.RecomendedMods {
			_, data := DestinyManifestQuery(id, "DestinyInventoryItemDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Gauntlet optional mods
		optional_mods = make([]Item, 0)
		for _, id := range build.Gauntlets.OptionalMods {
			_, data := DestinyManifestQuery(id, "DestinyInventoryItemDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods

		//Adding it to object
		buildData["gauntlets"] = manifestItemData

		//Chest armor
		manifestItemData = make(map[string]interface{})
		if build.ChestArmor.Item != "" {

			_, data := DestinyManifestQuery(build.ChestArmor.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Chest armor recomended mods
		recomended_mods = make([]Item, 0)
		for _, id := range build.ChestArmor.RecomendedMods {
			_, data := DestinyManifestQuery(id, "DestinyInventoryItemDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Chest armor optional mods
		optional_mods = make([]Item, 0)
		for _, id := range build.ChestArmor.OptionalMods {
			_, data := DestinyManifestQuery(id.(string), "DestinyInventoryItemDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods

		//Adding it to object
		buildData["chest_armor"] = manifestItemData

		//Leg armor
		manifestItemData = make(map[string]interface{})

		if build.LegArmor.Item != "" {
			_, data := DestinyManifestQuery(build.LegArmor.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Leg armor recomended mods
		recomended_mods = make([]Item, 0)
		for _, id := range build.LegArmor.RecomendedMods {
			_, data := DestinyManifestQuery(id, "DestinyInventoryItemDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Leg armor optional mods
		optional_mods = make([]Item, 0)
		for _, id := range build.LegArmor.OptionalMods {
			_, data := DestinyManifestQuery(id.(string), "DestinyInventoryItemDefinition")
			optional_mods = append(optional_mods, data)
		}
		manifestItemData["optional_mods"] = optional_mods

		//Adding it to object
		buildData["leg_armor"] = manifestItemData

		//Class armor
		manifestItemData = make(map[string]interface{})

		if build.ClassArmor.Item != "" {
			_, data := DestinyManifestQuery(build.ClassArmor.Item, "DestinyInventoryItemDefinition")
			manifestItemData["item"] = data
		}

		//Class armor recomended mods
		recomended_mods = make([]Item, 0)
		for _, id := range build.ClassArmor.RecomendedMods {
			_, data := DestinyManifestQuery(id, "DestinyInventoryItemDefinition")
			recomended_mods = append(recomended_mods, data)
		}
		manifestItemData["recomended_mods"] = recomended_mods

		//Class armor optional mods
		optional_mods = make([]Item, 0)
		for _, id := range build.ClassArmor.OptionalMods {
			_, data := DestinyManifestQuery(id, "DestinyInventoryItemDefinition")
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
