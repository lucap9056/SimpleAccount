SimpleAccount 是一個由Golang編寫的基本的帳戶系統，提供用戶註冊、登入及基本帳戶管理功能。
## 功能

---

### 1. 使用者註冊

使用以下資訊註冊帳戶：
- 使用者名稱 (username)
- 電子郵件 (email)
- 密碼 (password)

### 2. 使用者登入

使用者可以使用以下資訊登入帳戶：
- 電子郵件 (email)
- 密碼 (password)

### 3. 驗證電子郵件

系統會寄送驗證電子郵件至用戶註冊時提供的電子郵件地址。用戶需驗證其電子郵件地址以啟用帳戶。
### 4. 更改帳戶資訊

用戶可以更新以下帳戶資訊：
- 使用者名稱 (username)
- 電子郵件 (email)
- 密碼 (password)

## 前端處理

---

前端檔案需要透過 Nginx 或 Apache 伺服器來處理。建議使用這些伺服器來提供靜態檔案服務。

## 設定檔案 `config.json` 

---

- `listen_port`: 伺服器監聽的port。預設為 80。
- `logs_dir_path`: 日誌檔案的存放目錄。預設為 "logs"。
- `extension_channel_port`: 擴充程式通道的port，請避免對外開放。設定 0 以停用此通道。預設為 81。

#### Database 設定

- `database.source`: 連接資料庫的來源字串。格式為 `user:password@tcp(host:port)/Database`。

#### Email 設定

- `email.email_host`: SMTP 伺服器的主機地址。
- `email.port`: SMTP 伺服器的埠號。預設為 587。
- `email.user`: 發送郵件的使用者帳號。
- `email.password`: 發送郵件的使用者密碼。
- `email.api_host`: SimpleAccount API 主機地址。
- `email.template_files_path`: 電子郵件模板檔案的存放目錄。預設為 "html"。
- `email.verification_duration`: 驗證碼有效時間（分鐘）。預設為 5 分鐘。

#### Auth 設定

- `auth.total_number_of_keys`: 可用金鑰的總數。預設為 10。
- `auth.key_validity_duration`: 金鑰有效期限（天）。預設為 30 天。
- `auth.key_files_path`: 金鑰檔案的存放目錄，請嚴加保管。預設為 "./keys/"。
- `auth.token_validity_duration`: 記號有效期限（天）。預設為 30 天。
- `auth.token_auto_renew_time`: 記號自動續期的時間範圍（天）。預設為最後 6 天。
- `auth.temporary_token_validity_duration`: 臨時記號的有效時間（小時）。預設為 2 小時。

## 開發環境
---
- Ubuntu: 22.04
- Golang: 1.21.5
- MariaDB: 11.2.2

