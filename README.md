# 汽车摇号中签数据自动采集


[![Build Status](https://travis-ci.org/coderguang/xcxYaohaoServer.svg?branch=master)](https://travis-ci.org/coderguang/xcxYaohaoServer)
![](https://img.shields.io/badge/language-golang-orange.svg)
[![codebeat badge](https://codebeat.co/badges/f44324a8-0342-4607-8ad9-028e78eee224)](https://codebeat.co/projects/github-com-coderguang-xcxyaohaoserver-master)
[![](https://img.shields.io/badge/wp-@royalchen-blue.svg)](https://www.royalchen.com)

## 小程序体验
   ![card](https://github.com/coderguang/img/blob/master/xcxYaohaoServer/card.jpg)
   

## 功能
   * 历次数据自动采集（自动下载相应pdf,解析pdf数据到mysql数据库）
   * 自动爬取最新数据（爬虫**自动采集**每月最新数据，数据发布5分钟内即由爬虫爬取并解析）
   * 爬取**当月**数据后，推送中签数据到指定端口，方便做下一步拓展（短信通知）
   * 提供对外接口,只需输入姓名或编码即可查询中签情况
   * 只要摇号网站的模板是采用 **广州、深圳、杭州** 这种的，都可以自动爬取
   * ps：wx 只通过了我 “**广州小型汽车摇号**” 小程序的申请，其他地区由于主体与地区不符不给我通过，坑。开源，有想搞这方面的可以用起来。
   
## 使用
  1. 首先将sql/内的sql文件导入到数据库中
  2. 配置config/内的文件。
  ```json
     {
        "title":"shenzhen",
        "indexUrl":"http://xqctk.jtys.sz.gov.cn/gbl/",
        "allowUrls":["xqctk.jtys.sz.gov.cn"],
        "ignoreUrls":["http://xqctk.jtys.sz.gov.cn/attachment/2015212/1423707141859.pdf",
            "http://xqctk.jtys.sz.gov.cn"],
        "dbUrl":"your db url",
        "dbPort":"db port",
        "dbUser":"db user",
        "dbPwd":"db passwd",
        "dbName":"db name",
        "dbTable":"data table name",
        "historyTable":"history table name",
        "listenPort":"listen port",
        "finishTxt":"中签详细列表数据完成",
        "timeTxt":"本期编号",
        "totalNumTxt":"指标总数",
        "personTxt":"个人",
        "companyTxt":"单位",
        "normalTxt":"普通",
        "newEngineTxt":["新能源","电动"],
        "pageTxt":["增量指标摇号结果公告","指标配置结果"],
        "resultDate":26,
        "http":"http",
        "noticeUrl":"http://localhost:2000"
    }
   ```
 
  > 其中需要重点关注 title 字段的值，这个值关联到不少地方   
    >> 需要将导入数据库的数据表名字的 template 改为 title 的值   
     >> 对外接口访问时，需要提供该 title 值   
     
  3. 配置完 config文件后，执行  
  
      ```shell
        go run main.go config/config_xxx.json
      ```  
	  
  4. **执行流程图:**
   **启动后的输入如下图**:  
     ![init](https://github.com/coderguang/img/blob/master/xcxYaohaoServer/init.png)  

   **遇到需要采集的pdf数据时，输出如下**：  
     ![get_data](https://github.com/coderguang/img/blob/master/xcxYaohaoServer/get_data.png)  

   **解析完插入数据库**:  
     ![insertDb](https://github.com/coderguang/img/blob/master/xcxYaohaoServer/insertDb.png)  

   **查找完所有数据进入睡眠状态**:  
     ![findEnd](https://github.com/coderguang/img/blob/master/xcxYaohaoServer/findEnd.png)  

   **外部API查询时**:  
     ![search](https://github.com/coderguang/img/blob/master/xcxYaohaoServer/search.png)  


 
## About me

**Author** | _[royalchen](https://www.royalchen.com)_
---------- | -----------------
email  | royalchen@royalchen.com
qq  | royalchen@royalchen.com
website | [www.royalchen.com](https://www.royalchen.com)
