// Copyright 2017 The go-ethereum Authors
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
//
//

package params

// These are the multipliers for VXET denominations.
// Example: To get the wei value of an amount in 'gwei', use
//
//	new(big.Int).Mul(value, big.NewInt(params.GWei))
const (
	Wei    = 1
	KWei   = 1e3
	MWei   = 1e6
	GWei   = 1e9
	Szabo  = 1e12
	Finney = 1e15
	VXET   = 1e18 // 从Ether改为VXET
	KVXET  = 1e21 // 从KEther改为KVXET
	MVXET  = 1e24 // 从MEther改为MVXET
	GVXET  = 1e27 // 从GEther改为GVXET
)
