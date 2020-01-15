package dns

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	logger "github.com/cihub/seelog"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var ConfFileMd5 string

// 函数功能：添加DNS映射
// 输入		domain： mcloud的域名
// 			ip：	映射的relay ip 列表
// 返回： 	error（错误信息）
func DnsAdd(domainn string, ipList []string) (error) {
	var f *os.File
	var tmpFile string
	var err error

	if domainn == "" || len(ipList) == 0 {
		return nil
	}
	domain := strings.Split(domainn, ".")[0]

	if err := DnsDelByDomainAndIp(domain, ipList); err != nil {
		logger.Error("【dns】DnsDelByDomainAndIp err! err: %v\n", err)
		return err
	}

	if checkFileIsExist(cfg.GetDNS().NasDnsFileName) != true {
		logger.Errorf("【dns】Dns file not exits! filename: %s", cfg.GetDNS().NasDnsFileName)
		return errors.New("Dns file not exits!")
	}

	if tmpFile, err = makeTmpFile(cfg.GetDNS().NasDnsFileName); err != nil {   // 创建 tmp 文件， 在tmp文件中操作
		logger.Errorf("【dns】makeTmpFile error: %s", err.Error())
		return err
	}

	if f, err = os.OpenFile(tmpFile, os.O_APPEND | os.O_WRONLY, 0666); err != nil {   // 打开文件
		logger.Errorf("【dns】OpenFile error: %s", err.Error())
		return err
	}

	rulesBuffer := bytes.Buffer{}
	for _,v := range ipList {
		// 拼接域名解析规则
		rulesBuffer.WriteString(fmt.Sprintf("%s\tIN\tA\t%s\n", domain, v))
	}
	rules := rulesBuffer.String()

	if rules != "" {
		if _, err := io.WriteString(f, rules); err != nil { // 写入文件
			logger.Errorf("【dns】WriteString error: %s", err.Error())
			return err
		}
	}

	logger.Infof("【dns】add rules is : \n%s", rules)
	if err := MvFile(tmpFile, cfg.GetDNS().NasDnsFileName); err != nil{		// 移动tmp文件到主文件
		logger.Errorf("【dns】mvFile error: %s", err.Error())
		return err
	}
	return nil
}

// 函数功能：通过域名和ip删除dns映射
// 输入		domain： mcloud的域名
// 			ip：	映射的relay ip 列表
// 返回： 	error（错误信息）
func DnsDelByDomainAndIp(domainn string, ip []string) (error) {
	if domainn == "" || len(ip) == 0 {
		logger.Infof("【dns】domain or ip list is null!")
		return  nil
	}

	domain := strings.Split(domainn, ".")[0]
	for _, v := range ip {
		delParam := fmt.Sprintf("/^%s.*%s/d", domain, v)
		logger.Infof("【dns】del rule: %s", delParam)

		c := exec.Command("sed", "-i", delParam, cfg.GetDNS().NasDnsFileName)
		w := bytes.NewBuffer(nil)
		c.Stderr = w

		if err := c.Run(); err != nil {
			logger.Errorf("【dns】cmd error: %s", string(w.Bytes()))
			return errors.New(string(w.Bytes()))
		}
	}
	return  nil
}

// 函数功能：通过域名删除dns映射
// 输入		domain： mcloud的域名
// 返回： 	error（错误信息）
func DnsDelByDomain(domainn string) (error) {
	if domainn == "" {
		logger.Infof("【dns】domain is null!")
		return nil
	}

	domain := strings.Split(domainn, ".")[0]
	delParam := fmt.Sprintf("/^%s/d", domain)
	c := exec.Command("sed", "-i", delParam, cfg.GetDNS().NasDnsFileName)
	w := bytes.NewBuffer(nil)
	c.Stderr = w

	if err := c.Run(); err != nil {
		logger.Errorf("【dns】cmd error: %s", string(w.Bytes()))
		return errors.New(string(w.Bytes()))
	}
	return nil
}

