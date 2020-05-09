/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package exec

import (
	"bytes"

	"k8s.io/kubernetes/pkg/kubelet/util/ioutils"
	"k8s.io/kubernetes/pkg/probe"

	"k8s.io/klog"
	"k8s.io/utils/exec"
)

const (
	maxReadLength = 10 * 1 << 10 // 10KB
)

// New creates a Prober.
// 创建一个探针
func New() Prober {
	return execProber{}
}

// Prober is an interface defining the Probe object for container readiness/liveness checks.
// 探针接口
type Prober interface {
	// 输入：命令
	// 输出：结果、信息、错误
	Probe(e exec.Cmd) (probe.Result, string, error)
}

// 探针结构
type execProber struct{}

// Probe executes a command to check the liveness/readiness of container
// from executing a command. Returns the Result status, command output, and
// errors if any.
func (pr execProber) Probe(e exec.Cmd) (probe.Result, string, error) {
	// 缓存
	var dataBuffer bytes.Buffer
	writer := ioutils.LimitWriter(&dataBuffer, maxReadLength)

	// 重定向命令的输出
	e.SetStderr(writer)
	e.SetStdout(writer)

	// 执行命令
	err := e.Start()

	if err == nil {
		// 执行成功后，等待命令结束
		err = e.Wait()
	}

	// 取命令输出
	data := dataBuffer.Bytes()

	// 打个日志
	klog.V(4).Infof("Exec probe response: %q", string(data))

	if err != nil {
		// 尝试转成退出类型
		exit, ok := err.(exec.ExitError)
		if ok {
			// 返回值为退出类型
			if exit.ExitStatus() == 0 {
				// 返回成功
				return probe.Success, string(data), nil
			}

			// 返回失败
			return probe.Failure, string(data), nil
		}

		// 返回未知
		return probe.Unknown, "", err
	}

	// 命令执行成功
	return probe.Success, string(data), nil
}
