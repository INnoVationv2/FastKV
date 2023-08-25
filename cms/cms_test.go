// Copyright 2021 hardcore-os Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cms

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

func Test_Basic_CRUD(t *testing.T) {
	for i := 0; i < 100; i++ {
		cms := NewCMS(100)
		str := randString()
		cms.Increment(str)
		if cms.Frequency(str) == 0 {
			log.Fatalf("Frequency Equal 0")
		}
	}

	fmt.Printf("[Basic_CRUD]Basic CRUD is working fine\n")
}

// 测试并发是否正常
func Test_Concurrent_CRUD(t *testing.T) {
	visitTimes := 10
	dataSz := int(1e6)

	stringSet := getRandomStringSet(dataSz)

	group := &sync.WaitGroup{}
	cms1 := NewCMS(dataSz)
	cms2 := NewCMS(dataSz)
	for k := range stringSet {
		for i := 0; i < visitTimes; i++ {
			cms1.Increment(k)

			group.Add(1)
			go func() {
				cms2.Increment(k)
				group.Done()
			}()
		}
	}

	group.Wait()

	sum1, sum2 := 0, 0
	for k := range stringSet {
		sum1 += cms1.Frequency(k)
		sum2 += cms1.Frequency(k)
	}

	if sum1 != sum2 {
		log.Fatalf("Not Equal!")
	}

	fmt.Printf("[Concurrent_CRUD]%d==%d, Concurrent Access is working fine\n", sum1, sum2)
}

// 评估误差率
// 创建一个大小为x的CMS
// 生成x个不同的字符串
// 每个字符串访问y次
//
//	CMS返回的访问频率必然大于等于y，
//	用CMS返回的结果-y，然后除以x，得到平均误差
func Test_Error_Rate(t *testing.T) {
	visitTimes := 10
	dataSz := int(1e6)

	stringSet := getRandomStringSet(dataSz)

	cms := NewCMS(dataSz)
	for k := range stringSet {
		for i := 0; i < visitTimes; i++ {
			cms.Increment(k)
		}
	}

	sum := 0
	for k := range stringSet {
		sum += cms.Frequency(k) - visitTimes
	}

	fmt.Printf("[Error_Rate]Average Error Distance： %f\n", float64(sum)/float64(dataSz))
}

// 测试Reset是否正常
func Test_Reset(t *testing.T) {
	visitTimes := 10
	dataSz := int(1e6)

	stringSet := getRandomStringSet(dataSz)

	cms := NewCMS(1000000)
	for k := range stringSet {
		for i := 0; i < visitTimes; i++ {
			cms.Increment(k)
		}
	}

	sum := uint64(0)
	for k := range stringSet {
		sum += uint64(cms.Frequency(k))
	}

	fmt.Printf("[Reset]Average Frequency: %f\n", float64(sum)/1000000.0)

	cms.Reset()

	sum = uint64(0)
	for k := range stringSet {
		sum += uint64(cms.Frequency(k))
	}

	fmt.Printf("[Reset]Average Frequency After Reset: %f\n", float64(sum)/1000000.0)

}
