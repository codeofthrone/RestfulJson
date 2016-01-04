package dashboardlib

import (
	"database/sql"
	_ "github.com/GO-SQL-Driver/MySQL"
	"github.com/pmylund/sortutil"
	// "html/template"
	"log"
	"sort"
	"strconv"
	"strings"
	// "time"
)

type JiraIssueVersions struct {
	Sn   int
	Name string
}

type JiraIssuesStatusDiffVersion struct {
	Name          string `json:"Name"`
	CreateCounter int    `json:"CreateCounter"`
	CloseCounter  int    `json:"CloseCounter"`
	RemainCounter int    `json:"RemainCounter"`
}

type JiraIssuesDateRemain struct {
	Name          string `json:"Name"`
	CreateCounter int    `json:"CreateCounter"`
	CloseCounter  int    `json:"CloseCounter"`
	RemainCounter int    `json:"RemainCounter"`
}

type JiraIssuesPriority struct {
	Fatal  string `json:"Name"`
	Low    string `json:"Name"`
	Medium string `json:"Name"`
	Name   string `json:"Name"`
}

type JsonResult struct {
	Name          []string `json:"Name"`
	CreateCounter []int    `json:"CreateCounter"`
	CloseCounter  []int    `json:"CloseCounter"`
	RemainCounter []int    `json:"RemainCounter"`
}

type JsonResultStr struct {
	Name          []string `json:"Name"`
	CreateCounter []string `json:"CreateCounter"`
	CloseCounter  []string `json:"CloseCounter"`
	RemainCounter []string `json:"RemainCounter"`
}

type JsonResultTime struct {
	Week       string `json:"Week"`
	LastUpdate string `json:"LastUpdate"`
	Current    string `json:"Current"`
}

type JsonResultPie struct {
	Name          []string  `json:"Name"`
	RemainCounter []float64 `json:"RemainCounter"`
}

type JsonResultTable struct {
	Name     string `json:"Key"`
	Summary  string `json:"Summary"`
	Priority string `json:"Priority"`
	DiffDate int    `json:"date"`
	Assignee string `json:"assignee"`
}

type ProjectSummaryTable struct {
	Project       string  `json:"Project"`
	Version       string  `json:"Version"`
	FileName      string  `json:"FileName"`
	Scenario      float64 `json:"Scenario"`
	DataCheck     float64 `json:"DataCheck"`
	Auto          float64 `json:"Auto"`
	BDI           float64 `json:"BDI"`
	Compatibility float64 `json:"Compatibility"`
	Security      float64 `json:"Security"`
	Other         float64 `json:"Other"`
	Battery       float64 `json:"Battery"`
	Date          string  `json:"Date"`
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func checkerr(str string, err error) {
	if err != nil {
		log.Println(str)
		log.Fatal(err.Error())
		panic(err.Error())
	}
}

func QuerySingle(query string) string {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	InsertType, err := db.Query(query)
	checkerr(query, err)
	InsertType.Next()
	var tmp string
	InsertType.Scan(&tmp)
	InsertType.Close()
	db.Close()
	return tmp
}

// func ReleaseDayIssueTable(ProjectName string) RemainIssueResult {
// 	//TODO
// 	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
// 	RemainIssueSQL := "SELECT a.`Key`, a.`Summary`, a.`Priority`, a.DiffDate, a.`Assignee` from (SELECT *,DATEDIFF(DATE(`UpdatedTime`),DATE(CURDATE() )) AS DiffDate " +
// 		"FROM `Issues` WHERE `Key` LIKE '" + ProjectName + "-%' AND Status!='Closed' AND Resolution not like  '% Fix' ORDER BY UpdatedTime AND Status='Closed' )as a " +
// 		"ORDER by a.DiffDate"
// 	RemainIssues, err := db.Query(RemainIssueSQL)
// 	checkerr(RemainIssueSQL, err)
// 	var RemainIssueResult []JsonResultTable
// 	for RemainIssues.Next() {
// 		var tmpRemainIssue JsonResultTable
// 		RemainIssues.Scan(&tmpRemainIssue.Name, &tmpRemainIssue.Summary, &tmpRemainIssue.Priority, &tmpRemainIssue.DiffDate, &tmpRemainIssue.Assignee)
// 		// log.Println(tmpRemainIssue)
// 		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)

// 	}
// 	RemainIssues.Close()
// 	db.Close()
// 	return RemainIssueResult
// }

func WeekRemain(ProjectName string) JsonResult {
	var ReturnJson JsonResult
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	WeekCreateSQL := "SELECT  WEEK(`CreatedTime`)+1 as week , count(*)  FROM Issues  WHERE `Key` LIKE '" + ProjectName + "-%'  group by week"
	WeekRemains, err := db.Query(WeekCreateSQL)
	// log.Println(WeekCreateSQL)
	checkerr(WeekCreateSQL, err)
	var WeekReaminResult []JiraIssuesDateRemain
	WeekRemainCounter := 0
	for WeekRemains.Next() {
		var (
			WeekRemain       JiraIssuesDateRemain
			tmpCreateCounter int
			tmpName          string
		)
		WeekRemains.Scan(&tmpName, &tmpCreateCounter)
		ClosqSQL := "Select closecounter from ( SELECT count(*) as closecounter, WEEK(UpdatedTime)+1 as week FROM Issues " +
			" WHERE `Key` LIKE '" + ProjectName + "-%' AND Status='Closed'  group by week ) as a WHERE a.week='" + tmpName + "'"
		tmpCloseCounter := QuerySingle(ClosqSQL)

		WeekRemain.Name = tmpName
		WeekRemain.CreateCounter = tmpCreateCounter
		WeekRemain.CloseCounter, _ = strconv.Atoi(tmpCloseCounter)
		WeekRemainCounter = WeekRemainCounter + tmpCreateCounter - WeekRemain.CloseCounter
		WeekRemain.RemainCounter = WeekRemainCounter
		WeekReaminResult = append(WeekReaminResult, WeekRemain)
	}
	WeekRemains.Close()
	sortutil.AscByField(WeekReaminResult, "Name")
	for _, value := range WeekReaminResult {
		// log.Println(key, value)
		ReturnJson.Name = append(ReturnJson.Name, value.Name)
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, value.CreateCounter)
		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, value.CloseCounter)
		ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, value.RemainCounter)
	}
	db.Close()
	// log.Println(ReturnJson)
	return ReturnJson
}

