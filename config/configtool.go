package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mkideal/log"
	"io/ioutil"
	"os"
	"path/filepath"
)
var   (
	HConf ConfigTools

)
var GbConf ConfigInfomation = ConfigInfomation{}
var ChainId string
const (
	CoinEthereum string ="ETH"
	CoinUSDT string = "USDT"
	CoinBSC string ="BSC"
)


//10.30,目前预留
type WDCNodeConf struct{
	//RPC for WDC:
	RPCPort		int
	RPCHostPort		string
	RPCUser		string
	RPCPassWord		string

}

//1119adding:
type USDTConf struct{
	//RPC for BTC:
	RPCPort		int
	RPCHostPort		string
	RPCUser		string
	RPCPassWord		string
	RPCTestNet		int

}



type  MongoDatabaseInfo struct {
	Db string   	`josn:"dbase"`  //数据库
	Coll string   	`josn:"coll"` //数据集
	User string		`josn:"user"` //用户名
	Pass  string	`josn:"pass"` //密码
}
type MySqlConfig struct  {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Dbname string `json:"dbname"`
}

type SettleApiReq struct {
	SettlApiQuery string `json:"settleapiquery"`
	SettlApiUpdate string `json:"settleapiupdate"`
}

type SettleAccessKey struct {
	AccessComePubKey string `json:"AccessPubKey"`
	AccessPrivKey string `json:"AccessPrivKey"`
}
//sgj 1019 add
type ConfigInfomation struct {
	MySqlCfg        MySqlConfig         `json:"MySqlConfig"`

	SettleApiReq    SettleApiReq         `json:"SettleApiReq"`
	SettleApiQuery string `json:"SettleApiQuery"`
	SettleApiUpdate string `json:"SettleApiUpdate"`
	//1217 add for WDC,KTC
	SettleApiDepositQuery string `json:"SettleApiDepositQuery"`

	//1019 add
	WebPort			int `json:"WebPort"`
	//2022.0407add
	ChainId	string `json:"ChainId"`
	SettleAccessKey    SettleAccessKey         `json:"SettleAccessKey"`
	//请求RSMGroupSever分组URL
	BscAddrMapList    []BscAddrMap         `json:"bsc_pubkey_addr_maps"`
	TrustBSCnodeTMUrl    string         `json:"TrustBSCnodeTMUrl"`
	RequestTrustInterval		int 	`json:"request_trust_interval"`
	//2022。0114add:
	ReqRSMGroupSeverUrl    string         `json:"GroupServerUrl"`


}

//sgj202112add
type BscAddrMap struct {
	BscNodePubkeyAddr   string `json:"bsc_node_pubkey_addr"`
	TmNodePubkeyAddr string `json:"tm_node_pubkey_addr"`
}
func InitWithProviders(providers, dir string) error {
	return log.Init(providers, log.M{
		"rootdir":     dir,
		"suffix":      ".txt",
		"date_format": "%04d-%02d-%02d",
	})
}

func defaultInt(ptr *int, dft int) {
	if *ptr == 0 {
		*ptr = dft
	}
}

func defaultString(ptr *string, dft string) {
	if len(*ptr) == 0 {
		*ptr = dft
	}
}
type ConfigTools struct {

	// log
	LogProviders string
	LogLevel     string
	Logpath      string
	MgoAddrData		MongoDatabaseInfo   /*用户addr数据*/
	//1119adding:
	CurUSDTConf USDTConf
	//OrmEngine          *xorm.Engine

}

//sgj 1019 add
func InitConfigInfo() error {
	//*good conf:
	//log.SetFlags(log.Lshortfile | log.Ltime)
	var strConf string
	flag.StringVar(&strConf, "conf", "config.json", "config <file>")
	flag.Parse()
	byData, err := ioutil.ReadFile(strConf)
	if nil != err {
		log.Error("Read config file :::%v", err)
		return err
	}
	err = json.Unmarshal(byData, &GbConf)
	if nil != err {
		log.Error("Unmarshal config file :::%v", err)
		return err
	}
	log.Info("ConfigInfo:::%+v", GbConf)
	return nil
}

//0227add
var GbTrustConf ConfigTrustInfomation

//0507add
func GetAppPath() string {
	path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	dir := filepath.Dir(path)
	//fmt.Println("cur gethapp run path,dir is:",dir)
	return dir
}

var GbPeerIPMap Libp2pAddrInfo
func InitLibp2pAddrMapInfo() error {

	curPath := GetAppPath()
	addmapfilePath := fmt.Sprintf("%s/%s", curPath, "libp2p_netaddr_map.json")	//pwd
	fmt.Printf("loadconfpeerip(),cur addmapfilePath file is:%s ", addmapfilePath)

	byData, err := ioutil.ReadFile(addmapfilePath)
	if nil != err {
		fmt.Printf("cur addmapfilePath exist no libp2p_netaddr_map file is:%s ", addmapfilePath)
		//log.Error("Read config file :::%v", err)
		return nil	//err
	}
	err = json.Unmarshal(byData, &GbPeerIPMap)
	if nil != err {
		log.Error("Unmarshal config file :::%v", err)
		return err
	}
	log.Info("ConfigLibp2pAddrInfo:::%+v", GbPeerIPMap)
	return nil
}

func InitValidatorConfigInfo() error {
	//*good conf:
	//log.SetFlags(log.Lshortfile | log.Ltime)
	var strConf string
	flag.StringVar(&strConf, "conf", "config.json", "config <file>")
	flag.Parse()
	byData, err := ioutil.ReadFile(strConf)
	if nil != err {
		log.Error("Read config file :::%v", err)
		return err
	}
	err = json.Unmarshal(byData, &GbTrustConf)
	if nil != err {
		log.Error("Unmarshal config file :::%v", err)
		return err
	}
	log.Info("ConfigValidatorInfo:::%+v", GbTrustConf)
	return nil
}
type  ConfigTrustInfomation struct {
	BscAddrMapList       []BscAddrMap `json:"bsc_pubkey_addr_maps"`
	RequestTrustInterval int          `json:"request_token_key"`
	WebPort	int				`json:"WebPort"`
	GroupServerUrl	string `json:"group_server_url"`
}

type Libp2pAddrMap struct {
	InneNetAddr   string `json:"inne_net_addr"`
	OuterNetAddr   string `json:"outer_net_addr"`
	Libp2pListenPort	int	`json:"libp2p_listen_port"`
}
type Libp2pAddrInfo struct {
	Libp2pAddrMapList       []Libp2pAddrMap `json:"libp2p_addr_maps"`

}
