# 0xLibre — Minimalist, Self-Hosted File Hosting  
# 0xLibre — 极简、可自托管的文件托管服务

A lightweight, auditable pastebin-for-files service—designed for **privacy**, **simplicity**, and **decentralization**.  
轻量、可审计的“文件粘贴板”服务，专注**隐私**、**简洁**与**去中心化**。

- 📦 Zero external runtime deps (only Go stdlib + `glog`)  
  零外部运行时依赖（仅 Go 标准库 + `glog`）  
- 🧾 Files stored on disk, no database  
  文件直存磁盘，无需数据库  
- 🔗 Short, random URLs (12-char UUID-like IDs)  
  短链接，12 位随机 ID  
- 🖥️ Browser upload + CLI-friendly (`curl -F 'file=@…'`)  
  支持网页与命令行上传  

> “Centralization is bad. Clone it.”  
> “中心化有害。欢迎克隆。”

---

## Quick Start  
## 快速开始

```bash
go build -o 0xlibre .
./0xlibre -p 5000
```
