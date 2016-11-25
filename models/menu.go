package models

import (
	"beegostudy/util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

/**
*	pk		主键
*	auto 		自增值（限数值）
*	column(N)	指定字段名N
*	null		非空
*	index 		单个字段索引
* 	unique 		唯一键
* 	auto_now_add 	第一次插入数据时自动添加当前时间
* 	auto_now 	每一次保存时自动更新当前时间
* 	type(T)		对应数据库的指定类型
*	size(S)		类型长度S
*	default(D)	默认值D（需要对应类型）
**/
type Menu struct {
	Id        int       `orm:"pk;auto;column(id)"`
	Pid       int       `orm:"column(pid)"`
	MenuName  string    `orm:"column(menuname);size(64)"`
	Icon      string    `orm:"null;column(icon);size(32)"`
	IsLeaf    bool      `orm:"-"`
	Link      string    `orm:"null;column(link);size(128)"`
	InnerCode string    `orm:"column(innercode);size(128)"`
	Level     int       `orm:"column(level)"`
	OrderFlag int       `orm:"column(orderflag)"`
	AddTime   time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser   string    `orm:"column(adduser)"`
	Checked   bool      `orm:"-"`
	Expanded  bool      `orm:"-"`
	ChildNode []*Menu   `orm:"-"`
}

//自定义表名
func (m *Menu) TableName() string {
	return "menu"
}

func (menu *Menu) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	menuid, err := strconv.Atoi(tmpId)
	if err == nil {
		menu.Id = menuid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (menu *Menu) GetId() int {
	return menu.Id
}

func (menu *Menu) GetPid() int {
	return menu.Pid
}

func (menu *Menu) SetLevel(l int) {
	menu.Level = l
}

//是否为叶子节点
func (menu *Menu) GetIsLeaf() bool {
	o := orm.NewOrm()
	var count int
	o.Raw("SELECT COUNT(*) FROM menu WHERE pid=?", menu.Id).QueryRow(&count)
	return count == 0
}

func (menu *Menu) SetAddTime(t time.Time) {
	menu.AddTime = t
}

func (menu *Menu) SetCurrentTime() {
	menu.AddTime = time.Now()
}

func (menu *Menu) SetAddUser(uname string) {
	menu.AddUser = uname
}

func (menu *Menu) SetValue(data map[string]interface{}) error {
	return util.FillStruct(data, menu)
}

//查询数据库
func (menu *Menu) Fill() error {
	o := orm.NewOrm()
	if menu.Id > 0 {
		return o.Read(menu, "Id")
	}

	return fmt.Errorf("请确认是否传递了Id", "")

}

//插入
func (menu *Menu) Insert() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(menu)
}

//修改
func (menu *Menu) Update(column ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(menu, column...)
}

func init() {
	orm.RegisterModel(new(Menu))
}

func (menu *Menu) String() string {
	data, err := json.MarshalIndent(menu, "", "    ")
	if err != nil {
		beego.Warn("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

//相关函数

//根据ID得到菜单数据
func GetMenu(id int) (*Menu, error) {
	menu := Menu{Id: id}
	err := menu.Fill()
	if err != nil {
		beego.Error("菜单Id不存在")
		return &menu, err
	}
	return &menu, nil
}

//得到所有的菜单
func GetMenus() ([]*Menu, error) {
	var menus []*Menu
	o := orm.NewOrm()
	_, err := o.QueryTable("menu").OrderBy("innercode").All(&menus)
	return menus, err
}

//得到所有菜单并按级别排序
func GetMenusLevel(url string) ([]*Menu, error) {
	menus, err := GetMenus()
	if err == nil && menus != nil {
		var menuslevel []*Menu = make([]*Menu, len(menus), cap(menus))
		//用来记录menuslevel的当前位置
		idx := 0
		//线型遍历一遍
		for _, menu := range menus {
			menu.IsLeaf = true //默认全部是叶子节点
			if strings.EqualFold(menu.Link, url) {
				menu.Checked = true
			}
			if menuslevel[idx] == nil {
				menuslevel[idx] = menu
				continue
			}
			if menu.Pid == 0 {
				idx++
				menuslevel[idx] = menu
				continue
			}
			if menu.Pid == menuslevel[idx].Id {
				//如果当前元素是menuslevel第idx个元素的子集时就放入ChildNode中
				if menu.Checked {
					menuslevel[idx].Checked = true
					menuslevel[idx].IsLeaf = false //如果有子级则不为叶子节点
				}
				menuslevel[idx].ChildNode = append(menuslevel[idx].ChildNode, menu)
			}
			//循环到了这里的元素都是找不到父级的元素这里直接丢弃
		}
		return menuslevel[:idx+1], nil
	}
	return menus, err
}

//得到分页的菜单
/**
*	size	每页查询长度
*	index	查询的页码
*	ordercolumn	排序字段
*	orderby		升降序:desc\asc
**/
func GetMenusPage(size, index int, ordercolumn, orderby string) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "innercode"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var menus []*Menu
	o := orm.NewOrm()

	_, err := o.QueryTable("menu").OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&menus)

	if err == nil {
		cnt, err := o.QueryTable("menu").Count()

		pagetotal := cnt / int64(size)

		if cnt%int64(size) > 0 {
			pagetotal++
		}

		for _, menu := range menus {
			menu.Expanded = true
		}

		return GetDataGrid(menus, index, int(pagetotal), cnt), err
	}

	return nil, err
}

//返回对应级别的菜单组
func GetMenusByLevel(level int) ([]*Menu, error) {
	var menus []*Menu
	o := orm.NewOrm()
	_, err := o.QueryTable("menu").Filter("Level", level).OrderBy("innercode").All(&menus)
	return menus, err
}

type MenuSelectInit struct {
	Menu
	Select string
}

//返回所有顶级菜单，并指定当前父级
func GetTopMenus(pid int) ([]MenuSelectInit, error) {
	menus, err := GetMenusByLevel(1)
	if err != nil {
		return nil, err
	}

	menuselectinit := make([]MenuSelectInit, len(menus))
	for i, menu := range menus {
		menuselectinit[i].Menu.Id = menu.Id
		menuselectinit[i].Menu.MenuName = menu.MenuName
		if menu.Id == pid {
			menuselectinit[i].Select = "selected"
		}
	}
	return menuselectinit, nil
}
