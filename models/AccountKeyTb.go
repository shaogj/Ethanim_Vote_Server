
package models
var (
	//账户信息表 wdc表存放keystore对应的字段值表
	TableBTCAccount="gjc_account_key_tb"
	TableCoinPrivateKey= "coin_private_key"
	TableGGEXTranRecord="ggex_tran_state"
	//0116add
	TableBTCTranRecord="btc_tran_state"
	//0614
	TableClientVoteRecord="client_vote_rsm_record"


)

type GjcAccountKeyTb struct {
	Id          int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	//Uid           int    `json:"uid" xorm:"not null index INT(11)"`
	AccountName        string `json:"accountname" xorm:"default ''"`
	CoinType string    `json:"cointype"` //交易币种类
	Walletid  	int64 `json:"walletId" xorm:"BIGINT(20)"`
	//0926 add
	//Wallettype string    `xorm:"CHAR(32)"`
	PrivKey        string `json:"privkey" xorm:"not null TEXT"`
	PubKey      string `json:"pubkey" xorm:"not null TEXT"`
	AddressId        string `json:"addressid" xorm:"not null TEXT"`
	//Txid,is made by last txout, to pay to for next time
	Utxoid      string `json:"utxoid" xorm:"default '' TEXT"`
	CreatedTime   int64  `json:"created_time" xorm:"BIGINT(20)"`
	Status           int    `json:"status" xorm:"default 0 index INT(11)"`
	UpdatedTime   int64  `json:"updated_time" xorm:"BIGINT(20)"`
}
/*
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for client_vote_rsm_record
-- ----------------------------
DROP TABLE IF EXISTS `client_vote_rsm_record`;
CREATE TABLE `client_vote_rsm_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `clientid` varchar(120) DEFAULT NULL COMMENT 'clientid',
  `rsm_group_id` bigint(20) NOT NULL COMMENT '分组ID',
  `rsmid` varchar(40) DEFAULT NULL COMMENT 'rsmid',
  `vertify_result` varchar(4) DEFAULT NULL COMMENT '验证结果',
  `client_signstr` varchar(200) NOT NULL COMMENT '客户签名',
  `vote_time` bigint(20) NOT NULL COMMENT '投票时间',
  `time_create` bigint(20) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=75 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
*/
//10128add
type ClientVoteRsmRecord struct {
	Id         int       `xorm:"not null pk autoincr comment('ID') INT(11)"`
	Clientid   string    `xorm:"unique(clientid) CHAR(120)"`
	RsmGroupId     int       `xorm:"INT(15)"`
	RsmId   string    `xorm:"comment('币种类型') CHAR(40)"`
	VertifyResult     bool       `xorm:"INT(4)"`
	//Address    string    `xorm:"comment('账户公钥hash地址') unique(address) VARCHAR(200)"`
	//Status     int       `xorm:"INT(11)"`
	ClientSignstr     string    `xorm:"comment('client签名') VARCHAR(200)"`
	VoteTime   int64  `json:"vote_time" xorm:"BIGINT(20)"`
	TimeCreate   int64  `json:"created_time" xorm:"BIGINT(20)"`

}

type WdcTranRecord struct {
	Settleid       int64    `json:"orderid" `
	Txhash        string `json:"txhash" xorm:"default ''"`
	From        string `json:"from" xorm:"default ''"`
	To		 	string `json:"to" xorm:"default ''"`
	Amount  	float64 `json:"amount"`
	Amountfee  	float64 `json:"amountfee"`
	Coincode string    `json:"cointype" xorm:"default ''"` //交易币种类
	Status        string `json:"status" xorm:"default ''"`

	Verifystatus int    `json:"verifystatus" xorm:"default 0 index INT(4)"` //交易审核状态
	Errcode       int64    `json:"errorid" `
	Desc	       string `json:"desc" xorm:"default ''"`
	//改为varchar类型在mysql里；
	TimeCreate	   string  `json:"time_create" `
	TimeUpdate	   string  `json:"time_update" `
	Raw   string  `json:"raw" `
}