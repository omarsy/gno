// PKGPATH: gno.land/r/demo/boards_test
package boards_test

// SEND: 200000000ugnot

import (
	"std"

	"gno.land/p/demo/testutils"
	"gno.land/r/demo/boards"
	"gno.land/r/demo/users"
)

var bid boards.BoardID

func init() {
	caller := testutils.TestAddress("caller")
	std.TestSetRealm(std.NewUserRealm(caller))
	users.Register("", "gnouser", "my profile")
	bid = boards.CreateBoard("test_board")
	boards.CreateReply(bid, 0, 0, "Reply of the second post")
}

func main() {
	println(boards.Render("test_board"))
}

// Error:
// thread not exist
