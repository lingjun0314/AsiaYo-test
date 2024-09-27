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

使用技術：Go、Gin、go-micro、consul
此專案已部署在 AWS EC2 實例上，供面試官參考使用，具體 api 路徑為： http://3.25.235.191/asiayo/api/orders
根據循序圖可以發現題目要求檢查欄位與實際邏輯必須分離，因此建立了建議微服務架構，具體資料流可以參考下圖：
![資料流](https://github.com/lingjun0314/AsiaYo-test/blob/main/images/asiayo.png)