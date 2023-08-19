
# Hi~ this is *gophers* :wave:  
# Let's see our ***Simple—douyin***~.</center>


<!-- Introduction -->
## 目录
- [背景介绍](#背景介绍)
- [功能简介](#功能简介)
- [安装方法](#安装方法)
- [使用说明](#使用说明)
- [贡献](#贡献)
- [许可证](#许可证)

## 背景介绍
### 1.项目介绍：第六届字节跳动青训营后端任务————实现一个简易版的抖音。要求使用Go语言编程、常用框架、数据库、对象存储等内容。我们组使用Golang实现基于Redis的安全高效RPC通信。
### 2.项目要求：
<table class="MsoNormalTable" border="0" cellspacing="0" cellpadding="0" width="500" style="width:375.0pt;border-collapse:collapse;mso-yfti-tbllook:1184">
 <tbody><tr style="mso-yfti-irow:0;mso-yfti-firstrow:yes;height:29.25pt">
  <td style="border:solid #DEE0E3 1.0pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt"></td>
  <td colspan="2" style="border:solid #DEE0E3 1.0pt;border-left:none;mso-border-left-alt:
  solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;padding:.75pt .75pt .75pt .75pt;
  height:29.25pt">
  <p class="MsoNormal" align="center" style="text-align:center;mso-pagination:widow-orphan"><b><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">互动方向</span></b><span lang="EN-US" style="font-size:10.0pt;font-family:
  宋体;mso-bidi-font-family:宋体;mso-font-kerning:0pt"><o:p></o:p></span></p>
  </td>
  <td colspan="2" style="border:solid #DEE0E3 1.0pt;border-left:none;mso-border-left-alt:
  solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;padding:.75pt .75pt .75pt .75pt;
  height:29.25pt">
  <p class="MsoNormal" align="center" style="text-align:center;mso-pagination:widow-orphan"><b><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">社交方向</span></b><span lang="EN-US" style="font-size:10.0pt;font-family:
  宋体;mso-bidi-font-family:宋体;mso-font-kerning:0pt"><o:p></o:p></span></p>
  </td>
 </tr>
 <tr style="mso-yfti-irow:1;height:29.25pt">
  <td style="border:solid #DEE0E3 1.0pt;border-top:none;mso-border-top-alt:
  solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;padding:.75pt .75pt .75pt .75pt;
  height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">基础功能项<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td colspan="4" style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="center" style="text-align:center;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">视频<span lang="EN-US"> Feed </span>流、视频投稿、个人主页<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
 </tr>
 <tr style="mso-yfti-irow:2;height:29.25pt">
  <td style="border:solid #DEE0E3 1.0pt;border-top:none;mso-border-top-alt:
  solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;padding:.75pt .75pt .75pt .75pt;
  height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">基础功能项说明<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td colspan="4" style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">视频<span lang="EN-US">Feed</span>流：支持所有用户<span class="GramE">刷抖音</span>，视频按投稿时间倒序推出<span lang="EN-US"><o:p></o:p></span></span></p>
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">视频投稿：支持登录用户自己拍视频投稿<span lang="EN-US"><o:p></o:p></span></span></p>
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">个人主页：支持查看用户基本信息和投稿列表，注册用户流程简化<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
 </tr>
 <tr style="mso-yfti-irow:3;height:29.25pt">
  <td style="border:solid #DEE0E3 1.0pt;border-top:none;mso-border-top-alt:
  solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;padding:.75pt .75pt .75pt .75pt;
  height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">方向功能项<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="center" style="text-align:center;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">喜欢列表<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="center" style="text-align:center;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">用户评论<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="center" style="text-align:center;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">关系列表<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="center" style="text-align:center;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">消息<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
 </tr>
 <tr style="mso-yfti-irow:4;mso-yfti-lastrow:yes;height:29.25pt">
  <td style="border:solid #DEE0E3 1.0pt;border-top:none;mso-border-top-alt:
  solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;padding:.75pt .75pt .75pt .75pt;
  height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">方向功能项说明<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">登录用户可以对视频点赞，在个人主页喜欢<span lang="EN-US">Tab</span>下能够<span class="GramE">查看点赞视频</span>列表<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">支持未登录用户查看视频下的评论列表，登录用户能够发表评论<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">登录用户可以关注其他用户，能够在个人主页查看本人的关注数和粉丝数，查看关注列表和粉丝列表<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
  <td style="border-top:none;border-left:none;border-bottom:solid #DEE0E3 1.0pt;
  border-right:solid #DEE0E3 1.0pt;mso-border-top-alt:solid #DEE0E3 .75pt;
  mso-border-left-alt:solid #DEE0E3 .75pt;mso-border-alt:solid #DEE0E3 .75pt;
  padding:.75pt .75pt .75pt .75pt;height:29.25pt">
  <p class="MsoNormal" align="left" style="text-align:left;mso-pagination:widow-orphan"><span style="font-size:10.0pt;font-family:宋体;mso-bidi-font-family:宋体;mso-font-kerning:
  0pt">登录用户在消息<span class="GramE">页展示已</span>关注的用户列表，点击用户头像进入聊天页后可以发送消息<span lang="EN-US"><o:p></o:p></span></span></p>
  </td>
 </tr>
</tbody></table>

### 3.完成情况：

## 功能简介
### 👥 interaction

### 📳 message
### 🥳relation
### :selfie: user 
### 🎦vedio

## 安装方法
###

## 使用说明
###

## 贡献
### 贡献者
非常感谢这些人对该项目的付出~

## 许可证
