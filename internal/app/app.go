package app

import (
	"context"
	"fmt"
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/service"
	"github.com/HulkLiu/WtTools/internal/utils"
	"github.com/evercyan/brick/xfile"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"log"
	"os/exec"
)

type App struct {
	Ctx            context.Context
	Log            *logrus.Logger
	CfgFile        string
	LogFile        string
	DB             DbClient
	SetManage      SettingManage
	TaskManage     service.TaskManager
	ResourceManage service.ResourceManage
	CourseManage   service.CourseManage
}
type DbClient struct {
	EsClient    *elastic.Client
	MysqlClient *gorm.DB
	RedisClient *redis.Client
}

func NewApp() *App {
	return &App{}
}

func (a *App) initDB() {
	a.DB.EsClient, err = elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Printf("initEsData connect failed,err:%v", err)
		return
	}
	a.DB.MysqlClient, err = gorm.Open(mysql.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	a.DB.RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       7,  // use default DB
	})

}

// OnStartup 初始化
func (a *App) OnStartup(ctx context.Context) {
	a.Ctx = ctx
	a.initDB()

	cfgPath := utils.GetCfgPath()

	//课程管理
	a.CourseManage = service.CourseManage{
		ElasticIndex: config.ElasticIndex,
		Client:       a.DB.EsClient,
		Err:          nil,
		PageSize:     20,
	}

	//设置
	a.SetManage = NewSet()

	//日志
	a.LogFile = fmt.Sprintf(config.LogFile, cfgPath)
	a.Log = utils.NewLogger(a.LogFile)
	a.Log.Info("OnStartup begin")

	//任务管理初始化
	a.TaskManage = service.NewTask(a.DB.MysqlClient)

	//资源管理
	a.ResourceManage = service.NewResource(a.DB.MysqlClient)
}

// diag ...
func (a *App) diag(message string, buttons ...string) (string, error) {
	if len(buttons) == 0 {
		buttons = []string{
			config.BtnConfirmText,
		}
	}
	return runtime.MessageDialog(a.Ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         config.Title,
		Message:       message,
		CancelButton:  config.BtnConfirmText,
		DefaultButton: config.BtnConfirmText,
		Buttons:       buttons,
	})
}

// Menu 应用菜单
func (a *App) Menu() *menu.Menu {
	return menu.NewMenuFromItems(
		menu.SubMenu("文件", menu.NewMenuFromItems(
			menu.Text("关于", nil, func(_ *menu.CallbackData) {
				a.diag(config.Description)
			}),
			menu.Text("检查更新", nil, func(_ *menu.CallbackData) {
				a.diag(config.VersionNewMsg)
			}),
			menu.Separator(),

			menu.Separator(),
			menu.Text("退出", keys.CmdOrCtrl("Q"), func(_ *menu.CallbackData) {
				runtime.Quit(a.Ctx)
			}),
		)),
		menu.EditMenu(),
		menu.SubMenu("帮助", menu.NewMenuFromItems(
			menu.Text(
				"打开配置文件",
				keys.Combo("C", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					if !xfile.IsExist(a.CfgFile) {
						a.diag("文件不存在, 请先更新配置")
						return
					}
					_, err := exec.Command("open", a.CfgFile).Output()
					if err != nil {
						a.diag("操作失败: " + err.Error())
						return
					}
				},
			),
			menu.Text(
				"打开日志文件",
				keys.Combo("L", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					if !xfile.IsExist(a.LogFile) {
						a.diag("文件不存在, 请先更新配置")
						return
					}
					_, err := exec.Command("open", a.LogFile).Output()
					if err != nil {
						a.diag("操作失败: " + err.Error())
						return
					}
				},
			),
			menu.Separator(),
			menu.Text(
				"打开应用主页",
				keys.Combo("H", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					runtime.BrowserOpenURL(a.Ctx, config.GitRepoURL)
				},
			),
		)),
	)
}

// OnDomReady ...
func (a *App) OnDomReady(ctx context.Context) {
	a.Log.Info("OnDomReady")
	return
}

// OnShutdown ...
func (a *App) OnShutdown(ctx context.Context) {
	a.Log.Info("OnShutdown")
	return
}

// OnBeforeClose ...
func (a *App) OnBeforeClose(ctx context.Context) bool {
	a.Log.Info("OnBeforeClose")
	// 返回 true 将阻止程序关闭
	return false
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
