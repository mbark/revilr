// mgo - MongoDB driver for Go
//
// Copyright (c) 2010-2012 - Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package mgo_test

import (
	"fmt"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
	"net"
	"strings"
	"sync"
	"time"
)

func (s *S) TestNewSession(c *C) {
	session, err := mgo.Dial("localhost:40001")
	c.Assert(err, IsNil)
	defer session.Close()

	// Do a dummy operation to wait for connection.
	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"_id": 1})
	c.Assert(err, IsNil)

	// Tweak safety and query settings to ensure other has copied those.
	session.SetSafe(nil)
	session.SetBatch(-1)
	other := session.New()
	defer other.Close()
	session.SetSafe(&mgo.Safe{})

	// Clone was copied while session was unsafe, so no errors.
	otherColl := other.DB("mydb").C("mycoll")
	err = otherColl.Insert(M{"_id": 1})
	c.Assert(err, IsNil)

	// Original session was made safe again.
	err = coll.Insert(M{"_id": 1})
	c.Assert(err, NotNil)

	// With New(), each session has its own socket now.
	stats := mgo.GetStats()
	c.Assert(stats.MasterConns, Equals, 2)
	c.Assert(stats.SocketsInUse, Equals, 2)

	// Ensure query parameters were cloned.
	err = otherColl.Insert(M{"_id": 2})
	c.Assert(err, IsNil)

	// Ping the database to ensure the nonce has been received already.
	c.Assert(other.Ping(), IsNil)

	mgo.ResetStats()

	iter := otherColl.Find(M{}).Iter()
	c.Assert(err, IsNil)

	m := M{}
	ok := iter.Next(m)
	c.Assert(ok, Equals, true)
	err = iter.Close()
	c.Assert(err, IsNil)

	// If Batch(-1) is in effect, a single document must have been received.
	stats = mgo.GetStats()
	c.Assert(stats.ReceivedDocs, Equals, 1)
}

func (s *S) TestCloneSession(c *C) {
	session, err := mgo.Dial("localhost:40001")
	c.Assert(err, IsNil)
	defer session.Close()

	// Do a dummy operation to wait for connection.
	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"_id": 1})
	c.Assert(err, IsNil)

	// Tweak safety and query settings to ensure clone is copying those.
	session.SetSafe(nil)
	session.SetBatch(-1)
	clone := session.Clone()
	defer clone.Close()
	session.SetSafe(&mgo.Safe{})

	// Clone was copied while session was unsafe, so no errors.
	cloneColl := clone.DB("mydb").C("mycoll")
	err = cloneColl.Insert(M{"_id": 1})
	c.Assert(err, IsNil)

	// Original session was made safe again.
	err = coll.Insert(M{"_id": 1})
	c.Assert(err, NotNil)

	// With Clone(), same socket is shared between sessions now.
	stats := mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 1)
	c.Assert(stats.SocketRefs, Equals, 2)

	// Refreshing one of them should let the original socket go,
	// while preserving the safety settings.
	clone.Refresh()
	err = cloneColl.Insert(M{"_id": 1})
	c.Assert(err, IsNil)

	// Must have used another connection now.
	stats = mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 2)
	c.Assert(stats.SocketRefs, Equals, 2)

	// Ensure query parameters were cloned.
	err = cloneColl.Insert(M{"_id": 2})
	c.Assert(err, IsNil)

	// Ping the database to ensure the nonce has been received already.
	c.Assert(clone.Ping(), IsNil)

	mgo.ResetStats()

	iter := cloneColl.Find(M{}).Iter()
	c.Assert(err, IsNil)

	m := M{}
	ok := iter.Next(m)
	c.Assert(ok, Equals, true)
	err = iter.Close()
	c.Assert(err, IsNil)

	// If Batch(-1) is in effect, a single document must have been received.
	stats = mgo.GetStats()
	c.Assert(stats.ReceivedDocs, Equals, 1)
}

func (s *S) TestSetModeStrong(c *C) {
	session, err := mgo.Dial("localhost:40012")
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Monotonic, false)
	session.SetMode(mgo.Strong, false)

	c.Assert(session.Mode(), Equals, mgo.Strong)

	result := M{}
	cmd := session.DB("admin").C("$cmd")
	err = cmd.Find(M{"ismaster": 1}).One(&result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, true)

	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	// Wait since the sync also uses sockets.
	for len(session.LiveServers()) != 3 {
		c.Log("Waiting for cluster sync to finish...")
		time.Sleep(5e8)
	}

	stats := mgo.GetStats()
	c.Assert(stats.MasterConns, Equals, 1)
	c.Assert(stats.SlaveConns, Equals, 2)
	c.Assert(stats.SocketsInUse, Equals, 1)

	session.SetMode(mgo.Strong, true)

	stats = mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 0)
}

