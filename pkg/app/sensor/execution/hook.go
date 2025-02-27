package execution

import (
	"context"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/slimtoolkit/slim/pkg/util/errutil"
)

type kind string

const (
	sensorPostStart     kind = "sensor-post-start"
	sensorPreShutdown   kind = "sensor-pre-shutdown"
	monitorPreStart     kind = "monitor-pre-start"
	targetAppRunning    kind = "target-app-running"
	monitorPostShutdown kind = "monitor-post-shutdown"
	monitorFailed       kind = "monitor-failed"
)

type hookExecutor struct {
	ctx      context.Context
	cmd      string
	lastHook kind
}

func (h *hookExecutor) State() string {
	return fmt.Sprintf("%s/%s", h.cmd, string(h.lastHook))
}

func (h *hookExecutor) HookSensorPostStart() {
	h.doHook(sensorPostStart)
}

func (h *hookExecutor) HookSensorPreShutdown() {
	h.doHook(sensorPreShutdown)
}

func (h *hookExecutor) HookMonitorPreStart() {
	h.doHook(monitorPreStart)
}

func (h *hookExecutor) HookTargetAppRunning() {
	h.doHook(targetAppRunning)
}

func (h *hookExecutor) HookMonitorPostShutdown() {
	h.doHook(monitorPostShutdown)
}

func (h *hookExecutor) HookMonitorFailed() {
	h.doHook(monitorFailed)
}

func (h *hookExecutor) doHook(k kind) {
	if len(h.cmd) == 0 {
		return
	}

	h.lastHook = k
	cmd := exec.CommandContext(h.ctx, h.cmd, string(k))
	out, err := cmd.CombinedOutput()

	logger := log.
		WithField("command", h.cmd).
		WithField("exit_code", cmd.ProcessState.ExitCode()).
		WithField("output", string(out))

	// Some lifecycle hooks are really fast - hence, the IsNoChildProcesses() check.
	if err == nil || errutil.IsNoChildProcesses(err) {
		logger.Debugf("sensor: %s hook succeeded", k)
	} else {
		logger.WithError(err).Warnf("sensor: %s hook failed", k)
	}
}