func VersionData(ProjectName string) []JiraIssuesStatusDiffVersion {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	// ProjectName := "TYGH"
	versSQL := "SELECT Id,Name from ( " +
		"SELECT f.Id,f.Name,f.Project FROM Fixversions as f LEFT join Versions as v on f.Id = v.Id UNION DISTINCT " +
		"SELECT v.Id,v.Name,v.Project FROM Fixversions as f RIGHT join Versions as v on f.Id = v.Id) as a where  Project='" + ProjectName + "'  group by a.Id " +
		"ORDER BY `a`.`Id` ASC"
	Vers, err := db.Query(versSQL)
	checkerr(versSQL, err)
	var DiffVersionResult []JiraIssuesStatusDiffVersion
	for Vers.Next() {
		var (
			IssueVersions JiraIssuesStatusDiffVersion
			tmpVerId      string
			tmpVerName    string
		)
		Vers.Scan(&tmpVerId, &tmpVerName)
		IssueVersions.Name = tmpVerName
		DiffVersionResult = append(DiffVersionResult, IssueVersions)
	}
	Vers.Close()

	// // TODO no version
	sortutil.AscByField(DiffVersionResult, "Name")
	IssuesSQL := "SELECT Id,FixVersions,Status,Version FROM Issues WHERE `Key` like '" + ProjectName + "-%' AND Resolution not like '% Fix'"
	Issuess, err := db.Query(IssuesSQL)
	checkerr(IssuesSQL, err)
	for Issuess.Next() {
		var (
			tmpId             string
			tmpFixVer         string
			tmpVer            string
			tmpStatus         string
			CreateVersionName []string
		)
		Issuess.Scan(&tmpId, &tmpFixVer, &tmpStatus, &tmpVer)
		// log.Println(tmpId, tmpFixVer, tmpStatus, tmpVer)
		IssueCreateCounter := 0
		if tmpVer != "0" {
			VerDataSQL := " SELECT Data FROM `Version` WHERE `Id` = '" + tmpId + "' AND Enable ='1' "
			VerDatas, err := db.Query(VerDataSQL)
			checkerr(VerDataSQL, err)
			for VerDatas.Next() {
				var tmpData string
				VerDatas.Scan(&tmpData)
				VersNameSQL := " SELECT Name FROM `Versions` WHERE `Sn` = '" + tmpData + "' AND Project ='" + ProjectName + "' "
				VerName := QuerySingle(VersNameSQL)
				// log.Println("Create", VerName)
				CreateVersionName = append(CreateVersionName, VerName)
				for key, Value := range DiffVersionResult {
					if Value.Name == VerName {
						DiffVersionResult[key].CreateCounter++
						IssueCreateCounter++
					}
				}
			}
			VerDatas.Close()
		}
		// log.Println("Create counter", IssueCreateCounter)
		if tmpStatus == "Closed" {
			FixVerDataSQL := " SELECT Data FROM `Fixversion` WHERE `Id` = '" + tmpId + "' AND Enable ='1' "
			VerDatas, err := db.Query(FixVerDataSQL)
			checkerr(FixVerDataSQL, err)
			for VerDatas.Next() {
				var tmpData string
				VerDatas.Scan(&tmpData)
				VersNameSQL := " SELECT Name FROM `Fixversions` WHERE `Sn` = '" + tmpData + "' AND Project ='" + ProjectName + "' "
				VerName := QuerySingle(VersNameSQL)
				// log.Println("close", VerName)
				for key, Value := range DiffVersionResult {
					if Value.Name == VerName && IssueCreateCounter > 0 {
						DiffVersionResult[key].CloseCounter++
						IssueCreateCounter--
					}
				}
			}
		}
		// log.Println("Close counter", IssueCreateCounter, VersionRemainCounter)
	}
	Issuess.Close()
	db.Close()
	return DiffVersionResult
}

func TYGHDiffVersion(ProjectName string) (JsonResult, JsonResult) {
	DiffVersionResult := VersionData(ProjectName)
	var ReturnJsonAPP JsonResult
	var ReturnJsonWEB JsonResult
	APPRemain := 0
	WEBRemain := 0
	appflag := 0
	webflag := 0
	for _, Value := range DiffVersionResult {
		if strings.Contains(Value.Name, "APP") {
			if appflag == 0 {
				ReturnJsonAPP.RemainCounter = append(ReturnJsonAPP.RemainCounter, 0)
			} else {
				ReturnJsonAPP.RemainCounter = append(ReturnJsonAPP.RemainCounter, APPRemain)
			}
			APPRemain = APPRemain + Value.CreateCounter - Value.CloseCounter
			ReturnJsonAPP.Name = append(ReturnJsonAPP.Name, Value.Name)
			ReturnJsonAPP.CreateCounter = append(ReturnJsonAPP.CreateCounter, Value.CreateCounter)
			ReturnJsonAPP.CloseCounter = append(ReturnJsonAPP.CloseCounter, Value.CloseCounter)
			appflag++
		} else {
			if webflag == 0 {
				ReturnJsonWEB.RemainCounter = append(ReturnJsonWEB.RemainCounter, 0)
			} else {
				ReturnJsonWEB.RemainCounter = append(ReturnJsonWEB.RemainCounter, WEBRemain)
			}
			ReturnJsonWEB.Name = append(ReturnJsonWEB.Name, Value.Name)
			WEBRemain = WEBRemain + Value.CreateCounter - Value.CloseCounter
			ReturnJsonWEB.CreateCounter = append(ReturnJsonWEB.CreateCounter, Value.CreateCounter)
			ReturnJsonWEB.CloseCounter = append(ReturnJsonWEB.CloseCounter, Value.CloseCounter)
			webflag++
		}
	}
	return ReturnJsonAPP, ReturnJsonWEB
}

func DiffVersion(ProjectName string) JsonResult {
	DiffVersionResult := VersionData(ProjectName)
	var ReturnJson JsonResult
	VersionRemainCounter := 0
	for key, _ := range DiffVersionResult {
		if key == 0 {
			DiffVersionResult[key].RemainCounter = 0
			VersionRemainCounter = DiffVersionResult[key].CreateCounter - DiffVersionResult[key].CloseCounter
		} else {
			DiffVersionResult[key].RemainCounter = VersionRemainCounter
			VersionRemainCounter = VersionRemainCounter + DiffVersionResult[key].CreateCounter - DiffVersionResult[key].CloseCounter
		}
		ReturnJson.Name = append(ReturnJson.Name, DiffVersionResult[key].Name)
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, DiffVersionResult[key].CreateCounter)
		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, DiffVersionResult[key].CloseCounter)
		ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, DiffVersionResult[key].RemainCounter)
		// log.Println(key, Value)
	}
	// log.Println(ReturnJson)
	return ReturnJson
}

