package main

import (
	"blockchain/controller"
	logger "blockchain/logs"
	"blockchain/routes"
	"blockchain/setting"
	"blockchain/tools"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//@title 简单的区块链项目
//@version 1.0
//@description 第一个练手区块链项目

//@contact.name cyxxxxxx
//@contact.email 2642427375@qq.com

//@host 127.0.0.1
//@BasePath: 8080

func main() {
	//127.0.0.1
	//1.配置信息初始化
	err := setting.Init()
	if err != nil {
		fmt.Println("config load failed err: ", err)
	}
	//2.初始化日志
	if err = logger.Init(setting.Conf.LogConf, setting.Conf.Mode); err != nil {
		fmt.Println("logger load failed err: ", err)
	}
	defer zap.L().Sync() //延迟注册让缓冲区的数据也都写入日志core

	//初始化时区为CST
	tools.InitInMain()

	/* //3.初始化数据库连接
	if err := mysql.Init(setting.Conf.MysqlConf); err != nil {
		zap.L().Error("database connect failed err: ", zap.Error(err))
	} else {
		zap.L().Error("database connect succeeded!!!")
	}
	defer mysql.Close() */

	//4.服务前要

	//注册gin框架内置的校验器使用的翻译器
	if err = controller.InitTrans("zh"); err != nil {
		zap.L().Error("init trans failed, err:%v\n", zap.Error(err))
		return
	}
	controller.GenesisBlock()
	//5.注册路由
	r := routes.Setup(setting.Conf.Mode)
	r.Run()

	//6.启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		//开启一个goroutine启动服务
		if err = srv.ListenAndServe(); err != nil {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	//等待中断信号来优雅关闭服务器 为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) //创建一个接收信号的通道
	//kill 默认会发送 syscall.SIGTERM 信号
	//kill -2 发送 syscall.SIGINT 信号 我们常用的Ctrl+C就是出发系统SIGINT信号
	//kill -9 发送 syscall.SIGKILL 信号 但是不能被捕获 所以不需要添加它
	//signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不会阻塞
	<-quit                                               //这里阻塞 当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ......")
	//创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//5秒内优雅关闭服务(将未处理完的请求处理完再关闭服务) 超过5秒时就超时退出
	if err = srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	log.Println("Server Exiting!!!")

}

/* package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

// Message takes incoming JSON payload for writing heart rate
type Message struct {
	BPM int
}

var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())

}

// web server
func run() error {
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// write blockchain when we receive an http request
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for heart rate (BPM)
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var msg Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	prevBlock := Blockchain[len(Blockchain)-1]
	newBlock := generateBlock(prevBlock, msg.BPM)

	if isBlockValid(newBlock, prevBlock) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}
	mutex.Unlock()

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hasing
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, BPM int) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}
*/
