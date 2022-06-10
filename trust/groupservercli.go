package trust

import (
	"Ethanim_Vote_Server/proto"
	"Ethanim_Vote_Server/utils"
	"fmt"
	"github.com/mkideal/log"
	"net/http"
	"net/url"
	"strings"
	"time"

	//"github.com/trias-lab/tmware/types"
)


type GroupServerCli struct {
	cli http.Client
	url string
}

func NewRSMServerCli(urlStr string) *GroupServerCli {
	urlStr = strings.TrimSuffix(urlStr, "/")
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		urlStr = "http://" + urlStr
	} else {
		urlStr = urlObj.String()
	}
	urlStr += "/rsm/groups"
	//urlStr += "/trias/getranking"
	return &GroupServerCli{
		url: urlStr,
		cli: http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    20 * time.Second,
				DisableCompression: true,
			},
			Timeout: 5 * time.Second,
		},
	}
}

func (ac GroupServerCli) GetRsmGrouplist(ht utils.CHttpClientEx) (*[]proto.GroupItems, error) {
	//0607add
	//ht:=utils.CHttpClientEx{}
	//ht.Init()
	ht.HeaderSet("Content-Type", "text/json")

	//atteInfo := attestInformation{}
	atteInfo := proto.RsmServerGroupResq{}
	//groupsinfo :=proto.GroupAttachRSM{}
	//atteInfo.GroupAttach = &groupsinfo
	reqServerInfo :=  proto.RsmServerGroupReq{}
	_,statusCode,errorCode,err:=ht.RequestJsonResponseJson(ac.url,5000,&reqServerInfo,&atteInfo)
	if nil!=err {
		utils.LogErrorf("coinType :%s,ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ","coinType",statusCode,errorCode,err,ac.url)
		sttError :=fmt.Sprintf("GetReqtoGroupServer{Error=%d,Desc=%v}",errorCode,err)
		fmt.Printf("in GetRsmGrouplist(),sttError is:",sttError)
		return nil,err
	}
	utils.LogInfof("ac.url :%s,GetRsmGrouplist() res=%v",ac.url,atteInfo)
	/*
	resp, err := ac.cli.Get(ac.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resStr := string(respBytes)
	rankStr := strings.ReplaceAll(resStr, "'", "\"")
	err = json.Unmarshal([]byte(rankStr), &atteInfo)
	if err != nil {
		err = fmt.Errorf("cannot parse data on attestInformation[%s] from %s,%v", resStr, ac.url, err)
		return nil, err
	}
	ret := attestInformation{}
	//err = json.Unmarshal(respBytes, &ret)
	//Fix bug for strconv.Atoi err as nodeInfo[0] turn to "1.072842e+06
	d := json.NewDecoder(strings.NewReader(string(respBytes)))
	d.UseNumber()
	err = d.Decode(&ret)
	if err != nil {
		return nil, err
	}
	var trustData []TrustNodes
	*/

	/*for _, nodeInfo := range ret.Ranking {
		score, err := strconv.Atoi(fmt.Sprintf("%s", nodeInfo[0]))
		if err != nil {
			return nil, err
		}
		trustData = append(trustData, TrustNodes{
			Score: int64(score),
			IP:    fmt.Sprintf("%v", nodeInfo[1]),
		})

	}
	*/
	for _, RsmInfo := range atteInfo.GroupItems{
		log.Info("check001 in GetRsmGrouplist(),cur RsmInfo is:%s,Clients list is:%v",RsmInfo.RmsId,RsmInfo.Clients)
	}
	return &atteInfo.GroupItems, nil
}
