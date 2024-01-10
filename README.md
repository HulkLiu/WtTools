# README

## About

This is the official Wails Vue template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

    
1、 编译前端
    
    cd /frontend
    npn install

2。先初始化 ESData 并开启 docker 

    a) docker run -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.4.2

    b) cd ./initEsData  go run main.go

3、  MySQL 配置

    /internal/config 
    Dsn := "root:root@tcp(localhost:3306)/test"

4. 编译项目 

    
    wails build


## 功能需求
![](show/资源搜索引擎 WtTools.png)