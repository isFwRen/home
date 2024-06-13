create table project_bills
(
    id                   varchar(32)                 not null
        constraint project_bill_pkey
            primary key,
    bill_name            varchar(32)                 not null,
    bill_num             varchar(32)                 not null,
    pro_code             varchar(10)                 not null,
    stage                integer default 1           not null,
    download_path        varchar                     not null,
    download_at          timestamp(6) with time zone not null,
    batch_num            varchar,
    agency               varchar,
    scan_at              timestamp(6) with time zone,
    export_at            timestamp(6) with time zone,
    upload_at            timestamp with time zone,
    status               integer default 1,
    sale_channel         varchar,
    insurance_type       varchar(32),
    claim_type           integer,
    count_money          numeric,
    invoice_num          integer,
    question_num         integer,
    quality_user_code    varchar,
    is_auto_upload       boolean,
    stick_level          integer,
    pre_status           integer,
    edit_version         integer default 0,
    created_at           timestamp(6) with time zone,
    updated_at           timestamp(6) with time zone,
    del_remarks          varchar,
    wrong_note           varchar(1024000),
    images               character varying[],
    pictures             character varying[],
    template             varchar(255),
    last_upload_at       timestamp(6) with time zone,
    quality_user_name    varchar(255),
    crypto               varchar(255),
    is_timeout           boolean,
    deadline_upload_time timestamp with time zone
);

comment on column project_bills.bill_name is '单据号来源内部';

comment on column project_bills.bill_num is '单据号来源客户';

comment on column project_bills.pro_code is '项目编码';

comment on column project_bills.stage is '录入状态';

comment on column project_bills.download_path is '文件路径';

comment on column project_bills.download_at is '下载时间';

comment on column project_bills.batch_num is '批次号';

comment on column project_bills.agency is '机构号';

comment on column project_bills.scan_at is '扫描时间';

comment on column project_bills.export_at is '导出时间';

comment on column project_bills.upload_at is '回传时间';

comment on column project_bills.status is '案件状态';

comment on column project_bills.sale_channel is '销售渠道';

comment on column project_bills.insurance_type is '医保类型，改为单证类型';

comment on column project_bills.claim_type is '理赔类型';

comment on column project_bills.count_money is '账单金额';

comment on column project_bills.invoice_num is '发票数量';

comment on column project_bills.question_num is '问题件';

comment on column project_bills.quality_user_code is '质检人';

comment on column project_bills.is_auto_upload is '是否自动回传';

comment on column project_bills.stick_level is '加急件';

comment on column project_bills.pre_status is '删除前的状态';

comment on column project_bills.edit_version is '版本';

comment on column project_bills.created_at is '创建时间';

comment on column project_bills.updated_at is '更新时间';

comment on column project_bills.del_remarks is '删除备注';

comment on column project_bills.wrong_note is '导出校验';

comment on column project_bills.template is '模板';

comment on column project_bills.last_upload_at is '最新回传时间';

comment on column project_bills.quality_user_name is '质检人姓名';

comment on column project_bills.crypto is '图片加密秘钥';

comment on column project_bills.is_timeout is '是否超时';

comment on column project_bills.deadline_upload_time is '最晚回传时间';

alter table project_bills
    owner to postgres;

create index project_bills_stage_idx
    on project_bills (stage);

create index project_bills_status_idx
    on project_bills (status);

create table project_blocks
(
    id             varchar(32)           not null
        primary key,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    bill_id        varchar(32),
    name           varchar,
    code           varchar,
    f_eight        boolean,
    ocr            varchar,
    free_time      bigint,
    is_loop        boolean,
    pic_page       bigint,
    is_mobile      boolean,
    w_coordinate   varchar(100)[],
    m_coordinate   varchar(100)[],
    m_pic_page     bigint,
    stage          varchar,
    status         bigint,
    zero           bigint,
    op1_code       varchar,
    op1_apply_at   timestamp with time zone,
    op1_submit_at  timestamp with time zone,
    op2_code       varchar,
    op2_apply_at   timestamp with time zone,
    op2_submit_at  timestamp with time zone,
    opq_code       varchar,
    opq_apply_at   timestamp with time zone,
    opq_submit_at  timestamp with time zone,
    op0_code       varchar,
    op0_apply_at   timestamp with time zone,
    op0_submit_at  timestamp with time zone,
    pre_b_code     varchar(100)[],
    link_b_code    varchar(100)[],
    op1_pre_b_code varchar(100)[],
    is_practice    boolean default false not null,
    picture        varchar(100),
    level          integer,
    op0_stage      varchar(255),
    op1_stage      varchar(255),
    op2_stage      varchar(255),
    opq_stage      varchar(255),
    is_competitive boolean,
    crypto         varchar(255),
    temp           varchar(20)
);

alter table project_blocks
    owner to postgres;

create table project_fields
(
    id            varchar(32) not null
        primary key,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone,
    name          varchar,
    code          varchar,
    bill_id       varchar(32),
    block_id      varchar,
    block_index   bigint,
    field_index   bigint,
    op1_value     varchar,
    op1_input     varchar,
    op2_value     varchar,
    op2_input     varchar,
    opq_value     varchar,
    opq_input     varchar,
    op0_value     varchar,
    op0_input     varchar,
    result_value  varchar,
    result_input  varchar,
    final_value   varchar,
    final_input   varchar,
    right_value   varchar(255),
    feedback_date timestamp with time zone,
    is_practice   boolean default false,
    is_change     boolean
);

comment on column project_fields.final_value is '结果值内容';

comment on column project_fields.final_input is '结果值状态';

comment on column project_fields.right_value is '客户反馈对的值';

comment on column project_fields.feedback_date is '客户反馈日期';

comment on column project_fields.is_practice is '练习';

alter table project_fields
    owner to postgres;