func DiffVersionSoFarRemain(ProjectName string) JsonResult {
	ReturnJson := DiffVersion(ProjectName)
	soFarClose := 0
	for Key, _ := range ReturnJson.Name {
		tmp := ReturnJson.CloseCounter[Key]
		if Key == 0 {
			ReturnJson.CloseCounter[Key] = 0
		} else {
			ReturnJson.CloseCounter[Key] = soFarClose
		}
		soFarClose = soFarClose + tmp
	}
	return ReturnJson
}

func TYGHDiffVersionSoFarRemain(ProjectName string) (JsonResult, JsonResult) {
	ReturnJsonAPP, ReturnJsonWEB := TYGHDiffVersion(ProjectName)
	soFarAPPClose := 0
	soFarWEBClose := 0
	for Key, _ := range ReturnJsonAPP.Name {
		tmp := ReturnJsonAPP.CloseCounter[Key]
		if Key == 0 {
			ReturnJsonAPP.CloseCounter[Key] = 0
		} else {
			ReturnJsonAPP.CloseCounter[Key] = soFarAPPClose
		}
		soFarAPPClose = soFarAPPClose + tmp

		// soFarAPPClose = soFarAPPClose + ReturnJsonAPP.CloseCounter[Key]
		// ReturnJsonAPP.CloseCounter[Key] = soFarAPPClose
	}

	for Key, _ := range ReturnJsonWEB.Name {
		tmp := ReturnJsonWEB.CloseCounter[Key]
		if Key == 0 {
			ReturnJsonWEB.CloseCounter[Key] = 0
		} else {
			ReturnJsonWEB.CloseCounter[Key] = soFarWEBClose
		}
		soFarWEBClose = soFarWEBClose + tmp

		// soFarWEBClose = soFarWEBClose + ReturnJsonWEB.CloseCounter[Key]
		// ReturnJsonWEB.CloseCounter[Key] = soFarWEBClose
	}
	// log.Println(ReturnJsonAPP)
	// log.Println(ReturnJsonWEB)
	return ReturnJsonAPP, ReturnJsonWEB
}

func DiffDate(ProjectName string) JsonResult {
	var ReturnJson JsonResult
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	// Date Create
	DateRemainCreateSQL := "SELECT count(*),CreateDate FROM (SELECT *,DATE(CreatedTime) as CreateDate FROM Issues WHERE `Key` LIKE '" + ProjectName + "-%'  ORDER by CreatedTime) as b GROUP BY CreateDate"
	CreateDateRemains, err := db.Query(DateRemainCreateSQL)
	checkerr(DateRemainCreateSQL, err)
	var DateReaminResult []JiraIssuesDateRemain
	for CreateDateRemains.Next() {
		var (
			DateRemain    JiraIssuesDateRemain
			tmpCounter    int
			tmpCreateDate string
		)
		CreateDateRemains.Scan(&tmpCounter, &tmpCreateDate)
		DateRemain.Name = tmpCreateDate
		DateRemain.CreateCounter = tmpCounter
		DateReaminResult = append(DateReaminResult, DateRemain)
	}
	CreateDateRemains.Close()

	sortutil.AscByField(DateReaminResult, "Name")
	IssuesSQL := "SELECT Id FROM Issues WHERE `Status`='Closed' AND `Key` like '" + ProjectName + "-%' "
	Issuess, err := db.Query(IssuesSQL)
	checkerr(IssuesSQL, err)
	for Issuess.Next() {
		var tmpId string
		var tmpresult JiraIssuesDateRemain
		Issuess.Scan(&tmpId)
		HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Resolved' ORDER BY Sn DESC LIMIT 1"
		ResolveDate := QuerySingle(HistorySQL)
		appendflag := 0
		if ResolveDate == "" {
			HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Closed' ORDER BY Sn DESC LIMIT 1"
			ResolveDate = QuerySingle(HistorySQL)
		}
		// log.Println("Resolve date", ResolveDate, tmpId)
		for Key, Value := range DateReaminResult {
			if Value.Name == ResolveDate {
				appendflag++
				DateReaminResult[Key].CloseCounter++
			}
		}
		// log.Println(ResolveDate, tmpId, appendflag)
		if appendflag == 0 {

			tmpresult.Name = ResolveDate
			tmpresult.CloseCounter++
			DateReaminResult = append(DateReaminResult, tmpresult)
		}
	}
	Issuess.Close()
	db.Close()

	OldResultRemain := 0
	sortutil.AscByField(DateReaminResult, "Name")
	for key, Value := range DateReaminResult {
		if key == 0 {
			ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, 0)
		} else {
			ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, OldResultRemain)
		}
		OldResultRemain = (OldResultRemain + DateReaminResult[key].CreateCounter) - DateReaminResult[key].CloseCounter
		ReturnJson.Name = append(ReturnJson.Name, Value.Name)
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, Value.CreateCounter)
		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, Value.CloseCounter)
	}
	// for key, value := range ReturnJson {
	// 	log.Println(ProjectName, key, value)
	// }
	// log.Println(ReturnJson)
	return ReturnJson
}

