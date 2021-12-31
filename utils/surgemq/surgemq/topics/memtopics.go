// Copyright (c) 2014 The SurgeMQ Authors. All rights reserved.
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

package topics

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"

	"git.iflytek.com/HY_XIoT/core/utils/log"
	"git.iflytek.com/HY_XIoT/core/xfyun/surgemq/message"
)

var (
	// MaxQosAllowed is the maximum QOS supported by this server
	MaxQosAllowed = message.QosExactlyOnce
)

var _ TopicsProvider = (*memTopics)(nil)

type memTopics struct {
	// Sub/unsub mutex
	smu sync.RWMutex
	// Subscription tree
	sroot *snode

	// Retained message mutex
	rmu sync.RWMutex
	// Retained messages topic tree
	rroot *rnode
}

func init() {
	Register("mem", NewMemProvider())
}

var _ TopicsProvider = (*memTopics)(nil)

// NewMemProvider returns an new instance of the memTopics, which is implements the
// TopicsProvider interface. memProvider is a hidden struct that stores the topic
// subscriptions and retained messages in memory. The content is not persistend so
// when the server goes, everything will be gone. Use with care.
func NewMemProvider() *memTopics {
	return &memTopics{
		sroot: newSNode(),
		rroot: newRNode(),
	}
}

func (this *memTopics) Subscribe(topic []byte, qos byte, sub interface{}) (byte, error) {
	if !message.ValidQos(qos) {
		return message.QosFailure, fmt.Errorf("Invalid QoS %d", qos)
	}

	if sub == nil {
		return message.QosFailure, fmt.Errorf("Subscriber cannot be nil")
	}

	this.smu.Lock()
	defer this.smu.Unlock()

	if qos > MaxQosAllowed {
		qos = MaxQosAllowed
	}

	// 解析出group Name
	var group string
	log.Debugf("memTopics.Subscribe | topic:%s", string(topic))
	group, topic = getGroupName(topic)
	log.Debugf("memTopics.Subscribe | group:%s topic:%s", group, string(topic))
	if err := this.sroot.sinsert(group, topic, qos, sub); err != nil {
		return message.QosFailure, err
	}

	return qos, nil
}

func (this *memTopics) Unsubscribe(topic []byte, sub interface{}) error {
	this.smu.Lock()
	defer this.smu.Unlock()

	// 解析出group Name
	var group string
	group, topic = getGroupName(topic)
	return this.sroot.sremove(group, topic, sub)
}

// Returned values will be invalidated by the next Subscribers call
func (this *memTopics) Subscribers(topic []byte, qos byte, subs *[]interface{}, qoss *[]byte) error {
	log.Debugf("memtopics.Subscribers | topic:%s", string(topic))
	if !message.ValidQos(qos) {
		return fmt.Errorf("Invalid QoS %d", qos)
	}

	this.smu.RLock()
	defer this.smu.RUnlock()

	*subs = (*subs)[0:0]
	*qoss = (*qoss)[0:0]

	return this.sroot.smatch(topic, qos, subs, qoss)
}

func (this *memTopics) Retain(msg *message.PublishMessage) error {
	this.rmu.Lock()
	defer this.rmu.Unlock()

	// So apparently, at least according to the MQTT Conformance/Interoperability
	// Testing, that a payload of 0 means delete the retain message.
	// https://eclipse.org/paho/clients/testing/
	if len(msg.Payload()) == 0 {
		return this.rroot.rremove(msg.Topic())
	}

	return this.rroot.rinsert(msg.Topic(), msg)
}

func (this *memTopics) Retained(topic []byte, msgs *[]*message.PublishMessage) error {
	this.rmu.RLock()
	defer this.rmu.RUnlock()

	return this.rroot.rmatch(topic, msgs)
}

func (this *memTopics) Close() error {
	this.sroot = nil
	this.rroot = nil
	return nil
}

// subscrition nodes
type snode struct {
	// If this is the end of the topic string, then add subscribers here
	subs []interface{}
	// 按组订阅模式
	groupSubs map[string][]interface{}
	qos       []byte

	// Otherwise add the next topic level here
	snodes map[string]*snode
}

func newSNode() *snode {
	return &snode{
		groupSubs: make(map[string][]interface{}),
		snodes:    make(map[string]*snode),
	}
}

