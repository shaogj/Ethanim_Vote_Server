package handle

import (
	"Ethanim_Vote_Server/config"
	"Ethanim_Vote_Server/service"
	"Ethanim_Vote_Server/trust"
	"Ethanim_Vote_Server/utils"

	"sync"
	"time"

	//"Ethanim_Vote_Server/service/ethtranssign/ethclientrpc"
	"fmt"
	"io/ioutil"
	"net/http"

	"bytes"
	"encoding/json"

	transproto "Ethanim_Vote_Server/proto"

	"github.com/mkideal/log"
)

type ReturnInfo struct {
	//Cmd         string      `json:"cmd"`      // 命令名,具有协议类型的作用
	InvokeResultCode    int         `json:"invokeResultCode"`    // 返回码(参见枚举 ReturnStatus)
	InvokeResultMessage string      `json:"invokeResultMessage"` // 返回码描述
	Data                interface{} `json:"data"`                // 协议数据
}

// protocol: 返回: 生成数字支付地址
type GenerateAddressRes struct {
	Count         int64    `json:"count"`
	GeneratedAddr []string `json:"getNewAddr"` // 生成地址
	CoinType      string   `json:"coinType"`
}

func JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func JSONResponseWithStatus(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	JSONResponse(w, data)
}

func GeneJsonResultFin(w http.ResponseWriter, r *http.Request, protostruct interface{}, status int, description string) {

	res := ReturnInfo{}
	res.InvokeResultMessage = description
	res.InvokeResultCode = status
	//res.Cmd = cmdname
	res.Data = protostruct
	buf := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(buf)
	err := jsonEncoder.Encode(res)
	if err != nil {
		fmt.Fprintln(w, "command %s:  result to json error: %v", res, err)
		w.Write([]byte(`{"invokeResultCode":999999,"invokeResultMessage":""}`))
	} else {
		//w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		JSONResponseWithStatus(w, res, http.StatusOK)
	}
}

const StatusNewAddressErr = 201 //  生成账号地址错误

var GSettleAccessKey string

//跨域访问
func HttpExCrossDomainAccess(w *http.ResponseWriter) {
	// 允许跨域访问
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "content-type")
}

//请求数据json
//请求数据转为json数据
func HttpExRequestJson(w http.ResponseWriter, r *http.Request, v interface{}) (string, transproto.ErrorInfo) {
	HttpExCrossDomainAccess(&w)
	if r.Method != "POST" {
		return "", transproto.ErrorHttpost
	}
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sttErr := transproto.ErrorRequest
		sttErr.Desc = fmt.Sprintf("%s %s", transproto.ErrorRequest.Desc, err)
		return "", sttErr
	}
	err = json.Unmarshal(result, v)
	if nil != err {
		sttErr := transproto.ErrorRequest
		sttErr.Desc = fmt.Sprintf("%s %s", transproto.ErrorRequest.Desc, err)
		return string(result), sttErr
	}
	return string(result), transproto.ErrorSuccess
}

//0217add
func HttpExRequestMulpJson(w http.ResponseWriter, r *http.Request, v interface{}) (string, transproto.ErrorInfo) {
	HttpExCrossDomainAccess(&w)
	/*if r.Method != "POST" {
		return "", transproto.ErrorHttpost
	}
	*/
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sttErr := transproto.ErrorRequest
		sttErr.Desc = fmt.Sprintf("%s %s", transproto.ErrorRequest.Desc, err)
		return "", sttErr
	}
	err = json.Unmarshal(result, v)
	if nil != err {
		sttErr := transproto.ErrorRequest
		sttErr.Desc = fmt.Sprintf("%s %s", transproto.ErrorRequest.Desc, err)
		return string(result), sttErr
	}
	return string(result), transproto.ErrorSuccess
}

