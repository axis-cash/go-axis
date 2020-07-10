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
	"xnode://57c7e378a284ec7f196097be6540d99305381f80e6a6a1b951edb9775a856a8c0adc4d68253bc71eb049ee66944b7527da295bcccd47270f720da35ef879bfd6@39.100.244.135:33896",
	"xnode://7bd428aba2b833493e8b7f0d583ee82178581ec944d30a7a73a4ad38721bbe73424e6f96153ed2bc752edd0848a29113c0ad4cc4b4c83c2c121502e8d88b51b7@141.164.51.67:33896",
}

// AlphanetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var AlphanetBootnodes = []string{
	"xnode://29ad1a59fbb8984ec167b890ddb52e8690c6c1699c6796ca209cc2d0a5603760f25747fdae61ecb937260491f51632e6651a9292b2b3ce871e6c2e0c7c0f7c00@118.25.146.113:60609",
}

var DevBootnodes = []string{
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{

}
