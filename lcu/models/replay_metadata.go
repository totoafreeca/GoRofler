package models

type ReplayMetadata struct {
	DownloadProgress int64  `json:"downloadProgress"`
	GameID           int64  `json:"gameId"`
	State            string `json:"state"`
}
