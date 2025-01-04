package urlwhitelist

var Business = map[string]int{
	"/pb.BusinessExt/SignIn":                 0,
	"/pb.BusinessExt/TwitterSignIn":          1,
	"/pb.BusinessExt/GetTwitterAuthorizeURL": 2,
	"/pb.BusinessExt/WalletSignIn":           3,
}

var Logic = map[string]int{
	"/pb.LogicExt/RegisterDevice": 0,
}
