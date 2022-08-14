package ctl

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/saintfish/chardet"

	//	"github.com/arstd/ecies"

	"github.com/axgle/mahonia"
	"golang.org/x/crypto/ssh"
)

func GbkToUtf(gbk string) string {
	var dec mahonia.Decoder
	dec = mahonia.NewDecoder("gbk")
	ret := dec.ConvertString(gbk)
	return ret
}
func GbkToUtfCheck(str string) (string, error) {
	var (
		char    *chardet.Detector
		charres *chardet.Result
		res     string
		err     error
	)
	res = str
	char = chardet.NewTextDetector()
	charres = nil
	charres, err = char.DetectBest([]byte(str))
	if err != nil {
		return res, Errorf("错误:解析更新内容编码失败,%s", err.Error())
	}
	if Log.Logger != nil {
		Log.Debug("char", charres.Charset)
	}
	Debug("char", charres.Charset)

	switch {
	case charres.Charset == "UTF-8":
		res = str
	case charres.Charset == "GB-18030":
		res = GbkToUtf(str)
	case RegGb.MatchString(charres.Charset):
		res = GbkToUtf(str)
	case RegIso.MatchString(charres.Charset):
		res = GbkToUtf(str)
	case charres.Charset == "Shift_JIS":
		res = GbkToUtf(str)
	case charres.Charset == "EUC-KR":
		res = GbkToUtf(str)
	case strings.Contains(charres.Charset, "KOI"):
		res = GbkToUtf(str)
	case strings.Contains(charres.Charset, "windows-12"):
		res = GbkToUtf(str)
	case strings.Contains(charres.Charset, "IBM420"):
		res = GbkToUtf(str)
	default:
		Debug(charres.Charset, str, GbkToUtf(str))
		if Log.Logger != nil {
			Log.Debug(charres.Charset, str, GbkToUtf(str))
		}

		return res, Errorf("错误:解析更新内容编码失败," + charres.Charset)

	}
	return res, nil
}
func UtfToGbk(utf string) string {
	var enc mahonia.Encoder
	enc = mahonia.NewEncoder("gbk")
	ret := enc.ConvertString(utf)
	return ret
}

func IntTo64(n int) (int64, error) {
	var (
		str string
		err error
		i64 int64
	)
	str = strconv.Itoa(n)
	i64, err = strconv.ParseInt(str, 10, 64)
	return i64, err
}

func Md5sum(path string) (string, error) {
	var (
		f       *os.File
		err     error
		md5Hash hash.Hash
	)
	f, err = os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5Hash = md5.New()
	_, err = io.Copy(md5Hash, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5Hash.Sum(nil)), nil
}

func Md5Data(d []byte) string {
	var (
		md5Hash hash.Hash
	)
	md5Hash = md5.New()
	md5Hash.Write(d)
	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

type AesEncrypt struct {
}

func (this *AesEncrypt) Encrypt(key, strMesg string) ([]byte, error) {
	aesKey := this.getKey(key)
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(strMesg))
	aesBlockEncrypter, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strMesg))
	return encrypted, nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(key string, src []byte) (strDesc string, err error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	aesKey := this.getKey(key)
	var iv = []byte(key)[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}

func PublicKeyFile(path string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func PublicKeyByte(buffer []byte) ssh.AuthMethod {
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func (this *AesEncrypt) getKey(strKey string) []byte {
	keyLen := len(strKey)
	if keyLen < 16 {
		return nil
	}
	arrKey := []byte(strKey)
	if keyLen >= 32 {
		return arrKey[:32]
	}
	if keyLen >= 24 {
		return arrKey[:24]
	}
	return arrKey[:16]
}

func PublicKeyByte1(buffer []byte) ssh.AuthMethod {
	//var (
	//	signer ssh.Signer
	//)
	key, err := ssh.ParsePublicKey(buffer)
	if err != nil {
		return nil
	}
	signer, err := ssh.NewSignerFromKey(key)

	if err != nil {
		return nil
	}
	signer.PublicKey()

	return ssh.PublicKeys(signer)
}

func AesDecryptBase(key, src string) (string, error) {
	var (
		baseStr    string
		arrEncrypt []byte
		aesEnc     AesEncrypt
		err        error
	)
	arrEncrypt, err = base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	baseStr, err = aesEnc.Decrypt(key, arrEncrypt)
	if err != nil {
		return "", err
	}
	return baseStr, nil
}

func AesEncryptBase(key, src string) (string, error) {
	var (
		err        error
		aesEnc     AesEncrypt
		arrEncrypt []byte
		baseStr    string
	)
	arrEncrypt, err = aesEnc.Encrypt(key, src)
	if err != nil {
		return "", err
	}
	baseStr = base64.StdEncoding.EncodeToString(arrEncrypt)
	return baseStr, nil

}

func ZipList(src string) ([]string, error) {
	var (
		err      error
		file     *zip.ReadCloser
		fileList []string
	)

	file, err = zip.OpenReader(src)
	if err != nil {
		return nil, err
	}

	for _, v := range file.File {
		fileList = append(fileList, v.Name)
	}
	file.Close()
	return fileList, nil
}

func RsaEncrypt(pubKey []byte, data interface{}) ([]byte, error) {
	var (
		newData []byte
		block   *pem.Block
		pub     *rsa.PublicKey
		err     error
	)
	block, _ = pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	pubInt, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub = pubInt.(*rsa.PublicKey)
	switch value := data.(type) {
	case string:
		newData = []byte(value)
	case []byte:
		newData = value
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub, newData)
}

func RsaDecrypt(privateKey []byte, data interface{}) ([]byte, error) {
	var (
		newData []byte
		block   *pem.Block
		priv    *rsa.PrivateKey
		err     error
	)
	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//	Debug(block)
	priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//	Debug(priv)
	switch value := data.(type) {
	case string:
		newData = []byte(value)
	case []byte:
		newData = value
	}
	//	Debug("start rsa")

	return rsa.DecryptPKCS1v15(rand.Reader, priv, newData)

}

func RsaB64Encrypt(pubKey []byte, data interface{}) (string, error) {
	var (
		d   []byte
		err error
	)
	d, err = RsaEncrypt(pubKey, data)
	if err != nil {
		return "", err
	}
	//	Debug(d)
	return base64.StdEncoding.EncodeToString(d), err

}
func RsaB64Decrypt(privateKey []byte, data string) ([]byte, error) {
	var (
		d, d1 []byte
		err   error
	)
	d, err = base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	//	Debug(d)
	d1, err = RsaDecrypt(privateKey, d)
	if err != nil {
		return nil, err
	}
	return d1, err
}
