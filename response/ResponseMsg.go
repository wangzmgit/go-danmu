package response

const (
	OK                     = "ok"
	ParameterError         = "参数错误"
	RequestError           = "请求错误"
	SystemError            = "服务器出错了"
	OperationTooFrequently = "操作过于频繁"
	PageOrSizeError        = "页码或请求数量错误"
	TooManyRequests        = "请求数量过多"
	SendFail               = "发送失败"
	VerificationFail       = "验证失败"
	CreateFail             = "上传失败"
	ModifyFail             = "修改失败"
	DeleteFail             = "删除失败"
	UploadFail             = "上传失败"
	CheckUpdateFail        = "检查更新失败"
	PullFail               = "获取源码失败"
	UpdateFail             = "更新失败"
	UpdateStatusFail       = "视频状态更新失败"

	NickCheck            = "昵称不能为空"
	NameCheck            = "姓名不能为空"
	TitleCheck           = "标题不能为空"
	CoverCheck           = "请上传封面图片"
	ContentCheck         = "内容不能为空"
	LoginCheck           = "用户名和密码不能为空"
	CommentCheck         = "评论或回复内容不能为空"
	DanmakuCheck         = "不能发送空弹幕"
	GenderCheck          = "性别选择有误"
	SearchCheck          = "搜索的内容不能为空"
	PasswordCheck        = "密码不能少于6位"
	PartitionCheck       = "未选择分区"
	AuthorityCheck       = "权限选择有误"
	EmailFormatCheck     = "邮箱格式有误"
	TelephoneFormatCheck = "联系方式格式有误"
	BirthdayFormatCheck  = "请输入正确的出生日期"

	ReadFail          = "读取失败"
	FileUploadFail    = "文件上传失败"
	FileSaveFail      = "文件保存失败"
	FileCheckFail     = "文件不符合要求"
	FileSizeCheckFail = "文件大小不符合要求"

	UserNotExist              = "用户不存在"
	VideoNotExist             = "视频不存在"
	ImgNotExist               = "图片不存在"
	PartitionNotExist         = "分区不存在"
	CarouselNotExist          = "轮播图不存在"
	MessageNotExist           = "消息不存在"
	ResourceNotExist          = "视频资源不存在"
	CollectionNotExist        = "合集不存在"
	ParentPartitionNotExist   = "所属分区不存在"
	CollectionOrVideoNotExist = "合集或视频不存在"
	SkinNotExist              = "主题不存在或已被删除"
	CommentNotExist           = "评论不存在或已被删除"

	VideoTypeError        = "视频类型错误"
	IsCollect             = "已经收藏过了"
	IsLike                = "已经点过赞了"
	NotCollect            = "还没有收藏"
	NotLike               = "还没有点赞"
	NameOrPasswordError   = "用户名和密码错误"
	CantFollowYourself    = "不能关注自己"
	CantSendYourself      = "不能发送给自己"
	EmailRegistered       = "该邮箱已经被注册了"
	VerificationCodeError = "验证码错误"
	PleaseLoginFirst      = "请先登录"
	SystemNotSupported    = "暂不支持当前系统"
)
