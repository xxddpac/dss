package management

import (
	"bytes"
	"context"
	"dss/common/http"
	"dss/common/log"
	"dss/common/utils"
	"dss/core/config"
	"dss/core/dao"
	"dss/core/global"
	"dss/core/models"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	PortManager *_PortManager
)

type _PortManager struct {
}

type ErrInfo struct {
	// `fix typo in word with 2 json tags`
	//  Goland -> Settings -> Editor -> Natural Languages -> Spelling -> Accepted words
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (*_PortManager) Get(param models.ScanQuery) (interface{}, error) {
	var (
		result models.ScanQueryResult
		resp   []*models.ScanInsert
		query  = bson.M{}
	)
	if param.Location != "" {
		query["location"] = param.Location
	}
	if param.Date != "" {
		query["done_time"] = param.Date
	} else {
		query["done_time"] = time.Now().Format(utils.TimeLayout)
	}
	if param.Search != "" {
		query["$or"] = []bson.M{
			{"host": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"port": bson.M{"$regex": param.Search, "$options": "$i"}},
		}
	}
	if err := dao.Repo(global.PortScan).SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
		return nil, err
	}
	if param.Date == "" {
		if len(resp) == 0 {
			query["done_time"] = time.Now().AddDate(0, 0, -1).Format(utils.TimeLayout)
			if err := dao.Repo(global.PortScan).SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
				return nil, err
			}
		}
	}
	result.Size = param.Size
	result.Page = param.Page
	result.Total = dao.Repo(global.PortScan).Count(query)
	result.Items = models.ScanQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}

func (*_PortManager) FieldGroupBy(field string) (pipeline []bson.M) {
	group := bson.M{"$group": bson.M{"_id": fmt.Sprintf("$%s", field), "count": bson.M{"$sum": 1}}}
	orderBy := bson.M{"$sort": bson.M{"_id": 1}}
	pipeline = []bson.M{group, orderBy}
	return
}

func (*_PortManager) Parse(resp []bson.M) (result []string) {
	for _, item := range resp {
		if _, ok := item["_id"].(string); !ok {
			continue
		}
		result = append(result, item["_id"].(string))
	}
	return
}

func (*_PortManager) Location() (interface{}, error) {
	var (
		err            error
		field          = "location"
		resp, pipeline []bson.M
	)
	pipeline = PortManager.FieldGroupBy(field)
	if err = dao.Repo(global.PortScan).Aggregate(pipeline, &resp); err != nil {
		return nil, err
	}
	return PortManager.Parse(resp), nil
}

func (*_PortManager) Clear() {
	var (
		err            error
		field          = "done_time"
		resp, pipeline []bson.M
	)
	pipeline = PortManager.FieldGroupBy(field)
	if err = dao.Repo(global.PortScan).Aggregate(pipeline, &resp); err != nil {
		log.Errorf("group by field err:%v", err)
		return
	}
	result := PortManager.Parse(resp)
	if len(result) <= 7 {
		return
	}
	result = result[:len(result)-7]
	for _, item := range result {
		if err = dao.Repo(global.PortScan).RemoveAll(bson.M{field: item}); err != nil {
			log.Errorf("remove field err:%v", err)
		}
	}
}

func (*_PortManager) Trend() (interface{}, error) {
	var (
		err            error
		field          = "done_time"
		resp, pipeline []bson.M
		result         []map[string]interface{}
	)
	pipeline = PortManager.FieldGroupBy(field)
	if err = dao.Repo(global.PortScan).Aggregate(pipeline, &resp); err != nil {
		return nil, err
	}
	res := PortManager.Parse(resp)
	if len(res) > 7 {
		res = res[len(res)-7:]
	}
	for _, item := range res {
		tmp := make(map[string]interface{})
		tmp["date"] = item
		tmp["count"] = dao.Repo(global.PortScan).Count(bson.M{field: item})
		result = append(result, tmp)
	}
	return result, nil
}

