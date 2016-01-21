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

type AddGamerParam struct {
	Name      string
	Gamertype string
	Total     float64
}

type AddShareParam struct {
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

	gamer := AddGamerParam{}

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

	addGamer(gamer)

	// jsonSG, err := json.Marshal(ShareGame)
	// if err == nil {

	// }
	result.Msg = "Add Gamer Success"

	printInfo(ShareGame)
	httpResp(w, result)
}

func printInfo(sg model.ShareGame) {
	fmt.Println("______________________________________________________________________________")
	for _, si := range ShareInfos {
		s := fmt.Sprintf("%.3f", si.Price)
		fmt.Printf("股票信息:\n")
		fmt.Printf("股票代码: %v\t股票名称: %v\t股票现价: %v\n", si.ShareID, si.ShareName, s)
	}
	fmt.Println("______________________________________________________________________________")
	for _, i := range sg.Gamers {
		fmt.Printf("股票所有人: %v\n", i.Name)
		fmt.Printf("商户类型: %v\n", i.Gamertype)
		fmt.Printf("总资产: %v\n", i.Total)
		fmt.Printf("现金资产: %v\n", i.CashTotal)
		fmt.Printf("股票资产: %v\n", i.SharesTotals)
		fmt.Printf("明细: \n")
		for _, j := range i.Shares {
			fmt.Printf("股票代码: %v\t", j.ShareID)
			fmt.Printf("股票市值: %v\t", j.SharesTotal)
			fmt.Printf("持股数量: %v\t", j.Number)
			s1 := fmt.Sprintf("%.3f", j.AvgPrice)
			fmt.Printf("平均价格: %v\n", s1)
		}
		fmt.Println("______________________________________________________________________________")
	}

}

func addGamer(gamer AddGamerParam) {
	game := model.Gamer{}
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
}

func (server *webServer) addshare(w http.ResponseWriter, r *http.Request) {
	var errmsg string

	result := &result{}
	result.Code = 0
	asp := AddShareParam{}
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

		err = json.Unmarshal(reqBody, &asp)
		if err != nil {
			errmsg = fmt.Sprintf("json unmarshal err: (%s) ", err.Error())
			result.Code = -1
			result.Msg = errmsg
			httpResp(w, result)
			return
		}
	}

	si := model.ShareInfo{}
	si.ShareID = asp.ShareID
	si.ShareName = asp.ShareName
	si.Price = asp.Price

	name := asp.Name

	gamer, ok := ShareGame.Gamers[name]
	if !ok {
		errmsg = fmt.Sprintf("no gamer selected ")
		result.Code = -1
		result.Msg = errmsg
		httpResp(w, result)
		return
	}

	share, ok1 := ShareGame.Gamers[name].Shares[asp.ShareID]
	if !ok1 {
		newshare := model.Share{}
		ShareInfos[asp.ShareID] = si
		newshare.ShareID = asp.ShareID
		newshare.Number = asp.Number
		newshare.SharesTotal = asp.Price * newshare.Number
		newshare.AvgPrice = newshare.SharesTotal / newshare.Number
		gamer.Shares[newshare.ShareID] = newshare
		gamer.SharesTotals += gamer.Shares[asp.ShareID].SharesTotal
		gamer.Total = gamer.CashTotal + gamer.SharesTotals

		ShareGame.Gamers[name] = gamer

	} else {
		newshare := model.Share{}
		newshare.ShareID = asp.ShareID
		newshare.Number = asp.Number
		newshare.SharesTotal = asp.Price * newshare.Number

		share.ShareID = share.ShareID
		share.Number += newshare.Number
		share.SharesTotal += newshare.SharesTotal
		share.AvgPrice = share.SharesTotal / share.Number

		gamer.Shares[share.ShareID] = share
		gamer.SharesTotals += newshare.SharesTotal
		gamer.Total = gamer.CashTotal + gamer.SharesTotals
		ShareGame.Gamers[name] = gamer
	}

	ShareInfos[asp.ShareID] = si

	// jsonSG, err := json.Marshal(ShareGame)
	// if err == nil {

	// }
	// result.Msg = string(jsonSG)

	result.Msg = "Add Shares Success"

	printInfo(ShareGame)

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

	ShareInfos = make(map[string]model.ShareInfo, 0)
}
