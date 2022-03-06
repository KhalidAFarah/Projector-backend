package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var base string = "https://www.bungie.net/Platform"
var key string = "24101fc72fdd4f2f8ce7289e3e4ce847"
var auth string = "Bearer COL4AxKGAgAgGy1GRBIVzJJk2ND/NhqEhXjpsmmiogZ6OaCPW6/d9nvgAAAA9xlcHoj/0IJnH8uYLhkCiGNIaz+oCvqKhPSFG717lujHBB3sfM/7viVvd3gZxt7V/4M2mef+5oiG090mCY9gYCqqdxbwtfTDlylp9s9J/IxvAHbFZeTwLIV30EiV3CJloZU0mklOI1tgUWABW+HdY8xny7BI+rYN/c3pvzEyaVXvdg9qYQt+ja3iA73+VHgpTd2E5eq8pnF+1+B5MYnHKdaGKk/Zq0JkcbGkOOOCgKwMAp2KavFhgfV/lLqyROh5vEVtKgLZHlvp+WXd0xTEWua62PEU5wXX0bwnUDmN0do="

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

func (user *User) Init() {
	req := RequestHeader{APIKey: key, Authorization: auth}

	var data interface{}

	error := json.Unmarshal(req.Send(base+"/User/GetMembershipsForCurrentUser/", "GET"), &data)
	if error != nil {
		log.Fatal("Unable to read data")
	}

	fmt.Println(data)

}

type Character struct {
	ID               int    `json:"id"`
	EmblemBackground string `json:"emblemBackground"`
	EmblemIcon       string `json:"emblemIcon"`
	Gear             []Item `json:"items"`
}

type Item struct {
	HashId int `json:"hashId"`
}