//add for RemoteSignSendTransaction
/*
desc 服务器: 发交易
*/
var ETHTransferFromAddress string
//配置的传输额度
var TransAmountDef int64
//0217doing
//接收BSC传来的验证节点升级请求
type AddValidatorServerReq struct{
	ReqUpdateUrl string
	ReqBSCNodeIP string
	AccessToken	string
}
type BSCTMPubkeyPair struct {
	BSCPubkey	string
	TMPubkey string
}
type AddValidatorNodeResp struct {
	TotalBSCTMPubkeyPair []config.BscAddrMap	 `json:"bscTMPubkeyPairs"`
	PubkeyNum int	`json:"pubkeyNum"`
}
//0217doing end

//2022.0613add
type ClientVoteReq struct{
	ClientId string
	//ClientAddress string	`json:"client_address"`
	RsmGroupId int
	Rsmid string
	VertifyResult bool
	//投票时间
	VoteTime int64
	ClientSignStr string
}


/*
ReqUpdateUrl:"106.3.133.179:4444/remote/AddValidator"
CurBSCNodeIP="10.120.1.104"
AccessToken: "BSCValidatorIdentifyId008"
*/


//0217add
func AddNodeServerHandler(w http.ResponseWriter, r *http.Request) {

	log.Info("fun=AddNodeServerHandler() bef--,request=%v", "jReq")

}

func ReceiveClientVoteHandler(w http.ResponseWriter, r *http.Request) {
	//jReq := transproto.ClientVoteReq{} //SignTransactionReq{}
	jReq := ClientVoteReq{}
	_,sttErr := HttpExRequestJson(w, r, &jReq)//HttpExRequestJson
	if true != transproto.Success(sttErr) {
		log.Error("fun=ReceiveClientVoteHandler()erred! cur request=%v", jReq)
		GeneJsonResultFin(w, r, nil, sttErr.Code, sttErr.Desc)
		return
	}
	log.Info("fun=ReceiveClientVoteHandler() bef--,request=%v", jReq)
	curClientVoteResp :=  transproto.ClientVoteResq{}
	//tmp skip:
	//G_VoteServer.CollectGroupVoteInfo(jReq)
	curClientVoteResp.StatusStr = "StatusStr"
	curClientVoteResp.ResultCode = 200
	//to do compare with
	GeneJsonResultFin(w, r, curClientVoteResp, 200, "调用成功pp")

}
func AddValidatorNodeHandler(w http.ResponseWriter, r *http.Request) {
	jReq := AddValidatorServerReq{}	//SignTransactionReq{}
	//strreq,
	_,sttErr := HttpExRequestJson(w, r, &jReq)//HttpExRequestJson
	if true != transproto.Success(sttErr) {
		GeneJsonResultFin(w, r, nil, sttErr.Code, sttErr.Desc)
		return
	}
	/*_,*/
	log.Info("fun=AddValidatorNodeHandler() bef--,request=%v", jReq)

	totalBSCTMPubkeyPair := config.GbTrustConf.BscAddrMapList

	curAddValidatorNodeResp := AddValidatorNodeResp{}
	curAddValidatorNodeResp.TotalBSCTMPubkeyPair = totalBSCTMPubkeyPair
	curAddValidatorNodeResp.PubkeyNum = len(totalBSCTMPubkeyPair)

	GeneJsonResultFin(w, r, curAddValidatorNodeResp, 200, "调用成功pp")

}

