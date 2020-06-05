package config

// Error is the errors list
var Error []someError = []someError{
	someError{
		Error: "Can not load JSON data",
		Code:  0,
	},
	someError{
		Error: "Sing Confirmation Failure!",
		Code:  1,
	},
	someError{
		Error: "User with that vk_user_id doesn't exist",
		Code:  2,
	},
	someError{
		Error: "Wrong code",
		Code:  3,
	},
	someError{
		Error: "Are you a type of cool hacker?? :)",
		Code:  4,
	},
	someError{
		Error: "Not Enough money!",
		Code:  5,
	},
}