func (s *S) TestSetModeMonotonic(c *C) {
	// Must necessarily connect to a slave, otherwise the
	// master connection will be available first.
	session, err := mgo.Dial("localhost:40012")
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Monotonic, false)

	c.Assert(session.Mode(), Equals, mgo.Monotonic)

	result := M{}
	cmd := session.DB("admin").C("$cmd")
	err = cmd.Find(M{"ismaster": 1}).One(&result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, false)

	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	result = M{}
	err = cmd.Find(M{"ismaster": 1}).One(&result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, true)

	// Wait since the sync also uses sockets.
	for len(session.LiveServers()) != 3 {
		c.Log("Waiting for cluster sync to finish...")
		time.Sleep(5e8)
	}

	stats := mgo.GetStats()
	c.Assert(stats.MasterConns, Equals, 1)
	c.Assert(stats.SlaveConns, Equals, 2)
	c.Assert(stats.SocketsInUse, Equals, 2)

	session.SetMode(mgo.Monotonic, true)

	stats = mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 0)
}

func (s *S) TestSetModeMonotonicAfterStrong(c *C) {
	// Test that a strong session shifting to a monotonic
	// one preserves the socket untouched.

	session, err := mgo.Dial("localhost:40012")
	c.Assert(err, IsNil)
	defer session.Close()

	// Insert something to force a connection to the master.
	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	session.SetMode(mgo.Monotonic, false)

	// Wait since the sync also uses sockets.
	for len(session.LiveServers()) != 3 {
		c.Log("Waiting for cluster sync to finish...")
		time.Sleep(5e8)
	}

	// Master socket should still be reserved.
	stats := mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 1)

	// Confirm it's the master even though it's Monotonic by now.
	result := M{}
	cmd := session.DB("admin").C("$cmd")
	err = cmd.Find(M{"ismaster": 1}).One(&result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, true)
}

func (s *S) TestSetModeStrongAfterMonotonic(c *C) {
	// Test that shifting from Monotonic to Strong while
	// using a slave socket will keep the socket reserved
	// until the master socket is necessary, so that no
	// switch over occurs unless it's actually necessary.

	// Must necessarily connect to a slave, otherwise the
	// master connection will be available first.
	session, err := mgo.Dial("localhost:40012")
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Monotonic, false)

	// Ensure we're talking to a slave, and reserve the socket.
	result := M{}
	err = session.Run("ismaster", &result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, false)

	// Switch to a Strong session.
	session.SetMode(mgo.Strong, false)

	// Wait since the sync also uses sockets.
	for len(session.LiveServers()) != 3 {
		c.Log("Waiting for cluster sync to finish...")
		time.Sleep(5e8)
	}

	// Slave socket should still be reserved.
	stats := mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 1)

	// But any operation will switch it to the master.
	result = M{}
	err = session.Run("ismaster", &result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, true)
}

func (s *S) TestSetModeMonotonicWriteOnIteration(c *C) {
	// Must necessarily connect to a slave, otherwise the
	// master connection will be available first.
	session, err := mgo.Dial("localhost:40012")
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Monotonic, false)

	c.Assert(session.Mode(), Equals, mgo.Monotonic)

	coll1 := session.DB("mydb").C("mycoll1")
	coll2 := session.DB("mydb").C("mycoll2")

	ns := []int{40, 41, 42, 43, 44, 45, 46}
	for _, n := range ns {
		err := coll1.Insert(M{"n": n})
		c.Assert(err, IsNil)
	}

	// Release master so we can grab a slave again.
	session.Refresh()

	// Wait until synchronization is done.
	for {
		n, err := coll1.Count()
		c.Assert(err, IsNil)
		if n == len(ns) {
			break
		}
	}

	iter := coll1.Find(nil).Batch(2).Iter()
	i := 0
	m := M{}
	for iter.Next(&m) {
		i++
		if i > 3 {
			err := coll2.Insert(M{"n": 47 + i})
			c.Assert(err, IsNil)
		}
	}
	c.Assert(i, Equals, len(ns))
}

func (s *S) TestSetModeEventual(c *C) {
	// Must necessarily connect to a slave, otherwise the
	// master connection will be available first.
	session, err := mgo.Dial("localhost:40012")
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Eventual, false)

	c.Assert(session.Mode(), Equals, mgo.Eventual)

	result := M{}
	err = session.Run("ismaster", &result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, false)

	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	result = M{}
	err = session.Run("ismaster", &result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, false)

	// Wait since the sync also uses sockets.
	for len(session.LiveServers()) != 3 {
		c.Log("Waiting for cluster sync to finish...")
		time.Sleep(5e8)
	}

	stats := mgo.GetStats()
	c.Assert(stats.MasterConns, Equals, 1)
	c.Assert(stats.SlaveConns, Equals, 2)
	c.Assert(stats.SocketsInUse, Equals, 0)
}

