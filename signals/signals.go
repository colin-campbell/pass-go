/*
 * Copyright (c) 2020. Colin Stewart Campbell <colin.campbell@digitalistgroup.com>
 *  This file is part of Pass-Go.
 *
 *      Pass-Go is free software: you can redistribute it and/or modify
 *      it under the terms of the GNU General Public License as published by
 *      the Free Software Foundation, either version 3 of the License, or
 *      (at your option) any later version.
 *
 *      Pass-Go is distributed in the hope that it will be useful,
 *      but WITHOUT ANY WARRANTY; without even the implied warranty of
 *      MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *      GNU General Public License for more details.
 *
 *      You should have received a copy of the GNU General Public License
 *      along with Pass-Go.  If not, see <https://www.gnu.org/licenses/>.
 */

package signals

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Setup() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		//	syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
			// kill -SIGTERM XXXX
			case syscall.SIGTERM:
			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				fmt.Println(s.String())
			default:
				fmt.Println("Unknown signal.")
			}
			os.Exit(0)
		}
	}()
}
