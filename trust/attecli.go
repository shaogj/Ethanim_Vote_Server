package trust

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	//"github.com/trias-lab/tmware/types"
)

type TrustNodes struct {
	Pubkey string `json:"pubkey"`
	Score  int64  `json:"score"`
	IP     string `json:"ip"`
}

type attestInformation struct {
	Action  string          `json:"action"`
	Ranking [][]interface{} `json:"ranking"`
}


type AtteCli struct {
	cli http.Client
	url string
}

func NewAttestationHttpClient(urlStr string) *AtteCli {
	urlStr = strings.TrimSuffix(urlStr, "/")
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		urlStr = "http://" + urlStr
	} else {
		urlStr = urlObj.String()
	}
	urlStr += "/trias/getranking"
	return &AtteCli{
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

func (ac AtteCli) GetNodeRank() ([]TrustNodes, error) {
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
	atteInfo := attestInformation{}
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

	for _, nodeInfo := range ret.Ranking {
		score, err := strconv.Atoi(fmt.Sprintf("%s", nodeInfo[0]))
		if err != nil {
			return nil, err
		}
		trustData = append(trustData, TrustNodes{
			Score: int64(score),
			IP:    fmt.Sprintf("%v", nodeInfo[1]),
		})

	}
	return trustData, nil
}
