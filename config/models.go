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
	HostPort string
	Workers  int
	// economy settings
	MaxProfit int
	MinProfit int
	MinLoss int
	MaxLoss int
	MaxPrice  int
	MinPrice  int
	WinCount int
	DefeatCount int
	// Db Settings
	DbHostPort string
	DbName     string
	DbUser     string
	DbPassword string
	// Vk App Settings
	VkAppID           string
	VkAppSecret       string
	VkAppServiceToken string
	UseVkSignChecker  bool
}

type someError struct {
	Error string `json:"error"`
	Code  int8   `json:"code"`
}