func TYGHDiffDate(ProjectName string) (JsonResult, JsonResult) {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	// Date Create
	// start := time.Now()

	DateRemainCreateSQL := "SELECT * from (select date(UpdatedTime) as date from Issues  WHERE `Key` LIKE 'TYGH-%'  " +
		" union " +
		" select date(CreatedTime) as date  from Issues WHERE `Key` LIKE 'TYGH-%'  ) as a GROUP by a.date "
	CreateDateRemains, err := db.Query(DateRemainCreateSQL)
	checkerr(DateRemainCreateSQL, err)
	var APPDateReaminResult []JiraIssuesDateRemain
	var WEBDateReaminResult []JiraIssuesDateRemain
	for CreateDateRemains.Next() {
		var (
			DateRemain    JiraIssuesDateRemain
			tmpCreateDate string
		)
		CreateDateRemains.Scan(&tmpCreateDate)
		DateRemain.Name = tmpCreateDate
		APPDateReaminResult = append(APPDateReaminResult, DateRemain)
		WEBDateReaminResult = append(WEBDateReaminResult, DateRemain)
	}
	CreateDateRemains.Close()

	CreateIssueSQL := "SELECT Id,Version,DATE(CreatedTime) as createdate FROM `Issues` WHERE `Key` like '" + ProjectName + "-%' ORDER BY `Issues`.`Version` ASC"
	Issuess, err := db.Query(CreateIssueSQL)
	checkerr(CreateIssueSQL, err)
	for Issuess.Next() {
		var (
			tmpId    string
			tmpVerId string
			tmpDate  string
		)
		Issuess.Scan(&tmpId, &tmpVerId, &tmpDate)
		VersionsSQL := "SELECT `Versions`.`Name` FROM `Versions` INNER join `Version` on `Versions`.`Sn` = `Version`.`Data` AND `Version`.`Enable`='1' WHERE `Version`.`Id`='" + tmpId + "'"
		// VersionsSQL := "SELECT Data FROM `Version` WHERE Id='" + tmpId + "' AND Enable='1'"
		VersionsROW, err := db.Query(VersionsSQL)
		checkerr(VersionsSQL, err)
		CreateVersionCounter := 0
		for VersionsROW.Next() {
			appappendflag := 0
			webappendflag := 0
			CreateVersionCounter++
			var VersionName string
			// var versionID string
			var tmpresult JiraIssuesDateRemain
			VersionsROW.Scan(&VersionName)
			// VersionsROW.Scan(&versionID)
			// VersionNameSQL := "SELECT Name FROM `Versions` WHERE Sn='" + versionID + "'"
			// VersionName := QuerySingle(VersionNameSQL)
			// log.Println("Create", VersionName, tmpId, tmpDate)
			if strings.Contains(VersionName, "APP") {
				for Key, Value := range APPDateReaminResult {
					if Value.Name == tmpDate {
						appappendflag++
						APPDateReaminResult[Key].CreateCounter++
					}
				}
				if appappendflag == 0 {
					tmpresult.Name = tmpDate
					tmpresult.CreateCounter++
					APPDateReaminResult = append(APPDateReaminResult, tmpresult)
				}
			} else if strings.Contains(VersionName, "WEB") {
				for Key, Value := range WEBDateReaminResult {
					if Value.Name == tmpDate {
						webappendflag++
						WEBDateReaminResult[Key].CreateCounter++
					}
				}
				if webappendflag == 0 {
					log.Println("web", tmpDate)
					tmpresult.Name = tmpDate
					tmpresult.CreateCounter++
					WEBDateReaminResult = append(WEBDateReaminResult, tmpresult)
				}
			}
		}
		VersionsROW.Close()

		// log.Println(CreateVersionCounter)
		ResolvedIssueSQL := "SELECT Status as createdate FROM `Issues` WHERE `Key` like '" + ProjectName + "-%' AND Id='" + tmpId + "' ORDER BY `Issues`.`Version` ASC"
		Status := QuerySingle(ResolvedIssueSQL)
		if Status == "Closed" {
			HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Resolved' ORDER BY Sn DESC LIMIT 1"
			ResolveDate := QuerySingle(HistorySQL)
			if ResolveDate == "" {
				HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Closed' ORDER BY Sn DESC LIMIT 1"
				ResolveDate = QuerySingle(HistorySQL)
			}

			FixVersionsSQL := "SELECT `Fixversions`.`Name` FROM `Fixversions` INNER join `Fixversion` on `Fixversions`.`Sn` = `Fixversion`.`Data` AND `Fixversion`.`Enable`='1' WHERE `Fixversion`.`Id`='" + tmpId + "'"

			// FixVersionsSQL := "SELECT Data FROM `Fixversion` WHERE Id='" + tmpId + "' AND Enable='1'"
			// log.Println("tmpid", tmpId)
			FixVersionsROW, err := db.Query(FixVersionsSQL)
			checkerr(FixVersionsSQL, err)
			for FixVersionsROW.Next() {
				var tmpresult JiraIssuesDateRemain
				appappendflag := 0
				webappendflag := 0
				// var versionID string
				var FixVersionName string
				// FixVersionsROW.Scan(&versionID)
				FixVersionsROW.Scan(&FixVersionName)
				// log.Println("tmpid", tmpId, "versionID", versionID)
				// FixVersionSQL := "SELECT Name FROM `Fixversions` WHERE Sn='" + versionID + "'"
				// FixVersionName := QuerySingle(FixVersionSQL)
				// log.Println("Close", FixVersionName, tmpId, ResolveDate, CreateVersionCounter)
				if strings.Contains(FixVersionName, "APP") && CreateVersionCounter >= 0 {
					for Key, Value := range APPDateReaminResult {
						if Value.Name == ResolveDate {
							appappendflag++
							APPDateReaminResult[Key].CloseCounter++
						}
					}
					if appappendflag == 0 {
						tmpresult.Name = ResolveDate
						tmpresult.CloseCounter++
						APPDateReaminResult = append(APPDateReaminResult, tmpresult)
					}
				} else if strings.Contains(FixVersionName, "WEB") && CreateVersionCounter >= 0 {
					for Key, Value := range WEBDateReaminResult {
						if Value.Name == ResolveDate {
							webappendflag++
							WEBDateReaminResult[Key].CloseCounter++
						}
					}
					if webappendflag == 0 {
						tmpresult.Name = ResolveDate
						tmpresult.CloseCounter++
						WEBDateReaminResult = append(WEBDateReaminResult, tmpresult)
					}
				}
				CreateVersionCounter--
			}
			FixVersionsROW.Close()
		}
	}
	Issuess.Close()
	db.Close()

	sortutil.AscByField(APPDateReaminResult, "Name")
	sortutil.AscByField(WEBDateReaminResult, "Name")

	WEBRemain := 0
	APPRemain := 0
	for key, _ := range WEBDateReaminResult {
		WEBRemain = (WEBRemain + WEBDateReaminResult[key].CreateCounter) - WEBDateReaminResult[key].CloseCounter
		WEBDateReaminResult[key].RemainCounter = WEBRemain
	}
	for key, _ := range APPDateReaminResult {
		APPRemain = (APPRemain + APPDateReaminResult[key].CreateCounter) - APPDateReaminResult[key].CloseCounter
		APPDateReaminResult[key].RemainCounter = APPRemain
	}

	var ReturnJsonAPP JsonResult
	var ReturnJsonWEB JsonResult
	for key, _ := range WEBDateReaminResult {
		ReturnJsonWEB.Name = append(ReturnJsonWEB.Name, WEBDateReaminResult[key].Name)
		ReturnJsonWEB.CreateCounter = append(ReturnJsonWEB.CreateCounter, WEBDateReaminResult[key].CreateCounter)
		ReturnJsonWEB.CloseCounter = append(ReturnJsonWEB.CloseCounter, WEBDateReaminResult[key].CloseCounter)
		ReturnJsonWEB.RemainCounter = append(ReturnJsonWEB.RemainCounter, WEBDateReaminResult[key].RemainCounter)
	}
	for key, _ := range APPDateReaminResult {
		ReturnJsonAPP.Name = append(ReturnJsonAPP.Name, APPDateReaminResult[key].Name)
		ReturnJsonAPP.CreateCounter = append(ReturnJsonAPP.CreateCounter, APPDateReaminResult[key].CreateCounter)
		ReturnJsonAPP.CloseCounter = append(ReturnJsonAPP.CloseCounter, APPDateReaminResult[key].CloseCounter)
		ReturnJsonAPP.RemainCounter = append(ReturnJsonAPP.RemainCounter, APPDateReaminResult[key].RemainCounter)
	}

	// log.Println(ReturnJsonAPP, ReturnJsonWEB)
	// elapsed := time.Since(start)
	// log.Println(elapsed)
	return ReturnJsonAPP, ReturnJsonWEB
}

