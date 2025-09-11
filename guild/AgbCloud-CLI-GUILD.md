# AgbCloud CLI 使用教程

本教程将指导您如何使用 AgbCloud CLI 工具进行镜像管理操作。

## 目录

- [前提条件](#前提条件)
- [1. 登录认证](#1-登录认证)
- [2. 创建镜像](#2-创建镜像)
- [3. 激活镜像](#3-激活镜像)
- [4. 停止镜像](#4-停止镜像)
- [5. 镜像列表](#5-镜像列表)
- [常见问题](#常见问题)

## 前提条件

在开始使用之前，请确保：
- 已安装 AgbCloud CLI 工具
- 拥有有效的 AgbCloud 账户
- 网络连接正常

## 1. 登录认证

在使用任何镜像管理功能之前，您需要先登录到 AgbCloud。

### 命令语法

```bash
agbcloud login
```

### 使用步骤

1. **执行登录命令**：
   ```bash
   agbcloud login
   ```

2. **系统响应**：
   ```
   🔐 Starting AgbCloud authentication...
   📡 Using callback port: 8080
   🌐 Requesting OAuth login URL...
   ✅ Successfully retrieved OAuth URL!
   📋 Request ID: req-xxxxx
   🔍 Trace ID: trace-xxxxx
   
   🚀 Starting local callback server on port 8080...
   🔗 OAuth URL:
     https://agb.cloud/oauth/authorize?...
   
   🌐 Opening the browser for authentication...
   ```

3. **浏览器认证**：
   - CLI 会自动打开浏览器
   - 如果浏览器未自动打开，请手动复制 URL 到浏览器
   - 在浏览器中完成 Google 账户认证

4. **认证成功**：
   ```
   ✅ Authentication successful!
   🔑 Received authorization code: abcd1234...
   🔄 Exchanging authorization code for access token...
   ✅ Login successful!
   ```

### 注意事项

- 登录会话有效期为一定时间，过期后需要重新登录
- 登录信息会安全存储在本地配置文件中

## 2. 创建镜像

创建自定义镜像需要提供 Dockerfile 和基础镜像 ID。

### 命令语法

```bash
agbcloud image create <镜像名称> --dockerfile <Dockerfile路径> --imageId <基础镜像ID>
```

### 参数说明

- `<镜像名称>`：自定义镜像的名称（必需）
- `--dockerfile, -f`：Dockerfile 文件路径（必需）
- `--imageId, -i`：基础镜像 ID（必需）

### 使用示例

```bash
# 完整命令
agbcloud image create myCustomImage --dockerfile ./Dockerfile --imageId agb-code-space-1

# 使用短参数
agbcloud image create myCustomImage -f ./Dockerfile -i agb-code-space-1
```

### 执行流程

1. **开始创建**：
   ```
   🏗️  Creating image 'myCustomImage'...
   📡 Getting upload credentials...
   ✅ Upload credentials obtained (Task ID: task-xxxxx)
   ```

2. **上传 Dockerfile**：
   ```
   📤 Uploading Dockerfile...
   ✅ Dockerfile uploaded successfully
   ```

3. **创建镜像**：
   ```
   🔨 Creating image...
   ✅ Image creation initiated
   ```

4. **监控进度**：
   ```
   ⏳ Monitoring image creation progress...
   📊 Status: Creating
   📊 Status: Available
   ✅ Image creation completed successfully!
   ```

### 镜像状态说明

- **Creating**：镜像正在创建中
- **Create Failed**：镜像创建失败
- **Available**：镜像创建完成，可以使用

## 3. 激活镜像

激活镜像会启动一个运行实例，您可以指定 CPU 和内存资源。

### 命令语法

```bash
agbcloud image activate <镜像ID> [--cpu <核心数>] [--memory <内存GB>]
```

### 参数说明

- `<镜像ID>`：要激活的镜像 ID（必需）
- `--cpu, -c`：CPU 核心数（可选）
- `--memory, -m`：内存大小，单位 GB（可选）

### 使用示例

```bash
# 基本激活
agbcloud image activate img-7a8b9c1d0e

# 指定资源配置
agbcloud image activate img-7a8b9c1d0e --cpu 2 --memory 4

# 使用短参数
agbcloud image activate img-7a8b9c1d0e -c 2 -m 4
```

### 执行流程

1. **开始激活**：
   ```
   🚀 Activating image 'img-7a8b9c1d0e'...
   💾 CPU: 2 cores, Memory: 4 GB
   🔍 Checking current image status...
   ```

2. **状态检查**：
   ```
   📊 Current Status: Available
   ✅ Image is available, proceeding with activation...
   🔄 Starting image activation...
   ```

3. **激活成功**：
   ```
   ✅ Image activation initiated successfully!
   📊 Operation Status: true
   🔍 Request ID: req-xxxxx
   ```

4. **监控激活状态**：
   ```
   ⏳ Monitoring image activation status...
   📊 Status: Activating
   📊 Status: Activated
   ✅ Image activation completed successfully!
   ```

### 镜像激活状态说明

- **Available**：镜像可用，未激活
- **Activating**：镜像正在激活中
- **Activated**：镜像已激活，正在运行
- **Activate Failed**：镜像激活失败
- **Ceased Billing**：镜像已停止计费

### 特殊情况处理

- 如果镜像已经激活，系统会显示当前状态
- 如果镜像正在激活中，会自动加入监控流程
- 如果镜像处于失败状态，会尝试重新激活

## 4. 停止镜像

停止（停用）正在运行的镜像实例。

### 命令语法

```bash
agbcloud image deactivate <镜像ID>
```

### 参数说明

- `<镜像ID>`：要停止的镜像 ID（必需）

### 使用示例

```bash
agbcloud image deactivate img-7a8b9c1d0e
```

### 执行流程

1. **开始停止**：
   ```
   🛑 Deactivating image 'img-7a8b9c1d0e'...
   🔄 Deactivating image instance...
   ```

2. **停止成功**：
   ```
   ✅ Image deactivation initiated successfully!
   📊 Operation Status: true
   🔍 Request ID: req-xxxxx
   ```

### 注意事项

- 停止镜像会终止正在运行的实例
- 停止后的镜像状态会变为 "Available"
- 停止操作通常会立即生效

## 5. 镜像列表

查看您的镜像列表，支持分页和类型筛选。

### 命令语法

```bash
agbcloud image list [--type <类型>] [--page <页码>] [--size <每页数量>]
```

### 参数说明

- `--type, -t`：镜像类型，可选值：
  - `User`：用户自定义镜像（默认）
  - `System`：系统基础镜像
- `--page, -p`：页码，默认为 1
- `--size, -s`：每页显示数量，默认为 10

### 使用示例

```bash
# 查看用户镜像（默认）
agbcloud image list

# 查看系统镜像
agbcloud image list --type System

# 分页查看
agbcloud image list --page 2 --size 5

# 使用短参数
agbcloud image list -t User -p 1 -s 20
```

### 输出示例

```
📋 Listing User images (Page 1, Size 10)...
🔍 Fetching image list...
✅ Found 3 images (Total: 3)
📄 Page 1 of 1 (Page Size: 10)

IMAGE ID                  IMAGE NAME               STATUS               TYPE            UPDATED AT          
--------                  ----------               ------               ----            ----------          
img-7a8b9c1d0e           myCustomImage            Available            User            2025-01-15 10:30    
img-2f3g4h5i6j           webAppImage              Activated            User            2025-01-15 09:15    
img-8k9l0m1n2o           dataProcessImage         Creating             User            2025-01-15 11:45    
```

### 状态说明

镜像可能处于以下状态：

**创建相关状态：**
- **Creating**：镜像正在创建中
- **Create Failed**：镜像创建失败
- **Available**：镜像创建完成，可以使用

**激活相关状态：**
- **Activating**：镜像正在激活中
- **Activated**：镜像已激活，正在运行
- **Deactivating**：镜像正在停止中
- **Activate Failed**：镜像激活失败
- **Ceased Billing**：镜像已停止计费

## 常见问题

### Q: 如何查看命令帮助？

A: 在任何命令后添加 `--help` 或 `-h` 参数：

```bash
agbcloud --help
agbcloud image --help
agbcloud image create --help
```

### Q: 登录失败怎么办？

A: 请检查：
1. 网络连接是否正常
2. 浏览器是否能正常访问 agb.cloud
3. 是否有有效的 Google 账户
4. 防火墙是否阻止了回调端口

### Q: 镜像创建失败怎么办？

A: 请检查：
1. Dockerfile 语法是否正确
2. 基础镜像 ID 是否有效
3. 网络连接是否稳定
4. 查看错误信息中的 Request ID 以便技术支持

### Q: 如何查看详细的执行信息？

A: 使用 `--verbose` 或 `-v` 参数：

```bash
agbcloud -v image create myImage -f ./Dockerfile -i agb-code-space-1
```

### Q: 镜像激活很慢怎么办？

A: 镜像激活可能需要几分钟时间，特别是：
- 首次激活某个镜像
- 镜像较大
- 系统负载较高

请耐心等待，系统会自动监控激活状态。

### Q: 如何获取基础镜像 ID？

A: 使用镜像列表命令查看系统镜像：

```bash
agbcloud image list --type System
```

---

**技术支持**：如果遇到问题，请联系技术支持团队，并提供相关的 Request ID 和 Trace ID。 