func AddClientVoteRSMHandle(w http.ResponseWriter, r *http.Request) {
	//jReq := AddValidatorServerReq{}	//SignTransactionReq{}
	jReq := transproto.ClientVoteReq{}//ClientVoteReq{}
	_,sttErr := HttpExRequestJson(w, r, &jReq)//HttpExRequestJson
	if true != transproto.Success(sttErr) {
		GeneJsonResultFin(w, r, nil, sttErr.Code, sttErr.Desc)
		log.Error("fun=ReceiveClientVoteHandler() bef--,request=%v", jReq)
		return
	}
	/*_,*/
	log.Info("fun=ReceiveClientVoteHandler() good,cur request=%v", jReq)
	log.Info("cur to add ServerGroup map len is:%d",len(G_VoteServer.ServerGroup))
	//0613,,to upgrade to channel msg
	G_VoteServer.CollectGroupVoteInfo(jReq)

	getresq :=transproto.ClientVoteResq{}
	getresq.ResultCode = 0
	getresq.ResultMsg = "invoke check good!"
	//0614add
	err := service.ClientVoteRecordSave(service.GXormMysql,jReq.ClientId,jReq.Rsmid,jReq.RsmGroupId,jReq.ClientSignstr,true)
	if err != nil {
		log.Error("ClientVoteRecordSave(),Insert row is failed! err is-:%v \n",err)
	//return err

	}
	/*
	totalBSCTMPubkeyPair := config.GbTrustConf.BscAddrMapList

	curAddValidatorNodeResp := AddValidatorNodeResp{}
	curAddValidatorNodeResp.TotalBSCTMPubkeyPair = totalBSCTMPubkeyPair
	curAddValidatorNodeResp.PubkeyNum = len(totalBSCTMPubkeyPair)
	*/
	GeneJsonResultFin(w, r, getresq, 200, "调用成功pp")

}


type VoteServerConfig struct {
	RequestInterval int
	PrivKey  	string
	PublicKey	string
	RSMServerUrl string
	TMNodeUrl string
}
//RSM周期服务分组的信息
type VoteServer struct {
	VoteConfigParams *VoteServerConfig
	//ServerGroup map[int]map[transproto.Rsmnode][]string
	//每次从server端获取rms分组的client列表
	ServerGroup map[transproto.RsmNode][]ClientRSMInfo
	lock sync.Mutex
	//收集每个groupID中，client对rms投票的msg
	RMSGroupVotes	map[int]transproto.RSMVoteGroupMsgs

}
//0605,RSM分组的列表消息
type ClientRSMInfo struct {
	Startime int64
	Endtime int64
	ClientID string

}
//type RsmServerGroupResq struct{
//	startime int64	`json:"start_time"`
//	Endtime int64	`json:"end_time"`

func DefaultVoteConfig() *VoteServerConfig {
	ip, err := ExternalIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip.String())
	trustLocalIp := ip.String()
	if trustLocalIp == ""{
		trustLocalIp = "127.0.0.1"
	}
	trustUrlVerifyReal := fmt.Sprintf("http://%s:46657/tri_block_validators",trustLocalIp)
	fmt.Println("in DefaultTrustConfig(),cur trustUrlVerifyReal is:%s",trustUrlVerifyReal)
	//ReqTMGRPCUrlReal := fmt.Sprintf("%s:46658",trustLocalIp)

	return &VoteServerConfig{
		//RequestInterval:4,
		//0612add,30sec 请求一次分组服务
		RequestInterval:10,
		PrivKey:	"PrivKeystr",
		PublicKey:	"PubKeystr",
		RSMServerUrl:trustUrlVerifyReal,
	}

}
func (this *VoteServer) RequestRSMGroupInfo(rsmserurl string,accessKey string)(curRsmServerGroupInfo *transproto.RsmServerGroupInfo,err error){
	log.Info("Succ Invoke RequestRSMGroupInfo()' get getValidators!,RequestRSMGroupInfo is :%s","ssss444")
	var respinfo = &transproto.RsmServerGroupInfo{Rsmcount:14,ServerSignStr:"ServerSignStrnfo43"}
	return respinfo,nil
}

func NewVoteServer(curinfo *VoteServerConfig) *VoteServer {
	curVoteServer := &VoteServer{}
	curVoteServer.VoteConfigParams = curinfo
	curVoteServer.ServerGroup = make(map[transproto.RsmNode][]ClientRSMInfo,0)
	curVoteServer.RMSGroupVotes = make(map[int]transproto.RSMVoteGroupMsgs)
	return curVoteServer

}

//to insert value
//0613，get mysql data:

