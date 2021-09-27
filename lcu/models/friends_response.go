package models

type FriendsResponse []struct {
	SummonerID int    `json:"summonerId"`
	Name       string `json:"name"`
	Puuid      string `json:"puuid"`
}

type FriendResponseFull []struct {
	Availability            string `json:"availability"`
	DisplayGroupID          int    `json:"displayGroupId"`
	DisplayGroupName        string `json:"displayGroupName"`
	GameName                string `json:"gameName"`
	GameTag                 string `json:"gameTag"`
	GroupID                 int    `json:"groupId"`
	GroupName               string `json:"groupName"`
	Icon                    int    `json:"icon"`
	ID                      string `json:"id"`
	IsP2PConversationMuted  bool   `json:"isP2PConversationMuted"`
	LastSeenOnlineTimestamp string `json:"lastSeenOnlineTimestamp"`
	Lol                     struct {
		AdditionalProp1 string `json:"additionalProp1"`
		AdditionalProp2 string `json:"additionalProp2"`
		AdditionalProp3 string `json:"additionalProp3"`
	} `json:"lol"`
	Name          string `json:"name"`
	Note          string `json:"note"`
	Patchline     string `json:"patchline"`
	Pid           string `json:"pid"`
	PlatformID    string `json:"platformId"`
	Product       string `json:"product"`
	ProductName   string `json:"productName"`
	Puuid         string `json:"puuid"`
	StatusMessage string `json:"statusMessage"`
	Summary       string `json:"summary"`
	SummonerID    int    `json:"summonerId"`
	Time          int    `json:"time"`
}
