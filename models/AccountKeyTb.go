
package models

import "time"

var (
	//0614
	TableClientVoteRecord="client_vote_rsm_record"
	TableWalletUserPollsRecord="wallet_user_polls"


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
  `rsm_id` varchar(60) DEFAULT NULL COMMENT 'rsmid',
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
	RsmId   string    `xorm:"comment('币种类型') CHAR(60)"`
	VertifyResult     bool       `xorm:"INT(4)"`
	//Address    string    `xorm:"comment('账户公钥hash地址') unique(address) VARCHAR(200)"`
	//Status     int       `xorm:"INT(11)"`
	ClientSignstr     string    `xorm:"comment('client签名') VARCHAR(200)"`
	VoteTime   int64  `json:"vote_time" xorm:"BIGINT(20)"`
	TimeCreate   int64  `json:"created_time" xorm:"BIGINT(20)"`

}
type WalletUserPolls struct {
	Id         int       		`xorm:"not null pk autoincr comment('ID') INT(11)"`
	MiningGroupId     int   	`xorm:"comment('挖矿分组ID') INT(15)"`
	BlockNums     int      		`xorm:"comment('1个混洗周期所产生的区块数量') INT(15)"`
	AssociatedRsm   string    	`xorm:"unique(associated_rsm) CHAR(60)"`
	FinalResults     int       `xorm:"comment('rsm验证结果:1-可信 0-不可信') INT(4)"`
	Status     int       		`xorm:"comment('0-奖励待发放 1-奖励已发放') TINYINT(3)"`
	MinorityIds     string    `xorm:"comment('验证的人数少的客户端ID') VARCHAR(200)"`
	MajorityIds     string    `xorm:"comment('验证的人数多的客户端ID') VARCHAR(200)"`
	CreateTime     time.Time `xorm:"comment('创建时间') VARCHAR(200)"`
	UpdateTime     time.Time  `xorm:"comment('更新时间') VARCHAR(200)"`
}
//0614add
/*
CREATE TABLE `wallet_user_polls` (
    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
    `mining_group_id` integer COMMENT '挖矿分组ID',
    `block_nums` integer COMMENT '1个混洗周期所产生的区块数量',
    `associated_rsm` integer COMMENT '关联RSM',
    `final_results` integer COMMENT 'rsm验证结果', # 1-可信 0-不可信
    `status` TINYINT(3) NOT NULL DEFAULT 0 COMMENT '状态', # 0-奖励待发放 1-奖励已发放
    `minority_ids` text COMMENT '验证的人数少的客户端ID',
    `majority_ids` text COMMENT '验证的人数多的客户端ID',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `idx_group_associated_rsm_tx` (`mining_group_id`, `associated_rsm`)
) ENGINE=InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8 COMMENT='wallet user polls';
*/
//mining_group_id，associated_rsm，final_results，minority_ids （json（[ids]））, majority_ids
