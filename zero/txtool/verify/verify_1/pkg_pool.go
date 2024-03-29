package verify_1

import (
	"github.com/axis-cash/go-axis-import/c_superzk"
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/zero/utils"
	"github.com/axis-cash/go-axis/zero/zconfig"
)

var verify_pkg_procs_pool = utils.NewProcsPool(func() int { return zconfig.G_v_thread_num })

type verify_pkg_desc struct {
	asset_cm c_type.Uint256
	proof    c_type.Proof
	e        error
}

func (self *verify_pkg_desc) Run() error {
	if e := c_superzk.VerifyOutput(&self.asset_cm, &self.proof, true); e != nil {
		self.e = e
		return e
	} else {
		return nil
	}
}
