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

	//http://10.200.1.89:5000/rsm/groups

	groupserverurl := gbConf.GroupServerUrl
	requestInterval :=gbConf.RequestInterval
	fmt.Println("get from config.json.GroupServerUrl is :%s,requestInterval is:%d", groupserverurl,requestInterval)

	err =service.InitMysqlDB(gbConf.MySqlCfg)
	if  nil!=err {
		log.Error("cur InitMysqlDB() to conn err!,err is :%v",err)
		//os.Exit(0)
	}
	//err = service.InsertGroupRSMVotes(service.GXormMysql,groupId,"Rsmid00234",1,"minorityIdslistsss","majorityIdslistss")
	router :=mux.NewRouter().StrictSlash(true)
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