func (s *S) TestSetModeEventualAfterStrong(c *C) {
	// Test that a strong session shifting to an eventual
	// one preserves the socket untouched.

	session, err := mgo.Dial("localhost:40012")
	c.Assert(err, IsNil)
	defer session.Close()

	// Insert something to force a connection to the master.
	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	session.SetMode(mgo.Eventual, false)

	// Wait since the sync also uses sockets.
	for len(session.LiveServers()) != 3 {
		c.Log("Waiting for cluster sync to finish...")
		time.Sleep(5e8)
	}

	// Master socket should still be reserved.
	stats := mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 1)

	// Confirm it's the master even though it's Eventual by now.
	result := M{}
	cmd := session.DB("admin").C("$cmd")
	err = cmd.Find(M{"ismaster": 1}).One(&result)
	c.Assert(err, IsNil)
	c.Assert(result["ismaster"], Equals, true)

	session.SetMode(mgo.Eventual, true)

	stats = mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 0)
}

func (s *S) TestPrimaryShutdownStrong(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	session, err := mgo.Dial("localhost:40021")
	c.Assert(err, IsNil)
	defer session.Close()

	// With strong consistency, this will open a socket to the master.
	result := &struct{ Host string }{}
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)

	// Kill the master.
	host := result.Host
	s.Stop(host)

	// This must fail, since the connection was broken.
	err = session.Run("serverStatus", result)
	c.Assert(err, Equals, io.EOF)

	// With strong consistency, it fails again until reset.
	err = session.Run("serverStatus", result)
	c.Assert(err, Equals, io.EOF)

	session.Refresh()

	// Now we should be able to talk to the new master.
	// Increase the timeout since this may take quite a while.
	session.SetSyncTimeout(3 * time.Minute)

	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)
	c.Assert(result.Host, Not(Equals), host)

	// Insert some data to confirm it's indeed a master.
	err = session.DB("mydb").C("mycoll").Insert(M{"n": 42})
	c.Assert(err, IsNil)
}

func (s *S) TestPrimaryHiccup(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	session, err := mgo.Dial("localhost:40021")
	c.Assert(err, IsNil)
	defer session.Close()

	// With strong consistency, this will open a socket to the master.
	result := &struct{ Host string }{}
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)

	// Establish a few extra sessions to create spare sockets to
	// the master. This increases a bit the chances of getting an
	// incorrect cached socket.
	var sessions []*mgo.Session
	for i := 0; i < 20; i++ {
		sessions = append(sessions, session.Copy())
		err = sessions[len(sessions)-1].Run("serverStatus", result)
		c.Assert(err, IsNil)
	}
	for i := range sessions {
		sessions[i].Close()
	}

	// Kill the master, but bring it back immediatelly.
	host := result.Host
	s.Stop(host)
	s.StartAll()

	// This must fail, since the connection was broken.
	err = session.Run("serverStatus", result)
	c.Assert(err, Equals, io.EOF)

	// With strong consistency, it fails again until reset.
	err = session.Run("serverStatus", result)
	c.Assert(err, Equals, io.EOF)

	session.Refresh()

	// Now we should be able to talk to the new master.
	// Increase the timeout since this may take quite a while.
	session.SetSyncTimeout(3 * time.Minute)

	// Insert some data to confirm it's indeed a master.
	err = session.DB("mydb").C("mycoll").Insert(M{"n": 42})
	c.Assert(err, IsNil)
}

func (s *S) TestPrimaryShutdownMonotonic(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	session, err := mgo.Dial("localhost:40021")
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Insert something to force a switch to the master.
	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	// Wait a bit for this to be synchronized to slaves.
	time.Sleep(3 * time.Second)

	result := &struct{ Host string }{}
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)

	// Kill the master.
	host := result.Host
	s.Stop(host)

	// This must fail, since the connection was broken.
	err = session.Run("serverStatus", result)
	c.Assert(err, Equals, io.EOF)

	// With monotonic consistency, it fails again until reset.
	err = session.Run("serverStatus", result)
	c.Assert(err, Equals, io.EOF)

	session.Refresh()

	// Now we should be able to talk to the new master.
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)
	c.Assert(result.Host, Not(Equals), host)
}

