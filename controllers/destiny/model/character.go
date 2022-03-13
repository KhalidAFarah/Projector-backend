package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var base string = "https://www.bungie.net/Platform"
var key string = "24101fc72fdd4f2f8ce7289e3e4ce847"
var auth string = "Bearer CNf7AxKGAgAgK9jqxPZMYlpLd5j0coqa+mha3q79ww9zbpp9CwqblC3gAAAAVUURWDWrYgH5AEU79dUdaSMUqiHU1u9aSFDqKu39Z20B4T+0EV4Z+XB/VyCcM4SZ9cZwDTxd9TOfflw853mZlssB+ADT+2R7gY6yF3adZ4vbwCglKMu/63OP9aMG+zy6A6tSvm4nQqF/I5i9uSYBgU9kv8orC9Y7pzU/hm8DFZtgdg68dgCY9TSdLpmvYHArFIJqv67vyH3ofvfm67PLiHUD16/6vcPrDGa55h2gvjr9T7RvS1aKSwRxIqZzhI3yHND74GbHRN2h5ES2xac94V9FMcJdVZbjEuZKkrroCWY="

type RequestHeader struct {
	APIKey        string
	Authorization string
}

type Data struct {
	Id   string `json:"id"`
	Json string `json:"json"`
}

type Message struct {
	Response string `json:"response"`
}

func (requestheader *RequestHeader) Send(URL, method string) []byte { //
	client := http.Client{}
	request, error := http.NewRequest(method, URL, nil)
	if error != nil {
		log.Fatal("Unable to create request")
	}
	request.Header.Add("X-API-KEY", requestheader.APIKey)
	request.Header.Add("Authorization", requestheader.Authorization)

	response, error := client.Do(request)

	if error != nil {
		log.Fatal("Unable to send request to bungie")
	}

	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatal("Unable to read data")
	}

	//fmt.Println(string(body))
	//fmt.Println("------------------------------------")

	return body
}

type User struct {
	MembershipID   string `json:"membershipId"`
	MembershipType int    `json:"membershipType"`

	//custom
	Characters []Character `json:"characters"`
	Vault      []Item      `json:"vault"`

}
type Character struct {
	Light          			int    		`json:"light"`
	Stats          			interface{} `json:"stats"`
	ClassHash            	int64   	`json:"classHash"`
	EmblemPath           	string  	`json:"emblemPath"`
	EmblemBackgroundPath 	string  	`json:"emblemBackgroundPath"`
	EmblemHash           	int64   	`json:"emblemHash"`
	BaseCharacterLevel   	int     	`json:"baseCharacterLevel"`
	Inventory 				[]Item 		`json:"inventory"`
}

type UserData struct {
	Response struct {
		DestinyMemberships []struct {
			MembershipType int    `json:"membershipType"`
			MembershipID   string `json:"membershipId"`
		} `json:"destinyMemberships"`
		PrimaryMembershipID string `json:"primaryMembershipId"`
		BungieNetUser       struct {
			MembershipID string `json:"membershipId"`
			DisplayName  string `json:"displayName"`
		} `json:"bungieNetUser"`
	} `json:"Response"`
	ErrorCode       int    `json:"ErrorCode"`
	ThrottleSeconds int    `json:"ThrottleSeconds"`
	ErrorStatus     string `json:"ErrorStatus"`
	Message         string `json:"Message"`
	MessageData     struct {
	} `json:"MessageData"`
}

type ProfileData struct {
	Response struct {
		Profile struct {
			Data struct {
				CharacterIds []string `json:"characterIds"`
			} `json:"data"`
		} `json:"profile"`
	} `json:"Response"`
}

