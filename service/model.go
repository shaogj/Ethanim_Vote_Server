package service

import (
	"Ethanim_Vote_Server/config"
	"Ethanim_Vote_Server/models"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//mmysql "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/mkideal/log"
	"time"
)

var GXormMysql *xorm.Engine

func InitMysqlDB(conf config.MySqlConfig) error {
	//mmysql.Config{}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
	conf.User,conf.Password,conf.Host,conf.Port,conf.Dbname)
	log.Info("Mysql{%s}",dataSourceName)
	var err error
	GXormMysql,err= xorm.NewEngine("mysql", dataSourceName)
	if nil!=err {
		log.Error("InitMysqlDB() exec err!,errinfo is:%v",err)
		//os.Exit(0)
		return err
	}
	log.Info("InitMysqlDB() exec succ!,conf.Host is:%v",conf.Host)
	return nil
}

//请求API: 创建新的数字币(比特币,莱特币等)地址
func GenerateAccount(curengine *xorm.Engine,cointype string,privkey string,pubkey string,pubkeyaddr string) error {
	enginewrite := curengine
	curGjcAccountKeyTb := new(models.GjcAccountKeyTb)
	//2018.080301==sgj update:
	//accountsql := "select * from gjc_account_key_tb"

	curGjcAccountKeyTb.PrivKey = privkey
	//1019,                  pubkey comparess value is err
	//curGjcAccountKeyTb.PubKey = pubkey
	curGjcAccountKeyTb.AddressId = pubkeyaddr
	curGjcAccountKeyTb.CreatedTime = time.Now().Unix()
	curGjcAccountKeyTb.CoinType = cointype
	rows, err := enginewrite.Table("gjc_account_key_tb").Insert(curGjcAccountKeyTb)

	if err != nil {
		log.Error("GenerateAccount(),Insert row is :%v,rowsnum is:%d,err is-:%v \n", curGjcAccountKeyTb,rows,err)
		return err

	}
	log.Info("GenerateAccount(),Insert row success!,rec is :%v,rowsnum is:%d \n", curGjcAccountKeyTb,rows)
	return nil
}

//1028,for wdc
func  ClientVoteRecordSave(curengine *xorm.Engine,clientId string,rmsid string,rsmgroupid int,clientsignstr string,VertifyResult bool ) error {
	curClientVote := &models.ClientVoteRsmRecord{
		Clientid: clientId,
		RsmId:rmsid,
		RsmGroupId:rsmgroupid,
		ClientSignstr:clientsignstr,
		VertifyResult:VertifyResult,
		VoteTime:time.Now().Unix(),
		TimeCreate:time.Now().Unix(),
	}
	////time.Now().Unix()
	rows, err := curengine.Table(models.TableClientVoteRecord).Insert(curClientVote)
	if err != nil {
		log.Error("ClientVoteRecordSave(),Insert row is :%v,rowsnum is:%d,err is-:%v \n", curClientVote,rows,err)
		return err

	}
	log.Info("ClientVoteRecordSave(),Insert row success!,rec is :%v,rowsnum is:%d \n", curClientVote,rows)
	return nil
}

