/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/
package engine

import (
	"testing"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/utils"
	"github.com/cgrates/rpcclient"
	"github.com/nyaruka/phonenumbers"
)

func TestDynamicDpFieldAsInterface(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()
	ms := utils.MapStorage{}
	dDp := newDynamicDP([]string{}, []string{utils.ConcatenatedKey(utils.MetaInternal, utils.StatSConnsCfg)}, []string{}, "cgrates.org", ms)
	clientconn := make(chan rpcclient.ClientConnector, 1)
	clientconn <- &ccMock{
		calls: map[string]func(args interface{}, reply interface{}) error{
			utils.StatSv1GetQueueFloatMetrics: func(args, reply interface{}) error {
				rpl := &map[string]float64{
					"stat1": 31,
				}
				*reply.(*map[string]float64) = *rpl
				return nil
			},
		},
	}
	connMgr := NewConnManager(cfg, map[string]chan rpcclient.ClientConnector{
		utils.ConcatenatedKey(utils.MetaInternal, utils.StatSConnsCfg): clientconn,
	})
	SetConnManager(connMgr)
	if _, err := dDp.fieldAsInterface([]string{utils.MetaStats, "val", "val3"}); err == nil || err != utils.ErrNotFound {
		t.Error(err)
	} else if _, err := dDp.fieldAsInterface([]string{utils.MetaLibPhoneNumber, "+402552663", "val3"}); err != nil {
		t.Error(err)
	} else if _, err := dDp.fieldAsInterface([]string{utils.MetaLibPhoneNumber, "+402552663", "val3"}); err != nil {
		t.Error(err)
	} else if _, err := dDp.fieldAsInterface([]string{utils.MetaAsm, "+402552663", "val3"}); err == nil || err != utils.ErrNotFound {
		t.Error(err)
	}

}

func TestDpLibPhoneNumber(t *testing.T) {

	libphonenumber := &libphonenumberDP{
		pNumber: &phonenumbers.PhoneNumber{
			CountryCode: func(i int32) *int32 {

				return &i
			}(33),
			NationalNumber: func(i uint64) *uint64 {

				return &i
			}(121411111),
		},
		cache: utils.MapStorage{},
	}
	if val, err := libphonenumber.fieldAsInterface([]string{"CountryCode"}); err != nil {
		t.Error(err)
	} else if val != *libphonenumber.pNumber.CountryCode {
		t.Errorf("expected %v,received %v", libphonenumber.pNumber.CountryCode, val)
	}

	if val, err := libphonenumber.fieldAsInterface([]string{"NationalNumber"}); err != nil {
		t.Error(err)
	} else if val != *libphonenumber.pNumber.NationalNumber {
		t.Errorf("expected %v,received %v", libphonenumber.pNumber.CountryCode, val)
	}

	if val, err := libphonenumber.fieldAsInterface([]string{"Region"}); err != nil {
		t.Error(err)
	} else if val != "FR" {
		t.Errorf("expected %v,received %v", "FR", val)
	}
	if _, err := libphonenumber.fieldAsInterface([]string{"NumberType"}); err != nil {
		t.Error(err)
	}
}