//获取ETH的地址的私钥
func GetAddrPrivkeyETH(curaddress string) (addrPrikey string,err error){

	//engineread:= ETH_MOrmEngine
	engineread:= service.GXormMysql
	//get address 's [privkey]
	selectsql := "select * from bsc_account_key where address = '"  + curaddress + "'"
	addr_accountinfo, err := engineread.Query(selectsql)
	if err != nil || len(addr_accountinfo) <= 0{
		log.Error("when GetAddrPrivkeyETH(),curaddress' is:%s ,get privkey error: %v", curaddress,err)
		return "",err
	}

	curaddrprivkey := string(addr_accountinfo[0]["priv_key"])
	log.Info("when GetAddrPrivkeyETH(),curaddress is :%s,get addr's privkey succ ,info is: %v", curaddress,curaddrprivkey)
	return curaddrprivkey,nil
}


//0613
func (this *VoteServer) CollectGroupVoteInfo(curclientvoteinfo transproto.ClientVoteReq) {
	//this.RMSGroupVotes = this.RMSGroupVotes
	//根据ServerGroup中的groupid时间，和rmsid；
	//比较msg中的VoteTime，垃圾数据则扔掉
	//else插入相关的记录中
	rsmid := curclientvoteinfo.Rsmid
	rsmGroupId := curclientvoteinfo.RsmGroupId
	key := transproto.RsmNode{RsmId: rsmid, GroupId: rsmGroupId}
	//0613pm,tmp add:
	var curClientRSMInfo = []ClientRSMInfo{ClientRSMInfo{ClientID:"3323fsdfsaf"}}
	this.ServerGroup[key] = curClientRSMInfo
	_,ok := this.ServerGroup[key]
	if !ok {
		log.Info("In CollectGroupVoteInfo(), get cur key from clientmsg is exist no in ServerGroup.key is:%v,recv clientmsg req info is:%v",key,curclientvoteinfo)
	}else{
		//var curRSMVoteGroupMsgs  transproto.RSMVoteGroupMsgs
		getVoteGroupMsgs,ok :=this.RMSGroupVotes[rsmGroupId]
		if !ok{
			curRSMVoteGroupMsgs := transproto.RSMVoteGroupMsgs{}
			curRSMVoteGroupMsgs.GatherVertifyResult = curclientvoteinfo.VertifyResult
			if curclientvoteinfo.VertifyResult {
				curRSMVoteGroupMsgs.GatherTrustCount ++
			}
			curRSMVoteGroupMsgs.ClientVote = append(curRSMVoteGroupMsgs.ClientVote,curclientvoteinfo)
			this.RMSGroupVotes[rsmGroupId] = curRSMVoteGroupMsgs
			log.Info("In CollectGroupVoteInfo(), get cur key from clientmsg is new added to RMSGroupVotes. key is:%v,recv clientmsg req info is:%v",key,curclientvoteinfo)
		}else{
			getVoteGroupMsgs.ClientVote = append(getVoteGroupMsgs.ClientVote,curclientvoteinfo)
			this.RMSGroupVotes[rsmGroupId] = getVoteGroupMsgs
			log.Info("In CollectGroupVoteInfo(), get cur key from clientmsg is exist in RMSGroupVotes. key is:%v,recv clientmsg req info is:%v",key,curclientvoteinfo)

		}

	}
}

