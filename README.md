# Himawari

🌻🌻🌻Web Vulnerability Scanner🌻🌻🌻

## setup

```txt
git clone https://github.com/futabato/Himawari.git
cd Himawari/
bash setup.sh
```

## execute

以下のコマンドを実行すると、API ServerとWebUIが起動します。  
実行にはHimawariのバイナリが必要なので、`api`ディレクトリで`make setup`で`go build`コマンドを実行してください。
**apiのコードに変更があった場合は、Himawariのバイナリを更新する必要があります。**  

```txt
bash exec.sh
```

Ctrl + Cを押すとapi serverのprocessもkillされます。   

## API Server

以下のコマンドを実行すると<http://localhost:8080> にAPI Serverが起動します。

```txt
make run
```

## WebUI

以下のコマンドを実行すると<http://localhost:3000>にWebUIが起動します。

```txt
npm run dev
```