func (s *S) TestPrimaryShutdownMonotonicWithSlave(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	session, err := mgo.Dial("localhost:40021")
	c.Assert(err, IsNil)
	defer session.Close()

	ssresult := &struct{ Host string }{}
	imresult := &struct{ IsMaster bool }{}

	// Figure the master while still using the strong session.
	err = session.Run("serverStatus", ssresult)
	c.Assert(err, IsNil)
	err = session.Run("isMaster", imresult)
	c.Assert(err, IsNil)
	master := ssresult.Host
	c.Assert(imresult.IsMaster, Equals, true, Commentf("%s is not the master", master))

	// Create new monotonic session with an explicit address to ensure
	// a slave is synchronized before the master, otherwise a connection
	// with the master may be used below for lack of other options.
	var addr string
	switch {
	case strings.HasSuffix(ssresult.Host, ":40021"):
		addr = "localhost:40022"
	case strings.HasSuffix(ssresult.Host, ":40022"):
		addr = "localhost:40021"
	case strings.HasSuffix(ssresult.Host, ":40023"):
		addr = "localhost:40021"
	default:
		c.Fatal("Unknown host: ", ssresult.Host)
	}

	session, err = mgo.Dial(addr)
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Check the address of the socket associated with the monotonic session.
	c.Log("Running serverStatus and isMaster with monotonic session")
	err = session.Run("serverStatus", ssresult)
	c.Assert(err, IsNil)
	err = session.Run("isMaster", imresult)
	c.Assert(err, IsNil)
	slave := ssresult.Host
	c.Assert(imresult.IsMaster, Equals, false, Commentf("%s is not a slave", slave))

	c.Assert(master, Not(Equals), slave)

	// Kill the master.
	s.Stop(master)

	// Session must still be good, since we were talking to a slave.
	err = session.Run("serverStatus", ssresult)
	c.Assert(err, IsNil)

	c.Assert(ssresult.Host, Equals, slave,
		Commentf("Monotonic session moved from %s to %s", slave, ssresult.Host))

	// If we try to insert something, it'll have to hold until the new
	// master is available to move the connection, and work correctly.
	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	// Must now be talking to the new master.
	err = session.Run("serverStatus", ssresult)
	c.Assert(err, IsNil)
	err = session.Run("isMaster", imresult)
	c.Assert(err, IsNil)
	c.Assert(imresult.IsMaster, Equals, true, Commentf("%s is not the master", master))

	// ... which is not the old one, since it's still dead.
	c.Assert(ssresult.Host, Not(Equals), master)
}

func (s *S) TestPrimaryShutdownEventual(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	session, err := mgo.Dial("localhost:40021")
	c.Assert(err, IsNil)
	defer session.Close()

	result := &struct{ Host string }{}
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)
	master := result.Host

	session.SetMode(mgo.Eventual, true)

	// Should connect to the master when needed.
	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	// Wait a bit for this to be synchronized to slaves.
	time.Sleep(3 * time.Second)

	// Kill the master.
	s.Stop(master)

	// Should still work, with the new master now.
	coll = session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1})
	c.Assert(err, IsNil)

	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)
	c.Assert(result.Host, Not(Equals), master)
}

func (s *S) TestPreserveSocketCountOnSync(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	stats := mgo.GetStats()
	for stats.MasterConns+stats.SlaveConns != 3 {
		stats = mgo.GetStats()
		c.Log("Waiting for all connections to be established...")
		time.Sleep(5e8)
	}

	c.Assert(stats.SocketsAlive, Equals, 3)

	// Kill the master (with rs1, 'a' is always the master).
	s.Stop("localhost:40011")

	// Wait for the logic to run for a bit and bring it back.
	go func() {
		time.Sleep(5e9)
		s.StartAll()
	}()

	// Do an action to kick the resync logic in, and also to
	// wait until the cluster recognizes the server is back.
	result := struct{ Ok bool }{}
	err = session.Run("getLastError", &result)
	c.Assert(err, IsNil)
	c.Assert(result.Ok, Equals, true)

	for i := 0; i != 20; i++ {
		stats = mgo.GetStats()
		if stats.SocketsAlive == 3 {
			break
		}
		c.Logf("Waiting for 3 sockets alive, have %d", stats.SocketsAlive)
		time.Sleep(5e8)
	}

	// Ensure the number of sockets is preserved after syncing.
	stats = mgo.GetStats()
	c.Assert(stats.SocketsAlive, Equals, 3)
	c.Assert(stats.SocketsInUse, Equals, 1)
	c.Assert(stats.SocketRefs, Equals, 1)
}

// Connect to the master of a deployment with a single server,
// run an insert, and then ensure the insert worked and that a
// single connection was established.
func (s *S) TestTopologySyncWithSingleMaster(c *C) {
	// Use hostname here rather than IP, to make things trickier.
	session, err := mgo.Dial("localhost:40001")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"a": 1, "b": 2})
	c.Assert(err, IsNil)

	// One connection used for discovery. Master socket recycled for
	// insert. Socket is reserved after insert.
	stats := mgo.GetStats()
	c.Assert(stats.MasterConns, Equals, 1)
	c.Assert(stats.SlaveConns, Equals, 0)
	c.Assert(stats.SocketsInUse, Equals, 1)

	// Refresh session and socket must be released.
	session.Refresh()
	stats = mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 0)
}

func (s *S) TestTopologySyncWithSlaveSeed(c *C) {
	// That's supposed to be a slave. Must run discovery
	// and find out master to insert successfully.
	session, err := mgo.Dial("localhost:40012")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")
	coll.Insert(M{"a": 1, "b": 2})

	result := struct{ Ok bool }{}
	err = session.Run("getLastError", &result)
	c.Assert(err, IsNil)
	c.Assert(result.Ok, Equals, true)

	// One connection to each during discovery. Master
	// socket recycled for insert.
	stats := mgo.GetStats()
	c.Assert(stats.MasterConns, Equals, 1)
	c.Assert(stats.SlaveConns, Equals, 2)

	// Only one socket reference alive, in the master socket owned
	// by the above session.
	c.Assert(stats.SocketsInUse, Equals, 1)

	// Refresh it, and it must be gone.
	session.Refresh()
	stats = mgo.GetStats()
	c.Assert(stats.SocketsInUse, Equals, 0)
}