// 函数功能：通过ip删除dns映射
// 输入		ip： relay的ip
// 返回： 	error（错误信息）
func DnsDelByIp(ip string) (error) {
	if ip == "" {
		logger.Infof("【dns】ip is null!")
		return nil
	}

	delParam := fmt.Sprintf("/%s/d", ip)
	c := exec.Command("sed", "-i", delParam, cfg.GetDNS().NasDnsFileName)
	w := bytes.NewBuffer(nil)
	c.Stderr = w

	if err := c.Run(); err != nil {
		logger.Errorf("【dns】cmd error: %s", string(w.Bytes()))
		return errors.New(string(w.Bytes()))
	}
	return nil
}

// 通过Ip删除域名映射，忽略传的域名
func DnsDelByIpIgnoreDomain(ip string, ignoreDomain string) (error){
	if ignoreDomain == "" || ip == "" {
		return nil
	}

	domain := strings.Split(ignoreDomain, ".")[0]
	//	删除所有 包含该relay IP的记录
	if err := DnsDelByIp(ip); err != nil {
		logger.Errorf("【dns】DnsDelByIpIgnoreDomain:DnsDelByIp err: %v\n",err)
		return err
	}

	// 添加上本身的域名映射
	if err := DnsAdd(domain, []string{ip}); err != nil {
		logger.Errorf("【dns】DnsDelByIpIgnoreDomain:DnsAdd failed! err: %v\n", err)
		return err
	}
	return nil
}

// 函数功能：DNS重新加载配置
// 输入		无
// 返回： 	error（错误信息）
func dnsReload()  (error){
	// service bind9 reload
	c := exec.Command("service", "bind9", "restart")
	w := bytes.NewBuffer(nil)
	c.Stderr = w

	if err := c.Run(); err != nil {
		logger.Errorf("【dns】cmd error: %s", string(w.Bytes()))
		return errors.New(string(w.Bytes()))
	}
	logger.Infof("【dns】 dns reload success!\n")
	return nil
}

// 函数功能：检查文件是否存在
// 输入		filename： 文件路径
// 返回： 	bool(检查结果)
func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		logger.Debugf("【dns】IsNotExist: %s", err.Error())
		return false
	}
	return true
}

// 函数功能：生成tmp临时文件
// 输入		fileName： 源文件路径
// 返回： 	string(tmp文件路径)  error（错误信息）
func makeTmpFile(fileName string) (string, error)  {
	tmpFileName := fileName + strconv.FormatInt(time.Now().UnixNano(), 10)

	c := exec.Command("cp", fileName, tmpFileName)
	w := bytes.NewBuffer(nil)
	c.Stderr = w

	if err := c.Run(); err != nil {
		logger.Errorf("【dns】cmd error: %s", string(w.Bytes()))
		return "", errors.New(string(w.Bytes()))
	}
	return tmpFileName, nil
}

// 函数功能：移动文件
// 输入		src： 源文件
// 			dst： 目标文件
// 返回： 	error（错误信息）
func MvFile(src string, dst string) (error)  {
	c := exec.Command("mv", src, dst)
	w := bytes.NewBuffer(nil)
	c.Stderr = w

	if err := c.Run(); err != nil {
		logger.Errorf("【dns】cmd error: %s", string(w.Bytes()))
		return errors.New(string(w.Bytes()))
	}
	return nil
}

// cp 文件
func CpFile(src string, dst string) (error)  {
	c := exec.Command("cp", src, dst)
	w := bytes.NewBuffer(nil)
	c.Stderr = w

	if err := c.Run(); err != nil {
		logger.Errorf("【dns】cmd error: %s", string(w.Bytes()))
		return errors.New(string(w.Bytes()))
	}
	return nil
}

// 检查配置文件文件是否变化
func confFileChange() bool {
	file, err := os.Open(cfg.GetDNS().NasDnsFileName)
	if err != nil {
		logger.Error("checkFileChange open dns conf file error! err: %v\n", err)
		return false
	}
	md5h := md5.New()
	io.Copy(md5h, file)
	md5Value := fmt.Sprintf("%x", md5h.Sum([]byte(""))) //md5
	if md5Value != ConfFileMd5 {
		ConfFileMd5 = md5Value
		return true
	}
	return false
}

// Dns 配置文件检查
func ConfCheck()  {
	if confFileChange() == true {		// 若配置文件变化， 重新加载配置
		dnsReload()
	}
}

