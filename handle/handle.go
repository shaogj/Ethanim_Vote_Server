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
	//for Testing :接收client msg时间.tmp 手动 -10s
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
	//每次从server端获取rms分组的client列表
	ServerGroup map[transproto.RsmNode][]ClientRSMInfo
	clientMsgRwlock   sync.RWMutex

	//收集每个groupID中，client对rms投票的msg
	RMSGroupVotes	map[int]transproto.RSMVoteGroupMsgs
	//upgrade
	RMSGroupVotesSubRsm	map[int]map[string]transproto.RSMVoteGroupMsgs

	ClientVoteMsgQueue	chan *ClientVoteReq
	qtChan             chan struct{}
	//最新的RSM分组id
	LastestGroupId	int64


}
//0605,RSM分组的列表消息
type ClientRSMInfo struct {
	Startime int64
	Endtime int64
	ClientID string

}

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
	//0617,,enhance to replace above map
	curVoteServer.RMSGroupVotesSubRsm = make(map[int]map[string]transproto.RSMVoteGroupMsgs)

	curVoteServer.ClientVoteMsgQueue = make(chan *ClientVoteReq,100)
	curVoteServer.qtChan = make(chan struct{})

	return curVoteServer

}
func (this *VoteServer) RecvClientVoteMsg(curClientVoteReq *ClientVoteReq) error {
	this.ClientVoteMsgQueue <- curClientVoteReq
	log.Info("cur RecvClientVote() to queue,client votemsg  is:%v", curClientVoteReq)
	return nil
}


//to insert value
//0613，get mysql data:

type MajorityIds struct {
	GroupId	int  `json:"group_id"`
	ClientId []string `json:"majority_ids"`
}

