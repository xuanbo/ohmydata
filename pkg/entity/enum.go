package entity

// ParamLocation 请求位置
type ParamLocation uint8

const (
	// ParamPath 路径参数
	ParamPath ParamLocation = iota
	// ParamQuery query参数
	ParamQuery
	// ParamBody body参数
	ParamBody
)

// ParamType 参数类型
type ParamType uint

const (
	// Boolean 布尔
	Boolean ParamType = iota
	// Int 整形
	Int
	// Long 长整形
	Long
	// Float 单精度浮点型
	Float
	// Double 双精度浮点型
	Double
	// DateTime 时间
	DateTime
	// String 字符串
	String
	// Object 对象
	Object
	// Array 数组
	Array
)

// ConvertType 转换方式
type ConvertType uint8

const (
	// ConvertNone 无
	ConvertNone ConvertType = iota
	// ConvertRename 重命名
	ConvertRename
)
