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
	jReq := transproto.ClientVoteReq{}
	_,sttErr := HttpExRequestJson(w, r, &jReq)
	if true != transproto.Success(sttErr) {
		GeneJsonResultFin(w, r, nil, sttErr.Code, sttErr.Desc)
		log.Error("fun=ReceiveClientVoteHandler() bef--,request=%v", jReq)
		return
	}
	//接收client msg时间.tmp 手动 -10s
	jReq.VoteTime = time.Now().Unix() - 10
	log.Info("fun=ReceiveClientVoteHandler() good,cur request=%v", jReq)
	log.Info("cur to add ServerGroup map len is:%d",len(G_VoteServer.ServerGroup))
	//0613,,to upgrade to channel msg
	G_VoteServer.CollectGroupVoteInfo(jReq)

	getresq :=transproto.ClientVoteResq{}
	getresq.ResultCode = 0
	getresq.ResultMsg = "invoke check good!"
	//0614add
	err := service.ClientVoteRecordSave(service.GXormMysql,jReq.ClientId,jReq.Rsmid,jReq.RsmGroupId,jReq.ClientSignstr,jReq.VoteTime,true)
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
	clientMsgRwlock   sync.RWMutex

	//收集每个groupID中，client对rms投票的msg
	RMSGroupVotes	map[int]transproto.RSMVoteGroupMsgs
	ClientVoteMsgQueue	chan ClientVoteReq
	qtChan             chan struct{}


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

	return &VoteServerConfig{
		//RequestInterval:4,
		//0612add,30sec 请求一次分组服务
		RequestInterval:10,
		PrivKey:	"PrivKeystr",
		PublicKey:	"PubKeystr",
		RSMServerUrl:trustUrlVerifyReal,
	}

}


