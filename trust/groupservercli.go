package trust

import (
	"Ethanim_Vote_Server/proto"
	"Ethanim_Vote_Server/utils"
	"fmt"
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

//*[]proto.GroupItems,
func (ac GroupServerCli) GetRsmGrouplist(ht utils.CHttpClientEx) (servergroupresq *proto.RsmServerGroupResq ,startime int64,err error) {
	//0607add
	//ht:=utils.CHttpClientEx{}
	//ht.Init()
	ht.HeaderSet("Content-Type", "text/json")

	atteInfo := proto.RsmServerGroupResq{}
	//atteInfo.GroupAttach = &groupsinfo
	reqServerInfo :=  proto.RsmServerGroupReq{}
	_,statusCode,errorCode,err:=ht.RequestJsonResponseJson(ac.url,3000,&reqServerInfo,&atteInfo)
	if nil!=err {
		utils.LogErrorf("coinType :%s,ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ","coinType",statusCode,errorCode,err,ac.url)
		sttError :=fmt.Sprintf("GetReqtoGroupServer{Error=%d,Desc=%v}",errorCode,err)
		fmt.Printf("in GetRsmGrouplist(),sttError is:",sttError)
		return &atteInfo,0,err
	}
	utils.LogInfof("ac.url :%s,GetRsmGrouplist() res=%v",ac.url,atteInfo)

	/*
	for _, RsmInfo := range atteInfo.GroupItems{
		log.Info("check001 in GetRsmGrouplist(),cur RsmInfo is:%s,Clients list is:%v",RsmInfo.RmsId,RsmInfo.Clients)
	}
	*/
	return &atteInfo, atteInfo.Startime,nil
}