func (this *snode) sinsert(group string, topic []byte, qos byte, sub interface{}) error {

	// If there's no more topic levels, that means we are at the matching snode
	// to insert the subscriber. So let's see if there's such subscriber,
	// if so, update it. Otherwise insert it.
	if len(topic) == 0 {
		// Let's see if the subscriber is already on the list. If yes, update
		// QoS and then return.
		if group == "" {
			for i := range this.subs {
				if equal(this.subs[i], sub) {
					this.qos[i] = qos
					return nil
				}
			}

			// Otherwise add.
			this.subs = append(this.subs, sub)
			this.qos = append(this.qos, qos)
		} else {
			// 添加到指定group中
			this.groupSubs[group] = append(this.groupSubs[group], sub)
		}

		return nil
	}

	// Not the last level, so let's find or create the next level snode, and
	// recursively call it's insert().

	// ntl = next topic level
	ntl, rem, err := nextTopicLevel(topic)
	if err != nil {
		return err
	}

	level := string(ntl)

	// Add snode if it doesn't already exist
	n, ok := this.snodes[level]
	if !ok {
		n = newSNode()
		this.snodes[level] = n
	}

	return n.sinsert(group, rem, qos, sub)
}

// This remove implementation ignores the QoS, as long as the subscriber
// matches then it's removed
func (this *snode) sremove(group string, topic []byte, sub interface{}) error {

	//log.Infof("group:%s topic:%s sub:%v\n", group, string(topic), sub)
	// If the topic is empty, it means we are at the final matching snode. If so,
	// let's find the matching subscribers and remove them.
	if len(topic) == 0 {
		// If subscriber == nil, then it's signal to remove ALL subscribers
		if sub == nil {
			log.Infof("clear all subscribers\n")
			this.subs = this.subs[0:0]
			this.qos = this.qos[0:0]
			// 清空map
			for k := range this.groupSubs {
				delete(this.groupSubs, k)
			}
			return nil
		}

		// If we find the subscriber then remove it from the list. Technically
		// we just overwrite the slot by shifting all other items up by one.
		if group == "" {
			for i := range this.subs {
				if equal(this.subs[i], sub) {
					this.subs = append(this.subs[:i], this.subs[i+1:]...)
					this.qos = append(this.qos[:i], this.qos[i+1:]...)
					return nil
				}
			}
		} else {
			gsubs, ok := this.groupSubs[group]
			if ok {
				for i := range gsubs {
					if equal(gsubs[i], sub) {
						temp := append(gsubs[:i], gsubs[i+1:]...)
						if len(temp) > 0 {
							this.groupSubs[group] = temp
						} else {
							// 已无订阅者，删除map中该级key
							delete(this.groupSubs, group)
						}
						return nil
					}
				}
			}
		}

		return fmt.Errorf("memtopics/remove: No topic found for subscriber")
	}

	// Not the last level, so let's find the next level snode, and recursively
	// call it's remove().

	// ntl = next topic level
	ntl, rem, err := nextTopicLevel(topic)
	if err != nil {
		return err
	}

	level := string(ntl)

	// Find the snode that matches the topic level
	n, ok := this.snodes[level]
	if !ok {
		return fmt.Errorf("memtopics/remove: No topic found")
	}

	// Remove the subscriber from the next level snode
	if err := n.sremove(group, rem, sub); err != nil {
		return err
	}

	// If there are no more subscribers and snodes to the next level we just visited
	// let's remove it
	if len(n.groupSubs) == 0 && len(n.subs) == 0 && len(n.snodes) == 0 {
		delete(this.snodes, level)
	}

	return nil
}

// smatch() returns all the subscribers that are subscribed to the topic. Given a topic
// with no wildcards (publish topic), it returns a list of subscribers that subscribes
// to the topic. For each of the level names, it's a match
// - if there are subscribers to '#', then all the subscribers are added to result set
func (this *snode) smatch(topic []byte, qos byte, subs *[]interface{}, qoss *[]byte) error {
	// If the topic is empty, it means we are at the final matching snode. If so,
	// let's find the subscribers that match the qos and append them to the list.
	if len(topic) == 0 {
		this.matchQos(qos, subs, qoss)
		return nil
	}

	// ntl = next topic level
	ntl, rem, err := nextTopicLevel(topic)
	if err != nil {
		return err
	}

	level := string(ntl)
	log.Debugf("snode.smatch | level:%s, rem:%s", level, string(rem))

	for k, n := range this.snodes {
		// If the key is "#", then these subscribers are added to the result set
		if k == MWC {
			n.matchQos(qos, subs, qoss)
		} else if k == SWC || k == level {
			if err := n.smatch(rem, qos, subs, qoss); err != nil {
				return err
			}
		}
	}

	return nil
}

