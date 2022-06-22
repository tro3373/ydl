# ydl

## prd
```
make STAGE=prd up
```

```
All request handle

    / all request     :3000
--------------------> server/ydl
```

## dev
```
make up
```

```
SPA asssets embeded in api will be ignored

              /api  proxy    :3000
-- nginx ------------------> server/ydl
            | other (/*)
            |       proxy    :8080
            |--------------> client/back vue
```


---

# TODO
# DONE
- assets empbeded in api for prd
	- [GolangのGin/bindataでシングルバイナリを試してみた(+React) - Qiita](https://qiita.com/wadahiro/items/4173788d54f028936723)
	- [【GO】gin + statikのシングルバイナリファイルサーバ | Narumium Blog](https://blog.narumium.net/2019/06/07/%E3%80%90go%E3%80%91gin-statik%E3%81%AE%E3%82%B7%E3%83%B3%E3%82%B0%E3%83%AB%E3%83%90%E3%82%A4%E3%83%8A%E3%83%AA%E3%83%95%E3%82%A1%E3%82%A4%E3%83%AB%E3%82%B5%E3%83%BC%E3%83%90/)
    - [The easiest way to embed static files into a binary file in your Golang app (no external dependencies) - DEV Community](https://dev.to/koddr/the-easiest-way-to-embed-static-files-into-a-binary-file-in-your-golang-app-no-external-dependencies-43pc)
- gosec
- go generate

- gin static
    - [【Go】gin packageを使用して簡単なweb アプリケーションを作成する - Colorful Bullet](https://blog.bltinc.co.jp/entry/2020/02/04/141721)
    - [Go で作る SPA 用バックエンド - l12a](https://lnly.hatenablog.com/entry/2020/02/26/225722)
        - [golang gin api with static spa - Google 検索](https://www.google.com/search?q=golang+gin+api+with+static+spa&newwindow=1&sxsrf=ALiCzsZqi8ugQpt6ZO5mgmcNykWKy3K0bg%3A1655253833178&ei=SSupYuXCCpucseMPhLC6sA8&ved=0ahUKEwili_WFna74AhUbTmwGHQSYDvYQ4dUDCA4&uact=5&oq=golang+gin+api+with+static+spa&gs_lcp=Cgdnd3Mtd2l6EAM6BQgAEIAEOgUIABDLAToECAAQHjoFCCEQoAE6BggAEB4QCDoECCEQFUoECEEYAEoECEYYAFAAWPA0YKRAaAJwAXgAgAGcAYgB7RGSAQQ0LjE3mAEAoAEBwAEB&sclient=gws-wiz)

- routeing go
    - [Golang Ginを使った、すっきりルーティングのサンプル - Qiita](https://qiita.com/pon_maeda/items/c1fa3cf54ab432e8d45b)
- gin zap logging color
    - hack zap
    - color
        - [Go のコマンドラインツールのロギングとして zap を検証](https://zenn.dev/shunsuke_suzuki/scraps/542af5bd59863b)
    - what is sugar
        - [Go のロギングライブラリ zap について](https://zenn.dev/mima/articles/069b223d9b221f)
- is downloaded local storage mark
- removable client
- request download progress?
- override user id
- show empty message if client v-list-item is empty
- download movie,mp3 icon click propagation
- oembed error
    - Access to XMLHttpRequest at 'https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=zkZARKFuzNQ&format=json' from origin 'http://localhost' has been blocked by CORS policy: Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource.
