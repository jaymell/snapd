// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2014-2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package snappy

import (
	"github.com/ubuntu-core/snappy/snap"
)

// openSnapBlob opens a snap blob returning both a snap.Info completed
// with sideInfo (if not nil) and a corresponding snap.File.
func openSnapBlob(snapFile string, unsignedOk bool, sideInfo *snap.SideInfo) (*snap.Info, snap.File, error) {
	// TODO: what precautions to take if unsignedOk == false ?

	blobf, err := snap.Open(snapFile)
	if err != nil {
		return nil, nil, err
	}

	info, err := blobf.Info()
	if err != nil {
		return nil, nil, err
	}

	var snapInfo snap.Info
	snapInfo = *info
	if sideInfo != nil {
		snapInfo.SideInfo = *sideInfo
	}

	return &snapInfo, blobf, nil
}