func InitUser() {
	req := RequestHeader{APIKey: key, Authorization: auth}

	var userdata UserData

	error := json.Unmarshal(req.Send(base+"/User/GetMembershipsForCurrentUser/", "GET"), &userdata)
	if error != nil {
		log.Fatal("Unable to read data")
	}

	//"https://www.bungie.net/Platform/Destiny2/" + str(user['membershipType']) + "/Profile/" + str(user['membershipId']) + "/?components=100"
	var profiledata ProfileData
	newType := strconv.Itoa(userdata.Response.DestinyMemberships[0].MembershipType)
	error = json.Unmarshal(req.Send(base+"/Destiny2/"+newType+"/Profile/"+userdata.Response.DestinyMemberships[0].MembershipID+"/?components=100", "GET"), &profiledata)
	if error != nil {
		log.Fatal("Unable to read data")
	}

	db, error := sql.Open("sqlite3", "./controllers/destiny/manifest/world_sql_content_f5d265c7cb4dc5794bc2006c58a1f33b.content")
	if error != nil {
		log.Fatal("Unable to read data")
		//m := Message{Response: "Unable to load destiny manifest file"}
		//json.NewEncoder(w).Encode(m)
		//return
	}

	for _, characterID := range profiledata.Response.Profile.Data.CharacterIds {
		var characterdata CharacterData
		error = json.Unmarshal(req.Send(base+"/Destiny2/"+newType+"/Profile/"+userdata.Response.DestinyMemberships[0].MembershipID+"/Character/"+characterID+"/?components=200,205", "GET"), &characterdata)
		if error != nil {
			log.Fatal("Unable to read data")
		}
		
		//var items []Item
		for _, item := range characterdata.Response.Equipment.Data.Items {
			newHash := strconv.Itoa(int(item.ItemHash))
			rows, error := db.Query("SELECT * FROM DestinyInventoryItemDefinition WHERE id='" + newHash + "';")
			if error != nil {
				log.Fatal("Unable to read data asd")
				//m := Message{Response: "Unable to query the destiny manifest"}
				//json.NewEncoder(w).Encode(m)
				//return
			}
			
			var itemdata Item
			var itemid int
			
			for rows.Next() {
				rows.Scan(&itemid, &itemdata)
				//items = append(items, itemdata)
				
				fmt.Println(itemdata)
				fmt.Println("-----------------")
				//break
				//data := Data{Id: strconv.Itoa(id), Json: jsondata}
				//json.NewEncoder(w).Encode(data)
			}
			
			
			
			/*var index int = 0
			var readyPerk bool = true
			//readyStat := true

			lengthPerk := len(itemdata.Perks)

			//lengthStat := len(itemdata.)
			fmt.Println(lengthPerk)
			for readyPerk {//|| readyStat {
				if index < lengthPerk {
			

					perkhash := strconv.Itoa(itemdata.Perks[index].PerkHash)
					rows, error := db.Query("SELECT * FROM DestinySandboxPerkDefinition WHERE id='" + perkhash + "';")
					if error != nil {
						log.Fatal("Unable to read data asd")
						//m := Message{Response: "Unable to query the destiny manifest"}
						//json.NewEncoder(w).Encode(m)
						//return
					}
					
					var perkdata string
					var perkid int
					
					for rows.Next() {
						rows.Scan(&perkid, &perkdata)
						//items = append(items, perkdata)
						fmt.Println(perkdata)
						
						
				
					}
					
					//itemdata.Perks[index] = perkdata
				}else {
					readyPerk = false
				}
				index++


			character : = Character{characterdata.Response.Character.Light,
				 characterdata.Response.Character.Stats,
				 characterdata.Response.Character.ClassHash
				 characterdata.Response.Character.EmblemPath
				 characterdata.Response.Character.EmblemBackgroundPath
				 characterdata.Response.Character.EmblemHash
				 characterdata.Response.Character.BaseCharacterLevel
				 items
				}
			}*/

			



		}
		

	}

	//return userdata

}

type CharacterData struct {
	Response struct {
		Character struct {
			Data struct {
				MembershipID   string `json:"membershipId"`
				MembershipType int    `json:"membershipType"`
				CharacterID    string `json:"characterId"`
				Light          int    `json:"light"`
				Stats          interface{} `json:"stats"`
				ClassHash            int64   `json:"classHash"`
				EmblemPath           string  `json:"emblemPath"`
				EmblemBackgroundPath string  `json:"emblemBackgroundPath"`
				EmblemHash           int64   `json:"emblemHash"`
				BaseCharacterLevel   int     `json:"baseCharacterLevel"`
			} `json:"data"`
		} `json:"character"`
		Equipment struct {
			Data struct {
				Items []struct {
					ItemHash              int64  `json:"itemHash"`
					ItemInstanceID        string `json:"itemInstanceId"`
					Quantity              int    `json:"quantity"`
					BindStatus            int    `json:"bindStatus"`
					Location              int    `json:"location"`
					BucketHash            int    `json:"bucketHash"`
					TransferStatus        int    `json:"transferStatus"`
					Lockable              bool   `json:"lockable"`
					State                 int    `json:"state"`
					DismantlePermission   int    `json:"dismantlePermission"`
					VersionNumber         int    `json:"versionNumber,omitempty"`
					OverrideStyleItemHash int    `json:"overrideStyleItemHash,omitempty"`
				} `json:"items"`
			} `json:"data"`
		} `json:"equipment"`
		UninstancedItemComponents struct {
		} `json:"uninstancedItemComponents"`
	} `json:"Response"`
}



