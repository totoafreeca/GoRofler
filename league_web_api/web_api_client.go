package leaguewebapi

import (
	"GoRofler/league_web_api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	summonerV4 = "lol/summoner/v4/summoners/by-name"
	matchesV5  = "lol/match/v5/matches/by-puuid/{puuid}/ids"
)

//need to add rest of the regions
var regionsMapping = map[string]string{
	"EUN1": "europe",
	"EUW":  "europe",
}

type leagueWebApiClient struct {
	ApiKey     string
	PlatformID string
	HTTPClient *http.Client
}

func getSummonerUrl(platformId string, summonerNameString string) string {
	return fmt.Sprintf("https://%s.api.riotgames.com/%s/%s", platformId, summonerV4, summonerNameString)
}

func getMatchesUrl(platformId string, puuid string) string {
	return fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids", regionsMapping[platformId], puuid)
}

func NewClient(apiKey string, platformId string) *leagueWebApiClient {
	return &leagueWebApiClient{
		ApiKey:     apiKey,
		PlatformID: platformId,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (webClient *leagueWebApiClient) GetSummonerByName(summonerName string) (models.Summoner, error) {
	emptyJsonData := []byte("")
	summonerUrl := getSummonerUrl(webClient.PlatformID, summonerName)
	request, err := http.NewRequest("GET", summonerUrl, bytes.NewBuffer(emptyJsonData))
	request.Header.Add("Accept-Charset", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("X-Riot-Token", webClient.ApiKey)
	request.Header.Add("Accept", "application/json")
	var summoner models.Summoner

	response, err := webClient.HTTPClient.Do(request)

	if err != nil {
		return summoner, err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&summoner); err != nil {
		log.Fatal("Couldn't decode replay metadata response")
	}

	return summoner, nil
}

func (webClient *leagueWebApiClient) GetMatchHistoryByPuuid(puuid string) ([]string, error) {
	emptyJsonData := []byte("")
	matchesUrl := getMatchesUrl(webClient.PlatformID, puuid)
	request, err := http.NewRequest("GET", matchesUrl, bytes.NewBuffer(emptyJsonData))
	request.Header.Add("X-Riot-Token", webClient.ApiKey)
	request.Header.Add("Accept-Charset", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Add("Accept", "application/json")
	var matches []string

	response, err := webClient.HTTPClient.Do(request)

	if err != nil {
		return matches, err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&matches); err != nil {
		log.Fatal("Couldn't decode replay metadata response")
	}

	return matches, nil
}

func MatchIdsToInts(matchIds []string) []int {
	convertedIds := []int{}
	for _, matchId := range matchIds {
		matchNumber, err := strconv.Atoi(matchId[len(matchId)-10:])
		if err != nil {
			log.Fatal(err)
		}
		convertedIds = append(convertedIds, matchNumber)
	}
	return convertedIds
}
