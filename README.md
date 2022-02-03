### functions for now
<li>delete class:
curl -X DELETE http://localhost:9101/api/v1/class/7?key=2c78afaf-97da-4816-bbee-9ad239abb296</li>
<li>create class:
curl -X POST -H "Content-Type:application/json" -d "{"moduleid":"DF","classdate":"Monday","classstart":"11:00:00","classend":"13:00:00","classcap":30,"tutorname":"James_Lee"}" "http://localhost:9101/api/v1/class?key=2c78afaf-97da-4816-bbee-9ad239abb296"</li>
<li>update class(with class ID=1):
curl -H "Content-Type:application/json" -X PUT http://localhost:9101/api/v1/class/1?key=2c78afaf-97da-4816-bbee-9ad239abb296 -d "{\"ModuleID\":\"DF\"}"</li>

### final struct format would be:
    ModuleID   string `json:"moduleid"`
    ClassID    int    `json:"classid" gorm:"primaryKey"`
    ClassDate  string json:"classdate"
    ClassStart string json:"classstart"
    ClassEnd   string json:"classend"
    ClassCap   int    json:"classcap"
    TutorName string json:"tutorname"
    TutorID    int    json:tutorid
    Rating float32 json:rating
    ClassInfo string json:classinfo`

