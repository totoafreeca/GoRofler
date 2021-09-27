package lcu

import (
	"GoRofler/lcu/models"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	basePath             = ""
	matchHistoryEndpoint = "lol-match-history/v1/products/lol/current-summoner/matches"
	metadataEndpoint     = "lol-replays/v1/metadata/"
	downloadDemoEndpoint = ""
)

type lcuHttpClient struct {
	Port       int
	Password   string
	HTTPClient *http.Client
}

func NewClient(port int, password string) *lcuHttpClient {
	return &lcuHttpClient{
		Port:     port,
		Password: password,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Minute,
		},
	}
}

func (lc *lcuHttpClient) GetBaseUrl() string {
	return fmt.Sprintf("https://riot:%s@127.0.0.1:%d/", lc.Password, lc.Port)
}

func (lc *lcuHttpClient) GetMatchHistoryUrl() string {
	return fmt.Sprintf("%slol-match-history/v1/products/lol/current-summoner/matches", lc.GetBaseUrl())
}

func (lc *lcuHttpClient) GetMetadataEndpointUrl(gameId int) string {
	return fmt.Sprintf("%slol-replays/v1/metadata/%d", lc.GetBaseUrl(), gameId)
}

func (lc *lcuHttpClient) GetDownloadUrl(gameId int) string {
	return fmt.Sprintf("%slol-replays/v1/rofls/%d/download", lc.GetBaseUrl(), gameId)
}

func (lc *lcuHttpClient) GetMatchHistory() models.MatchHistoryResponse {
	response, err := lc.HTTPClient.Get(lc.GetMatchHistoryUrl())
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var matchHistory models.MatchHistoryResponse

	if err := json.NewDecoder(response.Body).Decode(&matchHistory); err != nil {
		log.Fatal("Couldn't decode match history response")
	}

	return matchHistory
}

func (lc *lcuHttpClient) GetReplayMetadata(gameId int) models.ReplayMetadata {
	response, err := lc.HTTPClient.Get(lc.GetMetadataEndpointUrl(gameId))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var replayMetadata models.ReplayMetadata
	if err := json.NewDecoder(response.Body).Decode(&replayMetadata); err != nil {
		log.Fatal("Couldn't decode replay metadata response")
	}

	return replayMetadata
}

func (lc *lcuHttpClient) DownloadReplay(gameId int) error {

	emptyJsonData := []byte("{}")
	request, err := http.NewRequest("POST", lc.GetDownloadUrl(gameId), bytes.NewBuffer(emptyJsonData))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	// headerss =  {'Content-Type': "application/json", 'Accept': "application/json"}
	if err != nil {
		log.Fatal(err)
	}
	response, err := lc.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode/200 == 1 {
		return nil
	}
	return err

}

func (lc *lcuHttpClient) GetFriendsSummonerIds() models.FriendsResponse {
	response, err := lc.HTTPClient.Get(lc.GetBaseUrl() + "lol-chat/v1/friends")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var chatFriendsIds models.FriendsResponse
	if err := json.NewDecoder(response.Body).Decode(&chatFriendsIds); err != nil {
		log.Fatal("Couldn't decode replay metadata response")
	}

	return chatFriendsIds
}

func (lc *lcuHttpClient) GetPlayerMatchHistory(summId int) models.MatchHistoryResponse {
	response, err := lc.HTTPClient.Get(lc.GetBaseUrl() + fmt.Sprintf("lol-match-history/v3/matchlist/account/%d", summId))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var matchHistory models.MatchHistoryResponse

	if err := json.NewDecoder(response.Body).Decode(&matchHistory); err != nil {
		log.Fatal("Couldn't decode match history response")
	}

	return matchHistory
}

func (lc *lcuHttpClient) GetPlatformId() string {
	response, err := lc.HTTPClient.Get(lc.GetBaseUrl() + fmt.Sprintf("lol-platform-config/v1/namespaces/LoginDataPacket/platformId"))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var platformId string

	if err := json.NewDecoder(response.Body).Decode(&platformId); err != nil {
		log.Fatal("Couldn't decode match history response")
	}

	return platformId
}

func (lc *lcuHttpClient) DownloadMatchHistory(matchHistory models.MatchHistoryResponse) error {

	for _, game := range matchHistory.Games.Games {

		replayMetadata := lc.GetReplayMetadata(int(game.GameID))
		//420 - ranked solo , 430 - ranked flex
		if replayMetadata.State == "download" && (game.QueueID == 420 || game.QueueID == 440) {
			lc.DownloadReplay(int(game.GameID))
		}
	}
	return nil
}
