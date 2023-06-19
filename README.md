# iptv-go
使用 [vercel](https://vercel.com/) 部署 [https://github.com/youshandefeiyang/LiveRedirect](https://github.com/youshandefeiyang/LiveRedirect) 的 [Golang脚本](https://github.com/youshandefeiyang/LiveRedirect/tree/main/Golang/liveurls)

## 部署步骤
1. Fork项目到自己的仓库
2. 在Vercel创建Project并选择`iptv-go`
3. Build&Deploy
4. Enjoy~

> 国内优化指南：
>
> 1. 自定义域名`CNAME`到[cname-china.vercel-dns.com](cname-china.vercel-dns.com)加速访问
>
> 2. Vercel的Function Region设置为香港服务器可以获得更低延迟
> ![Vercel设置](.github/asserts/region.png)
> 设置完需要重新部署生效

## 访问路径

`https://<你的域名>/live/平台/id`

> 注意路径多了一层`live`

## 一起看访问路径

`https://<你的域名>/yqk/xxx.m3u`

> 注意路径多了一层`yqk`

> 可选设置
> 1. 设置Vercel的环境变量`LIVE_URL`可以替换返回的m3u直播链接，就不用通过请求连接的`url`参数设置了
> 例如设置环境变量为: `https://www.baidu.com/live`
> ![环境变量设置](.github/asserts/env.png)
> 返回文件的效果
> ![环境变量设置](.github/asserts/envDemo.png)
> 优先级：链接的`url`参数 > 环境变量`LIVE_URL` > 默认域名

详细使用说明参考: [https://github.com/youshandefeiyang/LiveRedirect/blob/main/Golang/README.md](https://github.com/youshandefeiyang/LiveRedirect/blob/main/Golang/README.md)