func NewVoteServer(curinfo *VoteServerConfig) *VoteServer {
	curVoteServer := &VoteServer{}
	curVoteServer.VoteConfigParams = curinfo
	curVoteServer.ServerGroup = make(map[transproto.RsmNode][]ClientRSMInfo,0)
	curVoteServer.RMSGroupVotes = make(map[int]transproto.RSMVoteGroupMsgs)
	curVoteServer.ClientVoteMsgQueue = make(chan ClientVoteReq,100)
	curVoteServer.qtChan = make(chan struct{})

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


//0613,处理周期的分组结果，set VertifyResult
func (this *VoteServer) HandleGroupVotes(rsmGroupId int,rsmId string) error{

	gatherTrustCount := 0
	curservergroup,ok := this.RMSGroupVotes[rsmGroupId]
	if !ok {
		log.Error("In HandleGroupVotes(), cur rsmGroupId is exist no in RMSGroupVotes. rsmGroupId is:%d,rsmId is:%v", rsmGroupId, rsmId)
		return fmt.Errorf("cur rsmGroupId is exist no,rsmGroupId is :%s",rsmGroupId)
	}
	totalvotenum := len(this.RMSGroupVotes[rsmGroupId].ClientVote)
	for i := 0; i < totalvotenum; i++ {
		clientvoteresult := curservergroup.ClientVote[i].VertifyResult
		log.Info("In HandleGroupVotes(), cur rsmGroupId' client votemsg  is:%v", curservergroup.ClientVote[i])
		if clientvoteresult == true{
			gatherTrustCount += 1
		}
		curservergroup.ClientVertifyMap[curservergroup.ClientVote[i].ClientId] = true
	}
	log.Info("In HandleGroupVotes(), cur rsmGroupId is:%d,get gatherTrustCount is:%d,votemsg's ClientVertifyMap  is:%v", rsmGroupId,gatherTrustCount,curservergroup.ClientVertifyMap)

	curservergroup.GatherTrustCount = gatherTrustCount
	if this.RMSGroupVotes[rsmGroupId].GatherTrustCount > totalvotenum /2 {
		//this.RMSGroupVotes[rsmGroupId].GatherVertifyResult = true
		curservergroup.GatherVertifyResult = true
	}
	this.clientMsgRwlock.Lock()
	this.RMSGroupVotes[rsmGroupId] = curservergroup
	this.clientMsgRwlock.Unlock()
	//todo,save votesrecord to db
	dbVertifyResult :=0
	var clientMinorityIds string
	if curservergroup.GatherVertifyResult == true {
		dbVertifyResult =1
	}
	for clientitem,_ := range curservergroup.ClientVertifyMap {
		clientMinorityIds += clientitem
	}
	log.Info("cur rsmGroupId is :%d, after HandleGroupVotes(),get votemsg info is:%v", rsmGroupId,curservergroup)
	err := service.InsertGroupRSMVotes(service.GXormMysql,rsmGroupId,"01234",dbVertifyResult,"minorityIdslistsss",clientMinorityIds)
	if err != nil {
		log.Error("ClientVoteRecordSave(),Insert row is failed! err is-:%v \n",err) //return err
	}
	fmt.Printf("cur InsertGroupRSMVotes() Insert row finish!,groupId is :%v", rsmGroupId)

	return nil

}
func (this *VoteServer) CollectGroupVoteInfo(curclientvoteinfo transproto.ClientVoteReq) {
	//根据ServerGroup中的groupid时间，和rmsid；
	//比较msg中的VoteTime，垃圾数据则扔掉
	//else插入相关的记录中
	rsmid := curclientvoteinfo.Rsmid
	rsmGroupId := curclientvoteinfo.RsmGroupId
	voteTime := curclientvoteinfo.VoteTime
	validtimestatus := false
	key := transproto.RsmNode{GroupId: rsmGroupId,RsmId: rsmid}
	curservergroup,ok := this.ServerGroup[key]
	if !ok {
		log.Info("In CollectGroupVoteInfo(), cur key is exist no in ServerGroup.key is:%v,recv clientmsg req info is:%v",key,curclientvoteinfo)

	}else{
		//curservergroup
		log.Info("In CollectGroupVoteInfo(), cur ServerGroup' key is :%v,group's ClientRSMInfo list is:%v,recv clientmsg req info is:%v",key,"curservergroup",curclientvoteinfo)
		for _,rsmgroupitem :=range curservergroup{
			if rsmgroupitem.ClientID == curclientvoteinfo.ClientId {
				if rsmgroupitem.Startime < voteTime &&  rsmgroupitem.Endtime > voteTime{
					validtimestatus = true
					log.Info("In CollectGroupVoteInfo(), cur curclientvoteinfo is in voteTime. cur ServerGroup.key is:%v,recv clientmsg req info is:%v",key,curclientvoteinfo)
				}
			}
		}
		//非法的时间恶意投票，不处理
		/**/
		if validtimestatus == false {
			log.Error("In CollectGroupVoteInfo() judge validtimestatus is false, clientmsg is skipped.cur vote ClientId is:%v,votetime is:%v",curclientvoteinfo.ClientId,voteTime)
			return
		}
		//var curRSMVoteGroupMsgs  transproto.RSMVoteGroupMsgs
		log.Info("In CollectGroupVoteInfo(), judge validtimestatus is true, cur group is :%d, cur key is exist in ServerGroup is:%v,recv clientmsg req info is:%v",rsmGroupId,key,curclientvoteinfo)

		getVoteGroupMsgs,ok :=this.RMSGroupVotes[rsmGroupId]
		if !ok{
			curRSMVoteGroupMsgs := transproto.RSMVoteGroupMsgs{}
			curRSMVoteGroupMsgs.ClientVertifyMap = make(map[string]bool)
			curRSMVoteGroupMsgs.SeverGroupId = rsmGroupId
			curRSMVoteGroupMsgs.ServerSignStr ="ServerSignSthttpsstrinfoeee"
			log.Info("In sub RMSGroupVotes(), cur exist rsmGroupId is :%v ",rsmGroupId)
			//curRSMVoteGroupMsgs.ClientVertifyMap[curclientvoteinfo.ClientId] = true
			curRSMVoteGroupMsgs.ClientVote = append(curRSMVoteGroupMsgs.ClientVote,curclientvoteinfo)
			this.clientMsgRwlock.Lock()
			this.RMSGroupVotes[rsmGroupId] = curRSMVoteGroupMsgs
			this.clientMsgRwlock.Unlock()
			log.Info("In CollectGroupVoteInfo(), cur key is:%v ,get clientmsg is new added to RMSGroupVotes",key)
		}else{
			//todo 去重clientid
			if _,ok := getVoteGroupMsgs.ClientVertifyMap[curclientvoteinfo.ClientId];ok {
				log.Error("In CollectGroupVoteInfo(),group key is:%v, cur Client is repeat vote by RMSGroupVotes.cur ClientId is:%d",key,curclientvoteinfo.ClientId)
				return
			}
			getVoteGroupMsgs.ClientVote = append(getVoteGroupMsgs.ClientVote,curclientvoteinfo)
			this.clientMsgRwlock.Lock()
			this.RMSGroupVotes[rsmGroupId] = getVoteGroupMsgs
			this.clientMsgRwlock.Unlock()

			log.Info("In CollectGroupVoteInfo(), cur key is:%v ,get clientmsg is added to RMSGroupVotes.cur ClientVote num is:%d",key,len(getVoteGroupMsgs.ClientVote))

		}

	}
}

//transproto.RsmServerGroupResq
func (this *VoteServer) AddRSMGroupList(servergroupresq *transproto.RsmServerGroupResq,curgroupinfo []transproto.GroupItems) {
	var curClientRSMInfo ClientRSMInfo
	for _, groupitem := range curgroupinfo {
		//item.
		rmsid := groupitem.RmsId
		client := groupitem.Clients
		key := transproto.RsmNode{GroupId: servergroupresq.RsmGroupId,RsmId: rmsid}
		_,ok := this.ServerGroup[key]
		if !ok {
			for _, curclient := range client {
				curClientRSMInfo = ClientRSMInfo{Startime:servergroupresq.Startime,Endtime:servergroupresq.Endtime,ClientID: curclient}
				this.ServerGroup[key] = append(this.ServerGroup[key], curClientRSMInfo)
			}
			log.Info("in AddRSMGroupList(), key groupId with rmsid is new. to add ServerGroup' key is:%v,get this.ServerGroup[key] is:%v",key,this.ServerGroup[key])
		}else{
			log.Info("in AddRSMGroupList(), key groupId with rmsid is exist. cur groupid is:%v",key.GroupId)
			//groupId重复即为重复分组信息
			return
		}

		//this.ServerGroup[key] = append(this.ServerGroup[key], curClientRSMInfo)

	}


}

//获取分组请求
func (this *VoteServer) ReqGetRsmGrouplist() (groupid int,startime int64,err error){
	//log.Info("cur Invoke ReqGetRsmGrouplist(), to req and AddRSMGroupList is :%s","ssss444")
	ht:=utils.CHttpClientEx{}
	ht.Init()

	cli := trust.NewRSMServerCli(this.VoteConfigParams.RSMServerUrl)
	//127.0.0.1:8999
	servergroupresq, startime,err := cli.GetRsmGrouplist(ht)
	if servergroupresq == nil {
		log.Error("get cli.GetRsmGrouplist() err!,get err is:%v",err)
		return 0,0,err
	}
	//fmt.Println("after GetRsmGrouplist(),get grouprsmMap is :%v,err is:%v", servergroupresq,err)
	this.AddRSMGroupList(servergroupresq,servergroupresq.GroupItems)
	return servergroupresq.RsmGroupId,startime,nil
}

func (this *VoteServer) CalcRsmVotesProc(interval int) {
	log.Info("启动投票汇总服务，interval is:%d", interval)
	//go StartRobotRegisterTask(robotRegCfgBuffer)
	//20 *6 = 2min进行汇总计算
	ticker := time.NewTicker(time.Duration(interval) * time.Second * 6)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			curmin := time.Now().Minute()
			newcyclestart, _ := service.GetDayTodayLastMinute(time.Now().Unix(), curmin)
			log.Info("run 定时task CalcRsmVotesProc()，curmin is:%d,newcyclestart is:%d，this.ServerGroup len is:%d",curmin,newcyclestart,len(this.ServerGroup))
			//to do,,获取上一个周期时间的投票分组信息,to single队列任务
			for groupkey,groupitem := range this.ServerGroup {
				log.Info("cur in CalcRsmVotesProc(),groupkey is:%v,to HandleGroupVotes() for groupitem is:%v",groupkey,groupitem)
				error :=this.HandleGroupVotes(groupkey.GroupId,groupkey.RsmId)
				if error != nil {
					log.Error("cur in CalcRsmVotesProc(),GroupId: %v is exist no.newcyclestart is:%v",groupkey.GroupId,newcyclestart)
				}
			}

		case <-this.qtChan:
			ticker.Stop()
			return
		}


	}
}

