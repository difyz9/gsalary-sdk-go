# GSalary 对接文档 (v1.23)

## 目录
1. [更新记录](#更新记录)
2. [鉴权](#鉴权)
3. [错误码清单](#错误码清单)
4. [业务字典](#业务字典)
5. [规范](#规范)
6. [核心接口](#核心接口)
   - [钱包](#钱包)
   - [换汇](#换汇)

---

## 更新记录
| 日期       | 版本  | 详情                                                                 |
|------------|-------|----------------------------------------------------------------------|
| 2025-10-17 | 1.23  | 移除修改卡是否支持3DS交易API                                          |
| 2025-08-20 | 1.22  | 新增卡交易类型: CHARGEBACK                                           |
| 2025-08-18 | 1.21  | 新增收单服务: 便捷安全支付和卡自动代扣相关接口                        |
| 2025-08-08 | 1.20  | 新增卡交易类型: TDS_FEE，新增修改卡是否支持3DS交易API                 |
| 2025-07-30 | 1.19  | 新增实体卡的支持，新增卡交易类型、webhook事件类型和卡状态             |
| 2025-07-11 | 1.18  | 新增收单业务相关接口                                                |
| 2025-04-21 | 1.17  | 新增付款状态PROCESS_FAILED，新增推送事件类型REMITTANCE_REVERSE，补充python sdk |
| 2025-04-14 | 1.16  | 新增充值卡的支持                                                    |
| 2025-04-10 | 1.15  | 新增更新收款人账户（电子钱包）API                                    |
| 2025-03-24 | 1.14  | PAYEE_DEACTIVATED中deactivated_payee_ids字段弃用，改为deactivate_payee_id，旧字段保留兼容性 |
| 2025-02-07 | 1.13  | 卡信息查询增加card_holder_id字段，卡余额变更记录新增card_id字段       |
| 2025-01-23 | 1.12  | 调额申请记录新增gsalary_request_id，用于和卡余额变更记录的transaction_id进行关联 |
| 2024-12-25 | 1.11  | 新增查询支付方式支持付款国家和币种列表API，停用查询支持付款国家(银行账户)API |
| 2024-12-04 | 1.10  | 调整账务流水和卡余额变更记录业务类型                                  |
| 2024-11-27 | 1.9   | 卡交易详情和卡交易通知新增status_description字段                      |
| 2024-11-21 | 1.8   | 移除新增卡接口限额配置的必填要求                                    |
| 2024-10-17 | 1.7   | 银行卡付款相关API                                                    |
| 2024-09-04 | 1.6   | 卡交易流水增加手续费字段，账户余额接口返回转账可用金额字段            |
| 2024-09-03 | 1.5   | 流水业务类型新增OTHER类别                                            |
| 2024-08-23 | 1.4   | 修改部分卡交易类型、流水账务类型枚举项                              |
| 2024-08-01 | 1.3   | 新增支持多种收款账户数据结构，调整部分必填项                          |
| 2024-07-03 | 1.2   | 新增附件上传API以及付款人证件参数                                    |
| 2024-06-07 | 1.1   | 新增付款相关API                                                      |
| 2024-05-13 | 1.0   | 初版                                                                 |

---

## 鉴权
### 配置
开发者需生成 RSA 密钥对，完成以下步骤：
1. 生成 RSA 密钥对（4096 位）
2. 在 GSalary Portal 配置商户公钥（无需去除 `-----BEGIN PUBLIC KEY-----` 标签）
3. 获取平台分配的服务器公钥和 Appid
4. 妥善保存商户私钥

#### Bash 密钥生成示例
```bash
# 生成 RSA 私钥
openssl genrsa 4096 > private.pem
# 从私钥生成公钥
openssl rsa -in private.pem -pubout -out public.pem
# 转换私钥为 PKCS8 格式（Java 等语言需此格式）
openssl pkcs8 -topk8 -in private.pem -nocrypt -out private-pkcs8.pem
```

### 签名流程
#### 第一步：计算请求体 Hash
使用 SHA256 算法计算请求体哈希，并用 Base64 编码（无请求体时留空）。

##### Python 示例
```python
import hashlib
import base64

body = '{}'  # 实际请求体
body_hash = base64.b64encode(hashlib.sha256(body.encode('utf-8')).digest()).decode()
```

##### Java 示例
```java
import org.apache.commons.codec.digest.DigestUtils;
import java.util.Base64;

String body = "{}";
String bodyHash = Base64.getEncoder().encodeToString(DigestUtils.sha256(body.getBytes(StandardCharsets.UTF_8)));
```

#### 第二步：拼接签名字符串
格式如下（换行符为 `\n`，末尾需保留空行）：
```
<METHOD> <PATH>
<APPID>
<TIMESTAMP>
<BODY_HASH>

```
- `METHOD`：HTTP 请求方法（大写，如 GET/POST）
- `PATH`：请求路径（含查询参数，无需转义）
- `APPID`：平台分配的 Appid
- `TIMESTAMP`：毫秒级 Unix 时间戳（±5 分钟内有效）
- `BODY_HASH`：第一步计算的哈希值（无请求体时留空）

##### Python 拼接示例
```python
method = "POST"
sign_path = "/exchange/quotes"
appid = "your_appid"
timestamp = "1716888888888"

sign_base = f"""{method} {sign_path}
{appid}
{timestamp}
{body_hash}
"""
```

#### 第三步：私钥签名
使用 `SHA256WithRSA` 算法对签名字符串签名，结果用 Base64 编码。

#### 第四步：设置 HTTP 头
```http
X-Appid: {APPID}
Authorization: algorithm=RSA2,time={TIMESTAMP},signature={SIGNATURE_BASE64}
```
> 注意：Base64 签名需做 URL-Safe 处理（如 `+` 替换为 `%2B`）。

### 验签流程
1. 解析 `X-Appid` 和 `Authorization` 头
2. 计算响应体 SHA256-Base64 哈希
3. 按签名流程第二步拼接字符串
4. 使用平台公钥和 `SHA256WithRSA` 算法验签

### SDK/Demo
#### Java
- GitHub：https://github.com/gsalary-develop/gsalary-sdk-java
- Maven 依赖：
  ```xml
  <dependency>
    <groupId>com.gsalary</groupId>
    <artifactId>gsalary-sdk-java</artifactId>
    <version>1.0.1</version>
  </dependency>
  ```

#### PHP
- GitHub：https://github.com/gsalary-develop/gsalary-demo-php

#### Python
- PyPI：https://pypi.org/project/gsalary-sdk
- 安装：`pip install gsalary-sdk`
- GitHub：https://github.com/gsalary-develop/gsalary-sdk-python

---

## 错误码清单
| 错误码                  | HTTP 状态码 | 描述                                                                 |
|-------------------------|-------------|----------------------------------------------------------------------|
| SYSTEM_ERROR            | 500         | 系统错误                                                             |
| ADD_CARD_FAILED         | 500         | 申请开卡失败                                                         |
| CREATE_PAYEE_ACCOUNT_FAILED | 500      | 新增收款人收款账户失败                                               |
| UPDATE_PAYEE_ACCOUNT_FAILED | 500      | 更新收款人账户失败                                                   |
| SYSTEM_BUSY             | 423         | 系统繁忙，请稍后重试                                                 |
| NOT_FOUND               | 404         | 请求对象不存在                                                       |
| FORBIDDEN               | 403         | 无权限访问                                                           |
| BAD_REQUEST             | 400         | 请求参数不符合要求                                                   |
| MISSING_ARGUMENT        | 400         | 缺少参数                                                             |
| INVALID_ARGUMENT        | 400         | 参数不合法                                                           |
| INVALID_STATUS          | 400         | 请求对象状态不合法                                                   |
| DUPLICATED              | 400         | 重复请求                                                             |
| QUOTE_EXPIRE            | 400         | 锁汇超时                                                             |
| ORDER_EXPIRE            | 400         | 订单超时                                                             |
| INSUFFICIENT_BALANCE    | 400         | 账户余额不足                                                         |
| RISK_REJECT             | 400         | 达到付款次数限制（仅 ALIPAY 和 CNY 币种银行转账）                    |
| USER_AMOUNT_EXCEED_LIMIT | 400        | 支付金额超过用户支付限额                                             |
| USER_BALANCE_NOT_ENOUGH | 400         | 支付方式用户余额不足                                                 |

---

## 业务字典
### 汇款目的
| 枚举值                          | 描述                                   | 企业主体 | 个人主体 |
|---------------------------------|----------------------------------------|----------|----------|
| SALARY                          | 发薪                                   | Y        | Y        |
| FAMILY_SUPPORT                  | 家庭支出                               | N        | Y        |
| CLOTHES_BAGS_SHOES              | 服装 鞋帽 箱包购物                     | N        | Y        |
| DAILY_SUPPLIES_AND_COSMETICS    | 化妆 日用品购物                         | N        | Y        |
| ELECTRONICS_AND_HOME_APPLIANCES | 数码家电购物                           | N        | Y        |
| TOYS_KIDS_BABIES                | 玩具 婴幼儿用品购物                     | N        | Y        |
| INTERPRETATION_SERVICE          | 口译服务费用                           | Y        | Y        |
| TRANSLATION_SERVICE             | 笔译服务费用                           | Y        | Y        |
| HUMAN_RESOURCE_SERVICE          | 人才中介服务费用                       | Y        | Y        |
| ESTATE_AGENCY_SERVICE           | 房屋中介服务费用                       | Y        | Y        |
| SOFTWARE_DEVELOPMENT_SERVICE    | 软件开发者服务费用                     | Y        | Y        |
| WEB_DESIGN_OR_DEVELOPMENT_SERVICE | 网站开发/网页设计类服务费用           | Y        | Y        |
| DRAFTING_LEGAL_SERVICE          | 起草法务文件服务费用                   | Y        | Y        |
| LEGAL_RELATED_CERTIFICATION_SERVICE | 法律相关认证服务费用                 | Y        | Y        |
| ACCOUNTING_SERVICE              | 会计记录报表审计咨询规划服务费用       | Y        | Y        |
| TAX_SERVICE                     | 准备税务文件服务费用                   | Y        | Y        |
| ARCHITECTURAL_DECORATION_DESIGN_SERVICE | 建筑装潢设计服务费用                 | Y        | Y        |
| ADVERTISING_SERVICE             | 广告设计服务费用                       | Y        | Y        |
| MARKET_RESEARCH_SERVICE         | 市场调查服务费用                       | Y        | Y        |
| EXHIBITION_BOOTH_SERVICE        | 展会摊位租赁服务费用                   | Y        | Y        |
| PRODUCT_PROMOTION_SERVICE       | 产品内容推广服务费收入                 | Y        | N        |
| ECOMMERCE_PROMOTION_SERVICE     | 电商成交订单佣金服务费收入             | Y        | N        |

### 支付方式
| 国家/地区       | 支付方式         | 枚举值               | 类型       |
|----------------|------------------|----------------------|------------|
| 全球           | Card             | CARD                 | 卡         |
| 全球           | Google Pay       | GOOGLEPAY            | 电子钱包   |
| 全球           | Apple Pay        | APPLEPAY             | 电子钱包   |
| 中国           | Alipay           | ALIPAY_CN            | 电子钱包   |
| 中国香港       | AlipayHK         | ALIPAY_HK            | 电子钱包   |
| 巴西           | Pix              | PIX                  | 实时支付   |
| 新加坡         | PayNow           | PAYNOW               | 实时支付   |
| 泰国           | PromptPay        | PROMPTPAY            | 实时支付   |
| 日本           | Konbini          | KONBINI              | 现金支付   |

> 完整支付方式清单见官方文档，以上为常用示例。

---

## 规范
### 持卡人电话号码验证规则
| 国家/地区       | ISO2 | 国家/地区代码 | 手机号码长度 |
|----------------|------|---------------|--------------|
| 中国大陆       | CN   | 86            | 11           |
| 中国香港       | HK   | 852           | 8            |
| 中国澳门       | MO   | 853           | 8            |
| 中国台湾       | TW   | 886           | 9            |
| 美国           | US   | 1             | 10           |
| 新加坡         | SG   | 65            | 8            |
| 日本           | JP   | 81            | 10           |
| 韩国           | KR   | 82            | 10           |

> 注：长度验证针对不含国家/地区代码的手机号部分。

---

## 核心接口
### 钱包
#### 1. 查询钱包余额
- **请求方式**：GET
- **接口路径**：`/wallets/balance`
- **请求参数**
  | 参数位置 | 参数名   | 是否必填 | 类型   | 说明                     |
  |----------|----------|----------|--------|--------------------------|
  | query    | currency | 是       | string | 币种（参考 ISO-4217 标准）|
  | header   | x-appid  | 是       | string | 平台分配的 Appid         |
  | header   | authorization | 是    | string | 签名信息（按鉴权规则生成）|
- **响应示例**（200 OK）
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "currency": "USD",
    "amount": 1000.00,
    "share_card_account_balance": 500.00,
    "available": 800.00,
    "account_type": "BALANCE",
    "query_time": "2024-04-03T12:00:00Z"
  }
}
```

#### 2. 查询钱包流水
- **请求方式**：GET
- **接口路径**：`/wallets/transactions`
- **请求参数**
  | 参数位置 | 参数名     | 是否必填 | 类型    | 说明                                  |
  |----------|------------|----------|---------|---------------------------------------|
  | query    | page       | 是       | integer | 页码（默认 1）                        |
  | query    | limit      | 是       | integer | 每页条数（默认 20）                   |
  | query    | time_start | 否       | string  | 起始时间（ISO-8601，含）              |
  | query    | time_end   | 否       | string  | 截止时间（ISO-8601，不含）            |
  | query    | currency   | 否       | string  | 币种（ISO-4217）                      |
  | query    | txn_type   | 否       | string  | 流水类型（枚举值见下方）              |
  | header   | x-appid    | 是       | string  | 平台分配的 Appid                      |
  | header   | authorization | 是     | string  | 签名信息                              |
- **流水类型枚举**
  | 枚举值                  | 描述                                                                 |
  |-------------------------|----------------------------------------------------------------------|
  | ACCOUNT_PAY             | 账户支付（如转账）                                                   |
  | BALANCE_RECHARGE        | 充值到资金账户                                                       |
  | EXCHANGE_IN             | 换汇买入                                                             |
  | EXCHANGE_OUT            | 换汇卖出                                                             |
  | CARD_PAYMENT            | 卡交易结算                                                           |
  | CARD_REFUND


# 接口文档
## 一、换汇订单
### 1.1 查询换汇订单
- **请求方式**：GET
- **接口路径**：`/exchange/orders`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 默认值 | 说明 |
  |----------|--------|----------|------|--------|------|
  | query | page | 是 | integer | 1 | 页码，从1开始计算 |
  | query | limit | 是 | integer | 20 | 每页记录条数 |
  | query | time_start | 否 | string | - | 查询起始时间（含），使用ISO-8601时间格式，示例：2024-03-04T10:00:00Z |
  | query | time_end | 否 | string | - | 查询截止时间（不含），使用ISO-8601时间格式，示例：2024-03-05T10:00:00Z |
  | query | status | 否 | string | - | 换汇订单状态，枚举值：PENDING、SUCCESS、FAIL |
  | query | buy_currency | 否 | string | - | 购入币种，参考ISO-4217币种清单，示例：USD |
  | query | sell_currency | 否 | string | - | 卖出币种，参考ISO-4217币种清单，示例：JPY |
  | header | x-appid | 是 | string | - | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | - | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "query": {},
      "orders": [],
      "page": 0,
      "limit": 0,
      "total_count": 0,
      "total_page": 0
    }
  }
  ```

## 二、持卡人
### 2.1 概述
管理持卡人信息相关接口集合

### 2.2 查询持卡人
- **请求方式**：GET
- **接口路径**：`/card_holders`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 默认值 | 说明 |
  |----------|--------|----------|------|--------|------|
  | query | page | 是 | integer | 1 | 页码，从1开始计算 |
  | query | limit | 是 | integer | 20 | 每页记录条数 |
  | query | time_start | 否 | string | - | 查询创建起始时间（含），使用ISO-8601时间格式，示例：2024-03-04T02:00:00Z |
  | query | time_end | 否 | string | - | 查询创建截止时间（不含），使用ISO-8601时间格式，示例：2024-03-04T09:00:00Z |
  | header | x-appid | 是 | string | - | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | - | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "query": {},
      "card_holders": [],
      "page": 0,
      "limit": 0,
      "total_count": 0,
      "total_page": 0
    }
  }
  ```

### 2.3 添加持卡人
- **请求方式**：POST
- **接口路径**：`/card_holders`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | first_name | 是 | string | 持卡人名字。长度需在1至40个字符之间，仅可包含英文字母和空格 |
  | body | last_name | 是 | string | 持卡人姓氏。长度需在1至40个字符之间，仅可包含英文字母和空格 |
  | body | birth | 是 | string | 持卡人生日，使用ISO-8601日期格式，示例：2001-12-20 |
  | body | email | 是 | string | 请确保持卡人邮箱唯一性 |
  | body | mobile | 是 | object (MobileNumber) | 持卡人手机号，验证规则参考持卡人的电话号码验证规则 |
  | body | region | 是 | string | 2字符国家码，参考ISO-3166标准国家码清单，示例：US |
  | body | bill_address | 否 | object (Address) | 账单地址 |
- **请求示例**：
```json
{
  "first_name": "string",
  "last_name": "string",
  "birth": "2001-12-20",
  "email": "string",
  "mobile": {
    "nation_code": "86",
    "mobile": "string"
  },
  "region": "US",
  "bill_address": {
    "postcode": "string",
    "address": "string",
    "city": "string",
    "state": "string",
    "country": "string"
  }
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "card_holder_id": "string",
      "first_name": "string",
      "last_name": "string",
      "birth": "2001-12-20",
      "email": "string",
      "region": "string",
      "create_time": "2024-03-02T12:25:01Z",
      "bill_address": {},
      "status": "PENDING"
    }
  }
  ```

### 2.4 查看持卡人信息
- **请求方式**：GET
- **接口路径**：`/card_holders/{card_holder_id}`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_holder_id | 是 | string | 创建持卡人API返回的持卡人id |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "card_holder_id": "string",
      "first_name": "string",
      "last_name": "string",
      "birth": "2001-12-20",
      "email": "string",
      "region": "string",
      "create_time": "2024-03-02T12:25:01Z",
      "bill_address": {},
      "status": "PENDING"
    }
  }
  ```

### 2.5 修改持卡人信息
- **请求方式**：PUT
- **接口路径**：`/card_holders/{card_holder_id}`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_holder_id | 是 | string | 创建持卡人API返回的持卡人id |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | first_name | 否 | string | 持卡人名字。长度需在1至40个字符之间，仅可包含英文字母和空格 |
  | body | last_name | 否 | string | 持卡人姓氏。长度需在1至40个字符之间，仅可包含英文字母和空格 |
  | body | birth | 否 | string | 持卡人生日，使用ISO-8601日期格式，示例：2001-12-20 |
  | body | email | 否 | string | 请确保持卡人邮箱唯一性 |
  | body | mobile | 否 | object (MobileNumber) | 持卡人手机号，验证规则参考持卡人的电话号码验证规则 |
  | body | region | 否 | string | 2字符国家码，参考ISO-3166标准国家码清单 |
  | body | bill_address | 否 | object (Address) | 账单地址 |
- **请求示例**：
```json
{
  "first_name": "string",
  "last_name": "string",
  "birth": "2001-12-20",
  "email": "string",
  "mobile": {
    "nation_code": "86",
    "mobile": "string"
  },
  "region": "US",
  "bill_address": {
    "postcode": "string",
    "address": "string",
    "city": "string",
    "state": "string",
    "country": "string"
  }
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "card_holder_id": "string",
      "first_name": "string",
      "last_name": "string",
      "birth": "2001-12-20",
      "email": "string",
      "region": "string",
      "create_time": "2024-03-02T12:25:01Z",
      "bill_address": {},
      "status": "PENDING"
    }
  }
  ```

## 三、卡
### 3.1 概述
包含卡管理、卡交易查询相关接口，开卡前需先创建持卡人

### 3.2 业务指南
#### 3.2.1 虚拟卡开卡流程
1. 调用**申请开卡**接口申请开卡
2. 调用**查询开卡结果**接口查询开卡结果
3. 接收卡状态变更通知 (发行卡事件推送, business_type=CARD_STATUS_UPDATE)

#### 3.2.2 实体卡开卡流程
1. 调用**分配卡**接口，为持卡人分配一张实体卡
2. 接收卡激活码 (发行卡事件推送, business_type=CARD_ACTIVATION_CODE)
3. 调用**激活卡**接口，激活实体卡
4. 接收卡状态变更通知 (发行卡事件推送, business_type=CARD_STATUS_UPDATE)

### 3.3 查询卡可用配额
- **请求方式**：GET
- **接口路径**：`/cards/available_quotas`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | query | currency | 是 | string | 卡币种，参考ISO-4217币种清单(支持的币种：USD) |
  | query | accounting_card_type | 否 | string | 卡账务类型，枚举值：SHARE、RECHARGE，不填默认为SHARE |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "query": {},
      "quota": {}
    }
  }
  ```

### 3.4 查看可用的卡产品列表
- **请求方式**：GET
- **接口路径**：`/card_support/products`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | query | card_type | 否 | string | 卡类型，枚举值：PHYSICAL（实体卡）、VIRTUAL（虚拟卡） |
  | query | brand_code | 否 | string | 卡品牌，枚举值：VISA、MASTER |
  | query | currency | 否 | string | 卡币种，参考ISO-4217币种清单 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "query": {},
      "products": []
    }
  }
  ```

### 3.5 申请开卡
- **请求方式**：POST
- **接口路径**：`/card_applies`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | request_id | 是 | string | 唯一请求id，最多50个字符 |
  | body | product_code | 是 | string | 卡产品编码 |
  | body | currency | 是 | string | 卡币种，参考ISO-4217币种清单 |
  | body | card_holder_id | 是 | string | 注册的持卡人id |
  | body | limit_per_day | 否 | number | 每日交易限额，>=1，不填写则采用系统默认值，是否支持配置参考查看可用的卡产品列表返回的限额参数 |
  | body | limit_per_month | 否 | number | 每月交易限额，>=1，不填写则采用系统默认值，是否支持配置参考查看可用的卡产品列表返回的限额参数 |
  | body | limit_per_transaction | 否 | number | 单笔交易限额，>=1，不填写则采用系统默认值，是否支持配置参考查看可用的卡产品列表返回的限额参数 |
  | body | init_balance | 是 | number | 初始额度 |
- **请求示例**：
```json
{
  "request_id": "string",
  "product_code": "string",
  "currency": "string",
  "card_holder_id": "string",
  "limit_per_day": 1,
  "limit_per_month": 1,
  "limit_per_transaction": 1,
  "init_balance": 0
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "request_id": "string",
      "status": "PENDING"
    }
  }
  ```

### 3.6 查询开卡结果
- **请求方式**：GET
- **接口路径**：`/card_applies/{request_id}`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | request_id | 是 | string | 开卡请求时提交的唯一请求id |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "request_id": "string",
      "status": "PENDING",
      "create_time": "string",
      "card_detail": {}
    }
  }
  ```

### 3.7 查询卡列表
- **请求方式**：GET
- **接口路径**：`/cards`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 默认值 | 说明 |
  |----------|--------|----------|------|--------|------|
  | query | page | 是 | integer | 1 | 页码，从1开始计算 |
  | query | limit | 是 | integer | 20 | 每页记录条数 |
  | query | product_code | 否 | string | - | 卡产品编码，参考查询卡产品列表 |
  | query | brand_code | 否 | string | - | 卡品牌，枚举值：VISA、MASTER |
  | query | card_holder_id | 否 | string | - | 持卡人id，1.13新增 |
  | query | create_start | 否 | string | - | 查询创建卡起始时间（含），使用ISO-8601时间格式 |
  | query | create_end | 否 | string | - | 查询创建卡截止时间（不含），使用ISO-8601时间格式 |
  | query | status | 否 | string | - | 卡状态，枚举值及描述：<br/>PENDING-卡片正在激活<br/>INACTIVE-待激活<br/>ACTIVE-正常<br/>FREEZING-正在冻结<br/>FROZEN-已冻结<br/>UNFREEZING-正在解冻<br/>EXPIRED-已过期<br/>CANCELLING-正在销卡<br/>CANCELLED-已销卡 |
  | header | x-appid | 是 | string | - | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | - | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "query": {},
      "cards": [],
      "page": 0,
      "limit": 0,
      "total_count": 0,
      "total_page": 0
    }
  }
  ```

