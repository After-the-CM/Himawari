# Himawari

ð»ð»ð»Web Vulnerability Scannerð»ð»ð»

![Himawari_Gopher.png](Himawari_Gopher.png)

The Gopher character is based on the Go mascot designed by RenÃ©e French.

## åä½ä¿è¨¼ç°å¢

ä»¥ä¸ã® version ã§åä½ãããã¨ãç¢ºèªãã¦ãã¾ãã  
2022å¹´1æ21æ¥æç¹ã§ã®nodejsã®å®å®ç(v16.13.2)ã§åä½ã§ãããã¨ãç¢ºèªæ¸ã¿ã§ãã  

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

ä»¥ä¸ã®ã³ãã³ããå®è¡ããã¨ãAPI Server ã¨ WebUI ãèµ·åãã¾ãã  

```txt
bash exec.sh
```

â» `Ctrl + c` ãæ¼ãã¨ API server ã¨ WebUI ã®ä¸¡æ¹ã® process ã kill ããã¾ãã

## develop

ä»¥ä¸ã®ã³ãã³ããå®è¡ããã¨ WebUI ãããããªã­ã¼ãã§éçºã§ããç°å¢ãç¨æããã¾ãã  

```txt
cd webui/
npm run dev
```

API server å´ã®ã³ã¼ãã«å¤æ´ãå ããéã«ã¯ã`make run`ã³ãã³ãã§éæ¬¡å®è¡ãããã¨ãã§ãã¾ãã  

```
cd api/
make run
```
