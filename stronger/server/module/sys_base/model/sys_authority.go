package model

import (
	"time"
)

type SysAuthority struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AuthorityId   string `json:"authorityId" gorm:"not null;unique;primary_key" gorm:"comment:'角色ID'"`
	AuthorityName string `json:"authorityName" gorm:"comment:'角色名'"`
	ParentId      string `json:"parentId" gorm:"comment:'父角色ID'"`
	AuthorityDes  string `json:"authorityDes" gorm:"comment:'角色描述'"`
	//DataAuthorityId []SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;association_jointable_foreignkey:data_authority_id"`
	//Children        []SysAuthority `json:"children"`
	//SysBaseMenus    []SysBaseMenu  `json:"menus" gorm:"many2many:sys_authority_menus;"`
}

type AuthorityManager struct {
	Id             int    `json:"id"`
	CreateAt       string `json:"create_at"`
	AuthorityName  string `json:"authority_name"`
	AuthorityState *bool  `json:"authority_state"`
	UserCount      int    `json:"user_count"`
	Describe       string `json:"describe"`
	Id00R          *bool  `json:"id_00_r"`
	Id00w          *bool  `json:"id_00_w"`
	Id010R         *bool  `json:"id_010_r"`
	Id010w         *bool  `json:"id_010_w"`
	Id011R         *bool  `json:"id_011_r"`
	Id011w         *bool  `json:"id_011_w"`
	Id02R          *bool  `json:"id_02_r"`
	Id02w          *bool  `json:"id_02_w"`
	Id03R          *bool  `json:"id_03_r"`
	Id03w          *bool  `json:"id_03_w"`
	Id04R          *bool  `json:"id_04_r"`
	Id04w          *bool  `json:"id_04_w"`
	Id05R          *bool  `json:"id_05_r"`
	Id05w          *bool  `json:"id_05_w"`
	Id06R          *bool  `json:"id_06_r"`
	Id06w          *bool  `json:"id_06_w"`
	Id07R          *bool  `json:"id_07_r"`
	Id07w          *bool  `json:"id_07_w"`
	Id10R          *bool  `json:"id_10_r"`
	Id10w          *bool  `json:"id_10_w"`
	Id110R         *bool  `json:"id_110_r"`
	Id110w         *bool  `json:"id_110_w"`
	Id111R         *bool  `json:"id_111_r"`
	Id111w         *bool  `json:"id_111_w"`
	Id112R         *bool  `json:"id_112_r"`
	Id112w         *bool  `json:"id_112_w"`
	Id113R         *bool  `json:"id_113_r"`
	Id113w         *bool  `json:"id_113_w"`
	Id114R         *bool  `json:"id_114_r"`
	Id114w         *bool  `json:"id_114_w"`
	Id115R         *bool  `json:"id_115_r"`
	Id115w         *bool  `json:"id_115_w"`
	Id116R         *bool  `json:"id_116_r"`
	Id116w         *bool  `json:"id_116_w"`
	Id117R         *bool  `json:"id_117_r"`
	Id117w         *bool  `json:"id_117_w"`
	Id120R         *bool  `json:"id_120_r"`
	Id120w         *bool  `json:"id_120_w"`
	Id121R         *bool  `json:"id_121_r"`
	Id121w         *bool  `json:"id_121_w"`
	Id122R         *bool  `json:"id_122_r"`
	Id122w         *bool  `json:"id_122_w"`
	Id123R         *bool  `json:"id_123_r"`
	Id123w         *bool  `json:"id_123_w"`
	Id124R         *bool  `json:"id_124_r"`
	Id124w         *bool  `json:"id_124_w"`
	Id130R         *bool  `json:"id_130_r"`
	Id130w         *bool  `json:"id_130_w"`
	Id140R         *bool  `json:"id_140_r"`
	Id140w         *bool  `json:"id_140_w"`
	Id141R         *bool  `json:"id_141_r"`
	Id141w         *bool  `json:"id_141_w"`
	Id142R         *bool  `json:"id_142_r"`
	Id142w         *bool  `json:"id_142_w"`
	Id143R         *bool  `json:"id_143_r"`
	Id143w         *bool  `json:"id_143_w"`
	Id144R         *bool  `json:"id_144_r"`
	Id144w         *bool  `json:"id_144_w"`
	Id150R         *bool  `json:"id_150_r"`
	Id150w         *bool  `json:"id_150_w"`
	Id151R         *bool  `json:"id_151_r"`
	Id151w         *bool  `json:"id_151_w"`
	Id160R         *bool  `json:"id_160_r"`
	Id160w         *bool  `json:"id_160_w"`
	Id161R         *bool  `json:"id_161_r"`
	Id161w         *bool  `json:"id_161_w"`
	Id162R         *bool  `json:"id_162_r"`
	Id162w         *bool  `json:"id_162_w"`
	Id163R         *bool  `json:"id_163_r"`
	Id163w         *bool  `json:"id_163_w"`
	Id164R         *bool  `json:"id_164_r"`
	Id164w         *bool  `json:"id_164_w"`
	Id165R         *bool  `json:"id_165_r"`
	Id165w         *bool  `json:"id_165_w"`
}