### 3.8 查看卡信息
- **请求方式**：GET
- **接口路径**：`/cards/{card_id}`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一卡id，从查询开卡请求结果中获取 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "card_id": "string",
      "card_name": "string",
      "mask_card_number": "41******1111",
      "card_currency": "string",
      "available_balance": 0,
      "brand_code": "VISA",
      "status": "PENDING",
      "card_type": "PHYSICAL",
      "accounting_type": "SHARE",
      "card_region": "US",
      "card_holder_id": "string",
      "first_name": "string",
      "last_name": "string",
      "mobile": {},
      "email": "string",
      "limit_per_day": 0,
      "limit_per_month": 0,
      "limit_per_transaction": 0,
      "bill_address": {},
      "support_tds_trans": true,
      "create_time": "2024-04-02T13:03:12Z"
    }
  }
  ```

### 3.9 修改卡信息
- **接口说明**：修改卡基本信息
- **请求方式**：PUT
- **接口路径**：`/cards/{card_id}`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一卡id，从查询开卡请求结果中获取 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | card_name | 否 | string | 卡昵称 |
  | body | limit_per_day | 否 | number | 每日限额(需满足单笔<=每日<=每月) |
  | body | limit_per_month | 否 | number | 每月限额(需满足单笔<=每日<=每月) |
  | body | limit_per_transaction | 否 | number | 单笔交易限额(需满足单笔<=每日<=每月) |
- **请求示例**：
```json
{
  "card_name": "string",
  "limit_per_day": 0,
  "limit_per_month": 0,
  "limit_per_transaction": 0
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "card_id": "string",
      "card_name": "string",
      "mask_card_number": "41******1111",
      "card_currency": "string",
      "available_balance": 0,
      "brand_code": "VISA",
      "status": "PENDING",
      "card_type": "PHYSICAL",
      "accounting_type": "SHARE",
      "card_region": "US",
      "card_holder_id": "string",
      "first_name": "string",
      "last_name": "string",
      "mobile": {},
      "email": "string",
      "limit_per_day": 0,
      "limit_per_month": 0,
      "limit_per_transaction": 0,
      "bill_address": {},
      "support_tds_trans": true,
      "create_time": "2024-04-02T13:03:12Z"
    }
  }
  ```

### 3.10 销卡
- **请求方式**：DELETE
- **接口路径**：`/cards/{card_id}`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一卡id，从查询开卡请求结果中获取 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    }
  }
  ```
  - 响应说明：
    - result枚举值：
      - S: 请求接收成功，仅表示请求被成功受理，具体业务执行结果请以业务对象为准
      - F: 请求接收失败，无需重试
      - U: 未知异常，请稍后调用查询接口确认结果
    - code：错误码，参考错误码清单
    - message：错误信息