type MinorityIds struct {
	GroupId	int  `json:"group_id"`
	ClientId []string `json:"minority_ids"`
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
	dbVertifyResult :=0
	curMajorityIds := MajorityIds{GroupId:rsmGroupId}

	if curservergroup.GatherVertifyResult == true {
		dbVertifyResult =1
	}
	for clientitem,_ := range curservergroup.ClientVertifyMap {
		//clientMinorityIds += clientitem
		curMajorityIds.ClientId = append(curMajorityIds.ClientId,clientitem)
	}
	minorityIdsJson, err := json.Marshal(curMajorityIds)
	if err != nil {
		log.Error("cur rsmGroupId is :%d, curMajorityIds is:%s, Marshal err is:%v", rsmGroupId,curMajorityIds,err)
		return err
	}
	log.Info("cur rsmGroupId is :%d, after HandleGroupVotes(),get votemsg info is:%v", rsmGroupId,curservergroup)
	//0617add
	curservergrouprsm := this.RMSGroupVotesSubRsm[rsmGroupId]
	log.Info("cur rsmGroupId is :%d, after HandleGroupVotes(),cur RMSGroupVotesSubRsm info is:%v", rsmGroupId,curservergrouprsm)
	//0620
	for rsmid,curservergrouprsm:= range curservergrouprsm {
		//curservergrouprsm.ClientVertifyMap
		log.Info("cur rsmid :%s in RMSGroupVotesSubRsm(),cur curservergrouprsm info is:%v", rsmid,curservergrouprsm)
		//to move add handle
	}
	//save votesrecord to db
	err = service.InsertGroupRSMVotes(service.GXormMysql,rsmGroupId,rsmId,dbVertifyResult,"minorityIdslistsss",string(minorityIdsJson),"untrust_vote_idsinfo","333",4)//clientMinorityIds
	if err != nil {
		log.Error("ClientVoteRecordSave(),Insert row is failed! err is-:%v \n",err)
		return err
	}
	fmt.Printf("cur InsertGroupRSMVotes() Insert row finish!,groupId is :%v", rsmGroupId)

	return nil
}
//汇总处理分组的client投票数据
func (this *VoteServer) HandleGroupRSMVotes(rsmGroupId int) error {
	gatherTrustCount := 0
	curservergrouprsm,ok := this.RMSGroupVotesSubRsm[rsmGroupId]
	if !ok {
		log.Error("In HandleGroupRSMVotes(), client vote of rsmGroupId is exist no in RMSGroupVotesSubRsm. rsmGroupId is:%d", rsmGroupId)
		return fmt.Errorf("cur rsmGroupId is exist no,rsmGroupId is :%s",rsmGroupId)
	}
	log.Info("cur rsmGroupId is :%d, after HandleGroupRSMVotes(),cur RMSGroupVotesSubRsm for client voteinfo is:%v", rsmGroupId,curservergrouprsm)
	//0620
	for rsmid,curgrouprsmvotes:= range curservergrouprsm {
		totalvotenum := len(curgrouprsmvotes.ClientVote)
		//获取当前组对RSM未投票的clients
		getrsmclients,_ := this.GetGroupIdRsmClients(rsmGroupId,rsmid)
		rsmreqnoclients := getrsmclients
		for i := 0; i < totalvotenum; i++ {
			clientvoteresult := curgrouprsmvotes.ClientVote[i].VertifyResult
			log.Info("In HandleGroupRSMVotes(), cur rsmGroupId' client votemsg  is:%v", curgrouprsmvotes.ClientVote[i])
			if clientvoteresult == true{
				gatherTrustCount += 1
				curgrouprsmvotes.ClientVertifyMap[curgrouprsmvotes.ClientVote[i].ClientId] = true
			}else{
				curgrouprsmvotes.ClientVertifyMap[curgrouprsmvotes.ClientVote[i].ClientId] = false
			}
			_,ok := rsmreqnoclients[curgrouprsmvotes.ClientVote[i].ClientId]
			if ok {
				//对RSM投票的client做标记
				rsmreqnoclients[curgrouprsmvotes.ClientVote[i].ClientId] = 1
			}

		}
		log.Info("In HandleGroupRSMVotes(), cur rsmGroupId is:%d,get gatherTrustCount is:%d,votemsg's ClientVertifyMap  is:%v", rsmGroupId,gatherTrustCount,curgrouprsmvotes.ClientVertifyMap)

		curgrouprsmvotes.GatherTrustCount = gatherTrustCount
		if curgrouprsmvotes.GatherTrustCount >= totalvotenum /2 {
			curgrouprsmvotes.GatherVertifyResult = true
		}
		this.clientMsgRwlock.Lock()
		this.RMSGroupVotesSubRsm[rsmGroupId][rsmid] = curgrouprsmvotes
		this.clientMsgRwlock.Unlock()
		//MinorityIds
		dbVertifyResult :=0
		curMajorityIds := MajorityIds{GroupId:rsmGroupId}
		curMinorityIds := MinorityIds{GroupId:rsmGroupId}

		if curgrouprsmvotes.GatherVertifyResult == true {
			dbVertifyResult =1
		}
		for clientitem,vertifyresult:= range curgrouprsmvotes.ClientVertifyMap {
			//clientMinorityIds += clientitem
			if vertifyresult == true{
				curMajorityIds.ClientId = append(curMajorityIds.ClientId,clientitem)
			}else{
				curMinorityIds.ClientId = append(curMinorityIds.ClientId,clientitem)
			}
		}
		//未投票的clients的value为0
		var slack_vote_ids string = "["
		for clientid,vertifystatus:= range rsmreqnoclients {
			if vertifystatus == 0{
				//curMinorityIds.ClientId = append(curMinorityIds.ClientId,clientid)
				if len(slack_vote_ids) < 2 {
					slack_vote_ids = clientid
				}else{
					slack_vote_ids = fmt.Sprintf("%s,%s",slack_vote_ids,clientid)
				}
			}
		}
		slack_vote_ids = slack_vote_ids + "]"
		log.Info("cur groupId is :%d,check slack_vote_ids is:%s", rsmGroupId,slack_vote_ids)

		//0621
		majorityIdsJson, err := json.Marshal(curMajorityIds.ClientId)
		minorityIdsJson, err := json.Marshal(curMinorityIds.ClientId)

		log.Info("cur rsmId is :%s, dbVertifyResult is:%d,after HandleGroupVotes(),get votemsg info is:%v,err is:%v", rsmid,dbVertifyResult,curgrouprsmvotes,err)
		curwalletuserrecord := service.GetGroupRSMVotesMsgs(rsmGroupId,rsmid,dbVertifyResult,string(majorityIdsJson),string(minorityIdsJson),slack_vote_ids)
		//调用节点RPC写入数据上链
		curwalletuserrecordmsg, err := json.Marshal(*curwalletuserrecord)
		getrandnumstr := RandStr(5)
		fmt.Println("get randnumstr7777 is:%s",getrandnumstr)
		log.Info("cur to SendRMSVoteMsgToNode(),rsmGroupId is:%d,user votes' record msg is:%v",rsmGroupId,curwalletuserrecordmsg)

		sendrecordmsgs := "tx=%22" + string(curwalletuserrecordmsg) + getrandnumstr + "%22"
		getblockInfo,getresq,err := SendRMSVoteMsgToNode(this.VoteConfigParams.TMNodeUrl,sendrecordmsgs)
		if err != nil{
			log.Error("cur after CommitRMSVoteMsg(),get error! ,getRespInfo is :%v,err is:%v",getblockInfo,err)
		}
		log.Info("cur CommitRMSVoteMsg(),get getresq hash is:%s,height is:%d",getresq.Hash,getresq.Height)

		err = service.InsertGroupRSMVotes(service.GXormMysql,rsmGroupId,rsmid,dbVertifyResult,string(majorityIdsJson),string(minorityIdsJson),slack_vote_ids,string(getresq.Hash),getresq.Height)//clientMinorityIds
		if err != nil {
			log.Error("ClientVoteRecordSave(),Insert row is failed! err is-:%v \n",err)
		}
		log.Info("cur InsertGroupRSMVotes() Insert row finish!,groupId is :%v,rsmid is:%s", rsmGroupId,rsmid)


	}
	log.Info("cur groupid :%d in RMSGroupVotesSubRsm(),cur curservergrouprsm info is:%v",rsmGroupId,curservergrouprsm)
	return nil

}

