# AsiaYo-test

資料庫考題
===
## 題目一

```sql =
SELECT 
    bnbs.id AS bnb_id,
    bnbs.name AS bnb_name,
    SUM(order.amount) AS may_amount
FROM
    orders
JOIN
    bnbs ON orders.bnb_id = bnbs.id
WHERE
    orders.currency = 'TWD'
    AND orders.created_at >= '2023-05-01'
    AND orders.created_at < '2023-06-01'
GROUP BY
    bnbs.id,
    bnbs.name
ORDER BY
    may_amount DESC
LIMIT 10;
```

### 題目二
首先查看執行計劃，如上面的 SQL 語句（使用MySQL）：
```sql = 
EXPLAIN SELECT 
    bnbs.id AS bnb_id,
    bnbs.name AS bnb_name,
    SUM(order.amount) AS may_amount
FROM
    orders
JOIN
    bnbs ON orders.bnb_id = bnbs.id
WHERE
    orders.currency = 'TWD'
    AND orders.created_at >= '2023-05-01'
    AND orders.created_at < '2023-06-01'
GROUP BY
    bnbs.id,
    bnbs.name
ORDER BY
    may_amount DESC
LIMIT 10;
```
透過執行計劃可以看到下面的結果（資料為臨時資料）：
|id |select_type |table |partitions |type |possible_keys |key |key_len |ref |rows |filtered |Extra|
|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|
|1|SIMPLE|orders|(Null)|ALL|order_bnb_id|(Null)|(Null)|(Null)|3|33.33|Using where; Using temp|
|1|SIMPLE|bnbs|(Null)|eq_ref|PRIMARY|PRIMARY|4|asiayo.orders.bnb_id|1|100.00|(Null)|

透過第一條執行計劃可以發現 type 為 ALL 代表全表掃描，並且 filtered 指出資料庫找到的有效資料只佔掃描行數的 33.33%，可能是效率出現問題的點，透過在兩個 where 表達句用到的欄位 currency 和 created_at 添加索引可以避免全表掃描，下面是添加索引過後的執行計劃：
|id |select_type |table |partitions |type |possible_keys |key |key_len |ref |rows |filtered |Extra|
|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|
|1|SIMPLE|orders|(Null)|ref|order_bnb_id,currency,created_at|currency|13|const|3|100.00|Using index condition;Using where;Using temporary; Using filesort|
|1|SIMPLE|bnbs|(Null)|eq_ref|PRIMARY|PRIMARY|4|asiayo.orders.bnb_id|1|100.00|(Null)|

透過以上方法，成功將查詢數據的有效數據量提升，避免查詢過多的無效資料。
另外，SQL 語句可以遵從下面幾個原則來實現高效率的查詢：
1. 避免使用 SELECT *，只查詢需要的欄位可以減少數據傳輸量
2. 使用 JOIN 代替子查詢
3. 使用 Where 條件時盡量提供精準條件，避免使用 ！= 或 <>，可以使用 IN 才不會導致索引失效
4. 列上不使用聚合函數，這會導致列資料失去有序，從而使索引失效
5. 如果數據量太大，可以使用 Limit 來限制輸出的資料，並且使用分頁查詢查看後面的數據

如果各項優化手段都已經實作但查詢速度仍然很慢，可以考慮分庫分表將資料分散，或者讀寫分離來加快讀與寫的速度。如果內部資料不常變動，也可以考慮使用 Redis 快取資料來提高整體系統的回應速度。

API 實作測驗
===
## 題目一

使用技術：Go、Gin、go-micro、consul、protobuf、Nginx

此專案已使用 docker 部署在 AWS EC2 實例上，供面試官參考使用，具體 api 路徑為： http://3.25.235.191/asiayo/api/orders

單元測試檔案路徑：orderService/handler/order/orderHandler_test.go

根據循序圖可以發現題目要求檢查欄位與實際邏輯必須分離，因此建立了簡易的微服務架構，具體資料流可以參考下圖：