func (s *S) TestSyncTimeout(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	session, err := mgo.Dial("localhost:40001")
	c.Assert(err, IsNil)
	defer session.Close()

	s.Stop("localhost:40001")

	timeout := 3 * time.Second
	session.SetSyncTimeout(timeout)
	started := time.Now()

	// Do something.
	result := struct{ Ok bool }{}
	err = session.Run("getLastError", &result)
	c.Assert(err, ErrorMatches, "no reachable servers")
	c.Assert(started.Before(time.Now().Add(-timeout)), Equals, true)
	c.Assert(started.After(time.Now().Add(-timeout*2)), Equals, true)
}

func (s *S) TestDialWithTimeout(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	timeout := 2 * time.Second
	started := time.Now()

	// 40009 isn't used by the test servers.
	session, err := mgo.DialWithTimeout("localhost:40009", timeout)
	if session != nil {
		session.Close()
	}
	c.Assert(err, ErrorMatches, "no reachable servers")
	c.Assert(session, IsNil)
	c.Assert(started.Before(time.Now().Add(-timeout)), Equals, true)
	c.Assert(started.After(time.Now().Add(-timeout*2)), Equals, true)
}

func (s *S) TestDirect(c *C) {
	session, err := mgo.Dial("localhost:40012?connect=direct")
	c.Assert(err, IsNil)
	defer session.Close()

	// We know that server is a slave.
	session.SetMode(mgo.Monotonic, true)

	result := &struct{ Host string }{}
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)
	c.Assert(strings.HasSuffix(result.Host, ":40012"), Equals, true)

	stats := mgo.GetStats()
	c.Assert(stats.SocketsAlive, Equals, 1)
	c.Assert(stats.SocketsInUse, Equals, 1)
	c.Assert(stats.SocketRefs, Equals, 1)

	// We've got no master, so it'll timeout.
	session.SetSyncTimeout(5e8 * time.Nanosecond)

	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"test": 1})
	c.Assert(err, ErrorMatches, "no reachable servers")

	// Writing to the local database is okay.
	coll = session.DB("local").C("mycoll")
	defer coll.RemoveAll(nil)
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id})
	c.Assert(err, IsNil)

	// Data was stored in the right server.
	n, err := coll.Find(M{"_id": id}).Count()
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 1)

	// Server hasn't changed.
	result.Host = ""
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)
	c.Assert(strings.HasSuffix(result.Host, ":40012"), Equals, true)
}

func (s *S) TestDirectToUnknownStateMember(c *C) {
	session, err := mgo.Dial("localhost:40041?connect=direct")
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	result := &struct{ Host string }{}
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)
	c.Assert(strings.HasSuffix(result.Host, ":40041"), Equals, true)

	// We've got no master, so it'll timeout.
	session.SetSyncTimeout(5e8 * time.Nanosecond)

	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(M{"test": 1})
	c.Assert(err, ErrorMatches, "no reachable servers")

	// Slave is still reachable.
	result.Host = ""
	err = session.Run("serverStatus", result)
	c.Assert(err, IsNil)
	c.Assert(strings.HasSuffix(result.Host, ":40041"), Equals, true)
}

type OpCounters struct {
	Insert  int
	Query   int
	Update  int
	Delete  int
	GetMore int
	Command int
}

func getOpCounters(server string) (c *OpCounters, err error) {
	session, err := mgo.Dial(server + "?connect=direct")
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	result := struct{ OpCounters }{}
	err = session.Run("serverStatus", &result)
	return &result.OpCounters, err
}

