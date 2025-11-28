#!/usr/bin/env python3
"""
生成RSA公钥和私钥对
用于G-Salary API鉴权
"""

from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
import os

def generate_rsa_keys():
    """生成RSA密钥对"""
    # 生成私钥
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048,
        backend=default_backend()
    )
    
    # 导出私钥 (PEM格式)
    private_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption()
    )
    
    # 导出公钥 (PEM格式)
    public_key = private_key.public_key()
    public_pem = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )
    
    # 保存私钥到文件
    with open('private_key.pem', 'wb') as f:
        f.write(private_pem)
    print("✓ 私钥已保存到: private_key.pem")
    
    # 保存公钥到文件
    with open('public_key.pem', 'wb') as f:
        f.write(public_pem)
    print("✓ 公钥已保存到: public_key.pem")
    
    print("\n" + "="*60)
    print("公钥内容 (Public Key):")
    print("="*60)
    print(public_pem.decode('utf-8'))
    
    print("\n" + "="*60)
    print("私钥内容 (Private Key):")
    print("="*60)
    print(private_pem.decode('utf-8'))
    
    print("\n⚠️  请妥善保管私钥，不要泄露给他人！")
    print("💡 将公钥提供给G-Salary API进行鉴权配置")

if __name__ == "__main__":
    generate_rsa_keys()
