/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/18 10:04 上午
 */

package db_config

import (
//_ "github.com/mattn/go-sqlite3"
)

func InitSqlite() {

	//"gorm.io/driver/sqlite" 和postgres冲突
	//db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	//var err error
	//global.GTaskDb, err = sql.Open("sqlite3", "./db/"+global.GConfig.System.ProCode+"_DB.db?cache=shared&mode=memory")
	//
	//if err != nil {
	//	global.GLog.Error("连接错误，", zap.Any("err", err))
	//	os.Exit(0)
	//}
	//global.GLog.Info("sqlite is connecting")

	//生成内存库的表
	//_, _ = global.GTaskDb.Exec(sqlBillStr)
	//_, _ = global.GTaskDb.Exec(sqlBlockStr)
	//_, _ = global.GTaskDb.Exec(sqlFieldStr)
}

var sqlBillStr = `create table "project_bills" (
    created_at     DATETIME,
    updated_at     DATETIME,
    wrong_note     varchar(255),
    id             varchar(20)       not null,
    bill_name      varchar(32)       not null,
    bill_num       varchar(32)       not null,
    del_remarks    varchar,
    pro_code       varchar(10)       not null,
    stage          integer default 1 not null,
    download_path  varchar           not null,
    download_at    DATETIME,
    batch_num      varchar,
    agency         varchar,
    scan_at        DATETIME,
    export_at      DATETIME,
    upload_at      DATETIME,
    status         integer default 1,
    sale_channel   varchar,
    insurance_type integer,
    claim_type     integer,
    count_money    numeric,
    invoice_num    integer,
    question_num   integer,
    quality_user   varchar,
    is_auto_upload boolean,
    stick_level    integer,
    pre_status     integer,
    edit_version   integer default 0,
	CONSTRAINT "project_bill_pkey" PRIMARY KEY ("id")
);
`
var sqlBlockStr = `CREATE TABLE "project_blocks" (
  "id" VARCHAR(20) NOT NULL,
  "created_at" DATETIME,
  "updated_at" DATETIME,
  "bill_id" VARCHAR(20),
  "name" VARCHAR,
  "code" VARCHAR,
  "f_eight" VARCHAR(5),
  "ocr" VARCHAR,
  "free_time" integer(32),
  "is_loop" VARCHAR(5),
  "pic_page" integer(32),
  "is_mobile" VARCHAR(5),
  "w_coordinate" VARCHAR(100),
  "m_coordinate" VARCHAR(100),
  "m_pic_page" integer(32),
  "pre_b_code" VARCHAR,
  "link_b_code" VARCHAR,
  "op1_pre_b_code" VARCHAR,
  "stage" VARCHAR,
  "status" integer(32),
  "zero" integer(32),
  "op1_code" VARCHAR,
  "op1_apply_at" DATETIME,
  "op1_submit_at" DATETIME,
  "op2_code" VARCHAR,
  "op2_apply_at" DATETIME,
  "op2_submit_at" DATETIME,
  "op_q_code" VARCHAR,
  "op_q_apply_at" DATETIME,
  "op_q_submit_at" DATETIME,
  "op_d_code" VARCHAR,
  "op_d_apply_at" DATETIME,
  "op_d_submit_at" DATETIME,
  CONSTRAINT "project_blocks_pkey" PRIMARY KEY ("id")
);
`
var sqlFieldStr = `CREATE TABLE "project_fields" (
  "id" varchar(20) NOT NULL,
  "created_at" DATETIME,
  "updated_at" DATETIME,
  "name" varchar,
  "code" varchar,
  "bill_id" varchar(20),
  "block_id" varchar(20),
  "block_index" integer(32),
  "field_index" integer(32),
  "op1_value" varchar,
  "op1_input" varchar,
  "op2_value" varchar,
  "op2_input" varchar,
  "op_q_value" varchar,
  "op_q_input" varchar,
  "op_d_value" varchar,
  "op_d_input" varchar,
  "result_value" varchar,
  "result_input" varchar,
  CONSTRAINT "project_fields_pkey" PRIMARY KEY ("id")
);
`
