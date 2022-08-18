package mysql

import (
	"fmt"
	"gowstrader/entities"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (c *SqlConn) GetGame(gameID int64) (row entities.GameWithBet, err error) {
	sqlrows, err := c.conn.Query(`
		SELECT 	bt.id as id_time, bt.name as name_time, bt.time_done_seconds, bt.time_calc_seconds,
				bg.id as id_game, bg.id_moedas_pares as symbol_id, m.symbol as symbol_name, bg.game_id_status, 
				bg.game_date_start, bg.game_date_process, bg.game_date_finish, 
				bg.game_profit_percent

		FROM binary_option_game bg
		JOIN moedas_pares m ON m.id = bg.id_moedas_pares
		JOIN binary_option_game_time bt ON bt.id = bg.game_id_type_time
		WHERE bg.id = ?
	`, gameID)

	if err != nil {
		return
	}

	defer sqlrows.Close()

	for sqlrows.Next() {
		err := sqlrows.Scan(&row.TimeID, &row.NameTime, &row.TimeDoneSeconds, &row.TimeCalcSeconds, &row.GameID, &row.SymbolID, &row.SymbolName, &row.GameIDStatus, &row.GameDateStart, &row.GameDateProcess, &row.GameDateFinish, &row.GameProfitPercent)
		if err != nil {
			log.Println("fetch err: ", err)
		}
	}

	return
}

func (c *SqlConn) GetInfoGameWithBet(gameID int64, userID int64) (rows []entities.GameWithBet, err error) {
	sqlrows, err := c.conn.Query(`
	    SELECT bt.id as id_time, bt.name as name_time, bt.time_done_seconds, bt.time_calc_seconds

			,bg.id as id_game, bg.id_moedas_pares as symbol_id, m.symbol as symbol_name, bg.game_id_status
			,bg.game_date_start, bg.game_date_process, bg.game_date_finish, bg.game_profit_percent
			,(UNIX_TIMESTAMP(bg.game_date_start) * 1000) as game_date_start_timestamp
			,(UNIX_TIMESTAMP(bg.game_date_process) * 1000) as game_date_process_timestamp
			,(UNIX_TIMESTAMP(bg.game_date_finish) * 1000) as game_date_finish_timestamp
			,IF(bg.game_id_status = ?,bg.game_price_amount_selected_finish,0) as close_price

			,bb.id as id_bet, bb.id_choice, bb.bet_amount_dolar, bb.price_amount_selected
			,IF(bb.refund = 1,0,IF(bb.amount_win_dolar > 0,bb.amount_win_dolar,(bb.bet_amount_dolar * -1))) as amount_win_dolar
			,IF(bb.refund = 1,3,IF(bb.amount_win_dolar > 0,2,1)) as bet_status
			,IF(bb.id_balance = ?,0,1) as bet_practice, bb.date_register as bet_date_register
			
		FROM binary_option_game bg 
		JOIN moedas_pares m ON m.id = bg.id_moedas_pares
		JOIN binary_option_game_time bt ON bt.id = bg.game_id_type_time
		LEFT JOIN binary_option_game_bet bb 
		ON bb.id_game = bg.id 
		AND bb.id_usuario = ?

		WHERE bg.id = ?
	`, entities.GameStatusDone, entities.GameSaldoTipoJogar, userID, gameID)

	if err != nil {
		return
	}

	defer sqlrows.Close()

	for sqlrows.Next() {
		row := entities.GameWithBet{}
		err := sqlrows.Scan(&row.TimeID, &row.NameTime, &row.TimeDoneSeconds, &row.TimeCalcSeconds, &row.GameID, &row.SymbolID, &row.SymbolName, &row.GameIDStatus, &row.GameDateStart, &row.GameDateProcess, &row.GameDateFinish, &row.GameProfitPercent, &row.GameDateStartTimestamp, &row.GameDateProcessTimeStamp, &row.GameDateFinishTimeStamp, &row.GameClosePrice, &row.BetID, &row.BetIDChoice, &row.BetAmountDolar, &row.BetPriceClose, &row.BetAmountWinDolar, &row.BetStatus, &row.BetPractice, &row.BetDateRegister)
		if err != nil {
			break
		}
		rows = append(rows, row)
	}

	return
}

func (c *SqlConn) GetInfoGameListWithBet(gameIDs []int64, userID int64) (rows []*entities.GameWithBet, err error) {
	var in string
	for i, id := range gameIDs {
		if i+1 == len(gameIDs) {
			in += fmt.Sprintf("%v", id)
		} else {
			in += fmt.Sprintf("%v,", id)
		}
	}

	sqlrows, err := c.conn.Query(`
	    SELECT bt.id as id_time, bt.name as name_time, bt.time_done_seconds, bt.time_calc_seconds
			,bg.id as id_game, bg.id_moedas_pares as symbol_id, m.symbol as symbol_name, bg.game_id_status
			,bg.game_date_start, bg.game_date_process, bg.game_date_finish, bg.game_profit_percent
			,(UNIX_TIMESTAMP(bg.game_date_start) * 1000) as game_date_start_timestamp
			,(UNIX_TIMESTAMP(bg.game_date_process) * 1000) as game_date_process_timestamp
			,(UNIX_TIMESTAMP(bg.game_date_finish) * 1000) as game_date_finish_timestamp
			,IF(bg.game_id_status = ?,bg.game_price_amount_selected_finish,0) as close_price

			,bb.id as id_bet, bb.id_choice, bb.bet_amount_dolar, bb.price_amount_selected
			,IF(bb.refund = 1,0,IF(bb.amount_win_dolar > 0,bb.amount_win_dolar,(bb.bet_amount_dolar * -1))) as amount_win_dolar
			,IF(bb.refund = 1,3,IF(bb.amount_win_dolar > 0,2,1)) as bet_status
			,IF(bb.id_balance = ?,0,1) as bet_practice, bb.date_register as bet_date_register
			
		FROM binary_option_game bg 
		JOIN moedas_pares m ON m.id = bg.id_moedas_pares
		JOIN binary_option_game_time bt ON bt.id = bg.game_id_type_time
		LEFT JOIN binary_option_game_bet bb 
		ON bb.id_game = bg.id 
		AND bb.id_usuario = ?
		WHERE bg.id IN (?)
	`, entities.GameStatusDone, 3, userID, in)

	if err != nil {
		return
	}

	defer sqlrows.Close()

	for sqlrows.Next() {
		row := entities.GameWithBet{}
		err := sqlrows.Scan(&row.TimeID, &row.NameTime, &row.TimeDoneSeconds, &row.TimeCalcSeconds, &row.GameID,
			&row.SymbolID, &row.SymbolName, &row.GameIDStatus, &row.GameDateStart, &row.GameDateProcess, &row.GameDateFinish,
			&row.GameProfitPercent, &row.GameDateStartTimestamp, &row.GameDateProcessTimeStamp, &row.GameDateFinishTimeStamp,
			&row.GameClosePrice, &row.BetID, &row.BetIDChoice, &row.BetAmountDolar, &row.BetPriceClose, &row.BetAmountWinDolar,
			&row.BetStatus, &row.BetPractice, &row.BetDateRegister)

		if err != nil {
			break
		}

		rows = append(rows, &row)
	}

	return
}
