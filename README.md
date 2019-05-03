# yipan

Demo: https://p.halu.lu/

## 使用说明

1. Fork 本仓库，再创建一个叫`yipan-config`的空仓库

2. Clone 你 fork 的仓库，在项目根目录下创建`config`文件， 内容如下，token获取方式参考[yitu](https://github.com/lzjluzijie/yitu#authorization)

    ```json
    {
        "ClientID": "4caae01e-515a-490f-bde7-92cff3b895ac",
        "ClientSecret": "qohmO45%%-jtxUVCAGP372{",
        "AccessToken": "AccessToken",
        "RefreshToken": "RefreshToken",
        "RedirectURI": "http://127.0.0.1:23333"
     }
    ```
    
3. 准备长度64的为十六进制密钥，可由[sha256](https://tools.halu.lu/#/hash)生成，将其设为你本地的环境变量

    ```bash
    # linux
    export hexkey=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 
    # windows
    set hexkey=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 
    ```

4. 运行`go run yipan.go enc`，然后将加密过的`config`添加到空仓库`yipan-config`并push

    如果你本地没有go环境也没事，把下面这段代码复制到https://play.golang.org/，将`hexKey`与`raw`替换为你自己的值就可以直接运行，系统会自动输出base64编码的加密文件
    
    ```go
    package main
    
    import (
    	"crypto/aes"
    	"crypto/cipher"
    	"crypto/rand"
    	"encoding/base64"
    	"encoding/hex"
    	"fmt"
    	"io"
    )
    
    var block cipher.Block
    
    func Encrypt(raw []byte) (enc []byte) {
    	enc = make([]byte, aes.BlockSize+len(raw))
    	iv := enc[:aes.BlockSize]
    	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    		panic(err)
    	}
    
    	stream := cipher.NewCFBEncrypter(block, iv)
    	stream.XORKeyStream(enc[aes.BlockSize:], raw)
    	return
    }
    
    func main() {
    	hexKey := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    	raw := []byte(`{
            "ClientID": "4caae01e-515a-490f-bde7-92cff3b895ac",
            "ClientSecret": "qohmO45%%-jtxUVCAGP372{",
            "AccessToken": "AccessToken",
            "RefreshToken": "RefreshToken",
            "RedirectURI": "http://127.0.0.1:23333"
    	}`)
    
    	key, err := hex.DecodeString(hexKey)
    	if err != nil {
    		panic(err)
    	}
    
    	block, err = aes.NewCipher(key)
    	if err != nil {
    		panic(err)
    	}
    
    	enc := Encrypt(raw)
    
    	fmt.Println(base64.StdEncoding.EncodeToString(enc))
    }

    ```
    
    然后打开浏览器，把`data:application/octet-stream;base64,上面输出的字符串`复制到地址栏并回车，就可以下载加密好的config了

5. 生成新的ssh密钥对，将公钥添加至`yipan-config`的`deploy keys`，并给予push权限

6. 在 Netlify 添加你fork的仓库，`Build command`为`bash deploy.sh`，`Publish directory`为`public`，记得将构建日志设为私有

7. 添加三个环境变量

    - hexkey: 64为的十六进制密钥
    - sshkey: ssh私钥，**你需要把换行替换为汉字`無`**，因为环境变量不支持换行
    - config: `yipan-config`的git仓库url，例如`git@github.com:lzjluzijie/yipan-config.git`

8. 理论上就可以用了
