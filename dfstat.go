package nux

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/toolkits/file"
	"io"
	"io/ioutil"
	"strings"
	"syscall"
)

var FSSPEC_IGNORE = map[string]struct{}{
	"none":  struct{}{},
	"nodev": struct{}{},
	"tmpfs": struct{}{},
}

var FSTYPE_IGNORE = map[string]struct{}{
	"cgroup":     struct{}{},
	"debugfs":    struct{}{},
	"devtmpfs":   struct{}{},
	"rpc_pipefs": struct{}{},
	"rootfs":     struct{}{},
}

var FSFILE_PREFIX_IGNORE = []string{
	"/dev",
	"/sys",
	"/net",
	"/misc",
	"/proc",
	"/lib",
}

func IgnoreFsFile(fs_file string) bool {
	for _, prefix := range FSFILE_PREFIX_IGNORE {
		if strings.HasPrefix(fs_file, prefix) {
			return true
		}
	}

	return false
}

// return: [][$fs_spec, $fs_file, $fs_vfstype]
func ListMountPoint() ([][3]string, error) {
	contents, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return nil, err
	}

	ret := make([][3]string, 0)

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		line, err := file.ReadLine(reader)
		if err == io.EOF {
			break
		}

		fields := strings.Fields(string(line))
		// Docs come from the fstab(5)
		// fs_spec     # Mounted block special device or remote filesystem e.g. /dev/sda1
		// fs_file     # Mount point e.g. /data
		// fs_vfstype  # File system type e.g. ext4
		// fs_mntops   # Mount options
		// fs_freq     # Dump(8) utility flags
		// fs_passno   # Order in which filesystem checks are done at reboot time

		fs_spec := fields[0]
		fs_file := fields[1]
		fs_vfstype := fields[2]

		if _, exist := FSSPEC_IGNORE[fs_spec]; exist {
			continue
		}

		if _, exist := FSTYPE_IGNORE[fs_vfstype]; exist {
			continue
		}

		if strings.HasPrefix(fs_vfstype, "fuse") {
			continue
		}

		if IgnoreFsFile(fs_file) {
			continue
		}

		// keep /dev/xxx device with shorter fs_file (remove mount binds)
		if strings.HasPrefix(fs_spec, "/dev") {
			deviceFound := false
			for idx := range ret {
				if ret[idx][0] == fs_spec {
					deviceFound = true
					if len(fs_file) < len(ret[idx][1]) {
						ret[idx][1] = fs_file
					}
					break
				}
			}
			if !deviceFound {
				ret = append(ret, [3]string{fs_spec, fs_file, fs_vfstype})
			}
		} else {
			ret = append(ret, [3]string{fs_spec, fs_file, fs_vfstype})
		}
	}
	return ret, nil
}

type DeviceUsage struct {
	FsSpec            string
	FsFile            string
	FsVfstype         string
	BlocksAll         uint64
	BlocksUsed        uint64
	BlocksFree        uint64
	BlocksUsedPercent float64
	BlocksFreePercent float64
	InodesAll         uint64
	InodesUsed        uint64
	InodesFree        uint64
	InodesUsedPercent float64
	InodesFreePercent float64
}

func (this *DeviceUsage) String() string {
	return fmt.Sprintf("<FsSpec:%s, FsFile:%s, FsVfstype:%s, BAll:%d, BUsed:%d, BFree:%d, BPUsed:%f, BPFree:%f, IAll:%d, IUsed:%d, IFree:%d, IPUsed:%f, IPFree:%f>",
		this.FsSpec,
		this.FsFile,
		this.FsVfstype,
		this.BlocksAll,
		this.BlocksUsed,
		this.BlocksFree,
		this.BlocksUsedPercent,
		this.BlocksFreePercent,
		this.InodesAll,
		this.InodesUsed,
		this.InodesFree,
		this.InodesUsedPercent,
		this.InodesFreePercent)
}

func BuildDeviceUsage(_fsSpec, _fsFile, _fsVfstype string) (*DeviceUsage, error) {
	ret := &DeviceUsage{FsSpec: _fsSpec, FsFile: _fsFile, FsVfstype: _fsVfstype}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(_fsFile, &fs)
	if err != nil {
		return nil, err
	}

	// blocks
	used := fs.Blocks - fs.Bfree
	ret.BlocksAll = uint64(fs.Frsize) * fs.Blocks
	ret.BlocksUsed = uint64(fs.Frsize) * used
	ret.BlocksFree = uint64(fs.Frsize) * fs.Bavail
	if fs.Blocks == 0 {
		ret.BlocksUsedPercent = 100.0
	} else {
		ret.BlocksUsedPercent = float64(used) * 100.0 / float64(used+fs.Bavail)
	}
	ret.BlocksFreePercent = 100.0 - ret.BlocksUsedPercent

	// inodes
	ret.InodesAll = fs.Files
	ret.InodesFree = fs.Ffree
	ret.InodesUsed = fs.Files - fs.Ffree
	if fs.Files == 0 {
		ret.InodesUsedPercent = 100.0
	} else {
		ret.InodesUsedPercent = float64(float64(ret.InodesUsed) * 100.0 / float64(ret.InodesAll))
	}
	ret.InodesFreePercent = 100.0 - ret.InodesUsedPercent

	return ret, nil
}
