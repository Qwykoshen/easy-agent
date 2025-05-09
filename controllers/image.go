package controllers

import (
	"easy-agent/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ImageGenerationRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type ImageGenerationResponse struct {
	ImageURL string `json:"image_url"`
}

type BaiduImageResponse struct {
	Data struct {
		TaskProgress      int    `json:"task_progress"`
		TaskStatus        string `json:"task_status"`
		SubTaskResultList []struct {
			FinalImageList []struct {
				ImgURL string `json:"img_url"`
			} `json:"final_image_list"`
		} `json:"sub_task_result_list"`
	} `json:"data"`
	LogID int64 `json:"log_id"`
}

func GenerateImage(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.JSON(http.StatusOK, ImageGenerationResponse{
		ImageURL: "http://bj.bcebos.com/v1/ai-picture-creation/watermark/3_351567399_0_final.png?authorization=bce-auth-v1%2FALTAKBvI5HDpIAzJaklvFTUfAz%2F2025-05-08T08%3A24%3A17Z%2F2592000%2F%2F7b8ca2a84212a126ef18a48447cc5dad238887d0e97a684d1417721410782b02",
	},
	)
	return
	var req ImageGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("参数错误", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 第一步：调用百度API生成图片
	fmt.Println("access_token:", config.AccessToken)
	response, err := http.Post(
		fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/wenxin/v1/extreme/textToImage?access_token=%s", config.AccessToken),
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"prompt":"%s","width":1024,"height":1024}`, req.Prompt)),
	)

	if err != nil {
		fmt.Println("调用API失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "调用API失败"})
		return
	}
	defer response.Body.Close()

	// 解析响应获取task_id
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		fmt.Println("解析响应失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析响应失败"})
		return
	}

	taskID, ok := result["log_id"].(float64)
	if !ok {
		fmt.Println("获取任务ID失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务ID失败"})
		return
	}

	fmt.Println("task_id:", int64(taskID))

	// 第二步：查询图片生成结果
	queryURL := fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/wenxin/v1/extreme/getImg?access_token=%s", config.AccessToken)
	queryBody := fmt.Sprintf(`{"task_id":%d}`, int64(taskID))

	queryResp, err := http.Post(queryURL, "application/json", strings.NewReader(queryBody))
	if err != nil {
		fmt.Println("查询任务状态失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务状态失败"})
		return
	}
	defer queryResp.Body.Close()

	var baiduResp BaiduImageResponse
	if err := json.NewDecoder(queryResp.Body).Decode(&baiduResp); err != nil {
		fmt.Println("解析查询结果失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析查询结果失败"})
		return
	}
	fmt.Println("queryResp:", baiduResp)
	// 检查是否成功生成图片
	if len(baiduResp.Data.SubTaskResultList) > 0 &&
		len(baiduResp.Data.SubTaskResultList[0].FinalImageList) > 0 {
		imageURL := baiduResp.Data.SubTaskResultList[0].FinalImageList[0].ImgURL
		c.JSON(http.StatusOK, ImageGenerationResponse{
			ImageURL: imageURL,
		})
		return
	}
	fmt.Println("图片生成失败", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "图片生成失败"})
}