//获取分组rsm的所以投票clients
func (this *VoteServer) GetGroupIdRsmClients(rsmGroupId int,rsmid string) (clients map[string]int,err error){
	var rsmclients = make(map[string]int)
	key := transproto.RsmNode{GroupId: rsmGroupId,RsmId: rsmid}

	curservergroup,ok := this.ServerGroup[key]
	if !ok {
		return nil, fmt.Errorf("cur rsmGroupId is exist no in ServerGroup,rsmGroupId is :%s",rsmGroupId)
	}else{
		for _,rsmgroupitem :=range curservergroup{
			//rsmclients = append(rsmclients,rsmgroupitem.ClientID)
			rsmclients[rsmgroupitem.ClientID] = 0
		}
	}
	log.Info("In GetGroupIdRsmClients(), cur ServerGroup' key is :%v,get rms' rsmclients is:%v",key,rsmclients)
	return rsmclients,nil
}

//归集分组rsmid的client请求的投票数据
func (this *VoteServer) CollectGroupVoteInfo(curclientvoteinfo transproto.ClientVoteReq) {
	//根据ServerGroup中的groupid时间，和rmsid，比较msg中的VoteTime，垃圾数据则扔掉
	rsmid := curclientvoteinfo.Rsmid
	rsmGroupId := curclientvoteinfo.RsmGroupId
	voteTime := curclientvoteinfo.VoteTime
	validtimestatus := false
	key := transproto.RsmNode{GroupId: rsmGroupId,RsmId: rsmid}
	curservergroup,ok := this.ServerGroup[key]
	if !ok {
		log.Error("In CollectGroupVoteInfo(), cur clientVoteReq's key is exist no in ServerGroup.key is:%v,recv clientmsg req info is:%v",key,curclientvoteinfo)

	}else{
		log.Info("In CollectGroupVoteInfo(), cur ServerGroup' key is :%v,group's ClientRSMInfo list is:%v,recv clientmsg req info is:%v",key,"curservergroup",curclientvoteinfo)
		for _,rsmgroupitem :=range curservergroup{
			if rsmgroupitem.ClientID == curclientvoteinfo.ClientId {
				if rsmgroupitem.Startime < voteTime &&  rsmgroupitem.Endtime > voteTime{
					validtimestatus = true
					log.Info("In CollectGroupVoteInfo(), cur curclientvoteinfo is in voteTime. cur ServerGroup.key is:%v,recv clientmsg req info is:%v",key,curclientvoteinfo)
				}
			}
		}
		//非法时间的或恶意投票，不处理
		if validtimestatus == false {
			log.Error("In CollectGroupVoteInfo() judge validtimestatus is false, clientmsg is skipped.cur rsmgroupid is:%d,vote ClientId is:%v,votetime is:%v",rsmGroupId,curclientvoteinfo.ClientId,voteTime)
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
			curRSMVoteGroupMsgs.ClientVote = append(curRSMVoteGroupMsgs.ClientVote,curclientvoteinfo)
			this.clientMsgRwlock.Lock()
			this.RMSGroupVotes[rsmGroupId] = curRSMVoteGroupMsgs
			this.clientMsgRwlock.Unlock()
			log.Info("In CollectGroupVoteInfo(), cur key is:%v ,get clientmsg is new added to RMSGroupVotes",key)
			//0617
			this.RMSGroupVotesSubRsm[rsmGroupId] = make(map[string]transproto.RSMVoteGroupMsgs)
			subrsmvotemsgs :=this.MakeVoteGroupMsg(rsmGroupId,rsmid,curclientvoteinfo)
			log.Info("In CollectGroupVoteInfo(), cur key is:%v ,cur ClientVote info:%v",key,*subrsmvotemsgs)
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
			//0617enhance
			subrmsclientmsgs,ok :=this.RMSGroupVotesSubRsm[rsmGroupId][rsmid]
			if !ok {
				subrsmvotemsgs :=this.MakeVoteGroupMsg(rsmGroupId,rsmid,curclientvoteinfo)
				log.Info("In CollectGroupVoteInfo(), cur key is:%v ,cur ClientVote info:%v",key,*subrsmvotemsgs)
				//this.RMSGroupVotesSubRsm[rsmGroupId][rsmid] = getVoteGroupMsgs
			}else{
				if _,ok := subrmsclientmsgs.ClientVertifyMap[curclientvoteinfo.ClientId];ok {
					log.Error("In CollectGroupVoteInfo(),sub group key is:%v, cur Client is repeat vote by RMSGroupVotes.cur ClientId is:%d",key,curclientvoteinfo.ClientId)
					return
				}
				subrmsclientmsgs.ClientVote = append(subrmsclientmsgs.ClientVote,curclientvoteinfo)
				this.RMSGroupVotesSubRsm[rsmGroupId][rsmid] = subrmsclientmsgs
			}
		}

	}
}
func (this *VoteServer) MakeVoteGroupMsg(rsmgroupid int,rsmid string,curclientvoteinfo transproto.ClientVoteReq)(*transproto.RSMVoteGroupMsgs){
	curRSMVoteGroupMsgs := transproto.RSMVoteGroupMsgs{}
	curRSMVoteGroupMsgs.ClientVertifyMap = make(map[string]bool)
	curRSMVoteGroupMsgs.SeverGroupId = rsmgroupid
	curRSMVoteGroupMsgs.ServerSignStr ="subMakeVoteGroupMsgstrinfo"
	curRSMVoteGroupMsgs.ClientVote = append(curRSMVoteGroupMsgs.ClientVote,curclientvoteinfo)
	this.RMSGroupVotesSubRsm[rsmgroupid][rsmid] = curRSMVoteGroupMsgs
	log.Info("In sub MakeVoteGroupMsg(), cur rsmgroupid is:%d, rsmid is:%s",rsmgroupid,rsmid)
	return &curRSMVoteGroupMsgs
}