func (this *VoteServer) Quit() {
	this.qtChan <- struct{}{}
}

func (this *VoteServer) StartRequest() {
	iserion :=0
	for {
		//time.Sleep(time.Second * time.Duration(this.VoteConfigParams.RequestInterval))
		fmt.Printf("cur In StartRequest()，invoke interval is:%d",this.VoteConfigParams.RequestInterval)
		//step to 请求group server：请求应在每个节点出块流程commit后

		//step to 请求可信节点列表
		//getGroupInfo,err :=this.RequestRSMGroupInfo(this.VoteConfigParams.RSMServerUrl,this.VoteConfigParams.PublicKey)
		getGroupId,startime,err :=this.ReqGetRsmGrouplist()
		if err != nil || startime ==0 {
			log.Error(fmt.Sprintf("get ReqGetRsmGrouplist() err!,getGroupId is:%d,get err is:%v",getGroupId,err))
			//continue
		}
		var nowtime int64 = time.Now().Unix()
		log.Info(fmt.Sprintf("cur get ReqGetRsmGrouplist(),getGroupId is:%d,startime:%d，nowtime is：%d",getGroupId,startime,nowtime))
		//DurationOfTime:= time.Duration(3) * time.Second
		var curstrtime string = time.Unix(nowtime, 0).Format("2006-01-02 15:04:05")
		sendFunc2 := func() {
			var curstrtime2 string = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("checking task time.AfterFunc(),cur group is :%d,startime is:%d,curstrtime is:%d，invokeFunc time is:%d\n",getGroupId,startime,curstrtime,curstrtime2)
			//todo
			this.HandleGroupVotes(getGroupId,"")
		}
		//0616testing,to update:
		//waitTime := time.Duration(signaling.Timestamp-time.Now().UnixNano()/1e6) * time.Millisecond
		lastcyclestart, lastcycleend := service.GetDayTodayLastMinute(time.Now().Unix(), time.Now().Minute())
		realdelaytime := lastcycleend + 1 - startime + 60
		fmt.Printf("cur time is:%d,lastcyclestart is:%d,lastcycleend is :%d,realdelaytime is:%d", time.Now().Unix(),lastcyclestart,lastcycleend,realdelaytime)

		waitTime :=time.Duration(40)	* time.Second //time.Now().Unix()
		time.AfterFunc(waitTime, sendFunc2)
		log.Info(fmt.Sprintf("succ Invoke cur RequestTrustInfo(),get total group num is :%d,cur Groupid is: %d,err is:%v",len(this.ServerGroup),getGroupId,err))

		iserion ++
		//this.clientMsgRwlock.Lock()
		//11.25,返回的可信node列表，比现有的列表少，则设置不可信的节点score为-1，及为不可信的节点
		//this.TrustNodeMap = trustnodeaddr
		//this.clientMsgRwlock.Unlock()
		time.Sleep(time.Second * time.Duration(this.VoteConfigParams.RequestInterval))

	}
}
var G_VoteServer  *VoteServer

func  StartVoteServer(groupserverurl string,requestInterval int) *VoteServer{
	crustConfig :=DefaultVoteConfig()
	crustConfig.RSMServerUrl = groupserverurl
	crustConfig.RequestInterval = requestInterval
	curVoteRsmTask:= NewVoteServer(crustConfig)
	go curVoteRsmTask.StartRequest()
	//go curVoteRsmTask.CalcRsmVotesProc(requestInterval)
	fmt.Println("StartVoteServer is start!,,gbConf' gbTrustConf.crustConfig is %v", crustConfig)
	G_VoteServer = curVoteRsmTask
	return G_VoteServer
}