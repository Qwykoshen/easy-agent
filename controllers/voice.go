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

type VoiceGenerationRequest struct {
	Text string `json:"text" binding:"required"`
}

type VoiceGenerationResponse struct {
	AudioURL string `json:"audio_url"`
}

type BaiduTTSResponse struct {
	TaskId     string `json:"task_id"`
	TaskStatus string `json:"task_status"`
	Speech     string `json:"speech_url,omitempty"`
	Error      int    `json:"error_code,omitempty"`
}

type BadiduTTSQueryResponse struct {
	LogId int64 `json:"log_id"`
	Data  []struct {
		TaskStatus string `json:"task_status"`
		TaskResult struct {
			Speech string `json:"speech_url"`
		} `json:"task_result"`
	} `json:"tasks_info"`
}

func GenerateVoice(c *gin.Context) {
	var req VoiceGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 第一步：调用百度TTS API创建任务
	response, err := http.Post(
		fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/tts/v1/create?access_token=%s", config.AccessToken),
		"application/json",
		strings.NewReader(fmt.Sprintf(`{
			"text": ["%s"],
			"format": "mp3-16k",
			"voice": 5003,
			"lang": "zh",
			"speed": 5,
			"enable_subtitle": 2,
			"break":0
		}`, req.Text)),
	)

	if err != nil {
		fmt.Println("调用API失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "调用API失败"})
		return
	}
	defer response.Body.Close()

	var ttsResp BaiduTTSResponse
	if err := json.NewDecoder(response.Body).Decode(&ttsResp); err != nil {
		fmt.Println("解析响应失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析响应失败"})
		return
	}

	if ttsResp.Error != 0 {
		fmt.Println("语音合成请求失败", ttsResp.Error)
		fmt.Println("响应内容", ttsResp) // 打印响应内容以获取更多信息，如错误码和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": "语音合成请求失败"})
		return
	}

	// 第二步：轮询查询任务状态
	taskID := ttsResp.TaskId
	fmt.Println("taskID: ", taskID)
	maxRetries := 20
	for i := 0; i < maxRetries; i++ {
		time.Sleep(1 * time.Second)

		queryResp, err := http.Post(
			fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/tts/v1/query?access_token=%s", config.AccessToken),
			"application/json",
			strings.NewReader(fmt.Sprintf(`{"task_ids":["%s"]}`, taskID)),
		)

		if err != nil {
			fmt.Println("查询任务状态失败", err)
			continue
		}

		var queryResult BadiduTTSQueryResponse
		if err := json.NewDecoder(queryResp.Body).Decode(&queryResult); err != nil {
			fmt.Println("解析查询结果失败", err)
			queryResp.Body.Close()
			continue
		}
		queryResp.Body.Close()
		if queryResult.Data[0].TaskStatus == "Success" && queryResult.Data[0].TaskResult.Speech != "" {
			c.JSON(http.StatusOK, VoiceGenerationResponse{
				AudioURL: queryResult.Data[0].TaskResult.Speech,
			})
			return
		}
		fmt.Println("未查询到任务完成...")

		if queryResult.Data[0].TaskStatus == "Failed" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "语音合成失败"})
			return
		}
	}
	fmt.Println("语音合成超时, 退出")
	c.JSON(http.StatusInternalServerError, gin.H{"error": "语音合成超时"})
}
