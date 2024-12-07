package consts

const AppName = "go-auth-admin"

const PasswordHashCost int = 10
const ErrExitStatus int = 2

const StatusSuccess = "success"

// App consts
const (
	PhoneNumberMinLength = 2 + 8  // 1+9
	PhoneNumberMaxLength = 4 + 15 // 1+18

	EmailMinLength = 6

	PasswordMinLength = 8
	PasswordMaxLength = 50 // bcrypt 72

	SecretCodeLength      = 8
	LongTextLength        = 32767 //  int(int16(^uint16(0) >> 1)) // equivalent of short.MaxValue
	DefaultTextLength     = 100
	DefaultMapZoom        = 12
	DefaultMaxQty         = 12
	TitleTextLengthTiny   = 12
	TitleTextLengthSmall  = 25
	TitleTextLengthInfo   = 35
	TitleTextLengthMedium = 50
	TitleTextLengthLarge  = 100

	// WF_STATUS_NEW       = 0
	// WF_STATUS_PROGRESS  = 6
	// WF_STATUS_DELETE    = 7
	// WF_STATUS_ERROR     = 10
	// WF_STATUS_SUCCESS   = 15
	// WF_STATUS_VOID      = 17
	// WF_STATUS_SIGNED    = 4
	// WF_STATUS_DELIVERED = 5
	// WF_STATUS_OUTBOX    = 3
	// WF_STATUS_READONLY  = 32
	// WF_STATUS_UNPAID    = 19
	// WF_STATUS_PAID      = 21
	// WF_STATUS_INQUEUE   = 31
)

// const (
// 	LogLevelError = 0
// 	LogLevelWarn  = 1
// 	LogLevelInfo  = 2
// 	LogLevelDebug = 3
// )

const (
	// PathAPI represents the group of PathAPI.
	PathAPI = "/api"
)

const (
	RoleAdmin = "admin"
)
const (
	AuthRolePrefix  = "auth_"
	AuthRoleAccess  = "auth_access"
	AuthRoleAdd     = "auth_add"
	AuthRoleEdit    = "auth_edit"
	AuthRoleView    = "auth_view"
	AuthRoleDelete  = "auth_delete"
	AuthRolePublish = "auth_publish"
)

//nolint:gosec
const (
	PathSysMetricsAPI = "/sys/api/metrics"
)

//nolint:gosec
const (
	PathAuthAdminPingDebugAPI = "/auth-admin/api/ping"

	PathAuthAdmin               = "/auth-admin"
	PathAuthAdminAssets         = "/auth-admin/assets"
	PathAuthAdminAccounts       = "/auth-admin/accounts"
	PathAuthAdminAccountsEntity = "/auth-admin/accounts/:code" // GET
	PathAuthAdminStatusAPI      = "/auth-admin/api/status"     // get _csrf, user related, no-cache
	PathAuthAdminConfigAPI      = "/auth-admin/api/config"     // public

	PathAuthAdminAccountsAPI             = "/auth-admin/api/accounts"            // LIST POST
	PathAuthAdminAccountsEntityAPI       = "/auth-admin/api/accounts/:id"        // GET PUT DELETE
	PathAuthAdminAccountsEntityByCodeAPI = "/auth-admin/api/accounts/:code/code" // GET

	PathAuthAdminAccountsEntityPasswordAPI = "/auth-admin/api/accounts/:id/password" // GET
)
