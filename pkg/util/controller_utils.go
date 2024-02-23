// Tencent is pleased to support the open source community by making Polaris available.
//
// Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://opensource.org/licenses/BSD-3-Clause
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package util

import (
	"errors"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	keyFunc         = cache.DeletionHandlingMetaNamespaceKeyFunc
	NamespacePrefix = "Namespace~"
	ServicePrefix   = "Service~"
	ConfigMapPrefix = "ConfigMap~"
)

func GenObjectQueueKey(obj interface{}) (string, error) {
	key, err := keyFunc(obj)
	if err != nil {
		return "", err
	}
	return key, err
}

// GetOriginKeyWithResyncQueueKey 通过同步任务key生成原始key
func GetOriginKeyWithResyncQueueKey(key string) string {
	return key[:len(key)-len("~resync")]
}

// GenResourceResyncQueueKey 通过原始key生成用于同步任务的key便于区分不同的任务
func GenResourceResyncQueueKey(key string) string {
	return key + "~" + "resync"
}

// GenQueueKeyWithFlag 在 namespace 的事件流程中使用。
// 产生 service queue 中的 key，flag 表示添加时是否是北极星的服务
func GenQueueKeyWithFlag(svc interface{}, flag string) (string, error) {
	key, err := keyFunc(svc)
	if err != nil {
		return "", err
	}
	key += "~" + flag

	return key, nil
}

// GetResourceRealKeyWithFlag 从 service queue 中的 key ，解析出 namespace、service、flag
func GetResourceRealKeyWithFlag(queueKey string) (string, string, string, string, error) {
	if queueKey == "" {
		return "", "", "", "", nil
	}
	op := ""
	ss := strings.Split(queueKey, "~")
	namespace, service, err := cache.SplitMetaNamespaceKey(ss[0])
	if err != nil {
		return "", "", "", "", err
	}
	if len(ss) != 1 {
		op = ss[1]
	}
	return ss[0], namespace, service, op, nil
}

// GenConfigMapQueueKeyWithFlag 在 namespace 的事件流程中使用。
// 产生 service queue 中的 key，flag 表示添加时是否是北极星的服务
func GenConfigMapQueueKeyWithFlag(svc *v1.ConfigMap, flag string) (string, error) {
	key, err := keyFunc(svc)
	if err != nil {
		return "", err
	}
	key += "~" + flag

	return key, nil
}

// GenResourceMapQueueKey 产生 service 中 queue 中用的 key
func GenResourceMapQueueKey(val interface{}) (string, error) {
	key, err := keyFunc(val)
	if err != nil {
		return "", err
	}

	switch val.(type) {
	case *v1.Namespace:
		return NamespacePrefix + key, nil
	case *v1.ConfigMap:
		return ConfigMapPrefix + key, nil
	case *v1.Service:
		return ServicePrefix + key, nil
	}
	return "", errors.New("not invalid kubernetes resource type")
}

func IsServiceKey(key string) (string, bool) {
	if strings.HasPrefix(key, ServicePrefix) {
		return strings.TrimPrefix(key, ServicePrefix), true
	}
	return key, false
}

func IsConfigMapKey(key string) (string, bool) {
	if strings.HasPrefix(key, ConfigMapPrefix) {
		return strings.TrimPrefix(key, ConfigMapPrefix), true
	}
	return key, false
}

func GetNamespace(svr *v1.Service) string {
	if v, ok := svr.GetAnnotations()[PolarisOverideNamespace]; ok && v != "" {
		return v
	}

	return svr.GetNamespace()
}

func GetServiceName(svr *v1.Service) string {
	if v, ok := svr.GetAnnotations()[PolarisOverideService]; ok && v != "" {
		return v
	}
	return svr.GetName()
}
