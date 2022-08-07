// Copyright 2018-2020 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/ua"
)

func main() {

	var (
		endpoint = flag.String("endpoint", "opc.tcp://localhost:14840", "OPC UA Endpoint URL")
		nodeID   = flag.String("node", "ns=1;i=2345", "NodeID to read")
	)
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()
	log.SetFlags(0)

	ctx := context.Background()

	c := opcua.NewClient(*endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.CloseWithContext(ctx)

	id, err := ua.ParseNodeID(*nodeID)
	if err != nil {
		log.Fatalf("invalid node id: %v", err)
	}

	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	for range time.Tick(time.Second * 10) {
		log.Printf("Starting acquisition")
		resp, err := c.ReadWithContext(ctx, req)
		if err != nil {
			log.Printf("Read failed: %s", err)
			continue
		}
		if resp.Results[0].Status != ua.StatusOK {
			log.Printf("Status not OK: %v", resp.Results[0].Status)
		}
		log.Printf("%#v", resp.Results[0].Value.Value())
	}
}
