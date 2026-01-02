package node

import (
	"codim/internal/executors/drivers/cmd"
	"codim/internal/executors/drivers/models"
	"codim/internal/utils/logger"
	"context"

	"github.com/google/uuid"
)

const (
	nsjailConfigTemplate = `name: "JOB-{{JOB_ID}}"
mode: ONCE
hostname: "JOB-{{JOB_ID}}"

clone_newns: true

clone_newpid: true
clone_newipc: true
clone_newuts: true
clone_newuser: true

clone_newnet: true
iface_no_lo: true

cwd: "/work"
mount_proc: false

mount {
  src: "/opt/nsjail/rootfs"
  dst: "/"
  is_bind: true
  rw: false
}

mount { src: "/jobs/{{JOB_ID}}" dst: "/work" is_bind: true rw: true }

mount { src: "/usr/bin/node" 		dst: "/usr/bin/node" 		is_bind: true rw: false }
mount { src: "/usr/bin/time" 		dst: "/usr/bin/time" 		is_bind: true rw: false }
mount { src: "/usr/lib"         	dst: "/usr/lib"         	is_bind: true rw: false }
mount { src: "/usr/lib/nodejs"  	dst: "/usr/lib/nodejs"  	is_bind: true rw: false }
mount { src: "/lib"             	dst: "/lib"             	is_bind: true rw: false }

mount { dst: "/tmp" fstype: "tmpfs" rw: true options: "size=128m" }

mount { dst: "/dev" fstype: "tmpfs" rw: false }
mount { src: "/dev/null"    dst: "/dev/null"    is_bind: true rw: true }
mount { src: "/dev/urandom" dst: "/dev/urandom" is_bind: true rw: false }

rlimit_as: 512
rlimit_cpu: 1
rlimit_nofile: 64
rlimit_nproc: 16
time_limit: 1

exec_bin {
  path: "/usr/bin/time"
  arg: "-f"
  arg: "%e,%U,%S,%M"
  arg: "/usr/bin/node"
  arg: "/work/{{ENTRY_POINT}}"
}
`
)

type Driver struct {
	logger    *logger.Logger
	cmdPrefix string
}

func New(cmdPrefix string, logger *logger.Logger) *Driver {
	return &Driver{
		logger:    logger,
		cmdPrefix: cmdPrefix,
	}
}

func (d *Driver) Execute(ctx context.Context, jobID uuid.UUID, src models.Directory, entryPoint string) (models.ExecuteResponse, error) {
	return cmd.Execute(ctx, d.cmdPrefix, nsjailConfigTemplate, jobID, src, entryPoint)
}

func (d *Driver) SetCmdPrefix(prefix string) error {
	d.cmdPrefix = prefix
	return nil
}

func (d *Driver) CmdPrefix() string {
	return d.cmdPrefix
}
