package config

var (
	settings  *Settings
)

type Redis struct{
	Size			int
	Network			string
	Address			string
	Password		string
	KeyPairs		string
	SessionDB		int
	AccountLockDB	int
	SessionMaxAge	int
}
type Mysql struct{
	User			string
	Password		string
	Address			string
	DBName			string
	LogMode			bool
	MaxIdleConns	int
	MaxOpenConns	int
}
type WebSvr struct{
	Port			string
}
type DNS struct{
	DnsFilePath			string
	NasDnsFileName		string
	NasDnsBkFileName	string
	DnsConfUpdateTime	int64		// DNS配置更新时间
}

type Common struct{
	HeartTimeRelay		int64		//	relay心跳时间 (默认 30秒)
	HeartTimemcloud		int64		//	mcloud心跳时间	(默认 120秒)
	HTRelayTimeOut		int64		// 	relay超时时间 (默认 45秒)
	HTRelayDelete		int64		// 	relay断开删除时间 (默认 1800秒)
	HTMcloudTimeOut		int64		//	mcloud超时时间 (默认 180秒)
	HTMcloudDelete		int64		//	mcloud断开删除时间 (默认 5400秒)
	RelayForbidSwitch	int			// 	黑名单开关  (默认 1)
	ForbidCountTime		int64		// 	黑名单统计时间 (默认 600秒)
	ForbidFailedRatio	float32		//	统计时间内需要加入黑名单的失败比例 (默认 0.5)
	ForbidTime			int64		// 	封禁时间 (默认 600秒)
	MaxMcloudNum		int			//  relay允许连接的最大mcloud数量 (默认 100)
	MaxRelayCpuPercent	int			//	relay负载标准之一，表示relay 最大cpu使用率(1 - 100) (默认 90)
	MaxRelayMemPercent	int			//  relay负载标准之一， 表示relay最大内存使用率 (1 - 100) (默认 90)
	AllocRelayNum		int			//	分配relay个数，默认3个
	DeleteRalayRecordTime int64		// 	删除 mcloud连接的relay记录的时间
}

type Settings struct {
	Redis			Redis
	Mysql			Mysql
	WebSvr			WebSvr
	DNS				DNS
	Common			Common
}

