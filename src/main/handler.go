// Copyright (c) 2021 上海骞云信息科技有限公司. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"runtime/debug"
	"time"
)

type MessageHandler interface {
	Encode(msg interface{}) []byte
	Decode(buf []byte) (interface{}, int)
	MessageReceived(connHandler *ConnHandler, msg interface{})
	ConnSuccess(connHandler *ConnHandler)
	ConnError(connHandler *ConnHandler)
}

type ConnHandler struct {
	ReadTime       int64
	WriteTime      int64
	Active         bool
	NextConn       *ConnHandler
	conn           net.Conn
	readBuf        []byte
	messageHandler MessageHandler
}

func (connHandler *ConnHandler) Listen(conn net.Conn, messageHandler interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("run time panic: %v", err)
			debug.PrintStack()
			connHandler.messageHandler.ConnError(connHandler)
		}
	}()
	if conn == nil { // Client与Server建立连接
		return
	}
	connHandler.conn = conn
	connHandler.messageHandler = messageHandler.(MessageHandler)
	connHandler.Active = true
	connHandler.ReadTime = time.Now().Unix()
	connHandler.WriteTime = connHandler.ReadTime
	connHandler.messageHandler.ConnSuccess(connHandler) // messageHandler为LPMessageHandler时,RealServerHandler时
	for {
		buf := make([]byte, 1024*8)
		// 一个数据包大小不能超过2M
		if connHandler.readBuf != nil && len(connHandler.readBuf) > 1024*1024*2 {
			connHandler.conn.Close()
		}
		n, err := connHandler.conn.Read(buf) // Read没有设置超时时间，是阻塞的
		if err != nil || n == 0 {
			log.Infof("Disconnect [%s] to [%s].Error Message:%s", conn.LocalAddr(), conn.RemoteAddr(), err)
			connHandler.Active = false
			connHandler.messageHandler.ConnError(connHandler)
			break
		}
		connHandler.ReadTime = time.Now().Unix()
		if connHandler.readBuf == nil {
			connHandler.readBuf = buf[0:n]
		} else {
			connHandler.readBuf = append(connHandler.readBuf, buf[0:n]...)
		}

		for {
			msg, n := connHandler.messageHandler.Decode(connHandler.readBuf)
			if msg == nil {
				break
			}
			connHandler.messageHandler.MessageReceived(connHandler, msg)
			connHandler.readBuf = connHandler.readBuf[n:]
			if len(connHandler.readBuf) == 0 {
				break
			}
		}

		if len(connHandler.readBuf) > 0 {
			buf := make([]byte, len(connHandler.readBuf))
			copy(buf, connHandler.readBuf)
			connHandler.readBuf = buf
		}
	}
}

func (connHandler *ConnHandler) Write(msg interface{}) {
	if connHandler.messageHandler != nil {
		buf := connHandler.messageHandler.Encode(msg)
		connHandler.WriteTime = time.Now().Unix()
		connHandler.conn.Write(buf)
	}
}
