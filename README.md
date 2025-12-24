

# 本地配置

```shell
cp configs/config.dev.yaml config.yaml
```

# 本地启动

```shell
mkdir test_data
go run cmd/server/main.go -c config.yaml
```

# 健康检查

```shell
curl http://127.0.0.1:6065/health
```

# 常用命令行操作

```shell
# 1. 安装 xq
# macOS
brew install python-yq
brew install libxml2

# Ubuntu/Debian
sudo apt-get install libxml2-utils


# 2. 列出目录（PROPFIND）
curl -s -X PROPFIND \
  -u test:test \
  -H "Depth: 1" \
  http://127.0.0.1:6065/ | xq .

# 或者

curl -X PROPFIND \
  -u test:test \
  -H "Depth: 1" \
  http://127.0.0.1:6065/ | xmllint --format -

3. 上传文件（PUT）

echo "Test content" | curl -X PUT \
  -u test:test \
  --data-binary @- \
  http://127.0.0.1:6065/upload.txt

4. 下载文件（GET）

curl -u test:test \
  http://127.0.0.1:6065/upload.txt

5. 删除文件（DELETE）

curl -X DELETE \
  -u test:test \
  http://127.0.0.1:6065/upload.txt

6. 创建目录（MKCOL）

curl -X MKCOL \
  -u test:test \
  http://127.0.0.1:6065/new

7. 测试错误的密码

curl -u test:wrongpassword \
  http://127.0.0.1:6065/

```

# 常用的客户端操作

```text
MACOS
打开访达 -> 选择前往菜单 -> 连接服务器 -> 输入连接地址 -> 输入用户名和密码
```
