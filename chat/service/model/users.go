package model

import "time"

type Users struct {
	ID                  uint64    `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint unsigned;not null"`
	Nickname            string    `gorm:"index:users_nickname_index;column:nickname;type:varchar(255);not null;default:''"`
	Mobile              string    `gorm:"unique;column:mobile;type:varchar(255);default:null"`
	Avatar              string    `gorm:"column:avatar;type:varchar(255);not null;default:''"`
	Password            string    `gorm:"column:password;type:varchar(255);default:null"`
	Birthday            string    `gorm:"column:birthday;type:varchar(255);default:null;comment:'生日'"`                                     // 生日
	Gender              string    `gorm:"column:gender;type:enum('male','female','unknown');not null;default:unknown;comment:'性别'"`        // 性别
	Signature           string    `gorm:"column:signature;type:varchar(255);not null;default:'';comment:'个性签名'"`                           // 个性签名
	Locked              bool      `gorm:"index:users_locked_index;column:locked;type:tinyint(1);not null;default:0;comment:'账号是否被锁定'"`     // 账号是否被锁定
	RestrictLoginTime   time.Time `gorm:"column:restrict_login_time;type:timestamp;default:null;comment:'限制登录时间至'"`                        // 限制登录时间至
	IsAnchor            bool      `gorm:"index:users_is_anchor_index;column:is_anchor;type:tinyint(1);not null;default:0;comment:'是否为主播'"` // 是否为主播
	IsVip               bool      `gorm:"index:users_is_vip_index;column:is_vip;type:tinyint(1);not null;default:0;comment:'是否VIP用户'"`     // 是否VIP用户
	VipLevel            bool      `gorm:"column:vip_level;type:tinyint(1);not null;default:0;comment:'vip等级'"`                             // vip等级
	VipExpiration       time.Time `gorm:"column:vip_expiration;type:timestamp;default:null"`
	Accumulated         int64     `gorm:"column:accumulated;type:bigint;default:null;default:0;comment:'积累价值'"`                             // 积累价值
	Gold                uint64    `gorm:"column:gold;type:bigint unsigned;not null;default:0;comment:'金币:主要货币'"`                            // 金币:主要货币
	Silver              uint64    `gorm:"column:silver;type:bigint unsigned;not null;default:0;comment:'银币:附属货币,用于其他免费领取项'"`                // 银币:附属货币,用于其他免费领取项
	GoldEarnings        uint64    `gorm:"column:gold_earnings;type:bigint unsigned;not null;default:0;comment:'金币当前收益'"`                    // 金币当前收益
	GoldEarningsTotal   int64     `gorm:"column:gold_earnings_total;type:bigint;not null;default:0;comment:'金币总收益'"`                        // 金币总收益
	SilverEarnings      uint64    `gorm:"column:silver_earnings;type:bigint unsigned;not null;default:0;comment:'银币当前收益'"`                  // 银币当前收益
	SilverEarningsTotal uint64    `gorm:"column:silver_earnings_total;type:bigint unsigned;not null;default:0;comment:'银币总收益'"`             // 银币总收益
	GoldConsume         uint64    `gorm:"column:gold_consume;type:bigint unsigned;not null;default:0;comment:'金币总消费'"`                      // 金币总消费
	FollowNum           int       `gorm:"column:follow_num;type:int;not null;default:0;comment:'关注的人数'"`                                    // 关注的人数
	FansNum             int       `gorm:"column:fans_num;type:int;not null;default:0;comment:'粉丝人数'"`                                       // 粉丝人数
	Level               int       `gorm:"column:level;type:int;not null;default:1;comment:'全站等级'"`                                          // 全站等级
	Exp                 int64     `gorm:"column:exp;type:bigint;not null;default:0;comment:'全站经验'"`                                         // 全站经验
	LevelRank           int       `gorm:"column:level_rank;type:int;not null;default:0;comment:'全站等级排名'"`                                   // 全站等级排名
	InviteID            int       `gorm:"column:invite_id;type:int;not null;default:0;comment:'邀请我的用户ID'"`                                  // 邀请我的用户ID
	InviteType          int8      `gorm:"column:invite_type;type:tinyint;not null;default:0;comment:'邀请类型'"`                                // 邀请类型
	Province            string    `gorm:"column:province;type:varchar(255);not null;default:'';comment:'省份'"`                               // 省份
	City                string    `gorm:"column:city;type:varchar(255);not null;default:'';comment:'城市'"`                                   // 城市
	Address             string    `gorm:"column:address;type:varchar(255);not null;default:'';comment:'地址'"`                                // 地址
	Occupation          string    `gorm:"column:occupation;type:varchar(255);not null;default:'';comment:'职业'"`                             // 职业
	IsLiving            bool      `gorm:"index:users_is_living_index;column:is_living;type:tinyint(1);not null;default:0;comment:'是否正在直播'"` // 是否正在直播
	CreatedAt           time.Time `gorm:"column:created_at;type:timestamp;default:null"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:timestamp;default:null"`
	FamilyID            int       `gorm:"index:users_family_id_index;column:family_id;type:int;not null;default:0;comment:'家族ID'"` // 家族ID
	FamilyRole          int8      `gorm:"column:family_role;type:tinyint;not null;default:0;comment:'家族身份 0 普通 1 族长 2副族长'"`        // 家族身份 0 普通 1 族长 2副族长
	FamilyJoinedAt      time.Time `gorm:"column:family_joined_at;type:timestamp;default:null;comment:'加入家族时间'"`                    // 加入家族时间
	FamilyContribution  int       `gorm:"column:family_contribution;type:int;not null;default:0;comment:'家族贡献值'"`                  // 家族贡献值
	VehicleID           int       `gorm:"column:vehicle_id;type:int;not null;default:0;comment:'装备的坐骑ID'"`                         // 装备的坐骑ID
	NiceNumber          int       `gorm:"unique;column:nice_number;type:int;default:null;comment:'靓号'"`                            // 靓号
	LastAliveTime       time.Time `gorm:"column:last_alive_time;type:timestamp;default:null;comment:'最后活动时间'"`                     // 最后活动时间
	LastAlivePlatform   int8      `gorm:"column:last_alive_platform;type:tinyint;not null;default:0;comment:'0安卓 1 ios 2web'"`     // 0安卓 1 ios 2web
	LastLoginTime       time.Time `gorm:"column:last_login_time;type:timestamp;default:null;comment:'最后登陆时间'"`                     // 最后登陆时间
	LastLoginIP         string    `gorm:"column:last_login_ip;type:varchar(255);default:null;comment:'最后登陆IP'"`                    // 最后登陆IP
	LastLoginPlatform   int8      `gorm:"column:last_login_platform;type:tinyint;not null;default:0;comment:'0安卓 1 ios 2web'"`     // 0安卓 1 ios 2web
	LastLoginAppToken   string    `gorm:"column:last_login_app_token;type:varchar(500);default:null;comment:'APP端登陆token'"`        // APP端登陆token
	LastLoginWebToken   string    `gorm:"column:last_login_web_token;type:varchar(500);default:null;comment:'WEB端登陆token'"`        // WEB端登陆token
}

type UserModel struct {
}

func (Users) TableName() string {
	return "users"
}

func (a *UserModel) Get(id uint64) (*Users, error) {
	res := &Users{}
	db := mysqlHandles["live"].Handle().Where("id=?", id).First(res)
	if db.Error != nil {
		return nil, db.Error
	}

	return res, nil
}