func PieChart(ProjectName string, TYPE string) JsonResultPie {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	var ReturnJson JsonResult
	var ReturnJsonPie JsonResultPie
	PrioritySQL := "Select `" + TYPE + "` FROM `Issues` WHERE `Key` like '" + ProjectName + "-%' group by `" + TYPE + "` "
	// log.Println(PrioritySQL)
	PriorityRows, err := db.Query(PrioritySQL)
	checkerr(PrioritySQL, err)
	for PriorityRows.Next() {
		var (
			tmpName string
		)
		PriorityRows.Scan(&tmpName)
		// log.Println(tmpName)
		if tmpName != "" {
			ReturnJson.Name = append(ReturnJson.Name, tmpName)
			ReturnJsonPie.Name = append(ReturnJsonPie.Name, tmpName)
		} else if tmpName == "" && TYPE == "Resolution" {
			ReturnJson.Name = append(ReturnJson.Name, tmpName)
			ReturnJsonPie.Name = append(ReturnJsonPie.Name, "Unresolved")
		}
	}
	PriorityRows.Close()
	totalcounter := 0
	for _, Value := range ReturnJson.Name {
		// log.Println(key, Value)
		tmpValue := Value
		TmpCounter := 0
		SQL := "Select count(*) FROM `Issues` WHERE `Key` like '" + ProjectName + "-%' AND `" + TYPE + "` = ?  "
		// log.Println(SQL, "="+tmpValue+"=")
		rows, err := db.Query(SQL, tmpValue)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&TmpCounter)
			if err != nil {
				log.Fatal(err)
			}
			// log.Println(TmpCounter)
		}
		rows.Close()
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, TmpCounter)
		totalcounter = totalcounter + TmpCounter
	}

	for key, _ := range ReturnJson.Name {
		tmpremain := ((float64)(ReturnJson.CreateCounter[key]) * 100) / (float64)(totalcounter)
		// log.Println(tmpremain, ReturnJson.CreateCounter[key], totalcounter)
		ReturnJsonPie.RemainCounter = append(ReturnJsonPie.RemainCounter, tmpremain)
	}
	// log.Println(ReturnJsonPie)
	db.Close()
	// 	// log.Println(ReturnJson)
	return ReturnJsonPie
}

func IssueTimespent(ProjectName string) []JsonResultTable {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	RemainIssueSQL := "SELECT a.`Key`, a.`Summary`, a.`Priority`, a.DiffDate, a.`Assignee` from (SELECT *,DATEDIFF(DATE(`UpdatedTime`),DATE(CURDATE() )) AS DiffDate " +
		"FROM `Issues` WHERE `Key` LIKE '" + ProjectName + "-%' AND Status!='Closed' AND Resolution not like  '% Fix' ORDER BY UpdatedTime AND Status='Closed' )as a " +
		"ORDER by a.DiffDate"
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	var RemainIssueResult []JsonResultTable
	for RemainIssues.Next() {
		var tmpRemainIssue JsonResultTable
		RemainIssues.Scan(&tmpRemainIssue.Name, &tmpRemainIssue.Summary, &tmpRemainIssue.Priority, &tmpRemainIssue.DiffDate, &tmpRemainIssue.Assignee)
		// log.Println(tmpRemainIssue)
		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)

	}
	RemainIssues.Close()
	db.Close()
	return RemainIssueResult
}

func LastUpdateWeek() JsonResultTime {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	RemainIssueSQL := "SELECT week(FROM_UNIXTIME(UNIX_TIMESTAMP(MAX(UPDATE_TIME))))+1 as week , FROM_UNIXTIME(UNIX_TIMESTAMP(MAX(UPDATE_TIME))) as last_update , NOW() as now " +
		"FROM information_schema.tables  WHERE TABLE_SCHEMA='Jira_Data' GROUP BY TABLE_SCHEMA"
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	var returnJson JsonResultTime
	for RemainIssues.Next() {
		RemainIssues.Scan(&returnJson.Week, &returnJson.LastUpdate, &returnJson.Current)
	}
	RemainIssues.Close()
	db.Close()
	// log.Println(returnJson)
	return returnJson
}

func DueDateRemain(ProjectName string) []JsonResultTable {
	//SELECT * FROM `Issues` WHERE `DueDate`!= "0000-00-00" AND Status != 'Closed'  AND `Key` like 'BABY-%' ORDER BY `FixVersions` DESC
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	RemainIssueSQL := "SELECT `Key`, `Summary`, `Priority`, `Assignee` FROM `Issues` WHERE `DueDate`!= '0000-00-00' AND Status != 'Closed'  AND `Key` like '" + ProjectName + "-%' ORDER BY `FixVersions` DESC"
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	var RemainIssueResult []JsonResultTable
	for RemainIssues.Next() {
		var tmpRemainIssue JsonResultTable
		RemainIssues.Scan(&tmpRemainIssue.Name, &tmpRemainIssue.Summary, &tmpRemainIssue.Priority, &tmpRemainIssue.Assignee)
		// log.Println(tmpRemainIssue)
		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)
	}
	RemainIssues.Close()
	db.Close()
	return RemainIssueResult
}

