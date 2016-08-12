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

package osutil_test

import (
	"io/ioutil"
	"path/filepath"

	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/osutil"
)

type ReleaseTestSuite struct {
}

var _ = Suite(&ReleaseTestSuite{})

func (s *ReleaseTestSuite) TestSetup(c *C) {
	c.Check(osutil.Series, Equals, "16")
}

func mockOSRelease(c *C) string {
	// FIXME: use AddCleanup here once available so that we
	//        can do osutil.SetLSBReleasePath() here directly
	mockOSRelease := filepath.Join(c.MkDir(), "mock-os-release")
	s := `
NAME="Ubuntu"
VERSION="18.09 (Awesome Artichoke)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="I'm not real!"
VERSION_ID="18.09"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"
UBUNTU_CODENAME=awesome
`
	err := ioutil.WriteFile(mockOSRelease, []byte(s), 0644)
	c.Assert(err, IsNil)

	return mockOSRelease
}

func (s *ReleaseTestSuite) TestReadOSRelease(c *C) {
	reset := osutil.MockOSReleasePath(mockOSRelease(c))
	defer reset()

	os := osutil.ReadOSRelease()
	c.Check(os.ID, Equals, "ubuntu")
	c.Check(os.VersionID, Equals, "18.09")
}

func (s *ReleaseTestSuite) TestReadWonkyOSRelease(c *C) {
	mockOSRelease := filepath.Join(c.MkDir(), "mock-os-release")
	dump := `NAME="elementary OS"
VERSION="0.4 Loki"
ID="elementary OS"
ID_LIKE=ubuntu
PRETTY_NAME="elementary OS Loki"
VERSION_ID="0.4"
HOME_URL="http://elementary.io/"
SUPPORT_URL="http://elementary.io/support/"
BUG_REPORT_URL="https://bugs.launchpad.net/elementary/+filebug"`
	err := ioutil.WriteFile(mockOSRelease, []byte(dump), 0644)
	c.Assert(err, IsNil)

	reset := osutil.MockOSReleasePath(mockOSRelease)
	defer reset()

	os := osutil.ReadOSRelease()
	c.Check(os.ID, Equals, "elementary")
	c.Check(os.VersionID, Equals, "0.4")
}

func (s *ReleaseTestSuite) TestReadOSReleaseNotFound(c *C) {
	reset := osutil.MockOSReleasePath("not-there")
	defer reset()

	os := osutil.ReadOSRelease()
	c.Assert(os, DeepEquals, osutil.OS{ID: "linux", VersionID: "unknown"})
}

func (s *ReleaseTestSuite) TestOnClassic(c *C) {
	reset := osutil.MockOnClassic(true)
	defer reset()
	c.Assert(osutil.OnClassic, Equals, true)

	reset = osutil.MockOnClassic(false)
	defer reset()
	c.Assert(osutil.OnClassic, Equals, false)
}

func (s *ReleaseTestSuite) TestReleaseInfo(c *C) {
	reset := osutil.MockReleaseInfo(&osutil.OS{
		ID: "distro-id",
	})
	defer reset()
	c.Assert(osutil.ReleaseInfo.ID, Equals, "distro-id")
}

func (s *ReleaseTestSuite) TestForceDevMode(c *C) {
	// Restore real OS info at the end of this function.
	defer osutil.MockReleaseInfo(&osutil.OS{})()
	distros := []struct {
		id        string
		idVersion string
		devmode   bool
	}{
		// Please keep this list sorted
		{id: "arch", devmode: true},
		{id: "debian", devmode: true},
		{id: "elementary", devmode: true},
		{id: "elementary", idVersion: "0.4", devmode: false},
		{id: "fedora", devmode: true},
		{id: "gentoo", devmode: true},
		{id: "neon", devmode: false},
		{id: "opensuse", devmode: true},
		{id: "rhel", devmode: true},
		{id: "ubuntu", devmode: false},
	}
	for _, distro := range distros {
		rel := &osutil.OS{ID: distro.id, VersionID: distro.idVersion}
		c.Logf("checking distribution %#v", rel)
		osutil.MockReleaseInfo(rel)
		c.Assert(osutil.ReleaseInfo.ForceDevMode(), Equals, distro.devmode)
	}
}
