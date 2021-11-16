# Himawari

🌻🌻🌻Web Vulnerability Scanner🌻🌻🌻

## setup

クローンしてセットアップ

```txt
git clone https://github.com/futabato/Himawari.git
cd Himawari/api/
make setup
go build
```

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

## exec.sh

以下のコマンドを実行すると、API ServerとWebUIが起動します。  
実行にはHimawariのバイナリが必要なので、`api`ディレクトリで`make setup`で`go build`コマンドを実行してください。
apiのコードに変更があった場合は、Himawariのバイナリを更新する必要があります。  

```txt
bash exec.sh
```

Ctrl + Cを押すとapi serverのprocessもkillされます。 