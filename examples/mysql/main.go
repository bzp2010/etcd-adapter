// Copyright api7.ai
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import (
	"context"
	"net"
	"time"

	"go.uber.org/zap"

	adapter "github.com/api7/etcd-adapter"
	"github.com/api7/etcd-adapter/backends/mysql"
)

func main() {
	opts := &adapter.AdapterOptions{
		Logger:  zap.NewExample(),
		Backend: adapter.BackendMySQL,
		MySQLOptions: &mysql.Options{
			DSN: "root@tcp(127.0.0.1:4000)/apisix",
		},
	}
	a := adapter.NewEtcdAdapter(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	ln, err := net.Listen("tcp", "127.0.0.1:12379")
	if err != nil {
		panic(err)
	}
	go func() {
		if err := a.Serve(context.Background(), ln); err != nil {
			panic(err)
		}
	}()
	<-ctx.Done()
	if err := a.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
