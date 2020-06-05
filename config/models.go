package config

// ConfirmData is the confirm connection status struct
type ConfirmData struct {
	Status           string `json:"status"`
	Channel          string `json:"channel"`
	UseVkSignChecker bool   `json:"use_vksign_checker"`
}

// SignCheck is the confirm sign check success struct
type SignCheck struct {
	CheckSignStatus  string `json:"check_sign_status"`
	UseVksignChecker bool   `json:"use_vksign_checker"`
}

// Data is json deserialisation model for data request
type Data struct {
	VkUserID                  string `json:"vk_user_id"`
	VkAPIID                   string `json:"vk_app_id"`
	VkIsAppUser               string `json:"vk_is_app_user"`
	VkAreNotificationsEnabled string `json:"vk_are_notifications_enabled"`
	VkLanguage                string `json:"vk_language"`
	VkRef                     string `json:"vk_ref"`
	VkAccessTokenSettings     string `json:"vk_access_token_settings"`
	VkGroupID                 string `json:"vk_group_id"`
	VkViewerGroupRole         string `json:"vk_viewer_group_role"`
	VkPlatform                string `json:"vk_platform"`
	VkIsFavorite              string `json:"vk_is_favorite"`
	Sign                      string `json:"sign"`
}

// GameData is json, socket sends to the server
type GameData struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Coef   int8   `json:"coefficient"`
}

// Config is config
type Config struct {
	// server settings
	HostPort          string `json:"host_port"`
	BugReportHostPort string `json:"bug_report"`
	MyHostPort        string `json:"my_host_port"`
	// economy settings
	MaxProfit int `json:"max_profit"`
	MinProfit int `json:"min_profit"`
	MinLoss   int `json:"min_loss"`
	MaxLoss   int `json:"max_loss"`
	MaxPrice  int `json:"max_price"`
	MinPrice  int `json:"min_price"`
	// Db Settings
	DbHostPort string `json:"db_host_port"`
	DbName     string `json:"db_name"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	// Vk App Settings
	VkAppID           string `json:"vk_app_id"`
	VkAppSecret       string `json:"vk_app_secret"`
	VkAppServiceToken string `json:"vk_app_service_token"`
	UseVkSignChecker  bool   `json:"use_vk_sign_checker"`
}

type someError struct {
	Error string `json:"error"`
	Code  int8   `json:"code"`
}
