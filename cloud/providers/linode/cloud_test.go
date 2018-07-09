package linode

import (
	"testing"

	"github.com/taoh/linodego"

	//"fmt"
	"encoding/json"
	"fmt"
)

func TestInstances(t *testing.T) {
	//linodeClient := linodego.NewClient("", nil)
	//fmt.Println(linodes, err := client.Linode.List(-1))
}

func TestJS(t *testing.T) {
	cl := `[{"TOTALXFER":3000,"ALERT_BWQUOTA_ENABLED":1,"ALERT_DISKIO_ENABLED":1,"DISTRIBUTIONVENDOR":"Ubuntu","ALERT_BWOUT_ENABLED":1,"ALERT_CPU_THRESHOLD":90,"LINODEID":4188216,"ALERT_BWOUT_THRESHOLD":10,"BACKUPWINDOW":0,"DATACENTERID":3,"ALERT_BWIN_ENABLED":1,"STATUS":1,"PLANID":3,"LABEL":"lin3-045-079-066-185","ALERT_BWIN_THRESHOLD":10,"ALERT_CPU_ENABLED":1,"BACKUPSENABLED":0,"TOTALRAM":4096,"WATCHDOG":1,"CREATE_DT":"2017-11-01 04:53:39.0","ISKVM":1,"ALERT_BWQUOTA_THRESHOLD":80,"BACKUPWEEKLYDAY":0,"TOTALHD":49152,"LPM_DISPLAYGROUP":"","ALERT_DISKIO_THRESHOLD":10000,"ISXEN":0}]`
	tt := make([]linodego.Linode, 0)
	err := json.Unmarshal([]byte(cl), &tt)
	fmt.Println(err, tt)
}