func BABYIOSPieChart(ProjectName string, TYPE string) JsonResultPie {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	var ReturnJson JsonResult
	var ReturnJsonPie JsonResultPie
	PrioritySQL := "Select `" + TYPE + "` FROM `Issues` WHERE `Reporter`='Ming_Chiang' AND `Key` like '" + ProjectName + "-%' group by `" + TYPE + "` "
	// log.Println(PrioritySQL)
	PriorityRows, err := db.Query(PrioritySQL)
	checkerr(PrioritySQL, err)
	for PriorityRows.Next() {
		var (
			tmpName string
		)
		PriorityRows.Scan(&tmpName)
		// log.Println(tmpName)
		if tmpName != "" {
			ReturnJson.Name = append(ReturnJson.Name, tmpName)
			ReturnJsonPie.Name = append(ReturnJsonPie.Name, tmpName)
		} else if tmpName == "" && TYPE == "Resolution" {
			ReturnJson.Name = append(ReturnJson.Name, tmpName)
			ReturnJsonPie.Name = append(ReturnJsonPie.Name, "Unresolved")
		}
	}
	PriorityRows.Close()
	totalcounter := 0
	for _, Value := range ReturnJson.Name {
		// log.Println(key, Value)
		tmpValue := Value
		TmpCounter := 0
		SQL := "Select count(*) FROM `Issues` WHERE `Reporter`='Ming_Chiang' AND `Key` like '" + ProjectName + "-%' AND `" + TYPE + "` = ?  "
		// log.Println(SQL, "="+tmpValue+"=")
		rows, err := db.Query(SQL, tmpValue)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&TmpCounter)
			if err != nil {
				log.Fatal(err)
			}
			// log.Println(TmpCounter)
		}
		rows.Close()
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, TmpCounter)
		totalcounter = totalcounter + TmpCounter
	}

	for key, _ := range ReturnJson.Name {
		tmpremain := ((float64)(ReturnJson.CreateCounter[key]) * 100) / (float64)(totalcounter)
		// log.Println(tmpremain, ReturnJson.CreateCounter[key], totalcounter)
		ReturnJsonPie.RemainCounter = append(ReturnJsonPie.RemainCounter, tmpremain)
	}
	// log.Println(ReturnJsonPie)
	db.Close()
	// 	// log.Println(ReturnJson)
	return ReturnJsonPie
}

func BABYIOSIssueTimespent(ProjectName string) []JsonResultTable {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	RemainIssueSQL := "SELECT a.`Key`, a.`Summary`, a.`Priority`, a.DiffDate, a.`Assignee` from (SELECT *,DATEDIFF(DATE(`UpdatedTime`),DATE(CURDATE() )) AS DiffDate " +
		"FROM `Issues` WHERE `Reporter`='Ming_Chiang' AND `Key` LIKE '" + ProjectName + "-%' AND Status!='Done' AND Resolution not like  '% Fix' ORDER BY UpdatedTime AND Status='Closed' )as a " +
		"ORDER by a.DiffDate"
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	var RemainIssueResult []JsonResultTable
	for RemainIssues.Next() {
		var tmpRemainIssue JsonResultTable
		RemainIssues.Scan(&tmpRemainIssue.Name, &tmpRemainIssue.Summary, &tmpRemainIssue.Priority, &tmpRemainIssue.DiffDate, &tmpRemainIssue.Assignee)
		// log.Println(tmpRemainIssue)
		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)

	}
	RemainIssues.Close()
	db.Close()
	return RemainIssueResult
}

func BABYIOSVersionData(ProjectName string) []JiraIssuesStatusDiffVersion {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	// ProjectName := "TYGH"
	versSQL := "SELECT Id,Name from ( " +
		"SELECT f.Id,f.Name,f.Project FROM Fixversions as f LEFT join Versions as v on f.Id = v.Id UNION DISTINCT " +
		"SELECT v.Id,v.Name,v.Project FROM Fixversions as f RIGHT join Versions as v on f.Id = v.Id) as a where  Project='" + ProjectName + "'  group by a.Id " +
		"ORDER BY `a`.`Id` ASC"
	Vers, err := db.Query(versSQL)
	checkerr(versSQL, err)
	var DiffVersionResult []JiraIssuesStatusDiffVersion
	for Vers.Next() {
		var (
			IssueVersions JiraIssuesStatusDiffVersion
			tmpVerId      string
			tmpVerName    string
		)
		Vers.Scan(&tmpVerId, &tmpVerName)
		IssueVersions.Name = tmpVerName
		DiffVersionResult = append(DiffVersionResult, IssueVersions)
	}
	Vers.Close()

	// // TODO no version
	sortutil.AscByField(DiffVersionResult, "Name")
	IssuesSQL := "SELECT Id,FixVersions,Status,Version FROM Issues WHERE `Reporter`='Ming_Chiang' AND `Key` like '" + ProjectName + "-%' AND Resolution not like '% Fix'"
	Issuess, err := db.Query(IssuesSQL)
	checkerr(IssuesSQL, err)
	for Issuess.Next() {
		var (
			tmpId             string
			tmpFixVer         string
			tmpVer            string
			tmpStatus         string
			CreateVersionName []string
		)
		Issuess.Scan(&tmpId, &tmpFixVer, &tmpStatus, &tmpVer)
		// log.Println(tmpId, tmpFixVer, tmpStatus, tmpVer)
		IssueCreateCounter := 0
		if tmpVer != "0" {
			VerDataSQL := " SELECT Data FROM `Version` WHERE `Id` = '" + tmpId + "' AND Enable ='1' "
			VerDatas, err := db.Query(VerDataSQL)
			checkerr(VerDataSQL, err)
			for VerDatas.Next() {
				var tmpData string
				VerDatas.Scan(&tmpData)
				VersNameSQL := " SELECT Name FROM `Versions` WHERE `Sn` = '" + tmpData + "' AND Project ='" + ProjectName + "' "
				VerName := QuerySingle(VersNameSQL)
				// log.Println("Create", VerName)
				CreateVersionName = append(CreateVersionName, VerName)
				for key, Value := range DiffVersionResult {
					if Value.Name == VerName {
						DiffVersionResult[key].CreateCounter++
						IssueCreateCounter++
					}
				}
			}
			VerDatas.Close()
		}
		// log.Println("Create counter", IssueCreateCounter)
		if tmpStatus == "Done" {
			FixVerDataSQL := " SELECT Data FROM `Fixversion` WHERE `Id` = '" + tmpId + "' AND Enable ='1' "
			VerDatas, err := db.Query(FixVerDataSQL)
			checkerr(FixVerDataSQL, err)
			for VerDatas.Next() {
				var tmpData string
				VerDatas.Scan(&tmpData)
				VersNameSQL := " SELECT Name FROM `Fixversions` WHERE `Sn` = '" + tmpData + "' AND Project ='" + ProjectName + "' "
				VerName := QuerySingle(VersNameSQL)
				// log.Println("close", VerName)
				for key, Value := range DiffVersionResult {
					if Value.Name == VerName && IssueCreateCounter > 0 {
						DiffVersionResult[key].CloseCounter++
						IssueCreateCounter--
					}
				}
			}
		}
		// log.Println("Close counter", IssueCreateCounter, VersionRemainCounter)
	}
	Issuess.Close()
	db.Close()
	return DiffVersionResult
}