![資料流](https://github.com/lingjun0314/AsiaYo-test/blob/main/images/asiayo.png)

先使用 Nginx 反向代理到 API Gateway 伺服器處理請求，並檢查 json 輸入格式及資料是否正確，使用 ShouldBindJSON 方法來確保 json結構的正確，定義的結構由 proto 文件確定資料格式及形態：
```protobuf=
message OrderModule {
    string id = 1;
    string name = 2;
    AddressModule address = 3;
    string price = 4;
    string currency = 5;
}

message AddressModule {
    string city = 1;
    string district = 2;
    string street = 3;
}
```
當有資料類型不正確則回傳提示信息：Invalid JSON，格式正確則調用服務發現 consul 找到 orderService 實際 ip 地址並調用遠程方法。

微服務內部的實現為符合 SOLID 原則的設計方法，以下說明：

1. 單一職責原則（SRP）：程式在做資料檢查時僅檢查單一的部分，並不處理除了該欄位以外的邏輯，共分成 name currency price response，參考 nameHandler的程式部分：
```go=
func (n *NameHandler) IsEnglish(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name)
}

func (n *NameHandler) IsCapitalized(name string) bool {
	words := strings.Fields(name)
	for _, word := range words {
		if word[0] < 'A' || word[0] > 'Z' {
			return false
		}
	}
	return true
}
```
2. 開放封閉原則（OCP）：現有函式不需修改就可以達成功能，若要擴展功能可以再新增函數，實現擴展開放，修改封閉。可以參考 name 方法的 interface：
```go=
type NameLogic interface {
	IsEnglish(string) bool
	IsCapitalized(string) bool
    //  這裡寫新的函式以擴展功能
}
```
3. 里氏替換原則（LSP）：程式中的接口實現為結構體，如上方的接口實現的結構體如下：
```go=
type NameHandler struct{}

func (n *NameHandler) IsEnglish(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name)
}

func (n *NameHandler) IsCapitalized(name string) bool {
	words := strings.Fields(name)
	for _, word := range words {
		if word[0] < 'A' || word[0] > 'Z' {
			return false
		}
	}
	return true
}
```
當然也可以再寫一個不一樣的結構實現 interface，這並不影響整體程式碼結構，如：
```go=
type NewNameHandler struct{}

func (n *NewNameHandler) IsEnglish(name string) bool {
	return false
}

func (n *NewNameHandler) IsCapitalized(name string) bool {
	return true
}
```

4. 接口隔離原則（ISP）：這裡的實現在 API Gateway 調用服務時，不需要依賴微服務內部的接口，只需要依賴微服務提供的方法接口就可以調用函數。下面是微服務內部的接口：
```go=
type NameLogic interface {
	IsEnglish(string) bool
	IsCapitalized(string) bool
}

type PriceLogic interface {
	PriceIsNumber(string) bool
	PriceOverTwoThousands(string) bool
	PriceUSDToTWD(string) string
	GetPriceInt(string) int
}

type CurrencyLogic interface {
	IsCurrencyFormatValid(string) bool
	TransformUSDToTWD(currency, price string) (string, string)
}

type ResponseLogic interface {
	SetResponseFailure(res *pb.CheckAndTransformDataResponse, message string)
	SetResponseSuccess(req *pb.CheckAndTransformDataRequest, res *pb.CheckAndTransformDataResponse, message string)
}
```
在 API Gateway 依賴的接口是微服務提供的方法接口：
```go=
type OrderService interface {
	CheckAndTransformData(ctx context.Context, in *CheckAndTransformDataRequest, opts ...client.CallOption) (*CheckAndTransformDataResponse, error)
}
```
5. 依賴反轉原則（DIP）：高層模組並不依賴於低層模組，二者都依賴於抽象接口。此專案中高層模組為 order 模組，可以看到依賴於 interface：
```go=
type Order struct {
	name     NameLogic
	price    PriceLogic
	currency CurrencyLogic
	response ResponseLogic
}
```
而所有的低層模組都是接口的實現，也並不依賴於高層模組。

為了實現依賴反轉，在使用方法時必須進行依賴注入，先定義一個方法回傳包含依賴的結構體：
```go=
func NewOrder(name NameLogic, price PriceLogic, currency CurrencyLogic, response ResponseLogic) *Order {
	return &Order{
		name:     name,
		price:    price,
		currency: currency,
		response: response,
	}
}
```
當調用方法時進行依賴注入：
```go=
order := NewOrder(&NameHandler{}, &PriceHandler{}, &CurrencyHandler{}, &ResponseHandler{})
```