// retained message nodes
type rnode struct {
	// If this is the end of the topic string, then add retained messages here
	msg *message.PublishMessage
	buf []byte

	// Otherwise add the next topic level here
	rnodes map[string]*rnode
}

func newRNode() *rnode {
	return &rnode{
		rnodes: make(map[string]*rnode),
	}
}

func (this *rnode) rinsert(topic []byte, msg *message.PublishMessage) error {
	// If there's no more topic levels, that means we are at the matching rnode.
	if len(topic) == 0 {
		l := msg.Len()

		// Let's reuse the buffer if there's enough space
		if l > cap(this.buf) {
			this.buf = make([]byte, l)
		} else {
			this.buf = this.buf[0:l]
		}

		if _, err := msg.Encode(this.buf); err != nil {
			return err
		}

		// Reuse the message if possible
		if this.msg == nil {
			this.msg = message.NewPublishMessage()
		}

		if _, err := this.msg.Decode(this.buf); err != nil {
			return err
		}

		return nil
	}

	// Not the last level, so let's find or create the next level snode, and
	// recursively call it's insert().

	// ntl = next topic level
	ntl, rem, err := nextTopicLevel(topic)
	if err != nil {
		return err
	}

	level := string(ntl)

	// Add snode if it doesn't already exist
	n, ok := this.rnodes[level]
	if !ok {
		n = newRNode()
		this.rnodes[level] = n
	}

	return n.rinsert(rem, msg)
}

// Remove the retained message for the supplied topic
func (this *rnode) rremove(topic []byte) error {
	// If the topic is empty, it means we are at the final matching rnode. If so,
	// let's remove the buffer and message.
	if len(topic) == 0 {
		this.buf = nil
		this.msg = nil
		return nil
	}

	// Not the last level, so let's find the next level rnode, and recursively
	// call it's remove().

	// ntl = next topic level
	ntl, rem, err := nextTopicLevel(topic)
	if err != nil {
		return err
	}

	level := string(ntl)

	// Find the rnode that matches the topic level
	n, ok := this.rnodes[level]
	if !ok {
		return fmt.Errorf("memtopics/rremove: No topic found")
	}

	// Remove the subscriber from the next level rnode
	if err := n.rremove(rem); err != nil {
		return err
	}

	// If there are no more rnodes to the next level we just visited let's remove it
	if len(n.rnodes) == 0 {
		delete(this.rnodes, level)
	}

	return nil
}

