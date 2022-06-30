package main

import (
	"Ethanim_Vote_Server/config"
	"Ethanim_Vote_Server/handle"
	"Ethanim_Vote_Server/service"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mkideal/log"
	"net/http"
	"os"
	"time"
)


func main() {

	curmin := time.Now().Minute()
	newcyclestart, newcycleend := service.GetDayTodayLastMinute(time.Now().Unix(), curmin)
	fmt.Println("cur time is:%d",time.Now().Unix())
	fmt.Println("cur curmin is:%d,GetDayTodayLastMinute() newcyclestart is :%s,newcycleend is:%s", curmin,newcyclestart,newcycleend)

	nodeUrl := "http://101.251.211.201:21630/tri_broadcast_tx_commit?"
	getrandnumstr := handle.RandStr(5)
	fmt.Println("get randnumstr7777 is:%s",getrandnumstr)

	//sendmsgNew := fmt.Sprintf("%s:%s%22","ssA1710BBsssCCaaCC11",getrandnumstr)
	sendmsgNew := "tx=%22" + "ssA1710BBsssCCaadCCs11" + getrandnumstr + "%22"
	getblockInfo,getresq,err := handle.SendRMSVoteMsgToNode(nodeUrl,sendmsgNew)
	if err != nil{
		log.Error("cur after CommitRMSVoteMsg(),get error! ,getRespInfo is :%v,err is:%v",getblockInfo,err)
	}
	fmt.Println("cur CommitRMSVoteMsg(),get getresq hash is:%s,height is:%d",getresq.Hash,getresq.Height)

	//return
	//0630doing
	if err := config.InitWithProviders("multifile/console","./logs"); err != nil {
		panic("init log error: " + err.Error())
	}
	log.Info("log level: %v", log.SetLevelFromString("trace"))
	//mysql
	err = config.InitValidatorConfigInfo()
	if  nil!=err {
		log.Error("from config.json,get json conf err!")
		os.Exit(0)
	}
	gbConf := &config.GbTrustConf

	//log.Info
	fmt.Printf("check Ethanim--get config.json's conf info is %v",gbConf)

	strHost:=fmt.Sprintf(":%d",gbConf.WebPort)
	fmt.Println("strHost is :%s", strHost)

	//http://10.200.1.89:5000/rsm/groups

	groupserverurl := gbConf.GroupServerUrl
	requestInterval :=gbConf.RequestInterval
	nodeUrl = gbConf.NodeUrl
	fmt.Println("get from config.json.GroupServerUrl is :%s,requestInterval is:%d,nodeUrl is:%s", groupserverurl,requestInterval,nodeUrl)

	err =service.InitMysqlDB(gbConf.MySqlCfg)
	if  nil!=err {
		log.Error("cur InitMysqlDB() to conn err!,err is :%v",err)
		//os.Exit(0)
	}
	//err = service.InsertGroupRSMVotes(service.GXormMysql,334,"Rsmid00234",1,"minorityIdslistsss","majorityIdslistss","SlackVoteIds")
	router :=mux.NewRouter().StrictSlash(true)
	//0601add
	handle.StartVoteServer(groupserverurl,requestInterval,nodeUrl)
	router.HandleFunc("/remote/AddValidator", handle.AddValidatorNodeHandler)

	router.HandleFunc("/remote/ClientVoteRSM", handle.AddClientVoteRSMHandle)

	err =http.ListenAndServe(strHost, router)
	if nil!=err {
		fmt.Println("%+v",err)
		os.Exit(0)
	}


}