func (s *S) TestMonotonicSlaveOkFlagWithMongos(c *C) {
	session, err := mgo.Dial("localhost:40021")
	c.Assert(err, IsNil)
	defer session.Close()

	ssresult := &struct{ Host string }{}
	imresult := &struct{ IsMaster bool }{}

	// Figure the master while still using the strong session.
	err = session.Run("serverStatus", ssresult)
	c.Assert(err, IsNil)
	err = session.Run("isMaster", imresult)
	c.Assert(err, IsNil)
	master := ssresult.Host
	c.Assert(imresult.IsMaster, Equals, true, Commentf("%s is not the master", master))

	// Collect op counters for everyone.
	opc21a, err := getOpCounters("localhost:40021")
	c.Assert(err, IsNil)
	opc22a, err := getOpCounters("localhost:40022")
	c.Assert(err, IsNil)
	opc23a, err := getOpCounters("localhost:40023")
	c.Assert(err, IsNil)

	// Do a SlaveOk query through MongoS

	mongos, err := mgo.Dial("localhost:40202")
	c.Assert(err, IsNil)
	defer mongos.Close()

	mongos.SetMode(mgo.Monotonic, true)

	coll := mongos.DB("mydb").C("mycoll")
	result := &struct{}{}
	for i := 0; i != 5; i++ {
		err := coll.Find(nil).One(result)
		c.Assert(err, Equals, mgo.ErrNotFound)
	}

	// Collect op counters for everyone again.
	opc21b, err := getOpCounters("localhost:40021")
	c.Assert(err, IsNil)
	opc22b, err := getOpCounters("localhost:40022")
	c.Assert(err, IsNil)
	opc23b, err := getOpCounters("localhost:40023")
	c.Assert(err, IsNil)

	masterPort := master[strings.Index(master, ":")+1:]

	var masterDelta, slaveDelta int
	switch masterPort {
	case "40021":
		masterDelta = opc21b.Query - opc21a.Query
		slaveDelta = (opc22b.Query - opc22a.Query) + (opc23b.Query - opc23a.Query)
	case "40022":
		masterDelta = opc22b.Query - opc22a.Query
		slaveDelta = (opc21b.Query - opc21a.Query) + (opc23b.Query - opc23a.Query)
	case "40023":
		masterDelta = opc23b.Query - opc23a.Query
		slaveDelta = (opc21b.Query - opc21a.Query) + (opc22b.Query - opc22a.Query)
	default:
		c.Fatal("Uh?")
	}

	c.Check(masterDelta, Equals, 0) // Just the counting itself.
	c.Check(slaveDelta, Equals, 5)  // The counting for both, plus 5 queries above.
}

func (s *S) TestRemovalOfClusterMember(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	master, err := mgo.Dial("localhost:40021")
	c.Assert(err, IsNil)
	defer master.Close()

	// Wait for cluster to fully sync up.
	for i := 0; i < 10; i++ {
		if len(master.LiveServers()) == 3 {
			break
		}
		time.Sleep(5e8)
	}
	if len(master.LiveServers()) != 3 {
		c.Fatalf("Test started with bad cluster state: %v", master.LiveServers())
	}

	result := &struct {
		IsMaster bool
		Me       string
	}{}
	slave := master.Copy()
	slave.SetMode(mgo.Monotonic, true) // Monotonic can hold a non-master socket persistently.
	err = slave.Run("isMaster", result)
	c.Assert(err, IsNil)
	c.Assert(result.IsMaster, Equals, false)
	slaveAddr := result.Me

	defer func() {
		master.Refresh()
		master.Run(bson.D{{"$eval", `rs.add("` + slaveAddr + `")`}}, nil)
		master.Close()
		slave.Close()
	}()

	c.Logf("========== Removing slave: %s ==========", slaveAddr)

	master.Run(bson.D{{"$eval", `rs.remove("` + slaveAddr + `")`}}, nil)
	err = master.Ping()
	c.Assert(err, Equals, io.EOF)

	master.Refresh()

	// Give the cluster a moment to catch up by doing a roundtrip to the master.
	err = master.Ping()
	c.Assert(err, IsNil)

	time.Sleep(3e9)

	// This must fail since the slave has been taken off the cluster.
	err = slave.Ping()
	c.Assert(err, NotNil)

	for i := 0; i < 15; i++ {
		if len(master.LiveServers()) == 2 {
			break
		}
		time.Sleep(time.Second)
	}
	live := master.LiveServers()
	if len(live) != 2 {
		c.Errorf("Removed server still considered live: %#s", live)
	}

	c.Log("========== Test succeeded. ==========")
}

func (s *S) TestSocketLimit(c *C) {
	if *fast {
		c.Skip("-fast")
	}
	const socketLimit = 64
	restore := mgo.HackSocketsPerServer(socketLimit)
	defer restore()

	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	stats := mgo.GetStats()
	for stats.MasterConns+stats.SlaveConns != 3 {
		stats = mgo.GetStats()
		c.Log("Waiting for all connections to be established...")
		time.Sleep(5e8)
	}
	c.Assert(stats.SocketsAlive, Equals, 3)

	// Consume the whole limit for the master.
	var master []*mgo.Session
	for i := 0; i < socketLimit; i++ {
		s := session.Copy()
		defer s.Close()
		err := s.Ping()
		c.Assert(err, IsNil)
		master = append(master, s)
	}

	before := time.Now()
	go func() {
		time.Sleep(3e9)
		master[0].Refresh()
	}()

	// Now a single ping must block, since it would need another
	// connection to the master, over the limit. Once the goroutine
	// above releases its socket, it should move on.
	session.Ping()
	delay := time.Now().Sub(before)
	c.Assert(delay > 3e9, Equals, true)
	c.Assert(delay < 6e9, Equals, true)
}

