package agt

import (
	"ctl"
	"os"
	"regexp"
)

const (
	AppInfoFile   = "/home/app/info/app.json"
	NginxInfoFile = "/home/app/info/nginx.json"
	NginxUpDir    = "/usr/local/nginx/conf/upstream/"
	TomcatDir     = "/data/opt/tomcat/"
)

var (
	regNgxSrv, _     = regexp.Compile("^[^#][ \t]*server ")
	regPro, _        = regexp.Compile("^[0-9]+$")
	regTomcat, _     = regexp.Compile("/data1/opt/tomcat/[^/| ]+")
	regTomcatPort, _ = regexp.Compile(`port *= *"[0-9]+"`)
	regNum, _        = regexp.Compile(`[0-9]+`)
	regWar, _        = regexp.Compile("\\.war$")
)

func main() {
	//	ctl.Debug(regWar.FindString("test1-service-0.0.1-SNAPSHOT.war"))

	if false {
		r1 := regTomcatPort.FindString(`      <Connector port="20017" `)
		ctl.Debug(r1)
		r1 = regNum.FindString(r1)
		ctl.Debug(r1)
		ps := `app      17466  0.1 22.0 7461892 1808676 ?     Sl   May07   1:52 /usr/local/java/jdk/bin/java -Djava.util.logging.config.file=/data1/opt/tomcat/test1conf/logging.properties -Djava.util.logging.manager=org.apache.juli.ClassLoaderLogManager -server -Xms4096m -Xmx4096m -XX:MaxPermSize=512m -Djava.endorsed.dirs=/data1/opt/tomcat/test1/endorsed -classpath /data1/opt/tomcat/test1/bin/bootstrap.jar:/data1/opt/tomcat/test1/bin/tomcat-juli.jar -Dcatalina.base=/data1/opt/tomcat/test1 -Dcatalina.home=/data1/opt/tomcat/test1 -Djava.io.tmpdir=/data1/opt/tomcat/test1/temp org.apache.catalina.startup.Bootstrap start`
		for _, v := range regTomcat.FindAllString(ps, -1) {
			ctl.Debug(v)
		}
	}
	if len(os.Args) < 2 {
		ctl.FatalErr(ctl.Errorf("参数不正确"))

	}

	switch os.Args[1] {
	case "--service-get-info":
		serviceGetInfo()
	case "--get-app-info":
		AppGetInfo()
	case "--get-app-nginx-info":
		AppNginxGetInfo()
	case "--app-restart":
		AppRestart()
	case "--app-update":
		AppUpdate()
	case "--app-discovery":
		AppDiscovery()
	case "ngx-switch":
		AppSwitch()
	default:
		ctl.FatalErr(ctl.Errorf("参数不正确"))
	}
}
