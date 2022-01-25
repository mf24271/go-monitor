package monitor

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	"time"

	"github.com/shirou/gopsutil/process"
)

type monitor struct {
	config *config

	recordRing []recordInfo
	idx        int
}

func NewMonitor(config *config) *monitor {
	return &monitor{
		config:     config,
		recordRing: make([]recordInfo, 6),
	}
}

func (m *monitor) Start() error {
	// check config
	if m.config == nil {
		return errors.New("config is nil")
	}

	// check log path
	_, err := os.Stat(m.config.LogPath)
	if err != nil {
		err = os.Mkdir(m.config.LogPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	proce, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return err
	}
	go func() {
		for {
			cpuPercent, err := proce.Percent(5 * time.Second)
			if err != nil {
				continue
			}

			m.record(recordInfo{CPUPercent: float32(cpuPercent)})
		}
	}()
	return nil
}

func (m *monitor) GetLastRecordInfo() recordInfo {
	return m.recordRing[m.idx]
}

func (m *monitor) record(info recordInfo) {
	m.idx++
	if m.idx >= len(m.recordRing) {
		m.idx = 0
	}
	m.recordRing[m.idx].CPUPercent = info.CPUPercent

	cpuAverage := float32(0)

	for _, v := range m.recordRing {
		cpuAverage += v.CPUPercent
	}
	cpuAverage /= float32(len(m.recordRing))

	if info.CPUPercent > cpuAverage*1.3 && info.CPUPercent > 30 {
		fileName := fmt.Sprintf("%v-%v-%v.cpu.profile", time.Now().Format("2006-01-02_150405"), int(info.CPUPercent), int(cpuAverage))
		m.cpuProfile(path.Join(m.config.LogPath, fileName))
	}
}

func (m *monitor) cpuProfile(file string) {
	go func() {
		f, err := os.Create(file)
		if err != nil {
			fmt.Println(file, err)
			return
			// log.Fatal(err)
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			// StartCPUProfile failed, so no writes yet.
			f.Close()
			os.Remove(file)
			return
		}
		defer f.Close()
		time.Sleep(time.Duration(5) * time.Second)
		pprof.StopCPUProfile()
	}()
}
