package test

import (
	"core/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
)

func TestGorm(t *testing.T) {
	db, err := gorm.Open("mysql", "root:501124524@(116.62.177.68)/cloud-disk?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	var user models.UserBasic
	db.Where("id = ?", 1).Find(&user)
	if user.Id == 1 {
		fmt.Println("存在数据")
	}
	fmt.Println("不存在数据")

	//count := 0
	//l..Where("email = ?", "123456@qq.com").Table("user_basic").Count(&count)
	//fmt.Println("count: ", count)

	//b, err := json.Marshal(data)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//dst := new(bytes.Buffer)
	//err = json.Indent(dst, b, "", " ")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(dst.String())
	//user := models.UserBasic{}
	//db.First(&user)

	//fmt.Println("user: ", user)
}
