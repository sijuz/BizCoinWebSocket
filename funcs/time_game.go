package funcs

import (
	"database/sql"
	"github.com/gorilla/websocket"
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

type connClass struct {
	Connection *websocket.Conn
	UserID int
}

var checkUserID = make(chan connClass)
var connectionClose = make(chan int)

func ConnList(close chan struct{}) {

	connections := make(map[int]*websocket.Conn)

	for {
		select {
		case conn := <-checkUserID:
			was := false
			for key, _ := range connections {
				if key == conn.UserID {
					if err := connections[key].WriteJSON(config.Error[6]); err != nil {
						additionally.SendError("err := connections[key].WriteJSON(config.Error[6]); err != nil", err)
					}
					if err := connections[key].Close(); err != nil {
						additionally.SendError("if err := connections[key].Close(); err != nil {", err)
						//break
					}
					log.Println("User", key, "started new session! Previous connection interrupted!")
					delete(connections, conn.UserID)
					// it's something like go-crutch, instead of remove :)
					was = true
					break
				}
			}
			if !was {
				connections[conn.UserID] = conn.Connection
				log.Println("User", conn.UserID, "added successfully!")
			}

		case conn := <-connectionClose:
			if err := connections[conn].Close(); err != nil {
				additionally.SendError("if err := connections[conn].Close(); err != nil {\n\t\t\t\tadditionally.SendError(\"Error in sending disconnect status:\", err)\n\t\t\t\t//break\n\t\t\t}", err)
				//break
			}
			delete(connections, conn)
		case <-close:
			return
		}
	}
}


// TimeGame is Websocket /time_game path function
func TimeGame(w http.ResponseWriter, r *http.Request) {

	var upgrader = config.Upgrader
	//ws, err := upgrader.Upgrade(w, r, nil)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in TimeGame method: ", err)
		additionally.SendError("c, err := upgrader.Upgrade(w, r, nil) ", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	err = c.WriteJSON(config.ConfirmData{
		Status:           "connected",
		Channel:          "time_game",
		UseVkSignChecker: config.Conf.UseVkSignChecker,
	})

	if err != nil {
		log.Println("Error in sending connected status (time_game):", err)
		additionally.SendError("err = c.WriteJSON(config.ConfirmData{\n\t\tStatus:           \"connected\",\n\t\tChannel:          \"time_game\",\n\t\tUseVkSignChecker: config.Conf.UseVkSignChecker,\n\t})", err)
		return
	}

	log.Println("TimeGame Connection request!")
//}
//
//func timeGame(c *websocket.Conn) {

	vkUserID := 0
	signResult := false

	for {
		if err := c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
			log.Println("Error in Setting timeout!", err)
			additionally.SendError("if err := c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {\n\t\t\tlog.Println(\"Error in Setting timeout!\", err)\n\t\t\tadditionally.SendError(\"Error in Setting timeout:\", err)\n\t\t}", err)
		}

		var data config.Data
		err := c.ReadJSON(&data)
		if err != nil {
			log.Println("Message read error in update: ", err)
			additionally.SendError("err := c.ReadJSON(&data)", err)
			err = c.WriteJSON(config.Error[0])
			if err != nil {
				additionally.SendError("it's an error in error message sending: ", err)
				log.Println("err = c.WriteJSON(config.Error[0])", err)
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

						checkUserID <- connClass{
							Connection: c,
							UserID: vkUserID,
						}

						if err != nil {
							log.Println("error in sending singChecker status (time_game):", err)
							additionally.SendError("err = c.WriteJSON(config.Error[1])", err)
							return
						}

						break
					} else {
						signResult = false
						err = c.WriteJSON(config.Error[1])
						if err != nil {
							log.Println("error in sending singChecker status (time_game):", err)
							additionally.SendError("err = c.WriteJSON(config.Error[1])", err)
							return
						}

						break
					}
				} else {
					signResult = false
					err = c.WriteJSON(config.Error[1])
					if err != nil {
						log.Println("error in sending singChecker status (time_game):", err)
						additionally.SendError("err = c.WriteJSON(config.Error[1])", err)
						return
					}
					break
				}
			} else {
				signResult = false
				err = c.WriteJSON(config.Error[1])
				if err != nil {
					log.Println("error in sending singChecker status (time_game):", err)
					additionally.SendError("err = c.WriteJSON(config.Error[1])", err)
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
				additionally.SendError("err = c.WriteJSON(config.SignCheck{\n\t\t\t\tCheckSignStatus:  \"OK\",\n\t\t\t\tUseVksignChecker: config.Conf.UseVkSignChecker,\n\t\t\t})", err)
				return
			}

			checkUserID <- connClass{
				Connection: c,
				UserID: vkUserID,
			}

			break
		}
	}
	defer func(close chan<-int){
		close <-vkUserID
	}(nil)

	if !signResult {
		return
	}

	log.Println("New socket connected to TimeGame! ID:", vkUserID)

	db, err := sql.Open(
		"mysql",
		config.Conf.DbUser+":"+
			config.Conf.DbPassword+
			"@tcp("+config.Conf.DbHostPort+")/"+
			config.Conf.DbName,
	)
	if err != nil {
		log.Println("error in opening db (time_game):", err)
		additionally.SendError("db, err := sql.Open(\n\t\t\"mysql\",\n\t\tconfig.Conf.DbUser+\":\"+\n\t\t\tconfig.Conf.DbPassword+\n\t\t\t\"@tcp(\"+config.Conf.DbHostPort+\")/\"+\n\t\t\tconfig.Conf.DbName,\n\t)", err)
		return
	}

	defer func() {
		_ = db.Close()
	}()

	for signResult {
		var data config.GameData // data to load

		if err = c.SetReadDeadline(time.Now().Add(time.Hour)); err != nil {
			log.Println("Error in Setting timeout!", err)
			additionally.SendError("if err = c.SetReadDeadline(time.Now().Add(time.Hour)); err != nil {\n\t\t\tlog.Println(\"Error in Setting timeout!\", err)\n\t\t\tadditionally.SendError(\"Error in Setting timeout! (time_name)\", err)\n\t\t\treturn\n\t\t}", err)
			return
		}

		err = c.ReadJSON(&data) // read data

		if err != nil {
			log.Println("error in reading action JSON (time_game):", err)
			additionally.SendError("error in reading action JSON (time_game):", err)
			return
		}
		//clickTime := time.Now().Unix()

		if data.Action == "gamestart" && data.Status == "start_game" { // now game is starting

			var balance int // user's balance
			{
				row, err := db.Query(
					"select balance_bizcoin "+
						"from user "+
						"where user_id = ?",
					vkUserID,
				) // db request
				if err != nil {
					log.Println("error in getting balance (time_game):", err)
					return
				}
				was := false
				for row.Next() {
					err = row.Scan(&balance)
					if err != nil {
						log.Println("error in scanning balance (time_game):", err)
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

			rand.Seed(time.Now().UnixNano()) // randomise random :)

			if fir, sec := rand.Intn(10)+1, rand.Intn(10)+1; fir == sec {
				toPlus := int(config.Conf.MinProfit + rand.Intn(config.Conf.MaxProfit-config.Conf.MinProfit))
				balance += toPlus

				if _, err = db.Exec(
					"update user set balance_bizcoin = ? where user_id = ?",
					balance,
					vkUserID,
				); err != nil {
					additionally.SendError("error in changing balance", err)
					log.Print("error in changing balance")
				}

				if err = c.WriteJSON(map[string]interface{}{
					"action":  "win_checking",
					"status":  "win",
					"balance": balance,
					"to_plus": toPlus,
				}); err != nil {
					additionally.SendError("error in sending win status (time_game):", err)
					log.Println("error in sending win status (time_game):", err)
					return
				}
			} else if fir, sec := rand.Intn(20)+1, rand.Intn(20)+1; fir == sec {

				toMines := int(config.Conf.MinLoss + rand.Intn(config.Conf.MaxLoss-config.Conf.MinLoss))

				if balance < toMines {
					if err = c.WriteJSON(config.Error[5]); err != nil {
						additionally.SendError("error in sending an error(5) (time_game):", err)
						log.Println("error in sending an error(5) (time_game):", err)
						return
					}
					continue
				}

				balance -= toMines

				if _, err = db.Exec(
					"update user set balance_bizcoin = ? where user_id = ?",
					balance,
					vkUserID,
				); err != nil {
					log.Print("error in changing balance")
					additionally.SendError("error in changing balance", err)
					return
				}

				if err = c.WriteJSON(map[string]interface{}{
					"action":   "win_checking",
					"status":   "defeat",
					"balance":  balance,
					"to_mines": toMines,
				}); err != nil {
					log.Println("error in sending defeat status (time_game):", err)
					additionally.SendError("error in sending defeat status (time_game):", err)
					return
				}

			} else {
				if err = c.WriteJSON(map[string]interface{}{
					"action":  "win_checking",
					"status":  "none",
					"balance": balance,
				}); err != nil {
					log.Println("error in sending none status (time_game):", err)
					additionally.SendError("error in sending none status (time_game):", err)
					return
				}
			}

		} else {

			err := c.WriteJSON(config.Error[4])
			if err != nil {
				log.Println("error in sending an error(4) (time_game):", err)
				additionally.SendError("error in sending an error(4) (time_game):", err)
				return
			}

		}
	}
}
