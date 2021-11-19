# Himawari

🌻🌻🌻Web Vulnerability Scanner🌻🌻🌻

![Himawari_Gopher.png](Himawari_Gopher.png)

The Gopher character is based on the Go mascot designed by Renée French.

## 動作保証環境

Ubuntu-20.04  
Go 1.17  
Node.js 14.17.5

## setup

```txt
git clone https://github.com/futabato/Himawari.git
cd Himawari/
bash setup.sh
```

## execute

以下のコマンドを実行すると、API Server と WebUI が起動します。  
実行には Himawari のバイナリが必要なので、`api`ディレクトリで`make setup`で`go build`コマンドを実行してください。  
**api のコードに変更があった場合は、Himawari のバイナリを更新する必要があります。**  
(`api`ディレクトリに移り、`make setup`コマンド実行することで`go build`コマンドが実行され、Himawari のバイナリが更新されます。)

```txt
bash exec.sh
```

Ctrl + C を押すと api server の process も kill されます。
