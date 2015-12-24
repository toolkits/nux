package nux

import (
    "bufio"
    "bytes"
    "fmt"
    "github.com/toolkits/file"
    "io"
    "io/ioutil"
    "strconv"
    "strings"
    "errors"
)

type Proc struct {
    Pid     int
    Name    string
    Cmdline string
    CpuUser float64
    CpuSys  float64
    CpuAll  float64
    Swap    int64
    Fd      int64
    Mem     int64
}

func (this *Proc) String() string {
    return fmt.Sprintf("<Pid:%d, Name:%s, Cmdline:%s, CpuUser:%f, CpuSys:%f, CpuAll:%f, Swap:%d, Fd:%d, Mem: %d>", this.Pid, this.Name, this.Cmdline, this.CpuUser, this.CpuSys, this.CpuAll, this.Swap, this.Fd, this.Mem)
}

func AllProcs() (ps []*Proc, err error) {
    var dirs []string
    dirs, err = file.DirsUnder("/proc")
    if err != nil {
        return
    }

    size := len(dirs)
    if size == 0 {
        return
    }

    for i := 0; i < size; i++ {
        pid, e := strconv.Atoi(dirs[i])
        if e != nil {
            continue
        }

        statFile := fmt.Sprintf("/proc/%d/stat", pid)
        statusFile := fmt.Sprintf("/proc/%d/status", pid)
        cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
        if !file.IsExist(statusFile) || !file.IsExist(cmdlineFile) || !file.IsExist(statFile){
            continue
        }

        cpu_user, cpu_sys, cpu_all, swap, e := ParseStatFile(statFile)
        if e != nil {
            continue
        }

        name, mem, fd, e := ParseStatusFile(statusFile)
        if e != io.EOF {
            continue
        }

        cmdlineBytes, e := file.ToBytes(cmdlineFile)
        if e != nil {
            continue
        }

        cmdlineBytesLen := len(cmdlineBytes)
        if cmdlineBytesLen == 0 {
            continue
        }

        noNut := make([]byte, 0, cmdlineBytesLen)

        for j := 0; j < cmdlineBytesLen; j++ {
            if cmdlineBytes[j] != 0 {
                noNut = append(noNut, cmdlineBytes[j])
            }
        }

        p := Proc{Pid: pid, Name: name, Cmdline: string(noNut), CpuUser: cpu_user, CpuSys: cpu_sys, CpuAll: cpu_all, Swap: swap, Fd: fd, Mem: mem*1024}
        ps = append(ps, &p)
    }

    return
}

func ParseStatFile(path string) (cpu_user float64, cpu_sys float64, cpu_all float64, swap int64, err error) {
    data, err := file.ToTrimString(path)
    if err != nil {
        return
    }

    stat := strings.Split(data, " ")
    if len(stat) < 44 {
        err = errors.New("file content too short")
        return
    }

    stat_13, err := strconv.ParseFloat(stat[13], 64)
    if err != nil {
        return
    }
    stat_14, err := strconv.ParseFloat(stat[14], 64)
    if err != nil {
        return
    }
    stat_15, err := strconv.ParseFloat(stat[15], 64)
    if err != nil {
        return
    }
    stat_16, err := strconv.ParseFloat(stat[16], 64)
    if err != nil {
        return
    }
    stat_35, err := strconv.ParseInt(stat[35], 10, 64)
    if err != nil {
        return
    }
    stat_36, err := strconv.ParseInt(stat[36], 10, 64)
    if err != nil {
        return
    }

    cpu_user = stat_13+ stat_15
    cpu_sys = stat_14 + stat_16
    cpu_all = cpu_user + cpu_sys
    swap = stat_35 + stat_36

    return
}

func ParseStatusFile(path string) (name string, mem int64, fd int64, err error) {
    var content []byte
    content, err = ioutil.ReadFile(path)
    if err != nil {
        return
    }

    reader := bufio.NewReader(bytes.NewBuffer(content))

    for {
        var bs []byte
        bs, err = file.ReadLine(reader)
        if err == io.EOF {
            return
        }

        line := string(bs)
        colonIndex := strings.Index(line, ":")

        if strings.TrimSpace(line[0:colonIndex]) == "Name" {
            name = strings.TrimSpace(line[colonIndex+1:])
        }
        if strings.TrimSpace(line[0:colonIndex]) == "VmRSS" {
            mem_str := strings.TrimSpace(line[colonIndex+1:])
            mem_vals := strings.Split(mem_str, " ")
            if len(mem_vals) == 0 {
                err = errors.New("File Content Missed")
                return
            }
            mem, err = strconv.ParseInt(mem_vals[0], 10, 64)
            if err != nil {
                return
            }
        }
        if strings.TrimSpace(line[0:colonIndex]) == "FDSize" {
            fdsize := strings.TrimSpace(line[colonIndex+1:])
            fd, err = strconv.ParseInt(fdsize, 10, 64)
            if err != nil {
                return
            }
        }

    }

    return
}
