package handle

import (
	"Ethanim_Vote_Server/proto"
	tmjson "Ethanim_Vote_Server/utils/json"
	"github.com/mkideal/log"

	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"encoding/json"
	crand "crypto/rand"
	"sync"
	mrand "math/rand"
)

const (
	strChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // 62 characters
)
type Rand struct {
	sync.Mutex
	rand *mrand.Rand
}
var grand *Rand

func cRandBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
func InitRand() *mrand.Rand{
	bz := cRandBytes(8)
	var seed uint64
	for i := 0; i < 8; i++ {
		seed |= uint64(bz[i])
		seed <<= 8
	}
	rand := mrand.New(mrand.NewSource(int64(seed)))
	return rand
}
func RandStr(length int) string {
	chars := []byte{}
	//var rand mrand.Rand
	rand := InitRand()

MAIN_LOOP:
	for {
		i63 := rand.Int63()
		val := i63
		//r.Int63()
		for i := 0; i < 10; i++ {
			v := int(val & 0x3f) // rightmost 6 bits
			if v >= 62 {         // only 62 characters in strChars
				val >>= 6
				continue
			} else {
				chars = append(chars, strChars[v])
				if len(chars) == length {
					break MAIN_LOOP
				}
				val >>= 6
			}
		}
	}

	return string(chars)
}

//1129add,to auto get local server IP
func ExternalIP() (net.IP, error) {
ifaces, err := net.Interfaces()
if err != nil {
return nil, err
}
for _, iface := range ifaces {
if iface.Flags&net.FlagUp == 0 {
continue // interface down
}
if iface.Flags&net.FlagLoopback != 0 {
continue // loopback interface
}
addrs, err := iface.Addrs()
if err != nil {
return nil, err
}
for _, addr := range addrs {
ip := getIpFromAddr(addr)
if ip == nil {
continue
}
return ip, nil
}
}
return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
var ip net.IP
switch v := addr.(type) {
case *net.IPNet:
ip = v.IP
case *net.IPAddr:
ip = v.IP
}
if ip == nil || ip.IsLoopback() {
return nil
}
ip = ip.To4()
if ip == nil {
return nil // not an ipv4 address
}

return ip
}


func ReqGetUrl(vertify string)(respstr string,err error){
	resp, err := http.Get(vertify)
	if err != nil {
		fmt.Println("when trustQuery,Marshal to json error:%s", err.Error())
		return "",err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("when trustQuery,Marshal to json error:%s", err.Error())
		return "",err
	}
	return string(body),nil
}

func UnmarshalResponseBytes(responseBytes []byte, result interface{}) (interface{}, error) {
	var err error
	response := &proto.RPCResponse{}
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshalling rpc response: %v", err))
	}
	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("response error: %v", response.Error))
	}
	//err = cdc.UnmarshalJSON(response.Result, result)
	err = tmjson.Unmarshal(response.Result,result)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshalling rpc response result: %v", err))
	}
	return result, nil
}

func SendRMSVoteMsgToNode(nodeUrl,sendmsg string)(getblockInfo string,getresq *proto.ResultBroadcastTxCommit,err error){
	sendinfo := fmt.Sprintf("%s%s",nodeUrl,sendmsg)
	//0630
	getRespInfo,err :=ReqGetUrl(sendinfo)
	if err != nil{
		log.Error("cur after CommitRMSVoteMsg(),get error! ,getRespInfo is :%v,err is:%v",getRespInfo,err)
		return "",nil,err
	}
	log.Info("cur invoke CommitRMSVoteMsg() succ! getRespInfo is:%v",getRespInfo)
	respInfo :=proto.ResultBroadcastTxCommit{}
	getrrrInfo,err :=UnmarshalResponseBytes([]byte(getRespInfo),&respInfo)
	if  nil!=err  {
		log.Error("resp=%s,url=%s,err=%v",string(getRespInfo),nodeUrl,err.Error())
		return getRespInfo,&respInfo,err
	}
	log.Info("cur json.Unmarshal succ. respInfo is:%v,getrrrInfo is:%v",respInfo,getrrrInfo)

	return getRespInfo,&respInfo,nil
}