func BABYIOSDiffVersion(ProjectName string) JsonResult {
	DiffVersionResult := BABYIOSVersionData(ProjectName)
	var ReturnJson JsonResult
	VersionRemainCounter := 0
	for key, _ := range DiffVersionResult {
		if key == 0 {
			DiffVersionResult[key].RemainCounter = 0
			VersionRemainCounter = DiffVersionResult[key].CreateCounter - DiffVersionResult[key].CloseCounter
		} else {
			DiffVersionResult[key].RemainCounter = VersionRemainCounter
			VersionRemainCounter = VersionRemainCounter + DiffVersionResult[key].CreateCounter - DiffVersionResult[key].CloseCounter
		}
		ReturnJson.Name = append(ReturnJson.Name, DiffVersionResult[key].Name)
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, DiffVersionResult[key].CreateCounter)
		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, DiffVersionResult[key].CloseCounter)
		ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, DiffVersionResult[key].RemainCounter)
		// log.Println(key, Value)
	}
	// log.Println(ReturnJson)
	return ReturnJson
}

func BABYIOSDiffVersionSoFarRemain(ProjectName string) JsonResult {
	ReturnJson := BABYIOSDiffVersion(ProjectName)
	soFarClose := 0
	for Key, _ := range ReturnJson.Name {
		tmp := ReturnJson.CloseCounter[Key]
		if Key == 0 {
			ReturnJson.CloseCounter[Key] = 0
		} else {
			ReturnJson.CloseCounter[Key] = soFarClose
		}
		soFarClose = soFarClose + tmp
	}
	return ReturnJson
}

func BABYIOSDiffDate(ProjectName string) JsonResult {
	var ReturnJson JsonResult
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	// Date Create
	DateRemainCreateSQL := "SELECT count(*),CreateDate FROM (SELECT *,DATE(CreatedTime) as CreateDate FROM Issues WHERE `Reporter`='Ming_Chiang' AND `Key` LIKE '" + ProjectName + "-%'  ORDER by CreatedTime) as b GROUP BY CreateDate"
	// log.Println(DateRemainCreateSQL)
	CreateDateRemains, err := db.Query(DateRemainCreateSQL)
	checkerr(DateRemainCreateSQL, err)
	var DateReaminResult []JiraIssuesDateRemain
	for CreateDateRemains.Next() {
		var (
			DateRemain    JiraIssuesDateRemain
			tmpCounter    int
			tmpCreateDate string
		)
		CreateDateRemains.Scan(&tmpCounter, &tmpCreateDate)
		DateRemain.Name = tmpCreateDate
		DateRemain.CreateCounter = tmpCounter
		DateReaminResult = append(DateReaminResult, DateRemain)
	}
	CreateDateRemains.Close()

	sortutil.AscByField(DateReaminResult, "Name")
	IssuesSQL := "SELECT Id FROM Issues WHERE `Reporter`='Ming_Chiang' AND `Status`='Done' AND `Key` like '" + ProjectName + "-%' "
	Issuess, err := db.Query(IssuesSQL)
	checkerr(IssuesSQL, err)
	for Issuess.Next() {
		var tmpId string
		var tmpresult JiraIssuesDateRemain
		Issuess.Scan(&tmpId)
		HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Done' ORDER BY Sn DESC LIMIT 1"
		ResolveDate := QuerySingle(HistorySQL)
		appendflag := 0
		if ResolveDate == "" {
			HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Done' ORDER BY Sn DESC LIMIT 1"
			ResolveDate = QuerySingle(HistorySQL)
		}
		// log.Println("Resolve date", ResolveDate, tmpId)
		for Key, Value := range DateReaminResult {
			if Value.Name == ResolveDate {
				appendflag++
				DateReaminResult[Key].CloseCounter++
			}
		}
		// log.Println(ResolveDate, tmpId, appendflag)
		if appendflag == 0 {

			tmpresult.Name = ResolveDate
			tmpresult.CloseCounter++
			DateReaminResult = append(DateReaminResult, tmpresult)
		}
	}
	Issuess.Close()
	db.Close()

	OldResultRemain := 0
	sortutil.AscByField(DateReaminResult, "Name")
	for key, Value := range DateReaminResult {
		if key == 0 {
			ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, 0)
		} else {
			ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, OldResultRemain)
		}
		OldResultRemain = (OldResultRemain + DateReaminResult[key].CreateCounter) - DateReaminResult[key].CloseCounter
		ReturnJson.Name = append(ReturnJson.Name, Value.Name)
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, Value.CreateCounter)
		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, Value.CloseCounter)
	}
	// for key, value := range ReturnJson {
	// 	log.Println(ProjectName, key, value)
	// }
	// log.Println(ReturnJson)
	return ReturnJson
}

