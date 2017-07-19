package notifications

import (
	"github.com/chanxuehong/wechat.v2/mp/core"
	"github.com/chanxuehong/wechat.v2/mp/user"
)

// Follow 关注通知
/*
	新追XXX
*/
func Follow(OpenID string) (info *user.UserInfo, err error) {

	ats := core.NewDefaultAccessTokenServer(wxAppID, wxAppSecret, nil)
	clt := core.NewClient(ats, nil)
	// user,err := user.Get(clt,OpenID,"zh_CN")
	return user.Get(clt, OpenID, "zh_CN")
}
