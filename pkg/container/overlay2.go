package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type Overlay2Driver struct{}

func (overlay2 *Overlay2Driver) Name() string {
	return Overlay2
}

func (overlay2 *Overlay2Driver) Module() string {
	// Note: overlay2 module name is "overlay"
	return "overlay"
}

func (overlay2 *Overlay2Driver) MountRootfs(c *Container) error {
	workdir := path.Join(c.Rootfs.ContainerDir, DriverConfigs[Overlay2]["workDir"])

	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s",
		c.Rootfs.ImageDir, c.Rootfs.WriteDir, workdir)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", options, c.Rootfs.MergeDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to mount overlay2: %v", err)
	}
	return nil
}

func (overlay2 *Overlay2Driver) MountVolume(c *Container) error {
	for index, volume := range c.Volumes {
		volumeDir := path.Join(c.Rootfs.ContainerDir, "volumes",
			fmt.Sprintf("%03d", index+1))
		lowerDir := path.Join(volumeDir, "lower")
		workdir := path.Join(volumeDir, "work")

		for _, dir := range []string{lowerDir, workdir, volume.Source, volume.Target} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to mkdir %s: %v", dir, err)
			}
		}

		options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s",
			lowerDir, volume.Source, workdir)
		cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", options, volume.Target)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to mount local volume: %v", err)
		}
	}
	return nil
}