### 3.11 获取卡机密信息
- **请求方式**：GET
- **接口路径**：`/cards/{card_id}/secure_info`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一卡id，从查询开卡请求结果中获取 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "pan": "4111111111111111",
      "expire_year": "2025",
      "expire_month": "04",
      "cvv": "312"
    }
  }
  ```

### 3.12 卡片调额
- **请求方式**：POST
- **接口路径**：`/cards/balance_modifies`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | card_id | 是 | string | 卡id |
  | body | amount | 是 | number | 修改金额，>=0，必须是正数 |
  | body | type | 是 | string | 修改类型，枚举值：INCREASE、DECREASE |
  | body | request_id | 是 | string | 唯一请求id |
- **请求示例**：
```json
{
  "card_id": "string",
  "amount": 0,
  "type": "INCREASE",
  "request_id": "string"
}
```
- **响应结果（200 Success）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "gsalary_request_id": "string",
      "request_id": "string",
      "card_id": "string",
      "status": "PENDING",
      "create_time": "2024-04-05T12:32:08Z",
      "finish_time": "string",
      "amount": 0,
      "type": "INCREASE",
      "post_balance": 0
    }
  }
  ```

### 3.13 查询卡片调额结果
- **请求方式**：GET
- **接口路径**：`/cards/balance_modifies/{request_id}`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | request_id | 是 | string | 卡片调额请求时提交的唯一请求id |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "gsalary_request_id": "string",
      "request_id": "string",
      "card_id": "string",
      "status": "PENDING",
      "create_time": "2024-04-05T12:32:08Z",
      "finish_time": "string",
      "amount": 0,
      "type": "INCREASE",
      "post_balance": 0
    }
  }
  ```

### 3.14 冻结/解冻卡
- **请求方式**：PUT
- **接口路径**：`/cards/{card_id}/freeze_status`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一卡id，从查询开卡请求结果中获取 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | freeze | 是 | boolean | 要求变更的冻结状态 |
- **请求示例**：
```json
{
  "freeze": true
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    }
  }
  ```
  - 响应说明：
    - result枚举值：
      - S: 请求接收成功，仅表示请求被成功受理，具体业务执行结果请以业务对象为准
      - F: 请求接收失败，无需重试
      - U: 未知异常，请稍后调用查询接口确认结果
    - code：错误码，参考错误码清单
    - message：错误信息

### 3.15 卡交易列表
- **接口说明**：查询卡交易列表
- **请求方式**：GET
- **接口路径**：`/card_bill/card_transactions`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 默认值 | 说明 |
  |----------|--------|----------|------|--------|------|
  | query | page | 是 | integer | 1 | 页码，从1开始计算 |
  | query | limit | 是 | integer | 20 | 每页记录条数 |
  | query | transaction_id | 否 | string | - | 卡交易id |
  | query | mch_request_id | 否 | string | - | 商户请求id |
  | query | time_start | 否 | string | - | 查询起始时间（含），使用ISO-8601时间格式 |
  | query | time_end | 否 | string | - | 查询截止时间（不含），使用ISO-8601时间格式 |
  | query | card_id | 否 | string | - | 唯一卡id |
  | header | x-appid | 是 | string | - | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | - | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "query": {},
      "transactions": [],
      "page": 0,
      "limit": 0,
      "total_count": 0,
      "total_page": 0
    }
  }
  ```

