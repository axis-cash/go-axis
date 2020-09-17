// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// BetanetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var BetanetBootnodes = []string{
	"xnode://fbd124e321d20557010f79e5c6c0ff1c138e339d73a51b7a3564790b5a1cf922c8eebeda3d11a18770ee74c08e4ba5c6a5e1e57c5f1bbc18185306326c92d91a@141.164.51.67:33896",
}


// AlphanetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var AlphanetBootnodes = []string{
}

var DevBootnodes = []string{
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{

}
