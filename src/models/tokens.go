package models

type Token struct {
	Address       string `json:"address"`
	Name          string `json:"name"`
	CreationBlock int    `json:"creation_block"`
}

type Tokens []Token