### 3.16 卡余额变更记录
- **接口说明**：查询卡余额变更记录
- **请求方式**：GET
- **接口路径**：`/card_bill/balance_history`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 默认值 | 说明 |
  |----------|--------|----------|------|--------|------|
  | query | page | 是 | integer | 1 | 页码，从1开始计算 |
  | query | limit | 是 | integer | 20 | 每页记录条数 |
  | query | transaction_id | 否 | string | - | 卡交易id |
  | query | log_id | 否 | string | - | 入账id |
  | query | time_start | 否 | string | - | 查询起始时间（含），使用ISO-8601时间格式 |
  | query | time_end | 否 | string | - | 查询截止时间（不含），使用ISO-8601时间格式 |
  | query | card_id | 否 | string | - | 卡id |
  | header | x-appid | 是 | string | - | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | - | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "query": {},
      "history": [],
      "page": 0,
      "limit": 0,
      "total_count": 0,
      "total_page": 0
    }
  }
  ```

### 3.17 修改卡联系信息
- **接口说明**：电子邮件可用于接收ApplePay等绑卡时的验证码信息
- **请求方式**：PUT
- **接口路径**：`/cards/{card_id}/contact`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一卡id，从查询开卡请求结果中获取 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | email | 否 | string | 持卡人email |
- **请求示例**：
```json
{
  "email": "string"
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "card_id": "string",
      "card_name": "string",
      "mask_card_number": "41******1111",
      "card_currency": "string",
      "available_balance": 0,
      "brand_code": "VISA",
      "status": "PENDING",
      "card_type": "PHYSICAL",
      "accounting_type": "SHARE",
      "card_region": "US",
      "card_holder_id": "string",
      "first_name": "string",
      "last_name": "string",
      "mobile": {},
      "email": "string",
      "limit_per_day": 0,
      "limit_per_month": 0,
      "limit_per_transaction": 0,
      "bill_address": {},
      "support_tds_trans": true,
      "create_time": "2024-04-02T13:03:12Z"
    }
  }
  ```

### 3.18 删除卡的邮箱
- **接口说明**：只有卡产品G68796支持删除持卡人的电子邮件
- **请求方式**：DELETE
- **接口路径**：`/cards/{card_id}/email`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一卡id，从查询开卡请求结果中获取 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "card_id": "string",
      "card_name": "string",
      "mask_card_number": "41******1111",
      "card_currency": "string",
      "available_balance": 0,
      "brand_code": "VISA",
      "status": "PENDING",
      "card_type": "PHYSICAL",
      "accounting_type": "SHARE",
      "card_region": "US",
      "card_holder_id": "string",
      "first_name": "string",
      "last_name": "string",
      "mobile": {},
      "email": "string",
      "limit_per_day": 0,
      "limit_per_month": 0,
      "limit_per_transaction": 0,
      "bill_address": {},
      "support_tds_trans": true,
      "create_time": "2024-04-02T13:03:12Z"
    }
  }
  ```

### 3.19 分配卡
- **接口说明**：为持卡人分配一张实体卡
- **请求方式**：POST
- **接口路径**：`/cards/assign_card`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | card_number | 是 | string | 实体卡卡号 |
  | body | card_holder_id | 是 | string | 持卡人唯一id |
  | body | card_currency | 是 | string | 卡币种，参考ISO-4217币种清单 |
- **请求示例**：
```json
{
  "card_number": "string",
  "card_holder_id": "string",
  "card_currency": "string"
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    },
    "data": {
      "card_id": "string"
    }
  }
  ```

### 3.20 激活卡
- **接口说明**：激活实体卡
- **请求方式**：POST
- **接口路径**：`/cards/{card_id}/active_card`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一的卡片id，在卡片列表中提供或创建卡片结果 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | activation_code | 是 | string | 卡的激活码 |
  | body | pin | 是 | string | 卡PIN(交易密码&ATM取款密码)，必须是6位数字 |
  | body | no_pin_payment_amount | 否 | number | 无需PIN验证的卡交易的允许金额，>=0，默认值为200USD；想禁止无PIN交易可设为0；仅适用于BIN 49372410（生产环境） |
- **请求示例**：
```json
{
  "activation_code": "string",
  "pin": "string",
  "no_pin_payment_amount": null
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    }
  }
  ```
  - 响应说明：
    - result枚举值：
      - S: 请求接收成功，仅表示请求被成功受理，具体业务执行结果请以业务对象为准
      - F: 请求接收失败，无需重试
      - U: 未知异常，请稍后调用查询接口确认结果
    - code：错误码，参考错误码清单
    - message：错误信息

