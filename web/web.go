package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jd.com/cdn/cashgame/cfg"
	"jd.com/cdn/cashgame/model"
	"net/http"
	"strconv"
)

type webServer struct {
	Config *cfg.Config
}

type result struct {
	Code int
	Msg  string
}

// type taskResult struct {
// 	Code int
// 	Msg  string
// }

type reqaddgamer struct {
	Name      string
	Gamertype string
	Total     float64
}

type reqaddshare struct {
	Name      string
	ShareID   string
	ShareName string
	Price     float64
	Number    float64
}

func WebStart(config *cfg.Config) error {
	webServer := &webServer{Config: config}
	fmt.Printf("WebServer init success\n")
	err := webServer.Start()
	if err != nil {
		return err
	}
	fmt.Printf("WebServer is  running\n")
	return nil
}

func (server *webServer) Start() error {
	http.HandleFunc("/alive", server.alive)
	http.HandleFunc("/getall", server.getAll)
	http.HandleFunc("/addgamer", server.addgamer)
	http.HandleFunc("/addshare", server.addshare)

	fmt.Printf("server ready to run:%s\n", strconv.Itoa(server.Config.Server.Port))
	http.ListenAndServe(":"+strconv.Itoa(server.Config.Server.Port), nil)
	return nil
}

func (server *webServer) alive(w http.ResponseWriter, r *http.Request) {
	result := &result{}
	result.Code = 0
	result.Msg = "server is alive"
	httpResp(w, result)
}

func (server *webServer) getAll(w http.ResponseWriter, r *http.Request) {
	result := &result{}
	result.Code = 0
	jsonSG, err := json.Marshal(ShareGame)
	if err == nil {

	}
	result.Msg = string(jsonSG)

	httpResp(w, result)
}

func (server *webServer) addgamer(w http.ResponseWriter, r *http.Request) {
	var errmsg string
	result := &result{}
	result.Code = 0
	game := model.Gamer{}
	game.SharesTotals = 0
	gamer := reqaddgamer{}
	content_type := r.Header.Get("content-type")
	if content_type == "application/json" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errmsg = fmt.Sprintf("io util read all err (%s) \n", err.Error())
			result.Code = -1
			result.Msg = errmsg
			httpResp(w, result)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(reqBody, &gamer)
		if err != nil {
			errmsg = fmt.Sprintf("json unmarshal err: (%s) ", err.Error())
			result.Code = -1
			result.Msg = errmsg
			httpResp(w, result)
			return
		}
	}

	game.Name = gamer.Name           //arg "中国石油"
	game.Gamertype = gamer.Gamertype //arg "company"
	game.Total = gamer.Total         //arg 1000000
	game.CashTotal = gamer.Total
	game.SharesTotals = 0
	game.Shares = make(map[string]model.Share, 0)

	_, ok := ShareGame.Gamers[game.Name]
	if !ok {
		ShareGame.Gamers[game.Name] = game
	}

	jsonSG, err := json.Marshal(ShareGame)
	if err == nil {

	}
	result.Msg = string(jsonSG)

	httpResp(w, result)
}

func (server *webServer) addshare(w http.ResponseWriter, r *http.Request) {
	var errmsg string
	result := &result{}
	result.Code = 0
	s := reqaddshare{}
	content_type := r.Header.Get("content-type")
	if content_type == "application/json" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errmsg = fmt.Sprintf("io util read all err (%s) \n", err.Error())
			result.Code = -1
			result.Msg = errmsg
			httpResp(w, result)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(reqBody, &s)
		if err != nil {
			errmsg = fmt.Sprintf("json unmarshal err: (%s) ", err.Error())
			result.Code = -1
			result.Msg = errmsg
			httpResp(w, result)
			return
		}
	}

	name := s.Name
	sharename := s.ShareName
	gamer, ok := ShareGame.Gamers[name]
	if !ok {
		errmsg = fmt.Sprintf("no gamer selected ")
		result.Code = -1
		result.Msg = errmsg
		httpResp(w, result)
		return
	}
	share, ok1 := ShareGame.Gamers[name].Shares[sharename]
	if !ok1 {
		newshare := model.Share{}
		newshare.ShareID = s.ShareID
		newshare.ShareName = s.ShareName
		newshare.Price = s.Price
		newshare.Number = s.Number
		newshare.SharesTotal = newshare.Price * newshare.Number
		gamer.Shares[sharename] = newshare
		gamer.SharesTotals += gamer.Shares[sharename].SharesTotal
		gamer.Total += gamer.CashTotal + gamer.SharesTotals

		ShareGame.Gamers[name] = gamer
		// ShareGame.Gamers[name].SharesTotals = 0 //+= ShareGame.Gamers[name].Shares[sharename].SharesTotal
		// ShareGame.Gamers[name].Total = 0        //+= ShareGame.Gamers[name].CashTotal + ShareGame.Gamers[name].SharesTotals
		fmt.Println(ShareGame.Gamers[name])
	} else {
		newshare := model.Share{}
		newshare.ShareID = s.ShareID
		newshare.ShareName = s.ShareName
		newshare.Price = s.Price
		newshare.Number = s.Number
		newshare.SharesTotal = newshare.Price * newshare.Number

		share.Number += s.Number
		share.SharesTotal += newshare.SharesTotal
		share.Price = share.SharesTotal / share.Number

		gamer.Shares[sharename] = share
		gamer.SharesTotals += newshare.SharesTotal
		gamer.Total += gamer.CashTotal + gamer.SharesTotals

		ShareGame.Gamers[name] = gamer

		fmt.Println(ShareGame.Gamers[name])
	}

	// i.Shares[name].Price = s.Price
	// i.Shares[name].Number = s.Number
	// i.Shares[name].SharesTotal += i.Shares[name].Price * i.Shares[name].Number
	// i.CashTotal = i.CashTotal + i.Shares[name].SharesTotal
	// fmt.Println(name, i.Shares)

	jsonSG, err := json.Marshal(ShareGame)
	if err == nil {

	}
	result.Msg = string(jsonSG)

	httpResp(w, result)
}

func httpResp(w http.ResponseWriter, result interface{}) {
	jsonRt, err := json.Marshal(result)
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonRt)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func InitEnv(conf *cfg.Config) {
	ShareGame = model.ShareGame{}
	ShareGame.Country.Name = conf.Country.Name
	ShareGame.Country.Retio = conf.Country.Retio
	ShareGame.Country.Income = float64(0)
	ShareGame.Organization.Name = conf.Organization.Name
	ShareGame.Organization.Retio = conf.Organization.Retio
	ShareGame.Organization.Income = float64(0)
	ShareGame.Gamers = make(map[string]model.Gamer, 0)
}
