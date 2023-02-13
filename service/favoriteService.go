package service

import (
	"SimpleDouyin/Utils"
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"

	"SimpleDouyin/config"
)

var Vedio_like = config.Vedio_like
var User_like = config.User_like
var limit_ip = config.LIMIT_IP
var time_out = config.LIMIT_PERIOD

/*
返回当前用户的所有点赞视频列表
*/
func GetVedioLikeList(userId string) []impl.Result {
	var IdList []int64
	ctx := context.Background()
	result, err := Utils.RDB.SMembers(ctx, User_like+userId).Result()
	if err != nil {
		panic(err)
	}
	for i, s := range result {
		IdList[i], err = strconv.ParseInt(s, 10, 64)
	}
	//调用service方法
	//vedioList := []dao.Video
	results := impl.QueryListByVedionl(IdList)
	//UuserServiceImpl := UserServiceImpl{}
	//for _, r := range results {
	//	vedio := dao.Video{
	//		Id: int64(r.ID),
	//		Author:
	//	}
	//
	//}
	return results
}

//
func LimitIP(ip, vedioId string) bool {
	ctx := context.Background()
	key := limit_ip + ip + vedioId
	result, err := Utils.RDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			//如果是第一次访问，则set并设置过期时间
			Utils.RDB.Set(ctx, key, 0, time_out).Result()
			return true
		} else {
			panic(err)
		}
	}
	count, _ := strconv.ParseInt(result, 10, 64)
	if count <= 10 {
		Utils.RDB.Incr(ctx, key)
	} else {
		//过期
		return false
	}
	return true
}

/*
返回当前视频的点赞总数，视频不存在则返回0
*/
func GetVedioLikeCount(vedioId string) int64 {
	var count = int64(0)
	ctx := context.Background()
	result, err := Utils.RDB.SCard(ctx, Vedio_like+vedioId).Result()
	if err != nil {
		panic(err)
	}
	count = result
	return count
}

/**
点赞操作
返回值为int，点赞成功则为1，否则为0且报错
*/
func Like(vedioId, userId string) int {
	ctx := context.Background()
	result, err := Utils.RDB.SAdd(ctx, Vedio_like+vedioId, userId).Result()
	if err != nil {
		panic(err)
	}
	result2, err := Utils.RDB.SAdd(ctx, User_like+userId, vedioId).Result()
	if err != nil {
		panic(err)
	}
	return int(result) & int(result2)
}

/*
测试用，为redis添加数据
*/
func Add(vedioId, userId string) int64 {
	ctx := context.Background()
	result, err := Utils.RDB.SAdd(ctx, Vedio_like+vedioId, userId).Result()
	if err != nil {
		panic(err)
	}
	return result
}

/*
测试用，为redis添加数据
*/
func AdduserId(userId, vedioId string) int64 {
	ctx := context.Background()
	result, err := Utils.RDB.SAdd(ctx, User_like+userId, vedioId).Result()
	if err != nil {
		panic(err)
	}
	return result
}

/**
取消点赞操作
返回值为int，点赞成功则为1，否则为0且报错
*/
func DislikeVedio(vedioId, userId string) int {
	ctx := context.Background()
	result, err := Utils.RDB.SRem(ctx, Vedio_like+vedioId, userId).Result()
	if err != nil {
		panic(err)
	}
	result2, err := Utils.RDB.SRem(ctx, User_like+userId, vedioId).Result()
	if err != nil {
		panic(err)
	}
	return int(result) & int(result2)
}

//查询当前用户是否点赞
func LikeVedioOrNot(vedioId, userId string) bool {
	ctx := context.Background()
	result, err := Utils.RDB.SIsMember(ctx, Vedio_like+vedioId, userId).Result()
	if err != nil {
		panic(err)
	}
	return result
}
