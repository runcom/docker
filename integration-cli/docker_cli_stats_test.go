package main

import (
	"os/exec"
	"strings"
	"time"

	"github.com/go-check/check"
)

func (s *DockerSuite) TestStatsNoStream(c *check.C) {
	out, _, err := runCommandWithOutput(exec.Command(dockerBinary, "run", "-d", "busybox", "top"))
	c.Assert(err, check.IsNil)
	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), check.IsNil)

	out, _, err = runCommandWithOutput(exec.Command(dockerBinary, "stats", "--no-stream", id))
	c.Assert(err, check.IsNil)
	if !strings.Contains(out, id) {
		c.Fatalf("Expected output to contain %s, got instead: %s", id, out)
	}
}

func (*DockerSuite) TestStatsAllNoStream(c *check.C) {
	out, _, err := runCommandWithOutput(exec.Command(dockerBinary, "run", "-d", "busybox", "top"))
	c.Assert(err, check.IsNil)
	id1 := strings.TrimSpace(out)
	c.Assert(waitRun(id1), check.IsNil)

	out, _, err = runCommandWithOutput(exec.Command(dockerBinary, "run", "-d", "busybox", "top"))
	c.Assert(err, check.IsNil)
	id2 := strings.TrimSpace(out)
	c.Assert(waitRun(id2), check.IsNil)

	select {
	case err := <-chErr:
		if err != nil {
			c.Fatalf("Error running stats: %v", err)
		}
	case <-time.After(3 * time.Second):
		statsCmd.Process.Kill()
		c.Fatalf("stats did not return immediately when not streaming")
	}
}