type Item struct {
	DisplayProperties struct {
		Description string `json:"description"`
		Name        string `json:"name"`
		Icon        string `json:"icon"`
		HasIcon     bool   `json:"hasIcon"`
	} `json:"displayProperties"`
	CollectibleHash      int64  `json:"collectibleHash"`
	IconWatermark        string `json:"iconWatermark"`
	IconWatermarkShelved string `json:"iconWatermarkShelved"`

	Screenshot                 string `json:"screenshot"`
	ItemTypeDisplayName        string `json:"itemTypeDisplayName"`
	FlavorText                 string `json:"flavorText"`
	UIItemDisplayStyle         string `json:"uiItemDisplayStyle"`
	ItemTypeAndTierDisplayName string `json:"itemTypeAndTierDisplayName"`
	DisplaySource              string `json:"displaySource"`
	Action                     struct {
		VerbName                string        `json:"verbName"`
		VerbDescription         string        `json:"verbDescription"`
		IsPositive              bool          `json:"isPositive"`
		RequiredCooldownSeconds int           `json:"requiredCooldownSeconds"`
		RequiredItems           []interface{} `json:"requiredItems"`
		ProgressionRewards      []interface{} `json:"progressionRewards"`
		ActionTypeLabel         string        `json:"actionTypeLabel"`
		RewardSheetHash         int           `json:"rewardSheetHash"`
		RewardItemHash          int           `json:"rewardItemHash"`
		RewardSiteHash          int           `json:"rewardSiteHash"`
		RequiredCooldownHash    int           `json:"requiredCooldownHash"`
		DeleteOnAction          bool          `json:"deleteOnAction"`
		ConsumeEntireStack      bool          `json:"consumeEntireStack"`
		UseOnAcquire            bool          `json:"useOnAcquire"`
	} `json:"action"`
	Inventory struct {
		MaxStackSize                             int    `json:"maxStackSize"`
		BucketTypeHash                           int64  `json:"bucketTypeHash"`
		RecoveryBucketTypeHash                   int    `json:"recoveryBucketTypeHash"`
		TierTypeHash                             int64  `json:"tierTypeHash"`
		IsInstanceItem                           bool   `json:"isInstanceItem"`
		NonTransferrableOriginal                 bool   `json:"nonTransferrableOriginal"`
		TierTypeName                             string `json:"tierTypeName"`
		TierType                                 int    `json:"tierType"`
		ExpirationTooltip                        string `json:"expirationTooltip"`
		ExpiredInActivityMessage                 string `json:"expiredInActivityMessage"`
		ExpiredInOrbitMessage                    string `json:"expiredInOrbitMessage"`
		SuppressExpirationWhenObjectivesComplete bool   `json:"suppressExpirationWhenObjectivesComplete"`
	} `json:"inventory"`
	Stats struct {
		DisablePrimaryStatDisplay bool `json:"disablePrimaryStatDisplay"`
		StatGroupHash             int  `json:"statGroupHash"`
		Stats                     interface{} `json:"stats"`
		HasDisplayableStats bool `json:"hasDisplayableStats"`
		PrimaryBaseStatHash int  `json:"primaryBaseStatHash"`
	} `json:"stats"`
	EquippingBlock struct {
		UniqueLabelHash       int      `json:"uniqueLabelHash"`
		EquipmentSlotTypeHash int64    `json:"equipmentSlotTypeHash"`
		Attributes            int      `json:"attributes"`
		EquippingSoundHash    int      `json:"equippingSoundHash"`
		HornSoundHash         int      `json:"hornSoundHash"`
		AmmoType              int      `json:"ammoType"`
		DisplayStrings        []string `json:"displayStrings"`
	} `json:"equippingBlock"`
	TranslationBlock struct {
		WeaponPatternHash int64 `json:"weaponPatternHash"`
		Arrangements []struct {
			ClassHash          int `json:"classHash"`
			ArtArrangementHash int `json:"artArrangementHash"`
		} `json:"arrangements"`
		HasGeometry bool `json:"hasGeometry"`
	} `json:"translationBlock"`
	Preview struct {
		ScreenStyle         string `json:"screenStyle"`
		PreviewVendorHash   int    `json:"previewVendorHash"`
		PreviewActionString string `json:"previewActionString"`
	} `json:"preview"`
	Quality struct {
		ItemLevels                      []interface{} `json:"itemLevels"`
		QualityLevel                    int           `json:"qualityLevel"`
		InfusionCategoryName            string        `json:"infusionCategoryName"`
		InfusionCategoryHash            int64         `json:"infusionCategoryHash"`
		InfusionCategoryHashes          []int64       `json:"infusionCategoryHashes"`
		ProgressionLevelRequirementHash int64         `json:"progressionLevelRequirementHash"`
		CurrentVersion                  int           `json:"currentVersion"`
		Versions                        []struct {
			PowerCapHash int64 `json:"powerCapHash"`
		} `json:"versions"`
		DisplayVersionWatermarkIcons []string `json:"displayVersionWatermarkIcons"`
	} `json:"quality"`
	AcquireRewardSiteHash int `json:"acquireRewardSiteHash"`
	AcquireUnlockHash     int `json:"acquireUnlockHash"`
	Sockets               struct {
		Detail        string `json:"detail"`
		SocketEntries []struct {
			SocketTypeHash                        int64         `json:"socketTypeHash"`
			SingleInitialItemHash                 int           `json:"singleInitialItemHash"`
			ReusablePlugItems                     []interface{} `json:"reusablePlugItems"`
			PreventInitializationOnVendorPurchase bool          `json:"preventInitializationOnVendorPurchase"`
			PreventInitializationWhenVersioning   bool          `json:"preventInitializationWhenVersioning"`
			HidePerksInItemTooltip                bool          `json:"hidePerksInItemTooltip"`
			PlugSources                           int           `json:"plugSources"`
			ReusablePlugSetHash                   int           `json:"reusablePlugSetHash,omitempty"`
			OverridesUIAppearance                 bool          `json:"overridesUiAppearance"`
			DefaultVisible                        bool          `json:"defaultVisible"`
			RandomizedPlugSetHash                 int           `json:"randomizedPlugSetHash,omitempty"`
		} `json:"socketEntries"`
		IntrinsicSockets []struct {
			PlugItemHash   int64 `json:"plugItemHash"`
			SocketTypeHash int   `json:"socketTypeHash"`
			DefaultVisible bool  `json:"defaultVisible"`
		} `json:"intrinsicSockets"`
		SocketCategories []struct {
			SocketCategoryHash int64 `json:"socketCategoryHash"`
			SocketIndexes      []int `json:"socketIndexes"`
		} `json:"socketCategories"`
	} `json:"sockets"`
	TalentGrid struct {
		TalentGridHash   int    `json:"talentGridHash"`
		ItemDetailString string `json:"itemDetailString"`
		HudDamageType    int    `json:"hudDamageType"`
	} `json:"talentGrid"`
	InvestmentStats []struct {
		StatTypeHash          int  `json:"statTypeHash"`
		Value                 int  `json:"value"`
		IsConditionallyActive bool `json:"isConditionallyActive"`
	} `json:"investmentStats"`
	Perks []struct {
		RequirementDisplayString string `json:"requirementDisplayString"`
		PerkHash                 int    `json:"perkHash"`
		PerkVisibility           int    `json:"perkVisibility"`
	} `json:"perks"`
	LoreHash                          int      `json:"loreHash"`
	SummaryItemHash                   int64    `json:"summaryItemHash"`
	AllowActions                      bool     `json:"allowActions"`
	DoesPostmasterPullHaveSideEffects bool     `json:"doesPostmasterPullHaveSideEffects"`
	NonTransferrable                  bool     `json:"nonTransferrable"`
	ItemCategoryHashes                []int    `json:"itemCategoryHashes"`
	SpecialItemType                   int      `json:"specialItemType"`
	ItemType                          int      `json:"itemType"`
	ItemSubType                       int      `json:"itemSubType"`
	ClassType                         int      `json:"classType"`
	BreakerType                       int      `json:"breakerType"`
	Equippable                        bool     `json:"equippable"`
	DamageTypeHashes                  []int64  `json:"damageTypeHashes"`
	DamageTypes                       []int    `json:"damageTypes"`
	DefaultDamageType                 int      `json:"defaultDamageType"`
	DefaultDamageTypeHash             int64    `json:"defaultDamageTypeHash"`
	IsWrapper                         bool     `json:"isWrapper"`
	TraitIds                          []string `json:"traitIds"`
	TraitHashes                       []int64  `json:"traitHashes"`
	Hash                              int      `json:"hash"`
	Index                             int      `json:"index"`
	Redacted                          bool     `json:"redacted"`
	Blacklisted                       bool     `json:"blacklisted"`
}



/*type Character struct {
	ID               int    `json:"id"`
	EmblemBackground string `json:"emblemBackground"`
	EmblemIcon       string `json:"emblemIcon"`
	Gear             []Item `json:"items"`
}*/

/*type Item struct {
	HashId int `json:"hashId"`
}*/

				
type Stat struct {
	StatHash       int64 `json:"statHash"`
	Value          int   `json:"value"`
	Minimum        int   `json:"minimum"`
	Maximum        int   `json:"maximum"`
	DisplayMaximum int   `json:"displayMaximum"`
} 
/*

					Num144602215  int `json:"144602215"`
					Num392767087  int `json:"392767087"`
					Num1735777505 int `json:"1735777505"`
					Num1935470627 int `json:"1935470627"`
					Num1943323491 int `json:"1943323491"`
					Num2996146975 int `json:"2996146975"`
					Num4244567218 int `json:"4244567218"`
				
*/


