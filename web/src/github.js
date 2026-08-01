// 轻量混淆：仓库地址以 XOR + Base64 加密写入源码，运行时解密显示。
// 目的：防止被直接阅读 / 篡改源码中的仓库信息。
// 注意：这仅是「威慑性混淆」，并非真正的安全保护；若要彻底保护源码，
// 应将仓库设为私有或依靠授权协议与法律约束（克隆公开代码本身无法被加密阻止）。

const KEY = 'GIT_BACKUP_OPS_2026'
const ENC = 'Lz0gLzF7bGQyOSsnJTFxUV9fGSM4Lj0tOGwJFBMUEAMKEXFvdX8TAQEd'

function xorBase64Decode(b64, key) {
  const bin = atob(b64)
  let out = ''
  for (let i = 0; i < bin.length; i++) {
    out += String.fromCharCode(bin.charCodeAt(i) ^ key.charCodeAt(i % key.length))
  }
  return out
}

export const GITHUB_REPO_URL = xorBase64Decode(ENC, KEY)
export const GITHUB_REPO_NAME = GITHUB_REPO_URL.split('/').slice(-2).join('/')
