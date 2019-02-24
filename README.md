# hwpush

> Huawei HMS SDK: https://developer.huawei.com/consumer/cn/service/hms/catalog/huaweipush_agent.html?page=hmssdk_huaweipush_api_reference_agent_s2

---

```

// 示例代码
package main

import (
	"fmt"
	huawei "github.com/wangriyu/huaweiPush"
)

func main() {
	ClientId := "***"
	ClientSecret := "***"
	AppPkgName := "***"
	client := huawei.NewClient(ClientId, ClientSecret, AppPkgName)

    deviceToken := "***"
    extra := struct {
        SessionID   int `json:"session_id"`
        SessionType int `json:"session_type"`
    }{123456, 789}

	payload := huawei.NewMessage().SetContent("huawei-content").SetTitle("huawei-title").SetAppPkgName(client.AppPkgName).SetCustomize([]map[string]interface{}{{"extra": extra}})
	result := client.PushMsg(deviceToken, payload.Json())
	fmt.Printf("result: %+v\n", result)
}
```

---