//添加新的分组rsm信息列表
func (this *VoteServer) AddRSMGroupList(servergroupresq *transproto.RsmServerGroupResq,curgroupinfo []transproto.GroupItems) {
	var curClientRSMInfo ClientRSMInfo
	for _, groupitem := range curgroupinfo {
		rmsid := groupitem.RmsId
		client := groupitem.Clients
		key := transproto.RsmNode{GroupId: servergroupresq.RsmGroupId,RsmId: rmsid}
		_,ok := this.ServerGroup[key]
		if !ok {
			for _, curclient := range client {
				curClientRSMInfo = ClientRSMInfo{Startime:servergroupresq.Startime,Endtime:servergroupresq.Endtime,ClientID: curclient}
				this.ServerGroup[key] = append(this.ServerGroup[key], curClientRSMInfo)
			}
			this.LastestGroupId = int64(servergroupresq.RsmGroupId)
			log.Info("in AddRSMGroupList(), key groupId with rmsid is new. LastestGroupId is:%d,to add ServerGroup' key is:%v,get this.ServerGroup[key] is:%v",this.LastestGroupId,key,this.ServerGroup[key])
		}else{
			log.Info("in AddRSMGroupList(), key groupId with rmsid is exist. cur groupid is:%v",key.GroupId)
			//groupId重复即为重复分组信息
			return
		}
	}
}

//获取分组请求
func (this *VoteServer) ReqGetRsmGrouplist() (groupid int,startime int64,err error){
	ht:=utils.CHttpClientEx{}
	ht.Init()

	cli := trust.NewRSMServerCli(this.VoteConfigParams.RSMServerUrl)
	servergroupresq, startime,err := cli.GetRsmGrouplist(ht)
	if servergroupresq == nil {
		log.Error("get cli.GetRsmGrouplist() err!,get err is:%v",err)
		return 0,0,err
	}
	//fmt.Println("after GetRsmGrouplist(),get grouprsmMap is :%v,err is:%v", servergroupresq,err)
	this.AddRSMGroupList(servergroupresq,servergroupresq.GroupItems)
	return servergroupresq.RsmGroupId,startime,nil
}

