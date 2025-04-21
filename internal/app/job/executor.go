package job

import (
	"com.say.more.server/utils/helper"
	"github.com/dcron-contrib/commons/dlog"
	"github.com/dcron-contrib/redisdriver"
	"github.com/libi/dcron"
	"github.com/libi/dcron/cron"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type Job interface {
	Name() string
	Run()
	Cron() string
}

type Executor struct {
	cr  *dcron.Dcron
	log *dlog.StdLogger
}

func NewExecutor(addr, password string) *Executor {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	driver := redisdriver.NewDriver(redisCli)
	l := &dlog.StdLogger{
		Log:        log.New(os.Stdout, "[job] ", log.LstdFlags),
		LogVerbose: true,
	}

	executor := &Executor{
		log: l,
	}
	executor.cr = dcron.NewDcron("job", driver, cron.WithLogger(l), cron.WithSeconds(), cron.WithLocation(time.UTC))
	return executor
}

// Start start distribution job
func (e Executor) Start() {
	e.addJob(NewStartCourseJob())
	e.addJob(NewScanCourseJob())
	go e.cr.Start()
}

// addJob add job to dcron
func (e Executor) addJob(job Job) {
	recoverCmd := func() {
		helper.RunSafe(func() {
			begin := time.Now()
			job.Run()
			e.log.Infof("execute sync job cost: %s", time.Since(begin).String())
		}, func(err error) {
			e.log.Errorf("execute sync job failed: %s", err.Error())
		})
	}
	//skips an invocation of the Job if a previous invocation is still running
	jobWrapper := cron.SkipIfStillRunning(e.log)
	err := e.cr.AddJob(job.Name(), job.Cron(), jobWrapper(cron.FuncJob(recoverCmd)))
	if err != nil {
		e.log.Errorf("add job failed: %s", err.Error())
	}
}

func (e Executor) Info(msg string, keysAndValues ...interface{}) {
	e.log.Infof(msg, keysAndValues)
}
func (e Executor) Error(err error, msg string, keysAndValues ...interface{}) {
	e.log.Errorf(err.Error(), msg, keysAndValues)
}
