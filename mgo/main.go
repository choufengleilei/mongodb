package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

//测试结构体
type Person struct {
	Id    bson.ObjectId `bson:"_id"`
	Name  string        `bson:"tname"` //bson:"name" 表示mongodb数据库中对应的字段名称
	Phone string        `bson:"tphone"`
}

//如果是本机，并且MongoDB是默认端口27017启动的话，下面几种方式都可以。
//session, err := mgo.Dial("")
//session, err := mgo.Dial("localhost")
//session, err := mgo.Dial("127.0.0.1")
//session, err := mgo.Dial("localhost:27017")
//session, err := mgo.Dial("127.0.0.1:27017")
//如果不在本机或端口不同，传入相应的地址即可。如：
//mongodb://myuser:mypass@localhost:端口,otherhost:端口/mydb
const URL = "mongodb://test:123456@localhost:27017,远程服务器ip:27017/test" //mongodb连接字符串

var (
	mgoSession *mgo.Session
	dataBase   = "test"
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)

		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}
//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

/**
 * 添加person对象
 */
func AddPerson(p Person) string {
	p.Id = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection("person", query)
	if err != nil {
		return "false"
	}
	return p.Id.Hex()
}

/**
 * 获取一条记录通过objectid
 */
func GetPersonById(id string) *Person {
	objid := bson.ObjectIdHex(id)
	person := new(Person)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&person)
	}
	witchCollection("person", query)
	return person
}

//获取所有的person数据
func PagePerson() []Person {
	var persons []Person
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&persons)
	}
	err := witchCollection("person", query)
	if err != nil {
		return persons
	}
	return persons
}

//更新person数据
func UpdatePerson(query bson.M, change bson.M) string {
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := witchCollection("person", exop)
	if err != nil {
		return "true"
	}
	return "false"
}

/**
 * 执行查询，此方法可拆分做为公共方法
 * [SearchPerson description]
 * @param {[type]} collectionName string [description]
 * @param {[type]} query          bson.M [description]
 * @param {[type]} sort           bson.M [description]
 * @param {[type]} fields         bson.M [description]
 * @param {[type]} skip           int    [description]
 * @param {[type]} limit          int)   (results      []interface{}, err error [description]
 */
func SearchPerson(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = witchCollection(collectionName, exop)
	return
}


type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Name     string        `bson:"name"`
	PassWord string        `bson:"pass_word"`
	Age      int           `bson:"age"`
}

func main() {

	db, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	db.SetMode(mgo.Monotonic, true)
	c := db.DB("howie").C("person")

	//插入
	/*c.Insert(&User{
		Id:       bson.NewObjectId(),
		Name:     "leilei",
		PassWord: "123132",
		Age: 2,
	}, &User{
		Id:       bson.NewObjectId(),
		Name:     "leilei",
		PassWord: "qwer",
		Age: 5,
	}, &User{
		Id:       bson.NewObjectId(),
		Name:     "leilei",
		PassWord: "6666",
		Age: 7,
	})*/
	var users []User
	c.Find(nil).All(&users) //查询全部数据
	log.Println(users)

	c.FindId(users[0].Id).All(&users) //通过ID查询
	log.Println(users)

	c.Find(bson.M{"name": "leilei"}).All(&users) //单条件查询(=)
	log.Println(users)

	c.Find(bson.M{"name": bson.M{"$ne": "leilei"}}).All(&users) //单条件查询(!=)
	log.Println(users)

	c.Find(bson.M{"age": bson.M{"$gt": 5}}).All(&users) //单条件查询(>)
	log.Println(users)

	c.Find(bson.M{"age": bson.M{"$gte": 5}}).All(&users) //单条件查询(>=)
	log.Println(users)

	c.Find(bson.M{"age": bson.M{"$lt": 5}}).All(&users) //单条件查询(<)
	log.Println(users)

	c.Find(bson.M{"age": bson.M{"$lte": 5}}).All(&users) //单条件查询(<=)
	log.Println(users)

	/*c.Find(bson.M{"name": bson.M{"$in": []string{"JK_WEI", "JK_HE"}}}).All(&users) //单条件查询(in)
	log.Println(users)
	c.Find(bson.M{"$or": []bson.M{bson.M{"name": "JK_WEI"}, bson.M{"age": 7}}}).All(&users) //多条件查询(or)
	log.Println(users)
	c.Update(bson.M{"_id": users[0].Id}, bson.M{"$set": bson.M{"name": "JK_HOWIE", "age": 61}}) //修改字段的值($set)
	c.FindId(users[0].Id).All(&users)
	log.Println(users)
	c.Find(bson.M{"name": "JK_CHENG", "age": 66}).All(&users) //多条件查询(and)
	log.Println(users)
	c.Update(bson.M{"_id": users[0].Id}, bson.M{"$inc": bson.M{"age": -6,}}) //字段增加值($inc)
	c.FindId(users[0].Id).All(&users)
	log.Println(users)*/

	//c.Update(bson.M{"_id": users[0].Id}, bson.M{"$push": bson.M{"interests": "PHP"}}) //从数组中增加一个元素($push)

	c.Update(bson.M{"_id": users[0].Id}, bson.M{"$pull": bson.M{"interests": "go"}}) //从数组中删除一个元素($pull)

	c.FindId(users[0].Id).All(&users)
	log.Println(users)

	c.Remove(bson.M{"name": "leilei"})//删除


}
