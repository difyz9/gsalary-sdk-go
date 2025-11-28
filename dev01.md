开发者需自行生成RSA密钥对，打开GSalary Portal或GSalary Test Portal（测试环境）， 在开发者公钥配置面板粘贴生成的公钥无需去除头尾的public key标签，可以完全复制粘贴到后台商户公钥配置项，并将对应私钥妥善配置到对接系统中。

随后在GSalary Portal可以得到平台分配的服务器公钥，同样配置到对接系统中。完成双方的公钥交换。

Bash Demo:

# Generate a rsa private key with 4096 bit length
openssl genrsa 4096 > private.pem
# Generate rsa public key from private key
openssl rsa -in private.pem -pubout -out public.pem
# Convert private key into pkcs8 format.(Some language such as java require this format)
openssl pkcs8 -topk8 -in private.pem -nocrypt -out private-pkcs8.pem