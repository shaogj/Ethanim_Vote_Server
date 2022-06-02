package main

import (
	"Ethanim_Vote_Server/config"
	"Ethanim_Vote_Server/handle"
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
	fmt.Println("sss444!!")
	curtime := time.Now()
	///curtime.Before()
	//todo,时间戳比较；
	fmt.Println("curtime is:%v",curtime)
	strHostPort3 :=  4422//3333
	//0217ad
	fmt.Println("strHost0011 is :%s", strHostPort3)
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

	//0507addfor ipmap
	err = config.InitLibp2pAddrMapInfo()
	if  nil!=err {
		fmt.Println("from libp2p_netaddr_map.json.json,get json conf err!")
		log.Error("from libp2p_netaddr_map.json.json,get json conf err!")
		os.Exit(0)
	}
	gbConfIpMap := &config.GbPeerIPMap
	var ConfigPeerMapIP = make(map[string]string)
	xsrcmap := map[string]string{"192.168.77/10011": "111.153", "192.168.44/10022":"222.444"}
	for _,Libp2pAddrItem := range gbConfIpMap.Libp2pAddrMapList{
		nodePeerFoundAddr :=fmt.Sprintf("/ip4/%s/tcp/%d",Libp2pAddrItem.OuterNetAddr,Libp2pAddrItem.Libp2pListenPort)
		//fmt.Printf("check cur i:%d,Libp2pAddrItem info is %s,get nodePeerFoundAddr is:%s\n",i,Libp2pAddrItem,nodePeerFoundAddr)
		ConfigPeerMapIP[nodePeerFoundAddr] = Libp2pAddrItem.InneNetAddr
	}
	//for i,xsrcmapitem := range xsrcmap {
	for peerOuterAddr,addrItem := range ConfigPeerMapIP {
		fmt.Printf("checklibp2p--get addrItem conf peerouteraddr:is %s,inneraddr is:%s\n",peerOuterAddr,addrItem)
		xsrcmap[peerOuterAddr] = addrItem
	}
	fmt.Printf("checklibp2p--get xsrcmap map len is:%d,xsrcmap info:is %v\n",len(xsrcmap),xsrcmap)

	//交易发送节点的配置
	//ethclientrpc.ReqNodeUrl = gbConf.BSCNodeUrl
	fmt.Println("cur cfgparams: ReqNodeValidatorInfo is :%v", gbConf.BscAddrMapList)

	router :=mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/remote/GetTestCoinTx", handle.RemoteSignSendTransaction)
	//0601add
	handle.StartVoteServer()
	router.HandleFunc("/remote/AddValidator", handle.AddValidatorNodeHandler)

	err =http.ListenAndServe(strHost, router)
	if nil!=err {
		fmt.Println("%+v",err)
		os.Exit(0)
	}


}