func (s *S) TestSetModeEventualIterBug(c *C) {
	session1, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session1.Close()

	session1.SetMode(mgo.Eventual, false)

	coll1 := session1.DB("mydb").C("mycoll")

	const N = 100
	for i := 0; i < N; i++ {
		err = coll1.Insert(M{"_id": i})
		c.Assert(err, IsNil)
	}

	c.Logf("Waiting until secondary syncs")
	for {
		n, err := coll1.Count()
		c.Assert(err, IsNil)
		if n == N {
			c.Logf("Found all")
			break
		}
	}

	session2, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session2.Close()

	session2.SetMode(mgo.Eventual, false)

	coll2 := session2.DB("mydb").C("mycoll")

	i := 0
	iter := coll2.Find(nil).Batch(10).Iter()
	var result struct{}
	for iter.Next(&result) {
		i++
	}
	c.Assert(iter.Close(), Equals, nil)
	c.Assert(i, Equals, N)
}

func (s *S) TestCustomDial(c *C) {
	dials := make(chan bool, 16)
	dial := func(addr net.Addr) (net.Conn, error) {
		tcpaddr, ok := addr.(*net.TCPAddr)
		if !ok {
			return nil, fmt.Errorf("unexpected address type: %T", addr)
		}
		dials <- true
		return net.DialTCP("tcp", nil, tcpaddr)
	}
	info := mgo.DialInfo{
		Addrs: []string{"localhost:40012"},
		Dial:  dial,
	}

	// Use hostname here rather than IP, to make things trickier.
	session, err := mgo.DialWithInfo(&info)
	c.Assert(err, IsNil)
	defer session.Close()

	const N = 3
	for i := 0; i < N; i++ {
		select {
		case <-dials:
		case <-time.After(5 * time.Second):
			c.Fatalf("expected %d dials, got %d", N, i)
		}
	}
	select {
	case <-dials:
		c.Fatalf("got more dials than expected")
	case <-time.After(100 * time.Millisecond):
	}
}

func (s *S) TestPrimaryShutdownOnAuthShard(c *C) {
	if *fast {
		c.Skip("-fast")
	}

	// Dial the shard.
	session, err := mgo.Dial("localhost:40203")
	c.Assert(err, IsNil)
	defer session.Close()

	// Login and insert something to make it more realistic.
	session.DB("admin").Login("root", "rapadura")
	coll := session.DB("mydb").C("mycoll")
	err = coll.Insert(bson.M{"n": 1})
	c.Assert(err, IsNil)

	// Dial the replica set to figure the master out.
	rs, err := mgo.Dial("root:rapadura@localhost:40031")
	c.Assert(err, IsNil)
	defer rs.Close()

	// With strong consistency, this will open a socket to the master.
	result := &struct{ Host string }{}
	err = rs.Run("serverStatus", result)
	c.Assert(err, IsNil)

	// Kill the master.
	host := result.Host
	s.Stop(host)

	// This must fail, since the connection was broken.
	err = rs.Run("serverStatus", result)
	c.Assert(err, Equals, io.EOF)

	// This won't work because the master just died.
	err = coll.Insert(bson.M{"n": 2})
	c.Assert(err, NotNil)

	// Refresh session and wait for re-election.
	session.Refresh()
	for i := 0; i < 60; i++ {
		err = coll.Insert(bson.M{"n": 3})
		if err == nil {
			break
		}
		c.Logf("Waiting for replica set to elect a new master. Last error: %v", err)
		time.Sleep(500 * time.Millisecond)
	}
	c.Assert(err, IsNil)

	count, err := coll.Count()
	c.Assert(count > 1, Equals, true)
}

func (s *S) TestNearestSecondary(c *C) {
	defer mgo.HackPingDelay(3 * time.Second)()

	rs1a := "127.0.0.1:40011"
	rs1b := "127.0.0.1:40012"
	rs1c := "127.0.0.1:40013"
	s.Freeze(rs1b)

	session, err := mgo.Dial(rs1a)
	c.Assert(err, IsNil)
	defer session.Close()

	// Wait for the sync up to run through the first couple of servers.
	for len(session.LiveServers()) != 2 {
		c.Log("Waiting for two servers to be alive...")
		time.Sleep(100 * time.Millisecond)
	}

	// Extra delay to ensure the third server gets penalized.
	time.Sleep(500 * time.Millisecond)

	// Release third server.
	s.Thaw(rs1b)

	// Wait for it to come up.
	for len(session.LiveServers()) != 3 {
		c.Log("Waiting for all servers to be alive...")
		time.Sleep(100 * time.Millisecond)
	}

	session.SetMode(mgo.Monotonic, true)
	var result struct{ Host string }

	// See which slave picks the line, several times to avoid chance.
	for i := 0; i < 10; i++ {
		session.Refresh()
		err = session.Run("serverStatus", &result)
		c.Assert(err, IsNil)
		c.Assert(hostPort(result.Host), Equals, hostPort(rs1c))
	}

	if *fast {
		// Don't hold back for several seconds.
		return
	}

	// Now hold the other server for long enough to penalize it.
	s.Freeze(rs1c)
	time.Sleep(5 * time.Second)
	s.Thaw(rs1c)

	// Wait for the ping to be processed.
	time.Sleep(500 * time.Millisecond)

	// Repeating the test should now pick the former server consistently.
	for i := 0; i < 10; i++ {
		session.Refresh()
		err = session.Run("serverStatus", &result)
		c.Assert(err, IsNil)
		c.Assert(hostPort(result.Host), Equals, hostPort(rs1b))
	}
}

