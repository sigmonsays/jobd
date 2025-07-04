package polling

import (
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shurcooL/go/osutil"
	"github.com/sigmonsays/jobd/job"
)

// populate git env
// see https://git-scm.com/book/en/v2/Git-Internals-Environment-Variables
func populateEnv(e []string, j *job.JobSpec) []string {
	cfg := j.Upstream.Git

	e = append(e, "GIT_CONFIG_NOSYSTEM=1")
	e = append(e, "GIT_PAGER=cat")

	if cfg.IdentityFile != "" {
		e = env_ssh_command(j, e, cfg.IdentityFile)
		return e
	}

	slog.Debug("populateEnv new env updated", "env", e)
	return e
}

func env_ssh_command(
	j *job.JobSpec,
	e []string,
	identityFile string,
) []string {

	sshbin, err := exec.LookPath("ssh")
	if err != nil {
		sshbin = "ssh"
	}
	ssh_opts := makeControlSocket(j, identityFile)
	if identityFile != "" {
		ssh_opts += " -i " + identityFile
	}
	ssh_command := sshbin + " " + strings.Trim(ssh_opts, " ")
	slog.Debug("setting GIT_SSH_COMMAND", "jid", j.JobId, "GIT_SSH_COMMAND", ssh_command)
	env := osutil.Environ(e)
	env.Set("GIT_SSH_COMMAND", ssh_command)
	return env
}

func makeControlSocket(j *job.JobSpec, ident string) string {
	// todo: Do a better job with the identity file
	b := filepath.Base(ident)
	// maybe IdentitiesOnly=yes  ?
	ret := " -oControlMaster=auto "
	ret += " -oControlPersist=yes "
	ret += " -oControlPath=/tmp/ssh-git-jobd-%u-%h-%n-%p-" + b
	return ret
}