func BABYIOSDueDateRemain(ProjectName string) []JsonResultTable {
	//SELECT * FROM `Issues` WHERE `DueDate`!= "0000-00-00" AND Status != 'Closed'  AND `Key` like 'BABY-%' ORDER BY `FixVersions` DESC
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	RemainIssueSQL := "SELECT `Key`, `Summary`, `Priority`, `Assignee` FROM `Issues` WHERE `DueDate`!= '0000-00-00' AND `Reporter`='Ming_Chiang' AND Status != 'Done' AND `Key` like '" + ProjectName + "-%' ORDER BY `FixVersions` DESC"
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	var RemainIssueResult []JsonResultTable
	for RemainIssues.Next() {
		var tmpRemainIssue JsonResultTable
		RemainIssues.Scan(&tmpRemainIssue.Name, &tmpRemainIssue.Summary, &tmpRemainIssue.Priority, &tmpRemainIssue.Assignee)
		// log.Println(tmpRemainIssue)
		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)
	}
	RemainIssues.Close()
	db.Close()
	return RemainIssueResult
}

func BABYIssueTimespent() []JsonResultTable {
	db, _ := sql.Open("mysql", "eli:eli@/Jira_Data")

	RemainIssueSQL := "SELECT a.`Key`, a.`Summary`, a.`Priority`, a.DiffDate, a.`Assignee` from " +
		"( SELECT *,DATEDIFF(DATE(`UpdatedTime`),DATE(CURDATE() )) AS DiffDate 	FROM `Issues` WHERE  " +
		"(`Reporter`='Ming_Chiang' AND `Key` LIKE 'IOS-%' AND Status!='Done' AND Resolution not like  '% Fix' ) OR " +
		"(`Key` LIKE 'BABY-%' AND Status!='Closed' AND Resolution not like  '% Fix'  ) )as a  ORDER by a.DiffDate"

	var RemainIssueResult []JsonResultTable
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	for RemainIssues.Next() {
		var tmpRemainIssue JsonResultTable
		RemainIssues.Scan(&tmpRemainIssue.Name, &tmpRemainIssue.Summary, &tmpRemainIssue.Priority, &tmpRemainIssue.DiffDate, &tmpRemainIssue.Assignee)
		// log.Println(tmpRemainIssue)
		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)
	}
	RemainIssues.Close()

	sortutil.AscByField(RemainIssueResult, "DiffDate")
	// log.Println(RemainIssueResult)

	db.Close()
	return RemainIssueResult
}

func BABYDueDateRemain() []JsonResultTable {
	//SELECT * FROM `Issues` WHERE `DueDate`!= "0000-00-00" AND Status != 'Closed'  AND `Key` like 'BABY-%' ORDER BY `FixVersions` DESC
	db, _ := sql.Open("mysql", "eli:eli@/Jira_Data")
	RemainIssueSQL := "SELECT `Key`, `Summary`, `Priority`, `Assignee` FROM `Issues` WHERE (`DueDate`!= '0000-00-00' AND `Reporter`='Ming_Chiang' AND Status != 'Done' AND `Key` like 'IOS-%') OR (`DueDate`!= '0000-00-00' AND Status != 'Closed'  AND `Key` like 'BABY-%' )ORDER BY `FixVersions` DESC"
	var RemainIssueResult []JsonResultTable
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	for RemainIssues.Next() {
		var tmpRemainIssue JsonResultTable
		RemainIssues.Scan(&tmpRemainIssue.Name, &tmpRemainIssue.Summary, &tmpRemainIssue.Priority, &tmpRemainIssue.Assignee)
		// log.Println(tmpRemainIssue)
		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)
	}
	RemainIssues.Close()
	db.Close()
	return RemainIssueResult
}

func BABYDiffDate() (JsonResult, JsonResult) {
	JsonResult_I := BABYIOSDiffDate("IOS")
	JsonResult_A := DiffDate("BABY")

	var newandroidjson JsonResult
	var newiosjson JsonResult
	var tmpName []string

	for _, value := range JsonResult_A.Name {
		if !stringInSlice(value, tmpName) {
			tmpName = append(tmpName, value)
		}
	}

	for _, value := range JsonResult_I.Name {
		if !stringInSlice(value, tmpName) {
			tmpName = append(tmpName, value)
		}
	}

	sort.Strings(tmpName)

	newandroidjson.Name = append(newandroidjson.Name, tmpName...)
	newiosjson.Name = append(newiosjson.Name, tmpName...)

	Androidremain := 0
	Iosremain := 0

	for _, nvalue := range newandroidjson.Name {
		for key, value := range JsonResult_A.Name {
			if value == nvalue {
				Androidremain = JsonResult_A.RemainCounter[key]
			}
		}
		newandroidjson.RemainCounter = append(newandroidjson.RemainCounter, Androidremain)
	}

	for _, nvalue := range newiosjson.Name {
		for key, value := range JsonResult_I.Name {
			if value == nvalue {
				Iosremain = JsonResult_I.RemainCounter[key]
			}
		}
		newiosjson.RemainCounter = append(newiosjson.RemainCounter, Iosremain)
	}

	// for Key, value := range newandroidjson.Name {
	// 	log.Println(Key, value, newandroidjson.RemainCounter[Key])
	// 	log.Println(Key, value, newiosjson.RemainCounter[Key])

	// }

	return newandroidjson, newiosjson
}

func ProjectSummary(ProjectName string) []ProjectSummaryTable {
	db, _ := sql.Open("mysql", "eli:eli@/Jira_Data")

	RemainIssueSQL := "SELECT `Project`,`Version`,`FileName`,ROUND(`Scenario`,0),ROUND(`DataCheck`,0),ROUND(`Auto`,0),ROUND(`BDI`,0),ROUND(`Compatibility`,0),ROUND(`Security`,0),ROUND(`Other`,0),ROUND(`Battery`,0),DATE(`Date`) FROM `Projectsummary` WHERE  Project='" + ProjectName + "' order by `Date` DESC"

	var RemainIssueResult []ProjectSummaryTable
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	for RemainIssues.Next() {
		var tmpRemainIssue ProjectSummaryTable
		RemainIssues.Scan(&tmpRemainIssue.Project, &tmpRemainIssue.Version, &tmpRemainIssue.FileName, &tmpRemainIssue.Scenario, &tmpRemainIssue.DataCheck, &tmpRemainIssue.Auto, &tmpRemainIssue.BDI, &tmpRemainIssue.Compatibility, &tmpRemainIssue.Security, &tmpRemainIssue.Other, &tmpRemainIssue.Battery, &tmpRemainIssue.Date)
		// log.Println(tmpRemainIssue)
		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)
	}
	RemainIssues.Close()

	db.Close()
	return RemainIssueResult
}