func (s *S) TestConnectCloseConcurrency(c *C) {
	restore := mgo.HackPingDelay(1)
	defer restore()
	var wg sync.WaitGroup
	const n = 500
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			session, err := mgo.Dial("localhost:40001")
			if err != nil {
				c.Fatal(err)
			}
			time.Sleep(1)
			session.Close()
		}()
	}
	wg.Wait()
}

func (s *S) TestSelectServers(c *C) {
	if !s.versionAtLeast(2, 2) {
		c.Skip("read preferences introduced in 2.2")
	}

	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	session.SetMode(mgo.Eventual, true)

	var result struct{ Host string }

	session.Refresh()
	session.SelectServers(bson.D{{"rs1", "b"}})
	err = session.Run("serverStatus", &result)
	c.Assert(err, IsNil)
	c.Assert(hostPort(result.Host), Equals, "40012")

	session.Refresh()
	session.SelectServers(bson.D{{"rs1", "c"}})
	err = session.Run("serverStatus", &result)
	c.Assert(err, IsNil)
	c.Assert(hostPort(result.Host), Equals, "40013")
}

func (s *S) TestSelectServersWithMongos(c *C) {
	if !s.versionAtLeast(2, 2) {
		c.Skip("read preferences introduced in 2.2")
	}

	session, err := mgo.Dial("localhost:40021")
	c.Assert(err, IsNil)
	defer session.Close()

	ssresult := &struct{ Host string }{}
	imresult := &struct{ IsMaster bool }{}

	// Figure the master while still using the strong session.
	err = session.Run("serverStatus", ssresult)
	c.Assert(err, IsNil)
	err = session.Run("isMaster", imresult)
	c.Assert(err, IsNil)
	master := ssresult.Host
	c.Assert(imresult.IsMaster, Equals, true, Commentf("%s is not the master", master))

	var slave1, slave2 string
	switch hostPort(master) {
	case "40021":
		slave1, slave2 = "b", "c"
	case "40022":
		slave1, slave2 = "a", "c"
	case "40023":
		slave1, slave2 = "a", "b"
	}

	// Collect op counters for everyone.
	opc21a, err := getOpCounters("localhost:40021")
	c.Assert(err, IsNil)
	opc22a, err := getOpCounters("localhost:40022")
	c.Assert(err, IsNil)
	opc23a, err := getOpCounters("localhost:40023")
	c.Assert(err, IsNil)

	// Do a SlaveOk query through MongoS
	mongos, err := mgo.Dial("localhost:40202")
	c.Assert(err, IsNil)
	defer mongos.Close()

	mongos.SetMode(mgo.Monotonic, true)

	mongos.Refresh()
	mongos.SelectServers(bson.D{{"rs2", slave1}})
	coll := mongos.DB("mydb").C("mycoll")
	result := &struct{}{}
	for i := 0; i != 5; i++ {
		err := coll.Find(nil).One(result)
		c.Assert(err, Equals, mgo.ErrNotFound)
	}

	mongos.Refresh()
	mongos.SelectServers(bson.D{{"rs2", slave2}})
	coll = mongos.DB("mydb").C("mycoll")
	for i := 0; i != 7; i++ {
		err := coll.Find(nil).One(result)
		c.Assert(err, Equals, mgo.ErrNotFound)
	}

	// Collect op counters for everyone again.
	opc21b, err := getOpCounters("localhost:40021")
	c.Assert(err, IsNil)
	opc22b, err := getOpCounters("localhost:40022")
	c.Assert(err, IsNil)
	opc23b, err := getOpCounters("localhost:40023")
	c.Assert(err, IsNil)

	switch hostPort(master) {
	case "40021":
		c.Check(opc21b.Query - opc21a.Query, Equals, 0)
		c.Check(opc22b.Query - opc22a.Query, Equals, 5)
		c.Check(opc23b.Query - opc23a.Query, Equals, 7)
	case "40022":
		c.Check(opc21b.Query - opc21a.Query, Equals, 5)
		c.Check(opc22b.Query - opc22a.Query, Equals, 0)
		c.Check(opc23b.Query - opc23a.Query, Equals, 7)
	case "40023":
		c.Check(opc21b.Query - opc21a.Query, Equals, 5)
		c.Check(opc22b.Query - opc22a.Query, Equals, 7)
		c.Check(opc23b.Query - opc23a.Query, Equals, 0)
	default:
		c.Fatal("Uh?")
	}
}
