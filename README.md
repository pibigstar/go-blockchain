# go创建一个简单的区块链

## 使用示例

进入到cmd文件夹下
```
cd go-blockchain/src/cmd
```

创建一个新的区块链
```$xslt
go run main.go create -address pibigstar
```
![]()

查看余额

```cgo
go run main.go balance -address pibigstar
```
![]()

转账
```cgo
go run main.go send -from pibigstar -to lei -amount 3
```
![]()

再次查看余额
```cgo
go run main.go balance -address pibigstar
```
![]()

打印整个区块链

```cgo
go run main.go list
```
![]()