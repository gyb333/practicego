package server

import (
	"AdvancedResearch/sibo/protocol"
	"AdvancedResearch/sibo/share"
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"path"

	"runtime"
	"strings"
	"sync"
	"time"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

var ErrServerClosed = errors.New("sibo: Server closed")

var shutdownPollInterval = 500 * time.Millisecond

var schuduleInterval = 1 * time.Minute

// 读写缓冲区大小
const (
	ReaderBufferSize = 1024
	WriterBufferSize = 1024
)

type Server struct {
	ln                   net.Listener
	readTimeout          time.Duration
	writeTimeout         time.Duration

	mu            sync.RWMutex
	activeSession map[ISession]struct{} // 全局连接管理器
	doneChan      chan struct{} //
	inShutdown    int32 // 监听server的关闭状态

	// TLSConfig for creating tls tcp connection.
	tlsConfig *tls.Config // tls配置
	waitGroup *sync.WaitGroup

	cron *cron.Cron // 定时任务
	scheduleSaveDuration time.Duration // 内存数据保存到数据库的定时任务的周期
}

func NewServer() *Server {
	return &Server{
		waitGroup:            &sync.WaitGroup{},
		cron:                 cron.New(),
		scheduleSaveDuration: schuduleInterval,
	}
}
// 返回服务器地址
func (s *Server) Address() net.Addr {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ln == nil {
		return nil
	}
	return s.ln.Addr()
}

// 设置定时任务周期
func (s *Server) SetScheduleSaveDuration(t time.Duration) {
	s.scheduleSaveDuration = t
}

func (s *Server) Serve(network, address string) (err error) {
	var ln net.Listener
	if s.tlsConfig == nil {
		ln, err = net.Listen("tcp", address)
	} else {
		ln, err = tls.Listen("tcp", address, s.tlsConfig)
	}
	s.cron.Start()
	s.scheduleSavePlayer2DB() // start save job
	return s.serveListener(ln)
}

func (s *Server) serveListener(ln net.Listener) error {
	var tempDelay time.Duration

	s.mu.Lock()
	s.ln = ln
	if s.activeSession == nil {
		s.activeSession = make(map[ISession]struct{})
	}
	s.mu.Unlock()

	// 死循环监听来自client的连接
	for {
		conn, e := ln.Accept()
		if e != nil {
			select {
			case <-s.doneChan:
				return ErrServerClosed
			default:
			}
			if ne, ok := e.(net.Error); ok && ne.Temporary() { // temporary error, retry in one second
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("sibo: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			log.Println("closed listener")
			return e
		}
		tempDelay = 0

		if tc, ok := conn.(*net.TCPConn); ok {
			tc.SetKeepAlive(true)
			tc.SetKeepAlivePeriod(3 * time.Minute)
		}

		s.mu.Lock()
		session := NewSession(conn)
		s.activeSession[session] = struct{}{}
		s.mu.Unlock()
		go s.serveSession(session) // 每个连接对应一个goroutine
	}
}

func (s *Server) serveSession(session ISession) {
	log.Println("serve session -> ", session.SessionId())
	player := NewPlayer(session)
	// 客户端与服务器断开连接时执行
	defer func() { // execute when client disconnect from server
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			log.Printf("serving %s panic error: %s, stack:\n %s", session.Conn().RemoteAddr(), err, buf)
		}
		s.mu.Lock()
		delete(s.activeSession, session)
		s.mu.Unlock()
		session.Close(func() {
			log.Println("exec on session close.")
			/*player.Logout(func() {
			    log.Println("exec on player logout.")
			})*/
		})
	}()

	if tlsConn, ok := session.Conn().(*tls.Conn); ok {
		if d := s.readTimeout; d != 0 {
			session.Conn().SetReadDeadline(time.Now().Add(d))
		}
		if d := s.writeTimeout; d != 0 {
			session.Conn().SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.Handshake(); err != nil {
			log.Printf("sibo: TLS handshake error from %s: %v", session.Conn().RemoteAddr(), err)
			return
		}
	}
	r := bufio.NewReaderSize(session.Conn(), ReaderBufferSize)

	// client与server建立连接后，循环接受来自客户端的请求数据
	for {
		t0 := time.Now()
		if s.readTimeout != 0 {
			session.Conn().SetReadDeadline(t0.Add(s.readTimeout))
		}
		req, err := session.ReadRequest(r)
		if err != nil {
			if err == io.EOF { // client close conn
				//player.SaveAll()
				log.Printf("client has closed this connection: %s", session.Conn().RemoteAddr().String())
			} else if strings.Contains(err.Error(), "use of closed network connection") { // (服务器主动关闭conn)conn已被conn.Close()关闭,从conn读取会报这个错误
				log.Printf("sibo: connection %s is closed", session.Conn().RemoteAddr().String())
			} else {
				log.Printf("sibo: failed to read request: %v", err)
			}
			player.Logout(func() {
				log.Println("player logout on session close")
			})
			/*session.Close(func() {
			    player.SaveAll()
			    log.Println("session close : " + session.SessionId())
			})*/
			return
		}
		if s.writeTimeout != 0 {
			session.Conn().SetWriteDeadline(t0.Add(s.writeTimeout))
		}

		// 业务逻辑用单独的goroutine处理以防请求被阻塞(此处后期可使用类似fasthttp的goroutine池来优化)
		go func() {
			if req.IsHeartbeat() {
				req.SetMessageType(protocol.Response)
				session.SendResponse(req)
				return
			}

			res, err := s.handleRequest(player, req)
			if err != nil {
				log.Printf("sibo: failed to handle request: %v", err)
			}
			if !req.IsOneway() {
				session.SendResponse(res)
			}
			protocol.FreeMsg(req)
			protocol.FreeMsg(res)
		}()
	}
}
// 请求处理（业务逻辑）
func (s *Server) handleRequest(player IPlayer, req *protocol.Message) (res *protocol.Message, err error) {
	res = req.Clone()
	res.SetMessageType(protocol.Response)

	codec, ok := share.Codecs[req.SerializeType()]
	if !ok {
		err = fmt.Errorf("can not find codec for %d", req.SerializeType())
		handleError(res, err)
		return res, err
	}
	compression, ok := share.Compression[req.CompressType()]
	if !ok {
		err = fmt.Errorf("can not find compress type for %d", req.CompressType())
		handleError(res, err)
		return res, err
	}

	payload, err := compression.Decompress(req.Payload)
	if err != nil {
		err = fmt.Errorf("UnZip error for compress type %d", req.CompressType())
		handleError(res, err)
		return res, err
	}

	mmId := req.ModuleMessageID()
	if _, ok := proto.RequestMap[mmId]; !ok {
		handleError(res, errors.New("messageID not exist"))
	}
	request := proto.RequestMap[mmId]()
	err = codec.Decode(payload, request)
	if err != nil {
		return handleError(res, err)
	}
	responsePayload, err := ProcessorMap[mmId].Process(player, request)

	if !req.IsOneway() {
		data, err := codec.Encode(responsePayload)
		if err != nil {
			handleError(res, err)
			return res, err
		}
		payload, err := compression.Compress(data)
		if err != nil {
			handleError(res, err)
			return res, err
		}
		res.Payload = payload
	}
	return res, nil
}
// 关闭服务器
func (s *Server) Close() error {
	s.mu.Lock()
	s.closeDoneChan()
	var err error
	if s.ln != nil {
		err = s.ln.Close()
	}
	s.mu.Unlock()

	s.cron.Stop() // stop schedule job
	s.saveAllPlayer2DB() // 保存数据
	s.closeSessions() //关闭所有连接
	s.waitGroup.Wait()
	return err
}
// 关闭所有已建立的连接
func (s *Server) closeSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for session := range s.activeSession {
		session.Close(func() {
			log.Println("session close on by server: " + session.SessionId())
		})
		delete(s.activeSession, session)
	}
}

// 定时任务存储玩家数据
func (s *Server) scheduleSavePlayer2DB() {
	//spec := "0 */2 * * * ?"
	s.cron.Schedule(cron.Every(s.scheduleSaveDuration), new(SavePlayerJob))
}

// 保存玩家数据到数据库，shutdown graceful
func (s *Server) saveAllPlayer2DB() {
	if PlayerId2PlayerMap.Len() > 0 {
		for _, player := range PlayerId2PlayerMap.Values() {
			//player, ok := PlayerId2PlayerMap.Get(k)
			s.waitGroup.Add(1)
			go func(player IPlayer) {  // 此处每个玩家使用一个goroutine，玩家数多可能有问题，可修改为一个goroutine处理多个玩家
				defer s.waitGroup.Done()
				log.Println("task start")
				//time.Sleep(2 * time.Second) 模仿数据保存
				player.SaveAll()
				log.Println("task end")
			}(player)
		}
		PlayerId2PlayerMap.Clear()
	}
}
// 这里使用了一点代码的trick将chan当作sync.Cond使用
func (s *Server) closeDoneChan() {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	select {
	case <-s.doneChan:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(s.doneChan)
	}
}
// 日志配置
func (s *Server) ConfigLocalFilesystemLogger(logPath, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", err)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	})
	log.AddHook(lfHook)
}

func handleError(res *protocol.Message, err error) (*protocol.Message, error) {
	res.SetMessageStatusType(protocol.Error)
	return res, err
}