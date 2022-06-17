package main

import (
	"Ethanim_Vote_Server/config"
	"Ethanim_Vote_Server/handle"
	"Ethanim_Vote_Server/service"
	"Ethanim_Vote_Server/trust"
	"Ethanim_Vote_Server/utils"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mkideal/log"
	"net/http"
	"os"
	"time"
)

//0217add
func SetupNodeServer(webPort int) {
	//m := http.NewServeMux()
	router :=mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/remote/AddValidator3", handle.AddNodeServerHandler)
	//m.Handle("/debug/metrics/prometheus", prometheus.Handler(metrics.DefaultRegistry))
	log.Info("AddValidator3", "webPort", fmt.Sprintf("http://%s/debug/metrics", webPort))
	go func() {
		strHost:=fmt.Sprintf(":%d",webPort)
		fmt.Println("AddValidator3strHost is :%s", strHost)
		if err := http.ListenAndServe(strHost, router); err != nil {
			log.Error("Failure in running metrics server", "err", err)
		}
		fmt.Println("exec after ListenAndServe(), cur address is :%s", strHost)

	}()
}

func main() {

	strHostPort3 :=  4422//3333
	//0607testing
	//return
	//end 0607
	//0217ad
	fmt.Println("strHost0011 is :%s", strHostPort3)
	//0616testing:
	curmin := time.Now().Minute()
	newcyclestart, newcycleend := service.GetDayTodayLastMinute(time.Now().Unix(), curmin)
	fmt.Println("cur curmin is:%d,GetDayTodayLastMinute() newcyclestart is :%s,newcycleend is:%s", curmin,newcyclestart,newcycleend)
	var nowtime int64 = time.Now().Unix()
	var strtime string = time.Unix(nowtime, 0).Format("2006-01-02 15:04:05")
	fmt.Println("cur time is :%d, strtime is :%s",nowtime, strtime)


	//return
	SetupNodeServer(strHostPort3)

	if err := config.InitWithProviders("multifile/console","./logs"); err != nil {
		panic("init log error: " + err.Error())
	}
	log.Info("log level: %v", log.SetLevelFromString("trace"))
	//mysql
	err := config.InitValidatorConfigInfo()
	if  nil!=err {
		log.Error("from config.json,get json conf err!")
		os.Exit(0)
	}
	gbConf := &config.GbTrustConf

	//log.Info
	fmt.Printf("check Ethanim--get config.json's conf info is %v",gbConf)

	strHost:=fmt.Sprintf(":%d",gbConf.WebPort)
	fmt.Println("strHost is :%s", strHost)
	//0607add
	ht:=utils.CHttpClientEx{}
	ht.Init()
	//192.168.1.114
	//http://10.200.1.89:5000/rsm/groups
	//cururl := "http://10.200.1.24:5000"
	groupserverurl := gbConf.GroupServerUrl
	requestInterval :=gbConf.RequestInterval
	fmt.Println("get from config.json.GroupServerUrl is :%s,requestInterval is:%d", groupserverurl,requestInterval)
	cli := trust.NewRSMServerCli(groupserverurl)
	//127.0.0.1:8999
	GroupAttachRSMMap,startime, err := cli.GetRsmGrouplist(ht)
	fmt.Println("after GetRsmGrouplist(),startime is:%d,get grouprsmMap is :%v,err is:%v", startime,GroupAttachRSMMap,err)

	//交易发送节点的配置
	//ethclientrpc.ReqNodeUrl = gbConf.BSCNodeUrl
	err =service.InitMysqlDB(gbConf.MySqlCfg)
	if  nil!=err {
		log.Error("cur InitMysqlDB() to conn err!,err is :%v",err)
		//0117,,test to recover
		//0614PM:
		//os.Exit(0)
	}
	fmt.Println("cur cfgparams: ReqNodeValidatorInfo is :%v", gbConf.BscAddrMapList)
	//0616,trying
	/*
	ethaccount := "0xb42cb187D7738fA9c14dB86e0A25014D6c296bCd"

	addrPrikey,err :=handle.GetAddrPrivkeyETH(ethaccount)
	if  nil!=err {
		log.Error("cur InitMysqlDB() to conn err!,err is :%v", err)
	}

	fmt.Println("cur ethaccount is:%s,get addrPrikey is :%s", ethaccount,addrPrikey)
	*/
	router :=mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/remote/GetTestCoinTx", handle.RemoteSignSendTransaction)
	//0601add
	handle.StartVoteServer(groupserverurl,requestInterval)
	router.HandleFunc("/remote/AddValidator", handle.AddValidatorNodeHandler)

	router.HandleFunc("/remote/ClientVoteRSM", handle.AddClientVoteRSMHandle)

	err =http.ListenAndServe(strHost, router)
	if nil!=err {
		fmt.Println("%+v",err)
		os.Exit(0)
	}


}
