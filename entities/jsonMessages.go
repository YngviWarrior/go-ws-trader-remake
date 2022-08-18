package entities

type Read struct {
	Id         string `json:"id"`
	Endpoint   string `json:"endpoint"`
	Error      string `json:"error"`
	Parameters struct {
		Token  string `json:"auth-token"`
		IdGame string `json:"idGame"`
	} `json:"parameters"`
}

type Write struct {
	Id           string      `json:"id"`
	Endpoint     string      `json:"endpoint"`
	Error        bool        `json:"error"`
	ErrorCode    string      `json:"Error_code"`
	ErrorMessage string      `json:"Error_message"`
	Response     GameWithBet `json:"response"`
}

type Memcached struct {
	GameIdTypeTime string `json:"game_id_type_time"`
	IdSymbolPair   string `json:"idSymbolPair"`
	IdGames        string `json:"id_games"`
	GameIdStatus   string `json:"game_status"`
}
