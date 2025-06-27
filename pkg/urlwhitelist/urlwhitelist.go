package urlwhitelist

var Business = map[string]int{
	"/pb.BusinessExt/SignIn":                 0,
	"/pb.BusinessExt/TwitterSignIn":          1,
	"/pb.BusinessExt/GetTwitterAuthorizeURL": 2,
	"/pb.BusinessExt/WalletSignIn":           3,
	"/pb.BusinessExt/SearchTwitterUser":      4,

	"/pb.BusinessExt/CreateToken": 5,
	"/pb.BusinessExt/GetToken":    6,
}

var Logic = map[string]int{
	"/pb.LogicExt/RegisterDevice":          0,
	"/pb.LogicExt/CreateChatRoom":          1,
	"/pb.LogicExt/GetChatRoom":             2,
	"/pb.LogicExt/GetChatRooms":            3,
	"/pb.LogicExt/JoinChatRoom":            4,
	"/pb.LogicExt/LeaveChatRoom":           5,
	"/pb.LogicExt/GetChatRoomMembers":      6,
	"/pb.LogicExt/SendChatRoomMessage":     7,
	"/pb.LogicExt/GetUserChatRooms":        8,
	"/pb.LogicExt/GetUserCreatedChatRooms": 9,
	"/pb.LogicExt/GetChatRoomMessages":     10,
}
