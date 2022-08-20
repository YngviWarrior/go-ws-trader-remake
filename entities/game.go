package entities

import "database/sql"

const GameStatusDone = 4
const GameStatusCalc = 3
const GameStatusPending = 1
const GameSaldoTipoJogar = 3

type GameWithBet struct {
	TimeID                   int64           `json:"id_time"`
	NameTime                 string          `json:"name_time"`
	TimeDoneSeconds          int64           `json:"time_done_seconds"`
	TimeCalcSeconds          int64           `json:"time_calc_seconds"`
	GameID                   int64           `json:"id_game"`
	SymbolID                 int64           `json:"symbol_id"`
	SymbolName               string          `json:"symbol_name"`
	GameIDStatus             int64           `json:"game_id_status"`
	GameDateStart            string          `json:"game_date_start"`
	GameDateProcess          string          `json:"game_date_process"`
	GameDateFinish           string          `json:"game_date_finish"`
	GameProfitPercent        float64         `json:"game_profit_percent"`
	GameDateStartTimestamp   int64           `json:"game_date_start_timestamp"`
	GameDateProcessTimeStamp int64           `json:"game_date_process_timestamp"`
	GameDateFinishTimeStamp  int64           `json:"game_date_finish_timestamp"`
	GameClosePrice           float64         `json:"close_price"`
	BetID                    sql.NullInt64   `json:"id_bet"`
	BetIDChoice              sql.NullInt64   `json:"id_choice"`
	BetAmountDolar           sql.NullFloat64 `json:"bet_amount_dolar"`
	BetPriceClose            sql.NullFloat64 `json:"price_amount_selected"`
	BetAmountWinDolar        sql.NullFloat64 `json:"amount_win_dolar"`
	BetStatus                sql.NullInt64   `json:"bet_status"`
	BetPractice              sql.NullInt64   `json:"bet_practice"`
	BetDateRegister          sql.NullString  `json:"bet_date_register"`
}