### 3.21 重置卡PIN
- **接口说明**：重置实体卡pin
- **请求方式**：POST
- **接口路径**：`/cards/{card_id}/reset_card_pin`
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | path | card_id | 是 | string | 唯一的卡片id，在卡片列表中提供或创建卡片结果 |
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | pin | 是 | string | 卡的新PIN必须是6位数字 |
- **请求示例**：
```json
{
  "pin": "string"
}
```
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应示例：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    }
  }
  ```
  - 响应说明：
    - result枚举值：
      - S: 请求接收成功，仅表示请求被成功受理，具体业务执行结果请以业务对象为准
      - F: 请求接收失败，无需重试
      - U: 未知异常，请稍后调用查询接口确认结果
    - code：错误码，参考错误码清单
    - message：错误信息

## 四、收款人
### 4.1 概述
管理对外付款收款人信息相关接口集合

### 4.2 业务指南
#### 4.2.1 电子钱包类型
1. 调用新增收款人接口并设置account_type为E_WALLET
2. 调用新增收款人收款账户完成账户添加或获取注册URL（注册URL需发送给收款人完成注册）
3. 等待收款人账户状态变为ACTIVE

#### 4.2.2 银行账户类型
1. 调用新增收款人接口并设置account_type为BANK_ACCOUNT
2. 调用查询支持付款国家和币种获取当前收款人国家支持的付款币种
3. 调用获取收款人账户表单获得注册账户所需的表单项
4. 调用注册收款人账户提交银行账户信息
5. 等待收款人账户状态变为ACTIVE

### 4.3 新增收款人
- **请求方式**：POST
- **接口路径**：`/payees`（注：原文档未明确路径，按常规命名补充，若实际有差异需调整）
- **请求参数**
  | 参数位置 | 参数名 | 是否必填 | 类型 | 说明 |
  |----------|--------|----------|------|------|
  | header | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | header | authorization | 是 | string | 请求签名，请参考Signature Guide |
  | body | subject_type | 是 | string | 收款人主体类型，枚举值：INDIVIDUAL（个人）、ENTERPRISE（企业） |
  | body | account_type | 是 | string | 收款账户类型，枚举值：E_WALLET（电子钱包）、BANK_ACCOUNT（银行账户，1.7新增），为空时默认电子钱包 |
  | body | country | 是 | string | 国家、地区编码，遵循ISO-3166规范的2字符地区代码 |
  | body | first_name | 否 | string | 收款人名，类型为个人时必填，只允许字母和空格；当account_type是BANK_ACCOUNT且country是CN时，必须是中文 |
  | body | last_name | 否 | string | 收款人姓，类型为个人时必填，只允许字母和空格；当account_type是BANK_ACCOUNT且country是CN时，必须是中文 |
  | body | account_holder | 否 | string | 账户持有人（公司名称），类型为公司时必填 |
  | body | currency | 是 | string | 收款币种，遵循ISO-4217币种编码清单 |
- **响应结果（200 OK）**
  - 响应格式：application/json
  - 响应结构：
  ```json
  {
    "result": {
      "result": "S",
      "code": "string",
      "message": "string"
    }
  }
  ```
 （注：原文档未提供完整响应示例，按统一格式补充基础结构）

 # 收款人管理接口文档
## 概述
本文档详细描述了对外付款收款人信息管理的相关接口，包括收款人CRUD、收款账户管理、付款方式查询等功能，适用于电子钱包和银行账户两种收款类型。

## 核心接口
### 1. 新增收款人
- **请求方式**：POST
- **接口路径**：`/remittance/payees`
- **接口描述**：创建新的收款人信息，支持电子钱包（E_WALLET）和银行账户（BANK_ACCOUNT）两种类型
- **请求头参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | x-appid | 是 | string | 从Portal获取的appid，请参考Signature Guide |
  | authorization | 是 | string | 请求签名，请参考Signature Guide |
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | subject_type | 是 | string | 收款人主体类型，枚举值：INDIVIDUAL（个人）、ENTERPRISE（企业） |
  | account_type | 是 | string | 收款账户类型，枚举值：E_WALLET（电子钱包）、BANK_ACCOUNT（银行账户，1.7新增），为空时默认电子钱包 |
  | country | 是 | string | 国家/地区编码，遵循ISO-3166规范的2字符地区代码 |
  | first_name | 否 | string | 收款人名，个人类型必填；仅允许字母和空格；银行账户+中国（CN）时必须为中文 |
  | last_name | 否 | string | 收款人姓，个人类型必填；仅允许字母和空格；银行账户+中国（CN）时必须为中文 |
  | account_holder | 否 | string | 账户持有人（公司名称），企业类型必填 |
  | currency | 是 | string | 收款币种，遵循ISO-4217币种编码清单 |
- **请求示例（电子钱包类型）**
```json
{
  "subject_type": "INDIVIDUAL",
  "account_type": "E_WALLET",
  "country": "CN",
  "first_name": "LEI",
  "last_name": "LI",
  "account_holder": "",
  "currency": "USD"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payee_id": "string",
    "subject_type": "INDIVIDUAL",
    "account_type": "E_WALLET",
    "first_name": "LEI",
    "last_name": "LI",
    "country": "CN",
    "currencies": ["USD"],
    "account_holder": "",
    "mobile": {},
    "address": "string"
  }
}
```

### 2. 查询收款人列表
- **请求方式**：GET
- **接口路径**：`/remittance/payees`
- **接口描述**：分页查询收款人列表，支持多条件筛选
- **请求头参数**：同「新增收款人」
- **查询参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | page | 是 | integer | 页码，从1开始，默认值：1 |
  | limit | 是 | integer | 每页记录条数，默认值：20 |
  | payee_id | 否 | string | 收款人id |
  | name | 否 | string | 姓名或公司名称 |
  | country | 否 | string | 国家/地区编码（ISO-3166 2字符） |
  | currency | 否 | string | 币种编码（ISO-4217） |
  | mobile | 否 | string | 联系电话（不带国家区号） |
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "query": {},
    "payees": [],
    "page": 1,
    "limit": 20,
    "total_count": 100,
    "total_page": 5
  }
}
```

### 3. 更新收款人信息
- **请求方式**：PUT
- **接口路径**：`/remittance/payees/{payee_id}`
- **接口描述**：全量更新收款人信息（需完整提交所有字段，而非仅修改项）
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_id | 是 | string | 收款人id |
- **请求体参数**：同「新增收款人」，其中`account_type`需与新增时一致或不传（默认沿用原类型）
- **请求示例（电子钱包类型）**
```json
{
  "account_type": "E_WALLET",
  "country": "CN",
  "first_name": "LEI",
  "last_name": "LI",
  "account_holder": "",
  "currency": "USD"
}
```
- **响应结果（200 OK）**：同「新增收款人」响应结构

### 4. 停用收款人
- **请求方式**：DELETE
- **接口路径**：`/remittance/payees/{payee_id}`
- **接口描述**：停用指定收款人（仅受理请求，实际业务结果需以查询为准）
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_id | 是 | string | 收款人id |
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  }
}
```
- **响应说明**
  - `result`枚举值：
    - S：请求接收成功（需查询确认业务结果）
    - F：请求接收失败（无需重试）
    - U：未知异常（需后续查询确认）

### 5. 查询可用付款方式
- **请求方式**：GET
- **接口路径**：`/remittance/available_payment_methods`
- **接口描述**：查询当前支持的所有付款方式
- **请求头参数**：同「新增收款人」
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payment_methods": ["ALIPAY", "WECHATPAY", "BANK_TRANSFER"]
  }
}
```

### 6. 新增收款人收款账户（电子钱包）
- **请求方式**：POST
- **接口路径**：`/remittance/payees/{payee_id}/accounts`
- **接口描述**：为收款人添加电子钱包收款账户，同一付款类型仅保留最后配置
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_id | 是 | string | 收款人id |
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payment_method | 是 | string | 付款方式（如ALIPAY） |
  | account_no | 是 | string | 账号id（测试环境可用：86-13721473389，需配套姓名LI LEI） |
- **请求示例（支付宝）**
```json
{
  "payment_method": "ALIPAY",
  "account_no": "86-13721473389"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payment_method": "ALIPAY",
    "currencies": ["USD"],
    "account_id": "string",
    "account_no": "86-13721473389",
    "status": "PENDING"
  }
}
```

### 7. 查看收款人可用收款账户
- **请求方式**：GET
- **接口路径**：`/remittance/payees/{payee_id}/accounts`
- **接口描述**：查询指定收款人的所有可用收款账户
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_id | 是 | string | 收款人id |
- **查询参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | language | 否 | string | 表单语言，默认值：en（1.7新增） |
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "accounts": [
      {
        "payment_method": "ALIPAY",
        "currencies": ["USD"],
        "account_id": "string",
        "account_no": "86-13721473389",
        "status": "ACTIVE"
      }
    ]
  }
}
```

### 8. 更新收款人账户（电子钱包）
- **请求方式**：PUT
- **接口路径**：`/remittance/payees/{payee_id}/payee_accounts/{accountId}`
- **接口描述**：更新收款人电子钱包收款账户信息
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_id | 是 | string | 收款人id |
  | accountId | 是 | string | 收款账户id |
- **请求体参数**：同「新增收款人收款账户（电子钱包）」
- **响应结果（200 OK）**：同「新增收款人收款账户（电子钱包）」响应结构

### 9. 获取收款人账户表单（银行账户）
- **请求方式**：GET
- **接口路径**：`/remittance/payees/{payee_id}/account_register_format`
- **接口描述**：获取银行账户类型收款人的注册表单字段（1.7新增）
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_id | 是 | string | 收款人id |
- **查询参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payment_method | 是 | string | 付款方式，固定值：BANK_TRANSFER |
  | currency | 否 | string | 收款账户币种（收款人支持多币种时必填） |
  | language | 否 | string | 表单语言，默认值：en |
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "account_type": "BANK_ACCOUNT",
  "payment_method": "BANK_TRANSFER",
  "subject_type": "INDIVIDUAL",
  "currency": "CNY",
  "country": "CN",
  "fields": [
    {
      "field_name": "first_name",
      "required": true,
      "description": "收款人名（中文）"
    },
    {
      "field_name": "account_no",
      "required": true,
      "description": "银行卡账号（一类卡）"
    }
  ]
}
```

### 10. 新增收款人收款账户（银行账户）
- **请求方式**：POST
- **接口路径**：`/remittance/payees/{payee_id}/account_registry`
- **接口描述**：通过表单提交银行账户信息完成注册（1.7新增），需先调用「获取收款人账户表单」接口
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_id | 是 | string | 收款人id |
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payment_method | 是 | string | 付款方式，与表单获取时一致（如BANK_TRANSFER） |
  | currency | 否 | string | 账户币种，与表单获取时一致 |
  | fields | 是 | Array<object> | 表单字段集合（测试环境CNY+BANK_TRANSFER示例：first_name=三丰，last_name=张，account_no=622306492053631234，cert_number=330100198001010001） |