//transproto.RsmServerGroupResq
func (this *VoteServer) AddRSMGroupList(RsmGroupId int,Startime int64,curgroupinfo []transproto.GroupItems) {
	var curClientRSMInfo ClientRSMInfo
	for _, groupitem := range curgroupinfo {
		//item.
		rmsid := groupitem.RmsId
		client := groupitem.Clients
		key := transproto.RsmNode{RsmId: rmsid, GroupId: RsmGroupId}
		_,ok := this.ServerGroup[key]
		if !ok {
			for _, curclient := range client {
				curClientRSMInfo = ClientRSMInfo{Startime:Startime,ClientID: curclient}
				this.ServerGroup[key] = append(this.ServerGroup[key], curClientRSMInfo)
			}
			log.Info("in AddRSMGroupList(), get cur ServerGroup' key is:%v,get cur groupitem is:%v,err is:%v",key,groupitem)

		}

		//this.ServerGroup[key] = append(this.ServerGroup[key], curClientRSMInfo)

	}
	///rmsid := getcurgroupinfo.RSMGrouList[0][0]    //.(*string)
	//key := transproto.RsmNode{RsmId: rmsid.(string), GroupId: getcurgroupinfo.RsmGroupId}
	//ClientRSMInfoList,

}
func (this *VoteServer) ReqGetRsmGrouplist() (groupid int,err error){
	log.Info("cur Invoke ReqGetRsmGrouplist(), to req and AddRSMGroupList is :%s","ssss444")
	ht:=utils.CHttpClientEx{}
	ht.Init()

	cli := trust.NewRSMServerCli(this.VoteConfigParams.RSMServerUrl)
	//127.0.0.1:8999
	servergroupresq, err := cli.GetRsmGrouplist(ht)
	if servergroupresq == nil {
		log.Error("get cli.GetRsmGrouplist() err!,get err is:%v",err)
		return 0,err
	}
	fmt.Println("after GetRsmGrouplist(),get grouprsmMap is :%v,err is:%v", servergroupresq,err)
	this.AddRSMGroupList(servergroupresq.RsmGroupId,servergroupresq.Startime,servergroupresq.GroupItems)
	return servergroupresq.RsmGroupId,nil
}
func (this *VoteServer) StartRequest() {
	iserion :=0
	for {
		time.Sleep(time.Second * time.Duration(this.VoteConfigParams.RequestInterval))
		fmt.Println("cur In StartRequest()，invoke interval is:%d",this.VoteConfigParams.RequestInterval)
		//step to 请求group server：请求应在每个节点出块流程commit后

		//step to 请求可信节点列表
		//getGroupInfo,err :=this.RequestRSMGroupInfo(this.VoteConfigParams.RSMServerUrl,this.VoteConfigParams.PublicKey)
		getGroupCount,err :=this.ReqGetRsmGrouplist()
		if err != nil{
			fmt.Println("cur RequestTrustInfo() num no matched!,get getGroupCount is:%d,get err is:%v",getGroupCount,err)
			//continue
		}

		log.Info(fmt.Sprintf("succ Invoke cur RequestTrustInfo(),get total group num is :%d,cur Groupid is: %d,err is:%v",len(this.ServerGroup),getGroupCount,err))
		/*trustnodeaddr := make(map[string]TrustScore,trustNum)
		//1127fix
		if getTrustInfo == nil {
			continue
		}
		getValidators := getTrustInfo.(*proto.ResultValidators)
		*/
		//for id,nodeaddr:= range getValidators.Validators{
			//pkCurValidator := fmt.Sprintf("%s", nodeaddr.PubKey)
			//if !(strings.HasPrefix(pkCurValidator, "PubKeyEd25519") && len(pkCurValidator) == 79) {


			/*trustdatapubkey :=fmt.Sprintf("%s/%d", pkCurValidator[14:78], 300)
			parts := bytes.Split([]byte(trustdatapubkey), []byte("/"))
			if len(parts) == 2 {
				trustpubkey = string(parts[0])
				score,_ = strconv.Atoi(string(parts[1]))
			}
			*/
			//var curTrustScore TrustScore

			//trustnodeaddr[envcodeStr] = curTrustScore

			//log.Info("Succ Invoke RequestTrustInfo()' get getValidators!,trustpubkey is :%s",trustpubkey)

		//}

		iserion ++
		this.lock.Lock()
		//11.25,返回的可信node列表，比现有的列表少，则设置不可信的节点score为-1，及为不可信的节点
		//this.TrustNodeMap = trustnodeaddr
		this.lock.Unlock()

	}
}
var G_VoteServer  *VoteServer

func  StartVoteServer(groupserverurl string,requestInterval int) *VoteServer{
	crustConfig :=DefaultVoteConfig()
	crustConfig.RSMServerUrl = groupserverurl
	crustConfig.RequestInterval = requestInterval
	curVoteRsmTask:= NewVoteServer(crustConfig)
	go curVoteRsmTask.StartRequest()
	fmt.Println("StartVoteServer is start!,,gbConf' gbTrustConf.crustConfig is %v", crustConfig)
	G_VoteServer = curVoteRsmTask
	return G_VoteServer
}