// rmatch() finds the retained messages for the topic and qos provided. It's somewhat
// of a reverse match compare to match() since the supplied topic can contain
// wildcards, whereas the retained message topic is a full (no wildcard) topic.
func (this *rnode) rmatch(topic []byte, msgs *[]*message.PublishMessage) error {
	// If the topic is empty, it means we are at the final matching rnode. If so,
	// add the retained msg to the list.
	if len(topic) == 0 {
		if this.msg != nil {
			*msgs = append(*msgs, this.msg)
		}
		return nil
	}

	// ntl = next topic level
	ntl, rem, err := nextTopicLevel(topic)
	if err != nil {
		return err
	}

	level := string(ntl)

	if level == MWC {
		// If '#', add all retained messages starting this node
		this.allRetained(msgs)
	} else if level == SWC {
		// If '+', check all nodes at this level. Next levels must be matched.
		for _, n := range this.rnodes {
			if err := n.rmatch(rem, msgs); err != nil {
				return err
			}
		}
	} else {
		// Otherwise, find the matching node, go to the next level
		if n, ok := this.rnodes[level]; ok {
			if err := n.rmatch(rem, msgs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *rnode) allRetained(msgs *[]*message.PublishMessage) {
	if this.msg != nil {
		*msgs = append(*msgs, this.msg)
	}

	for _, n := range this.rnodes {
		n.allRetained(msgs)
	}
}

const (
	stateCHR byte = iota // Regular character
	stateMWC             // Multi-level wildcard
	stateSWC             // Single-level wildcard
	stateSEP             // Topic level separator
	stateSYS             // System level topic ($)
)

// Returns topic level, remaining topic levels and any errors
func nextTopicLevel(topic []byte) ([]byte, []byte, error) {
	s := stateCHR

	for i, c := range topic {
		switch c {
		case '/':
			if s == stateMWC {
				return nil, nil, fmt.Errorf("memtopics/nextTopicLevel: Multi-level wildcard found in topic and it's not at the last level")
			}

			if i == 0 {
				return []byte(SWC), topic[i+1:], nil
			}

			return topic[:i], topic[i+1:], nil

		case '#':
			if i != 0 {
				return nil, nil, fmt.Errorf("memtopics/nextTopicLevel: Wildcard character '#' must occupy entire topic level")
			}

			s = stateMWC

		case '+':
			if i != 0 {
				return nil, nil, fmt.Errorf("memtopics/nextTopicLevel: Wildcard character '+' must occupy entire topic level")
			}

			s = stateSWC

		case '$':
			if i == 0 {
				return nil, nil, fmt.Errorf("memtopics/nextTopicLevel: Cannot publish to $ topics")
			}

			s = stateSYS

		default:
			if s == stateMWC || s == stateSWC {
				return nil, nil, fmt.Errorf("memtopics/nextTopicLevel: Wildcard characters '#' and '+' must occupy entire topic level")
			}

			s = stateCHR
		}
	}

	// If we got here that means we didn't hit the separator along the way, so the
	// topic is either empty, or does not contain a separator. Either way, we return
	// the full topic
	return topic, nil, nil
}

// The QoS of the payload messages sent in response to a subscription must be the
// minimum of the QoS of the originally published message (in this case, it's the
// qos parameter) and the maximum QoS granted by the server (in this case, it's
// the QoS in the topic tree).
//
// It's also possible that even if the topic matches, the subscriber is not included
// due to the QoS granted is lower than the published message QoS. For example,
// if the client is granted only QoS 0, and the publish message is QoS 1, then this
// client is not to be send the published message.
func (this *snode) matchQos(qos byte, subs *[]interface{}, qoss *[]byte) {

	// 非group的subs全取出
	for i, sub := range this.subs {
		// If the published QoS is higher than the subscriber QoS, then we skip the
		// subscriber. Otherwise, add to the list.

		// 改为标准协议，以sub和pub中qos小的qos下发消息
		//  change by haojin 2018-5-14
		if qos <= this.qos[i] {
			*subs = append(*subs, sub)
			*qoss = append(*qoss, qos)
		} else {
			*subs = append(*subs, sub)
			*qoss = append(*qoss, this.qos[i])
		}

	}
	// group中的sub，每个key中对应的list中取一个
	// 先忽略qos，基于group订阅
	for k, gsubs := range this.groupSubs {
		log.Debugf("snode.matchQos | key:%s", k)
		*qoss = append(*qoss, 1)
		if len(gsubs) == 1 {
			*subs = append(*subs, gsubs[0])
			continue
		}
		idx := rand.Intn(len(gsubs))
		*subs = append(*subs, gsubs[idx])
	}
}

func equal(k1, k2 interface{}) bool {
	if reflect.TypeOf(k1) != reflect.TypeOf(k2) {
		return false
	}

	if reflect.ValueOf(k1).Kind() == reflect.Func {
		return &k1 == &k2
	}

	if k1 == k2 {
		return true
	}

	switch k1 := k1.(type) {
	case string:
		return k1 == k2.(string)

	case int64:
		return k1 == k2.(int64)

	case int32:
		return k1 == k2.(int32)

	case int16:
		return k1 == k2.(int16)

	case int8:
		return k1 == k2.(int8)

	case int:
		return k1 == k2.(int)

	case float32:
		return k1 == k2.(float32)

	case float64:
		return k1 == k2.(float64)

	case uint:
		return k1 == k2.(uint)

	case uint8:
		return k1 == k2.(uint8)

	case uint16:
		return k1 == k2.(uint16)

	case uint32:
		return k1 == k2.(uint32)

	case uint64:
		return k1 == k2.(uint64)

	case uintptr:
		return k1 == k2.(uintptr)
	}

	return false
}

// 根据topic前缀来确定是否是基于group模式的topic
// group格式 $group/groupname/topic
var GroupPrefix = "$group/"

func getGroupName(topic []byte) (groupName string, subTopic []byte) {

	// 非$group开头，不是基于group订阅模式
	if !strings.HasPrefix(string(topic), GroupPrefix) {
		return "", topic
	}

	// 去除group前缀
	subTopic = topic[len(GroupPrefix):]
	// 第一次出现/位置
	idx := strings.Index(string(subTopic), "/")
	if idx < 0 {
		return "", topic
	}

	return string(subTopic[:idx]), []byte(subTopic[idx:])
}