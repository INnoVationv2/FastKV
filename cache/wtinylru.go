package cache

import (
	"FastKV/cache/cms"
	"FastKV/cache/lru"
	"FastKV/cache/util"
	"fmt"
	"math"
	"sync"
)

type WindowTinyLfu struct {
	data        map[string]*lru.Node
	maximumSize int

	cms *cms.CMS

	window        *lru.Node
	windowMaxSize int
	windowSize    int
	windowLock    *sync.RWMutex

	probation        *lru.Node
	probationMaxSize int
	probationLock    *sync.RWMutex

	protected        *lru.Node
	protectedMaxSize int
	protectedSize    int
	protectedLock    *sync.RWMutex
}

func NewWindowTinyLfu(cap int) *WindowTinyLfu {
	wTinyLfu := &WindowTinyLfu{}
	wTinyLfu.data = make(map[string]*lru.Node)

	wTinyLfu.cms = cms.NewCMS(cap)
	wTinyLfu.maximumSize = cap

	wTinyLfu.windowMaxSize = int(math.Ceil(float64(wTinyLfu.maximumSize) * 0.01))
	wTinyLfu.probationMaxSize = int(math.Ceil(float64(wTinyLfu.maximumSize) * 0.2))
	wTinyLfu.protectedMaxSize = int(math.Ceil(float64(wTinyLfu.maximumSize) * 0.8))

	wTinyLfu.window = lru.NewEmptyNode()
	wTinyLfu.probation = lru.NewEmptyNode()
	wTinyLfu.protected = lru.NewEmptyNode()

	//wTinyLfu.windowLock = &sync.RWMutex{}
	//wTinyLfu.probationLock = &sync.RWMutex{}
	//wTinyLfu.protectedLock = &sync.RWMutex{}

	return wTinyLfu
}

func (w *WindowTinyLfu) record(key string, val interface{}) (*lru.Node, bool) {
	node, exist := w.data[key]
	if !exist {
		if val == nil {
			return nil, false
		}
		w.onMiss(key, val.(string))
		return nil, true
	} else if node.Type == lru.Window {
		w.onWindowHit(node)
	} else if node.Type == lru.Probation {
		w.onProbationHit(node)
	} else if node.Type == lru.Protected {
		w.onProtectedHit(node)
	}
	return node, true
}

func (w *WindowTinyLfu) onMiss(key, val string) {
	w.cms.Increment(key)

	node := lru.NewNode(key, val)
	node.AppendToFront(w.window)
	w.data[key] = node
	w.windowSize++
	w.evict()
}

func (w *WindowTinyLfu) onWindowHit(node *lru.Node) {
	w.cms.Increment(node.Key)

	node.MoveToFront(w.window)
}

func (w *WindowTinyLfu) onProbationHit(node *lru.Node) {
	w.cms.Increment(node.Key)

	node.Remove()
	node.Type = lru.Protected
	node.AppendToFront(w.protected)

	w.protectedSize++
	if w.protectedSize > w.protectedMaxSize {
		victim := w.protected.GetTail()
		victim.Remove()
		victim.Type = lru.Probation
		victim.AppendToFront(w.probation)
		w.protectedSize--
	}
}

func (w *WindowTinyLfu) onProtectedHit(node *lru.Node) {
	w.cms.Increment(node.Key)

	node.MoveToFront(w.protected)
}

func (w *WindowTinyLfu) evict() {
	if w.windowSize <= w.windowMaxSize {
		return
	}

	candidate := w.window.GetTail()
	w.windowSize--

	candidate.Remove()
	candidate.Type = lru.Probation
	candidate.AppendToFront(w.probation)

	if len(w.data) > w.maximumSize {
		victim := w.probation.GetTail()
		evict := w.estimate(candidate, victim)

		evict.Remove()
		delete(w.data, evict.Key)
		util.Debug("Pass %s\n", evict.Key)
	}
}

// candidate是从window来的新节点
// victim是从probation来的旧节点
// 倾向于保留新节点，去掉旧节点
func (w *WindowTinyLfu) estimate(candidate *lru.Node, victim *lru.Node) *lru.Node {
	candidateScore := w.cms.Frequency(candidate.Key)
	victimScore := w.cms.Frequency(victim.Key)

	if candidateScore < victimScore {
		return candidate
	}
	return victim
}

func (w *WindowTinyLfu) Set(key, val string) {
	w.record(key, val)
}

func (w *WindowTinyLfu) Get(key string) (interface{}, bool) {
	node, exist := w.record(key, nil)
	if !exist {
		return nil, false
	}
	return node.Value, true
}

func (w *WindowTinyLfu) String() string {
	str := fmt.Sprintf("%s | %s | %s", w.window, w.probation, w.protected)
	return str
}
