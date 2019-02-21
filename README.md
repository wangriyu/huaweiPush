# hw_push

> Huawei Push

---

```

华为 push 接口升级, 导致原先的引用库无法继续使用, 故...

```


```

// 示例代码
package main

import (
	"fmt"
	huawei "github.com/stevelee-me/huaweiPush"
)

func main() {
	ClientId := "***"
	ClientSecret := "***"
	client := huawei.NewClient(ClientId, ClientSecret)

	// getToken
	// accessToken := client.GetToken()
	// fmt.Println("accessToken", accessToken)

	// push msg(会自己去 getToken 再请求推送)
	token := "***"
	payload := huawei.NewMessage().SetContent("huawei-content").SetTitle("huawei-title").Json()
	result := client.PushMsg(token, payload)
	fmt.Println("result", result)
}

```

---
