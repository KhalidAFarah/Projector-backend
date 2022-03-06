package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var base string = "https://www.bungie.net/Platform"
var key string = "24101fc72fdd4f2f8ce7289e3e4ce847"
var auth string = "Bearer CPn4AxKGAgAgoZUenQOWQitwkc/U1bwzgrQj8WmUlfIhKGBFKbrKusfgAAAAsCazmcDHAkO7IZAZWed73vf4CsguYbfKCzsJo2mnkpRbQibLDNQlYSYwM5v1RtW3Fua5cwOT/4l3Zi7b7t67WV1FosJL+L6ONkmWpwUTd7V+g/HLEdEZBfGnlpMfJimWYOHW4oob+A+ftiOUsV4x7161ncUAsE4wEH2vdJh58uG4teN2hHA6QFgnsUIhojwrT5hF+90QnyabHyd7v7uNLSjykuYLaVPhJVEM6ISY9Z3rDRGz3ypZ7P7ouon/gKhUh9zPIAA6KsiQMXFVfh/imGa/dqrzbRI30x1abQg9Kks="

type RequestHeader struct {
	APIKey        string
	Authorization string
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

	fmt.Println(string(body))
	fmt.Println("------------------------------------")

	
	return body
}

type User struct {
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

	//custom
	Characters []Character `json:"characters"`
	Vault      []Item      `json:"vault"`
}

type Profile struct {
	Response struct {
		Profile struct {
			Data struct {
				CharacterIds                []string      `json:"characterIds"`
			} `json:"data"`
		} `json:"profile"`
	} `json:"Response"`
	
}

func (user *User) InitUser() (User){
	req := RequestHeader{APIKey: key, Authorization: auth}

	var userdata User

	error := json.Unmarshal(req.Send(base+"/User/GetMembershipsForCurrentUser/", "GET"), &userdata)
	if error != nil {
		log.Fatal("Unable to read data")
	}

	//"https://www.bungie.net/Platform/Destiny2/" + str(user['membershipType']) + "/Profile/" + str(user['membershipId']) + "/?components=100"
	var profiledata Profile
	newType := strconv.Itoa(userdata.Response.DestinyMemberships[0].MembershipType)
	error = json.Unmarshal(req.Send(base + "/Destiny2/" + newType + "/Profile/" + userdata.Response.DestinyMemberships[0].MembershipID + "/?components=100", "GET"), &profiledata)
	if error != nil {
		log.Fatal("Unable to read data")
	}

	for _, characterID := range profiledata.Response.Profile.Data.CharacterIds {
		var characterdata interface{}
		error = json.Unmarshal(req.Send(base + "/Destiny2/" + newType + "/Profile/" + userdata.Response.DestinyMemberships[0].MembershipID + "/Character/" + characterID + "/?components=200,205", "GET"), &characterdata)
		if error != nil {
			log.Fatal("Unable to read data")
		}
	} 


	return userdata

}

type Character struct {
	Response struct {
		Character struct {
			Data struct {
				MembershipID             string    `json:"membershipId"`
				MembershipType           int       `json:"membershipType"`
				CharacterID              string    `json:"characterId"`
				Light                    int       `json:"light"`
				Stats                    struct {
					Num144602215  int `json:"144602215"`
					Num392767087  int `json:"392767087"`
					Num1735777505 int `json:"1735777505"`
					Num1935470627 int `json:"1935470627"`
					Num1943323491 int `json:"1943323491"`
					Num2996146975 int `json:"2996146975"`
					Num4244567218 int `json:"4244567218"`
				} `json:"stats"`
				
				ClassHash            int64  `json:"classHash"`
				EmblemPath           string `json:"emblemPath"`
				EmblemBackgroundPath string `json:"emblemBackgroundPath"`
				EmblemHash           int64  `json:"emblemHash"`
				BaseCharacterLevel int     `json:"baseCharacterLevel"`
				PercentToNextLevel float64 `json:"percentToNextLevel"`
			} `json:"data"`
			Privacy int `json:"privacy"`
		} `json:"character"`
		Equipment struct {
			Data struct {
				Items []struct {
					ItemHash                   int64         `json:"itemHash"`
					ItemInstanceID             string        `json:"itemInstanceId"`
					Quantity                   int           `json:"quantity"`
					BindStatus                 int           `json:"bindStatus"`
					Location                   int           `json:"location"`
					BucketHash                 int           `json:"bucketHash"`
					TransferStatus             int           `json:"transferStatus"`
					Lockable                   bool          `json:"lockable"`
					State                      int           `json:"state"`
					DismantlePermission        int           `json:"dismantlePermission"`
					VersionNumber              int           `json:"versionNumber,omitempty"`
					OverrideStyleItemHash      int           `json:"overrideStyleItemHash,omitempty"`
				} `json:"items"`
			} `json:"data"`
		} `json:"equipment"`
		UninstancedItemComponents struct {
		} `json:"uninstancedItemComponents"`
	} `json:"Response"`
	
}

/*type Character struct {
	ID               int    `json:"id"`
	EmblemBackground string `json:"emblemBackground"`
	EmblemIcon       string `json:"emblemIcon"`
	Gear             []Item `json:"items"`
}*/

type Item struct {
	HashId int `json:"hashId"`
}
