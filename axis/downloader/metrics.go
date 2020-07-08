// copyright 2018 The go-ethereum Authors
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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/axis-cash/go-axis/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("axis/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("axis/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("axis/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("axis/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("axis/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("axis/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("axis/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("axis/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("axis/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("axis/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("axis/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("axis/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("axis/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("axis/downloader/states/drop", nil)
)
