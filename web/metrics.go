package web

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/juho05/log"
)

type Metrics struct {
	lock sync.RWMutex

	startTime    time.Time
	successCount int
	visitCount   int

	failureTooLargeCount  int
	failureWrongFileCount int
	failureParseFormCount int
	failureFormatCount    int
	failureRateLimitCount int
	failureNoFilesCount   int
	failureOtherCount     int
}

func NewMetrics() *Metrics {
	return &Metrics{
		startTime: time.Now().UTC(),
	}
}

func (m *Metrics) Success() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.successCount++
}

func (m *Metrics) Visit() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.visitCount++
}

func (m *Metrics) FailureTooLarge() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.failureTooLargeCount++
}

func (m *Metrics) FailureWrongFile() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.failureWrongFileCount++
}

func (m *Metrics) FailureParseForm() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.failureParseFormCount++
}

func (m *Metrics) FailureFormat() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.failureFormatCount++
}

func (m *Metrics) FailureRateLimit() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.failureRateLimitCount++
}

func (m *Metrics) FailureNoFiles() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.failureNoFilesCount++
}

func (m *Metrics) FailureOther() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.failureOtherCount++
}

func (m *Metrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(map[string]map[string]any{
		"metrics": {
			"success":   m.successCount,
			"visits":    m.visitCount,
			"startTime": m.startTime,
			"failure": map[string]int{
				"total":            m.failureTooLargeCount + m.failureWrongFileCount + m.failureParseFormCount + m.failureFormatCount + m.failureOtherCount,
				"filesTooLarge":    m.failureTooLargeCount,
				"wrongFileFormat":  m.failureWrongFileCount,
				"parseFormError":   m.failureParseFormCount,
				"formattingFailed": m.failureFormatCount,
				"rateLimitReached": m.failureRateLimitCount,
				"noFilesUploaded":  m.failureNoFilesCount,
				"other":            m.failureOtherCount,
			},
		},
	})
	if err != nil {
		log.Errorf("failed to encode metrics: %s", err)
		return
	}
}
