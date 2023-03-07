package management

import (
	"dss/common/log"
	"dss/common/utils"
	"dss/core/dao"
	"dss/core/global"
	"dss/core/models"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"math"
	"strings"
	"time"
)

var (
	TaskManager *_TaskManager
)

type _TaskManager struct {
}

func parseRule(r models.RuleInsert, id bson.ObjectId) {
	_ = dao.Repo(global.ScanTask).SetField(id, bson.M{"status": global.Running})
	portStart, portEnd, _ := utils.ParsePortRange(r.TargetPort)
	switch r.Type {
	case global.Single:
		//192.168.1.1
		for i := portStart; i <= portEnd; i++ {
			pushRedis(models.Scan{
				RuleId:   r.Id.Hex(),
				TaskId:   id.Hex(),
				Host:     r.TargetHost,
				Port:     fmt.Sprintf("%v", i),
				Location: r.Location,
			})
		}
	case global.Range:
		//192.168.1.10-30
		startIp, startIpEndSuffix, ipRangeEndSuffix, _ := utils.ParseIpRange(r.TargetHost)
		resp := strings.Split(startIp, ".")
		prefix := fmt.Sprintf("%v.%v.%v.", resp[0], resp[1], resp[2])
		for i := startIpEndSuffix; i <= ipRangeEndSuffix; i++ {
			for p := portStart; p <= portEnd; p++ {
				pushRedis(models.Scan{
					RuleId:   r.Id.Hex(),
					TaskId:   id.Hex(),
					Host:     fmt.Sprintf("%v%v", prefix, i),
					Port:     fmt.Sprintf("%v", p),
					Location: r.Location,
				})
			}
		}
	case global.Cidr:
		//192.168.1.0/20
		ipSlice := utils.GetIpListByCidr(r.TargetHost)
		for _, ip := range ipSlice {
			for i := portStart; i <= portEnd; i++ {
				pushRedis(models.Scan{
					RuleId:   r.Id.Hex(),
					TaskId:   id.Hex(),
					Host:     ip,
					Port:     fmt.Sprintf("%v", i),
					Location: r.Location,
				})
			}
		}
	}
}

func (*_TaskManager) Post(query models.QueryID) error {
	var (
		err error
		r   models.RuleInsert
	)
	if !bson.IsObjectIdHex(query.ID) {
		return fmt.Errorf("invalid ObjectIdHex")
	}
	if err = dao.Repo(global.ScanRule).SelectById(dao.BsonId(query.ID), &r); err != nil {
		return err
	}
	if !r.Status {
		log.WarnF("rule %s status disable,skip scan...", r.Name)
		return nil
	}
	task := models.TaskInsertFunc(models.Task{
		RuleId: query.ID,
		Name:   r.Name,
		Status: global.Waiting,
		Count:  r.Count,
	})
	if err = dao.Repo(global.ScanTask).Insert(task); err != nil {
		return err
	}
	go parseRule(r, task.Id)
	RunTimeTaskStatusCheck(task.Id)
	return nil
}

func pushRedis(scan models.Scan) {
	val, _ := utils.Marshal(scan)
	if err := dao.Redis.LPush(global.PortScanQueue, val); err != nil {
		log.Errorf("push msg to redis err:%v", err)
	}
}

func RunTimeTaskStatusCheck(ids ...bson.ObjectId) {
	log.Info("start runTimeTaskStatusCheck...")
	var (
		query                = bson.M{}
		taskSlice            []models.TaskInsert
		incompleteTaskStatus = []global.TaskStatus{global.Running, global.Waiting}
	)
	if len(ids) != 0 {
		query["_id"] = ids[0]
	}
	query["status"] = bson.M{"$in": incompleteTaskStatus}
	if err := dao.Repo(global.ScanTask).Select(query, &taskSlice); err != nil {
		log.Errorf("query scan task info err:%v", err)
		return
	}
	if len(taskSlice) == 0 {
		return
	}
	log.InfoF("check incomplete task count:%d,taskSlice:%v", len(taskSlice), taskSlice)
	for _, task := range taskSlice {
		go func(task models.TaskInsert) {
			ticker := time.NewTicker(time.Second * 10)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					result, err := dao.Redis.Get(task.Id.Hex())
					if err != nil {
						log.Errorf("query key:%s from redis err:%s", task.Id.Hex(), err)
						continue
					}
					taskFinishedCount := utils.StrToInt(result)
					_ = dao.Repo(global.ScanTask).SetField(task.Id, bson.M{"executed_time": utils.ExecutedTimeFormat(time.Now().Unix() - task.CreatedTime),
						"progress": fmt.Sprintf("%v%%", math.Round(float64(taskFinishedCount)/float64(task.Count)*100))})
					log.InfoF("task id:%s,total count :%d,completed :%d", task.Id.Hex(), task.Count, taskFinishedCount)
					if taskFinishedCount == task.Count {
						_ = dao.Repo(global.ScanTask).SetField(task.Id, bson.M{"status": global.Finished})
						return
					}
				case <-global.Ctx.Done():
					return
				}
			}
		}(task)
	}
}

func (*_TaskManager) Get(param models.TaskQuery) (interface{}, error) {
	var (
		result models.TaskQueryResult
		resp   []*models.TaskInsert
		query  = bson.M{}
	)
	if param.Status != 0 {
		query["status"] = param.Status
	}
	if param.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": param.Search, "$options": "$i"}},
		}
	}
	if err := dao.Repo(global.ScanTask).SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
		return nil, err
	}
	result.Size = param.Size
	result.Page = param.Page
	result.Total = dao.Repo(global.ScanTask).Count(query)
	result.Items = models.TaskQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}
