// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package stopdaemons

import (
	"fmt"
	"os"
	"time"

	"github.com/EGaaS/go-egaas-mvp/packages/system"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
	"github.com/EGaaS/go-egaas-mvp/packages/utils/sql"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("stop_daemons")

// WaitStopTime closes the database and stop daemons
func WaitStopTime() {
	var first bool
	for {
		if sql.DB == nil || sql.DB.DB == nil {
			time.Sleep(time.Second * 3)
			continue
		}
		if !first {
			err := sql.DB.ExecSQL(`DELETE FROM stop_daemons`)
			if err != nil {
				log.Error(utils.ErrInfo(err).Error())
			}
			first = true
		}
		dExists, err := sql.DB.Single(`SELECT stop_time FROM stop_daemons`).Int64()
		if err != nil {
			log.Error(utils.ErrInfo(err).Error())
		}
		log.Debug("dExtit: %d", dExists)
		if dExists > 0 {
			fmt.Println("Stop_daemons from DB!")
			utils.CancelFunc()
			for i := 0; i < utils.DaemonsCount; i++ {
				name := <-utils.ReturnCh
				log.Debugf("daemon %s stopped", name)
			}
			fmt.Println("Daemons killed")
			err := sql.DB.Close()
			if err != nil {
				log.Error(utils.ErrInfo(err).Error())
			}
			fmt.Println("DB Closed")
			err = os.Remove(*utils.Dir + "/daylight.pid")
			if err != nil {
				log.Error(utils.ErrInfo(err).Error())
				panic(err)
			}
			fmt.Println("removed " + *utils.Dir + "/daylight.pid")
			system.FinishThrust(1)
		}
		time.Sleep(time.Second)
	}
}
