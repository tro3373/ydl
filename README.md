# ydl

## prd
```
make STAGE=prd up
```

```
All request handle by proxied api

              /api  proxy    :3000
-- nginx ------------------> server/ydl
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