- **请求示例**
```json
{
  "payment_method": "BANK_TRANSFER",
  "currency": "CNY",
  "fields": [
    {
      "field_name": "first_name",
      "field_value": "三丰"
    },
    {
      "field_name": "last_name",
      "field_value": "张"
    },
    {
      "field_name": "account_no",
      "field_value": "622306492053631234"
    },
    {
      "field_name": "cert_number",
      "field_value": "330100198001010001"
    }
  ]
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payment_method": "BANK_TRANSFER",
    "form_fields": [],
    "account_id": "string",
    "account_no": "622306492053631234",
    "status": "PENDING",
    "require_clearing_network": true,
    "currencies": ["CNY"]
  }
}
```

### 11. 更新收款人账户（银行账户）
- **请求方式**：PUT
- **接口路径**：`/remittance/payee_accounts/{account_id}`
- **接口描述**：通过表单更新银行账户信息（1.7新增）
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | account_id | 是 | string | 收款账户id |
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | currency | 否 | string | 账户币种，与表单获取时一致 |
  | fields | 是 | Array<object> | 表单字段集合（需完整提交） |
- **响应结果（200 OK）**：同「新增收款人收款账户（银行账户）」响应结构

### 12. 移除收款账户
- **请求方式**：DELETE
- **接口路径**：`/remittance/payee_accounts/{account_id}`
- **接口描述**：删除指定收款账户（仅受理请求，实际业务结果需以查询为准）
- **请求头参数**：同「新增收款人」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | account_id | 是 | string | 收款账户id |
- **响应结果（200 OK）**：同「停用收款人」响应结构

### 13. 查询支持付款国家和币种列表
- **请求方式**：GET
- **接口路径**：`/remittance/payout_currencies`
- **接口描述**：查询指定付款方式支持的国家和币种（1.11新增）
- **请求头参数**：同「新增收款人」
- **查询参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payment_method | 是 | string | 付款方式





  # 对外付款管理接口文档
## 一、付款人管理
### 概述
用于管理对外付款的付款人信息，包含附件上传、付款人新增/查询/更新/删除等功能。

### 1. 上传附件
- **请求方式**：POST
- **接口路径**：`/attachments`
- **接口描述**：上传付款人证件类附件，仅支持图片格式（jpg/png）
- **请求头参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | x-appid | 是 | string | 从Portal获取的appid（参考Signature Guide） |
  | authorization | 是 | string | 请求签名（参考Signature Guide） |
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | type | 是 | string | 附件类型，固定值：CERT_FILE（证件） |
  | filename | 是 | string | 文件名 |
  | base64 | 是 | string | Base64编码的文件内容 |
- **请求示例**
```json
{
  "type": "CERT_FILE",
  "filename": "passport.jpg",
  "base64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAAAAAAAD/..."
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "file_id": "file_123456789"
}
```

### 2. 新增付款人
- **请求方式**：POST
- **接口路径**：`/remittance/payers`
- **接口描述**：创建新的付款人信息，支持个人/企业两种主体类型
- **请求头参数**：同「上传附件」
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | subject_type | 是 | string | 主体类型，枚举值：INDIVIDUAL（个人）、ENTERPRISE（企业） |
  | first_name | 否 | string | 付款人名（个人类型必填），仅允许英文字母和空格 |
  | last_name | 否 | string | 付款人姓（个人类型必填），仅允许英文字母和空格 |
  | cert_type | 是 | string | 证件类型，枚举值：PASSPORT（护照）、DRIVING_LICENSE（驾照）、ID_CARD（身份证）、BUSINESS_LICENSE（营业执照） |
  | cert_number | 是 | string | 证件号码 |
  | cert_files | 是 | Array<string> | 证件附件id列表（通过附件上传API获取） |
  | birthday | 否 | string | 生日（个人类型必填），ISO-8601日期格式（如1991-06-22） |
  | region | 是 | string | 国家/地区编码（ISO-3166 2字符） |
  | company_name | 否 | string | 公司名称（企业类型必填），仅允许英文字母、空格和数字 |
  | register_number | 否 | string | 公司注册号（企业类型必填） |
  | business_scopes | 否 | Array<string> | 企业业务类型（企业类型必填），枚举值见下方 |
  | address | 是 | object | 付款人地址信息，包含postcode/address/city/state/country字段 |
- **企业业务类型枚举**
  | 枚举值 | 描述 |
  |--------|------|
  | MISCELLANEOUS_SERVICES | 杂项事务 |
  | AUTOMOBILE_RENTAL_SERVICES | 汽车租赁服务 |
  | RESTAURANTS_LEISURE | 餐厅休闲 |
  | PROFESSIONAL_CONSULTING | 专业咨询 |
  | EDUCATION | 教育 |
  | DATA_PROCESSING_SERVICES | 数据处理服务 |
  | HUMAN_RESOURCE_EMPLOYMENT_SERVICES | 人力资源就业服务 |
  | ENVIRONMENTAL_FACILITIES_SERVICES | 环境设施服务 |
  | OTHER_SERVICES | 其他服务 |
  | AGRICULTURAL | 农业 |
  | FORESTRY | 林业 |
  | FISHING_HUNTING_AND_TRAPPING | 渔猎和诱捕 |
  | TRANSPORTATION | 运输 |
  | LOGISTICS_WAREHOUSING | 物流仓储 |
  | AIRLINES_AIR_CARRIERS | 航空公司 |
  | TRAVEL_ACCOMMODATION | 旅行住宿 |
  | AUTOMOBILES_AND_VEHICLES | 汽车和车辆 |
  | OFFICE_SUPPLIES | 办公用品 |
  | DISTRIBUTORS | 分销商 |
  | APPAREL_RETAIL | 服装零售 |
  | COMPUTER_ELECTRONICS_RETAIL | 电脑电子零售 |
  | HOME_IMPROVEMENT_HOMEFURNISHING_RETAIL | 家居装修家居零售 |
  | CULTURE_AMUSEMENT_PETS | 文化娱乐宠物 |
  | OTHER_RETAIL | 其他零售 |
  | CONSTRUCTION_MATERIALS | 建筑材料 |
  | CONTAINERS_PACKAGING | 集装箱包装 |
  | BUILDING_PRODUCTS | 建筑产品 |
  | CONSTRUCTION_ENGINEERING | 建筑工程 |
  | ELECTRICAL_EQUIPMENT | 电气设备 |
  | INDUSTRIAL_CONGLOMERATES | 工业集团 |
  | MACHINERY | 机械 |
  | TRADING_COMPANIES_DISTRIBUTORS | 贸易公司经销商 |
  | AUTOMOBILE_COMPONENTS | 汽车零部件 |
  | AUTOMOBILES | 汽车 |
  | HOUSEHOLD_DURABLES | 家庭耐用品 |
  | LEISURE_PRODUCTS | 休闲产品 |
  | TEXTILES_APPAREL_LUXURY_GOODS | 纺织品服装奢侈品 |
  | CONSUMER_STAPLES_DISTRIBUTION_RETAIL | 消费必需品分销零售 |
  | BEVERAGES | 饮料 |
  | FOOD_PRODUCTS | 食品 |
  | HOUSEHOLD_PRODUCTS | 家用产品 |
  | PERSONAL_CARE_PRODUCTS | 个人护理产品 |
  | HEALTH_CARE_TECHNOLOGY | 医疗保健技术 |
  | BIOTECHNOLOGY | 生物技术 |
  | SOFTWARE | 软件 |
  | TECHNOLOGY_HARDWARE_STORAGE_PERIPHERALS | 技术、硬件存储、外围设备 |
  | ELECTRONIC_EQUIPMENT_INSTRUMENTS_COMPONENTS | 电子设备、仪器元件 |
  | SEMICONDUCTORS_SEMICONDUCTOR_EQUIPMENT | 半导体设备 |
  | MEDIA | 媒体 |
  | ENTERTAINMENT | 娱乐 |
  | INTERACTIVE_MEDIA_SERVICES | 互动媒体服务 |
  | OTHERS | 其他 |
- **请求示例（个人类型）**
```json
{
  "subject_type": "INDIVIDUAL",
  "first_name": "John",
  "last_name": "Doe",
  "cert_type": "PASSPORT",
  "cert_number": "P12345678",
  "cert_files": ["file_123456789"],
  "birthday": "1991-06-22",
  "region": "US",
  "address": {
    "postcode": "10001",
    "address": "123 Main St",
    "city": "New York",
    "state": "NY",
    "country": "US"
  }
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payer_id": "payer_123456789",
    "subject_type": "INDIVIDUAL",
    "first_name": "John",
    "last_name": "Doe",
    "cert_type": "PASSPORT",
    "cert_number": "P12345678",
    "cert_files": ["file_123456789"],
    "birthday": "1991-06-22",
    "region": "US",
    "company_name": "",
    "register_number": "",
    "business_scopes": [],
    "address": {
      "postcode": "10001",
      "address": "123 Main St",
      "city": "New York",
      "state": "NY",
      "country": "US"
    }
  }
}
```

