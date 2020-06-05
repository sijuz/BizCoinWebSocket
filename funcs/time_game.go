package funcs

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"BizCoinWebSocket/additionally"
	"BizCoinWebSocket/config"

	// MySql Driver
	_ "github.com/go-sql-driver/mysql"
)

// TimeGame is Websocket /time_game path function
func TimeGame(w http.ResponseWriter, r *http.Request) {
	var upgrader = config.Upgrader
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in TimeGame method: ", err)
		return
	}
	defer c.Close()

	// тут пойдёт проверка подписи, наверное
	// статус отправлю такой же, потом переправлю, если чё

	err = c.WriteJSON(config.ConfirmData{
		Status:           "connected",
		Channel:          "time_game",
		UseVkSignChecker: config.Conf.UseVkSignChecker,
	})
	if err != nil {
		log.Println("Error in sending connected status (time_game):", err)
		return
	}

	log.Println("TimeGame Connection request!")

	vkUserID := 0
	signResult := false

	for {
		var data config.Data

		err = c.ReadJSON(&data)
		if err != nil {
			log.Println("Message read error in update: ", err)
			err = c.WriteJSON(config.Error[0])
			if err != nil {
				log.Println("it's an error in error message sending: ", err)
			}
			return
		}

		if config.Conf.UseVkSignChecker == true {
			if data.Sign != "" {
				if additionally.IsValid(data, config.Conf.VkAppSecret) {
					if data.VkUserID != "" && data.VkUserID != "0" {
						vkUserID, err = strconv.Atoi(data.VkUserID)

						signResult = true
						err = c.WriteJSON(config.SignCheck{
							CheckSignStatus:  "OK",
							UseVksignChecker: config.Conf.UseVkSignChecker,
						})

						if err != nil {
							log.Println("error in sending singChecker status (time_game):", err)
							return
						}

						break
					} else {
						signResult = false
						err = c.WriteJSON(config.Error[1])
						if err != nil {
							log.Println("error in sending singChecker status (time_game):", err)
							return
						}

						break
					}
				} else {
					signResult = false
					err = c.WriteJSON(config.Error[1])
					if err != nil {
						log.Println("error in sending singChecker status (time_game):", err)
						return
					}
					break
				}
			} else {
				signResult = false
				err = c.WriteJSON(config.Error[1])
				if err != nil {
					log.Println("error in sending singChecker status (time_game):", err)
					return
				}
				break
			}
		} else {
			signResult = true
			vkUserID, err = strconv.Atoi(data.VkUserID)
			err = c.WriteJSON(config.SignCheck{
				CheckSignStatus:  "OK",
				UseVksignChecker: config.Conf.UseVkSignChecker,
			})
			if err != nil {
				log.Println("error in sending singChecker status (time_game):", err)
				return
			}
			break
		}
	}

	log.Println("New socket connected to TimeGame!")

	db, err := sql.Open(
		"mysql",
		config.Conf.DbUser+":"+
			config.Conf.DbPassword+
			"@tcp("+config.Conf.DbHostPort+")/"+
			config.Conf.DbName,
	)
	if err != nil {
		log.Println("error in opening db (time_game):", err)
		return
	}

	defer func() {
		_ = db.Close()
	}()

	for signResult {

		var data config.GameData // data to load

		err = c.ReadJSON(&data) // read data
		if err != nil {
			log.Println("error in reading action JSON (time_game):", err)
			return
		}

		if data.Action == "gamestart" && data.Status == "start_game" { // now game is starting

			log.Println("Game Started!")

			var balance int // user's balance
			{
				row, err := db.Query(
					"select balance_bizcoin "+
						"from user "+
						"where user_id = ?",
					vkUserID,
				) // db request
				if err != nil {
					// panic(err)
					log.Println("error in getting balance (time_game):", err)
					return
				}
				was := false
				for row.Next() {
					err = row.Scan(&balance)
					if err != nil {
						// panic(err)
						log.Println("error in scaning balance (time_game):", err)
						return
					}
					was = true
				} // get user's balance
				if !was { // send an error in case, it's an unknown user
					err = c.WriteJSON(config.Error[2])
					if err != nil {
						log.Println("error in sending an error(2) (time_game):", err)
						return
					}
				}
			}

			log.Println(balance)

			rand.Seed(time.Now().UnixNano()) // randomise random :)

			if fir, sec := rand.Intn(10)+1, rand.Intn(10)+1; fir == sec {
				toPlus := int(config.Conf.MinProfit + rand.Intn(config.Conf.MaxProfit-config.Conf.MinProfit))
				balance += toPlus

				if _, err = db.Exec(
					"update user set balance_bizcoin = ? where user_id = ?",
					balance,
					vkUserID,
				); err != nil {
					log.Print("error in changing balance")
				}

				if err = c.WriteJSON(map[string]interface{}{
					"action":  "win_checking",
					"status":  "win",
					"balance": balance,
					"to_plus": toPlus,
				}); err != nil {
					log.Println("error in sending win status (time_game):", err)
					return
				}
			} else if fir, sec := rand.Intn(20)+1, rand.Intn(20)+1; fir == sec {

				toMines := int(config.Conf.MinLoss + rand.Intn(config.Conf.MaxLoss-config.Conf.MinLoss))

				if balance < toMines {
					if err = c.WriteJSON(config.Error[5]); err != nil {
						log.Println("error in sending an error(5) (time_game):", err)
						return
					}
					log.Println("not enough money")
					continue
				}

				balance -= toMines

				if _, err = db.Exec(
					"update user set balance_bizcoin = ? where user_id = ?",
					balance,
					vkUserID,
				); err != nil {
					log.Print("error in changing balance")
				}

				if err = c.WriteJSON(map[string]interface{}{
					"action":   "win_checking",
					"status":   "defeat",
					"balance":  balance,
					"to_mines": toMines,
				}); err != nil {
					log.Println("error in sending defeat status (time_game):", err)
					return
				}

			} else {
				if err = c.WriteJSON(map[string]interface{}{
					"action":  "win_checking",
					"status":  "none",
					"balance": balance,
				}); err != nil {
					log.Println("error in sending none status (time_game):", err)
					return
				}
			}

		} else {

			err := c.WriteJSON(config.Error[4])
			if err != nil {
				log.Println("error in sending an error(4) (time_game):", err)
				return
			}

		}
	}
}
