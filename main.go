package main

import (
	"GoRofler/lcu"
	leaguewebapi "GoRofler/league_web_api"
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	programFiles := flag.String("client_path", "", "path to client")
	apiKey := flag.String("api_key", "", "Development api key")

	flag.Parse()

	file, err := os.Open(*programFiles + "\\lockfile")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	lockfileInfo, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}

	splittedLock := strings.Split(string(lockfileInfo), ":")

	clientPort, err := strconv.Atoi(splittedLock[2])
	if err != nil {
		log.Fatal(err)
	}
	clientPassword := splittedLock[3]

	lcuClient := lcu.NewClient(clientPort, clientPassword)
	platformId := lcuClient.GetPlatformId()

	webClient := leaguewebapi.NewClient(*apiKey, platformId)

	for _, friend := range lcuClient.GetFriendsSummonerIds() {

		summoner, err := webClient.GetSummonerByName(friend.Name)
		if err != nil {
			log.Fatal(err)
		}
		matches, err := webClient.GetMatchHistoryByPuuid(summoner.Puuid)
		intMatches := leaguewebapi.MatchIdsToInts(matches)

		for _, match := range intMatches {
			lcuClient.DownloadReplay(match)
		}
	}

}
