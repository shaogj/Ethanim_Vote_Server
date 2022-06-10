package handle

import (
	"Ethanim_Vote_Server/config"

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
	jReq := transproto.ClientVoteReq{} //SignTransactionReq{}
	log.Info("fun=ReceiveClientVoteHandler() bef--,request=%v", jReq)
	curClientVoteResp :=  transproto.ClientVoteResq{}
	G_VoteServer.CollectGroupVoteInfo(jReq)
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
		RequestInterval:4,
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
	return curVoteServer

}

//to insert value
func (this *VoteServer) CollectGroupVoteInfo(curclientvoteinfo transproto.ClientVoteReq) {
	//this.RMSGroupVotes = this.RMSGroupVotes
}
func (this *VoteServer) AddRSMGroupList(curgroupinfo transproto.RsmServerGroupResq) {
	getcurgroupinfo := curgroupinfo
	for _, _ = range getcurgroupinfo.RSMGrouList {
		//item.
	}
	gettime :=getcurgroupinfo.Startime
	rmsid := getcurgroupinfo.RSMGrouList[0][0]    //.(*string)
	clientid := getcurgroupinfo.RSMGrouList[0][1] //.(*string)
	key := transproto.RsmNode{RsmId: rmsid.(string), GroupId: getcurgroupinfo.RsmGroupId}
	var curClientRSMInfo ClientRSMInfo
	//ClientRSMInfoList,
	_,ok := this.ServerGroup[key]
	if !ok {
		curClientRSMInfo = ClientRSMInfo{Startime:gettime,ClientID: clientid.(string)}
	}
	this.ServerGroup[key] = append(this.ServerGroup[key], curClientRSMInfo)

}
func (this *VoteServer) StartRequest() {
	iserion :=0
	for {
		time.Sleep(time.Second * time.Duration(this.VoteConfigParams.RequestInterval))
		fmt.Println("cur In StartRequest()，invoke interval is:%d",this.VoteConfigParams.RequestInterval)
		//step to 请求group server：请求应在每个节点出块流程commit后

		//step to 请求可信节点列表
		getGroupInfo,err :=this.RequestRSMGroupInfo(this.VoteConfigParams.RSMServerUrl,this.VoteConfigParams.PublicKey)
		if err != nil{
			fmt.Println("cur RequestTrustInfo() num no matched!,get trustNum is:%d,get err is:%v",getGroupInfo.Rsmcount,err)
			//continue
		}

		log.Info(fmt.Sprintf("succ Invoke cur RequestTrustInfo(),get trustNum is%d,get getTrustInfo is:%v,err is:%v",getGroupInfo.Rsmcount,getGroupInfo,err))
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

func  StartVoteServer() *VoteServer{
	crustConfig :=DefaultVoteConfig()
	curVoteRsmTask:= NewVoteServer(crustConfig)
	go curVoteRsmTask.StartRequest()
	fmt.Println("StartVoteServer is start!,,gbConf' gbTrustConf.crustConfig is %v", crustConfig)
	G_VoteServer = curVoteRsmTask
	return G_VoteServer
}