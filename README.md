# Himawari

🌻🌻🌻Web Vulnerability Scanner🌻🌻🌻

![Himawari_Gopher.png](Himawari_Gopher.png)

The Gopher character is based on the Go mascot designed by Renée French.

## 動作保証環境

以下の version で動作することを確認しています。  
2021年1月21日時点でのnodejsの安定版(v16.13.2)で動作できることも確認済みです。  

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

```txt
bash exec.sh
```

※ `Ctrl + C` を押すと  server と WebUI の両方の process が kill されます。

## develop

以下のコマンドを実行すると WebUI をホットリロードで開発できる環境が用意されます。  

```txt
cd webui/
npm run dev
```

※ `Ctrl + C` を押すと  server と WebUI の両方の process が kill されます。

API server 側のコードに変更を加える際には、`make run`コマンドで逐次実行することができます。  

```
cd api/
make run
```
