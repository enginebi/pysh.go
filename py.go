package pysh

// go:generate export PYTHONPATH=.
import (
	"encoding/json"
	"errors"
	// "os"
	"os/exec"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/toukii/bytes"
	"github.com/toukii/goutils"
)

type resp [][]float64

var (
	module string
	PyHome string
)

func Init(module_ string) {
	module = module_
}

func GoPyFuncV2(funcname string, args [][]float64, params map[string]int32) ([][]float64, error) {
	log.Debugf("GoPyFuncV2, %s", funcname)

	pyargs, err1 := marshal(args)
	pyparams, err2 := marshal(params)
	// log.Infof("pyargs:%+v", pyargs)
	// log.Infof("pyparams:%+v", pyparams)
	log.Infof("%+v %+v", err1, err2)

	out, err := call(funcname, pyargs, pyparams)
	// log.Infof("out:%+v", out)
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}

	var r resp
	err = json.Unmarshal([]byte(out), &r)
	if err != nil {
		log.Errorf("err:%+v", err)
		return nil, err
	}
	// log.Infof("resp:%+v", r)
	return [][]float64(r), nil
}

func call(funcname, input, params string) (string, error) {
	fmtargs := `import {{.Module}}; 
params = {{.Params}}
I = {{.Input}}
resp = {{.Module}}.{{.Funcname}}(*I,**params)
print("[resp]: %s" % (resp))`

	at, err := template.New("callpy").Parse(fmtargs)
	if err != nil {
		log.Errorf("%+v", err)
		return "", err
	}
	data := map[string]string{
		"Module":   module,
		"Funcname": funcname,
		"Input":    input,
		"Params":   params,
	}
	wr := bytes.NewWriter(make([]byte, 0, 1024))
	if err := at.Execute(wr, data); err != nil {
		log.Errorf("%+v", err)
		return "", err
	}
	args := goutils.ToString(wr.Bytes())
	// log.Infof("%s", args)

	// out, err := exc.Bash(args).Do()
	out, err := execute("python", "-c", args)

	outstr := goutils.ToString(out)
	// log.Infof("outstr:%+v", outstr)

	ss := strings.Split(outstr, "[resp]: ")
	if len(ss) < 2 {
		return "", errors.New(outstr)
	}
	return ss[1], nil
}

func GoPyFuncV4(funcname string, args interface{}, params string) (string, error) {
	log.Debugf("GoPyFuncV3, %s", funcname)

	pyargs, err1 := marshal(args)
	pyparams, err2 := marshal(params)
	// log.Infof("pyargs:%+v", pyargs)
	// log.Infof("pyparams:%+v", pyparams)
	log.Infof("%+v %+v", err1, err2)

	out, err := calltrain(funcname, pyargs, pyparams)
	// log.Infof("out:%+v", out)
	return out, err
}

func GoPyFuncV3(funcname string, args [][]float64, params string) (string, error) {
	log.Debugf("GoPyFuncV3, %s", funcname)

	pyargs, err1 := marshal(args)
	pyparams, err2 := marshal(params)
	// log.Infof("pyargs:%+v", pyargs)
	// log.Infof("pyparams:%+v", pyparams)
	log.Infof("%+v %+v", err1, err2)

	out, err := calltrain(funcname, pyargs, pyparams)
	// log.Infof("out:%+v", out)
	return out, err
}

func calltrain(funcname, input, params string) (string, error) {
	fmtargs := `import {{.Module}};
params = {{.Params}}
I = {{.Input}}
resp = {{.Module}}.{{.Funcname}}(I, params)
print("[resp]: %s" % (resp))`

	at, err := template.New("callpy").Parse(fmtargs)
	if err != nil {
		log.Errorf("%+v", err)
		return "", err
	}
	data := map[string]string{
		"Module":   module,
		"Funcname": funcname,
		"Input":    input,
		"Params":   params,
	}
	wr := bytes.NewWriter(make([]byte, 0, 1024))
	if err := at.Execute(wr, data); err != nil {
		log.Errorf("%+v", err)
		return "", err
	}
	args := goutils.ToString(wr.Bytes())
	// log.Infof("%s", args)

	// out, err := exc.Bash(args).Do()
	out, err := execute("python", "-c", args)
	// log.Infof("%+v", args)

	outstr := goutils.ToString(out)
	// log.Infof("outstr:%+v", outstr)

	ss := strings.Split(outstr, "[resp]: ")
	if len(ss) < 2 {
		return "", errors.New(outstr)
	}
	return ss[1], nil
}

func execute(executer, c, args string) ([]byte, error) {
	cmd := exec.Command(executer, c, args)
	// cmd.Dir = os.Getenv("PYTHONPATH")
	// log.Infof("PYTHONPATH:%s", cmd.Dir)
	if cmd.Dir == "" {
		cmd.Dir = PyHome
	}
	// log.Infof("PYTHONPATH final: %s", cmd.Dir)
	return cmd.CombinedOutput()
}

func GoPyFunc(funcname string, args ...float64) []float64 {
	return nil
}

func marshal(i interface{}) (string, error) {
	bs, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return goutils.ToString(bs), nil
}
