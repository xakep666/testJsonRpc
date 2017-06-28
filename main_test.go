package main

import "testing"
import (
	"os"
	"net"
	"net/rpc/jsonrpc"
	"github.com/stretchr/testify/assert"
	"time"
	"testJsonRpc/dao"
	"testJsonRpc/model"
)

func TestMain(m *testing.M) {
	if err:=dao.SetupDb("users_test.db"); err!=nil {
		panic(err)
	}
	go startServer()
	defer dao.GetDb().Close()
	time.Sleep(time.Millisecond)
	code:=m.Run()
	// cleanup
	os.RemoveAll("users_test.db")
	os.Exit(code)
}

func newConn() net.Conn{
	conn, err := net.Dial("tcp", "localhost:8222")

	if err != nil {
		panic(err)
	}
	return conn
}

func TestRegister(t *testing.T) {
	conn:=newConn()
	defer conn.Close()
	c:=jsonrpc.NewClient(conn)
	defer c.Close()
	var reply bool
	assert.NoError(t, c.Call("Users.Register","testuserReg", &reply))
	assert.True(t, reply)
	var reply1 bool
	assert.Error(t, c.Call("Users.Register","testuserReg", &reply1), "user testuserReg already registered")
}

func TestGet(t *testing.T) {
	conn:=newConn()
	defer conn.Close()
	c:=jsonrpc.NewClient(conn)
	defer c.Close()
	var user model.User
	var reply bool
	assert.NoError(t, c.Call("Users.Register","testuserGet", &reply))
	assert.True(t, reply)
	assert.NoError(t, c.Call("Users.GetByLogin", "testuserGet",&user))
	assert.NotNil(t, user)
	assert.Equal(t, "testuserGet", user.Login)
	var user1 model.User
	assert.Error(t, c.Call("Users.GetByLogin", "testuserGetFail",&user1), "user testuserGetFail is not registered")
}

func TestEdit(t *testing.T) {
	conn:=newConn()
	defer conn.Close()
	c:=jsonrpc.NewClient(conn)
	defer c.Close()
	var reply bool
	assert.NoError(t, c.Call("Users.Register","testuserEdit1", &reply))
	assert.True(t, reply)
	assert.NoError(t, c.Call("Users.Register","testuserEdit2", &reply))
	assert.True(t, reply)
	assert.NoError(t, c.Call("Users.Register","testuserEdit3", &reply))
	assert.True(t, reply)
	assert.NoError(t, c.Call("Users.Edit", dao.UserEditArgs{
		Login:"testuserEdit1",
		NewData:model.User{
			Login:"testuserEdit4",
		},
	},&reply))
	assert.True(t, reply)
	var user model.User
	assert.NoError(t, c.Call("Users.GetByLogin","testuserEdit4", &user))
	assert.Equal(t, "testuserEdit4", user.Login)
	var reply1 bool
	assert.Error(t, c.Call("Users.Edit", dao.UserEditArgs{
		Login:"testuserEdit2",
		NewData:model.User{
			Login:"testuserEdit3",
		},
	},&reply1), "user testuserEdit3 is already registered")
}