### 3. 列出付款人列表
- **请求方式**：GET
- **接口路径**：`/remittance/payers`
- **接口描述**：查询所有付款人信息列表
- **请求头参数**：同「上传附件」
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payers": [
      {
        "payer_id": "payer_123456789",
        "subject_type": "INDIVIDUAL",
        "first_name": "John",
        "last_name": "Doe",
        "region": "US"
      }
    ]
  }
}
```

### 4. 查看付款人详情
- **请求方式**：GET
- **接口路径**：`/remittance/payers/{payer_id}`
- **接口描述**：查询指定付款人的详细信息
- **请求头参数**：同「上传附件」
- **路径参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payer_id | 是 | string | 付款人id |
- **响应结果（200 OK）**：同「新增付款人」响应结构

### 5. 更新付款人信息
- **请求方式**：PUT
- **接口路径**：`/remittance/payers/{payer_id}`
- **接口描述**：全量更新付款人信息（需提交所有字段，而非仅修改项）
- **请求头参数**：同「上传附件」
- **路径参数**：同「查看付款人详情」
- **请求体参数**：同「新增付款人」
- **响应结果（200 OK）**：同「新增付款人」响应结构

### 6. 移除付款人信息
- **请求方式**：DELETE
- **接口路径**：`/remittance/payers/{payer_id}`
- **接口描述**：删除指定付款人信息（仅受理请求，实际结果需查询确认）
- **请求头参数**：同「上传附件」
- **路径参数**：同「查看付款人详情」
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  }
}
```
- **响应说明**
  - `result`枚举值：
    - S：请求接收成功（需查询确认业务结果）
    - F：请求接收失败（无需重试）
    - U：未知异常（需后续查询确认）

## 二、对外付款
### 概述
用于对外付款的申请、锁汇、订单提交及查询等操作，包含清算网络查询、锁汇申请、订单管理等功能。

### 1. 查询可用清算网络
- **请求方式**：GET
- **接口路径**：`/remittance/clearing_networks`
- **接口描述**：查询收款账户可用的清算网络并估算手续费（1.7新增），收款账户`require_clearing_network=true`时必填
- **请求头参数**：同「上传附件」
- **查询参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_account_id | 是 | string | 收款人账户id |
  | pay_currency | 是 | string | 付款币种（ISO-4217） |
  | amount | 是 | number | 金额（IDR不支持小数） |
  | amount_type | 是 | string | 金额类型，枚举值：PAY_AMOUNT（支付总金额）、RECEIVE_AMOUNT（收款人实收金额） |
  | receive_currency | 是 | string | 收款币种（ISO-4217） |
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "clearing_networks": [
    {
      "network": "SWIFT",
      "fee": 10.0,
      "estimated_arrival_time": "2024-06-08T12:00:00Z"
    }
  ],
  "payee_account_id": "account_123456789",
  "pay_currency": "USD",
  "receive_currency": "CNY",
  "amount": 1000.0,
  "amount_type": "PAY_AMOUNT",
  "country": "CN"
}
```

### 2. 申请锁汇
- **请求方式**：POST
- **接口路径**：`/remittance/quotes`
- **接口描述**：申请对外付款的汇率锁定，生成锁汇单
- **请求头参数**：同「上传附件」
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | payee_account_id | 是 | string | 收款人账户id |
  | payer_id | 否 | string | 付款人id（默认使用商户主体）；中国个人银行账户收款时付款人必须为个人；Alipay付款时企业付款人需提前申请B2C权限 |
  | purpose | 是 | string | 汇款目的，枚举值参考业务字典；中国个人银行账户收款时仅允许SALARY |
  | pay_currency | 是 | string | 付款币种（ISO-4217） |
  | receive_currency | 否 | string | 收款币种（ISO-4217，1.7新增）；默认收款人币种，收款人多币种时必填 |
  | amount | 是 | number | 金额（IDR不支持小数） |
  | amount_type | 是 | string | 金额类型，枚举值：PAY_AMOUNT、RECEIVE_AMOUNT |
  | clearing_network | 否 | string | 清算网络（SWIFT/ACH等），收款账户`require_clearing_network=true`时必填 |
  | aba_number | 否 | string | ABA码（ACH/FedWire网络必填，1.7新增） |
  | fps_bank_id | 否 | string | FPS码（FPS网络必填，1.7新增） |
  | ifs_code | 否 | string | IFS码（印度银行账户必填，1.7新增） |
  | intermediary_swift_code | 否 | string | 中间行Swift码（Swift网络可选，1.7新增） |
  | remark | 否 | string | 汇款备注，仅允许字母、数字、空格，最长100字符 |
- **请求示例**
```json
{
  "payee_account_id": "account_123456789",
  "payer_id": "payer_123456789",
  "purpose": "SALARY",
  "pay_currency": "USD",
  "receive_currency": "CNY",
  "amount": 1000.0,
  "amount_type": "PAY_AMOUNT",
  "clearing_network": "SWIFT",
  "remark": "Salary Payment"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "quote_id": "quote_123456789",
    "payment_method": "BANK_TRANSFER",
    "pay_amount": {
      "currency": "USD",
      "amount": 1000.0
    },
    "receive_amount": {
      "currency": "CNY",
      "amount": 7200.0
    },
    "surcharge": {
      "currency": "USD",
      "amount": 10.0
    },
    "exchange_rate": 7.2,
    "expire_at": "2024-06-07T11:23:22Z"
  }
}
```

### 3. 提交付款订单
- **请求方式**：POST
- **接口路径**：`/remittance/orders`
- **接口描述**：基于锁汇单提交付款订单
- **请求头参数**：同「上传附件」
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | quote_id | 是 | string | 锁汇id |
  | client_order_id | 是 | string | 客户系统唯一订单id |
- **请求示例**
```json
{
  "quote_id": "quote_123456789",
  "client_order_id": "order_987654321"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "order_id": "order_123456789",
    "order_source": "API",
    "client_order_id": "order_987654321",
    "create_time": "2024-06-07T12:33:23Z",
    "finish_time": "",
    "status": "CREATED",
    "payee_id": "payee_123456789",
    "payee_account_id": "account_123456789",
    "payment_method": "BANK_TRANSFER",
    "pay_amount": {
      "currency": "USD",
      "amount": 1000.0
    },
    "receive_amount": {
      "currency": "CNY",
      "amount": 7200.0
    },
    "surcharge": {
      "currency": "USD",
      "amount": 10.0
    },
    "exchange_rate": 7.2,
    "error_message": ""
  }
}
```

### 4. 查询付款单列表
- **请求方式**：GET
- **接口路径**：`/remittance/orders`
- **接口描述**：分页查询付款订单列表，支持多条件筛选
- **请求头参数**：同「上传附件」
- **查询参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | page | 是 | integer | 页码，从1开始，默认值：1 |
  | limit | 是 | integer | 每页条数，默认值：20 |
  | payee_id | 否 | string | 收款人id |
  | payer_id | 否 | string | 付款人id |
  | time_start | 否 | string | 查询起始时间（含），ISO-8601格式（如2024-04-03T12:00:00Z） |
  | time_end | 否 | string | 查询截止时间（不含），ISO-8601格式（如2024-05-03T12:00:00Z） |
  | order_id | 否 | string | 平台付款单号 |
  | client_order_id | 否 | string | 客户付款单号 |
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "query": {
      "page": 1,
      "limit": 20
    },
    "orders": [
      {
        "order_id": "order_123456789",
        "status": "SUCCESS",
        "pay_amount": {
          "currency": "USD",
          "amount": 1000.0
        },
        "create_time": "2024-06-07T12:33:23Z"
      }
    ],
    "page": 1,
    "limit": 20,
    "total_count": 1,
    "total_page": 1
  }
}
```

# 收单服务接口文档
## 一、概述
GSalary 集成 Antom 全球收单能力，支持收银台支付、钱包授权支付（EasySafePay）、卡授权支付（CardAutoDebit）等多种收款方式，适用于桌面网站、移动网站、移动应用等多终端场景。提供两种集成方案：Checkout Page 集成和 SDKs 集成，开发者可根据业务需求选择。

## 二、支付相关接口
### 2.1 支付咨询
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/pay_consult`
- **接口描述**：查询可用支付方式、限额、国家/货币支持等信息，支持自动化选择和排序支付方式
- **请求头参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | x-appid | 是 | string | 从 Portal 获取的 appid（参考 Signature Guide） |
  | authorization | 是 | string | 请求签名（参考 Signature Guide） |
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | mch_app_id | 是 | string | 商户应用 ID（开通收单业务后从商户后台获取） |
  | payment_currency | 是 | string | 支付币种（ISO 4217 三位代码，如 USD） |
  | payment_amount | 是 | number | 支付金额（单位：元） |
  | settlement_currency | 是 | string | 结算币种（ISO 4217 三位代码，如 USD） |
  | allowed_payment_method_regions | 否 | Array<string> | 允许的支付方式所属国家/地区（ISO 两位代码） |
  | allowed_payment_methods | 否 | Array<string> | 允许的支付方式列表 |
  | user_region | 否 | string | 用户所在国家/地区（ISO 两位代码），用于支付方式排序 |
  | env_terminal_type | 是 | string | 终端类型，枚举值：WEB、WAP、APP、MINI_APP |
  | env_os_type | 否 | string | 操作系统类型，枚举值：IOS、ANDROID |
  | env_client_ip | 否 | string | 客户端设备 IP 地址 |
- **请求示例**
```json
{
  "mch_app_id": "123456789",
  "payment_currency": "USD",
  "payment_amount": 10,
  "settlement_currency": "USD",
  "allowed_payment_method_regions": [],
  "allowed_payment_methods": [],
  "user_region": "US",
  "env_terminal_type": "APP",
  "env_os_type": "ANDROID",
  "env_client_ip": "192.168.1.1"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payment_options": [
      {
        "payment_method_type": "ALIPAY_CN",
        "currency": "USD",
        "limit": {
          "min": 0.01,
          "max": 10000
        },
        "country": "CN"
      }
    ]
  }
}
```

### 2.2 支付会话创建（收银台）
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/pay_session`
- **接口描述**：创建收银台支付会话，返回加密会话数据用于初始化客户端 SDK
- **请求头参数**：同「支付咨询」
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | mch_app_id | 是 | string | 商户应用 ID |
  | payment_request_id | 是 | string | 商户自定义支付请求 ID（唯一） |
  | payment_currency | 是 | string | 支付币种（ISO 4217 三位代码） |
  | payment_amount | 是 | number | 支付金额（单位：元） |
  | payment_method_type | 是 | string | 支付方式类型（如 ALIPAY_CN、WECHATPAY_CN，参考支付方式列表） |
  | payment_session_expiry_time | 否 | string | 会话过期时间（ISO 8601 格式，默认 1 小时） |
  | payment_redirect_url | 是 | string | 支付完成后用户重定向地址 |
  | order | 是 | object | 订单信息，包含 reference_order_id、order_description 等字段 |
  | settlement_currency | 是 | string | 结算币种（ISO 4217 三位代码） |
  | env_client_ip | 否 | string | 客户端设备 IP 地址 |
  | product_scene | 是 | string | 产品场景，枚举值：CHECKOUT_PAYMENT（Checkout Page 集成）、ELEMENT_PAYMENT（Element 集成） |