func (*_PortManager) Remind() {
	var (
		s                                    []models.Scan
		err                                  error
		todayScanResult, yesterdayScanResult []string
		check                                func(todayScanResult, yesterdayScanResult []string)
		query                                = bson.M{}
		todayTimeLayout                      = time.Now().Format(utils.TimeLayout)
		yesterdayTimeLayout                  = time.Now().AddDate(0, 0, -1).Format(utils.TimeLayout)
	)
	query["done_time"] = bson.M{"$in": []string{todayTimeLayout, yesterdayTimeLayout}}
	if err = dao.Repo(global.PortScan).Select(query, &s); err != nil {
		log.Errorf("select data err:%v", err)
		return
	}
	for _, item := range s {
		hostPort := fmt.Sprintf("%s-%s", item.Host, item.Port)
		if todayTimeLayout == item.DoneTime {
			todayScanResult = append(todayScanResult, hostPort)
		} else {
			yesterdayScanResult = append(yesterdayScanResult, hostPort)
		}
	}
	if len(todayScanResult) == 0 || len(yesterdayScanResult) == 0 {
		log.InfoF("no data found in today or yesterday")
		return
	}
	defer func() {
		check(todayScanResult, yesterdayScanResult)
	}()
	check = func(todayScanResult, yesterdayScanResult []string) {
		var (
			key              string
			newOpenItemSlice []string
			preInsertXlsx    [][]string
			header           = []string{"host", "port", "date"}
			headers          map[string]string
			body             io.Reader
			xlsxName         = fmt.Sprintf("new_port_open_%s.xlsx", time.Now().Format(utils.TimeLayout))
			filePath         = filepath.Join(os.TempDir(), xlsxName)
		)
		for _, item := range todayScanResult {
			if !utils.IsStrExists(yesterdayScanResult, item) {
				newOpenItemSlice = append(newOpenItemSlice, item)
			}
		}
		if len(newOpenItemSlice) == 0 {
			return
		}
		for _, item := range newOpenItemSlice {
			data := strings.Split(item, "-")
			if len(data) != 2 {
				continue
			}
			host := data[0]
			port := data[1]
			preInsertXlsx = append(preInsertXlsx, []string{host, port, time.Now().Format(utils.TimeLayout)})
		}
		if err = utils.WriteToXlsx(filePath, header, preInsertXlsx); err != nil {
			log.Errorf("write data to xlsx err:%v", err)
			return
		}
		defer func() {
			if err = os.Remove(filePath); err != nil {
				log.Errorf("remove tmp file %v err:%v", filePath, err)
			}
		}()
		headers, body, err = PortManager.buildMultipartFormData(filePath)
		if err != nil {
			log.Errorf("build multipart form_data err:%v", err)
			return
		}
		key, err = PortManager.uploadFileToWorkChat(headers, body)
		if err != nil {
			log.Errorf("uploadFileToWorkChat err:%v", err)
			return
		}
		if err = PortManager.notify(key); err != nil {
			log.Errorf("notify err:%v", err)
		}
	}
}

func (*_PortManager) notify(key string) error {
	var (
		e              ErrInfo
		byteStr        []byte
		err            error
		workChatBotUrl = config.CoreConf.Producer.WorkChatBotUrl
		request        = http.NewClient(workChatBotUrl)
	)
	body := fmt.Sprintf(`{"msgtype":"file","file": {"media_id": "%v"}}`, key)
	_, byteStr, err = request.Post(context.Background(), "", body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(byteStr, &e); err != nil {
		return err
	}
	if e.ErrCode != 0 {
		return fmt.Errorf(e.ErrMsg)
	}
	return nil
}

func (*_PortManager) uploadFileToWorkChat(headers map[string]string, body io.Reader) (string, error) {
	type uploadFile struct {
		MediaId string `json:"media_id"`
		ErrInfo
	}
	var (
		u                 uploadFile
		workChatUploadUrl = config.CoreConf.Producer.WorkChatUploadUrl
		request           = http.NewClient(workChatUploadUrl, headers)
	)
	_, byteStr, err := request.Post(context.Background(), "", body)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(byteStr, &u); err != nil {
		return "", err
	}
	if u.ErrCode != 0 {
		return "", fmt.Errorf("get media_id err,err_code:%d,err_msg:%v", u.ErrCode, u.ErrMsg)
	}
	return u.MediaId, nil
}

func (*_PortManager) buildMultipartFormData(path string) (map[string]string, io.Reader, error) {
	var (
		err       error
		file      *os.File
		part      io.Writer
		fieldName = "media"
	)
	file, err = os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err = writer.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return nil, nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, nil, err
	}
	return map[string]string{"Content-Type": writer.FormDataContentType()}, body, nil
}
