# Himawari

🌻🌻🌻Web Vulnerability Scanner🌻🌻🌻

## setup

クローンしてセットアップ

```txt
git clone https://github.com/futabato/Himawari.git
cd Himawari/api/
make setup
```

## api server

ビルドして実行

```txt
make run
```

<http://localhost:8080> にHimawariが起動します。

## webui

以下のコマンドを実行すると<http://localhost:3000>にWebUIが起動します。

```txt
npm run dev
```

## exec.sh

以下のコマンドを実行すると、api serverとWebUIが起動します。  
実行にはHimawariのバイナリが必要なので、`api`ディレクトリで`go build`コマンドを実行してください。
apiのコードに変更があった場合は、Himawariのバイナリを更新する必要があります。  

```txt
bash exec.sh
```

Ctrl + Cを押すとapi serverのprocessもkillされます。  