- **请求示例**
```json
{
  "mch_app_id": "1234567890",
  "payment_request_id": "PAY_9876543210",
  "payment_currency": "USD",
  "payment_amount": 10,
  "payment_method_type": "ALIPAY_CN",
  "payment_session_expiry_time": "2023-10-01T12:00:00Z",
  "payment_redirect_url": "https://merchant.com/callback",
  "order": {
    "reference_order_id": "ORD_1234567890",
    "order_description": "Order for product XYZ",
    "order_currency": "USD",
    "order_amount": 10,
    "order_buyer_id": "BUY_1234567890",
    "order_buyer_email": "buyer@example.com"
  },
  "settlement_currency": "USD",
  "env_client_ip": "192.168.1.1",
  "product_scene": "CHECKOUT_PAYMENT"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payment_session_data": "encrypted_session_data",
    "payment_session_expiry_time": "2023-04-30T12:34:56Z",
    "payment_session_id": "session_123456",
    "normal_url": "https://payment-gateway.com/pay"
  }
}
```

## 三、钱包授权支付（EasySafePay）
### 3.1 接入步骤
1. 第一次支付：调用「支付会话创建」接口，获取支付成功状态
2. 接收授权 Token 通知（`ACQUIRING_AUTH_TOKEN`），获取 `access_token`
3. 第二次支付：调用「支付」接口，传入 `access_token` 完成代扣
4. 关注 `access_token` 有效期，过期前调用「刷新令牌」接口更新
5. 提供用户取消授权入口，调用「取消授权」接口失效令牌
6. 接收授权取消消息（`ACQUIRING_AUTH_TOKEN` 推送）

### 3.2 支付会话创建（第一次支付）
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/easy_safe_pay/pay_session`
- **接口描述**：创建钱包授权支付会话（首次支付），返回 SDK 初始化所需会话数据
- **请求头参数**：同「支付咨询」
- **请求体参数**：在「支付会话创建（收银台）」基础上新增以下参数
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | auth_state | 是 | string | 商户生成的快捷授权专属 ID（仅首次支付需指定） |
  | user_login_id | 否 | string | 用户支付方式登录 ID（邮箱/手机号，可选，免用户手动输入） |
- **请求示例**：参考「支付会话创建（收银台）」，新增 `auth_state` 和 `user_login_id` 字段
- **响应结果**：同「支付会话创建（收银台）」

### 3.3 支付（第二次支付）
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/easy_safe_pay/pay`
- **接口描述**：使用授权 `access_token` 发起代扣支付
- **请求头参数**：同「支付咨询」
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | mch_app_id | 是 | string | 商户应用 ID |
  | payment_request_id | 是 | string | 商户自定义支付请求 ID（唯一） |
  | payment_currency | 是 | string | 支付币种（ISO 4217 三位代码） |
  | payment_amount | 是 | number | 支付金额（单位：元） |
  | payment_method_id | 是 | string | 支付方式 ID（从授权通知中获取的 `access_token`） |
  | payment_method_type | 是 | string | 支付方式类型（如 ALIPAY_CN） |
  | payment_redirect_url | 是 | string | 支付完成后重定向地址 |
  | order | 是 | object | 订单信息（同「支付会话创建」） |
  | settlement_currency | 是 | string | 结算币种（ISO 4217 三位代码） |
  | env_client_ip | 否 | string | 客户端设备 IP 地址 |
  | payment_expiry_time | 否 | string | 支付有效期（ISO 8601 格式，默认 14 分钟，银行转账默认 48 小时） |
  | env_terminal_type | 是 | string | 终端类型，枚举值：WEB、WAP、APP、MINI_APP |
  | env_os_type | 否 | string | 操作系统类型，枚举值：IOS、ANDROID |
- **请求示例**
```json
{
  "mch_app_id": "1234567890",
  "payment_request_id": "PAY_9876543210",
  "payment_currency": "USD",
  "payment_amount": 10,
  "payment_method_id": "access_token_123456",
  "payment_method_type": "ALIPAY_CN",
  "payment_redirect_url": "https://merchant.com/callback",
  "order": {
    "reference_order_id": "ORD_1234567890",
    "order_description": "Order for product XYZ",
    "order_currency": "USD",
    "order_amount": 10
  },
  "settlement_currency": "USD",
  "env_client_ip": "192.168.1.1",
  "payment_expiry_time": "2019-08-24T14:15:22Z",
  "env_terminal_type": "WEB",
  "env_os_type": "IOS"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "payment_request_id": "PAY_9876543210",
    "payment_id": "PAY_123456789",
    "payment_currency": "USD",
    "payment_amount": 10,
    "normal_url": "https://payment-gateway.com/pay",
    "scheme_url": "alipay://pay",
    "applink_url": "https://app.alipay.com/pay",
    "app_identifier": "com.alipay"
  }
}
```

### 3.4 刷新令牌
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/auth_refresh_token`
- **接口描述**：使用刷新令牌获取新的访问令牌（`access_token` 即将过期时调用）
- **请求头参数**：同「支付咨询」
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | mch_app_id | 是 | string | 商户应用 ID |
  | refresh_token | 是 | string | 刷新令牌（从授权通知中获取） |
  | merchant_region | 否 | string | 商户业务国家/地区（ISO 两位代码，支持 US、JP、PK、SG） |
- **请求示例**
```json
{
  "mch_app_id": "1234567890",
  "refresh_token": "refresh_token_123456",
  "merchant_region": "US"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  },
  "data": {
    "access_token": "new_access_token_123456",
    "access_token_expiry_time": "2023-12-31T23:59:59Z",
    "refresh_token": "new_refresh_token_123456",
    "refresh_token_expiry_time": "2024-01-31T23:59:59Z",
    "user_login_id": "buyer@example.com"
  }
}
```

### 3.5 取消授权
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/auth_revoke_token`
- **接口描述**：取消用户对商户的授权，`access_token` 和 `refresh_token` 失效
- **请求头参数**：同「支付咨询」
- **请求体参数**
  | 参数名 | 是否必填 | 类型 | 说明 |
  |--------|----------|------|------|
  | mch_app_id | 是 | string | 商户应用 ID |
  | access_token | 是 | string | 访问令牌 |
- **请求示例**
```json
{
  "mch_app_id": "1234567890",
  "access_token": "access_token_123456"
}
```
- **响应结果（200 OK）**
```json
{
  "result": {
    "result": "S",
    "code": "string",
    "message": "string"
  }
}
```
- **响应说明**
  - `result` 枚举值：
    - S：请求接收成功（需查询确认业务结果）
    - F：请求接收失败（无需重试）
    - U：未知异常（需后续查询确认）

## 四、卡授权支付（CardAutoDebit）
### 4.1 接入步骤
1. 第一次支付：调用「支付会话创建」接口并完成支付
2. 从支付结果通知中获取 `card_token`（`payment_result_info` 字段）
3. 第二次支付：调用「支付」接口，传入 `card_token` 完成代扣（仅支持全球卡）

### 4.2 支付会话创建（第一次支付）
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/card_auto_debit/pay_session`
- **接口描述**：创建卡授权支付会话（首次支付），返回 SDK 初始化所需会话数据
- **请求头参数**：同「支付咨询」
- **请求体参数**：与「支付会话创建（收银台）」完全一致
- **请求示例**：参考「支付会话创建（收银台）」
- **响应结果**：同「支付会话创建（收银台）」

### 4.3 支付（第二次支付）
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/card_auto_debit/pay`
- **接口描述**：使用 `card_token` 发起代扣支付
- **请求头参数**：同「支付咨询」
- **请求体参数**：与「钱包授权支付（第二次支付）」完全一致，仅 `payment_method_id` 为 `card_token`
- **请求示例**：参考「钱包授权支付（第二次支付）」，`payment_method_id` 传入 `card_token`
- **响应结果**：同「钱包授权支付（第二次支付）」

## 五、支付后相关接口
### 5.1 取消支付
- **请求方式**：POST
- **接口路径**：`/gateway/v1/acquiring/cancel`
- **接口描述**：取消未完成的支付（超出合同约定可取消期限则无法取消）
- **请求头参数**：同