//清理历史数据
func (this *VoteServer) ResetGroupVotesMap(interval int) {
	log.Info("启动投票汇总服务，interval is:%d", interval)
	//20 *RequestInterval = 20个分组周期清理历史数据
	ticker := time.NewTicker(time.Second * time.Duration(this.VoteConfigParams.RequestInterval) * 20)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			curmin := time.Now().Minute()
			newcyclestart, _ := service.GetDayTodayLastMinute(time.Now().Unix(), curmin)
			log.Info("run ticker task ResetGroupVotesMap()，curmin is:%d,newcyclestart is:%d，this.ServerGroup len is:%d,cur LastestGroupId is:%d",curmin,newcyclestart,len(this.ServerGroup),this.LastestGroupId)
			//to do,,获取上一个周期时间的投票分组信息,to single队列任务
			for groupkey,groupitem := range this.ServerGroup {
				if groupkey.GroupId != int(this.LastestGroupId) {
					delete(this.ServerGroup, groupkey)
					log.Debug("cur in ResetGroupVotesMap(), old groupkey is:%v,ServerGroup map,RMSGroupVotesSubRsm map for the groupitem is removed ",groupkey)
					delete(this.RMSGroupVotes,groupkey.GroupId)
					delete(this.RMSGroupVotesSubRsm,groupkey.GroupId)
				}
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

//请求可信RSM分组服务
func (this *VoteServer) StartRequestRsmGroup() {
	for {
		log.Info("cur In StartRequestRsmGroup()，invoke interval is:%d",this.VoteConfigParams.RequestInterval)
		//to enhance by VoteConfigParams.PublicKey
		getGroupId,startime,err :=this.ReqGetRsmGrouplist()
		if err != nil || startime ==0 {
			log.Error(fmt.Sprintf("get ReqGetRsmGrouplist() err!,getGroupId is:%d,get err is:%v",getGroupId,err))
			time.Sleep(time.Second * time.Duration(this.VoteConfigParams.RequestInterval) /2 )
			continue
		}
		var nowtime int64 = time.Now().Unix()
		log.Info(fmt.Sprintf("cur get ReqGetRsmGrouplist(),getGroupId is:%d,startime:%d，nowtime is：%d",getGroupId,startime,nowtime))
		var curstrtime string = time.Unix(nowtime, 0).Format("2006-01-02 15:04:05")

		//分组周期结束后，写入统计执行任务
		sendFunc2 := func() {
			var curstrtime2 string = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
			log.Info("checking task time.AfterFunc(),cur group is :%d,startime is:%d,curstrtime is:%s，invokeFunc time is:%s\n",getGroupId,startime,curstrtime,curstrtime2)
			//update to rmsid
			this.HandleGroupRSMVotes(getGroupId)

		}
		//waitTime := time.Duration(signaling.Timestamp-time.Now().UnixNano()/1e6) * time.Millisecond
		lastcyclestart, lastcycleend := service.GetDayTodayLastMinute(time.Now().Unix(), time.Now().Minute())
		realdelaytime := lastcycleend + 1 - startime + 60
		log.Info("cur time is:%d,lastcyclestart is:%d,lastcycleend is :%d,realdelaytime is:%d", time.Now().Unix(),lastcyclestart,lastcycleend,realdelaytime)

		//waitTime :=time.Duration(40)	* time.Second
		waitTime :=time.Duration(this.VoteConfigParams.RequestInterval)	* time.Second
		//启动任务在分组周期后
		time.AfterFunc(waitTime, sendFunc2)
		log.Info(fmt.Sprintf("succ Invoke cur RequestTrustInfo(),get total group num is :%d,cur Groupid is: %d,err is:%v",len(this.ServerGroup),getGroupId,err))

		time.Sleep(time.Second * time.Duration(this.VoteConfigParams.RequestInterval))

	}
}
var G_VoteServer  *VoteServer

//开启服务分组数据请求
func  StartVoteServer(groupserverurl string,requestInterval int,nodeUrl string) *VoteServer{
	trustConfig :=DefaultVoteConfig()
	trustConfig.RSMServerUrl = groupserverurl
	trustConfig.RequestInterval = requestInterval
	trustConfig.TMNodeUrl = nodeUrl
	curVoteRsmTask:= NewVoteServer(trustConfig)
	go curVoteRsmTask.StartRequestRsmGroup()
	go curVoteRsmTask.ResetGroupVotesMap(requestInterval)
	fmt.Println("StartVoteServer is start!,,gbConf' gbTrustConf.trus;tConfig is %v", trustConfig)
	G_VoteServer = curVoteRsmTask
	return G_VoteServer
}